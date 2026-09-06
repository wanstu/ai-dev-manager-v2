package gateway

import (
	"context"
	"errors"
	"net"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGatewayDevelopsPlainDirectoryWithoutGit(t *testing.T) {
	root := t.TempDir()
	if _, err := os.Stat(filepath.Join(root, ".git")); !os.IsNotExist(err) {
		t.Fatalf("fixture must not be a git repository")
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "plain")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	names := toolNames(tools.Tools)
	for _, required := range []string{"workspace_list", "workspace_add", "workspace_inspect", "workspace_rename", "workspace_remove", "exec_allow", "exec_allow_remove", "environment_create", "environment_remove", "environment_writer_acquire", "environment_writer_heartbeat", "mcp_list", "mcp_add", "environment_mcp_set", "skill_list", "skill_add", "environment_skill_set", "memory_global_write", "memory_environment_write", "tree", "read", "search", "write", "edit", "delete", "exec", "git_status"} {
		if !contains(names, required) {
			t.Fatalf("missing gateway tool %q in %v", required, names)
		}
	}

	created, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_create",
		Arguments: map[string]any{"workspace_id": ws.ID, "name": "mcp-task"},
	})
	if err != nil || created.IsError {
		t.Fatalf("environment_create failed: err=%v result=%+v", err, created)
	}
	envs, err := service.Environments.List()
	if err != nil || len(envs) != 1 {
		t.Fatalf("environment persistence after MCP create: envs=%+v err=%v", envs, err)
	}
	envID := envs[0].ID

	acquired, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_writer_acquire",
		Arguments: map[string]any{"environment_id": envID, "owner": "mcp-session"},
	})
	if err != nil || acquired.IsError {
		t.Fatalf("writer acquire failed: err=%v result=%+v", err, acquired)
	}
	if !strings.Contains(toolText(t, acquired), "expires_at") {
		t.Fatalf("writer acquire result must expose lease expiry: %s", toolText(t, acquired))
	}

	heartbeat, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_writer_heartbeat",
		Arguments: map[string]any{"environment_id": envID, "owner": "mcp-session"},
	})
	if err != nil || heartbeat.IsError {
		t.Fatalf("writer heartbeat failed: err=%v result=%+v", err, heartbeat)
	}

	written, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "write",
		Arguments: map[string]any{
			"environment_id": envID,
			"writer_owner":   "mcp-session",
			"path":           "plain.txt",
			"content":        "works without git\n",
		},
	})
	if err != nil || written.IsError {
		t.Fatalf("write failed: err=%v result=%+v", err, written)
	}

	read, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "read",
		Arguments: map[string]any{"environment_id": envID, "path": "plain.txt"},
	})
	if err != nil || read.IsError {
		t.Fatalf("read failed: err=%v result=%+v", err, read)
	}
	if !strings.Contains(toolText(t, read), "works without git") {
		t.Fatalf("unexpected read result: %s", toolText(t, read))
	}

	gitResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "git_status",
		Arguments: map[string]any{"environment_id": envID},
	})
	if err != nil {
		t.Fatalf("git_status transport error = %v", err)
	}
	if !gitResult.IsError {
		t.Fatalf("git_status should be a local tool error for non-git root: %+v", gitResult)
	}

	data, err := os.ReadFile(filepath.Join(root, "plain.txt"))
	if err != nil || string(data) != "works without git\n" {
		t.Fatalf("MCP write did not reach plain directory: data=%q err=%v", data, err)
	}

	released, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_writer_release",
		Arguments: map[string]any{"environment_id": envID, "owner": "mcp-session"},
	})
	if err != nil || released.IsError {
		t.Fatalf("writer release failed: err=%v result=%+v", err, released)
	}
	removed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_remove",
		Arguments: map[string]any{"environment_id": envID},
	})
	if err != nil || removed.IsError {
		t.Fatalf("environment_remove failed: err=%v result=%+v", err, removed)
	}
	if envs, err := service.Environments.List(); err != nil || len(envs) != 0 {
		t.Fatalf("environment_remove did not remove ADM record: envs=%+v err=%v", envs, err)
	}
	if data, err := os.ReadFile(filepath.Join(root, "plain.txt")); err != nil || string(data) != "works without git\n" {
		t.Fatalf("environment_remove must not delete project files: data=%q err=%v", data, err)
	}
}

func TestGatewayWorkspaceManagementGuardsProjectData(t *testing.T) {
	root := t.TempDir()
	sentinel := filepath.Join(root, "keep.txt")
	if err := os.WriteFile(sentinel, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	ws, err := service.Workspaces.Add(root, "before")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	inspected, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_inspect",
		Arguments: map[string]any{"workspace_id": ws.ID},
	})
	if err != nil || inspected.IsError || !strings.Contains(toolText(t, inspected), ws.ID) {
		t.Fatalf("workspace_inspect failed: err=%v result=%+v", err, inspected)
	}
	renamed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_rename",
		Arguments: map[string]any{"workspace_id": ws.ID, "name": "after"},
	})
	if err != nil || renamed.IsError || !strings.Contains(toolText(t, renamed), "after") {
		t.Fatalf("workspace_rename failed: err=%v result=%+v", err, renamed)
	}

	env, err := service.Environments.Create(ws.ID, "guard", "")
	if err != nil {
		t.Fatal(err)
	}
	blocked, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_remove",
		Arguments: map[string]any{"workspace_id": ws.ID},
	})
	if err != nil {
		t.Fatalf("workspace_remove transport error: %v", err)
	}
	if !blocked.IsError {
		t.Fatalf("workspace_remove must be rejected while Environment references Workspace: %+v", blocked)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("blocked workspace_remove touched project data: %v", err)
	}

	if _, err := service.Environments.Remove(env.ID); err != nil {
		t.Fatal(err)
	}
	removed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_remove",
		Arguments: map[string]any{"workspace_id": ws.ID},
	})
	if err != nil || removed.IsError {
		t.Fatalf("workspace_remove failed after Environment removal: err=%v result=%+v", err, removed)
	}
	if _, err := os.Stat(sentinel); err != nil {
		t.Fatalf("workspace_remove must not delete project data: %v", err)
	}
}

func TestGatewayExecAllowlistRemoveRevokesEntry(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	if err := service.AllowExecutable("go"); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	session := connectInMemory(t, ctx, New(service))
	defer session.Close()

	removed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "exec_allow_remove",
		Arguments: map[string]any{"executable": "GO"},
	})
	if err != nil || removed.IsError {
		t.Fatalf("exec_allow_remove failed: err=%v result=%+v", err, removed)
	}
	allowed, err := service.AllowedExecutables()
	if err != nil {
		t.Fatal(err)
	}
	if len(allowed) != 0 {
		t.Fatalf("allowlist after gateway removal = %v; want empty", allowed)
	}

	missing, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "exec_allow_remove",
		Arguments: map[string]any{"executable": "go"},
	})
	if err != nil {
		t.Fatalf("second exec_allow_remove transport error: %v", err)
	}
	if !missing.IsError {
		t.Fatalf("removing a missing allowlist entry must be a tool error: %+v", missing)
	}
}

func TestHTTPGatewayPersistsContextAcrossProcessRestart(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	projectRoot := t.TempDir()
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
	startGateway := func() *exec.Cmd {
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
		return cmd
	}
	stopGateway := func(cmd *exec.Cmd) {
		if err := cmd.Process.Kill(); err != nil {
			t.Fatalf("kill temporary Gateway: %v", err)
		}
		if err := cmd.Wait(); err != nil {
			var exitErr *exec.ExitError
			if !errors.As(err, &exitErr) {
				t.Fatalf("wait for temporary Gateway: %v", err)
			}
		}
		running = nil
	}

	first := startGateway()
	ctx := context.Background()
	firstSession := connectHTTPWithRetry(t, ctx, endpoint)
	added, err := firstSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "workspace_add",
		Arguments: map[string]any{"path": projectRoot, "name": "restart-persist"},
	})
	if err != nil || added.IsError {
		t.Fatalf("workspace_add before restart failed: err=%v result=%+v", err, added)
	}
	stateService := app.New(statePath)
	workspaces, err := stateService.Workspaces.List()
	if err != nil || len(workspaces) != 1 {
		t.Fatalf("Workspace state before restart = %+v err=%v", workspaces, err)
	}
	created, err := firstSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_create",
		Arguments: map[string]any{"workspace_id": workspaces[0].ID, "name": "restart-context"},
	})
	if err != nil || created.IsError {
		t.Fatalf("environment_create before restart failed: err=%v result=%+v", err, created)
	}
	environments, err := stateService.Environments.List()
	if err != nil || len(environments) != 1 {
		t.Fatalf("Environment state before restart = %+v err=%v", environments, err)
	}
	environmentID := environments[0].ID
	written, err := firstSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "memory_environment_write",
		Arguments: map[string]any{
			"environment_id": environmentID,
			"key":            "restart-check",
			"value":          "persisted",
		},
	})
	if err != nil || written.IsError {
		t.Fatalf("memory_environment_write before restart failed: err=%v result=%+v", err, written)
	}
	_ = firstSession.Close()
	stopGateway(first)

	second := startGateway()
	secondSession := connectHTTPWithRetry(t, ctx, endpoint)
	defer secondSession.Close()
	inspected, err := secondSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "environment_inspect",
		Arguments: map[string]any{"environment_id": environmentID},
	})
	if err != nil || inspected.IsError || !strings.Contains(toolText(t, inspected), environmentID) {
		t.Fatalf("environment_inspect after restart failed: err=%v result=%+v", err, inspected)
	}
	memoryResult, err := secondSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "memory_environment_read",
		Arguments: map[string]any{
			"environment_id": environmentID,
			"key":            "restart-check",
		},
	})
	if err != nil || memoryResult.IsError || !strings.Contains(toolText(t, memoryResult), "persisted") {
		t.Fatalf("Environment-private Memory after restart failed: err=%v result=%+v", err, memoryResult)
	}
	_ = secondSession.Close()
	stopGateway(second)
}

func TestHTTPGatewayRestartHelperProcess(t *testing.T) {
	if os.Getenv("ADM_TEST_GATEWAY_RESTART_HELPER") != "1" {
		return
	}
	statePath := os.Getenv("ADM_TEST_GATEWAY_RESTART_STATE")
	listen := os.Getenv("ADM_TEST_GATEWAY_RESTART_LISTEN")
	if statePath == "" || listen == "" {
		t.Fatal("restart helper requires state path and listen address")
	}
	if err := RunHTTP(context.Background(), app.New(statePath), listen); err != nil {
		t.Fatal(err)
	}
}

func connectHTTPWithRetry(t *testing.T, ctx context.Context, endpoint string) *mcp.ClientSession {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-restart-test", Version: "dev"}, nil)
		session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: endpoint}, nil)
		if err == nil {
			return session
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("connect temporary HTTP Gateway %s: %v", endpoint, lastErr)
	return nil
}

func TestHTTPGatewayUsesStreamableMCPAtMCPPath(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	httpServer := httptest.NewServer(NewHTTPHandler(service))
	defer httpServer.Close()

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-http-test", Version: "dev"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: httpServer.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect streamable HTTP gateway: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(toolNames(tools.Tools), "gateway_info") || !contains(toolNames(tools.Tools), "workspace_add") {
		t.Fatalf("unexpected HTTP gateway tools: %v", toolNames(tools.Tools))
	}
	for _, tool := range tools.Tools {
		if (tool.Name == "gateway_info" || tool.Name == "workspace_add" || tool.Name == "environment_remove" || tool.Name == "environment_writer_acquire" || tool.Name == "environment_writer_heartbeat") && tool.OutputSchema != nil {
			t.Fatalf("generic tool %q must omit outputSchema for broad MCP client compatibility; got %#v", tool.Name, tool.OutputSchema)
		}
		if tool.Name == "environment_inspect" && tool.OutputSchema == nil {
			t.Fatal("environment_inspect must retain its structured outputSchema")
		}
	}
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "gateway_info", Arguments: map[string]any{}})
	if err != nil || result.IsError || !strings.Contains(toolText(t, result), "ai-dev-manager-v2") {
		t.Fatalf("gateway_info over HTTP failed: err=%v result=%+v", err, result)
	}
}

func connectInMemory(t *testing.T, ctx context.Context, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect: %v", err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })
	client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-test", Version: "dev"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	return clientSession
}

func toolNames(tools []*mcp.Tool) []string {
	out := make([]string, 0, len(tools))
	for _, tool := range tools {
		out = append(out, tool.Name)
	}
	sort.Strings(out)
	return out
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func toolText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("tool result has no content")
	}
	text, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("tool result content type = %T", result.Content[0])
	}
	return text.Text
}
