//go:build !windows
// +build !windows

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

package services

import (
	"os/exec"
	"strings"

	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
)

var containerRuntimeUnitName = ""

// updateKubeletConfig will update platfrom specific kubelet configuration.
func adjustPlatformSpecificKubeletConfig(kc *kubeletconfig.KubeletConfiguration) (*kubeletconfig.KubeletConfiguration, error) {
	return kc, nil
}

// updateCmdArgs will update platform specific command arguments.
func adjustPlatformSpecificKubeletArgs(cmdArgs []string, isSystemd bool, killCommand *exec.Cmd, restartCommand *exec.Cmd) ([]string, *exec.Cmd, *exec.Cmd, error) {
	// Adjust the args if we are running kubelet with systemd.
	if isSystemd {
		adjustArgsForSystemd(cmdArgs)
	}

	return cmdArgs, killCommand, restartCommand, nil
}

// adjustArgsForSystemd escape special characters in kubelet arguments for systemd. Systemd
// may try to do auto expansion without escaping.
func adjustArgsForSystemd(args []string) {
	for i := range args {
		args[i] = strings.Replace(args[i], "%", "%%", -1)
		args[i] = strings.Replace(args[i], "$", "$$", -1)
	}
}
