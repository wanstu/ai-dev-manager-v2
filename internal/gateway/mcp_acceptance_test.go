package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHealthLifecycleEndToEnd(t *testing.T) {
	t.Setenv("ADM_MCP_ACCEPTANCE_TOKEN", "acceptance-secret")

	external := mcp.NewServer(&mcp.Implementation{Name: "external-acceptance", Version: "dev"}, nil)
	mcp.AddTool(external, &mcp.Tool{Name: "external_first", Description: "First external tool."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			return nil, map[string]any{"value": "first"}, nil
		})
	mcp.AddTool(external, &mcp.Tool{Name: "external_ping", Description: "Return an acceptance pong."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			return nil, map[string]any{"pong": "external-pong"}, nil
		})

	base := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return external }, &mcp.StreamableHTTPOptions{
		Stateless:                  true,
		DisableLocalhostProtection: true,
	})
	externalHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer acceptance-secret" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		base.ServeHTTP(w, r)
	}))
	closed := false
	t.Cleanup(func() {
		if !closed {
			externalHTTP.Close()
		}
	})

	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "marker.txt"), []byte("gateway-survives-broken-mcp\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "mcp-acceptance")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := service.MCPs.AddMCPConfig("external", catalog.MCPConfig{
		Endpoint:   externalHTTP.URL,
		Transport:  catalog.MCPTransportStreamableHTTP,
		HeaderRefs: map[string]string{"Authorization": "Bearer ${ADM_MCP_ACCEPTANCE_TOKEN}"},
	})
	if err != nil {
		t.Fatal(err)
	}
	env, err := service.Environments.Create(ws.ID, "mcp-runtime", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(env.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	status := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
		"environment_id": env.ID,
		"mcp_id":         entry.ID,
	})
	if status.IsError || !strings.Contains(toolText(t, status), "healthy") {
		t.Fatalf("healthy MCP status failed: %+v text=%s", status, toolText(t, status))
	}

	listed := callGatewayTool(t, ctx, session, "environment_mcp_tools", map[string]any{
		"environment_id": env.ID,
		"mcp_id":         entry.ID,
	})
	if listed.IsError {
		t.Fatalf("healthy MCP tool listing failed: %s", toolText(t, listed))
	}
	listedText := toolText(t, listed)
	if !strings.Contains(listedText, "external_first") || !strings.Contains(listedText, "external_ping") {
		t.Fatalf("external tool listing missing tools: %s", listedText)
	}

	called := callGatewayTool(t, ctx, session, "environment_mcp_call", map[string]any{
		"environment_id": env.ID,
		"mcp_id":         entry.ID,
		"tool":           "external_ping",
	})
	if called.IsError || !strings.Contains(toolText(t, called), "external-pong") {
		t.Fatalf("external MCP call failed: %s", toolText(t, called))
	}

	missingTool := callGatewayTool(t, ctx, session, "environment_mcp_call", map[string]any{
		"environment_id": env.ID,
		"mcp_id":         entry.ID,
		"tool":           "   ",
	})
	if !missingTool.IsError || !strings.Contains(toolText(t, missingTool), "error_kind=missing_tool_name") {
		t.Fatalf("empty MCP tool name must be structured error: %+v text=%s", missingTool, toolText(t, missingTool))
	}

	if _, err := service.SetEnvironmentMCP(env.ID, entry.ID, false); err != nil {
		t.Fatal(err)
	}
	disabled := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
		"environment_id": env.ID,
		"mcp_id":         entry.ID,
	})
	if disabled.IsError || !strings.Contains(toolText(t, disabled), "disabled") {
		t.Fatalf("disabled MCP status failed: %+v text=%s", disabled, toolText(t, disabled))
	}
	for _, name := range []string{"environment_mcp_tools", "environment_mcp_call"} {
		arguments := map[string]any{"environment_id": env.ID, "mcp_id": entry.ID}
		if name == "environment_mcp_call" {
			arguments["tool"] = "external_ping"
		}
		result := callGatewayTool(t, ctx, session, name, arguments)
		if !result.IsError || !strings.Contains(toolText(t, result), "error_kind=not_enabled") {
			t.Fatalf("%s must reject disabled MCP with not_enabled: %+v text=%s", name, result, toolText(t, result))
		}
	}

	if _, err := service.SetEnvironmentMCP(env.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}
	externalHTTP.Close()
	closed = true

	brokenStatus := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
		"environment_id": env.ID,
		"mcp_id":         entry.ID,
	})
	if brokenStatus.IsError {
		t.Fatalf("health status tool itself must return structured status, not a tool error: %s", toolText(t, brokenStatus))
	}
	brokenStatusText := toolText(t, brokenStatus)
	if !strings.Contains(brokenStatusText, "error") || !strings.Contains(brokenStatusText, "connection_refused") {
		t.Fatalf("broken MCP status must report structured error: %s", brokenStatusText)
	}

	for _, name := range []string{"environment_mcp_tools", "environment_mcp_call"} {
		arguments := map[string]any{"environment_id": env.ID, "mcp_id": entry.ID}
		if name == "environment_mcp_call" {
			arguments["tool"] = "external_ping"
		}
		result := callGatewayTool(t, ctx, session, name, arguments)
		text := toolText(t, result)
		if !result.IsError || !strings.Contains(text, "error_kind=connection_refused") {
			t.Fatalf("%s must return structured connection error: %+v text=%s", name, result, text)
		}
		if strings.Contains(text, externalHTTP.URL) || strings.Contains(text, "acceptance-secret") {
			t.Fatalf("%s leaked endpoint or secret: %s", name, text)
		}
	}

	gatewayInfo := callGatewayTool(t, ctx, session, "gateway_info", map[string]any{})
	if gatewayInfo.IsError {
		t.Fatalf("broken external MCP must not break gateway_info: %s", toolText(t, gatewayInfo))
	}
	read := callGatewayTool(t, ctx, session, "read", map[string]any{
		"environment_id": env.ID,
		"path":           "marker.txt",
	})
	if read.IsError || !strings.Contains(toolText(t, read), "gateway-survives-broken-mcp") {
		t.Fatalf("broken external MCP must not break read: %+v text=%s", read, toolText(t, read))
	}
}

func callGatewayTool(t *testing.T, ctx context.Context, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatalf("%s transport error: %v", name, err)
	}
	return result
}
