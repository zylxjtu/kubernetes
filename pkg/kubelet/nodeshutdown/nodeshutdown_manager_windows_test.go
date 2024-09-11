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

package nodeshutdown

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/record"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	"k8s.io/klog/v2/ktesting"
	_ "k8s.io/klog/v2/ktesting/init" // activate ktesting command line flags
	pkgfeatures "k8s.io/kubernetes/pkg/features"
	probetest "k8s.io/kubernetes/pkg/kubelet/prober/testing"
	"k8s.io/kubernetes/pkg/kubelet/volumemanager"
)

func makePod(name string, priority int32, terminationGracePeriod *int64) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			UID:  types.UID(name),
		},
		Spec: v1.PodSpec{
			Priority:                      &priority,
			TerminationGracePeriodSeconds: terminationGracePeriod,
		},
	}
}

func TestWindowsFeatureEnabled(t *testing.T) {
	var tests = []struct {
		desc                         string
		shutdownGracePeriodRequested time.Duration
		featureGateEnabled           bool
		expectEnabled                bool
	}{
		{
			desc:                         "shutdownGracePeriodRequested 0; disables feature",
			shutdownGracePeriodRequested: time.Duration(0 * time.Second),
			featureGateEnabled:           true,
			expectEnabled:                false,
		},
		{
			desc:                         "feature gate disabled; disables feature",
			shutdownGracePeriodRequested: time.Duration(100 * time.Second),
			featureGateEnabled:           false,
			expectEnabled:                false,
		},
		{
			desc:                         "feature gate enabled; shutdownGracePeriodRequested > 0; enables feature",
			shutdownGracePeriodRequested: time.Duration(100 * time.Second),
			featureGateEnabled:           true,
			expectEnabled:                true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			logger, _ := ktesting.NewTestContext(t)
			activePodsFunc := func() []*v1.Pod {
				return nil
			}
			killPodsFunc := func(pod *v1.Pod, evict bool, gracePeriodOverride *int64, fn func(*v1.PodStatus)) error {
				return nil
			}
			featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, pkgfeatures.WindowsGracefulNodeShutdown, tc.featureGateEnabled)

			proberManager := probetest.FakeManager{}
			fakeRecorder := &record.FakeRecorder{}
			fakeVolumeManager := volumemanager.NewFakeVolumeManager([]v1.UniqueVolumeName{}, 0, nil)
			nodeRef := &v1.ObjectReference{Kind: "Node", Name: "test", UID: types.UID("test"), Namespace: ""}

			manager, _ := NewManager(&Config{
				Logger:                          logger,
				ProbeManager:                    proberManager,
				VolumeManager:                   fakeVolumeManager,
				Recorder:                        fakeRecorder,
				NodeRef:                         nodeRef,
				GetPodsFunc:                     activePodsFunc,
				KillPodFunc:                     killPodsFunc,
				SyncNodeStatusFunc:              func() {},
				ShutdownGracePeriodRequested:    tc.shutdownGracePeriodRequested,
				ShutdownGracePeriodCriticalPods: 0,
				StateDirectory:                  os.TempDir(),
			})
			assert.Equal(t, tc.expectEnabled, manager != managerStub{})
		})
	}
}
