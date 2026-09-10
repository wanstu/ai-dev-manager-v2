package desktop

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/adminmcp"
	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/verifier"
)

type WorkspaceInput struct {
	Path string `json:"path"`
	Name string `json:"name,omitempty"`
}

type EnvironmentInput struct {
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Root        string `json:"root,omitempty"`
}

type MCPInput struct {
	Name           string                `json:"name"`
	Transport      string                `json:"transport"`
	AuthMode       string                `json:"auth_mode"`
	Endpoint       string                `json:"endpoint,omitempty"`
	HeaderRefs     map[string]string     `json:"header_refs,omitempty"`
	Executable     string                `json:"executable,omitempty"`
	Args           []string              `json:"args,omitempty"`
	EnvRefs        map[string]string     `json:"env_refs,omitempty"`
	HealthPolicy   model.MCPHealthPolicy `json:"health_policy"`
	DefaultInclude bool                  `json:"default_include_in_environment"`
}

type MCPImportInput struct {
	Format         string   `json:"format"`
	Content        string   `json:"json_or_jsonc"`
	SelectedNames  []string `json:"selected_names,omitempty"`
	ConflictPolicy string   `json:"conflict_policy,omitempty"`
	DefaultInclude bool     `json:"default_include,omitempty"`
	SourceScope    string   `json:"source_scope,omitempty"`
}
type SkillInput struct {
	Root           string `json:"root"`
	SupportRoot    string `json:"support_root,omitempty"`
	DefaultInclude bool   `json:"default_include_in_environment"`
}
type SkillSourceInput struct {
	Root           string   `json:"root"`
	SupportRoots   []string `json:"support_roots,omitempty"`
	DefaultInclude bool     `json:"default_include_in_environment"`
}

type ADMConnectionInput struct {
	BaseURL string `json:"base_url"`
}

type ADMConnectionStatus struct {
	State                  string `json:"state"`
	BaseURL                string `json:"base_url"`
	HealthURL              string `json:"health_url"`
	AgentMCPURL            string `json:"agent_mcp_url"`
	AdminMCPURL            string `json:"admin_mcp_url"`
	PID                    int    `json:"pid,omitempty"`
	Version                string `json:"version,omitempty"`
	OwnerID                string `json:"owner_id,omitempty"`
	Detail                 string `json:"detail,omitempty"`
	LocalBootstrapEligible bool   `json:"local_bootstrap_eligible"`
}

type Adapter struct {
	profilesPath string
	management   managementBackend
	runtime      runtimeBackend
}

// NewAdapter retains an explicit local backend for tests and offline/recovery callers.
// Production Desktop uses NewClientAdapter and connects through Admin MCP.
func NewAdapter(service *management.Service) *Adapter {
	return &Adapter{management: service}
}

func NewClientAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) GetSnapshot() (management.Snapshot, error) {
	if err := a.ready(); err != nil {
		return management.Snapshot{}, err
	}
	return a.management.Snapshot()
}

func (a *Adapter) GetGatewayStatus() (gateway.HTTPStatus, error) {
	if err := a.ready(); err != nil {
		return gateway.HTTPStatus{}, err
	}
	return gateway.InspectHTTP(gateway.DefaultHTTPListen)
}

func (a *Adapter) InspectADMConnection(input ADMConnectionInput) (ADMConnectionStatus, error) {
	if a == nil {
		return ADMConnectionStatus{}, errors.New("desktop adapter is not initialized")
	}
	baseURL := strings.TrimSpace(input.BaseURL)
	if baseURL == "" {
		baseURL = defaultADMBaseURL()
	}
	status, err := gateway.InspectHTTPBaseURL(baseURL)
	if err != nil {
		return ADMConnectionStatus{}, err
	}
	return desktopConnectionStatus(status), nil
}

func (a *Adapter) ConnectADM(input ADMConnectionInput) (ADMConnectionStatus, error) {
	status, err := a.InspectADMConnection(input)
	if err != nil {
		if a != nil {
			a.management = nil
			a.runtime = nil
		}
		return ADMConnectionStatus{}, err
	}
	if status.State == gateway.HTTPStateRunning {
		client := adminmcp.New(status.AdminMCPURL)
		a.management = client
		a.runtime = client
	} else {
		a.management = nil
		a.runtime = nil
	}
	return status, nil
}

func (a *Adapter) StartLocalADM(input ADMConnectionInput) (ADMConnectionStatus, error) {
	if a == nil {
		return ADMConnectionStatus{}, errors.New("desktop adapter is not initialized")
	}
	status, err := a.InspectADMConnection(input)
	if err != nil {
		return ADMConnectionStatus{}, err
	}
	listen, err := localBootstrapListen(status.BaseURL)
	if err != nil {
		return status, err
	}
	switch status.State {
	case gateway.HTTPStateRunning:
		client := adminmcp.New(status.AdminMCPURL)
		a.management = client
		a.runtime = client
		return status, nil
	case gateway.HTTPStateIncompatible:
		return status, fmt.Errorf("refusing to start local ADM because %s is incompatible: %s", status.BaseURL, status.Detail)
	}
	if err := gateway.CheckHTTPListenAvailable(listen); err != nil {
		return status, err
	}
	process, err := startDetachedGatewayProcess(listen)
	if err != nil {
		return status, fmt.Errorf("start detached Gateway: %w", err)
	}
	ready, err := gateway.WaitHTTPReady(listen, 5*time.Second)
	if err != nil {
		_ = process.Kill()
		_ = process.Release()
		return status, err
	}
	_ = process.Release()
	connected := desktopConnectionStatus(ready)
	client := adminmcp.New(connected.AdminMCPURL)
	a.management = client
	a.runtime = client
	return connected, nil
}

func (a *Adapter) StopLocalADM(input ADMConnectionInput) (ADMConnectionStatus, error) {
	if a == nil {
		return ADMConnectionStatus{}, errors.New("desktop adapter is not initialized")
	}
	status, err := a.InspectADMConnection(input)
	if err != nil {
		return ADMConnectionStatus{}, err
	}
	listen, err := localBootstrapListen(status.BaseURL)
	if err != nil {
		return status, err
	}
	stopped, err := gateway.StopHTTP(listen)
	if err != nil {
		return status, err
	}
	a.management = nil
	a.runtime = nil
	return desktopConnectionStatus(stopped), nil
}

func (a *Adapter) StartGateway() (gateway.HTTPStatus, error) {
	if err := a.ready(); err != nil {
		return gateway.HTTPStatus{}, err
	}
	status, err := gateway.InspectHTTP(gateway.DefaultHTTPListen)
	if err != nil {
		return gateway.HTTPStatus{}, err
	}
	switch status.State {
	case gateway.HTTPStateRunning:
		return status, nil
	case gateway.HTTPStateIncompatible:
		return status, fmt.Errorf("refusing to start Gateway because %s is incompatible: %s", gateway.DefaultHTTPListen, status.Detail)
	}
	if err := gateway.CheckHTTPListenAvailable(gateway.DefaultHTTPListen); err != nil {
		return gateway.HTTPStatus{}, err
	}

	process, err := startDetachedGatewayProcess(gateway.DefaultHTTPListen)
	if err != nil {
		return gateway.HTTPStatus{}, fmt.Errorf("start detached Gateway: %w", err)
	}
	ready, err := gateway.WaitHTTPReady(gateway.DefaultHTTPListen, 5*time.Second)
	if err != nil {
		_ = process.Kill()
		_ = process.Release()
		return gateway.HTTPStatus{}, err
	}
	_ = process.Release()
	return ready, nil
}

func (a *Adapter) StopGateway() (gateway.HTTPStatus, error) {
	if err := a.ready(); err != nil {
		return gateway.HTTPStatus{}, err
	}
	return gateway.StopHTTP(gateway.DefaultHTTPListen)
}

func (a *Adapter) InspectWorkspace(id string) (model.Workspace, error) {
	if err := a.ready(); err != nil {
		return model.Workspace{}, err
	}
	return a.management.WorkspaceInspect(id)
}

func (a *Adapter) AddWorkspace(input WorkspaceInput) (model.Workspace, error) {
	if err := a.ready(); err != nil {
		return model.Workspace{}, err
	}
	return a.management.WorkspaceAdd(input.Path, input.Name)
}

func (a *Adapter) RenameWorkspace(id, name string) (model.Workspace, error) {
	if err := a.ready(); err != nil {
		return model.Workspace{}, err
	}
	return a.management.WorkspaceRename(id, name)
}

func (a *Adapter) RemoveWorkspace(id string) (model.Workspace, error) {
	if err := a.ready(); err != nil {
		return model.Workspace{}, err
	}
	return a.management.WorkspaceRemove(id)
}

func (a *Adapter) InspectEnvironment(id string) (app.EnvironmentInspection, error) {
	if err := a.ready(); err != nil {
		return app.EnvironmentInspection{}, err
	}
	return a.management.EnvironmentInspect(id)
}

func (a *Adapter) CreateEnvironment(input EnvironmentInput) (app.EnvironmentSummary, error) {
	if err := a.ready(); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return a.management.EnvironmentCreate(input.WorkspaceID, input.Name, input.Root)
}

func (a *Adapter) RenameEnvironment(id, name string) (app.EnvironmentSummary, error) {
	if err := a.ready(); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return a.management.EnvironmentRename(id, name)
}

func (a *Adapter) RemoveEnvironment(id string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.EnvironmentRemove(id)
}

func (a *Adapter) AllowExecutable(executable string) ([]string, error) {
	if err := a.ready(); err != nil {
		return nil, err
	}
	return a.management.ExecAllow(executable)
}

func (a *Adapter) RemoveExecutable(executable string) ([]string, error) {
	if err := a.ready(); err != nil {
		return nil, err
	}
	return a.management.ExecRemove(executable)
}

func (a *Adapter) AddMCP(input MCPInput) (model.MCPDefinition, error) {
	if err := a.ready(); err != nil {
		return model.MCPDefinition{}, err
	}
	return a.management.MCPAddConfig(input.Name, catalog.MCPConfig{
		Transport:      input.Transport,
		AuthMode:       input.AuthMode,
		Endpoint:       input.Endpoint,
		HeaderRefs:     input.HeaderRefs,
		Executable:     input.Executable,
		Args:           input.Args,
		EnvRefs:        input.EnvRefs,
		HealthPolicy:   input.HealthPolicy,
		DefaultInclude: input.DefaultInclude,
	})
}
func (a *Adapter) UpdateMCP(id string, input MCPInput) (model.MCPDefinition, error) {
	if err := a.ready(); err != nil {
		return model.MCPDefinition{}, err
	}
	return a.management.MCPUpdateConfig(id, input.Name, catalog.MCPConfig{
		Transport:      input.Transport,
		AuthMode:       input.AuthMode,
		Endpoint:       input.Endpoint,
		HeaderRefs:     input.HeaderRefs,
		Executable:     input.Executable,
		Args:           input.Args,
		EnvRefs:        input.EnvRefs,
		HealthPolicy:   input.HealthPolicy,
		DefaultInclude: input.DefaultInclude,
	})
}

func (a *Adapter) PreviewMCPImport(input MCPImportInput) (app.MCPImportPreview, error) {
	if err := a.ready(); err != nil {
		return app.MCPImportPreview{}, err
	}
	return a.management.MCPImportPreview(app.MCPImportInput{
		Format: input.Format, Content: input.Content, SelectedNames: input.SelectedNames,
		ConflictPolicy: input.ConflictPolicy, DefaultInclude: input.DefaultInclude, SourceScope: input.SourceScope,
	})
}

func (a *Adapter) ApplyMCPImport(input MCPImportInput) (app.MCPImportApplyResult, error) {
	if err := a.ready(); err != nil {
		return app.MCPImportApplyResult{}, err
	}
	return a.management.MCPImportApply(app.MCPImportInput{
		Format: input.Format, Content: input.Content, SelectedNames: input.SelectedNames,
		ConflictPolicy: input.ConflictPolicy, DefaultInclude: input.DefaultInclude, SourceScope: input.SourceScope,
	})
}

func (a *Adapter) SetMCPDefault(id string, enabled bool) (model.MCPDefinition, error) {
	if err := a.ready(); err != nil {
		return model.MCPDefinition{}, err
	}
	return a.management.MCPSetDefault(id, enabled)
}

func (a *Adapter) ProbeMCPHealth(environmentID, mcpID string) (app.MCPHealthStatus, error) {
	if err := a.ready(); err != nil {
		return app.MCPHealthStatus{}, err
	}
	return a.management.MCPHealth(context.Background(), environmentID, mcpID)
}

func (a *Adapter) RemoveMCP(id string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.MCPRemove(id)
}

func (a *Adapter) AddSkill(input SkillInput) ([]model.CatalogEntry, error) {
	if err := a.ready(); err != nil {
		return nil, err
	}
	return a.management.SkillAdd(input.Root, input.SupportRoot, input.DefaultInclude)
}
func (a *Adapter) AddSkillSource(input SkillSourceInput) (model.SkillSource, error) {
	if err := a.ready(); err != nil {
		return model.SkillSource{}, err
	}
	return a.management.SkillSourceAdd(input.Root, input.SupportRoots, input.DefaultInclude)
}

func (a *Adapter) ListSkillSources() ([]model.SkillSource, error) {
	if err := a.ready(); err != nil {
		return nil, err
	}
	return a.management.SkillSourceList()
}

func (a *Adapter) RefreshSkillSource(id string) (catalog.SkillSourceRefreshResult, error) {
	if err := a.ready(); err != nil {
		return catalog.SkillSourceRefreshResult{}, err
	}
	return a.management.SkillSourceRefresh(id)
}

func (a *Adapter) RemoveSkillSource(id string) (catalog.SkillSourceRefreshResult, error) {
	if err := a.ready(); err != nil {
		return catalog.SkillSourceRefreshResult{}, err
	}
	return a.management.SkillSourceRemove(id)
}

func (a *Adapter) ListEnvironmentSkills(environmentID string) (app.SkillAvailabilityList, error) {
	if err := a.ready(); err != nil {
		return app.SkillAvailabilityList{}, err
	}
	return a.management.EnvironmentSkillList(environmentID)
}

func (a *Adapter) InspectEnvironmentSkill(environmentID, skillID string) (app.SkillAvailability, error) {
	if err := a.ready(); err != nil {
		return app.SkillAvailability{}, err
	}
	return a.management.EnvironmentSkillInspect(environmentID, skillID)
}

func (a *Adapter) SetSkillDefault(id string, enabled bool) (model.CatalogEntry, error) {
	if err := a.ready(); err != nil {
		return model.CatalogEntry{}, err
	}
	return a.management.SkillSetDefault(id, enabled)
}

func (a *Adapter) RemoveSkill(id string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.SkillRemove(id)
}

func (a *Adapter) SetEnvironmentMCP(environmentID, mcpID string, enabled bool) (app.EnvironmentSummary, error) {
	if err := a.ready(); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return a.management.EnvironmentMCPSet(environmentID, mcpID, enabled)
}

func (a *Adapter) SetEnvironmentSkill(environmentID, skillID string, enabled bool) (app.EnvironmentSummary, error) {
	if err := a.ready(); err != nil {
		return app.EnvironmentSummary{}, err
	}
	return a.management.EnvironmentSkillSet(environmentID, skillID, enabled)
}

func (a *Adapter) ListVerifiers(environmentID string) ([]model.VerifierDefinition, error) {
	if err := a.readyRuntime(); err != nil {
		return nil, err
	}
	return a.runtime.VerifierList(environmentID)
}

func (a *Adapter) RunVerifier(environmentID, writerOwner, verifierID string) (verifier.Result, error) {
	if err := a.readyRuntime(); err != nil {
		return verifier.Result{}, err
	}
	return a.runtime.VerifierRun(environmentID, writerOwner, verifierID, 64*1024)
}

func (a *Adapter) ListProcesses(environmentID string) ([]adminmcp.ProcessStatus, error) {
	if err := a.readyRuntime(); err != nil {
		return nil, err
	}
	return a.runtime.ProcessList(environmentID)
}

func (a *Adapter) GetProcessLogs(environmentID, processID string) (adminmcp.ProcessLogs, error) {
	if err := a.readyRuntime(); err != nil {
		return adminmcp.ProcessLogs{}, err
	}
	return a.runtime.ProcessLogs(environmentID, processID)
}

func (a *Adapter) StopProcess(environmentID, writerOwner, processID string) (adminmcp.ProcessStatus, error) {
	if err := a.readyRuntime(); err != nil {
		return adminmcp.ProcessStatus{}, err
	}
	return a.runtime.ProcessStop(environmentID, writerOwner, processID)
}

func (a *Adapter) ListRuns(environmentID string) ([]adminmcp.RunStatus, error) {
	if err := a.readyRuntime(); err != nil {
		return nil, err
	}
	return a.runtime.RunList(environmentID)
}

func (a *Adapter) CancelRun(environmentID, writerOwner, runID string) (adminmcp.RunStatus, error) {
	if err := a.readyRuntime(); err != nil {
		return adminmcp.RunStatus{}, err
	}
	return a.runtime.RunCancel(environmentID, writerOwner, runID)
}

func (a *Adapter) ListGlobalMemory() ([]memory.Entry, error) {
	if err := a.ready(); err != nil {
		return nil, err
	}
	return a.management.GlobalMemoryList()
}

func (a *Adapter) ReadGlobalMemory(key string) (memory.Entry, error) {
	if err := a.ready(); err != nil {
		return memory.Entry{}, err
	}
	return a.management.GlobalMemoryRead(key)
}

func (a *Adapter) WriteGlobalMemory(key, value string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.GlobalMemoryWrite(key, value)
}

func (a *Adapter) DeleteGlobalMemory(key string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.GlobalMemoryDelete(key)
}

func (a *Adapter) ListEnvironmentMemory(environmentID string) ([]memory.Entry, error) {
	if err := a.ready(); err != nil {
		return nil, err
	}
	return a.management.EnvironmentMemoryList(environmentID)
}

func (a *Adapter) ReadEnvironmentMemory(environmentID, key string) (memory.Entry, error) {
	if err := a.ready(); err != nil {
		return memory.Entry{}, err
	}
	return a.management.EnvironmentMemoryRead(environmentID, key)
}

func (a *Adapter) WriteEnvironmentMemory(environmentID, key, value string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.EnvironmentMemoryWrite(environmentID, key, value)
}

func (a *Adapter) DeleteEnvironmentMemory(environmentID, key string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.EnvironmentMemoryDelete(environmentID, key)
}

func defaultADMBaseURL() string {
	baseURL, _ := gateway.HTTPBaseURL(gateway.DefaultHTTPListen)
	return baseURL
}

func desktopConnectionStatus(status gateway.HTTPStatus) ADMConnectionStatus {
	return ADMConnectionStatus{
		State:                  status.State,
		BaseURL:                status.BaseURL,
		HealthURL:              strings.TrimRight(status.BaseURL, "/") + "/healthz",
		AgentMCPURL:            status.MCPURL,
		AdminMCPURL:            status.AdminMCPURL,
		PID:                    status.PID,
		Version:                status.Version,
		OwnerID:                status.OwnerID,
		Detail:                 status.Detail,
		LocalBootstrapEligible: isLocalBootstrapBaseURL(status.BaseURL),
	}
}

func isLocalBootstrapBaseURL(raw string) bool {
	_, err := localBootstrapListen(raw)
	return err == nil
}

func localBootstrapListen(raw string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("invalid ADM base URL: %w", err)
	}
	if parsed.Scheme != "http" {
		return "", fmt.Errorf("local ADM bootstrap requires an http URL")
	}
	if strings.Trim(parsed.Path, "/") != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("local ADM bootstrap requires a root base URL without path/query/fragment")
	}
	host := parsed.Hostname()
	port := parsed.Port()
	if host == "" || port == "" {
		return "", fmt.Errorf("local ADM bootstrap requires an explicit loopback host and port")
	}
	local := strings.EqualFold(host, "localhost")
	if !local {
		ip := net.ParseIP(host)
		local = ip != nil && ip.IsLoopback()
	}
	if !local {
		return "", fmt.Errorf("local ADM bootstrap is only available for loopback addresses")
	}
	return net.JoinHostPort(host, port), nil
}

func (a *Adapter) ready() error {
	if a == nil {
		return errors.New("desktop management adapter is not initialized")
	}
	if a.management == nil {
		return errors.New("ADM Admin MCP is not connected")
	}
	return nil
}

func (a *Adapter) readyRuntime() error {
	if a == nil {
		return errors.New("desktop runtime adapter is not initialized")
	}
	if a.runtime == nil {
		return errors.New("ADM Admin MCP runtime is not connected")
	}
	return nil
}
