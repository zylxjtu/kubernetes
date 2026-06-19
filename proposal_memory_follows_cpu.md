# Proposal: Windows Memory Manager Should Follow the CPU Manager's NUMA Placement

> **TL;DR** — On Windows, memory placement *follows CPU affinity* (the kernel's best‑effort
> "ideal‑node" policy; there is no way to pin memory to a NUMA node). Yet the Memory Manager still
> picks NUMA nodes **independently** of the CPU Manager, so its per‑node ledger can describe a
> placement that never happens — drifting from reality and corrupting admission and topology hints.
> **Proposal: make the Windows Memory Manager derive its NUMA affinity from the CPU Manager's
> decision instead of computing its own.** This keeps the two managers in sync, models what the OS
> actually does, and structurally eliminates the CPU‑affinity "union" problem described in
> [`proposal_CPU_union.md`](./proposal_CPU_union.md).

---

## 1. Relationship to `proposal_CPU_union.md`

`proposal_CPU_union.md` addresses the CPU‑affinity **union** that can grab CPUs owned by other
containers, and recommends (**Option 1**) always enforcing the CPU Manager's set. That is correct
for the **enforcement** layer — the affinity actually sent to the runtime — but it is **silent on
the Memory Manager's bookkeeping**.

This document argues Option 1 is **necessary but not sufficient**: it makes the *enforced* affinity
honest, but the Memory Manager's own state remains an independent, unenforced, drifting model on
Windows. Fixing the **root cause** — independent NUMA selection by a manager that cannot enforce it
— makes the union *structurally impossible* rather than merely suppressed.

## 2. The gap: the Memory Manager drifts from reality on Windows

Two facts combine:

1. **Memory NUMA placement is not enforceable on Windows.** There is no equivalent of Linux's
   per‑cgroup memory‑node confinement. The kernel serves a thread's pages from the NUMA node of the
   CPU it runs on, with best‑effort fallback. So **memory follows the CPU affinity**, not the Memory
   Manager's chosen node.

2. **The Memory Manager still chooses nodes independently.** It takes the Topology Manager's hint
   and, when that hint cannot satisfy the request, **extends it on its own** toward nodes with free
   *memory* — with no knowledge of which nodes the CPU Manager chose. The CPU Manager independently
   extends toward nodes with free *CPUs*.

Because the two extend independently, their node sets can **diverge**. And because only the CPU set
is enforced and memory follows it, the Memory Manager's ledger describes a placement that **never
occurs**.

## 3. Windows memory management is best‑effort by construction

On Windows there is no way to *enforce* memory onto a NUMA node — the kernel places a thread's pages
on the node of the CPU it runs on, best‑effort. The Memory Manager can therefore only *predict*
where memory lands, never pin it, and the correct prediction is "the node of the assigned CPUs."
(Two further platform facts — the per‑node capacity baseline and the request‑accounting currency —
are also only approximate; see the [Appendix](#appendix-windows-memory-semantics).)

**This is why the proposal costs nothing.** "Memory follows CPU" sacrifices no guarantee that ever
existed on Windows — there was none. It makes the *advisory* model **consistent with the one thing
Windows does honor (CPU affinity)** instead of letting it drift independently.
Best‑effort‑but‑consistent strictly dominates best‑effort‑and‑divergent.

## 4. Worked example

Three NUMA nodes, 4 CPUs each. A Guaranteed container requests **8 exclusive CPUs + 8 GB**. Topology
policy is best‑effort, so a multi‑node allocation is admitted.

- **CPU Manager:** 8 CPUs don't fit one node → picks nodes `{0, 1}`.
- **Memory Manager:** 8 GB doesn't fit node 0's free memory → extends toward free memory on node 2 →
  records its 8 GB across `{0, 2}`.

| | CPU Manager | Memory Manager | Enforced | Actual (kernel) |
|---|---|---|---|---|
| NUMA nodes | `{0,1}` | `{0,2}` | `{0,1}` | memory near CPUs → `{0,1}` |

The ledger claims 8 GB across nodes 0 and 2, but the memory physically lands on 0 and 1 (following
the CPUs). It is wrong in two directions: node 2 looks reserved but is free; node 1 looks free but
holds this container's memory.

## 5. Why it matters

The per‑node ledger drives **admission** and **topology hints** for future pods. Drift propagates:

- A later pod wanting node‑2 memory may be **needlessly rejected** (node 2 looks occupied but is
  free) → under‑utilization.
- A later pod placed on node 1 (looks free) competes with this container's memory there → **degraded
  locality**. On Windows this doesn't OOM (the kernel spills best‑effort), but the NUMA locality the
  Memory Manager exists to provide is lost.

So on Windows the Memory Manager, as currently designed, can **make locality worse** than doing
nothing — spending admission budget and emitting hints from a model the OS contradicts. (On Linux
this never arises: memory enforcement *follows* the memory ledger, so the ledger stays truthful.)

## 6. The proposal: memory follows CPU

**On Windows, the Memory Manager should derive its NUMA affinity from the CPU Manager's allocation
rather than computing its own.** After the CPU Manager picks a container's nodes, the Memory Manager
records its memory against **those same nodes** — exactly, with no independent extension.

This makes the model match the platform:

- **Bookkeeping = reality.** Memory is recorded where the kernel will actually place it (the CPU
  nodes). No drift; admission and hints become truthful.
- **One NUMA decision, not two.** The CPU Manager is the single NUMA authority on Windows; the
  Memory Manager becomes a follower.

A consequence to accept (§8): when the CPU's nodes lack enough free memory for the request, the
Memory Manager records the full request there anyway and lets the overflow fall where the kernel
best‑effort places it, rather than re‑deciding onto a different node.

## 7. How it subsumes the union fix

If the Memory Manager's nodes are *always* the CPU Manager's nodes, the memory node set is by
construction a superset of (or equal to) the CPU set — so reconciling the two at container‑create
time is a no‑op that returns the CPU set. The harmful union case **can never be reached**. This
proposal therefore **delivers Option 1's outcome as a side effect**, and additionally makes the
Memory Manager's state truthful, which Option 1 alone does not.

| Property | Option 1 alone | This proposal |
|---|---|---|
| Enforced affinity == CPU decision | ✅ | ✅ |
| Union removed | ✅ (suppressed) | ✅ (impossible) |
| Memory ledger matches reality | ❌ | ✅ |
| Admission / hints truthful on Windows | ❌ | ✅ |

## 8. Trade‑offs

1. **One‑way coupling (Windows‑only).** The Memory Manager depends on the CPU Manager's decision and
   on the CPU Manager deciding first. The dependency direction (memory → follows → CPU) is sane and
   confined to Windows.

2. **Soft accounting on tight nodes.** When the CPU's nodes lack free memory, the overflow is
   recorded on those nodes even though the kernel spills it elsewhere; the spill node's free figure
   is left optimistic. This is *honest about where, soft about how much* — and it is the same
   best‑effort nature the platform already imposes (§3 and the [Appendix](#appendix-windows-memory-semantics)),
   not new debt. It never drives a node's free figure negative, and it stays internally consistent
   across kubelet restart and container teardown.

3. **Does the Memory Manager still add value? Yes.** Following the CPU Manager removes only the
   Memory Manager's *independent re‑selection of nodes at allocation time* — the drift source. It
   remains a full participant in the *hint/admission* phase: its per‑NUMA memory‑availability hints
   are what make the (CPU‑committed) NUMA decision **memory‑aware** in the first place, and it still
   gates admission on memory capacity and tracks per‑node memory pressure for future pods. So the
   proposal keeps the Memory Manager's genuinely useful role and drops only the harmful one. Eliding
   it entirely (Alternative B) would forfeit that memory‑awareness — the more minimal but strictly
   less capable option.

## 9. Alternatives considered

**A. Fix it in the Topology Manager.** Have the Topology Manager produce a single affinity both
managers honor, without divergent extension. Correct and platform‑neutral, but a larger change to
shared admission logic. "Memory follows CPU" is the Windows‑specific shortcut to the same end‑state,
justified because on Windows CPU affinity is the only NUMA lever that exists.

**B. Don't run a Memory Manager on Windows at all.** Since memory cannot be NUMA‑enforced, one could
elide the Memory Manager on Windows and let CPU affinity plus the kernel's ideal‑node policy provide
best‑effort locality directly. Simpler, but it is a real trade‑down, not a free simplification: it
forfeits what the Memory Manager still contributes even when it follows CPU (§8.3) — **memory‑aware
NUMA placement** (without its hints, node selection is blind to per‑node memory), memory‑capacity
admission, and observability. The more minimal option, but strictly less capable.

**C. Write the union back to CPU state** (`proposal_CPU_union.md` Option 3). Rejected there for sound
reasons (breaks the CPU‑state invariant, races with reconcile, corrupts the meaning of "allocated").
This proposal achieves consistency from the other side — fix the *memory* decision — without writing
the enforcement layer's expansion back into CPU state.

## 10. Recommendation

Adopt **Option 1** (drop the union) **and** this proposal (Windows memory follows CPU):

- Option 1 alone makes the enforced affinity honest but leaves the Memory Manager's state a fiction
  on Windows.
- Making memory **follow** the CPU Manager makes the layers consistent, makes the union structurally
  impossible, and models what the OS actually does.
- **Alternative B** (elide the Memory Manager on Windows) is the more minimal option, but it gives
  up memory‑aware NUMA placement and admission (§8.3) — pursue it only if that contribution is
  judged unnecessary, not as a free simplification.

**In one sentence:** on Windows the CPU affinity is the only NUMA lever, so the CPU Manager should be
the sole NUMA decision‑maker and the Memory Manager should follow it — not allocate independently
against a placement the OS will never honor.

---

## Appendix: Windows memory semantics

Three independent facts — none fixable in the kubelet — make the Windows Memory Manager **inherently
advisory**: even with perfect bookkeeping it cannot guarantee placement or an exact per‑node
capacity. Fact 1 is the premise of the proposal (§3); facts 2–3 reinforce that the model is
approximate regardless of bookkeeping quality.

1. **No NUMA memory enforcement.** Placement is the kernel's best‑effort ideal‑node policy: pages
   are served from the NUMA node of the CPU the thread runs on, with fallback (spill) to other nodes
   under pressure. The Memory Manager can only *predict* where memory will land, not pin it — and the
   correct prediction is "the node of the assigned CPUs."

2. **The per‑node baseline is *available* memory, not installed total.** On Windows the manager's
   per‑node capacity is seeded from the node's *available* (free) physical memory captured once at
   startup — not the installed total (as on Linux), and never refreshed. So the baseline is itself a
   stale approximation.

3. **Requests are *committed* memory; the baseline is *physical* memory.** A pod's memory request is
   charged by Windows as committed memory (backed by RAM **plus** the pagefile), while the baseline
   is available physical RAM. The manager reconciles two different accounting currencies — a node can
   read "out of free physical RAM" while the workload is still satisfiable via commit, and vice‑versa.

Together these mean the model is advisory regardless of bookkeeping quality, which is why making it
**consistent with the CPU decision** (rather than independently "accurate") is the right goal.
