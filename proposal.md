# Proposal: Fix Case 3 CPU Affinity in computeFinalCpuSet (Windows)

## What is Case 3?

When both CPU Manager and Memory Manager are active on Windows, `computeFinalCpuSet()`
in `internal_container_lifecycle_windows.go` merges their results. There are three cases:

| Case | Condition | Resolution |
|---|---|---|
| Case 1 | CPU Manager CPUs > NUMA CPUs | Use CPU Manager's set |
| Case 2 | CPU Manager CPUs are a subset of NUMA CPUs | Use CPU Manager's set |
| **Case 3** | **Partial overlap — some CPU Manager CPUs are outside NUMA CPUs** | **Union of both sets** |

**Case 3 happens when the CPU Manager and Memory Manager disagree on NUMA placement.**

For example, a container requests 8 CPUs and 8GB memory:
- CPU Manager can't fit all 8 CPUs on NUMA node 0 (only 4 free), so it allocates
  CPUs 0-3 (node 0) + CPUs 16-19 (node 1)
- Memory Manager picks NUMA node 0 (has enough free memory)
- Partial overlap: CPUs 0-3 are on NUMA node 0, but CPUs 16-19 are not
- Case 3 unions both sets: CPU Manager's {0-3, 16-19} ∪ NUMA node 0's {0-15}

This is the problematic case — the union can include CPUs that don't belong to this container.

## Problem

### Isolation Violation in Original Code

Original code (before our commit) in `pkg/kubelet/cm/internal_container_lifecycle_windows.go`:

```go
// Case 3: union everything
return cpuManagerAffinityCPUSet.Union(numaNodeAffinityCPUSet)
```

This blindly unions all NUMA node CPUs, including ones allocated to other containers.
The container's affinity mask would include CPUs that the CPU Manager exclusively assigned
to other containers. Those other containers also have affinity masks pinning them to those
CPUs — so both containers can run on the same "exclusive" CPUs. **The isolation guarantee
is broken.**

### Example

```
System: NUMA node 0 has CPUs 0-15, NUMA node 1 has CPUs 16-31

Container A (already running): CPU Manager assigned exclusive CPUs {4,5,6,7,8,9,10,11}
Container B (being created):   CPU Manager assigned CPUs {0,1,16,17}
                                Memory Manager selected NUMA node 0 (CPUs 0-15)

Current Case 3 result:
  {0,1,16,17} ∪ {0-15} = {0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17}

Container B's affinity mask now includes CPUs 4-11, which are exclusively
allocated to Container A. Both containers can run on the same CPUs.
```

## Additional Problem: Reconciliation Reverts Case 3

Even if we fix the overlap issue, there is a deeper problem. The CPU Manager's
`reconcileState()` loop (`cpu_manager.go:485-496`) periodically reads the CPU set
from its bookkeeping state and re-applies it to running containers via
`updateContainerCPUSet()`.

`computeFinalCpuSet()` expands the CPU set at `PreCreateContainer` time but
**never updates the CPU Manager state**. So reconciliation reverts the expansion:

```
Step 5:  CPU Manager books {0,1,16,17}          (4 CPUs in state)
Step 7:  computeFinalCpuSet expands to           {0,1,2,3,12-17} (10 CPUs)
Step 8:  Container created with 10-CPU affinity
  ...
Later:   reconcileState() reads state          → {0,1,16,17} (4 CPUs)
         updateContainerCPUSet({0,1,16,17})    → overwrites affinity
         Container affinity SHRINKS to 4 CPUs  → NUMA-local CPUs lost
```

This makes the Case 3 union logic **effectively a no-op in steady state**.

## Options

### Option 1: Always use CPU Manager CPUs (simplest)

Treat Case 3 the same as Cases 1 and 2 — always return CPU Manager's set.

```go
if cpuManagerAffinityCPUSet.Len() > 0 {
    return cpuManagerAffinityCPUSet
}
```

**Pros:**
- Simplest change, no extra parameters needed
- Bookkeeping and enforcement always aligned
- No isolation violation
- No reconciliation mismatch
- Container gets exactly the CPUs it was allocated

**Cons:**
- Some CPUs may be on a non-optimal NUMA node for memory access
- But this is already the situation the CPU Manager created, and it's a
  performance characteristic, not a correctness issue

### Option 2: Filter NUMA CPUs to exclude already-allocated ones

```go
allAllocatedSet := sets.New[int](allAllocatedCPUs.List()...)
availableNumaCPUs := numaNodeAffinityCPUSet.Difference(allAllocatedSet)
return cpuManagerAffinityCPUSet.Union(availableNumaCPUs)
```

**Pros:**
- Preserves NUMA locality benefit (adds free NUMA-local CPUs)
- No isolation violation with other containers

**Cons:**
- Still has bookkeeping/enforcement mismatch
- Reconciliation loop reverts the expansion anyway
- Container gets more CPUs than requested (temporarily)
- More complex, requires passing allAllocatedCPUs through

### Option 3: Option 2 + update CPU Manager state

Same as Option 2 but also update the CPU Manager state to reflect the expanded set,
and update reconciliation to be aware of the Windows-specific expansion.

**Pros:**
- Preserves NUMA locality benefit permanently
- No reconciliation revert
- Bookkeeping matches enforcement

**Cons:**
- Most complex change
- CPU Manager state no longer reflects what was actually "allocated" vs "expanded for NUMA"
- Other containers see fewer available CPUs in the shared pool
- Crosses component boundaries (PreCreateContainer modifying CPU Manager state)

## Recommendation

**Option 1** is the safest and simplest choice:
- Case 3 is an edge case (CPU Manager and Memory Manager disagree on NUMA placement)
- The NUMA locality benefit from the union is lost to reconciliation anyway
- Keeping bookkeeping and enforcement aligned is more important than temporary NUMA optimization
- The CPU Manager's allocation is authoritative and final — trust it

The real fix for Case 3 scenarios is upstream: better Topology Manager coordination
so CPU Manager and Memory Manager agree on NUMA placement in the first place.

---

## Background: How We Get to Case 3

### Topology Manager Workflow (Per-Container)

The Topology Manager coordinates NUMA placement **before** allocation, using a hint-based system.
Both CPU Manager and Memory Manager are "hint providers." The workflow runs **per-container**, not per-pod.

**Full workflow for a Guaranteed pod (e.g., 4 CPUs + 8GB memory) on a 2-NUMA node system:**

```
Step 1: Pod arrives → Topology Manager starts admission
        Iterates over each container independently.

Step 2: Collect hints from all providers (for THIS container)
  CPU Manager returns multiple options (enumerates all valid NUMA combinations):
    - {NUMA 0} preferred    ← node 0 has enough free CPUs
    - {NUMA 1} preferred    ← node 1 has enough free CPUs
    - {NUMA 0,1} not preferred ← works but spans nodes
  Memory Manager returns:
    - {NUMA 0} preferred    ← node 0 has enough free memory
    - {NUMA 0,1} not preferred ← works but spans nodes

Step 3: Topology Manager merges hints (bitwise AND of all combinations)
    CPU {NUMA 0} ∩ Memory {NUMA 0}     = {NUMA 0} ← both preferred ✓ BEST
    CPU {NUMA 1} ∩ Memory {NUMA 0}     = empty     ← rejected
    CPU {NUMA 1} ∩ Memory {NUMA 0,1}   = {NUMA 1}  ← memory not preferred
    CPU {NUMA 0,1} ∩ Memory {NUMA 0,1} = {NUMA 0,1} ← neither preferred
  Winner: NUMA node 0
  Result is a NUMA node bitmask, not CPUs. It decides "where", not "what".

Step 4: Store decision — saves {NUMA 0} as affinity for this container

Step 5: Call Allocate() on each provider (within NUMA constraint)
    CPU Manager:    picks exact CPUs 0,1,2,3 from NUMA node 0, saves to state
    Memory Manager: reserves 8GB from NUMA node 0, saves to state
    These are reservations in Go maps — nothing enforced yet.

Step 6: Pod admitted, container creation starts

Step 7: PreCreateContainer runs (PLATFORM-SPECIFIC — see Part 2 and Part 3)
    Reads reservations, writes into ContainerConfig (the CRI request)

Step 8: CRI CreateContainer sent to container runtime (PLATFORM-SPECIFIC)
    OS-level enforcement applied. Container is now restricted.
```

**Key distinctions:**
- **Hints (Step 2):** Flexible, multiple options — "node 0 OR node 1 works for me"
- **Allocation (Step 5):** Final, exact resources — "CPUs 0,1,2,3" and "8GB from node 0"
- **Enforcement (Step 8):** OS actually restricts the container process

### Bookkeeping vs Enforcement Mismatch (Case 3 Bug)

In Case 3 (partial overlap), `computeFinalCpuSet()` expands the CPU set beyond what the
CPU Manager booked. But the CPU Manager state is **never updated** to reflect this expansion.
The reconciliation loop then **reverts** the expanded set back to the original bookkeeping.

**Proof from code:**

1. `PreCreateContainer()` (`internal_container_lifecycle_windows.go:70`):
   Computes `finalCPUSet` (expanded union), writes to `containerConfig.Windows.Resources.AffinityCpus`.
   **Does NOT call `state.SetCPUSet()` to update CPU Manager state.**

2. `reconcileState()` (`cpu_manager.go:485-496`):
   Reads `cset` from `m.state.GetCPUSetOrDefault()` — the **original** bookkeeping.
   Calls `updateContainerCPUSet(ctx, containerID, cset)` with the bookkeeping set.

**Timeline showing the revert:**
```
Step 5:  CPU Manager books {0,1,16,17}          (4 CPUs in state)
Step 7:  computeFinalCpuSet expands to           {0,1,2,3,12-17} (10 CPUs)
Step 8:  Container created with 10-CPU affinity
  ...
Later:   reconcileState() reads state          → {0,1,16,17} (4 CPUs)
         updateContainerCPUSet({0,1,16,17})    → overwrites affinity
         Container affinity SHRINKS to 4 CPUs  → NUMA-local CPUs lost
```
