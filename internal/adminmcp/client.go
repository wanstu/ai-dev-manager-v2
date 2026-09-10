package adminmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/verifier"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Client struct {
	endpoint string
}

type ProcessStatus struct {
	ID             string     `json:"id"`
	EnvironmentID  string     `json:"environment_id"`
	State          string     `json:"state"`
	PID            int        `json:"pid,omitempty"`
	StartedAt      time.Time  `json:"started_at"`
	ExitedAt       *time.Time `json:"exited_at,omitempty"`
	ExitCode       *int       `json:"exit_code,omitempty"`
	ListeningPorts []int      `json:"listening_ports,omitempty"`
	ErrorKind      string     `json:"error_kind,omitempty"`
}

type ProcessLogs struct {
	ProcessID       string `json:"process_id"`
	Stdout          string `json:"stdout,omitempty"`
	Stderr          string `json:"stderr,omitempty"`
	StdoutTruncated bool   `json:"stdout_truncated,omitempty"`
	StderrTruncated bool   `json:"stderr_truncated,omitempty"`
}

type RunStatus struct {
	ID            string     `json:"id"`
	EnvironmentID string     `json:"environment_id"`
	State         string     `json:"state"`
	Executable    string     `json:"executable"`
	Args          []string   `json:"args,omitempty"`
	Cwd           string     `json:"cwd,omitempty"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	ExitCode      *int       `json:"exit_code,omitempty"`
	Stdout        string     `json:"stdout,omitempty"`
	Stderr        string     `json:"stderr,omitempty"`
	ErrorKind     string     `json:"error_kind,omitempty"`
	Message       string     `json:"message,omitempty"`
}

func New(endpoint string) *Client {
	return &Client{endpoint: strings.TrimSpace(endpoint)}
}

func callAdmin[T any](client *Client, ctx context.Context, tool string, arguments map[string]any) (T, error) {
	var zero T
	if client == nil || client.endpoint == "" {
		return zero, fmt.Errorf("ADM Admin MCP is not connected")
	}
	mcpClient := mcp.NewClient(&mcp.Implementation{Name: "ai-dev-manager-v2-admin-client", Version: "v1.0.0-rc.1"}, nil)
	session, err := mcpClient.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: client.endpoint}, nil)
	if err != nil {
		return zero, fmt.Errorf("connect Admin MCP %s: %w", client.endpoint, err)
	}
	defer session.Close()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: compactArguments(arguments)})
	if err != nil {
		return zero, fmt.Errorf("call Admin MCP tool %s: %w", tool, err)
	}
	if result == nil {
		return zero, fmt.Errorf("Admin MCP tool %s returned no result", tool)
	}
	if result.IsError {
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok && strings.TrimSpace(text.Text) != "" {
				return zero, fmt.Errorf("Admin MCP tool %s failed: %s", tool, strings.TrimSpace(text.Text))
			}
		}
		if toolErr := result.GetError(); toolErr != nil {
			return zero, fmt.Errorf("Admin MCP tool %s failed: %w", tool, toolErr)
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

func compactArguments(arguments map[string]any) map[string]any {
	if len(arguments) == 0 {
		return map[string]any{}
	}
	compacted := make(map[string]any, len(arguments))
	for key, value := range arguments {
		if value == nil {
			continue
		}
		reflected := reflect.ValueOf(value)
		switch reflected.Kind() {
		case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface:
			if reflected.IsNil() {
				continue
			}
		}
		compacted[key] = value
	}
	return compacted
}

func callAdminNoResult(client *Client, ctx context.Context, tool string, arguments map[string]any) error {
	_, err := callAdmin[map[string]any](client, ctx, tool, arguments)
	return err
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func nonNilStringMap(values map[string]string) map[string]string {
	if values == nil {
		return map[string]string{}
	}
	return values
}

func (c *Client) Snapshot() (management.Snapshot, error) {
	return callAdmin[management.Snapshot](c, context.Background(), "management_snapshot", map[string]any{})
}

func (c *Client) WorkspaceList() ([]model.Workspace, error) {
	snapshot, err := c.Snapshot()
	return snapshot.Workspaces, err
}

func (c *Client) EnvironmentList() ([]app.EnvironmentSummary, error) {
	snapshot, err := c.Snapshot()
	return snapshot.Environments, err
}

func (c *Client) ExecList() ([]string, error) {
	snapshot, err := c.Snapshot()
	return snapshot.AllowedExecutables, err
}

func (c *Client) MCPList() ([]model.MCPDefinition, error) {
	snapshot, err := c.Snapshot()
	return snapshot.MCPs, err
}

func (c *Client) SkillList() ([]model.CatalogEntry, error) {
	snapshot, err := c.Snapshot()
	return snapshot.Skills, err
}

func (c *Client) EnvironmentCapabilityReport(environmentID string) (model.CapabilityReport, error) {
	return callAdmin[model.CapabilityReport](c, context.Background(), "environment_capability_report", map[string]any{"environment_id": environmentID})
}

func (c *Client) CapabilityReport(environmentID string) (model.CapabilityReport, error) {
	return c.EnvironmentCapabilityReport(environmentID)
}

func (c *Client) VerifierAdd(environmentID string, definition model.VerifierDefinition) (model.VerifierDefinition, error) {
	return callAdmin[model.VerifierDefinition](c, context.Background(), "environment_verifier_add", map[string]any{"environment_id": environmentID, "definition": definition})
}

func (c *Client) VerifierList(environmentID string) ([]model.VerifierDefinition, error) {
	return callAdmin[[]model.VerifierDefinition](c, context.Background(), "environment_verifier_list", map[string]any{"environment_id": environmentID})
}

func (c *Client) VerifierRemove(environmentID, verifierID string) (model.VerifierDefinition, error) {
	type removed struct {
		Removed model.VerifierDefinition `json:"removed"`
	}
	result, err := callAdmin[removed](c, context.Background(), "environment_verifier_remove", map[string]any{"environment_id": environmentID, "verifier_id": verifierID})
	return result.Removed, err
}

func (c *Client) VerifierRun(environmentID, writerOwner, verifierID string, maxOutputBytes int) (verifier.Result, error) {
	return callAdmin[verifier.Result](c, context.Background(), "environment_verifier_run", map[string]any{
		"environment_id": environmentID, "writer_owner": writerOwner, "verifier_id": verifierID, "max_output_bytes": maxOutputBytes,
	})
}

func (c *Client) ProcessList(environmentID string) ([]ProcessStatus, error) {
	return callAdmin[[]ProcessStatus](c, context.Background(), "process_list", map[string]any{"environment_id": environmentID})
}

func (c *Client) ProcessStatus(environmentID, processID string) (ProcessStatus, error) {
	return callAdmin[ProcessStatus](c, context.Background(), "process_status", map[string]any{"environment_id": environmentID, "process_id": processID})
}

func (c *Client) ProcessLogs(environmentID, processID string) (ProcessLogs, error) {
	return callAdmin[ProcessLogs](c, context.Background(), "process_logs", map[string]any{"environment_id": environmentID, "process_id": processID})
}

func (c *Client) ProcessStop(environmentID, writerOwner, processID string) (ProcessStatus, error) {
	return callAdmin[ProcessStatus](c, context.Background(), "process_stop", map[string]any{"environment_id": environmentID, "writer_owner": writerOwner, "process_id": processID})
}

func (c *Client) RunList(environmentID string) ([]RunStatus, error) {
	return callAdmin[[]RunStatus](c, context.Background(), "run_list", map[string]any{"environment_id": environmentID})
}

func (c *Client) RunStatus(environmentID, runID string) (RunStatus, error) {
	return callAdmin[RunStatus](c, context.Background(), "run_status", map[string]any{"environment_id": environmentID, "run_id": runID})
}

func (c *Client) RunCancel(environmentID, writerOwner, runID string) (RunStatus, error) {
	return callAdmin[RunStatus](c, context.Background(), "run_cancel", map[string]any{"environment_id": environmentID, "writer_owner": writerOwner, "run_id": runID})
}

func (c *Client) WriterAcquire(environmentID, owner string) (model.Environment, error) {
	return callAdmin[model.Environment](c, context.Background(), "environment_writer_acquire", map[string]any{"environment_id": environmentID, "owner": owner})
}

func (c *Client) WriterHeartbeat(environmentID, owner string) (model.Environment, error) {
	return callAdmin[model.Environment](c, context.Background(), "environment_writer_heartbeat", map[string]any{"environment_id": environmentID, "owner": owner})
}

func (c *Client) WriterRelease(environmentID, owner string, force bool) (model.Environment, error) {
	return callAdmin[model.Environment](c, context.Background(), "environment_writer_release", map[string]any{"environment_id": environmentID, "owner": owner, "force": force})
}
func (c *Client) WorkspaceInspect(id string) (model.Workspace, error) {
	return callAdmin[model.Workspace](c, context.Background(), "workspace_inspect", map[string]any{"workspace_id": id})
}
func (c *Client) WorkspaceAdd(path, name string) (model.Workspace, error) {
	return callAdmin[model.Workspace](c, context.Background(), "workspace_add", map[string]any{"path": path, "name": name})
}
func (c *Client) WorkspaceRename(id, name string) (model.Workspace, error) {
	return callAdmin[model.Workspace](c, context.Background(), "workspace_rename", map[string]any{"workspace_id": id, "name": name})
}
func (c *Client) WorkspaceRemove(id string) (model.Workspace, error) {
	type removed struct {
		Removed model.Workspace `json:"removed"`
	}
	result, err := callAdmin[removed](c, context.Background(), "workspace_remove", map[string]any{"workspace_id": id})
	return result.Removed, err
}
func (c *Client) EnvironmentInspect(id string) (app.EnvironmentInspection, error) {
	return callAdmin[app.EnvironmentInspection](c, context.Background(), "environment_inspect", map[string]any{"environment_id": id})
}
func (c *Client) EnvironmentCreate(workspaceID, name, root string) (app.EnvironmentSummary, error) {
	created, err := callAdmin[model.Environment](c, context.Background(), "environment_create", map[string]any{"workspace_id": workspaceID, "name": name, "root": root})
	if err != nil {
		return app.EnvironmentSummary{}, err
	}
	return c.environmentSummary(created.ID)
}
func (c *Client) EnvironmentRename(id, name string) (app.EnvironmentSummary, error) {
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
func (c *Client) EnvironmentRemove(id string) error {
	_, err := c.EnvironmentRemoveResult(id)
	return err
}

func (c *Client) EnvironmentRemoveResult(id string) (model.Environment, error) {
	type removed struct {
		Removed model.Environment `json:"removed"`
	}
	result, err := callAdmin[removed](c, context.Background(), "environment_remove", map[string]any{"environment_id": id})
	return result.Removed, err
}
func (c *Client) ExecAllow(executable string) ([]string, error) {
	return callAdmin[[]string](c, context.Background(), "exec_allow", map[string]any{"executable": executable})
}
func (c *Client) ExecRemove(executable string) ([]string, error) {
	return callAdmin[[]string](c, context.Background(), "exec_allow_remove", map[string]any{"executable": executable})
}
func (c *Client) MCPAddConfig(name string, config catalog.MCPConfig) (model.MCPDefinition, error) {
	return callAdmin[model.MCPDefinition](c, context.Background(), "mcp_add", map[string]any{
		"name": name, "transport": config.Transport, "auth_mode": config.AuthMode, "endpoint": config.Endpoint,
		"header_refs": nonNilStringMap(config.HeaderRefs), "executable": config.Executable, "args": nonNilStrings(config.Args), "env_refs": nonNilStringMap(config.EnvRefs),
		"health_policy": config.HealthPolicy, "default_include_in_environment": config.DefaultInclude,
	})
}
func (c *Client) MCPImportPreview(input app.MCPImportInput) (app.MCPImportPreview, error) {
	return callAdmin[app.MCPImportPreview](c, context.Background(), "mcp_import_preview", map[string]any{
		"format": input.Format, "json_or_jsonc": input.Content, "selected_names": nonNilStrings(input.SelectedNames), "conflict_policy": input.ConflictPolicy,
		"default_include": input.DefaultInclude, "source_scope": input.SourceScope,
	})
}
func (c *Client) MCPImportApply(input app.MCPImportInput) (app.MCPImportApplyResult, error) {
	return callAdmin[app.MCPImportApplyResult](c, context.Background(), "mcp_import_apply", map[string]any{
		"format": input.Format, "json_or_jsonc": input.Content, "selected_names": nonNilStrings(input.SelectedNames), "conflict_policy": input.ConflictPolicy,
		"default_include": input.DefaultInclude, "source_scope": input.SourceScope,
	})
}
func (c *Client) MCPRemove(id string) error {
	return callAdminNoResult(c, context.Background(), "mcp_remove", map[string]any{"id": id})
}
func (c *Client) MCPSetDefault(id string, value bool) (model.MCPDefinition, error) {
	return callAdmin[model.MCPDefinition](c, context.Background(), "mcp_set_default", map[string]any{"id": id, "default_include_in_environment": value})
}
func (c *Client) MCPHealth(ctx context.Context, environmentID, mcpID string) (app.MCPHealthStatus, error) {
	return callAdmin[app.MCPHealthStatus](c, ctx, "environment_mcp_status", map[string]any{"environment_id": environmentID, "mcp_id": mcpID})
}
func (c *Client) SkillAdd(root, supportRoot string, defaultInclude bool) ([]model.CatalogEntry, error) {
	supportRoots := []string{}
	if strings.TrimSpace(supportRoot) != "" {
		supportRoots = append(supportRoots, supportRoot)
	}
	return callAdmin[[]model.CatalogEntry](c, context.Background(), "skill_add", map[string]any{"root": root, "support_roots": nonNilStrings(supportRoots), "default_include_in_environment": defaultInclude})
}
func (c *Client) SkillSourceAdd(root string, supportRoots []string, defaultInclude bool) (model.SkillSource, error) {
	return callAdmin[model.SkillSource](c, context.Background(), "skill_source_add", map[string]any{"root": root, "support_roots": nonNilStrings(supportRoots), "default_include_in_environment": defaultInclude})
}
func (c *Client) SkillSourceList() ([]model.SkillSource, error) {
	return callAdmin[[]model.SkillSource](c, context.Background(), "skill_source_list", map[string]any{})
}
func (c *Client) SkillSourceRefresh(id string) (catalog.SkillSourceRefreshResult, error) {
	return callAdmin[catalog.SkillSourceRefreshResult](c, context.Background(), "skill_source_refresh", map[string]any{"id": id})
}
func (c *Client) SkillSourceRemove(id string) (catalog.SkillSourceRefreshResult, error) {
	return callAdmin[catalog.SkillSourceRefreshResult](c, context.Background(), "skill_source_remove", map[string]any{"id": id})
}
func (c *Client) SkillRemove(id string) error {
	return callAdminNoResult(c, context.Background(), "skill_remove", map[string]any{"id": id})
}
func (c *Client) SkillSetDefault(id string, value bool) (model.CatalogEntry, error) {
	return callAdmin[model.CatalogEntry](c, context.Background(), "skill_set_default", map[string]any{"id": id, "default_include_in_environment": value})
}
func (c *Client) EnvironmentMCPSet(environmentID, mcpID string, enabled bool) (app.EnvironmentSummary, error) {
	if _, err := callAdmin[model.Environment](c, context.Background(), "environment_mcp_set", map[string]any{"environment_id": environmentID, "id": mcpID, "enabled": enabled}); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return c.environmentSummary(environmentID)
}
func (c *Client) EnvironmentSkillSet(environmentID, skillID string, enabled bool) (app.EnvironmentSummary, error) {
	if _, err := callAdmin[model.Environment](c, context.Background(), "environment_skill_set", map[string]any{"environment_id": environmentID, "id": skillID, "enabled": enabled}); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return c.environmentSummary(environmentID)
}
func (c *Client) environmentSummary(id string) (app.EnvironmentSummary, error) {
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
func (c *Client) EnvironmentSkillList(environmentID string) (app.SkillAvailabilityList, error) {
	return callAdmin[app.SkillAvailabilityList](c, context.Background(), "environment_skill_list", map[string]any{"environment_id": environmentID})
}
func (c *Client) EnvironmentSkillInspect(environmentID, skillID string) (app.SkillAvailability, error) {
	return callAdmin[app.SkillAvailability](c, context.Background(), "environment_skill_inspect", map[string]any{"environment_id": environmentID, "skill_id": skillID})
}
func (c *Client) GlobalMemoryList() ([]memory.Entry, error) {
	return callAdmin[[]memory.Entry](c, context.Background(), "memory_global_list", map[string]any{})
}
func (c *Client) GlobalMemoryRead(key string) (memory.Entry, error) {
	return callAdmin[memory.Entry](c, context.Background(), "memory_global_read", map[string]any{"key": key})
}
func (c *Client) GlobalMemoryWrite(key, value string) error {
	return callAdminNoResult(c, context.Background(), "memory_global_write", map[string]any{"key": key, "value": value})
}
func (c *Client) GlobalMemoryDelete(key string) error {
	return callAdminNoResult(c, context.Background(), "memory_global_delete", map[string]any{"key": key})
}
func (c *Client) EnvironmentMemoryList(environmentID string) ([]memory.Entry, error) {
	return callAdmin[[]memory.Entry](c, context.Background(), "memory_environment_list", map[string]any{"environment_id": environmentID})
}
func (c *Client) EnvironmentMemoryRead(environmentID, key string) (memory.Entry, error) {
	return callAdmin[memory.Entry](c, context.Background(), "memory_environment_read", map[string]any{"environment_id": environmentID, "key": key})
}
func (c *Client) EnvironmentMemoryWrite(environmentID, key, value string) error {
	return callAdminNoResult(c, context.Background(), "memory_environment_write", map[string]any{"environment_id": environmentID, "key": key, "value": value})
}
func (c *Client) EnvironmentMemoryDelete(environmentID, key string) error {
	return callAdminNoResult(c, context.Background(), "memory_environment_delete", map[string]any{"environment_id": environmentID, "key": key})
}
