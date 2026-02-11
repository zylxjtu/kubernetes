# GracefulNodeShutdown Feature Analysis

## Feature Gate Status

Both feature gates are **still in Beta** — confirmed at `pkg/features/kube_features.go:1300-1308`:

- **`GracefulNodeShutdown`**: Alpha in 1.20, Beta since 1.21
- **`GracefulNodeShutdownBasedOnPodPriority`**: Alpha in 1.23, Beta since 1.24

Neither has been promoted to GA.

## Key GA Blockers

The KEP tracking issue ([kubernetes/enhancements#2000](https://github.com/kubernetes/enhancements/issues/2000)) lists several blockers.

### 1. DaemonSet controller disagrees with shutdown manager ([#122912](https://github.com/kubernetes/kubernetes/issues/122912))

This is the **core blocker**. During graceful shutdown:

1. The shutdown manager sets the node to `NotReady` and starts killing pods
2. The `Admit()` handler at `pkg/kubelet/nodeshutdown/nodeshutdown_manager_linux.go:115-126` **rejects ALL new pods** — no exceptions
3. The DaemonSet controller sees its pod terminated and tries to reschedule it back onto the same node (DaemonSets are supposed to run on every eligible node)
4. Kubelet rejects the new DaemonSet pod, leaving it stuck in **Pending**

The code has **zero awareness of pod ownership**. Looking at `pkg/kubelet/nodeshutdown/nodeshutdown_manager.go:129-219`, the `killPods` function groups pods purely by priority — no special handling for static pods, DaemonSet pods, or any owner references.

### 2. Static pods / system-critical pods not properly handled ([#124448](https://github.com/kubernetes/kubernetes/issues/124448))

The shutdown manager fails to update Pod status for system-critical pods. Static pods (etcd, kube-apiserver in kubeadm clusters) go through the same kill path as regular pods. After shutdown, their status can be left in an inconsistent state.

### 3. DaemonSet pods stuck in `Completed` state after reboot ([#117018](https://github.com/kubernetes/kubernetes/issues/117018))

After graceful shutdown + reboot, DaemonSet pods (like calico-node) can get stuck in `Completed` state and never recover, breaking CNI networking on the node.

### 4. Other notable blockers

- **Pods stuck in Error/Completed phases never cleaned up** ([#122122](https://github.com/kubernetes/kubernetes/issues/122122)) — deployments can stall for ~10 minutes waiting for progress
- **Endpoints not updated for terminating pods** ([#116965](https://github.com/kubernetes/kubernetes/issues/116965)) — traffic continues to be sent to pods being shut down
- **Volume teardown not awaited** ([#115148](https://github.com/kubernetes/kubernetes/issues/115148)) — PV-backed pods take 6+ minutes to restart after shutdown

## Root Cause in the Code

The fundamental issue is that the shutdown path at `pkg/kubelet/nodeshutdown/nodeshutdown_manager_linux.go:313-346` treats every active pod identically:

```go
func (m *managerImpl) processShutdownEvent() error {
    activePods := m.getPods()  // ALL active pods, no filtering
    // ...
    return m.podManager.killPods(activePods)  // kill them all by priority
}
```

And the admission handler is a blanket reject with no exceptions:

```go
func (m *managerImpl) Admit(attrs *lifecycle.PodAdmitAttributes) lifecycle.PodAdmitResult {
    if m.ShutdownStatus() != nil {
        return lifecycle.PodAdmitResult{Admit: false, ...}  // rejects everything
    }
}
```

There is no logic anywhere in `pkg/kubelet/nodeshutdown/` to check `pod.OwnerReferences` for DaemonSet ownership, static pod annotations, or system namespace membership.

## Implementation Overview

### Key Files

| Component | File Path |
|-----------|-----------|
| Feature Gate Definition | `pkg/features/kube_features.go:343, 347, 1300-1308` |
| Linux Manager | `pkg/kubelet/nodeshutdown/nodeshutdown_manager_linux.go` |
| Core Manager (kill logic) | `pkg/kubelet/nodeshutdown/nodeshutdown_manager.go` |
| Configuration | `pkg/kubelet/apis/config/types.go:449-470` |
| E2E Tests | `test/e2e_node/node_shutdown_linux_test.go` |

### Shutdown Flow

1. Kubelet registers a systemd inhibit lock via D-Bus
2. When systemd signals shutdown, `processShutdownEvent()` is called
3. All active pods are fetched via `GetActivePods()`
4. Pods are grouped by priority using `groupByPriority()`
5. Each group is killed concurrently, then volumes are unmounted (best-effort)
6. The inhibit lock is released, allowing the system to proceed with shutdown

### Additional Known Limitation

- Requires **systemd version 245+** — tests skip on older versions (see [#107043](https://github.com/kubernetes/kubernetes/issues/107043))

## Summary

The feature has been stuck in beta since **1.21 (April 2021)** — over 4 years. The core problem is the DaemonSet controller and shutdown manager working at cross-purposes: the shutdown manager kills DaemonSet pods and rejects their replacements, while the DaemonSet controller keeps trying to reschedule them. Fixing this likely requires the shutdown manager to either exempt certain pod types from termination, or coordinate with the DaemonSet controller about node shutdown state.

## References

- [KEP-2000: Graceful Node Shutdown](https://github.com/kubernetes/enhancements/issues/2000)
- [#122912 - DaemonSet controller disagreement](https://github.com/kubernetes/kubernetes/issues/122912)
- [#124448 - System critical pod status](https://github.com/kubernetes/kubernetes/issues/124448)
- [#117018 - DaemonSet stuck in Completed](https://github.com/kubernetes/kubernetes/issues/117018)
- [#131163 - Deployment stall after shutdown](https://github.com/kubernetes/kubernetes/issues/131163)
- [#107043 - Systemd version incompatibility](https://github.com/kubernetes/kubernetes/issues/107043)
