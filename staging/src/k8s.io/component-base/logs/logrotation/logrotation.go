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
	"strings"
	"sync"
	"time"
)

const timeLayout string = "20060102-150405"

type RotationFile struct {
	// required, the max size of the log file in MB, 0 means no rotation
	maxSize uint
	// required, the max age of the log file in days, 0 means no cleanup
	maxAge      uint
	filePath    string
	mut         sync.Mutex
	file        *os.File
	currentSize int64
	isFlushed   bool
	stop        chan struct{}
}

func Open(filePath string, maxSize uint, maxAge uint) (w *RotationFile, err error) {
	w = &RotationFile{
		filePath: filePath,
		maxSize:  maxSize,
		maxAge:   maxAge,
	}

	logFile, err := os.OpenFile(w.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	w.file = logFile

	w.isFlushed = true
	interval := time.NewTicker(time.Second * 5)
	w.stop = make(chan struct{})

	go func() {
		defer interval.Stop()
		for {
			select {
			case <-interval.C:
				// This block runs every time the ticker ticks
				w.isFlushed = false
			case <-w.stop:
				// Stop signal received, exiting the goroutine
				return
			}
		}
	}()

	if w.maxSize > 0 {
		info, err := os.Stat(w.filePath)
		if err != nil {
			return nil, err
		}
		w.currentSize = info.Size()
	}

	return w, nil
}

// write func to satisfy io.writer interface
func (w *RotationFile) Write(p []byte) (n int, err error) {
	w.mut.Lock()
	defer w.mut.Unlock()

	n, err = w.file.Write(p)
	if err != nil {
		return 0, err
	}

	if !w.isFlushed {
		err = w.file.Sync()
		if err != nil {
			return 0, err
		}
		w.isFlushed = true
	}

	if w.maxSize > 0 {
		w.currentSize += int64(len(p))

		// if file size over maxsize rotate the log file
		if w.currentSize >= int64(w.maxSize)*1024*1024 {
			// Explicitly call file.Sync() to ensure data is written to disk
			err = w.file.Sync()
			if err != nil {
				return 0, err
			}
			err = w.rotate()
			if err != nil {
				return 0, err
			}
		}
	}

	return n, nil
}

func (w *RotationFile) rotate() (err error) {
	// Get the file extension
	ext := filepath.Ext(w.filePath)

	// Remove the extension from the filename
	pathWithoutExt := strings.TrimSuffix(w.filePath, ext)

	rotateFilePath := pathWithoutExt + "-" + time.Now().Format(timeLayout) + ext

	err = w.file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(w.filePath, rotateFilePath)
	if err != nil {
		return err
	}

	w.file, err = os.OpenFile(w.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.currentSize = 0

	if w.maxAge > 0 {
		go func() {
			err = w.clean(pathWithoutExt, ext)
		}()
	}

	return
}

// Clean up the old log files in the format of
// <basename>-<timestamp><ext>
// This should be safe enough to avoid false deletion
// This will work for multiple restarts of the same program
func (w *RotationFile) clean(pathWithoutExt string, ext string) (err error) {
	ageTime := time.Now().AddDate(0, 0, -int(w.maxAge))

	directory := filepath.Dir(pathWithoutExt)
	basename := filepath.Base(pathWithoutExt) + "-"

	dir, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	err = nil
	for _, v := range dir {
		if strings.HasPrefix(v.Name(), basename) && strings.HasSuffix(v.Name(), ext) {
			// Remove the prefix and suffix
			trimmed := strings.TrimPrefix(v.Name(), basename)
			trimmed = strings.TrimSuffix(trimmed, ext)

			_, err = time.Parse(timeLayout, trimmed)
			if err == nil {
				info, errInfo := v.Info()
				if errInfo != nil {
					err = errInfo
					// Ignore the error while continue with the next clenup
					continue
				}

				if ageTime.After(info.ModTime()) {
					err = os.Remove(filepath.Join(directory, v.Name()))
					if err != nil {
						// Ignore the error while continue with the next clenup
						continue
					}
				}
			}

		}
	}

	return err
}

func (w *RotationFile) Close() (err error) {
	w.mut.Lock()
	defer w.mut.Unlock()

	// Explicitly call file.Sync() to ensure data is written to disk
	err = w.file.Sync()
	if err != nil {
		return err
	}

	err = w.file.Close()
	if err != nil {
		return err
	}

	close(w.stop)

	return
}
