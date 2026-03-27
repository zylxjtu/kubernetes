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
```

### Issues vs. Windows Server 2022+ Multi-Group NUMA

| Criteria | Status |
|---|---|
| Supports 64+ processors | ❌ Single `uint64` mask = max 64 CPUs |
| Supports multi-group NUMA nodes | ❌ Explicitly rejected |
| Uses CPU Sets (modern API) | ❌ Uses legacy affinity mask |
| Supports non-zero groups | ❌ Hard-fails if Group != 0 |
| Works on large VMs (96+ vCPUs) | ❌ Fails or only covers first 64 CPUs |

### Suggestions for fixing

#### 1. Accept multiple `GROUP_AFFINITY` entries

Remove the single-group restriction:

```go
// BEFORE (broken):
if len(affinity) != 1 {
    return nil, fmt.Errorf("cpu affinity with multiple processor groups is not supported")
}
if affinity[0].Group != 0 {
    return nil, fmt.Errorf("cpu affinity processor group %d is not supported", affinity[0].Group)
}
cpuAffinity = affinity[0].Mask

// AFTER: accept all entries, allow any group
for _, a := range affinity {
    if a.Mask == 0 {
        return nil, fmt.Errorf("cpu affinity mask must be non-zero")
    }
}
```

#### 2. Change `JobLimits.CPUAffinity` from `uint64` to `[]GROUP_AFFINITY`

A single `uint64` can only represent 64 processors in one group. The struct should hold an array:

```go
type JobLimits struct {
    // ...
    CPUAffinity []winnt.GROUP_AFFINITY  // was: uint64
}
```

#### 3. Use Job Object CPU Sets instead of legacy affinity

The legacy `JOB_OBJECT_LIMIT_AFFINITY` (set via `JOBOBJECT_BASIC_LIMIT_INFORMATION.Affinity`) is single-group only. Two options:

**Option A — CPU Sets (preferred):**
- Convert each `GROUP_AFFINITY` (group + mask) to CPU Set IDs
- Call `SetInformationJobObject` with `JobObjectCpuRateControlInformation` or assign CPU Sets to the job
- CPU Sets are inherently multi-group aware

**Option B — Per-group affinity via `JOBOBJECT_GROUP_INFORMATION`:**
- Use `SetInformationJobObject` with `JobObjectGroupInformationEx`
- Pass an array of `GROUP_AFFINITY` structs, one per group
- This allows the job object to span multiple processor groups

#### 4. Update tests

Current tests only validate group 0 with a single mask. Add:
- Test with multiple `GROUP_AFFINITY` entries (e.g., group 0 + group 1)
- Test with non-zero group index
- Test on a VM with 64+ vCPUs to verify cross-group behavior

