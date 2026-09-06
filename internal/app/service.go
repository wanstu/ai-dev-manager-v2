package app

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/environment"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/runtime"
	"ai-dev-manager-v2/internal/store"
	"ai-dev-manager-v2/internal/workspace"
)

type Service struct {
	Store        *store.Store
	Workspaces   *workspace.Service
	Environments *environment.Service
	MCPs         *catalog.Service
	Skills       *catalog.Service
	Memory       *memory.Service
}

type EnvironmentSummary struct {
	model.Environment
	PrivateMemoryCount int `json:"private_memory_count"`
}

type EnvironmentInspection struct {
	Environment        EnvironmentSummary   `json:"environment"`
	Workspace          model.Workspace      `json:"workspace"`
	Capabilities       []string             `json:"capabilities"`
	EnabledMCPs        []model.CatalogEntry `json:"enabled_mcps,omitempty"`
	EnabledSkills      []model.CatalogEntry `json:"enabled_skills,omitempty"`
	UnresolvedMCPIDs   []string             `json:"unresolved_mcp_ids,omitempty"`
	UnresolvedSkillIDs []string             `json:"unresolved_skill_ids,omitempty"`
}

func New(statePath string) *Service {
	s := store.New(statePath)
	ws := workspace.New(s)
	return &Service{
		Store:        s,
		Workspaces:   ws,
		Environments: environment.New(s, ws),
		MCPs:         catalog.New(s, catalog.KindMCP),
		Skills:       catalog.New(s, catalog.KindSkill),
		Memory:       memory.New(s),
	}
}

func (s *Service) EnvironmentSummaries() ([]EnvironmentSummary, error) {
	environments, err := s.Environments.List()
	if err != nil {
		return nil, err
	}
	items := make([]EnvironmentSummary, 0, len(environments))
	for _, env := range environments {
		items = append(items, environmentSummary(env))
	}
	return items, nil
}

func (s *Service) EnvironmentSummary(environmentID string) (EnvironmentSummary, error) {
	env, err := s.Environments.Get(environmentID)
	if err != nil {
		return EnvironmentSummary{}, err
	}
	return environmentSummary(env), nil
}

func (s *Service) InspectEnvironment(ctx context.Context, environmentID string) (EnvironmentInspection, error) {
	env, err := s.Environments.Get(environmentID)
	if err != nil {
		return EnvironmentInspection{}, err
	}
	ws, err := s.Workspaces.Get(env.WorkspaceID)
	if err != nil {
		return EnvironmentInspection{}, err
	}
	capabilities, err := s.Capabilities(ctx, env.ID)
	if err != nil {
		return EnvironmentInspection{}, err
	}
	mcps, err := s.MCPs.List()
	if err != nil {
		return EnvironmentInspection{}, err
	}
	skills, err := s.Skills.List()
	if err != nil {
		return EnvironmentInspection{}, err
	}
	enabledMCPs, unresolvedMCPs := resolveCatalogSelections(env.EnabledMCPIDs, mcps)
	enabledSkills, unresolvedSkills := resolveCatalogSelections(env.EnabledSkillIDs, skills)
	return EnvironmentInspection{
		Environment:        environmentSummary(env),
		Workspace:          ws,
		Capabilities:       capabilities,
		EnabledMCPs:        enabledMCPs,
		EnabledSkills:      enabledSkills,
		UnresolvedMCPIDs:   unresolvedMCPs,
		UnresolvedSkillIDs: unresolvedSkills,
	}, nil
}

func environmentSummary(env model.Environment) EnvironmentSummary {
	count := len(env.PrivateMemory)
	env.PrivateMemory = nil
	return EnvironmentSummary{Environment: env, PrivateMemoryCount: count}
}

func resolveCatalogSelections(ids []string, entries []model.CatalogEntry) ([]model.CatalogEntry, []string) {
	byID := make(map[string]model.CatalogEntry, len(entries))
	for _, entry := range entries {
		byID[entry.ID] = entry
	}
	resolved := make([]model.CatalogEntry, 0, len(ids))
	unresolved := make([]string, 0)
	for _, id := range ids {
		if entry, ok := byID[id]; ok {
			resolved = append(resolved, entry)
			continue
		}
		unresolved = append(unresolved, id)
	}
	return resolved, unresolved
}

func (s *Service) SetEnvironmentMCP(environmentID, mcpID string, enabled bool) (model.Environment, error) {
	exists, err := s.MCPs.Exists(mcpID)
	if err != nil {
		return model.Environment{}, err
	}
	if !exists {
		return model.Environment{}, fmt.Errorf("mcp %q not found", mcpID)
	}
	return s.Environments.SetMCP(environmentID, mcpID, enabled)
}

func (s *Service) SetEnvironmentSkill(environmentID, skillID string, enabled bool) (model.Environment, error) {
	exists, err := s.Skills.Exists(skillID)
	if err != nil {
		return model.Environment{}, err
	}
	if !exists {
		return model.Environment{}, fmt.Errorf("skill %q not found", skillID)
	}
	return s.Environments.SetSkill(environmentID, skillID, enabled)
}

func (s *Service) Runtime(environmentID string) (*runtime.Runtime, model.Environment, error) {
	env, err := s.Environments.Get(environmentID)
	if err != nil {
		return nil, model.Environment{}, err
	}
	state, err := s.Store.Load()
	if err != nil {
		return nil, model.Environment{}, err
	}
	rt, err := runtime.New(env.Root, state.AllowedExecutables)
	if err != nil {
		return nil, model.Environment{}, err
	}
	return rt, env, nil
}

func (s *Service) Capabilities(ctx context.Context, environmentID string) ([]string, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.Capabilities(ctx), nil
}

func (s *Service) AllowExecutable(executable string) error {
	executable = strings.TrimSpace(executable)
	if executable == "" {
		return fmt.Errorf("executable is required")
	}
	if filepath.IsAbs(executable) {
		executable = filepath.Clean(executable)
	}
	return s.Store.Update(func(state *model.State) error {
		for _, current := range state.AllowedExecutables {
			if strings.EqualFold(current, executable) {
				return nil
			}
		}
		state.AllowedExecutables = append(state.AllowedExecutables, executable)
		sort.Slice(state.AllowedExecutables, func(i, j int) bool {
			return strings.ToLower(state.AllowedExecutables[i]) < strings.ToLower(state.AllowedExecutables[j])
		})
		return nil
	})
}

func (s *Service) RemoveAllowedExecutable(executable string) error {
	executable = strings.TrimSpace(executable)
	if executable == "" {
		return fmt.Errorf("executable is required")
	}
	if filepath.IsAbs(executable) {
		executable = filepath.Clean(executable)
	}
	return s.Store.Update(func(state *model.State) error {
		for i, current := range state.AllowedExecutables {
			if strings.EqualFold(current, executable) {
				state.AllowedExecutables = append(state.AllowedExecutables[:i], state.AllowedExecutables[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("executable %q is not allowlisted", executable)
	})
}

func (s *Service) AllowedExecutables() ([]string, error) {
	state, err := s.Store.Load()
	if err != nil {
		return nil, err
	}
	return append([]string(nil), state.AllowedExecutables...), nil
}

func (s *Service) Tree(environmentID, path string, maxDepth, maxEntries int) (any, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.Tree(path, maxDepth, maxEntries)
}

func (s *Service) Read(environmentID, path string, maxBytes int) (any, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.Read(path, maxBytes)
}

func (s *Service) Search(environmentID, path, query string, maxFiles, maxMatches, maxBytesPerFile int) (any, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.Search(path, query, maxFiles, maxMatches, maxBytesPerFile)
}

func (s *Service) Write(environmentID, owner, path, content string, createParents bool) (any, error) {
	if _, err := s.Environments.RequireWriter(environmentID, owner); err != nil {
		return nil, err
	}
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	result, err := rt.Write(path, content, createParents)
	if err != nil {
		return nil, err
	}
	if err := s.Environments.Touch(environmentID, owner); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Edit(environmentID, owner, path, oldText, newText string, expected int) (any, error) {
	if _, err := s.Environments.RequireWriter(environmentID, owner); err != nil {
		return nil, err
	}
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	result, err := rt.Edit(path, oldText, newText, expected)
	if err != nil {
		return nil, err
	}
	if err := s.Environments.Touch(environmentID, owner); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Delete(environmentID, owner, path string) (any, error) {
	if _, err := s.Environments.RequireWriter(environmentID, owner); err != nil {
		return nil, err
	}
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	result, err := rt.Delete(path)
	if err != nil {
		return nil, err
	}
	if err := s.Environments.Touch(environmentID, owner); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) Exec(ctx context.Context, environmentID, owner, executable string, args []string, cwd string, timeoutMS int64, maxOutputBytes int) (any, error) {
	if _, err := s.Environments.RequireWriter(environmentID, owner); err != nil {
		return nil, err
	}
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}

	commandCtx, cancel := context.WithCancel(ctx)
	heartbeatDone := make(chan error, 1)
	go func() {
		interval := s.Environments.WriterLeaseTTL() / 3
		if interval <= 0 {
			interval = time.Second
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-commandCtx.Done():
				heartbeatDone <- nil
				return
			case <-ticker.C:
				if _, heartbeatErr := s.Environments.HeartbeatWriter(environmentID, owner); heartbeatErr != nil {
					heartbeatDone <- heartbeatErr
					cancel()
					return
				}
			}
		}
	}()

	result, execErr := rt.Exec(commandCtx, executable, args, cwd, timeoutMS, maxOutputBytes)
	cancel()
	heartbeatErr := <-heartbeatDone
	if heartbeatErr != nil {
		return nil, fmt.Errorf("writer heartbeat failed: %w", heartbeatErr)
	}
	if execErr != nil {
		return nil, execErr
	}
	if err := s.Environments.Touch(environmentID, owner); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) GitStatus(ctx context.Context, environmentID string) (any, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.GitStatus(ctx)
}

func (s *Service) GitDiff(ctx context.Context, environmentID string) (any, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.GitDiff(ctx)
}

func (s *Service) GitBranch(ctx context.Context, environmentID string) (any, error) {
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.GitBranch(ctx)
}
