package desktop

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/management"
	"ai-dev-manager-v2/internal/model"
	"ai-dev-manager-v2/internal/pathutil"
)

func TestAdapterInspectsConfigurableADMBaseURL(t *testing.T) {
	gatewayService := app.New(filepath.Join(t.TempDir(), "gateway-state.json"))
	server := httptest.NewServer(http.StripPrefix("/control", gateway.NewHTTPHandler(gatewayService)))
	defer server.Close()

	adapter := NewAdapter(management.New(app.New(filepath.Join(t.TempDir(), "desktop-state.json"))))
	status, err := adapter.InspectADMConnection(ADMConnectionInput{BaseURL: server.URL + "/control/"})
	if err != nil {
		t.Fatal(err)
	}
	if status.State != gateway.HTTPStateRunning || status.BaseURL != server.URL+"/control" {
		t.Fatalf("status=%+v", status)
	}
	if status.HealthURL != server.URL+"/control/healthz" || status.AgentMCPURL != server.URL+"/control/mcp" || status.AdminMCPURL != server.URL+"/control/admin/mcp" {
		t.Fatalf("derived URLs=%+v", status)
	}
	if status.LocalBootstrapEligible {
		t.Fatalf("base-path profile must not be eligible for local process bootstrap: %+v", status)
	}
}

func TestClientAdapterUsesAdminMCPAndDoesNotFallbackAfterDisconnect(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "gateway-state.json")
	gatewayService := app.New(statePath)
	server := httptest.NewServer(gateway.NewHTTPHandler(gatewayService))

	adapter := NewClientAdapter()
	if _, err := adapter.GetSnapshot(); err == nil {
		t.Fatal("disconnected client adapter unexpectedly exposed local management state")
	}
	status, err := adapter.ConnectADM(ADMConnectionInput{BaseURL: server.URL})
	if err != nil || status.State != gateway.HTTPStateRunning {
		t.Fatalf("ConnectADM status=%+v err=%v", status, err)
	}
	root := t.TempDir()
	workspace, err := adapter.AddWorkspace(WorkspaceInput{Path: root, Name: "through-admin-mcp"})
	if err != nil {
		t.Fatal(err)
	}
	persisted, err := gatewayService.Workspaces.Get(workspace.ID)
	if err != nil || !pathutil.Same(persisted.Path, root) {
		t.Fatalf("Admin MCP workspace=%+v err=%v", persisted, err)
	}
	snapshot, err := adapter.GetSnapshot()
	if err != nil || len(snapshot.Workspaces) != 1 || snapshot.Workspaces[0].ID != workspace.ID {
		t.Fatalf("Admin MCP snapshot=%+v err=%v", snapshot, err)
	}
	mcpDefinition, err := adapter.AddMCP(MCPInput{Name: "editable", Transport: "streamable-http", AuthMode: "none", Endpoint: "http://127.0.0.1:9000/mcp"})
	if err != nil {
		t.Fatal(err)
	}
	updatedMCP, err := adapter.UpdateMCP(mcpDefinition.ID, MCPInput{Name: "editable-renamed", Transport: "streamable-http", AuthMode: "none", Endpoint: "http://127.0.0.1:9001/mcp"})
	if err != nil {
		t.Fatal(err)
	}
	persistedMCP, err := gatewayService.MCPs.Get(mcpDefinition.ID)
	if err != nil || updatedMCP.ID != mcpDefinition.ID || persistedMCP.ID != mcpDefinition.ID || persistedMCP.Name != "editable-renamed" || persistedMCP.Endpoint != "http://127.0.0.1:9001/mcp" {
		t.Fatalf("Admin MCP update result=%+v persisted=%+v err=%v", updatedMCP, persistedMCP, err)
	}
	environment, err := adapter.CreateEnvironment(EnvironmentInput{WorkspaceID: workspace.ID, Name: "runtime-ui"})
	if err != nil {
		t.Fatal(err)
	}
	definition, err := gatewayService.AddVerifier(environment.ID, model.VerifierDefinition{Kind: "test", Executable: "go", Enabled: true})
	if err != nil {
		t.Fatal(err)
	}
	verifiers, err := adapter.ListVerifiers(environment.ID)
	if err != nil || len(verifiers) != 1 || verifiers[0].ID != definition.ID {
		t.Fatalf("Admin MCP verifier list=%+v err=%v", verifiers, err)
	}

	server.Close()
	status, err = adapter.ConnectADM(ADMConnectionInput{BaseURL: server.URL})
	if err != nil {
		t.Fatalf("stopped Admin MCP should return stopped status, got %v", err)
	}
	if status.State != gateway.HTTPStateStopped {
		t.Fatalf("disconnected status=%+v", status)
	}
	if _, err := adapter.GetSnapshot(); err == nil || !strings.Contains(err.Error(), "not connected") {
		t.Fatalf("disconnected management fallback error=%v", err)
	}
	if _, err := adapter.ListProcesses(environment.ID); err == nil || !strings.Contains(err.Error(), "not connected") {
		t.Fatalf("disconnected runtime fallback error=%v", err)
	}
}

func TestLocalBootstrapListenRequiresLoopbackHTTPRoot(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		want    string
		wantErr bool
	}{
		{name: "ipv4", baseURL: "http://127.0.0.1:41137", want: "127.0.0.1:41137"},
		{name: "localhost", baseURL: "http://localhost:41137", want: "localhost:41137"},
		{name: "ipv6", baseURL: "http://[::1]:41137", want: "[::1]:41137"},
		{name: "https", baseURL: "https://127.0.0.1:41137", wantErr: true},
		{name: "remote", baseURL: "http://adm.example.test:41137", wantErr: true},
		{name: "base path", baseURL: "http://127.0.0.1:41137/control", wantErr: true},
		{name: "missing port", baseURL: "http://127.0.0.1", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := localBootstrapListen(tt.baseURL)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("localBootstrapListen(%q)=%q, want error", tt.baseURL, got)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("localBootstrapListen(%q)=%q err=%v want=%q", tt.baseURL, got, err, tt.want)
			}
		})
	}
}
