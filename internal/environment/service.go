package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/identity"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/store"
	"ai-dev-manager-v2/internal/workspace"
)

const StateReady = "ready"

type Service struct {
	store      *store.Store
	workspaces *workspace.Service
}

func New(s *store.Store, workspaces *workspace.Service) *Service {
	return &Service{store: s, workspaces: workspaces}
}

func (s *Service) Create(workspaceID, name, root string) (model.Environment, error) {
	ws, err := s.workspaces.Get(strings.TrimSpace(workspaceID))
	if err != nil {
		return model.Environment{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Environment{}, fmt.Errorf("environment name is required")
	}
	if strings.TrimSpace(root) == "" {
		root = ws.Path
	}
	root, err = canonicalDir(root)
	if err != nil {
		return model.Environment{}, err
	}
	if !within(ws.Path, root) {
		return model.Environment{}, fmt.Errorf("environment root must stay inside workspace %s", ws.Path)
	}

	var result model.Environment
	err = s.store.Update(func(state *model.State) error {
		id, err := identity.New("env")
		if err != nil {
			return err
		}
		now := time.Now().UTC()
		result = model.Environment{
			ID:             id,
			WorkspaceID:    ws.ID,
			Name:           name,
			Root:           root,
			State:          StateReady,
			CreatedAt:      now,
			UpdatedAt:      now,
			LastActivityAt: now,
			PrivateMemory:  map[string]string{},
		}
		for _, entry := range state.MCPs {
			if entry.DefaultIncludeInEnv {
				result.EnabledMCPIDs = append(result.EnabledMCPIDs, entry.ID)
			}
		}
		for _, entry := range state.Skills {
			if entry.DefaultIncludeInEnv {
				result.EnabledSkillIDs = append(result.EnabledSkillIDs, entry.ID)
			}
		}
		sort.Strings(result.EnabledMCPIDs)
		sort.Strings(result.EnabledSkillIDs)
		state.Environments = append(state.Environments, result)
		return nil
	})
	return result, err
}

func (s *Service) List() ([]model.Environment, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	return append([]model.Environment(nil), state.Environments...), nil
}

func (s *Service) Get(id string) (model.Environment, error) {
	state, err := s.store.Load()
	if err != nil {
		return model.Environment{}, err
	}
	for _, env := range state.Environments {
		if env.ID == id {
			return env, nil
		}
	}
	return model.Environment{}, fmt.Errorf("environment %q not found", id)
}

func (s *Service) SetMCP(id, mcpID string, enabled bool) (model.Environment, error) {
	return s.setSelection(id, mcpID, enabled, true)
}

func (s *Service) SetSkill(id, skillID string, enabled bool) (model.Environment, error) {
	return s.setSelection(id, skillID, enabled, false)
}

func (s *Service) setSelection(id, value string, enabled, isMCP bool) (model.Environment, error) {
	var result model.Environment
	err := s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		env := &state.Environments[idx]
		values := &env.EnabledSkillIDs
		if isMCP {
			values = &env.EnabledMCPIDs
		}
		found := -1
		for i, current := range *values {
			if current == value {
				found = i
				break
			}
		}
		if enabled && found < 0 {
			*values = append(*values, value)
			sort.Strings(*values)
		}
		if !enabled && found >= 0 {
			*values = append((*values)[:found], (*values)[found+1:]...)
		}
		env.UpdatedAt = time.Now().UTC()
		result = *env
		return nil
	})
	return result, err
}

func (s *Service) AcquireWriter(id, owner string) (model.Environment, error) {
	owner = strings.TrimSpace(owner)
	if owner == "" {
		return model.Environment{}, fmt.Errorf("writer owner is required")
	}
	var result model.Environment
	err := s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		target := &state.Environments[idx]
		for i := range state.Environments {
			other := &state.Environments[i]
			if other.Writer == nil || !samePath(other.Root, target.Root) {
				continue
			}
			if other.ID == target.ID && other.Writer.Owner == owner {
				now := time.Now().UTC()
				other.Writer.LastSeenAt = now
				other.UpdatedAt = now
				result = *other
				return nil
			}
			return fmt.Errorf("physical root already has writer %q through environment %s", other.Writer.Owner, other.ID)
		}
		now := time.Now().UTC()
		target.Writer = &model.WriterLease{Owner: owner, AcquiredAt: now, LastSeenAt: now}
		target.UpdatedAt = now
		result = *target
		return nil
	})
	return result, err
}

func (s *Service) ReleaseWriter(id, owner string, force bool) (model.Environment, error) {
	var result model.Environment
	err := s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		target := &state.Environments[idx]
		if target.Writer == nil {
			result = *target
			return nil
		}
		if !force && target.Writer.Owner != strings.TrimSpace(owner) {
			return fmt.Errorf("writer is held by %q", target.Writer.Owner)
		}
		target.Writer = nil
		target.UpdatedAt = time.Now().UTC()
		result = *target
		return nil
	})
	return result, err
}

func (s *Service) RequireWriter(id, owner string) (model.Environment, error) {
	env, err := s.Get(id)
	if err != nil {
		return model.Environment{}, err
	}
	if env.Writer == nil || env.Writer.Owner != strings.TrimSpace(owner) {
		return model.Environment{}, fmt.Errorf("environment %s is not owned by writer %q", id, owner)
	}
	return env, nil
}

func (s *Service) Touch(id, owner string) error {
	return s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		env := &state.Environments[idx]
		if env.Writer == nil || env.Writer.Owner != strings.TrimSpace(owner) {
			return fmt.Errorf("environment %s is not owned by writer %q", id, owner)
		}
		now := time.Now().UTC()
		env.Writer.LastSeenAt = now
		env.UpdatedAt = now
		env.LastActivityAt = now
		return nil
	})
}

func canonicalDir(path string) (string, error) {
	abs, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("environment root is not a directory: %s", abs)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func within(base, target string) bool {
	rel, err := filepath.Rel(filepath.Clean(base), filepath.Clean(target))
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func findEnvironment(values []model.Environment, id string) int {
	for i := range values {
		if values[i].ID == id {
			return i
		}
	}
	return -1
}

func samePath(a, b string) bool { return strings.EqualFold(filepath.Clean(a), filepath.Clean(b)) }
