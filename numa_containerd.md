  ## hcsshim PR #2603 — CPU Affinity for Job Containers

  **PR**: https://github.com/microsoft/hcsshim/pull/2603

  ### What the PR Does

  Adds CPU affinity support to job containers by:
  - Parsing `Windows.Resources.CPU.Affinity` from the OCI spec
  - Storing it as a single `uint64` mask in `JobLimits.CPUAffinity`
  - Applying it via `SetCPUAffinity()` on the job object

  ### Limitations (Single Group Only)

  The implementation explicitly rejects multi-group configurations:

  ```go
  if len(affinity) != 1 {
      return nil, fmt.Errorf("cpu affinity with multiple processor groups is not supported")
  }
  if affinity[0].Group != 0 {
      return nil, fmt.Errorf("cpu affinity processor group %d is not supported", affinity[0].Group)
  }

  Issues vs. Windows Server 2022+ Multi-Group NUMA

  ┌─────────────────────────────────┬───────────────────────────────────────┐
  │            Criteria             │                Status                 │
  ├─────────────────────────────────┼───────────────────────────────────────┤
  │ Supports 64+ processors         │ ❌ Single uint64 mask = max 64 CPUs   │
  ├─────────────────────────────────┼───────────────────────────────────────┤
  │ Supports multi-group NUMA nodes │ ❌ Explicitly rejected                │
  ├─────────────────────────────────┼───────────────────────────────────────┤
  │ Uses CPU Sets (modern API)      │ ❌ Uses legacy affinity mask          │
  ├─────────────────────────────────┼───────────────────────────────────────┤
  │ Supports non-zero groups        │ ❌ Hard-fails if Group != 0           │
  ├─────────────────────────────────┼───────────────────────────────────────┤
  │ Works on large VMs (96+ vCPUs)  │ ❌ Fails or only covers first 64 CPUs │
  └─────────────────────────────────┴───────────────────────────────────────┘

  What Would Need to Change

  1. Accept multiple GROUP_AFFINITY entries — don't reject len(affinity) != 1
  2. Allow non-zero groups — remove the Group != 0 gate
  3. Use Job Object CPU Sets (via SetInformationJobObject with CPU set IDs) instead of legacy JOB_OBJECT_LIMIT_AFFINITY which is single-group
  4. Alternatively, apply per-group affinity with SetThreadGroupAffinity on threads within the job
