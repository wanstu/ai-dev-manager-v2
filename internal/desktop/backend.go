package desktop

import (
	"context"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
)

type managementBackend interface {
	Snapshot() (management.Snapshot, error)
	WorkspaceInspect(string) (model.Workspace, error)
	WorkspaceAdd(string, string) (model.Workspace, error)
	WorkspaceRename(string, string) (model.Workspace, error)
	WorkspaceRemove(string) (model.Workspace, error)
	EnvironmentInspect(string) (app.EnvironmentInspection, error)
	EnvironmentCreate(string, string, string) (app.EnvironmentSummary, error)
	EnvironmentRename(string, string) (app.EnvironmentSummary, error)
	EnvironmentRemove(string) error
	ExecAllow(string) ([]string, error)
	ExecRemove(string) ([]string, error)
	MCPAddConfig(string, catalog.MCPConfig) (model.MCPDefinition, error)
	MCPImportPreview(app.MCPImportInput) (app.MCPImportPreview, error)
	MCPImportApply(app.MCPImportInput) (app.MCPImportApplyResult, error)
	MCPRemove(string) error
	MCPSetDefault(string, bool) (model.MCPDefinition, error)
	MCPHealth(context.Context, string, string) (app.MCPHealthStatus, error)
	SkillAdd(string, string, bool) ([]model.CatalogEntry, error)
	SkillSourceAdd(string, []string, bool) (model.SkillSource, error)
	SkillSourceList() ([]model.SkillSource, error)
	SkillSourceRefresh(string) (catalog.SkillSourceRefreshResult, error)
	SkillSourceRemove(string) (catalog.SkillSourceRefreshResult, error)
	SkillRemove(string) error
	SkillSetDefault(string, bool) (model.CatalogEntry, error)
	EnvironmentMCPSet(string, string, bool) (app.EnvironmentSummary, error)
	EnvironmentSkillSet(string, string, bool) (app.EnvironmentSummary, error)
	EnvironmentSkillList(string) (app.SkillAvailabilityList, error)
	EnvironmentSkillInspect(string, string) (app.SkillAvailability, error)
	GlobalMemoryList() ([]memory.Entry, error)
	GlobalMemoryRead(string) (memory.Entry, error)
	GlobalMemoryWrite(string, string) error
	GlobalMemoryDelete(string) error
	EnvironmentMemoryList(string) ([]memory.Entry, error)
	EnvironmentMemoryRead(string, string) (memory.Entry, error)
	EnvironmentMemoryWrite(string, string, string) error
	EnvironmentMemoryDelete(string, string) error
}
