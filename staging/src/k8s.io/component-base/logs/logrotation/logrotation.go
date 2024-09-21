package logrotation

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotationFile struct {
	// required, the max size of the log file in MB, 0 means no rotation
	maxSize uint
	// required, the max age of the log file in days, 0 means no cleanup
	maxAge        uint
	filePath      string
	mut           sync.Mutex
	file          *os.File
	currentSize   int64
	isFlushed     bool
	rotationFiles []string
	stop          chan struct{}
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
		if (w.currentSize / 1000000) >= int64(w.maxSize) {
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

	rotateFilePath := pathWithoutExt + "-" + time.Now().Format("20060102-150405") + ext

	err = w.file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(w.filePath, rotateFilePath)
	if err != nil {
		return err
	}

	w.rotationFiles = append(w.rotationFiles, rotateFilePath)

	w.file, err = os.OpenFile(w.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	w.currentSize = 0

	if w.maxAge > 0 {
		go func() {
			err = w.clean()
		}()
	}

	return
}

func (w *RotationFile) clean() (err error) {
	ageTime := time.Now().AddDate(0, 0, -int(w.maxAge))

	j := 0
	for i := 0; i < len(w.rotationFiles); i++ {
		info, err := os.Stat(w.rotationFiles[i])
		if err != nil {
			return err
		}

		if ageTime.After(info.ModTime()) {
			err = os.Remove(w.rotationFiles[i])
			if err != nil {
				return err
			}
			w.rotationFiles[j] = w.rotationFiles[i]
			j++
		}
	}

	w.rotationFiles = w.rotationFiles[:j]

	return
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
