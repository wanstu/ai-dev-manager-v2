package desktop

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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

type adminManagementClient struct {
	endpoint string
}

func newAdminManagementClient(endpoint string) *adminManagementClient {
	return &adminManagementClient{endpoint: strings.TrimSpace(endpoint)}
}

func callAdmin[T any](client *adminManagementClient, ctx context.Context, tool string, arguments map[string]any) (T, error) {
	var zero T
	if client == nil || client.endpoint == "" {
		return zero, fmt.Errorf("ADM Admin MCP is not connected")
	}
	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "ai-dev-manager-v2-desktop", Version: "v0.1.0-dev"}, nil)
	session, err := mcpClient.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: client.endpoint}, nil)
	if err != nil {
		return zero, fmt.Errorf("connect Admin MCP %s: %w", client.endpoint, err)
	}
	defer session.Close()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: arguments})
	if err != nil {
		return zero, fmt.Errorf("call Admin MCP tool %s: %w", tool, err)
	}
	if result == nil {
		return zero, fmt.Errorf("Admin MCP tool %s returned no result", tool)
	}
	if result.IsError {
		if toolErr := result.GetError(); toolErr != nil {
			return zero, toolErr
		}
		return zero, fmt.Errorf("Admin MCP tool %s failed", tool)
	}
	if result.StructuredContent == nil {
		return zero, fmt.Errorf("Admin MCP tool %s returned no structured content", tool)
	}
	raw, err := json.Marshal(result.StructuredContent)
	if err != nil {
		return zero, fmt.Errorf("encode Admin MCP tool %s result: %w", tool, err)
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
	}
	if err := json.Unmarshal(raw, &envelope); err == nil && len(envelope.Result) > 0 && string(envelope.Result) != "null" {
		raw = envelope.Result
	}
	var output T
	if err := json.Unmarshal(raw, &output); err != nil {
		return zero, fmt.Errorf("decode Admin MCP tool %s result: %w", tool, err)
	}
	return output, nil
}

func callAdminNoResult(client *adminManagementClient, ctx context.Context, tool string, arguments map[string]any) error {
	_, err := callAdmin[map[string]any](client, ctx, tool, arguments)
	return err
}

func (c *adminManagementClient) Snapshot() (management.Snapshot, error) {
	return callAdmin[management.Snapshot](c, context.Background(), "management_snapshot", map[string]any{})
}
func (c *adminManagementClient) WorkspaceInspect(id string) (model.Workspace, error) {
	return callAdmin[model.Workspace](c, context.Background(), "workspace_inspect", map[string]any{"workspace_id": id})
}
func (c *adminManagementClient) WorkspaceAdd(path, name string) (model.Workspace, error) {
	return callAdmin[model.Workspace](c, context.Background(), "workspace_add", map[string]any{"path": path, "name": name})
}
func (c *adminManagementClient) WorkspaceRename(id, name string) (model.Workspace, error) {
	return callAdmin[model.Workspace](c, context.Background(), "workspace_rename", map[string]any{"workspace_id": id, "name": name})
}
func (c *adminManagementClient) WorkspaceRemove(id string) (model.Workspace, error) {
	type removed struct {
		Removed model.Workspace `json:"removed"`
	}
	result, err := callAdmin[removed](c, context.Background(), "workspace_remove", map[string]any{"workspace_id": id})
	return result.Removed, err
}
func (c *adminManagementClient) EnvironmentInspect(id string) (app.EnvironmentInspection, error) {
	return callAdmin[app.EnvironmentInspection](c, context.Background(), "environment_inspect", map[string]any{"environment_id": id})
}
func (c *adminManagementClient) EnvironmentCreate(workspaceID, name, root string) (app.EnvironmentSummary, error) {
	created, err := callAdmin[model.Environment](c, context.Background(), "environment_create", map[string]any{"workspace_id": workspaceID, "name": name, "root": root})
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	return c.environmentSummary(created.ID)
}
func (c *adminManagementClient) EnvironmentRename(id, name string) (app.EnvironmentSummary, error) {
	if _, err := callAdmin[model.Environment](c, context.Background(), "environment_rename", map[string]any{"environment_id": id, "name": name}); err != nil {
		return app.EnvironmentSummary{}, err
	}
	snapshot, err := c.Snapshot()
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	for _, environment := range snapshot.Environments {
		if environment.ID == id {
			return environment, nil
		}
	}
	return app.EnvironmentSummary{}, fmt.Errorf("renamed Environment %s was not present in management snapshot", id)
}
func (c *adminManagementClient) EnvironmentRemove(id string) error {
	return callAdminNoResult(c, context.Background(), "environment_remove", map[string]any{"environment_id": id})
}
func (c *adminManagementClient) ExecAllow(executable string) ([]string, error) {
	return callAdmin[[]string](c, context.Background(), "exec_allow", map[string]any{"executable": executable})
}
func (c *adminManagementClient) ExecRemove(executable string) ([]string, error) {
	return callAdmin[[]string](c, context.Background(), "exec_allow_remove", map[string]any{"executable": executable})
}
func (c *adminManagementClient) MCPAddConfig(name string, config catalog.MCPConfig) (model.MCPDefinition, error) {
	return callAdmin[model.MCPDefinition](c, context.Background(), "mcp_add", map[string]any{
		"name": name, "transport": config.Transport, "auth_mode": config.AuthMode, "endpoint": config.Endpoint,
		"header_refs": config.HeaderRefs, "executable": config.Executable, "args": config.Args, "env_refs": config.EnvRefs,
		"health_policy": config.HealthPolicy, "default_include_in_environment": config.DefaultInclude,
	})
}
func (c *adminManagementClient) MCPImportPreview(input app.MCPImportInput) (app.MCPImportPreview, error) {
	return callAdmin[app.MCPImportPreview](c, context.Background(), "mcp_import_preview", map[string]any{
		"format": input.Format, "json_or_jsonc": input.Content, "selected_names": input.SelectedNames, "conflict_policy": input.ConflictPolicy,
		"default_include": input.DefaultInclude, "source_scope": input.SourceScope,
	})
}
func (c *adminManagementClient) MCPImportApply(input app.MCPImportInput) (app.MCPImportApplyResult, error) {
	return callAdmin[app.MCPImportApplyResult](c, context.Background(), "mcp_import_apply", map[string]any{
		"format": input.Format, "json_or_jsonc": input.Content, "selected_names": input.SelectedNames, "conflict_policy": input.ConflictPolicy,
		"default_include": input.DefaultInclude, "source_scope": input.SourceScope,
	})
}
func (c *adminManagementClient) MCPRemove(id string) error {
	return callAdminNoResult(c, context.Background(), "mcp_remove", map[string]any{"id": id})
}
func (c *adminManagementClient) MCPSetDefault(id string, value bool) (model.MCPDefinition, error) {
	return callAdmin[model.MCPDefinition](c, context.Background(), "mcp_set_default", map[string]any{"id": id, "default_include_in_environment": value})
}
func (c *adminManagementClient) MCPHealth(ctx context.Context, environmentID, mcpID string) (app.MCPHealthStatus, error) {
	return callAdmin[app.MCPHealthStatus](c, ctx, "environment_mcp_status", map[string]any{"environment_id": environmentID, "mcp_id": mcpID})
}
func (c *adminManagementClient) SkillAdd(root, supportRoot string, defaultInclude bool) ([]model.CatalogEntry, error) {
	supportRoots := []string{}
	if strings.TrimSpace(supportRoot) != "" {
		supportRoots = append(supportRoots, supportRoot)
	}
	return callAdmin[[]model.CatalogEntry](c, context.Background(), "skill_add", map[string]any{"root": root, "support_roots": supportRoots, "default_include_in_environment": defaultInclude})
}
func (c *adminManagementClient) SkillSourceAdd(root string, supportRoots []string, defaultInclude bool) (model.SkillSource, error) {
	return callAdmin[model.SkillSource](c, context.Background(), "skill_source_add", map[string]any{"root": root, "support_roots": supportRoots, "default_include_in_environment": defaultInclude})
}
func (c *adminManagementClient) SkillSourceList() ([]model.SkillSource, error) {
	return callAdmin[[]model.SkillSource](c, context.Background(), "skill_source_list", map[string]any{})
}
func (c *adminManagementClient) SkillSourceRefresh(id string) (catalog.SkillSourceRefreshResult, error) {
	return callAdmin[catalog.SkillSourceRefreshResult](c, context.Background(), "skill_source_refresh", map[string]any{"id": id})
}
func (c *adminManagementClient) SkillSourceRemove(id string) (catalog.SkillSourceRefreshResult, error) {
	return callAdmin[catalog.SkillSourceRefreshResult](c, context.Background(), "skill_source_remove", map[string]any{"id": id})
}
func (c *adminManagementClient) SkillRemove(id string) error {
	return callAdminNoResult(c, context.Background(), "skill_remove", map[string]any{"id": id})
}
func (c *adminManagementClient) SkillSetDefault(id string, value bool) (model.CatalogEntry, error) {
	return callAdmin[model.CatalogEntry](c, context.Background(), "skill_set_default", map[string]any{"id": id, "default_include_in_environment": value})
}
func (c *adminManagementClient) EnvironmentMCPSet(environmentID, mcpID string, enabled bool) (app.EnvironmentSummary, error) {
	if _, err := callAdmin[model.Environment](c, context.Background(), "environment_mcp_set", map[string]any{"environment_id": environmentID, "id": mcpID, "enabled": enabled}); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return c.environmentSummary(environmentID)
}
func (c *adminManagementClient) EnvironmentSkillSet(environmentID, skillID string, enabled bool) (app.EnvironmentSummary, error) {
	if _, err := callAdmin[model.Environment](c, context.Background(), "environment_skill_set", map[string]any{"environment_id": environmentID, "id": skillID, "enabled": enabled}); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return c.environmentSummary(environmentID)
}
func (c *adminManagementClient) environmentSummary(id string) (app.EnvironmentSummary, error) {
	snapshot, err := c.Snapshot()
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	for _, environment := range snapshot.Environments {
		if environment.ID == id {
			return environment, nil
		}
	}
	return app.EnvironmentSummary{}, fmt.Errorf("Environment %s was not present in management snapshot", id)
}
func (c *adminManagementClient) EnvironmentSkillList(environmentID string) (app.SkillAvailabilityList, error) {
	return callAdmin[app.SkillAvailabilityList](c, context.Background(), "environment_skill_list", map[string]any{"environment_id": environmentID})
}
func (c *adminManagementClient) EnvironmentSkillInspect(environmentID, skillID string) (app.SkillAvailability, error) {
	return callAdmin[app.SkillAvailability](c, context.Background(), "environment_skill_inspect", map[string]any{"environment_id": environmentID, "skill_id": skillID})
}
func (c *adminManagementClient) GlobalMemoryList() ([]memory.Entry, error) {
	return callAdmin[[]memory.Entry](c, context.Background(), "memory_global_list", map[string]any{})
}
func (c *adminManagementClient) GlobalMemoryRead(key string) (memory.Entry, error) {
	return callAdmin[memory.Entry](c, context.Background(), "memory_global_read", map[string]any{"key": key})
}
func (c *adminManagementClient) GlobalMemoryWrite(key, value string) error {
	return callAdminNoResult(c, context.Background(), "memory_global_write", map[string]any{"key": key, "value": value})
}
func (c *adminManagementClient) GlobalMemoryDelete(key string) error {
	return callAdminNoResult(c, context.Background(), "memory_global_delete", map[string]any{"key": key})
}
func (c *adminManagementClient) EnvironmentMemoryList(environmentID string) ([]memory.Entry, error) {
	return callAdmin[[]memory.Entry](c, context.Background(), "memory_environment_list", map[string]any{"environment_id": environmentID})
}
func (c *adminManagementClient) EnvironmentMemoryRead(environmentID, key string) (memory.Entry, error) {
	return callAdmin[memory.Entry](c, context.Background(), "memory_environment_read", map[string]any{"environment_id": environmentID, "key": key})
}
func (c *adminManagementClient) EnvironmentMemoryWrite(environmentID, key, value string) error {
	return callAdminNoResult(c, context.Background(), "memory_environment_write", map[string]any{"environment_id": environmentID, "key": key, "value": value})
}
func (c *adminManagementClient) EnvironmentMemoryDelete(environmentID, key string) error {
	return callAdminNoResult(c, context.Background(), "memory_environment_delete", map[string]any{"environment_id": environmentID, "key": key})
}
