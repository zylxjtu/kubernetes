//go:build !linux && !windows
// +build !linux,!windows

/*
Copyright 2020 The Kubernetes Authors.

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

package utils

import (
	"context"

	"k8s.io/kubernetes/test/e2e/framework"
)

// IsCgroup2UnifiedMode returns whether we are running in cgroup v2 unified mode.
func IsCgroup2UnifiedMode() bool {
	return false
}

func RestartKubelet(ctx context.Context, running bool) {
}

// mustStopKubelet will kill the running kubelet, and returns a func that will restart the process again
func MustStopKubelet(ctx context.Context, f *framework.Framework) func(ctx context.Context) {
	return func(ctx context.Context) {}
}

func StopContainerRuntime() error {
	return nil
}

func StartContainerRuntime() error {
	return nil
}
