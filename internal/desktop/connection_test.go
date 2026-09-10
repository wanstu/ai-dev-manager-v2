package desktop

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/gateway"
	"ai-dev-manager-v2/internal/management"
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
