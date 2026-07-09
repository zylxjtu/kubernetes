# Kubernetes CPU/Memory Affinity Knowledge

## Feature Overview: WindowsCPUAndMemoryAffinity (KEP-4885)

- **Feature Gate:** `WindowsCPUAndMemoryAffinity`
- **KEP:** [4885](https://github.com/kubernetes/enhancements/tree/master/keps/sig-windows/4885-windows-cpu-and-memory-affinity)
- **Status:** Alpha since v1.32 (default: disabled)
- **Owner:** @jsturtevant
- **Note:** Code comment in `kube_features.go:1086` incorrectly references `kep.k8s.io/4888` — the real KEP is 4885

### What the gate controls on Windows (all-or-nothing)

On Windows this is a **single** gate that governs all three managers together — there is
**no** separate per-manager gate (no `Windows*CPUManager` / `*MemoryManager` / `*Topology`
gate exists). The Windows container manager wires it as one switch
(`container_manager_windows.go:138-186`):

| `WindowsCPUAndMemoryAffinity` | Topology Manager | CPU Manager | Memory Manager |
|---|---|---|---|
| **off** (default) | fake | fake | fake |
| **on** | real | real | real |

- When **off**, all three are fake managers (`NewFakeManager`), so `--topology-manager-policy`,
  `--cpu-manager-policy`, etc. are **inert**, and `PreCreateContainer` early-returns without
  setting any affinity (`internal_container_lifecycle_windows.go:35`). You cannot enable just
  one (e.g. a real CPU Manager with a fake Topology Manager) — it is the whole trio or nothing.
- When **on**, the configured policy / scope / options are passed straight through to the same
  cross-platform managers used on Linux.

(On Linux these managers are GA and not feature-gated; the fake-by-default pattern is
Windows-specific.)

---

## Part 1: Platform-Agnostic Workflow (Same on Linux and Windows)

### High-Level Flow

```
Hint → Topology Manager Decision → Reservation/Bookkeeping → Container Creation + Enforcement
```

All steps up to container creation are platform-agnostic. The same Go code runs on both Linux and Windows.

### NUMA Hardware Basics

```
Physical Socket (CPU package, e.g. an Intel Xeon chip)
  └── NUMA Node (usually 1:1 with socket; some AMD EPYC have 2+ per socket)
       ├── Local Memory (fast, ~100ns)
       └── CPUs (physical cores → logical processors via Hyperthreading/SMT)
```

NUMA = Non-Uniform Memory Access. Each node has its own CPUs and locally-attached RAM.
Accessing local memory is fast (~100ns), remote memory is 1.5-3x slower.

#### Socket ↔ NUMA Node Relationship

There is **no fixed ratio** between physical sockets and NUMA nodes — every cardinality
is possible, depending on the CPU microarchitecture and firmware (BIOS) configuration.
Never infer one from the other; always query the actual topology at runtime.

| Mapping | When it happens | Example |
|---|---|---|
| **1:1** | Classic one-memory-controller-per-socket (textbook NUMA) | Older 2-socket Xeon → 2 nodes |
| **1:N** (nodes > sockets) | Chiplet CPUs or sub-NUMA clustering — **common today** | AMD EPYC with NPS2/NPS4, Intel Sub-NUMA Clustering (SNC) → 1 socket exposes 2–4 nodes |
| **N:1** (sockets > nodes) | Firmware interleaves memory across sockets, collapsing node count | Multi-socket box configured as fewer NUMA nodes |

- A **socket** (package) is physical — the chip in the motherboard slot.
- A **NUMA node** is topological — a set of CPUs with equal-latency access to one memory bank.
- They are correlated but live on different axes, which is why all of 1:1, 1:N, and N:1 occur.
- Modern reality skews toward **1:N**: a 2-socket EPYC at NPS4 presents **2 sockets but 8 NUMA nodes**.

**Windows adds a third axis — processor groups** (≤64 logical CPUs each), which also do
*not* map 1:1 to NUMA nodes: a node >64 CPUs is split across multiple groups, and small
nodes can share a group. So on Windows you juggle **socket ⟷ NUMA node ⟷ processor group**,
all three pairwise many-to-many. This is why `winstats` resolves topology from the OS at
runtime rather than assuming any ratio (see Part 3).

### Configuring the Managers

The Topology, CPU, and Memory managers are **three independent components**, each with its
own policy setting — there is no single "affinity policy." The Topology Manager is the
*coordinator* (it sits above the other two and aligns them via hints); CPU and Memory
managers are the *hint providers / allocators*.

Each has its own kubelet flag and its own `KubeletConfiguration` field:

| Setting | CLI flag | Config field | Allowed values | Default |
|---|---|---|---|---|
| Topology policy | `--topology-manager-policy` | `topologyManagerPolicy` | `none`, `best-effort`, `restricted`, `single-numa-node` | `none` |
| Topology scope | `--topology-manager-scope` | `topologyManagerScope` | `container`, `pod` | `container` |
| Topology options | `--topology-manager-policy-options` | `topologyManagerPolicyOptions` | key=value (feature-gated) | — |
| CPU policy | `--cpu-manager-policy` | `cpuManagerPolicy` | `none`, `static` | `none` |
| CPU options | `--cpu-manager-policy-options` | `cpuManagerPolicyOptions` | key=value (feature-gated) | — |
| Memory policy | `--memory-manager-policy` | `memoryManagerPolicy` | `None`, `Static` | `None` |

> **Casing gotcha (all platforms, not Windows-specific):** CPU and Topology values are
> lowercase (`static`, `best-effort`); Memory Manager values are **capitalized** (`None`,
> `Static`). The flag help text reflects this. Wrong casing → kubelet fails validation.

**Example — `KubeletConfiguration`** (the recommended way; many of the flags above are deprecated):

```yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
# each manager set independently:
topologyManagerPolicy: single-numa-node      # coordinator
topologyManagerScope: pod
cpuManagerPolicy: static                      # exclusive CPUs for Guaranteed pods
memoryManagerPolicy: Static                   # note the capital S
# static CPU policy needs a non-zero reservation:
reservedSystemCPUs: "0"
```

> **On Windows:** these settings only take effect when the `WindowsCPUAndMemoryAffinity`
> gate is enabled — otherwise all three managers are fakes and every setting above is inert
> (see [What the gate controls on Windows](#what-the-gate-controls-on-windows-all-or-nothing)).

### Topology Manager Policies

The Topology Manager has four policies that determine what happens when NUMA alignment
can't be achieved. Configured via `--topology-manager-policy`. Default is `none`.
All four are available on both Linux and Windows — though on Windows they only take effect
when the `WindowsCPUAndMemoryAffinity` gate is enabled (otherwise the Topology Manager is a
no-op fake; see [Feature Overview](#what-the-gate-controls-on-windows-all-or-nothing)).

| Policy | No preferred hint found | Pod admitted? | Use case |
|---|---|---|---|
| `none` (default) | No topology checks at all | **Yes** — always | Topology awareness not needed |
| `best-effort` | Uses best available non-preferred hint | **Yes** — always | Try to align, but don't block workloads |
| `restricted` | Non-preferred hint found but rejected | **No** — TopologyAffinityError | Must have NUMA alignment, reject if impossible |
| `single-numa-node` | Pre-filters to single-node hints only | **No** — TopologyAffinityError | All resources must fit on one NUMA node |

> **`none` ≠ managers off.** With topology `none`, the CPU and Memory managers still fully
> run — the none scope calls each provider's `Allocate()` per container (`scope.go:139`), so
> exclusive CPUs and memory blocks are allocated and enforced as normal. What's missing is
> *coordination*: no hints are gathered, each manager allocates with `NUMANodeAffinity = nil`
> (unconstrained), and the pod is never rejected for topology reasons. So `none` + `static`
> CPU + `Static` memory is valid and common — you just opt out of co-locating CPUs and memory
> on the same NUMA node, which makes cross-NUMA placement (and Windows Case 3) **most** likely.
> Each manager is still individually topology-aware for its *own* resource; only inter-manager
> alignment is skipped.

#### `none` vs `best-effort` (Windows)

These are the two policies that **never reject a pod** — the admission outcome is identical.
The difference is purely whether the Topology Manager *coordinates* the CPU and Memory managers
before they allocate. On Windows this matters more than on Linux, because memory locality is
achieved *indirectly* (there is no `cpuset.mems`): the Memory Manager runs the **BestEffort**
policy and memory follows the CPU NUMA decision, so a coordinated hint is the only thing nudging
CPUs and memory onto the same node.

| Aspect | `none` | `best-effort` |
|---|---|---|
| Hint collection | **Skipped** — `admitPolicyNone` calls each provider's `Allocate()` directly (`scope.go:139`) | **Performed** — gathers hints from CPU/Memory/device managers and merges them (bitwise AND) |
| NUMA affinity passed to `Allocate()` | `nil` (each manager allocates independently) | The merged best hint (a specific node when one satisfies all providers) |
| Stored topology affinity | None | Best hint stored per container |
| Pod ever rejected for topology? | No | No (falls back to the best non-preferred hint and still admits) |
| CPU ↔ memory co-location on Windows | Not coordinated → cross-NUMA placement most likely → drives Case 3 in `computeFinalCpuSet()` | Coordinated → CPU and memory prefer the same node → cross-NUMA less likely |
| Topology admission metrics (`admission_requests_total`, `admission_duration_ms`) | Not recorded on the none path | Recorded per admission |

> **Both admit; only `best-effort` coordinates.** Think of `none` as "run the managers, don't
> talk to each other" and `best-effort` as "try to align, but never block." Neither rejects —
> that behavior starts at `restricted` / `single-numa-node`. This is also why the metrics e2e
> tests use `none`: admission goes straight to `Allocate` (no hint-merge step in between and no
> rejection), so the pinning counters increment on the pure allocate path.

**What happens when a pod is admitted with a non-preferred hint (best-effort)?**

The Topology Manager stores a fallback hint (typically all NUMA nodes). The CPU Manager
and Memory Manager allocate from wherever they can, without strong NUMA alignment.
This is the scenario that leads to Case 3 (partial overlap) in Windows `computeFinalCpuSet()`:

```
best-effort policy + no preferred hint
  → pod admitted with weak/no NUMA alignment
  → CPU Manager picks CPUs from node 0 (only free ones there)
  → Memory Manager picks node 1 (only free memory there)
  → Result on BOTH OSes: cross-NUMA memory access (sub-optimal performance —
    every load/store from node-0 CPUs hits node-1 memory over the inter-NUMA
    interconnect). The user opted into this trade-off by choosing best-effort
    over restricted/single-numa-node.

  Where Linux and Windows differ is in how the two managers' decisions are
  expressed to the kernel:
  → Linux: cpuset.cpus and cpuset.mems are two independent cgroup files, each
    enforced separately by the kernel. No code has to reconcile them — the
    kernel just runs the process on the chosen CPUs and serves memory from
    the chosen nodes (with the cross-NUMA penalty).
  → Windows: there is only one enforcement primitive in this path — the
    job-object affinity (AffinityCpus). Both intents have to fit into that
    single field, which is what computeFinalCpuSet() does (Case 3 widens the
    affinity to the union so threads can also run on CPUs near the memory).
```

> **Note:** "No conflict on Linux" refers to the **code/representation** layer,
> not to NUMA performance. The cross-NUMA performance penalty is the same on
> both OSes — only the implementation mechanism differs.

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
    Reads reservations from CPU/Memory manager state, writes them into the
    ContainerConfig (the CRI request).

    Linux:   pass-through. The reservations map directly to two independent
             cgroup paths — cpuset.cpus and cpuset.mems — that the kernel
             enforces separately. No reconciliation step exists.
    Windows: reconcile via computeFinalCpuSet(). The two managers' decisions
             must be expressed through a single field (AffinityCpus) since
             Windows has no separate "memory affinity" primitive. The 3-case
             logic (Case 1/2/3 — see Part 3) merges allocatedCPUs and
             allNumaNodeCPUs into one final mask.

Step 8: CRI CreateContainer sent to container runtime (PLATFORM-SPECIFIC)
    OS-level enforcement applied. Container is now restricted.
```

**Key distinctions:**
- **Hints (Step 2):** Flexible, multiple options — "node 0 OR node 1 works for me"
- **Allocation (Step 5):** Final, exact resources — "CPUs 0,1,2,3" and "8GB from node 0"
- **Reconcile (Step 7, Windows only):** Merges CPU and memory decisions into a single AffinityCpus field; Linux skips this because two cgroup files express the two intents independently
- **Enforcement (Step 8):** OS actually restricts the container process

**Per-container allocation call chain:**
```
Topology Manager: Admit(pod)
├─ for each container in pod:
│   ├─ collectHints(pod, &container)        ← hints for THIS container
│   ├─ mergeHints → bestHint                ← NUMA decision for THIS container
│   ├─ setTopologyHints(podUID, container.Name, bestHint)
│   └─ allocateAlignedResources(pod, &container)
│       ├─ cpuManager.Allocate(pod, container)
│       │   └─ state.SetCPUSet(podUID, container.Name, cpuset)
│       └─ memoryManager.Allocate(pod, container)
│           └─ state.SetMemoryBlocks(podUID, container.Name, blocks)
```

State is stored in two-level maps keyed by `podUID → containerName`:
```
CPU Manager:    map[podUID]map[containerName]cpuset.CPUSet
Memory Manager: map[podUID]map[containerName][]Block
```

### CPU Manager Policies

There are only **two** policies. Configured via kubelet flag `--cpu-manager-policy`.

**`none` (default):**
- `Allocate()` is a no-op — does nothing
- All containers share all CPUs equally via normal OS scheduling
- No exclusive pinning, no topology awareness
- Use when CPU pinning isn't needed

**`static`:**
- Grants exclusive CPUs to containers that qualify; everyone else shares the remaining pool
- Use when CPU isolation is needed for latency-sensitive workloads
- Requires `--cpu-manager-policy=static` and non-zero CPU reservation (`systemreserved.cpu + kubereserved.cpu > 0`)

A container gets **exclusive CPUs** under `static` policy only if ALL conditions are met:

| Condition | Why |
|---|---|
| Pod QoS is **Guaranteed** (requests == limits) | Only pods with guaranteed resources deserve pinning |
| CPU request is a **whole integer** (e.g., `4`, not `3500m`) | You can't exclusively pin half a CPU |

If either condition fails, the container runs in the **shared pool** — same as `none` policy.

| QoS Class | CPU Request | Exclusive CPUs? | CPU Pool |
|---|---|---|---|
| Guaranteed | Integral (e.g., 2) | YES | Exclusive allocation |
| Guaranteed | Non-integral (e.g., 2.5) | No | Shared pool |
| Burstable | Any | No | Shared pool |
| BestEffort | Any | No | Shared pool |

**How `static` picks CPUs:**
```
1. Check if container qualifies (Guaranteed + integer CPU)
   → No: shared pool, done
   → Yes: continue
2. Get NUMA affinity from Topology Manager → "prefer NUMA node 0"
3. Find available CPUs aligned with that NUMA node
4. Pick CPUs using topology-aware packing (default: packed onto same core/socket first)
5. Remove those CPUs from shared pool → other containers can no longer use them
6. Save to state: map[podUID][containerName] = {2,3,4,5}
```

**Static policy tuning options** (set via `--cpu-manager-policy-options`):

| Option | Status | Feature gate | What it does |
|---|---|---|---|
| `full-pcpus-only` | Stable (GA) | none | Only allocate full physical cores (request must be a multiple of the SMT level) |
| `distribute-cpus-across-numa` | Beta | `CPUManagerPolicyBetaOptions` | Spread CPUs evenly across NUMA nodes instead of packing |
| `distribute-cpus-across-cores` | Alpha | `CPUManagerPolicyAlphaOptions` | Spread across cores instead of packing SMT siblings first |
| `align-by-socket` | Alpha | `CPUManagerPolicyAlphaOptions` | Align at socket boundary, not NUMA |
| `strict-cpu-reservation` | Stable (GA) | none | Remove reserved CPUs from the shared pool entirely |
| `prefer-align-cpus-by-uncorecache` | Stable (GA) | none | Best-effort align a container's CPUs to an uncore/LLC (L3) cache boundary |

> Alpha options require the `CPUManagerPolicyAlphaOptions` feature gate; Beta options require `CPUManagerPolicyBetaOptions`. Stable options need no extra gate. All options are only valid with `--cpu-manager-policy=static`.

**Applicability on Windows nodes**

Windows CPU pinning is gated by `WindowsCPUAndMemoryAffinity` and enforced via Windows CPU group affinity through the CRI. Because the Windows CRI does not expose SMT-sibling or L3/uncore cache topology, several of the Linux-oriented tuning options have no effect on Windows even though they parse successfully.

| Option | Windows | Notes |
|---|---|---|
| `static` policy (base) | Supported | Exclusive CPU affinity set via Windows CPU group affinity |
| `strict-cpu-reservation` | Supported | Reserved-CPU semantics are topology-independent; covered by Windows node e2e |
| `full-pcpus-only` | Not supported | No SMT/HT sibling information via Windows CRI |
| `prefer-align-cpus-by-uncorecache` | Not supported | No L3/uncore cache topology via Windows CRI |
| `distribute-cpus-across-numa` | Not validated | Depends on multi-NUMA topology; not exercised on Windows |
| `align-by-socket` | Not validated | Depends on socket topology exposure on Windows |
| `distribute-cpus-across-cores` | Not applicable | Requires physical-core/HT sibling data unavailable on Windows |

> Note: the `physical_cpu` and `uncore_cache` alignment metric series still initialize to `0` on Windows (the static policy zero-initializes them at startup), even though the corresponding alignment behavior is Linux-only.

### CPU Manager Shared Pool

The CPU Manager's static policy divides all CPUs on a node into pools:

```
All CPUs on the node (e.g., CPUs 0-15)
 │
 ├── RESERVED (static, never changes at runtime)
 │     Sized by: ceil(--system-reserved cpu + --kube-reserved cpu)
 │     e.g., CPUs {0, 1} reserved for kubelet + OS
 │     These stay in the shared pool but can NEVER be exclusively assigned
 │
 ├── SHARED (= "default CPU set", shrinks/grows dynamically)
 │     Initially: ALL CPUs
 │     Used by: BestEffort, Burstable, and non-integer Guaranteed pods
 │     All these containers share these CPUs via normal OS scheduling
 │
 ├── ASSIGNABLE (= SHARED - RESERVED)
 │     The CPUs that CAN be exclusively allocated
 │
 └── EXCLUSIVE ALLOCATIONS (carved out of ASSIGNABLE)
       Pinned 1:1 to specific Guaranteed containers (integer CPU requests only)
```

**Example:**
```
Node: 16 CPUs, kube-reserved=2 CPUs

Initial:  RESERVED={0,1}  SHARED={0..15}  ASSIGNABLE={2..15}

Pod A (Guaranteed, cpu:4):  assigned {2,3,4,5}   SHARED={0,1,6..15}
Pod B (Guaranteed, cpu:2):  assigned {6,7}        SHARED={0,1,8..15}
Pod C (BestEffort, cpu:500m): runs in SHARED      uses {0,1,8..15}
Pod A deleted:              {2,3,4,5} returned     SHARED={0..5,8..15}
```

#### `reservedSystemCPUs` (explicitly reserved cores)

Instead of sizing the reserved set implicitly from `--system-reserved`/`--kube-reserved`, you can pin a specific set of CPU IDs with `--reserved-cpus` (kubelet flag) / `reservedSystemCPUs` (kubelet config). This is what the Windows node e2e tests use (e.g., `reservedSystemCPUs: "0"`).

Key properties:
- **Explicit CPU IDs**, not a quantity: e.g., `reservedSystemCPUs: "0"` reserves CPU 0 specifically, whereas `--kube-reserved=cpu=1` just reserves *some* CPU.
- **Never exclusively assignable**: reserved CPUs are removed from the assignable pool, so a Guaranteed integer-CPU container is never pinned to them.
- **Counted against node Allocatable**: `allocatable = capacity − reservedSystemCPUs`. Tests that size a request against usable CPUs should use Allocatable (which already excludes the reserved set).
- **Required for the static policy**: when `--cpu-manager-policy=static`, the kubelet needs a non-zero reservation (via `reservedSystemCPUs` or `kube-reserved`/`system-reserved` cpu > 0), otherwise it refuses to start.

Interaction with the shared pool (governed by `strict-cpu-reservation`):

| `strict-cpu-reservation` | Reserved CPUs in shared pool? | Who can run there |
|---|---|---|
| `false` (default) | Yes | Reserved CPUs stay in the shared pool; BestEffort/Burstable and non-integer Guaranteed containers may still run on them (just never exclusively pinned) |
| `true` | No | Reserved CPUs are removed from the shared pool entirely and dedicated to system/kubelet |

Example (`reservedSystemCPUs: "0"`, 12 logical CPUs):
```
capacity   = 12
reserved   = {0}                 # explicit CPU ID
allocatable = 11                 # capacity − reserved
assignable = {1..11}             # exclusive allocations come from here only
CPU 0      = shared pool only    # never pinned exclusively (strict-cpu-reservation=false)
```

---

## Part 2: Linux vs Windows Comparison

### Workflow Step Comparison

| Workflow Step | Linux | Windows | Same? |
|---|---|---|---|
| Hint collection | Same hint providers, same TopologyHint struct | Same | Yes |
| Topology Manager merges hints | Same bitwise-AND logic | Same | Yes |
| CPU Manager.Allocate() | Same static policy, same state maps | Same | Yes |
| Memory Manager.Allocate() | None or Static policies | **BestEffort only** | **No** |
| PreCreateContainer() | Direct pass-through, no merging | **`computeFinalCpuSet` merging logic** | **No** |
| CRI CreateContainer | `cpuset.cpus` + `cpuset.mems` (two independent fields) | **`GROUP_AFFINITY` only (single field)** | **No** |
| OS enforcement | cgroups | **Job Objects** | **No** |

### Why the Divergence Exists

**Linux has two independent enforcement mechanisms:**
- `cpuset.cpus` — restricts which CPUs the container can run on (from CPU Manager)
- `cpuset.mems` — restricts which NUMA nodes the container can allocate memory from (from Memory Manager)
- These are separate cgroup controllers. No merging needed — `PreCreateContainer` just sets both fields directly.

**Windows has only one mechanism:**
- `GROUP_AFFINITY` — restricts which CPUs the container can run on
- There is no equivalent of `cpuset.mems`. No way to directly pin memory to a NUMA node.
- Memory locality is achieved **indirectly**: by pinning CPUs to a NUMA node, the Windows memory allocator naturally prefers that node's local memory.
- This means CPU selection must serve double duty (CPU isolation + memory locality), requiring `computeFinalCpuSet()` to reconcile both managers' decisions through a single mechanism.

### Consequences of the Divergence

1. **Memory Manager policy (BestEffort on Windows)** — weaker NUMA guarantees during reservation, so CPU Manager and Memory Manager are more likely to disagree on NUMA placement
2. **On Linux**, the Topology Manager's stricter policies (Restricted, SingleNUMANode) can reject the pod entirely if alignment fails
3. **On Windows**, mismatches are tolerated and handled at merge time in `computeFinalCpuSet()`

### Detailed Comparison Table

| Aspect | Linux | Windows |
|---|---|---|
| CPU Manager policies | `none`, `static` | `none`, `static` |
| Allocation logic | Same (topology-aware packing) | Same (topology-aware packing) |
| Shared pool management | Same | Same |
| Enforcement mechanism | cgroups (`cpuset.cpus`) | Job Objects (`GROUP_AFFINITY`) |
| How CPUs are communicated | cgroup cpuset writes | CRI `WindowsCpuGroupAffinity` |
| Memory Manager policies | `None`, `Static` | `BestEffort` only |
| Memory enforcement | cgroups (memory controller) | Job Object `JOB_OBJECT_LIMIT_JOB_MEMORY` |
| NUMA memory pinning | Direct via cgroups `cpuset.mems` | Indirect (pin CPUs to NUMA node, OS prefers local memory) |
| Feature gate needed | No (GA) | Yes (`WindowsCPUAndMemoryAffinity`, Alpha v1.32) |

---

## Part 3: Windows-Specific Details

### Windows Hardware Concepts

**Processor Groups:**
Windows originally used a single 64-bit bitmask (`KAFFINITY`) for CPU affinity — limiting systems to 64 CPUs. When servers grew beyond 64 cores, Microsoft added **Processor Groups** (Server 2008 R2) rather than rewriting every API. Each group holds up to 64 logical processors. A 128-CPU system gets 2 groups.

Starting with Windows 11 / Server 2022, processes can span all groups by default without explicit group management.

**GROUP_AFFINITY** — the fundamental structure for identifying CPUs on Windows:
```c
struct GROUP_AFFINITY {
    uint64 Mask;    // bitmask: which CPUs within this group (bits 0-63)
    uint16 Group;   // which processor group (0, 1, 2, ...)
}
```
A flat "CPU ID" maps to this as: `group = cpuID / 64`, `bit = cpuID % 64`, `mask |= (1 << bit)`.

A NUMA node must fit within one processor group (pre-Build 20348) or can span multiple groups (Build 20348+).

**Job Objects** — the Windows equivalent of Linux cgroups for containers:
- **CPU affinity** via `SetInformationJobObject(JobObjectGroupInformationEx, GROUP_AFFINITY[])` — pins processes to specific CPUs
- **CPU rate** via hard cap (percentage of CPU time, 1-10000 range)
- **Memory limit** via `JOB_OBJECT_LIMIT_JOB_MEMORY`

Container runtimes (containerd -> hcsshim -> HCS) create a job object per container and apply these controls.

**Key Windows APIs for topology discovery:**
- `GetNumaNodeProcessorMaskEx(node)` — returns GROUP_AFFINITY of CPUs on that node
- `GetNumaAvailableMemoryNodeEx(node)` — free memory on that node
- `GetLogicalProcessorInformationEx()` — full topology (cores, sockets, NUMA, caches)

### Windows PreCreateContainer + Enforcement Flow

Steps 7-8 from the common workflow, Windows-specific:

```
Step 7: PreCreateContainer (Windows)
  cpuManager.GetCPUAffinity(podUID, containerName) → exact CPUs
  memoryManager.GetMemoryNUMANodes(pod, container)  → NUMA nodes
  GetCPUsforNUMANode() for each NUMA node           → NUMA CPUs
  computeFinalCpuSet(allocatedCPUs, numaNodeCPUs, allAllocatedCPUs) → final CPU set
  groupMasks(finalCPUSet)                           → GROUP_AFFINITY structs
  Set containerConfig.Windows.Resources.AffinityCpus

Step 8: CRI CreateContainer (Windows)
  kubelet → containerd → hcsshim → HCS
  Job Object created with GROUP_AFFINITY applied
  SetInformationJobObject(JobObjectGroupInformationEx, GROUP_AFFINITY)
  hcsshim PR that wires CPU affinity to JobObject limits:
    https://github.com/microsoft/hcsshim/pull/2603 (commit 7710aa52)
```

### Does an affinity update restart the container? (No)

Affinity changes are applied **in-place** — the container keeps running, same PID. Both
platforms route through the CRI `UpdateContainerResources` verb, which mutates a *running*
container's resource constraints:

- **Linux** (`cpu_manager_others.go`): sets `LinuxContainerResources.CpusetCpus` → rewrites
  the running cgroup `cpuset.cpus`.
- **Windows** (`cpu_manager_windows.go`): sets `WindowsContainerResources.AffinityCpus` →
  updates the live Job Object group affinity (`SetInformationJobObject`). No restart.

| Manager | Updates a running container? | Restart? |
|---|---|---|
| **CPU Manager** | Yes — the periodic `reconcileState()` re-pushes the cpuset via `UpdateContainerResources` whenever it changes (e.g. as other pods come/go and reshape the shared pool) | **No** — in-place |
| **Memory Manager** | No — it has no reconcile loop; NUMA affinity is decided at admission and applied once at creation, never updated mid-life | **No** — never changes after start |

Caveats:
- This is the kubelet-managed **affinity** path. Changing a pod's resource **requests/limits**
  in the spec is a *different* mechanism (In-Place Pod Vertical Scaling, governed by each
  container's `resizePolicy`) and *can* restart a container — but that is not the affinity reconcile.
- The reconcile loop is change-gated (only updates when the cpuset differs from
  `lastUpdateState`) and, on Windows, a no-op when `WindowsCPUAndMemoryAffinity` is off.

### CPU + Memory Conflict Resolution (`computeFinalCpuSet`)

When both CPU Manager and Memory Manager provide results (`internal_container_lifecycle_windows.go`):

| Scenario | Resolution | Why |
|---|---|---|
| CPU Manager CPUs > NUMA CPUs | Use CPU Manager's set | CPU count guarantee takes priority over memory locality |
| CPU Manager CPUs subset of NUMA CPUs | Use CPU Manager's set | CPU Manager's tighter set is more specific; locality already satisfied |
| Partial overlap | Union of both sets (excluding CPUs allocated to other containers) | Preserves CPU allocation + adds NUMA-local CPUs for memory access |
| Only one manager active | Use that manager's set | No conflict to resolve |

**When does partial overlap (Case 3) happen?**
When the CPU Manager can't fit all CPUs on the NUMA node that the Memory Manager selected
(e.g., not enough free CPUs due to other containers). The CPU Manager takes CPUs from another
node, creating a mismatch that `computeFinalCpuSet()` reconciles.

### Kubernetes-to-Windows Mapping

```
1. Topology Discovery
   kubelet calls Windows APIs → builds cadvisor-compatible topology map

2. CPU/Memory Allocation (platform-agnostic, same as Linux)
   CPU Manager (static policy) → assigns flat CPU IDs [0, 1, 2, 3]
   Memory Manager (best-effort) → selects NUMA node(s)
   Topology Manager → ensures CPUs and memory are co-located

3. Translation to Windows Primitives
   CPUs [0,1,2,3] → GROUP_AFFINITY{Mask: 0xF, Group: 0}
   (winstats.CpusToGroupAffinity: group = cpu/64, bit = cpu%64)

4. CRI Call to Container Runtime
   WindowsContainerResources{
     affinity_cpus: [{cpu_mask: 0xF, cpu_group: 0}],
     memory_limit_in_bytes: 8589934592
   }

5. Runtime Enforcement
   containerd → hcsshim → Job Object
   SetInformationJobObject(JobObjectGroupInformationEx, GROUP_AFFINITY)
```

| Kubernetes Concept | Windows Primitive | Implementation |
|---|---|---|
| CPU Request/Limit | Job Object CPU Limits | cpu_count, cpu_maximum in WindowsContainerResources |
| CPU Affinity | GROUP_AFFINITY struct | WindowsCpuGroupAffinity (cpu_mask + cpu_group) |
| NUMA Node Memory | GetNumaAvailableMemoryNodeEx() | MemoryManager.GetMemoryNUMANodes() |
| NUMA→CPU Mapping | GetNumaNodeProcessorMaskEx() | winstats.GetCPUsforNUMANode() |
| Processor Groups | Up to 64 CPUs per group | CPU ID = (group * 64) + bit_position |
| System Topology | GetLogicalProcessorInformationEx() | winstats.processorInfo() |
| Container CPU Config | Job Object SetInformation() | Container runtime UpdateContainerResources() |

### Windows CPU Enforcement: Two Independent Mechanisms

Windows has **two independent** CPU control mechanisms, both applied through Job Objects:

**1. CPU Affinity (pinning) — `AffinityCpus`**
- Restricts **which CPUs** the container can run on
- Mechanism: `GROUP_AFFINITY` bitmask → `SetInformationJobObject(JobObjectGroupInformationEx)`
- Set by: CPU Manager + Memory Manager via `PreCreateContainer()` hook
- Requires `WindowsCPUAndMemoryAffinity` feature gate
- Source: `pkg/kubelet/cm/internal_container_lifecycle_windows.go`

**2. CPU Rate Limiting — `CpuMaximum`**
- Caps the container's **total CPU usage** as a percentage of system CPU
- Mechanism: Job object CPU rate control, range 1-10000 (0.01% to 100%)
- Formula: `cpuMaximum = 10 * cpuLimit.MilliValue() / totalCpuCount`
- Always available, no feature gate needed
- Source: `pkg/kubelet/kuberuntime/kuberuntime_container_windows.go`

**How they interact:**
- `AffinityCpus` and `CpuMaximum` are **independent** — both can be applied to the same container
- `CpuCount` and `CpuMaximum` are **mutually exclusive** (both are rate-limiting; `CpuCount` takes precedence)
- Affinity says "only run on these CPUs"; rate limit says "use at most this much CPU time"

**Example:** Container with `cpu: 4` on a 16-CPU system with static CPU Manager:
```
CPU Rate Limiting (from pod spec limits):
  cpuMaximum = 10 * 4000 / 16 = 2500  (25% of total system CPU)

CPU Affinity (from CPU Manager, if feature gate enabled):
  AffinityCpus = [{Group: 0, Mask: 0b1111}]  (pinned to CPUs 0-3)

Both applied to the same Job Object independently.
```

**When each gets set during container creation:**
```
1. generateWindowsContainerResources()
   → Sets CpuMaximum from pod spec limits (rate limiting)
   → Sets MemoryLimitInBytes

2. PreCreateContainer() hook
   → Reads CPU Manager state (exact CPUs)
   → Reads Memory Manager state (NUMA nodes → CPUs)
   → computeFinalCpuSet() merges them
   → Sets AffinityCpus on the same ContainerConfig

3. CRI CreateContainer() sends config to containerd
   → Both rate limit AND affinity applied to Job Object

4. Post-creation reconciliation (periodic, ticks every --cpu-manager-reconcile-period, default 10s)
   → cpu_manager.go reconcileState() compares current cpusets to last-applied
   → skips containers whose cpuset hasn't changed since last tick (no CRI call)
   → for changed containers: cpu_manager_windows.go updateContainerCPUSet()
   → CRI UpdateContainerResources() with new AffinityCpus only
   → disabled entirely when CPU manager policy is `none`
```

### Windows Memory: Committed Memory vs Working Set

The CRI `memory_limit_in_bytes` maps to `JOB_OBJECT_LIMIT_JOB_MEMORY`, which limits **committed memory** (not working set).

| | Committed Memory | Working Set |
|---|---|---|
| What it is | Virtual memory allocated and backed by page file or physical RAM. The OS has *promised* this memory is available. | Physical RAM pages currently resident — what's actually in RAM right now. |
| Includes | All allocated virtual memory (heap, stack, mapped files), whether paged in or paged out to disk | Only pages in physical RAM. Excludes anything paged out to disk. |
| Can exceed physical RAM? | Yes — backed by page file on disk | No — bounded by physical RAM |
| Fluctuates? | Grows as app allocates, rarely shrinks | Constantly changes as OS pages in/out |

**Job object memory controls:**
- `JOB_OBJECT_LIMIT_JOB_MEMORY` → limits **committed bytes** (what containers use)
- `JOB_OBJECT_LIMIT_WORKINGSET` → limits **working set** min/max (not used by containers)

**Eviction monitoring:**
Kubelet on Windows monitors **working set** (`PrivateWorkingSet`) for eviction decisions, not committed bytes. See `pkg/kubelet/winstats/perfcounter_nodestats_windows.go`. This is analogous to Linux where eviction watches RSS.

### Memory Eviction
- `pkg/kubelet/eviction/defaults_windows.go` — Default threshold: **500Mi** (vs 100Mi on Linux)
- Uses Windows perf counters (CommitLimitPages, CommitTotalPages) instead of cgroup notifications

### Stats Collection
- `pkg/kubelet/winstats/perfcounter_nodestats_windows.go` — Tracks private working set, committed bytes, page file, NUMA-aware memory

### Concrete Example

A 2-socket server with 128 cores (256 threads with HT):
```
NUMA Node 0: CPUs 0-127   → Group 0 (CPUs 0-63) + Group 1 (CPUs 64-127)
NUMA Node 1: CPUs 128-255 → Group 2 (CPUs 128-191) + Group 3 (CPUs 192-255)
```

A pod requesting `cpu: 4, memory: 8Gi`:
- CPU Manager allocates CPUs [0, 1, 2, 3] on NUMA Node 0
- Translated to: GROUP_AFFINITY{Mask: 0xF, Group: 0}
- Memory naturally allocated from Node 0's local RAM (fast path)

---

## Part 4: Key Source Files

### Platform-Agnostic (shared by Linux and Windows)
- `pkg/kubelet/cm/topologymanager/scope_container.go` — Per-container hint collection and allocation
- `pkg/kubelet/cm/topologymanager/policy.go` — Hint merging (bitwise AND)
- `pkg/kubelet/cm/cpumanager/cpu_manager.go` — CPU Manager interface and allocation
- `pkg/kubelet/cm/cpumanager/policy_static.go` — Static policy, topology-aware CPU packing, hint generation
- `pkg/kubelet/cm/cpumanager/state/state.go` — CPU Manager state (ContainerCPUAssignments)
- `pkg/kubelet/cm/memorymanager/memory_manager.go` — Memory Manager interface and allocation
- `pkg/kubelet/cm/memorymanager/state/state.go` — Memory Manager state (ContainerMemoryAssignments)

### Windows-Specific
- `pkg/kubelet/cm/internal_container_lifecycle_windows.go` — PreCreateContainer(), computeFinalCpuSet(), groupMasks()
- `pkg/kubelet/cm/cpumanager/cpu_manager_windows.go` — Converts CPU sets to WindowsCpuGroupAffinity
- `pkg/kubelet/winstats/cpu_topology.go` — Windows topology discovery APIs
- `pkg/kubelet/cm/memorymanager/policy_best_effort.go` — Windows-only BestEffort memory policy
- `pkg/kubelet/kuberuntime/kuberuntime_container_windows.go` — CPU limit to job object conversion
- `pkg/kubelet/eviction/defaults_windows.go` — Windows eviction defaults
- `pkg/kubelet/winstats/perfcounter_nodestats_windows.go` — Windows stats collection

### CRI API
- `staging/src/k8s.io/cri-api/pkg/apis/runtime/v1/api.proto` — WindowsContainerResources, WindowsCpuGroupAffinity

---

## Part 5: Tests

### Unit Tests (Windows CPU/Memory Affinity)
- `pkg/kubelet/cm/internal_container_lifecycle_windows_test.go` — Tests computeFinalCpuSet merging logic and GROUP_AFFINITY conversion
- `pkg/kubelet/cm/cpumanager/cpu_manager_test.go` — CPU Manager allocation with WindowsCPUAndMemoryAffinity feature gate

### E2E Tests (General Windows, not affinity-specific)
- `test/e2e/windows/cpu_limits.go` — CPU limit enforcement (decimal and millicpu)
- `test/e2e/windows/memory_limits.go` — Memory limits, allocatable memory, eviction

### Test Coverage vs Workflow

The unit tests cover Steps 7-8 (Windows-specific merging + GROUP_AFFINITY conversion).
Steps 1-6 are platform-agnostic and tested in Topology Manager / CPU Manager / Memory Manager test suites.

```
Workflow Step                             Test Coverage
──────────────────────────────────────    ─────────────────────────────────────
Step 1: Pod arrives                       (platform-agnostic, tested elsewhere)
Step 2: Hints collected                   (platform-agnostic, tested elsewhere)
Step 3: Topology Manager merges hints     (platform-agnostic, tested elsewhere)
Step 4: Store decision                    (platform-agnostic, tested elsewhere)
Step 5: Allocate                          Simulated by test inputs
Step 6: Container creation starts         (platform-agnostic, tested elsewhere)
Step 7: PreCreateContainer                ✅ TestComputeFinalCpuSet
Step 8: Set GROUP_AFFINITY                ✅ TestGroupMasks
```

### Test Gap
No end-to-end test exercises the full chain from pod admission through Topology Manager
hint merging to actual container config on Windows.

---

## Other Windows Feature Gates

| Feature Gate | Status | Purpose |
|---|---|---|
| WindowsCPUAndMemoryAffinity | Alpha v1.32 | CPU/Memory affinity, Topology Manager |
| WindowsGracefulNodeShutdown | Alpha v1.32, Beta v1.34 | Graceful node shutdown |
| WindowsHostNetwork | Deprecated v1.33 | Host network namespace |

---

## Part 6: The Persisted State Checkpoint

The CPU and Memory managers each keep an on-disk **checkpoint** file in the kubelet root
directory:

- `cpu_manager_state` (CPU manager)
- `memory_manager_state` (Memory manager)

Default location: `/var/lib/kubelet/` on Linux, `C:\var\lib\kubelet\` (or the configured
`--root-dir`) on Windows. Each is a JSON file recording:

- The **policy name** active when it was written (`None`, `Static`, `BestEffort`, …).
- The **machine state** — per-NUMA-node bookkeeping plus the **assignments** already handed
  out to running containers (which CPUs / NUMA nodes / how much memory each container was
  pinned to), keyed by `podUID → containerName`.

### Why it exists

The CPU/Memory managers make *stateful, exclusive* allocation decisions, and the checkpoint
solves two problems across a kubelet restart:

1. **Survive restarts without disturbing running pods.** When the kubelet restarts (upgrade,
   crash, config reload), the containers it manages keep running. The manager must remember
   "container X already owns CPUs {2,3} on NUMA node 0" so that after restart it rebuilds its
   in-memory bookkeeping from disk instead of starting blind and double-allocating the same
   resource to a different container.

2. **Detect an incompatible policy change.** The assignments in the file were computed under
   one policy's rules. If the kubelet restarts with a *different* policy, those old
   assignments may be unsafe. On startup the manager compares the checkpoint's stored policy
   against the currently configured policy, and on mismatch **refuses to start**:

   > configured policy "BestEffort" differs from state checkpoint policy "None" — please
   > drain this node and delete the memory manager checkpoint file `memory_manager_state`
   > before restarting Kubelet.

   Rather than silently discarding data (and risking overlapping exclusive allocations), it
   fails fast. This is a deliberate safety fail-stop, not a bug.

### Effects on end users

- **Changing the CPU/memory manager policy is not a hot config change.** An admin who edits
  `memoryManagerPolicy` / `cpuManagerPolicy` and simply bounces the kubelet will find it
  **fails to start** with the mismatch error, and the node goes `NotReady`.
- **Documented safe procedure:**
  1. `kubectl drain <node>` (evict pods holding the old allocations),
  2. delete the checkpoint file (`cpu_manager_state` / `memory_manager_state`),
  3. restart the kubelet,
  4. `kubectl uncordon <node>`.

  Draining first matters: deleting the file *without* draining loses the record of what
  running containers were pinned to, so after restart the manager could reassign that
  CPU/memory to new pods, causing overlapping exclusive allocations and NUMA misplacement for
  the still-running workloads.
- **Corruption / manual edits** also make the kubelet refuse to start; recovery is again
  "drain + delete + restart".
- **Normal restarts are transparent.** As long as the policy is unchanged, users notice
  nothing — the checkpoint is what makes an in-place kubelet upgrade preserve the exact
  CPU/NUMA pinning of already-running Guaranteed pods.

### Relevance to the Windows node e2e metrics tests

The Windows CPU / memory / topology manager metrics tests intentionally flip policies
(`None` ⇄ `BestEffort`, `none` ⇄ `static`) in their `BeforeEach`/`AfterEach`. That is exactly
the mismatch the checkpoint guards against, so those tests call
`updateWindowsKubeletConfigClearState` (which deletes `cpu_manager_state` and
`memory_manager_state` before restarting) instead of a plain config update — otherwise the
kubelet would refuse to start on the policy change.