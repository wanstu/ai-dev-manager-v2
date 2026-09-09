package management

import (
	"context"
	"strings"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
)

type Snapshot struct {
	Workspaces         []model.Workspace        `json:"workspaces"`
	Environments       []app.EnvironmentSummary `json:"environments"`
	AllowedExecutables []string                 `json:"allowed_executables"`
	MCPs               []model.MCPDefinition    `json:"mcps"`
	Skills             []model.CatalogEntry     `json:"skills"`
	GlobalMemoryCount  int                      `json:"global_memory_count"`
}

type Service struct {
	app *app.Service
}

func New(application *app.Service) *Service {
	return &Service{app: application}
}

func (s *Service) WorkspaceInspect(id string) (model.Workspace, error) {
	return s.app.Workspaces.Get(id)
}

func (s *Service) EnvironmentInspect(id string) (app.EnvironmentInspection, error) {
	return s.app.InspectEnvironment(context.Background(), id)
}

func (s *Service) GlobalMemoryList() ([]memory.Entry, error) {
	return s.app.Memory.GlobalList()
}

func (s *Service) GlobalMemoryRead(key string) (memory.Entry, error) {
	return s.app.Memory.GlobalRead(key)
}

func (s *Service) EnvironmentMemoryList(environmentID string) ([]memory.Entry, error) {
	return s.app.Memory.EnvironmentList(environmentID)
}

func (s *Service) EnvironmentMemoryRead(environmentID, key string) (memory.Entry, error) {
	return s.app.Memory.EnvironmentRead(environmentID, key)
}

func (s *Service) WorkspaceAdd(path, name string) (model.Workspace, error) {
	return s.app.Workspaces.Add(path, name)
}

func (s *Service) WorkspaceRename(id, name string) (model.Workspace, error) {
	return s.app.Workspaces.Rename(id, name)
}

func (s *Service) WorkspaceRemove(id string) (model.Workspace, error) {
	return s.app.Workspaces.Remove(id)
}

func (s *Service) EnvironmentCreate(workspaceID, name, root string) (app.EnvironmentSummary, error) {
	env, err := s.app.Environments.Create(workspaceID, name, root)
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	return s.app.EnvironmentSummary(env.ID)
}

func (s *Service) EnvironmentRename(id, name string) (app.EnvironmentSummary, error) {
	env, err := s.app.Environments.Rename(id, name)
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	return s.app.EnvironmentSummary(env.ID)
}

func (s *Service) EnvironmentRemove(id string) error {
	_, err := s.app.Environments.Remove(id)
	return err
}

func (s *Service) ExecAllow(executable string) ([]string, error) {
	if err := s.app.AllowExecutable(executable); err != nil {
		return nil, err
	}
	return s.app.AllowedExecutables()
}

func (s *Service) ExecRemove(executable string) ([]string, error) {
	if err := s.app.RemoveAllowedExecutable(executable); err != nil {
		return nil, err
	}
	return s.app.AllowedExecutables()
}

func (s *Service) MCPAdd(name, endpoint string, defaultInclude bool) (model.MCPDefinition, error) {
	return s.app.MCPs.AddMCP(name, endpoint, defaultInclude)
}

func (s *Service) MCPAddConfig(name string, config catalog.MCPConfig) (model.MCPDefinition, error) {
	return s.app.MCPs.AddMCPConfig(name, config)
}

func (s *Service) MCPImportPreview(input app.MCPImportInput) (app.MCPImportPreview, error) {
	return s.app.PreviewMCPImport(input)
}

func (s *Service) MCPImportApply(input app.MCPImportInput) (app.MCPImportApplyResult, error) {
	return s.app.ApplyMCPImport(input)
}

func (s *Service) MCPRemove(id string) error {
	return s.app.MCPs.Remove(id)
}

func (s *Service) MCPSetDefault(id string, value bool) (model.MCPDefinition, error) {
	return s.app.MCPs.SetDefault(id, value)
}

func (s *Service) MCPHealth(ctx context.Context, environmentID, mcpID string) (app.MCPHealthStatus, error) {
	return s.app.ProbeMCPHealth(ctx, environmentID, mcpID)
}

func (s *Service) SkillAdd(root, supportRoot string, defaultInclude bool) ([]model.CatalogEntry, error) {
	supportRoots := []string{}
	if value := strings.TrimSpace(supportRoot); value != "" {
		supportRoots = append(supportRoots, value)
	}
	return s.app.Skills.AddSkillRoot(root, supportRoots, defaultInclude)
}

func (s *Service) SkillRemove(id string) error {
	return s.app.Skills.Remove(id)
}

func (s *Service) SkillSetDefault(id string, value bool) (model.CatalogEntry, error) {
	return s.app.Skills.SetDefault(id, value)
}

func (s *Service) EnvironmentMCPSet(environmentID, mcpID string, enabled bool) (app.EnvironmentSummary, error) {
	env, err := s.app.SetEnvironmentMCP(environmentID, mcpID, enabled)
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	return s.app.EnvironmentSummary(env.ID)
}

func (s *Service) EnvironmentSkillSet(environmentID, skillID string, enabled bool) (app.EnvironmentSummary, error) {
	env, err := s.app.SetEnvironmentSkill(environmentID, skillID, enabled)
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	return s.app.EnvironmentSummary(env.ID)
}

func (s *Service) GlobalMemoryWrite(key, value string) error {
	return s.app.Memory.GlobalWrite(key, value)
}

func (s *Service) GlobalMemoryDelete(key string) error {
	return s.app.Memory.GlobalDelete(key)
}

func (s *Service) EnvironmentMemoryWrite(environmentID, key, value string) error {
	return s.app.Memory.EnvironmentWrite(environmentID, key, value)
}

func (s *Service) EnvironmentMemoryDelete(environmentID, key string) error {
	return s.app.Memory.EnvironmentDelete(environmentID, key)
}

func (s *Service) Snapshot() (Snapshot, error) {
	workspaces, err := s.app.Workspaces.List()
	if err != nil {
		return Snapshot{}, err
	}
	environments, err := s.app.EnvironmentSummaries()
	if err != nil {
		return Snapshot{}, err
	}
	allowed, err := s.app.AllowedExecutables()
	if err != nil {
		return Snapshot{}, err
	}
	mcps, err := s.app.MCPs.List()
	if err != nil {
		return Snapshot{}, err
	}
	skills, err := s.app.Skills.List()
	if err != nil {
		return Snapshot{}, err
	}
	globalMemory, err := s.app.Memory.GlobalList()
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		Workspaces:         nonNilWorkspaces(workspaces),
		Environments:       nonNilEnvironments(environments),
		AllowedExecutables: nonNilStrings(allowed),
		MCPs:               nonNilMCP(mcps),
		Skills:             nonNilCatalog(skills),
		GlobalMemoryCount:  len(globalMemory),
	}, nil
}

func nonNilWorkspaces(values []model.Workspace) []model.Workspace {
	if values == nil {
		return []model.Workspace{}
	}
	return values
}

func nonNilEnvironments(values []app.EnvironmentSummary) []app.EnvironmentSummary {
	if values == nil {
		return []app.EnvironmentSummary{}
	}
	return values
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func nonNilMCP(values []model.MCPDefinition) []model.MCPDefinition {
	if values == nil {
		return []model.MCPDefinition{}
	}
	return values
}

func nonNilCatalog(values []model.CatalogEntry) []model.CatalogEntry {
	if values == nil {
		return []model.CatalogEntry{}
	}
	return values
}
