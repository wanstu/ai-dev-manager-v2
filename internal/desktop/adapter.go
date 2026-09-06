package desktop

import (
	"errors"
	"fmt"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/memory"
	"ai-dev-manager-v2/internal/model"
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

type CatalogInput struct {
	Name           string `json:"name"`
	Endpoint       string `json:"endpoint,omitempty"`
	Instructions   string `json:"instructions,omitempty"`
	DefaultInclude bool   `json:"default_include_in_environment"`
}

type Adapter struct {
	management *management.Service
}

func NewAdapter(service *management.Service) *Adapter {
	return &Adapter{management: service}
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

func (a *Adapter) AddMCP(input CatalogInput) (model.CatalogEntry, error) {
	if err := a.ready(); err != nil {
		return model.CatalogEntry{}, err
	}
	return a.management.MCPAdd(input.Name, input.Endpoint, input.DefaultInclude)
}

func (a *Adapter) SetMCPDefault(id string, enabled bool) (model.CatalogEntry, error) {
	if err := a.ready(); err != nil {
		return model.CatalogEntry{}, err
	}
	return a.management.MCPSetDefault(id, enabled)
}

func (a *Adapter) RemoveMCP(id string) error {
	if err := a.ready(); err != nil {
		return err
	}
	return a.management.MCPRemove(id)
}

func (a *Adapter) AddSkill(input CatalogInput) (model.CatalogEntry, error) {
	if err := a.ready(); err != nil {
		return model.CatalogEntry{}, err
	}
	return a.management.SkillAdd(input.Name, input.Instructions, input.DefaultInclude)
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

func (a *Adapter) ready() error {
	if a == nil || a.management == nil {
		return errors.New("desktop management adapter is not initialized")
	}
	return nil
}
