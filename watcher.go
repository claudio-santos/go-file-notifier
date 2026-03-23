package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileMonitor monitors files using simple polling
type FileMonitor struct {
	config   *Config
	state    *StateManager
	notifier *Notifier
	files    map[string]FileInfo
	mu       sync.Mutex
	stopChan chan struct{}
}

// FileInfo stores information about a file
type FileInfo struct {
	Size    int64
	ModTime time.Time
}

// NewFileMonitor creates a new monitor
func NewFileMonitor(config *Config, state *StateManager, notifier *Notifier) *FileMonitor {
	return &FileMonitor{
		config:   config,
		state:    state,
		notifier: notifier,
		files:    make(map[string]FileInfo),
		stopChan: make(chan struct{}),
	}
}

// Start begins monitoring
func (m *FileMonitor) Start() error {
	for _, filePath := range m.config.Files {
		info, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf("error accessing %q: %w", filePath, err)
		}

		// Check if changed since last run
		changed, _ := m.state.CheckForChanges(filePath, info.ModTime(), info.Size())
		if changed {
			fmt.Printf("File changed since last run: %s\n", filepath.Base(filePath))
			m.notifier.Notify(filePath, "File Changed")
		}

		// Save current state
		m.files[filePath] = FileInfo{
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}
		m.state.UpdateFileState(filePath, info.ModTime(), info.Size())
	}

	m.state.Save()

	fmt.Printf("Monitoring %d file(s), interval: %ds\n", len(m.config.Files), m.config.IntervalS)

	return nil
}

// Run executes the monitoring loop
func (m *FileMonitor) Run() error {
	ticker := time.NewTicker(time.Duration(m.config.IntervalS) * time.Second)
	defer ticker.Stop()

	lastNotified := make(map[string]time.Time)

	for {
		select {
		case <-m.stopChan:
			return nil
		case <-ticker.C:
			for _, filePath := range m.config.Files {
				info, err := os.Stat(filePath)
				if err != nil {
					fmt.Printf("Error: %s - %v\n", filepath.Base(filePath), err)
					continue
				}

				currentSize := info.Size()
				currentModTime := info.ModTime()

				m.mu.Lock()
				oldInfo, exists := m.files[filePath]
				m.mu.Unlock()

				// Check if changed
				if exists && (currentSize != oldInfo.Size || !currentModTime.Equal(oldInfo.ModTime)) {
					// Prevent duplicate notifications in short period
					if lastTime, ok := lastNotified[filePath]; ok && time.Since(lastTime) < time.Second {
						continue
					}

					fmt.Printf("File changed: %s\n", filepath.Base(filePath))
					m.notifier.Notify(filePath, "File Changed")

					// Update state
					m.mu.Lock()
					m.files[filePath] = FileInfo{
						Size:    currentSize,
						ModTime: currentModTime,
					}
					m.mu.Unlock()

					m.state.UpdateFileState(filePath, currentModTime, currentSize)
					m.state.Save()
					lastNotified[filePath] = time.Now()
				}
			}
		}
	}
}

// Stop stops monitoring
func (m *FileMonitor) Stop() error {
	close(m.stopChan)
	return m.state.Save()
}
