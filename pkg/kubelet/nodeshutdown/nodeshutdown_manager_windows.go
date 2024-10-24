//go:build windows
// +build windows

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

// Package nodeshutdown can watch for node level shutdown events and trigger graceful termination of pods running on the node prior to a system shutdown.
package nodeshutdown

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/features"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	kubeletevents "k8s.io/kubernetes/pkg/kubelet/events"
	"k8s.io/kubernetes/pkg/kubelet/eviction"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	"k8s.io/kubernetes/pkg/kubelet/prober"
	"k8s.io/kubernetes/pkg/kubelet/volumemanager"
	"k8s.io/kubernetes/pkg/windows/service"
	"k8s.io/utils/clock"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	// Kubelet service name
	serviceKubelet           = "kubelet"
	shutdownOrderRegPath     = `SYSTEM\CurrentControlSet\Control`
	shutdownOrderStringValue = "PreshutdownOrder"
)

const (
	nodeShutdownReason             = "Terminated"
	nodeShutdownMessage            = "Pod was terminated in response to imminent node shutdown."
	nodeShutdownNotAdmittedReason  = "NodeShutdown"
	nodeShutdownNotAdmittedMessage = "Pod was rejected as the node is shutting down."
	localStorageStateFile          = "graceful_node_shutdown_state"
)

// managerImpl has functions that can be used to interact with the Node Shutdown Manager.
type managerImpl struct {
	logger       klog.Logger
	recorder     record.EventRecorder
	nodeRef      *v1.ObjectReference
	probeManager prober.Manager

	volumeManager volumemanager.VolumeManager

	shutdownGracePeriodByPodPriority []kubeletconfig.ShutdownGracePeriodByPodPriority

	getPods        eviction.ActivePodsFunc
	killPodFunc    eviction.KillPodFunc
	syncNodeStatus func()

	nodeShuttingDownMutex sync.Mutex
	nodeShuttingDownNow   bool

	clock clock.Clock

	enableMetrics bool
	storage       storage
}

// NewManager returns a new node shutdown manager.
func NewManager(conf *Config) (Manager, lifecycle.PodAdmitHandler) {
	if !utilfeature.DefaultFeatureGate.Enabled(features.WindowsGracefulNodeShutdown) {
		m := managerStub{}
		return m, m
	}

	shutdownGracePeriodByPodPriority := conf.ShutdownGracePeriodByPodPriority
	// Migration from the original configuration
	if len(shutdownGracePeriodByPodPriority) == 0 {
		shutdownGracePeriodByPodPriority = migrateConfig(conf.ShutdownGracePeriodRequested, conf.ShutdownGracePeriodCriticalPods)
	}

	// Disable if the configuration is empty
	if len(shutdownGracePeriodByPodPriority) == 0 {
		m := managerStub{}
		return m, m
	}

	// Sort by priority from low to high
	sort.Slice(shutdownGracePeriodByPodPriority, func(i, j int) bool {
		return shutdownGracePeriodByPodPriority[i].Priority < shutdownGracePeriodByPodPriority[j].Priority
	})

	if conf.Clock == nil {
		conf.Clock = clock.RealClock{}
	}
	manager := &managerImpl{
		logger:       conf.Logger,
		probeManager: conf.ProbeManager,
		recorder:     conf.Recorder,
		//volumeManager:                    conf.VolumeManager,
		nodeRef:                          conf.NodeRef,
		getPods:                          conf.GetPodsFunc,
		killPodFunc:                      conf.KillPodFunc,
		syncNodeStatus:                   conf.SyncNodeStatusFunc,
		shutdownGracePeriodByPodPriority: shutdownGracePeriodByPodPriority,
		clock:                            conf.Clock,
		enableMetrics:                    utilfeature.DefaultFeatureGate.Enabled(features.GracefulNodeShutdownBasedOnPodPriority),
		storage: localStorage{
			Path: filepath.Join(conf.StateDirectory, localStorageStateFile),
		},
	}
	manager.logger.Info("Creating node shutdown manager",
		"shutdownGracePeriodRequested", conf.ShutdownGracePeriodRequested,
		"shutdownGracePeriodCriticalPods", conf.ShutdownGracePeriodCriticalPods,
		"shutdownGracePeriodByPodPriority", shutdownGracePeriodByPodPriority,
	)
	return manager, manager
}

// Admit rejects all pods if node is shutting
func (m *managerImpl) Admit(attrs *lifecycle.PodAdmitAttributes) lifecycle.PodAdmitResult {
	nodeShuttingDown := m.ShutdownStatus() != nil

	if nodeShuttingDown {
		return lifecycle.PodAdmitResult{
			Admit:   false,
			Reason:  nodeShutdownNotAdmittedReason,
			Message: nodeShutdownNotAdmittedMessage,
		}
	}
	return lifecycle.PodAdmitResult{Admit: true}
}

// setMetrics sets the metrics for the node shutdown manager.
func (m *managerImpl) setMetrics() {
	if m.enableMetrics && m.storage != nil {
		sta := state{}
		err := m.storage.Load(&sta)
		if err != nil {
			m.logger.Error(err, "Failed to load graceful shutdown state")
		} else {
			if !sta.StartTime.IsZero() {
				metrics.GracefulShutdownStartTime.Set(timestamp(sta.StartTime))
			}
			if !sta.EndTime.IsZero() {
				metrics.GracefulShutdownEndTime.Set(timestamp(sta.EndTime))
			}
		}
	}
}

// Start starts the node shutdown manager and will start watching the node for shutdown events.
func (m *managerImpl) Start() error {
	m.logger.V(1).Info("Shutdown manager get started")

	_, err := m.start()

	if err != nil {
		return err
	}

	service.SetPreShutdownHandler(m)

	m.setMetrics()

	return nil
}

// Check the return value, maybe return error only
func (m *managerImpl) start() (chan struct{}, error) {
	// Process the shutdown only when it is running as a windows service
	m.logger.V(1).Info("Shutdown manager entering start")

	isServiceIntialized := service.IsServiceInitialized()
	m.logger.V(1).Info("is service check", "isServiceIntialized", isServiceIntialized)
	if !isServiceIntialized {
		return nil, errors.Errorf("%s is NOT running as a Windows service", serviceKubelet)
	}
	m.logger.V(1).Info("Shutdown manage finished checking if running as a service")
	// Update the registry key to add the kubelet dependencies to the existing order
	mgr, err := mgr.Connect()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not connect to service manager")
	}
	defer mgr.Disconnect()
	m.logger.V(1).Info("Shutdown manage finished connecting to service manager")
	s, err := mgr.OpenService(serviceKubelet)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not access service %s", serviceKubelet)
	}
	defer s.Close()

	m.logger.V(1).Info("Shutdown manager trying to get reg key")
	config, err := s.Config()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not access config of service %s", serviceKubelet)
	}

	// Open the registry key
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, shutdownOrderRegPath, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not access registry")
	}
	defer key.Close()

	// Read the existing values
	existingOrders, _, err := key.GetStringsValue(shutdownOrderStringValue)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not access registry value %s", shutdownOrderStringValue)
	}

	m.logger.V(1).Info("Shutdown manager get existing orders", "existingorders", existingOrders)

	// Add the kubelet dependencies to the existing order
	newOrders := addToExistingOrder(config.Dependencies, existingOrders)

	m.logger.V(1).Info("Shutdown manager created new orders", "neworders", newOrders)

	err = key.SetStringsValue("PreshutdownOrder", newOrders)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not set registry %s to be new value %s", shutdownOrderStringValue, newOrders)
	}

	// If the preshutdown timeout is less than periodRequested, attempt to update the value to periodRequested.
	preshutdownInfo, err := service.QueryPreShutdownInfo(s.Handle)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not query preshutdown info")
	}

	m.logger.V(1).Info("Shutdown manager get current preshutdown info", "PreshutdownTimeout", preshutdownInfo.PreshutdownTimeout)
	if periodRequested := m.periodRequested().Milliseconds(); periodRequested > int64(preshutdownInfo.PreshutdownTimeout) {
		m.logger.V(1).Info("Shutdown manager override preshutdown info", "ShutdownGracePeriod", periodRequested)
		err := service.UpdatePreShutdownInfo(s.Handle, uint32(periodRequested))
		if err != nil {
			return nil, fmt.Errorf("Unable to override preshoutdown config by shutdown manager: %v", err)
		}

		// Read the preshutdownInfo again, if the override was successful, preshutdownInfo will be equal to shutdownGracePeriodRequested.
		preshutdownInfo, err := service.QueryPreShutdownInfo(s.Handle)
		if err != nil {
			return nil, err
		}

		if periodRequested > int64(preshutdownInfo.PreshutdownTimeout) {
			return nil, fmt.Errorf("Shutdown manager was unable to update preshutdown info to %v (ShutdownGracePeriod), current value of preshutdown info (%v) is less than requested ShutdownGracePeriod", periodRequested, preshutdownInfo.PreshutdownTimeout)
		}
	}

	return nil, nil
}

// ShutdownStatus will return an error if the node is currently shutting down.
func (m *managerImpl) ShutdownStatus() error {
	m.nodeShuttingDownMutex.Lock()
	defer m.nodeShuttingDownMutex.Unlock()

	if m.nodeShuttingDownNow {
		return fmt.Errorf("node is shutting down")
	}
	return nil
}

func (m *managerImpl) ProcessShutdownEvent() error {
	m.logger.V(1).Info("Shutdown manager detected new preshutdown event", "event", "preshutdown")

	m.recorder.Event(m.nodeRef, v1.EventTypeNormal, kubeletevents.NodeShutdown, "Shutdown manager detected preshutdown event")

	m.nodeShuttingDownMutex.Lock()
	m.nodeShuttingDownNow = true
	m.nodeShuttingDownMutex.Unlock()

	go m.syncNodeStatus()

	m.logger.V(1).Info("Shutdown manager processing preshutdown event")
	activePods := m.getPods()

	defer func() {
		m.logger.V(1).Info("Shutdown manager completed processing preshutdown event, node will shutdown shortly")
	}()

	if m.enableMetrics && m.storage != nil {
		startTime := time.Now()
		err := m.storage.Store(state{
			StartTime: startTime,
		})
		if err != nil {
			m.logger.Error(err, "Failed to store graceful shutdown state")
		}
		metrics.GracefulShutdownStartTime.Set(timestamp(startTime))
		metrics.GracefulShutdownEndTime.Set(0)

		defer func() {
			endTime := time.Now()
			err := m.storage.Store(state{
				StartTime: startTime,
				EndTime:   endTime,
			})
			if err != nil {
				m.logger.Error(err, "Failed to store graceful shutdown state")
			}
			metrics.GracefulShutdownEndTime.Set(timestamp(endTime))
		}()
	}

	groups := groupByPriority(m.shutdownGracePeriodByPodPriority, activePods)
	for _, group := range groups {
		// If there are no pods in a particular range,
		// then do not wait for pods in that priority range.
		if len(group.Pods) == 0 {
			continue
		}

		var wg sync.WaitGroup
		wg.Add(len(group.Pods))
		for _, pod := range group.Pods {
			go func(pod *v1.Pod, group podShutdownGroup) {
				defer wg.Done()

				gracePeriodOverride := group.ShutdownGracePeriodSeconds

				// If the pod's spec specifies a termination gracePeriod which is less than the gracePeriodOverride calculated, use the pod spec termination gracePeriod.
				if pod.Spec.TerminationGracePeriodSeconds != nil && *pod.Spec.TerminationGracePeriodSeconds <= gracePeriodOverride {
					gracePeriodOverride = *pod.Spec.TerminationGracePeriodSeconds
				}

				m.logger.V(1).Info("Shutdown manager killing pod with gracePeriod", "pod", klog.KObj(pod), "gracePeriod", gracePeriodOverride)

				if err := m.killPodFunc(pod, false, &gracePeriodOverride, func(status *v1.PodStatus) {
					// set the pod status to failed (unless it was already in a successful terminal phase)
					if status.Phase != v1.PodSucceeded {
						status.Phase = v1.PodFailed
					}
					status.Message = nodeShutdownMessage
					status.Reason = nodeShutdownReason
					podutil.UpdatePodCondition(status, &v1.PodCondition{
						Type:    v1.DisruptionTarget,
						Status:  v1.ConditionTrue,
						Reason:  v1.PodReasonTerminationByKubelet,
						Message: nodeShutdownMessage,
					})
				}); err != nil {
					m.logger.V(1).Info("Shutdown manager failed killing pod", "pod", klog.KObj(pod), "err", err)
				} else {
					m.logger.V(1).Info("Shutdown manager finished killing pod", "pod", klog.KObj(pod))
				}
			}(pod, group)
		}

		// This duration determines how long the shutdown manager will wait for the pods in this group
		// to terminate before proceeding to the next group.
		var groupTerminationWaitDuration = time.Duration(group.ShutdownGracePeriodSeconds) * time.Second
		var (
			doneCh       = make(chan struct{})
			timer        = m.clock.NewTimer(groupTerminationWaitDuration)
			_, ctxCancel = context.WithTimeout(context.Background(), groupTerminationWaitDuration)
		)
		go func() {
			defer close(doneCh)
			defer ctxCancel()
			wg.Wait()
			// The signal to kill a Pod was sent successfully to all the pods,
			// let's wait until all the volumes are unmounted from all the pods before
			// continuing to the next group. This is done so that the CSI Driver (assuming
			// that it's part of the highest group) has a chance to perform unmounts.
			// if err := m.volumeManager.WaitForAllPodsUnmount(ctx, group.Pods); err != nil {
			// 	var podIdentifiers []string
			// 	for _, pod := range group.Pods {
			// 		podIdentifiers = append(podIdentifiers, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
			// 	}

			// 	// Waiting for volume teardown is done on a best basis effort,
			// 	// report an error and continue.
			// 	//
			// 	// Depending on the user provided kubelet configuration value
			// 	// either the `timer` will tick and we'll continue to shutdown the next group, or,
			// 	// WaitForAllPodsUnmount will timeout, therefore this goroutine
			// 	// will close doneCh and we'll continue to shutdown the next group.
			// 	m.logger.Error(err, "Failed while waiting for all the volumes belonging to Pods in this group to unmount", "pods", podIdentifiers)
			//}
		}()

		select {
		case <-doneCh:
			timer.Stop()
		case <-timer.C():
			ctxCancel()
			m.logger.V(1).Info("Shutdown manager pod killing time out", "gracePeriod", group.ShutdownGracePeriodSeconds, "priority", group.Priority)
		}
	}

	return nil
}

func (m *managerImpl) periodRequested() time.Duration {
	var sum int64
	for _, period := range m.shutdownGracePeriodByPodPriority {
		sum += period.ShutdownGracePeriodSeconds
	}
	return time.Duration(sum) * time.Second
}

// Helper function to remove all occurrences of a specific item from a string list
func removeItemFromList(stringlist []string, item string) []string {
	writeIndex := 0

	// Iterate through the list and only keep those don't match the item (case-insensitive)
	for _, listItem := range stringlist {
		if !strings.EqualFold(listItem, item) {
			stringlist[writeIndex] = listItem
			writeIndex++
		}
	}

	// Return the modified slice, trimmed to the valid length
	return stringlist[:writeIndex]
}

// Helper function to insert an element into a slice at a specified index
func insertAt(slice []string, index int, value string) []string {
	// If the index is greater than or equal to the length, append the element to the end
	if index >= len(slice) {
		return append(slice, value)
	}

	// Ensure there's enough space in the slice by appending a zero-value element first
	slice = append(slice, "")
	copy(slice[index+1:], slice[index:])
	slice[index] = value

	return slice
}

// Dependencies: ["a", "b", "c", "d"]
// ExistingOrder: ["x", "b", "y", "c", "z"]
// The output will be:
// Modified List: ["x", "kubelet", "b", "y", "c", "z", "a", "d"]
func addToExistingOrder(dependencies []string, existingOrder []string) []string {
	//Remove "Kubelet" from existing order if any
	existingOrder = removeItemFromList(existingOrder, serviceKubelet)

	// Append items from dependencies to existingOrder if not already present
	existingOrderMap := make(map[string]bool)
	for _, item := range existingOrder {
		existingOrderMap[item] = true
	}
	// Append non-duplicate items from dependencies to existingOrder
	for _, item := range dependencies {
		if !existingOrderMap[item] {
			existingOrder = append(existingOrder, item)
		}
	}

	// Insert "kubelet" before the first common item
	// Honor the order of existing order
	firstCommonIndex := -1
	for i, item := range existingOrder {
		for _, item1 := range dependencies {
			if item == item1 {
				firstCommonIndex = i
				break
			}
		}
		if firstCommonIndex != -1 {
			break
		}
	}

	// If a common item is found, insert "kubelet" before it
	if firstCommonIndex != -1 {
		existingOrder = insertAt(existingOrder, firstCommonIndex, serviceKubelet)
	}

	return existingOrder
}
