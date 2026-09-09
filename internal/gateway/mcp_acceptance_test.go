package gateway

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMCPHealthLifecycleEndToEnd(t *testing.T) {
	t.Setenv("ADM_MCP_ACCEPTANCE_TOKEN", "acceptance-secret")

	newExternal := func() *mcp.Server {
		external := mcp.NewServer(&mcp.Implementation{Name: "external-acceptance", Version: "dev"}, nil)
		mcp.AddTool(external, &mcp.Tool{Name: "external_first", Description: "First external tool."},
			func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
				return nil, map[string]any{"value": "first"}, nil
			})
		mcp.AddTool(external, &mcp.Tool{Name: "external_ping", Description: "Return an acceptance pong."},
			func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
				return nil, map[string]any{"pong": "external-pong"}, nil
			})
		return external
	}

	external := newExternal()
	base := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return external }, &mcp.StreamableHTTPOptions{
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
		Endpoint:   externalHTTP.URL + "/mcp",
		Transport:  catalog.MCPTransportStreamableHTTP,
		AuthMode:   catalog.MCPAuthHeaders,
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

func TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance(t *testing.T) {
	var broken atomic.Bool
	external := mcp.NewServer(&mcp.Implementation{Name: "background-reconnect-upstream", Version: "dev"}, nil)
	mcp.AddTool(external, &mcp.Tool{Name: "background_ping", Description: "Return background reconnect pong."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			return nil, map[string]any{"pong": "background-pong"}, nil
		})
	base := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return external }, &mcp.StreamableHTTPOptions{
		DisableLocalhostProtection: true,
	})
	externalHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if broken.Load() {
			http.Error(w, "temporarily unavailable", http.StatusServiceUnavailable)
			return
		}
		base.ServeHTTP(w, r)
	}))
	defer externalHTTP.Close()

	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	workspace, err := service.Workspaces.Add(t.TempDir(), "background-reconnect")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := service.MCPs.AddMCPConfig("background-reconnect", catalog.MCPConfig{
		Endpoint:  externalHTTP.URL + "/mcp",
		Transport: catalog.MCPTransportStreamableHTTP,
		HealthPolicy: model.MCPHealthPolicy{
			HealthCheckEnabled:       true,
			CheckIntervalSeconds:     1,
			ProbeTimeoutSeconds:      1,
			AutoReconnect:            true,
			ReconnectIntervalSeconds: 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	environment, err := service.Environments.Create(workspace.ID, "background-reconnect", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(environment.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	ctx := context.Background()
	gatewaySession := connectInMemory(t, ctx, newServer(service, owner))
	defer gatewaySession.Close()
	status, err := owner.Status(ctx, environment.ID, entry.ID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial owner status=%+v err=%v", status, err)
	}
	if tools, err := owner.ListTools(ctx, environment.ID, entry.ID); err != nil || len(tools) != 1 || tools[0].Name != "background_ping" {
		t.Fatalf("initial tool inventory tools=%+v err=%v", tools, err)
	}
	initialInspect := callGatewayTool(t, ctx, gatewaySession, "environment_mcp_inspect", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	initialInspectText := toolText(t, initialInspect)
	for _, required := range []string{"healthy", "background_ping", "inventory_fetched_at", "health_check_enabled", "auto_reconnect"} {
		if initialInspect.IsError || !strings.Contains(initialInspectText, required) {
			t.Fatalf("initial API inspect missing %q: error=%v text=%s", required, initialInspect.IsError, initialInspectText)
		}
	}

	broken.Store(true)
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(ctx)
	inspection, err := owner.Inspect(ctx, environment.ID, entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.State != app.MCPHealthError ||
		inspection.Observation.FailureStage != "ping" ||
		inspection.Observation.NextReconnectAt == nil ||
		len(inspection.Observation.Inventory) != 0 ||
		owner.Info().OwnedMCPSessions != 0 {
		t.Fatalf("background monitor did not mark broken upstream unhealthy: info=%+v observation=%+v", owner.Info(), inspection.Observation)
	}
	brokenInspect := callGatewayTool(t, ctx, gatewaySession, "environment_mcp_inspect", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	brokenInspectText := toolText(t, brokenInspect)
	for _, required := range []string{"error", "ping", "next_reconnect_at", "consecutive_failures"} {
		if brokenInspect.IsError || !strings.Contains(brokenInspectText, required) {
			t.Fatalf("broken API inspect missing %q: error=%v text=%s", required, brokenInspect.IsError, brokenInspectText)
		}
	}
	if strings.Contains(brokenInspectText, "background_ping") {
		t.Fatalf("broken API inspect retained stale inventory: %s", brokenInspectText)
	}

	broken.Store(false)
	owner.MonitorOnce(ctx)
	if owner.Info().OwnedMCPSessions != 0 {
		t.Fatalf("reconnect ran before fixed interval: info=%+v", owner.Info())
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(ctx)
	inspection, err = owner.Inspect(ctx, environment.ID, entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.State != app.MCPHealthHealthy ||
		inspection.Observation.NextReconnectAt != nil ||
		len(inspection.Observation.Inventory) != 1 ||
		inspection.Observation.Inventory[0].Name != "background_ping" ||
		owner.Info().OwnedMCPSessions != 1 {
		t.Fatalf("background reconnect did not restore upstream: info=%+v observation=%+v", owner.Info(), inspection.Observation)
	}
	recoveredInspect := callGatewayTool(t, ctx, gatewaySession, "environment_mcp_inspect", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	recoveredInspectText := toolText(t, recoveredInspect)
	for _, required := range []string{"healthy", "background_ping", "inventory_fetched_at"} {
		if recoveredInspect.IsError || !strings.Contains(recoveredInspectText, required) {
			t.Fatalf("recovered API inspect missing %q: error=%v text=%s", required, recoveredInspect.IsError, recoveredInspectText)
		}
	}
	called, err := owner.CallTool(ctx, environment.ID, entry.ID, "background_ping", nil)
	if err != nil || called == nil || called.IsError {
		t.Fatalf("explicit tool call after reconnect failed: result=%+v err=%v", called, err)
	}
}

func TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance(t *testing.T) {
	root := t.TempDir()
	executable, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADM_STDIO_BACKGROUND_TOKEN", "stdio-background-secret")

	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	workspace, err := service.Workspaces.Add(root, "stdio-background")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := service.MCPs.AddMCPConfig("stdio-background", catalog.MCPConfig{
		Transport:  catalog.MCPTransportStdio,
		Executable: executable,
		Args:       []string{"-test.run=^TestStdioMCPHelperProcess$"},
		EnvRefs: map[string]string{
			"ADM_TEST_STDIO_HELPER": "1",
			"ADM_TEST_STDIO_ROOT":   root,
			"ADM_TEST_STDIO_TOKEN":  "${ADM_STDIO_BACKGROUND_TOKEN}",
		},
		HealthPolicy: model.MCPHealthPolicy{
			HealthCheckEnabled:       true,
			CheckIntervalSeconds:     1,
			ProbeTimeoutSeconds:      1,
			AutoReconnect:            true,
			ReconnectIntervalSeconds: 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	environment, err := service.Environments.Create(workspace.ID, "stdio-background", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(environment.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}
	if err := service.AllowExecutable(executable); err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	ctx := context.Background()
	status, err := owner.Status(ctx, environment.ID, entry.ID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial stdio status=%+v err=%v", status, err)
	}
	if tools, err := owner.ListTools(ctx, environment.ID, entry.ID); err != nil || len(tools) == 0 || tools[0].Name == "" {
		t.Fatalf("initial stdio tools=%+v err=%v", tools, err)
	}
	if starts := fileLineCount(t, filepath.Join(root, "stdio-start-count.txt")); starts != 1 {
		t.Fatalf("initial stdio starts=%d, want 1", starts)
	}

	called, err := owner.CallTool(ctx, environment.ID, entry.ID, "stdio_echo", nil)
	if err != nil || called == nil || called.IsError {
		t.Fatalf("initial stdio call failed: result=%+v err=%v", called, err)
	}
	if calls := fileLineCount(t, filepath.Join(root, "stdio-echo-called.txt")); calls != 1 {
		t.Fatalf("initial stdio echo calls=%d, want 1", calls)
	}
	stop, err := owner.CallTool(ctx, environment.ID, entry.ID, "stdio_stop", nil)
	if err != nil || stop == nil || stop.IsError {
		t.Fatalf("stdio stop tool failed: result=%+v err=%v", stop, err)
	}
	waitForFile(t, filepath.Join(root, "stdio-stop-requested.txt"))

	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(ctx)
	inspection, err := owner.Inspect(ctx, environment.ID, entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.State != app.MCPHealthError ||
		inspection.Observation.FailureStage != "ping" ||
		inspection.Observation.NextReconnectAt == nil ||
		len(inspection.Observation.Inventory) != 0 ||
		owner.Info().OwnedMCPSessions != 0 {
		t.Fatalf("stdio monitor did not mark stopped helper unhealthy: info=%+v observation=%+v", owner.Info(), inspection.Observation)
	}

	owner.MonitorOnce(ctx)
	if starts := fileLineCount(t, filepath.Join(root, "stdio-start-count.txt")); starts != 1 {
		t.Fatalf("stdio reconnect ran before fixed interval: starts=%d", starts)
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(ctx)
	inspection, err = owner.Inspect(ctx, environment.ID, entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if starts := fileLineCount(t, filepath.Join(root, "stdio-start-count.txt")); starts != 2 {
		t.Fatalf("stdio background reconnect did not start a second helper: starts=%d observation=%+v", starts, inspection.Observation)
	}
	if inspection.Observation.State != app.MCPHealthHealthy ||
		inspection.Observation.NextReconnectAt != nil ||
		len(inspection.Observation.Inventory) == 0 ||
		owner.Info().OwnedMCPSessions != 1 {
		t.Fatalf("stdio background reconnect did not restore health: info=%+v observation=%+v", owner.Info(), inspection.Observation)
	}
	if calls := fileLineCount(t, filepath.Join(root, "stdio-echo-called.txt")); calls != 1 {
		t.Fatalf("background reconnect must not replay stdio tool calls: calls=%d", calls)
	}
	called, err = owner.CallTool(ctx, environment.ID, entry.ID, "stdio_echo", nil)
	if err != nil || called == nil || called.IsError {
		t.Fatalf("explicit stdio call after reconnect failed: result=%+v err=%v", called, err)
	}
	if calls := fileLineCount(t, filepath.Join(root, "stdio-echo-called.txt")); calls != 2 {
		t.Fatalf("explicit stdio call after reconnect did not run exactly once: calls=%d", calls)
	}
}

func TestStdioMCPAcceptanceAndProcessCleanup(t *testing.T) {
	root := t.TempDir()
	executable, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ADM_STDIO_ACCEPTANCE_TOKEN", "stdio-secret")

	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	workspace, err := service.Workspaces.Add(root, "stdio-acceptance")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := service.MCPs.AddMCPConfig("stdio", catalog.MCPConfig{
		Transport:  catalog.MCPTransportStdio,
		Executable: executable,
		Args:       []string{"-test.run=^TestStdioMCPHelperProcess$"},
		EnvRefs: map[string]string{
			"ADM_TEST_STDIO_HELPER": "1",
			"ADM_TEST_STDIO_ROOT":   root,
			"ADM_TEST_STDIO_TOKEN":  "${ADM_STDIO_ACCEPTANCE_TOKEN}",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	environment, err := service.Environments.Create(workspace.ID, "stdio-runtime", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(environment.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	ctx := context.Background()
	session := connectInMemory(t, ctx, newServer(service, owner))
	defer session.Close()

	rejected := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	if rejected.IsError || !strings.Contains(toolText(t, rejected), "executable_not_allowed") {
		t.Fatalf("stdio MCP must reject a non-allowlisted executable: %+v text=%s", rejected, toolText(t, rejected))
	}

	if err := service.AllowExecutable(executable); err != nil {
		t.Fatal(err)
	}
	status := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	if status.IsError || !strings.Contains(toolText(t, status), "healthy") {
		t.Fatalf("allowlisted stdio MCP status failed: %+v text=%s", status, toolText(t, status))
	}
	listed := callGatewayTool(t, ctx, session, "environment_mcp_tools", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	if listed.IsError || !strings.Contains(toolText(t, listed), "stdio_echo") {
		t.Fatalf("stdio MCP tool listing failed: %+v text=%s", listed, toolText(t, listed))
	}
	called := callGatewayTool(t, ctx, session, "environment_mcp_call", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
		"tool":           "stdio_echo",
	})
	calledText := toolText(t, called)
	if called.IsError || !strings.Contains(calledText, "stdio-secret") {
		t.Fatalf("stdio MCP tool call failed: %+v text=%s", called, calledText)
	}
	started, err := os.ReadFile(filepath.Join(root, "stdio-started.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(started), root) || !strings.Contains(string(started), "stdio-secret") {
		t.Fatalf("stdio helper did not inherit Environment cwd and resolved env ref: %q", started)
	}

	disabled := callGatewayTool(t, ctx, session, "environment_mcp_set", map[string]any{
		"environment_id": environment.ID,
		"id":             entry.ID,
		"enabled":        false,
	})
	if disabled.IsError {
		t.Fatalf("disable stdio MCP failed: %s", toolText(t, disabled))
	}
	waitForFile(t, filepath.Join(root, "stdio-stopped.txt"))

	if err := os.Remove(filepath.Join(root, "stdio-stopped.txt")); err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(environment.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}
	if status, err := owner.Status(ctx, environment.ID, entry.ID); err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("restart stdio MCP status=%+v err=%v", status, err)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
	waitForFile(t, filepath.Join(root, "stdio-stopped.txt"))
}

func TestStdioMCPHelperProcess(t *testing.T) {
	if os.Getenv("ADM_TEST_STDIO_HELPER") != "1" {
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		os.Exit(2)
	}
	root := os.Getenv("ADM_TEST_STDIO_ROOT")
	started := cwd + "\n" + os.Getenv("ADM_TEST_STDIO_TOKEN") + "\n"
	if err := os.WriteFile(filepath.Join(root, "stdio-started.txt"), []byte(started), 0o644); err != nil {
		os.Exit(2)
	}
	if err := appendTestLine(filepath.Join(root, "stdio-start-count.txt"), "started\n"); err != nil {
		os.Exit(2)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "stdio-acceptance", Version: "dev"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "stdio_echo", Description: "Return stdio acceptance context."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			if err := appendTestLine(filepath.Join(root, "stdio-echo-called.txt"), "echo\n"); err != nil {
				return nil, nil, err
			}
			return nil, map[string]any{"cwd": cwd, "token": os.Getenv("ADM_TEST_STDIO_TOKEN")}, nil
		})
	mcp.AddTool(server, &mcp.Tool{Name: "stdio_stop", Description: "Stop the stdio acceptance helper after returning."},
		func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
			if err := appendTestLine(filepath.Join(root, "stdio-stop-requested.txt"), "stop\n"); err != nil {
				return nil, nil, err
			}
			go func() {
				time.Sleep(50 * time.Millisecond)
				os.Exit(0)
			}()
			return nil, map[string]any{"stopping": true}, nil
		})
	err = server.Run(context.Background(), &mcp.StdioTransport{})
	_ = os.WriteFile(filepath.Join(root, "stdio-stopped.txt"), []byte("stopped\n"), 0o644)
	if err != nil {
		os.Exit(2)
	}
	os.Exit(0)
}

func waitForFile(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", path)
}

func appendTestLine(path, line string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(line)
	return err
}

func fileLineCount(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(data), "\n")
}

func TestRuntimeOwnerRealHTTPRestartReconciliation(t *testing.T) {
	newExternal := func() *mcp.Server {
		external := mcp.NewServer(&mcp.Implementation{Name: "phase5-owned-upstream", Version: "dev"}, nil)
		mcp.AddTool(external, &mcp.Tool{Name: "owner_ping", Description: "Return owner pong."},
			func(context.Context, *mcp.CallToolRequest, EmptyInput) (*mcp.CallToolResult, any, error) {
				return nil, map[string]any{"pong": "owner-pong"}, nil
			})
		return external
	}
	external := newExternal()
	base := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return external }, &mcp.StreamableHTTPOptions{
		DisableLocalhostProtection: true,
	})
	var requests atomic.Int64
	externalHTTP := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		base.ServeHTTP(w, r)
	}))
	externalClosed := false
	t.Cleanup(func() {
		if !externalClosed {
			externalHTTP.Close()
		}
	})

	statePath := filepath.Join(t.TempDir(), "state.json")
	service := app.New(statePath)
	workspace, err := service.Workspaces.Add(t.TempDir(), "phase5-owner")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := service.MCPs.AddMCPConfig("phase5-owned-upstream", catalog.MCPConfig{
		Endpoint:  externalHTTP.URL + "/mcp",
		Transport: catalog.MCPTransportStreamableHTTP,
	})
	if err != nil {
		t.Fatal(err)
	}
	environment, err := service.Environments.Create(workspace.ID, "phase5-owner", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(environment.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}
	freshService := app.New(statePath)
	freshEnvironment, err := freshService.Environments.Get(environment.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(freshEnvironment.EnabledMCPIDs) != 1 || freshEnvironment.EnabledMCPIDs[0] != entry.ID {
		t.Fatalf("fresh service did not load desired MCP selection: %+v", freshEnvironment.EnabledMCPIDs)
	}
	activation, desiredStatus, err := freshService.ResolveMCPActivation(environment.ID, entry.ID)
	if err != nil || activation == nil || desiredStatus.State != app.MCPHealthConfigured {
		t.Fatalf("fresh service activation=%+v status=%+v err=%v", activation, desiredStatus, err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listen := listener.Addr().String()
	_ = listener.Close()
	endpoint := "http://" + listen + "/mcp"

	var running *exec.Cmd
	t.Cleanup(func() {
		if running != nil && running.Process != nil {
			_ = running.Process.Kill()
			_ = running.Wait()
		}
	})
	startGateway := func() (*exec.Cmd, HTTPStatus) {
		cmd := exec.Command(os.Args[0], "-test.run=^TestHTTPGatewayRestartHelperProcess$")
		cmd.Env = append(os.Environ(),
			"ADM_TEST_GATEWAY_RESTART_HELPER=1",
			"ADM_TEST_GATEWAY_RESTART_STATE="+statePath,
			"ADM_TEST_GATEWAY_RESTART_LISTEN="+listen,
		)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		running = cmd
		status, err := WaitHTTPReady(listen, 5*time.Second)
		if err != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			running = nil
			t.Fatal(err)
		}
		if status.OwnerID == "" {
			t.Fatalf("started Gateway has no runtime owner: %+v", status)
		}
		return cmd, status
	}
	stopGateway := func(cmd *exec.Cmd) {
		stopped, err := StopHTTP(listen)
		if err != nil {
			t.Fatal(err)
		}
		if stopped.State != HTTPStateStopped {
			t.Fatalf("Gateway did not stop cleanly: %+v", stopped)
		}
		if err := cmd.Wait(); err != nil {
			t.Fatalf("graceful Gateway process exit: %v", err)
		}
		running = nil
	}
	waitUpstreamRequests := func(minimum int64) {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) && requests.Load() < minimum {
			time.Sleep(20 * time.Millisecond)
		}
		if got := requests.Load(); got < minimum {
			t.Fatalf("upstream requests=%d want at least %d", got, minimum)
		}
	}

	firstProcess, firstStatus := startGateway()
	waitUpstreamRequests(1)
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		session := connectHTTPWithRetry(t, ctx, endpoint)
		info := callGatewayTool(t, ctx, session, "gateway_info", map[string]any{})
		infoText := toolText(t, info)
		if info.IsError || !strings.Contains(infoText, firstStatus.OwnerID) || !strings.Contains(infoText, `"owned_mcp_sessions":1`) {
			_ = session.Close()
			t.Fatalf("client %d did not observe owner %q: %s", i+1, firstStatus.OwnerID, infoText)
		}
		status := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
			"environment_id": environment.ID,
			"mcp_id":         entry.ID,
		})
		if status.IsError || !strings.Contains(toolText(t, status), "healthy") {
			_ = session.Close()
			t.Fatalf("client %d owner status failed: %s", i+1, toolText(t, status))
		}
		if i == 1 {
			called := callGatewayTool(t, ctx, session, "environment_mcp_call", map[string]any{
				"environment_id": environment.ID,
				"mcp_id":         entry.ID,
				"tool":           "owner_ping",
			})
			if called.IsError || !strings.Contains(toolText(t, called), "owner-pong") {
				_ = session.Close()
				t.Fatalf("owned upstream call failed: %s", toolText(t, called))
			}
		}
		_ = session.Close()
	}
	firstRunRequests := requests.Load()
	stopGateway(firstProcess)

	secondProcess, secondStatus := startGateway()
	if secondStatus.OwnerID == firstStatus.OwnerID {
		t.Fatalf("Gateway restart reused runtime owner %q", firstStatus.OwnerID)
	}
	waitUpstreamRequests(firstRunRequests + 1)
	secondSession := connectHTTPWithRetry(t, ctx, endpoint)
	secondStatusResult := callGatewayTool(t, ctx, secondSession, "environment_mcp_status", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	if secondStatusResult.IsError || !strings.Contains(toolText(t, secondStatusResult), "healthy") {
		_ = secondSession.Close()
		t.Fatalf("reconciled MCP after restart is not healthy: %s", toolText(t, secondStatusResult))
	}
	inspected := callGatewayTool(t, ctx, secondSession, "environment_inspect", map[string]any{"environment_id": environment.ID})
	if inspected.IsError || !strings.Contains(toolText(t, inspected), entry.ID) {
		_ = secondSession.Close()
		t.Fatalf("persisted desired MCP selection missing after restart: %s", toolText(t, inspected))
	}

	externalHTTP.Close()
	externalClosed = true
	broken := callGatewayTool(t, ctx, secondSession, "environment_mcp_status", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	brokenText := toolText(t, broken)
	if broken.IsError || !strings.Contains(brokenText, "error") || strings.Contains(brokenText, "healthy") {
		_ = secondSession.Close()
		t.Fatalf("dead upstream retained healthy observation: %s", brokenText)
	}
	ownerAfterFailure := callGatewayTool(t, ctx, secondSession, "gateway_info", map[string]any{})
	if ownerAfterFailure.IsError || !strings.Contains(toolText(t, ownerAfterFailure), `"owned_mcp_sessions":0`) {
		_ = secondSession.Close()
		t.Fatalf("dead upstream session was not evicted: %s", toolText(t, ownerAfterFailure))
	}
	_ = secondSession.Close()
	stopGateway(secondProcess)

	thirdProcess, thirdStatus := startGateway()
	if thirdStatus.OwnerID == secondStatus.OwnerID || thirdStatus.OwnerID == firstStatus.OwnerID {
		t.Fatalf("second restart reused runtime owner: first=%q second=%q third=%q", firstStatus.OwnerID, secondStatus.OwnerID, thirdStatus.OwnerID)
	}
	thirdSession := connectHTTPWithRetry(t, ctx, endpoint)
	thirdResult := callGatewayTool(t, ctx, thirdSession, "environment_mcp_status", map[string]any{
		"environment_id": environment.ID,
		"mcp_id":         entry.ID,
	})
	thirdText := toolText(t, thirdResult)
	if thirdResult.IsError || !strings.Contains(thirdText, "error") || strings.Contains(thirdText, "healthy") {
		_ = thirdSession.Close()
		t.Fatalf("restart resurrected stale healthy state: %s", thirdText)
	}
	persistedEnvironment, err := service.Environments.Get(environment.ID)
	if err != nil {
		_ = thirdSession.Close()
		t.Fatal(err)
	}
	if len(persistedEnvironment.EnabledMCPIDs) != 1 || persistedEnvironment.EnabledMCPIDs[0] != entry.ID {
		_ = thirdSession.Close()
		t.Fatalf("observed failure rewrote desired selection: %+v", persistedEnvironment.EnabledMCPIDs)
	}
	_ = thirdSession.Close()
	stopGateway(thirdProcess)
}

func callGatewayTool(t *testing.T, ctx context.Context, session *mcp.ClientSession, name string, arguments map[string]any) *mcp.CallToolResult {
	t.Helper()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: arguments})
	if err != nil {
		t.Fatalf("%s transport error: %v", name, err)
	}
	return result
}
