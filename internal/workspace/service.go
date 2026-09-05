package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
)

type Service struct{ store *store.Store }

func New(s *store.Store) *Service { return &Service{store: s} }

func (s *Service) Add(path, name string) (model.Workspace, error) {
	root, err := canonicalDir(path)
	if err != nil {
		return model.Workspace{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = filepath.Base(root)
		if name == "." || name == string(filepath.Separator) || name == "" {
			name = root
		}
	}
	var result model.Workspace
	err = s.store.Update(func(state *model.State) error {
		for _, ws := range state.Workspaces {
			if samePath(ws.Path, root) {
				return fmt.Errorf("workspace path already registered as %s", ws.ID)
			}
		}
		id, err := identity.New("ws")
		if err != nil {
			return err
		}
		result = model.Workspace{ID: id, Name: name, Path: root, CreatedAt: time.Now().UTC()}
		state.Workspaces = append(state.Workspaces, result)
		return nil
	})
	return result, err
}

func (s *Service) List() ([]model.Workspace, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	return append([]model.Workspace(nil), state.Workspaces...), nil
}

func (s *Service) Get(id string) (model.Workspace, error) {
	state, err := s.store.Load()
	if err != nil {
		return model.Workspace{}, err
	}
	for _, ws := range state.Workspaces {
		if ws.ID == id {
			return ws, nil
		}
	}
	return model.Workspace{}, fmt.Errorf("workspace %q not found", id)
}

func canonicalDir(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("workspace path is required")
	}
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("workspace path is not a directory: %s", abs)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func samePath(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}
