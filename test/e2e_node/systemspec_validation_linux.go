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

package e2enode

import (
	"encoding/json"
	"os"
	"syscall"

	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	klog "k8s.io/klog/v2"
	"k8s.io/kubernetes/test/e2e/framework"
	system "k8s.io/system-validators/validators"
)

// When running the containerized conformance test, we'll mount the
// host root filesystem as readonly to /rootfs.
const rootfs = "/rootfs"

func systemValidation(systemSpecFile *string) {
	spec := &system.DefaultSysSpec
	if *systemSpecFile != "" {
		var err error
		spec, err = loadSystemSpecFromFile(*systemSpecFile)
		if err != nil {
			klog.Exitf("Failed to load system spec: %v", err)
		}
	}
	if framework.TestContext.NodeConformance {
		// Chroot to /rootfs to make system validation can check system
		// as in the root filesystem.
		// TODO(random-liu): Consider to chroot the whole test process to make writing
		// // test easier.
		if err := syscall.Chroot(rootfs); err != nil {
			klog.Exitf("chroot %q failed: %v", rootfs, err)
		}
	}
	warns, errs := system.ValidateSpec(*spec, "remote")
	if len(warns) != 0 {
		klog.Warningf("system validation warns: %v", warns)
	}
	if len(errs) != 0 {
		klog.Exitf("system validation failed: %v", errs)
	}
	return
}

// loadSystemSpecFromFile returns the system spec from the file with the
// filename.
func loadSystemSpecFromFile(filename string) (*system.SysSpec, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data, err := utilyaml.ToJSON(b)
	if err != nil {
		return nil, err
	}
	spec := new(system.SysSpec)
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, err
	}
	return spec, nil
}
