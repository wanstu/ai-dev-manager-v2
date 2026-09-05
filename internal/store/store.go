package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"ai-dev-manager-v2/internal/model"
)

const stateVersion = 1

type Store struct {
	mu   sync.Mutex
	path string
}

func DefaultPath() (string, error) {
	if root := os.Getenv("ADM_V2_HOME"); root != "" {
		return filepath.Join(root, "state.json"), nil
	}
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "ai-dev-manager-v2", "state.json"), nil
}

func New(path string) *Store { return &Store{path: path} }

func (s *Store) Path() string { return s.path }

func (s *Store) Load() (model.State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked()
}

func (s *Store) Update(fn func(*model.State) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.loadLocked()
	if err != nil {
		return err
	}
	if err := fn(&state); err != nil {
		return err
	}
	return s.saveLocked(state)
}

func (s *Store) loadLocked() (model.State, error) {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultState(), nil
	}
	if err != nil {
		return model.State{}, err
	}
	var state model.State
	if err := json.Unmarshal(data, &state); err != nil {
		return model.State{}, fmt.Errorf("decode state: %w", err)
	}
	if state.Version != stateVersion {
		return model.State{}, fmt.Errorf("unsupported state version %d", state.Version)
	}
	if state.GlobalMemory == nil {
		state.GlobalMemory = map[string]string{}
	}
	return state, nil
}

func (s *Store) saveLocked(state model.State) error {
	state.Version = stateVersion
	if state.GlobalMemory == nil {
		state.GlobalMemory = map[string]string{}
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmp, s.path); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func defaultState() model.State {
	return model.State{Version: stateVersion, GlobalMemory: map[string]string{}}
}
