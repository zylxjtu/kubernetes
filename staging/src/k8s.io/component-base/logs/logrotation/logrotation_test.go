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

package logrotation

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogrotationWrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "logrotation_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to remove temp directory: %v", err)
		}
	})

	logFilePath := filepath.Join(tmpDir, "test.log")
	maxSize := uint(1) // 1 MB
	maxAge := uint(1)  // 1 day

	rotationFile, err := Open(logFilePath, maxSize, maxAge)
	if err != nil {
		t.Fatalf("Failed to open RotationFile: %v", err)
	}

	t.Cleanup(func() {
		if err := rotationFile.Close(); err != nil {
			t.Errorf("Failed to close rotationFile: %v", err)
		}
	})

	testData := []byte("This is a test log entry.")
	n, err := rotationFile.Write(testData)
	if err != nil {
		t.Fatalf("Failed to write to RotationFile: %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, but wrote %d bytes", len(testData), n)
	}

	// Check if data is written to the file
	content, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if string(content) != string(testData) {
		t.Errorf("Expected log file content to be %q, but got %q", string(testData), string(content))
	}

	// Write more data to trigger rotation
	largeData := make([]byte, 1024*1024) // 1 MB
	n, err = rotationFile.Write(largeData)
	if err != nil {
		t.Fatalf("Failed to write large data to RotationFile: %v", err)
	}
	if n != len(largeData) {
		t.Errorf("Expected to write %d bytes, but wrote %d bytes", len(largeData), n)
	}

	// Check if rotation happened
	rotatedFiles, err := filepath.Glob(filepath.Join(tmpDir, "test-*.log"))
	if err != nil {
		t.Fatalf("Failed to list rotated files: %v", err)
	}
	if len(rotatedFiles) != 1 {
		t.Errorf("Expected rotated log files, but found none")
	}

	// Check if new log file is created
	newContent, err := os.ReadFile(logFilePath)
	if err != nil {
		t.Fatalf("Failed to read new log file: %v", err)
	}
	// Rotation after the write, new log file should be empty
	if len(newContent) != 0 {
		t.Errorf("Expected new log file content to be %d bytes, but got %d bytes", len(largeData), len(newContent))
	}

	// Check if rotated file was able to be cleaned up
	err = os.Chtimes(rotatedFiles[0], time.Now(), time.Now().AddDate(0, 0, -10))
	if err != nil {
		t.Fatalf("Failed to change access time of rotated file: %v", err)
	}

	// Trigger another rotation and check if the old rotated file has been cleaned up
	time.Sleep(1 * time.Second)
	largeData2 := make([]byte, 1024*1025) // 1 MB
	_, err = rotationFile.Write(largeData2)
	if err != nil {
		t.Fatalf("Failed to write to RotationFile: %v", err)
	}
	rotatedFiles2, err := filepath.Glob(filepath.Join(tmpDir, "test-*.log"))
	if err != nil {
		t.Fatalf("Failed to list rotated files: %v", err)
	}
	if len(rotatedFiles2) != 1 {
		t.Errorf("Expected rotated log files to be 1, but found %d", len(rotatedFiles))
	}
	if rotatedFiles[0] == rotatedFiles2[0] {
		t.Errorf("Expected rotated log files to be different, but found same")
	}
}
