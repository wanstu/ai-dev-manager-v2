package main

import (
	"context"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
)

type cliManagementBackend interface {
	WorkspaceList() ([]model.Workspace, error)
	WorkspaceInspect(string) (model.Workspace, error)
	WorkspaceAdd(string, string) (model.Workspace, error)
	WorkspaceRename(string, string) (model.Workspace, error)
	WorkspaceRemove(string) (model.Workspace, error)

	EnvironmentCreate(string, string, string) (app.EnvironmentSummary, error)
	EnvironmentList() ([]app.EnvironmentSummary, error)
	EnvironmentInspect(string) (app.EnvironmentInspection, error)
	CapabilityReport(string) (model.CapabilityReport, error)
	EnvironmentRename(string, string) (app.EnvironmentSummary, error)
	EnvironmentRemoveResult(string) (model.Environment, error)
	EnvironmentMCPSet(string, string, bool) (app.EnvironmentSummary, error)
	EnvironmentSkillSet(string, string, bool) (app.EnvironmentSummary, error)

	VerifierAdd(string, model.VerifierDefinition) (model.VerifierDefinition, error)
	VerifierList(string) ([]model.VerifierDefinition, error)
	VerifierRemove(string, string) (model.VerifierDefinition, error)
	WriterAcquire(string, string) (model.Environment, error)
	WriterHeartbeat(string, string) (model.Environment, error)
	WriterRelease(string, string, bool) (model.Environment, error)

	ExecAllow(string) ([]string, error)
	ExecRemove(string) ([]string, error)
	ExecList() ([]string, error)

	MCPAddConfig(string, catalog.MCPConfig) (model.MCPDefinition, error)
	MCPList() ([]model.MCPDefinition, error)
	MCPImportPreview(app.MCPImportInput) (app.MCPImportPreview, error)
	MCPImportApply(app.MCPImportInput) (app.MCPImportApplyResult, error)
	MCPHealth(context.Context, string, string) (app.MCPHealthStatus, error)
	MCPRemove(string) error
	MCPSetDefault(string, bool) (model.MCPDefinition, error)

	SkillAdd(string, string, bool) ([]model.CatalogEntry, error)
	SkillList() ([]model.CatalogEntry, error)
	SkillSourceAdd(string, []string, bool) (model.SkillSource, error)
	SkillSourceList() ([]model.SkillSource, error)
	SkillSourceRefresh(string) (catalog.SkillSourceRefreshResult, error)
	SkillSourceRemove(string) (catalog.SkillSourceRefreshResult, error)
	SkillRemove(string) error
	SkillSetDefault(string, bool) (model.CatalogEntry, error)

	GlobalMemoryList() ([]memory.Entry, error)
	GlobalMemoryRead(string) (memory.Entry, error)
	GlobalMemoryWrite(string, string) error
	GlobalMemoryDelete(string) error
	EnvironmentMemoryList(string) ([]memory.Entry, error)
	EnvironmentMemoryRead(string, string) (memory.Entry, error)
	EnvironmentMemoryWrite(string, string, string) error
	EnvironmentMemoryDelete(string, string) error
}
