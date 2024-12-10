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

package windows

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/kubernetes/pkg/apis/scheduling"
	"k8s.io/kubernetes/pkg/features"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/windows/service"
	"k8s.io/kubernetes/test/e2e/framework"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	. "k8s.io/kubernetes/test/e2e_node/utils"
	admissionapi "k8s.io/pod-security-admission/api"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/kubectl/pkg/util/podutils"
	kubelettypes "k8s.io/kubernetes/pkg/kubelet/types"
	testutils "k8s.io/kubernetes/test/utils"
)

var _ = SIGDescribe("GracefulNodeShutdown", func() {
	f := framework.NewDefaultFramework("windows-node-graceful-shutdown")
	f.NamespacePodSecurityLevel = admissionapi.LevelBaseline
	ginkgo.Context("when windows node gracefully shutting down", func() {
		const (
			pollInterval                        = 1 * time.Second
			podStatusUpdateTimeout              = 90 * time.Second
			nodeStatusUpdateTimeout             = 90 * time.Second
			nodeShutdownGracePeriod             = 30 * time.Second
			nodeShutdownGracePeriodCriticalPods = 20 * time.Second
		)

		TempSetCurrentKubeletConfig(f, func(ctx context.Context, initialConfig *kubeletconfig.KubeletConfiguration) {
			initialConfig.FeatureGates = map[string]bool{
				string(features.WindowsGracefulNodeShutdown):            true,
				string(features.GracefulNodeShutdownBasedOnPodPriority): false,
				string(features.PodReadyToStartContainersCondition):     true,
			}
			initialConfig.ShutdownGracePeriod = metav1.Duration{Duration: nodeShutdownGracePeriod}
			initialConfig.ShutdownGracePeriodCriticalPods = metav1.Duration{Duration: nodeShutdownGracePeriodCriticalPods}
		})

		ginkgo.BeforeEach(func(ctx context.Context) {
			ginkgo.By("Wait for the node to be ready")
			WaitForNodeReady(ctx)
		})

		ginkgo.It("should be able to gracefully shutdown pods with various grace periods", func(ctx context.Context) {
			nodeName := GetNodeName(ctx, f)
			nodeSelector := fields.Set{
				"spec.nodeName": nodeName,
			}.AsSelector().String()

			// Define test pods
			pods := []*v1.Pod{
				getGracePeriodOverrideTestPod("period-20-"+string(uuid.NewUUID()), nodeName, 20, ""),
				getGracePeriodOverrideTestPod("period-25-"+string(uuid.NewUUID()), nodeName, 25, ""),
				getGracePeriodOverrideTestPod("period-critical-5-"+string(uuid.NewUUID()), nodeName, 5, scheduling.SystemNodeCritical),
				getGracePeriodOverrideTestPod("period-critical-10-"+string(uuid.NewUUID()), nodeName, 10, scheduling.SystemNodeCritical),
			}

			ginkgo.By("Creating batch pods")
			e2epod.NewPodClient(f).CreateBatch(ctx, pods)

			list, err := e2epod.NewPodClient(f).List(ctx, metav1.ListOptions{
				FieldSelector: nodeSelector,
			})
			framework.ExpectNoError(err)
			gomega.Expect(list.Items).To(gomega.HaveLen(len(pods)), "the number of pods is not as expected")

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			go func() {
				defer ginkgo.GinkgoRecover()
				w := &cache.ListWatch{
					WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
						return f.ClientSet.CoreV1().Pods(f.Namespace.Name).Watch(ctx, options)
					},
				}

				// Setup watch to continuously monitor any pod events and detect invalid pod status updates
				_, err = watchtools.Until(ctx, list.ResourceVersion, w, func(event watch.Event) (bool, error) {
					if pod, ok := event.Object.(*v1.Pod); ok {
						if isPodStatusAffectedByIssue108594(pod) {
							return false, fmt.Errorf("failing test due to detecting invalid pod status")
						}
						// Watch will never terminate (only when the test ends due to context cancellation)
						return false, nil
					}
					return false, nil
				})

				// Ignore timeout error since the context will be explicitly cancelled and the watch will never return true
				if err != nil && err != wait.ErrWaitTimeout {
					framework.Failf("watch for invalid pod status failed: %v", err.Error())
				}
			}()

			ginkgo.By("Verifying batch pods are running")
			for _, pod := range list.Items {
				if podReady, err := testutils.PodRunningReady(&pod); err != nil || !podReady {
					framework.Failf("Failed to start batch pod: %v", pod.Name)
				}
			}

			ginkgo.By("Emitting shutdown signal")
			err = sendControlToService("kubelet", svc.Cmd(service.PreshutdownSimulationCode))
			framework.ExpectNoError(err)

			ginkgo.By("Verifying that non-critical pods are shutdown")
			// Not critical pod should be shutdown
			gomega.Eventually(ctx, func(ctx context.Context) error {
				list, err = e2epod.NewPodClient(f).List(ctx, metav1.ListOptions{
					FieldSelector: nodeSelector,
				})
				if err != nil {
					return err
				}
				gomega.Expect(list.Items).To(gomega.HaveLen(len(pods)), "the number of pods is not as expected")

				for _, pod := range list.Items {
					if kubelettypes.IsCriticalPod(&pod) {
						if isPodShutdown(&pod) {
							framework.Logf("Expecting critical pod (%v/%v) to be running, but it's not currently. Pod Status %+v", pod.Namespace, pod.Name, pod.Status)
							return fmt.Errorf("critical pod (%v/%v) should not be shutdown, phase: %s", pod.Namespace, pod.Name, pod.Status.Phase)
						}
					} else {
						if !isPodShutdown(&pod) {
							framework.Logf("Expecting non-critical pod (%v/%v) to be shutdown, but it's not currently. Pod Status %+v", pod.Namespace, pod.Name, pod.Status)
							return fmt.Errorf("pod (%v/%v) should be shutdown, phase: %s", pod.Namespace, pod.Name, pod.Status.Phase)
						}
					}
				}
				return nil
			}, podStatusUpdateTimeout, pollInterval).Should(gomega.Succeed())

			ginkgo.By("Verifying that all pods are shutdown")
			// All pod should be shutdown
			gomega.Eventually(ctx, func(ctx context.Context) error {
				list, err = e2epod.NewPodClient(f).List(ctx, metav1.ListOptions{
					FieldSelector: nodeSelector,
				})
				if err != nil {
					return err
				}
				gomega.Expect(list.Items).To(gomega.HaveLen(len(pods)), "the number of pods is not as expected")

				for _, pod := range list.Items {
					if !isPodShutdown(&pod) {
						framework.Logf("Expecting pod (%v/%v) to be shutdown, but it's not currently: Pod Status %+v", pod.Namespace, pod.Name, pod.Status)
						return fmt.Errorf("pod (%v/%v) should be shutdown, phase: %s", pod.Namespace, pod.Name, pod.Status.Phase)
					}
				}
				return nil
			},
				// Critical pod starts shutdown after (nodeShutdownGracePeriod-nodeShutdownGracePeriodCriticalPods)
				podStatusUpdateTimeout+(nodeShutdownGracePeriod-nodeShutdownGracePeriodCriticalPods),
				pollInterval).Should(gomega.Succeed())

			ginkgo.By("Verify that all pod ready to start condition are set to false after terminating")
			// All pod ready to start condition should set to false
			gomega.Eventually(ctx, func(ctx context.Context) error {
				list, err = e2epod.NewPodClient(f).List(ctx, metav1.ListOptions{
					FieldSelector: nodeSelector,
				})
				if err != nil {
					return err
				}
				gomega.Expect(list.Items).To(gomega.HaveLen(len(pods)))
				for _, pod := range list.Items {
					if !isPodReadyToStartConditionSetToFalse(&pod) {
						framework.Logf("Expecting pod (%v/%v) 's ready to start condition set to false, "+
							"but it's not currently: Pod Condition %+v", pod.Namespace, pod.Name, pod.Status.Conditions)
						return fmt.Errorf("pod (%v/%v) 's ready to start condition should be false, condition: %s, phase: %s",
							pod.Namespace, pod.Name, pod.Status.Conditions, pod.Status.Phase)
					}
				}
				return nil
			},
			).Should(gomega.Succeed())
		})
	})
})

func sendControlToService(serviceName string, controlCode svc.Cmd) error {
	// Open the service control manager
	mgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to open service control manager: %w", err)
	}
	defer mgr.Disconnect()

	s, err := mgr.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("failed to get the pointer for service name: %w", err)
	}
	defer s.Close()

	_, err = s.Control(controlCode)
	if err != nil {
		return fmt.Errorf("failed to send control code: %w", err)
	}
	// scmHandle, err := windows.OpenSCManager(nil, nil, windows.SC_MANAGER_CONNECT)
	// if err != nil {
	// 	return fmt.Errorf("failed to open service control manager: %w", err)
	// }
	// defer windows.CloseServiceHandle(scmHandle)

	// Open the service
	// servieNamePointer, err := windows.UTF16PtrFromString(serviceName)
	// if err != nil {
	// 	return fmt.Errorf("failed to get the pointer for service name: %w", err)
	// }
	// serviceHandle, err := windows.OpenService(scmHandle, servieNamePointer, windows.SERVICE_USER_DEFINED_CONTROL|windows.SERVICE_QUERY_STATUS)
	// if err != nil {
	// 	return fmt.Errorf("failed to open service: %w", err)
	// }
	// defer windows.CloseServiceHandle(serviceHandle)

	// Send the control code
	// err = windows.ControlService(s, controlCode, &windows.SERVICE_STATUS{})
	// if err != nil {
	// 	return fmt.Errorf("failed to send control code: %w", err)
	// }

	return nil
}

// getGracePeriodOverrideTestPod returns a new Pod object containing a container
// runs a shell script, hangs the process until a SIGTERM signal is received.
// The script waits for $PID to ensure that the process does not exist.
// If priorityClassName is scheduling.SystemNodeCritical, the Pod is marked as critical and a comment is added.
func getGracePeriodOverrideTestPod(name string, node string, gracePeriod int64, priorityClassName string) *v1.Pod {

	pod := &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    name,
					Image:   AgnhostImage,
					Command: []string{"/agnhost", "netexec", "--delay-shutdown", "9999"},
				},
			},
			TerminationGracePeriodSeconds: &gracePeriod,
			NodeName:                      node,
		},
	}
	if priorityClassName == scheduling.SystemNodeCritical {
		pod.ObjectMeta.Annotations = map[string]string{
			kubelettypes.ConfigSourceAnnotationKey: kubelettypes.FileSource,
		}
		pod.Spec.PriorityClassName = priorityClassName
		if !kubelettypes.IsCriticalPod(pod) {
			framework.Failf("pod %q should be a critical pod", pod.Name)
		}
	} else {
		pod.Spec.PriorityClassName = priorityClassName
		if kubelettypes.IsCriticalPod(pod) {
			framework.Failf("pod %q should not be a critical pod", pod.Name)
		}
	}
	return pod
}

const (
	// https://github.com/kubernetes/kubernetes/blob/1dd781ddcad454cc381806fbc6bd5eba8fa368d7/pkg/kubelet/nodeshutdown/nodeshutdown_manager_linux.go#L43-L44
	podShutdownReason  = "Terminated"
	podShutdownMessage = "Pod was terminated in response to imminent node shutdown."
)

func isPodShutdown(pod *v1.Pod) bool {
	if pod == nil {
		return false
	}

	hasContainersNotReadyCondition := false
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.ContainersReady && cond.Status == v1.ConditionFalse {
			hasContainersNotReadyCondition = true
		}
	}

	return pod.Status.Message == podShutdownMessage && pod.Status.Reason == podShutdownReason && hasContainersNotReadyCondition && pod.Status.Phase == v1.PodFailed
}

// Pods should never report failed phase and have ready condition = true (https://github.com/kubernetes/kubernetes/issues/108594)
func isPodStatusAffectedByIssue108594(pod *v1.Pod) bool {
	return pod.Status.Phase == v1.PodFailed && podutils.IsPodReady(pod)
}

func isPodReadyToStartConditionSetToFalse(pod *v1.Pod) bool {
	if pod == nil {
		return false
	}
	readyToStartConditionSetToFalse := false
	for _, cond := range pod.Status.Conditions {
		if cond.Status == v1.ConditionFalse {
			readyToStartConditionSetToFalse = true
		}
	}

	return readyToStartConditionSetToFalse
}
