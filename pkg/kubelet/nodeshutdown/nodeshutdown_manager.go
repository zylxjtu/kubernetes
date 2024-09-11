/*
Copyright 2021 The Kubernetes Authors.

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

package nodeshutdown

import (
	"sort"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/apis/scheduling"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/prober"
	"k8s.io/utils/clock"
)

// Manager interface provides methods for Kubelet to manage node shutdown.
type Manager interface {
	Admit(attrs *lifecycle.PodAdmitAttributes) lifecycle.PodAdmitResult
	Start() error
	ShutdownStatus() error
}

// Config represents Manager configuration
type Config struct {
	Logger                           klog.Logger
	ProbeManager                     prober.Manager
	Recorder                         record.EventRecorder
	NodeRef                          *v1.ObjectReference
	GetPodsFunc                      eviction.ActivePodsFunc
	KillPodFunc                      eviction.KillPodFunc
	SyncNodeStatusFunc               func()
	ShutdownGracePeriodRequested     time.Duration
	ShutdownGracePeriodCriticalPods  time.Duration
	ShutdownGracePeriodByPodPriority []kubeletconfig.ShutdownGracePeriodByPodPriority
	StateDirectory                   string
	Clock                            clock.Clock
}

// managerStub is a fake node shutdown managerImpl .
type managerStub struct{}

// Admit returns a fake Pod admission which always returns true
func (managerStub) Admit(attrs *lifecycle.PodAdmitAttributes) lifecycle.PodAdmitResult {
	return lifecycle.PodAdmitResult{Admit: true}
}

// Start is a no-op always returning nil for non linux platforms.
func (managerStub) Start() error {
	return nil
}

// ShutdownStatus is a no-op always returning nil for non linux platforms.
func (managerStub) ShutdownStatus() error {
	return nil
}

func migrateConfig(shutdownGracePeriodRequested, shutdownGracePeriodCriticalPods time.Duration) []kubeletconfig.ShutdownGracePeriodByPodPriority {
	if shutdownGracePeriodRequested == 0 {
		return nil
	}
	defaultPriority := shutdownGracePeriodRequested - shutdownGracePeriodCriticalPods
	if defaultPriority < 0 {
		return nil
	}
	criticalPriority := shutdownGracePeriodRequested - defaultPriority
	if criticalPriority < 0 {
		return nil
	}
	return []kubeletconfig.ShutdownGracePeriodByPodPriority{
		{
			Priority:                   scheduling.DefaultPriorityWhenNoDefaultClassExists,
			ShutdownGracePeriodSeconds: int64(defaultPriority / time.Second),
		},
		{
			Priority:                   scheduling.SystemCriticalPriority,
			ShutdownGracePeriodSeconds: int64(criticalPriority / time.Second),
		},
	}
}

func groupByPriority(shutdownGracePeriodByPodPriority []kubeletconfig.ShutdownGracePeriodByPodPriority, pods []*v1.Pod) []podShutdownGroup {
	groups := make([]podShutdownGroup, 0, len(shutdownGracePeriodByPodPriority))
	for _, period := range shutdownGracePeriodByPodPriority {
		groups = append(groups, podShutdownGroup{
			ShutdownGracePeriodByPodPriority: period,
		})
	}

	for _, pod := range pods {
		var priority int32
		if pod.Spec.Priority != nil {
			priority = *pod.Spec.Priority
		}

		// Find the group index according to the priority.
		index := sort.Search(len(groups), func(i int) bool {
			return groups[i].Priority >= priority
		})

		// 1. Those higher than the highest priority default to the highest priority
		// 2. Those lower than the lowest priority default to the lowest priority
		// 3. Those boundary priority default to the lower priority
		// if priority of pod is:
		//   groups[index-1].Priority <= pod priority < groups[index].Priority
		// in which case we want to pick lower one (i.e index-1)
		if index == len(groups) {
			index = len(groups) - 1
		} else if index < 0 {
			index = 0
		} else if index > 0 && groups[index].Priority > priority {
			index--
		}

		groups[index].Pods = append(groups[index].Pods, pod)
	}
	return groups
}

type podShutdownGroup struct {
	kubeletconfig.ShutdownGracePeriodByPodPriority
	Pods []*v1.Pod
}
