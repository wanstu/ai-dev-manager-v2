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
	"ai-dev-manager-v2/internal/pathutil"
	"ai-dev-manager-v2/internal/store"
	"ai-dev-manager-v2/internal/workspace"
)

const (
	StateReady            = "ready"
	DefaultWriterLeaseTTL = 5 * time.Minute
)

type Service struct {
	store          *store.Store
	workspaces     *workspace.Service
	now            func() time.Time
	writerLeaseTTL time.Duration
}

func New(s *store.Store, workspaces *workspace.Service) *Service {
	return &Service{
		store:          s,
		workspaces:     workspaces,
		now:            time.Now,
		writerLeaseTTL: DefaultWriterLeaseTTL,
	}
}

func (s *Service) WriterLeaseTTL() time.Duration { return s.leaseTTL() }

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
		now := s.nowUTC()
		for _, existing := range state.Environments {
			if existing.WorkspaceID == ws.ID && strings.EqualFold(existing.Name, name) && samePath(existing.Root, root) {
				result = s.environmentView(existing, now)
				return nil
			}
		}
		created, createErr := s.newEnvironment(state, ws.ID, name, root, now)
		if createErr != nil {
			return createErr
		}
		result = created
		state.Environments = append(state.Environments, result)
		return nil
	})
	return result, err
}

func (s *Service) CreateManaged(workspaceID, name, root string, managed model.ManagedWorktree) (model.Environment, model.ManagedWorktree, error) {
	ws, err := s.workspaces.Get(strings.TrimSpace(workspaceID))
	if err != nil {
		return model.Environment{}, model.ManagedWorktree{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Environment{}, model.ManagedWorktree{}, fmt.Errorf("environment name is required")
	}
	root, err = canonicalDir(root)
	if err != nil {
		return model.Environment{}, model.ManagedWorktree{}, err
	}
	managed.ID = strings.TrimSpace(managed.ID)
	if managed.ID == "" {
		return model.Environment{}, model.ManagedWorktree{}, fmt.Errorf("managed worktree id is required")
	}
	if strings.TrimSpace(managed.Branch) == "" || strings.TrimSpace(managed.BaseCommit) == "" || strings.TrimSpace(managed.GitCommonDir) == "" {
		return model.Environment{}, model.ManagedWorktree{}, fmt.Errorf("managed worktree git identity is incomplete")
	}

	var result model.Environment
	var managedResult model.ManagedWorktree
	err = s.store.Update(func(state *model.State) error {
		for _, existing := range state.ManagedWorktrees {
			if existing.ID == managed.ID || samePath(existing.Root, root) {
				return fmt.Errorf("managed worktree %q already exists", managed.ID)
			}
		}
		now := s.nowUTC()
		created, createErr := s.newEnvironment(state, ws.ID, name, root, now)
		if createErr != nil {
			return createErr
		}
		managed.EnvironmentID = created.ID
		managed.WorkspaceID = ws.ID
		managed.Root = root
		managed.CreatedAt = now
		result = created
		managedResult = managed
		state.Environments = append(state.Environments, result)
		state.ManagedWorktrees = append(state.ManagedWorktrees, managedResult)
		return nil
	})
	return result, managedResult, err
}

func (s *Service) newEnvironment(state *model.State, workspaceID, name, root string, now time.Time) (model.Environment, error) {
	id, err := identity.New("env")
	if err != nil {
		return model.Environment{}, err
	}
	result := model.Environment{
		ID:             id,
		WorkspaceID:    workspaceID,
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
	return result, nil
}

func (s *Service) List() ([]model.Environment, error) {
	state, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	items := make([]model.Environment, len(state.Environments))
	now := s.nowUTC()
	for i, env := range state.Environments {
		items[i] = s.environmentView(env, now)
	}
	return items, nil
}

func (s *Service) Get(id string) (model.Environment, error) {
	state, err := s.store.Load()
	if err != nil {
		return model.Environment{}, err
	}
	now := s.nowUTC()
	for _, env := range state.Environments {
		if env.ID == id {
			return s.environmentView(env, now), nil
		}
	}
	return model.Environment{}, fmt.Errorf("environment %q not found", id)
}

func (s *Service) Rename(id, name string) (model.Environment, error) {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if name == "" {
		return model.Environment{}, fmt.Errorf("environment name is required")
	}
	var result model.Environment
	err := s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		env := &state.Environments[idx]
		env.Name = name
		env.UpdatedAt = s.nowUTC()
		result = *env
		return nil
	})
	return result, err
}

func (s *Service) Remove(id string) (model.Environment, error) {
	id = strings.TrimSpace(id)
	var removed model.Environment
	err := s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}

		now := s.nowUTC()
		target := state.Environments[idx]
		for _, managed := range state.ManagedWorktrees {
			if managed.EnvironmentID == target.ID {
				return fmt.Errorf("environment %s is backed by managed worktree %s; use managed worktree destroy", target.ID, managed.ID)
			}
		}
		if target.State != StateReady {
			return fmt.Errorf("environment %s cannot be removed while state is %q", target.ID, target.State)
		}
		if target.Writer != nil && !s.writerExpired(target.Writer, now) {
			return fmt.Errorf("environment %s cannot be removed while writer %q is active", target.ID, target.Writer.Owner)
		}
		target.Writer = nil
		removed = target
		state.Environments = append(state.Environments[:idx], state.Environments[idx+1:]...)
		return nil
	})
	return removed, err
}

func (s *Service) RemoveManaged(id string) (model.Environment, model.ManagedWorktree, error) {
	id = strings.TrimSpace(id)
	var removed model.Environment
	var managedRemoved model.ManagedWorktree
	err := s.store.Update(func(state *model.State) error {
		envIdx := findEnvironment(state.Environments, id)
		if envIdx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		managedIdx := findManagedWorktreeByEnvironment(state.ManagedWorktrees, id)
		if managedIdx < 0 {
			return fmt.Errorf("environment %s is not backed by a managed worktree", id)
		}
		target := state.Environments[envIdx]
		if target.State != StateReady {
			return fmt.Errorf("environment %s cannot be removed while state is %q", target.ID, target.State)
		}
		removed = target
		removed.Writer = nil
		managedRemoved = state.ManagedWorktrees[managedIdx]
		state.Environments = append(state.Environments[:envIdx], state.Environments[envIdx+1:]...)
		state.ManagedWorktrees = append(state.ManagedWorktrees[:managedIdx], state.ManagedWorktrees[managedIdx+1:]...)
		return nil
	})
	return removed, managedRemoved, err
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
		env.UpdatedAt = s.nowUTC()
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
		now := s.nowUTC()
		for i := range state.Environments {
			other := &state.Environments[i]
			if other.Writer == nil || !samePath(other.Root, target.Root) {
				continue
			}
			if s.writerExpired(other.Writer, now) {
				other.Writer = nil
				other.UpdatedAt = now
				continue
			}
			if other.ID == target.ID && other.Writer.Owner == owner {
				s.renewWriter(other, now, false)
				result = *other
				return nil
			}
			return fmt.Errorf("physical root already has writer %q through environment %s", other.Writer.Owner, other.ID)
		}
		target.Writer = &model.WriterLease{Owner: owner, AcquiredAt: now}
		s.renewWriter(target, now, false)
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
		now := s.nowUTC()
		if s.writerExpired(target.Writer, now) {
			target.Writer = nil
			target.UpdatedAt = now
			result = *target
			return nil
		}
		if !force && target.Writer.Owner != strings.TrimSpace(owner) {
			return fmt.Errorf("writer is held by %q", target.Writer.Owner)
		}
		target.Writer = nil
		target.UpdatedAt = now
		result = *target
		return nil
	})
	return result, err
}

func (s *Service) RequireWriter(id, owner string) (model.Environment, error) {
	owner = strings.TrimSpace(owner)
	var result model.Environment
	err := s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		env := &state.Environments[idx]
		now := s.nowUTC()
		if env.Writer == nil || s.writerExpired(env.Writer, now) || env.Writer.Owner != owner {
			return fmt.Errorf("environment %s is not owned by writer %q", id, owner)
		}
		s.renewWriter(env, now, false)
		result = *env
		return nil
	})
	return result, err
}

func (s *Service) HeartbeatWriter(id, owner string) (model.Environment, error) {
	return s.RequireWriter(id, owner)
}

func (s *Service) Touch(id, owner string) error {
	return s.store.Update(func(state *model.State) error {
		idx := findEnvironment(state.Environments, id)
		if idx < 0 {
			return fmt.Errorf("environment %q not found", id)
		}
		env := &state.Environments[idx]
		now := s.nowUTC()
		if env.Writer == nil || s.writerExpired(env.Writer, now) || env.Writer.Owner != strings.TrimSpace(owner) {
			return fmt.Errorf("environment %s is not owned by writer %q", id, owner)
		}
		s.renewWriter(env, now, true)
		return nil
	})
}

func (s *Service) nowUTC() time.Time {
	if s.now == nil {
		return time.Now().UTC()
	}
	return s.now().UTC()
}

func (s *Service) writerExpired(lease *model.WriterLease, now time.Time) bool {
	if lease == nil {
		return false
	}
	if lease.ExpiresAt.IsZero() {
		return true
	}
	return !now.Before(lease.ExpiresAt)
}

func (s *Service) environmentView(env model.Environment, now time.Time) model.Environment {
	if env.Writer == nil {
		return env
	}
	lease := *env.Writer
	if s.writerExpired(&lease, now) {
		env.Writer = nil
		return env
	}
	env.Writer = &lease
	return env
}

func (s *Service) renewWriter(env *model.Environment, now time.Time, activity bool) {
	env.Writer.LastSeenAt = now
	env.Writer.ExpiresAt = now.Add(s.leaseTTL())
	env.UpdatedAt = now
	if activity {
		env.LastActivityAt = now
	}
}

func (s *Service) leaseTTL() time.Duration {
	if s.writerLeaseTTL <= 0 {
		return DefaultWriterLeaseTTL
	}
	return s.writerLeaseTTL
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
	rel, err := filepath.Rel(pathutil.ForCompare(base), pathutil.ForCompare(target))
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

func findManagedWorktreeByEnvironment(values []model.ManagedWorktree, environmentID string) int {
	for i := range values {
		if values[i].EnvironmentID == environmentID {
			return i
		}
	}
	return -1
}

func samePath(a, b string) bool { return pathutil.Same(a, b) }
