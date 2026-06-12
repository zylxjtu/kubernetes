# Proposal: Fix Case 3 CPU Affinity in `computeFinalCpuSet` (Windows)

> **TL;DR** — When CPU Manager and Memory Manager disagree on NUMA placement,
> `computeFinalCpuSet()` (Windows) unions both sets ("Case 3"). That union both
> **breaks CPU isolation** (it can grab CPUs owned by other containers) and is
> **reverted by the reconcile loop** anyway, making it a no-op in steady state.
> **Recommendation: Option 1 — always use the CPU Manager's set.** It is correct,
> simple, and keeps bookkeeping aligned with enforcement. The real locality fix
> belongs upstream in Topology Manager coordination.

**Code:** `pkg/kubelet/cm/internal_container_lifecycle_windows.go`

## Contents

1. [What is Case 3?](#1-what-is-case-3)
2. [Background: how we reach Case 3](#2-background-how-we-reach-case-3)
3. [Lifecycle of CPU Manager state — who writes it, and when](#3-lifecycle-of-cpu-manager-state--who-writes-it-and-when)
4. [The two problems](#4-the-two-problems)
5. [Options](#5-options)
6. [Assessment: is Option 3 better?](#6-assessment-is-option-3-better-no)
7. [Recommendation](#7-recommendation)

---

## 1. What is Case 3?

When both CPU Manager and Memory Manager are active on Windows, `computeFinalCpuSet()`
merges their results. There are three cases:

| Case | Condition | Resolution |
|---|---|---|
| Case 1 | CPU Manager CPUs > NUMA CPUs | Use CPU Manager's set |
| Case 2 | CPU Manager CPUs are a subset of NUMA CPUs | Use CPU Manager's set |
| **Case 3** | **Partial overlap — some CPU Manager CPUs are outside NUMA CPUs** | **Union of both sets** ⚠️ |

**Case 3 happens when the CPU Manager and Memory Manager disagree on NUMA placement.**
Example — a container requests 8 CPUs and 8 GB memory on a 2-node system:

- CPU Manager can't fit all 8 CPUs on node 0 (only 4 free), so it allocates
  CPUs 0-3 (node 0) + CPUs 16-19 (node 1).
- Memory Manager picks node 0 (enough free memory).
- Partial overlap: CPUs 0-3 are on node 0, CPUs 16-19 are not.
- Case 3 unions both: CPU Manager's `{0-3, 16-19}` ∪ node 0's `{0-15}`.

The union can include CPUs that don't belong to this container — see [§4](#4-the-two-problems).

---

## 2. Background: how we reach Case 3

The Topology Manager coordinates NUMA placement **before** allocation, using hints. Both
CPU Manager and Memory Manager are "hint providers." The workflow runs **per-container**:

```
Step 1: Pod arrives → Topology Manager starts admission (per container).

Step 2: Collect hints from all providers (for THIS container)
  CPU Manager:    {NUMA 0} preferred, {NUMA 1} preferred, {NUMA 0,1} not preferred
  Memory Manager: {NUMA 0} preferred, {NUMA 0,1} not preferred

Step 3: Topology Manager merges hints (bitwise AND of all combinations)
  CPU {NUMA 0} ∩ Memory {NUMA 0}     = {NUMA 0}  ← both preferred ✓ BEST
  CPU {NUMA 1} ∩ Memory {NUMA 0}     = empty     ← rejected
  CPU {NUMA 0,1} ∩ Memory {NUMA 0,1} = {NUMA 0,1} ← neither preferred
  Winner: NUMA node 0. Result is a NUMA-node bitmask — it decides "where", not "what".

Step 4: Store decision — saves {NUMA 0} as affinity for this container.

Step 5: Call Allocate() on each provider (within the NUMA constraint)
  CPU Manager:    picks exact CPUs 0,1,2,3 from node 0, saves to state
  Memory Manager: reserves 8GB from node 0, saves to state
  (reservations in Go maps — nothing enforced yet)

Step 6: Pod admitted, container creation starts.

Step 7: PreCreateContainer runs (PLATFORM-SPECIFIC)
  Reads reservations, writes affinity into ContainerConfig (the CRI request).

Step 8: CRI CreateContainer → OS-level enforcement applied.
```

**Key distinctions:**
- **Hints (Step 2):** flexible, multiple options — "node 0 OR node 1 works for me."
- **Allocation (Step 5):** final, exact resources — "CPUs 0,1,2,3" + "8 GB from node 0."
- **Enforcement (Step 8):** the OS actually restricts the container process.

---

## 3. Lifecycle of CPU Manager state — who writes it, and when

The CPU Manager's state (`assignments` + `defaultCPUSet`) is the **single source of truth**.
It is written **only** in the allocation path, **never** at container-create time. This is
the crux of the bug in [§4](#4-the-two-problems), so it is worth making explicit:

| When | Trigger → call chain | Effect on state |
|---|---|---|
| **Pod admission** | Topology Manager scope `Admit` → `provider.Allocate(pod, container)` (`topologymanager/scope.go`) → `manager.Allocate` (`cpu_manager.go:269`) → `policy.Allocate` → `s.SetCPUSet(...)` (`policy_static.go:615`) | **WRITE** — decides the container's exclusive CPUs; removes them from `defaultCPUSet`. The only place "what CPUs" is decided. |
| **Container create** | `kuberuntime` → `PreCreateContainer` (`internal_container_lifecycle_windows.go:62`) → `computeFinalCpuSet` → CRI `AffinityCpus` | **READ only** (`GetCPUAffinity`). Computes the mask for the CRI request. **Does not touch state.** ← Case 3 expansion lives here. |
| **Container start** | `PreStartContainer` (`internal_container_lifecycle.go:41`) → `cpuManager.AddContainer` | Associates the runtime containerID with the already-allocated pod/container. No CPU decision. |
| **Reconcile (periodic)** | timer (`CPUManagerReconcilePeriod`) → `reconcileState()` (`cpu_manager.go:448`) → reads `GetCPUSetOrDefault` (`:512`) → `updateContainerCPUSet` (`:523`) | **READ + enforce.** Pushes state to running containers via CRI. Reverts any mask that disagrees with state. |
| **Container removal** | `RemoveContainer` → `policy.RemoveContainer` → `SetDefaultCPUSet` (`policy_static.go:667-669`) | **WRITE** — returns the CPUs to the shared pool. |

**Takeaway:** allocation (state) and the Windows Case 3 expansion (CRI mask) live in
**different layers that never sync**. The reconcile loop trusts state, so it overwrites
the expansion.

---

## 4. The two problems

### 4.1 Isolation violation

Original Case 3 code:

```go
// Case 3: union everything
return cpuManagerAffinityCPUSet.Union(numaNodeAffinityCPUSet)
```

This unions **all** NUMA-node CPUs, including ones the CPU Manager exclusively assigned to
**other** containers. Those containers are pinned to the same CPUs — so both can run on the
same "exclusive" CPUs. **The isolation guarantee is broken.**

```
System: node 0 = CPUs 0-15, node 1 = CPUs 16-31

Container A (running):  exclusive CPUs {4,5,6,7,8,9,10,11}
Container B (creating): CPU Manager {0,1,16,17}; Memory Manager node 0 {0-15}

Case 3 result for B: {0,1,16,17} ∪ {0-15} = {0..15, 16, 17}
  → B's mask now includes CPUs 4-11, exclusively owned by A. Both run on the same CPUs.
```

### 4.2 Reconciliation reverts the expansion

Even with the overlap filtered out, the expansion does not survive. As shown in
[§3](#3-lifecycle-of-cpu-manager-state--who-writes-it-and-when), `PreCreateContainer`
expands the mask but **never updates state**, while `reconcileState()` reads state as the
source of truth and re-applies it:

```
Admission:       policy.Allocate → state = {0,1,16,17}        (4 CPUs, authoritative)
CreateContainer: computeFinalCpuSet → CRI mask {0,1,2,3,12-17} (10 CPUs; state untouched)
... container runs with the expanded mask ...
Reconcile tick:  reads state {0,1,16,17} → updateContainerCPUSet
                 → mask overwritten back to 4 CPUs            (NUMA-local CPUs lost)
```

This makes the Case 3 union **effectively a no-op in steady state**.

---

## 5. Options

### Option 1 — Always use CPU Manager CPUs *(recommended, simplest)*

Treat Case 3 like Cases 1 and 2 — always return the CPU Manager's set:

```go
if cpuManagerAffinityCPUSet.Len() > 0 {
    return cpuManagerAffinityCPUSet
}
```

| Pros | Cons |
|---|---|
| Simplest change, no extra parameters | Some CPUs may sit on a non-optimal NUMA node for memory... |
| Bookkeeping and enforcement always aligned | ...but that is the situation the CPU Manager already created — a |
| No isolation violation; no reconcile mismatch | performance characteristic, not a correctness issue |
| Container gets exactly the CPUs it was allocated | |

### Option 2 — Filter NUMA CPUs to exclude already-allocated ones

```go
availableNumaCPUs := numaNodeAffinityCPUSet.Difference(allAllocatedSet)
return cpuManagerAffinityCPUSet.Union(availableNumaCPUs)
```

| Pros | Cons |
|---|---|
| Preserves NUMA locality (adds free NUMA-local CPUs) | Still has the bookkeeping/enforcement mismatch |
| No isolation violation with other containers | Reconcile reverts the expansion anyway |
| | Container temporarily gets more CPUs than requested |
| | More complex; must thread `allAllocatedCPUs` through |

### Option 3 — Option 2 + update CPU Manager state

Same as Option 2 but also update CPU Manager state to reflect the expanded set and teach
reconciliation about the Windows-specific expansion.

| Pros | Cons |
|---|---|
| Preserves NUMA locality permanently | Most complex change |
| No reconciliation revert | State no longer distinguishes "allocated" vs "expanded for NUMA" |
| Bookkeeping matches enforcement | Other containers see fewer CPUs in the shared pool |
| | Crosses component boundaries (`PreCreateContainer` writing CPU Manager state) |

---

## 6. Assessment: is Option 3 better? (No)

Option 3 buys a marginal, partly illusory benefit at a serious architectural cost.

**Why the benefit is weak:**
- It does **not** fix Case 3's root cause. The container's exclusive CPUs are already split
  across NUMA nodes; Option 3 adds *more* local CPUs but never moves the remote ones, so the
  locality split remains.
- The CPUs it adds aren't free: they're either allocated to other containers (the isolation
  bug — only avoided by also doing Option 2's filtering) or in the shared pool. Adding
  shared-pool CPUs to a guaranteed container lets it run on contended CPUs *and* shrinks the
  pool other pods depend on.
- On Windows the affinity mask is the only NUMA lever (no kernel memory pinning), so the
  locality story is softer than on Linux anyway.

**Why "crosses component boundaries" is the disqualifying con:**

The intended data flow is one-directional (see [§3](#3-lifecycle-of-cpu-manager-state--who-writes-it-and-when)) —
Topology Manager hints → CPU Manager `Allocate()` **writes** state → `reconcileState()` /
`PreCreateContainer` **read** state and translate it into the CRI request. Option 3 inverts
this by having the enforcement layer write back into the allocator's state. That is broken
on four levels:

1. **Breaks a state invariant.** The static policy keeps `assignments` (exclusive CPUs) and
   `defaultCPUSet` (shared pool) complementary; `GetAvailableCPUs = defaultCPUSet − reserved`.
   A bare `SetCPUSet` with the expanded set leaves the NUMA CPUs in *both* the assignment and
   the shared pool. Doing it correctly means also editing `defaultCPUSet` and reversing it on
   removal — i.e. reimplementing policy allocation outside the policy.
2. **Concurrency.** State is lock-guarded and the reconcile loop runs continuously; a second
   async writer races with it and with the `lastUpdateState` diff logic.
3. **Corrupts the meaning of "allocated."** State becomes a mix of real allocations and
   Windows NUMA expansions, so metrics/accounting/debugging that trust it read a fiction.
4. **Platform leakage.** `cpu_manager.go` / `policy_static.go` are cross-platform; Option 3
   requires Windows-specific logic in this shared core, whereas today all Windows logic is
   isolated in `internal_container_lifecycle_windows.go`.

**Conclusion:** Option 3 cannot be implemented correctly without dissolving the clean
separation that makes CPU Manager state trustworthy — a high price for a benefit that does
not even close the NUMA-span gap.

---

## 7. Recommendation

Go with **Option 1**:

- Case 3 is an edge case (CPU Manager and Memory Manager disagree on NUMA placement).
- The NUMA-locality benefit of the union is lost to reconciliation anyway.
- Keeping bookkeeping aligned with enforcement matters more than a temporary NUMA optimization.
- The CPU Manager's allocation is authoritative and final — trust it.
