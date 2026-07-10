/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package memorymanager

import (
	"context"
	"fmt"

	cadvisorapi "github.com/google/cadvisor/lib/model"

	v1 "k8s.io/api/core/v1"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	resourcehelper "k8s.io/component-helpers/resource"
	"k8s.io/klog/v2"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/kubelet/cm/memorymanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
)

// cpuFollowingChecker is optionally implemented by the affinity Store. On Windows
// the memory manager's affinity store (cpuFollowingStore) reports, via
// HasExclusiveCPUs, whether a container is following the CPU manager's NUMA
// decision because it owns exclusive CPUs. When it is, the BestEffort policy must
// not extend the hint to additional NUMA nodes: memory placement has to stay
// aligned with CPU affinity (Windows has no cpuset.mems, so the OS serves a
// thread's pages from the NUMA node of the CPU it runs on).
type cpuFollowingChecker interface {
	HasExclusiveCPUs(podUID, containerName string) bool
}

// On Windows we want to use the same logic as the StaticPolicy to compute the memory topology hints
// but unlike linux based systems, on Windows systems numa nodes cannot be directly assigned or guaranteed via Windows APIs
// (windows scheduler will use the numa node that is closest to the cpu assigned therefor respecting the numa node assignment as a best effort). Because of this we don't want to have users specify "StaticPolicy" for the memory manager
// policy via kubelet configuration. Instead we want to use the "BestEffort" policy which will use the same logic as the StaticPolicy
// and doing so will reduce code duplication.
const policyTypeBestEffort policyType = "BestEffort"

// bestEffortPolicy is implementation of the policy interface for the BestEffort policy
type bestEffortPolicy struct {
	static *staticPolicy
}

var _ Policy = &bestEffortPolicy{}

func NewPolicyBestEffort(logger klog.Logger, machineInfo *cadvisorapi.MachineInfo, reserved systemReservedMemory, affinity topologymanager.Store) (Policy, error) {
	p, err := NewPolicyStatic(logger, machineInfo, reserved, affinity)

	if err != nil {
		return nil, err
	}

	return &bestEffortPolicy{
		static: p.(*staticPolicy),
	}, nil
}

func (p *bestEffortPolicy) Name() string {
	return string(policyTypeBestEffort)
}

func (p *bestEffortPolicy) Start(logger klog.Logger, s state.State) error {
	return p.static.Start(logger, s)
}

// Allocate mirrors staticPolicy.Allocate, with one Windows-specific difference:
// when the container follows the CPU manager's NUMA decision (it owns exclusive
// CPUs), the hint is NOT extended to additional NUMA nodes. Windows has no
// cpuset.mems, so the OS serves a thread's pages from the NUMA node of the CPU it
// runs on, and spreading the memory hint across nodes would diverge from the CPU
// manager's exclusive-CPU binding. Windows memory affinity is best-effort: the
// CPU's node is kept even when the per-NUMA bookkeeping cannot fully satisfy the
// request here.
// Keep this in sync with staticPolicy.Allocate.
func (p *bestEffortPolicy) Allocate(ctx context.Context, s state.State, pod *v1.Pod, container *v1.Container) (rerr error) {
	st := p.static
	logger := klog.LoggerWithValues(klog.FromContext(ctx), "pod", klog.KObj(pod), "containerName", container.Name)

	podUID := string(pod.UID)
	// Allocate the memory only for guaranteed pods.
	qos := v1qos.GetPodQOS(pod)
	if qos != v1.PodQOSGuaranteed {
		logger.V(5).Info("Exclusive memory allocation skipped, pod QoS is not guaranteed", "qos", qos)
		return nil
	}

	if (!utilfeature.DefaultFeatureGate.Enabled(features.PodLevelResourceManagers) || !utilfeature.DefaultFeatureGate.Enabled(features.PodLevelResources)) && resourcehelper.IsPodLevelResourcesSet(pod) {
		logger.V(2).Info("Allocation skipped, pod is using pod-level resources which are not supported by the BestEffort Memory manager policy", "podUID", podUID)
		return nil
	}

	logger.Info("Allocate")
	// Container belongs in an exclusively allocated pool.
	metrics.MemoryManagerPinningRequestTotal.Inc()
	defer func() {
		if rerr != nil {
			metrics.MemoryManagerPinningErrorsTotal.Inc()
			metrics.ResourceManagerAllocationErrorsTotal.WithLabelValues(metrics.ResourceManagerMemory, metrics.ResourceManagerNode).Inc()
		}
	}()
	if blocks := s.GetMemoryBlocks(podUID, container.Name); blocks != nil {
		st.updatePodReusableMemory(pod, container, blocks)

		logger.Info("Container already present in state, skipping")
		return nil
	}

	// Call Topology Manager to get the aligned affinity across all hint providers.
	hint := st.affinity.GetAffinity(logger, podUID, container.Name)
	logger.Info("Got topology affinity", "hint", hint)

	requestedResources, err := getContainerRequestedResources(logger, pod, container)
	if err != nil {
		return err
	}

	machineState := s.GetMachineState()
	bestHint := &hint
	// topology manager returned the hint with NUMA affinity nil
	// we should use the default NUMA affinity calculated the same way as for the topology manager
	if hint.NUMANodeAffinity == nil {
		defaultHint, err := st.getDefaultHint(machineState, pod, requestedResources)
		if err != nil {
			return err
		}

		if !defaultHint.Preferred && bestHint.Preferred {
			return fmt.Errorf("[memorymanager] failed to find the default preferred hint")
		}
		bestHint = defaultHint
	}

	// topology manager returns the hint that does not satisfy completely the container request
	// we should extend this hint to the one who will satisfy the request and include the current hint
	if !isAffinitySatisfyRequest(machineState, bestHint.NUMANodeAffinity, requestedResources) {
		if p.followsCPUManager(podUID, container.Name) {
			// Windows: the container is following the CPU manager's NUMA decision,
			// so keep the CPU's node(s) and do not extend. Memory locality is
			// best-effort and provided through CPU affinity.
			logger.V(4).Info("Following CPU manager NUMA decision; skipping memory hint extension", "hint", bestHint)
		} else {
			extendedHint, err := st.extendTopologyManagerHint(machineState, pod, requestedResources, bestHint.NUMANodeAffinity)
			if err != nil {
				return err
			}

			if !extendedHint.Preferred && bestHint.Preferred {
				return fmt.Errorf("[memorymanager] failed to find the extended preferred hint")
			}
			bestHint = extendedHint
		}
	}

	// the best hint might violate the NUMA allocation rule on which
	// NUMA node cannot have both single and cross NUMA node allocations
	// https://kubernetes.io/blog/2021/08/11/kubernetes-1-22-feature-memory-manager-moves-to-beta/#single-vs-cross-numa-node-allocation
	if isAffinityViolatingNUMAAllocations(machineState, bestHint.NUMANodeAffinity) {
		return fmt.Errorf("[memorymanager] preferred hint violates NUMA node allocation")
	}

	var containerBlocks []state.Block
	maskBits := bestHint.NUMANodeAffinity.GetBits()
	for resourceName, requestedSize := range requestedResources {
		// update memory blocks
		containerBlocks = append(containerBlocks, state.Block{
			NUMAAffinity: maskBits,
			Size:         requestedSize,
			Type:         resourceName,
		})

		podReusableMemory := st.getPodReusableMemory(pod, bestHint.NUMANodeAffinity, resourceName)
		if podReusableMemory >= requestedSize {
			requestedSize = 0
		} else {
			requestedSize -= podReusableMemory
		}

		// Update nodes memory state
		st.updateMachineState(machineState, maskBits, resourceName, requestedSize)
		metrics.ResourceManagerAllocationsTotal.WithLabelValues(metrics.ResourceManagerMemory, metrics.ResourceManagerNode).Inc()
	}

	st.updatePodReusableMemory(pod, container, containerBlocks)

	s.SetMachineState(machineState)
	s.SetMemoryBlocks(podUID, container.Name, containerBlocks)

	st.updateInitContainersMemoryBlocks(logger, s, pod, container, containerBlocks)
	metrics.ResourceManagerContainerAssignments.WithLabelValues(metrics.ResourceManagerMemory, metrics.ResourceManagerExclusiveNode).Inc()

	logger.V(4).Info("Allocated exclusive memory")
	return nil
}

// followsCPUManager reports whether the container is following the CPU manager's
// NUMA decision: its affinity store is the Windows cpuFollowingStore and the
// container owns exclusive CPUs. When true, the memory hint must not be extended.
func (p *bestEffortPolicy) followsCPUManager(podUID, containerName string) bool {
	checker, ok := p.static.affinity.(cpuFollowingChecker)
	return ok && checker.HasExclusiveCPUs(podUID, containerName)
}

func (p *bestEffortPolicy) RemoveContainer(logger klog.Logger, s state.State, podUID string, containerName string) {
	p.static.RemoveContainer(logger, s, podUID, containerName)
}

func (p *bestEffortPolicy) GetPodTopologyHints(_ klog.Logger, s state.State, pod *v1.Pod) map[string][]topologymanager.TopologyHint {
	// Pod-level resources are not supported on Windows.
	return nil
}

func (p *bestEffortPolicy) GetTopologyHints(logger klog.Logger, s state.State, pod *v1.Pod, container *v1.Container) map[string][]topologymanager.TopologyHint {
	return p.static.GetTopologyHints(logger, s, pod, container)
}

func (p *bestEffortPolicy) GetAllocatableMemory(s state.State) []state.Block {
	return p.static.GetAllocatableMemory(s)
}

func (p *bestEffortPolicy) AllocatePod(_ klog.Logger, s state.State, pod *v1.Pod) error {
	// Pod-level resources are not supported on Windows.
	return nil
}
