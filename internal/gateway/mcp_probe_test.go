package gateway

import (
	"ai-dev-manager-v2/internal/app"
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlobalMCPProbeAdminOnly(t *testing.T) {
	ctx := context.Background()
	s := app.New(filepath.Join(t.TempDir(), "state.json"))
	external := mcp.NewServer(&mcp.Implementation{Name: "global-probe", Version: "dev"}, nil)
	remote := httptest.NewServer(mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return external }, &mcp.StreamableHTTPOptions{Stateless: true, DisableLocalhostProtection: true}))
	defer remote.Close()
	entry, err := s.MCPs.AddMCP("global", remote.URL, false)
	if err != nil {
		t.Fatal(err)
	}
	admin := connectInMemory(t, ctx, newServerForSurface(s, nil, serverSurfaceAdmin))
	defer admin.Close()
	result := callGatewayTool(t, ctx, admin, "mcp_probe", map[string]any{"id": entry.ID})
	if result.IsError || !strings.Contains(toolText(t, result), `"state":"healthy"`) {
		t.Fatalf("probe=%s", toolText(t, result))
	}
	agent := connectInMemory(t, ctx, newServer(s, nil))
	defer agent.Close()
	inventory, err := agent.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tool := range inventory.Tools {
		if tool.Name == "mcp_probe" {
			t.Fatal("global management probe exposed to Agent")
		}
	}
	denied, err := agent.CallTool(ctx, &mcp.CallToolParams{Name: "mcp_probe", Arguments: map[string]any{"id": entry.ID}})
	if err == nil && denied != nil && !denied.IsError {
		t.Fatal("Agent invoked Admin-only probe")
	}
	ws, err := s.Workspaces.Add(t.TempDir(), "plain")
	if err != nil {
		t.Fatal(err)
	}
	env, err := s.Environments.Create(ws.ID, "plain", "")
	if err != nil {
		t.Fatal(err)
	}
	blocked := callGatewayTool(t, ctx, agent, "environment_mcp_tools", map[string]any{"environment_id": env.ID, "mcp_id": entry.ID})
	if !blocked.IsError {
		t.Fatal("global healthy probe granted disabled Environment access")
	}
}
