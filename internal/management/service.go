package management

import (
	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/model"
)

type Snapshot struct {
	Workspaces         []model.Workspace        `json:"workspaces"`
	Environments       []app.EnvironmentSummary `json:"environments"`
	AllowedExecutables []string                 `json:"allowed_executables"`
	MCPs               []model.CatalogEntry     `json:"mcps"`
	Skills             []model.CatalogEntry     `json:"skills"`
	GlobalMemoryCount  int                      `json:"global_memory_count"`
}

type Service struct {
	app *app.Service
}

func New(application *app.Service) *Service {
	return &Service{app: application}
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
		MCPs:               nonNilCatalog(mcps),
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

func nonNilCatalog(values []model.CatalogEntry) []model.CatalogEntry {
	if values == nil {
		return []model.CatalogEntry{}
	}
	return values
}
