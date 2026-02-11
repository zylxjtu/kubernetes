# Windows E2E Test Images

## Images with Windows Variants — Upgrade Breakdown

### Busybox-based (Windows base: `REGISTRY/busybox`)

| Image | Version | Used in Windows e2e? | ltsc2025 Status |
|-------|---------|---------------------|-----------------|
| nginx | 1.27.0-0 → 1.27.0-1 | Yes (hybrid_network.go, reboot_node.go) | In this PR |
| nginx-new | 1.28.0-0 → 1.28.0-1 | No (not referenced in test/e2e/windows/) | In this PR |
| glibc-dns-testing | 2.0.0 → 2.1.0 | No (not referenced in test/e2e/windows/) | In this PR |

Note: busybox base updated from `1.29-2` to `1.37.0-2` (the version with ltsc2025).

### Agnhost-based (Windows base: `REGISTRY/agnhost`)

| Image | Version | Used in Windows e2e? | ltsc2025 Status |
|-------|---------|---------------------|-----------------|
| kitten | 1.8 → 1.9 | No (not referenced in test/e2e/windows/) | In this PR |
| nautilus | 1.8.0 → 1.9.0 | No (not referenced in test/e2e/windows/) | In this PR |

### Nanoserver-based (Windows base: `mcr.microsoft.com/windows/nanoserver`)

| Image | Version | Used in Windows e2e? | ltsc2025 Status |
|-------|---------|---------------------|-----------------|
| nonroot | 1.5.0 | Yes (security_context.go) | Not yet — nanoserver:ltsc2025 not published |
| resource-consumer | 1.14.0 | Yes (eviction.go) | Not yet — nanoserver:ltsc2025 not published |
| sample-apiserver | 1.33.7 | No (not referenced in test/e2e/windows/) | Not yet — nanoserver:ltsc2025 not published |

## Images Used in Windows E2E Tests (test/e2e/windows/)

| Image | Test Files | Count |
|-------|-----------|-------|
| Agnhost | node_shutdown, gmsa_full, security_context, hybrid_network, service, reboot_node, host_process, dns | 8 |
| BusyBox | gmsa_kubelet, hyperv, gmsa_full, security_context, host_process, kubelet_stats, cpu_limits | 7 |
| Pause | volumes, gmsa_full, memory_limits, security_context, density | 5 |
| Nginx | hybrid_network, reboot_node | 2 |
| NonRoot | security_context | 1 |
| ResourceConsumer | eviction | 1 |
