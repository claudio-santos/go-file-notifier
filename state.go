package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// FileState represents the state of a file
type FileState struct {
	LastModified time.Time `json:"last_modified"`
	Size         int64     `json:"size"`
}

// State represents the persistent state of the monitor
type State struct {
	Files map[string]FileState `json:"files"`
}

// StateManager manages persistent state
type StateManager struct {
	path  string
	state *State
}

// NewStateManager creates a new state manager
func NewStateManager(path string) *StateManager {
	return &StateManager{
		path: path,
		state: &State{
			Files: make(map[string]FileState),
		},
	}
}

// Load reads the state from file
func (sm *StateManager) Load() error {
	data, err := os.ReadFile(sm.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading state.json: %w", err)
	}

	if err := json.Unmarshal(data, &sm.state); err != nil {
		return fmt.Errorf("error parsing state.json: %w", err)
	}

	return nil
}

// Save saves the state to file
func (sm *StateManager) Save() error {
	data, err := json.MarshalIndent(sm.state, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing state: %w", err)
	}

	if err := os.WriteFile(sm.path, data, 0644); err != nil {
		return fmt.Errorf("error saving state.json: %w", err)
	}

	return nil
}

// UpdateFileState updates the state of a file
func (sm *StateManager) UpdateFileState(path string, modTime time.Time, size int64) {
	state, exists := sm.state.Files[path]
	if !exists {
		state = FileState{}
	}

	state.LastModified = modTime
	state.Size = size

	sm.state.Files[path] = state
}

// CheckForChanges checks if a file has changed since last check
func (sm *StateManager) CheckForChanges(path string, modTime time.Time, size int64) (bool, error) {
	storedState, exists := sm.state.Files[path]
	if !exists {
		return false, nil
	}

	if !modTime.Equal(storedState.LastModified) || size != storedState.Size {
		return true, nil
	}

	return false, nil
}
