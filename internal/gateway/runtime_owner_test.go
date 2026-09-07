package gateway

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"ai-dev-manager-v2/internal/app"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type fakeOwnedMCPSession struct {
	mu        sync.Mutex
	listCalls int
	callCalls int
	closes    int
	failAfter int
	listErr   error
}

func (s *fakeOwnedMCPSession) ListTools(context.Context, *mcp.ListToolsParams) (*mcp.ListToolsResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listCalls++
	if s.failAfter > 0 && s.listCalls > s.failAfter {
		if s.listErr != nil {
			return nil, s.listErr
		}
		return nil, errors.New("connection refused")
	}
	return &mcp.ListToolsResult{Tools: []*mcp.Tool{{Name: "owned_fake"}}}, nil
}

func (s *fakeOwnedMCPSession) CallTool(context.Context, *mcp.CallToolParams) (*mcp.CallToolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callCalls++
	return &mcp.CallToolResult{}, nil
}

func (s *fakeOwnedMCPSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closes++
	return nil
}

func (s *fakeOwnedMCPSession) counts() (list, call, closes int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.listCalls, s.callCalls, s.closes
}

func TestRuntimeOwnerReusesSessionAndClosesIt(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	fake := &fakeOwnedMCPSession{}
	connects := 0
	owner.connect = func(context.Context, string, string, map[string]string) (ownedMCPSession, error) {
		connects++
		return fake, nil
	}

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		status, err := owner.Status(ctx, environmentID, mcpID)
		if err != nil {
			t.Fatal(err)
		}
		if status.State != app.MCPHealthHealthy {
			t.Fatalf("status = %+v", status)
		}
	}
	if _, err := owner.ListTools(ctx, environmentID, mcpID); err != nil {
		t.Fatal(err)
	}
	if _, err := owner.CallTool(ctx, environmentID, mcpID, "owned_fake", nil); err != nil {
		t.Fatal(err)
	}
	if connects != 1 {
		t.Fatalf("connects = %d, want 1", connects)
	}
	if info := owner.Info(); info.OwnedMCPSessions != 1 || info.ID == "" || info.PID <= 0 {
		t.Fatalf("owner info = %+v", info)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
	_, _, closes := fake.counts()
	if closes != 1 {
		t.Fatalf("session closes = %d, want 1", closes)
	}
	if err := owner.Close(); err != nil {
		t.Fatal(err)
	}
	_, _, closes = fake.counts()
	if closes != 1 {
		t.Fatalf("idempotent close count = %d, want 1", closes)
	}
}

func TestRuntimeOwnerRestartRebuildsPersistedDesiredState(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)

	first := newRuntimeOwner(service)
	firstSession := &fakeOwnedMCPSession{}
	firstConnects := 0
	first.connect = func(context.Context, string, string, map[string]string) (ownedMCPSession, error) {
		firstConnects++
		return firstSession, nil
	}
	first.Reconcile(context.Background())
	if firstConnects != 1 || first.Info().OwnedMCPSessions != 1 {
		t.Fatalf("first owner connects=%d info=%+v", firstConnects, first.Info())
	}
	firstID := first.Info().ID
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	second := newRuntimeOwner(service)
	secondSession := &fakeOwnedMCPSession{}
	secondConnects := 0
	second.connect = func(context.Context, string, string, map[string]string) (ownedMCPSession, error) {
		secondConnects++
		return secondSession, nil
	}
	second.Reconcile(context.Background())
	defer second.Close()
	if secondConnects != 1 || second.Info().OwnedMCPSessions != 1 {
		t.Fatalf("second owner connects=%d info=%+v", secondConnects, second.Info())
	}
	if second.Info().ID == firstID {
		t.Fatalf("restart reused owner id %q", firstID)
	}
	status, err := second.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("rebuilt status=%+v err=%v", status, err)
	}
	if secondConnects != 1 {
		t.Fatalf("rebuilt owner reconnected unnecessarily: %d", secondConnects)
	}
}

func TestRuntimeOwnerDeadSessionCannotRemainHealthy(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	defer owner.Close()
	first := &fakeOwnedMCPSession{failAfter: 1, listErr: errors.New("connection refused")}
	connects := 0
	owner.connect = func(context.Context, string, string, map[string]string) (ownedMCPSession, error) {
		connects++
		if connects == 1 {
			return first, nil
		}
		return nil, errors.New("connection refused")
	}

	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	status, err = owner.Status(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if status.State != app.MCPHealthError || status.ErrorKind != "connection_refused" {
		t.Fatalf("dead session status=%+v", status)
	}
	if owner.Info().OwnedMCPSessions != 0 {
		t.Fatalf("dead session remained owned: %+v", owner.Info())
	}
	_, _, closes := first.counts()
	if closes != 1 {
		t.Fatalf("dead session closes=%d, want 1", closes)
	}
}

func TestRuntimeOwnerDisabledDesiredStateEvictsSession(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	defer owner.Close()
	fake := &fakeOwnedMCPSession{}
	owner.connect = func(context.Context, string, string, map[string]string) (ownedMCPSession, error) {
		return fake, nil
	}
	if status, err := owner.Status(context.Background(), environmentID, mcpID); err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	if _, err := service.SetEnvironmentMCP(environmentID, mcpID, false); err != nil {
		t.Fatal(err)
	}
	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if status.State != app.MCPHealthDisabled {
		t.Fatalf("disabled status=%+v", status)
	}
	if owner.Info().OwnedMCPSessions != 0 {
		t.Fatalf("disabled session remained owned: %+v", owner.Info())
	}
	_, _, closes := fake.counts()
	if closes != 1 {
		t.Fatalf("disabled session closes=%d, want 1", closes)
	}
}

func TestRuntimeOwnerSharedAcrossIndependentAgentSessions(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	defer owner.Close()
	fake := &fakeOwnedMCPSession{}
	connects := 0
	owner.connect = func(context.Context, string, string, map[string]string) (ownedMCPSession, error) {
		connects++
		return fake, nil
	}

	server := newServer(service, owner)
	ctx := context.Background()
	first := connectInMemory(t, ctx, server)
	defer first.Close()
	second := connectInMemory(t, ctx, server)
	defer second.Close()
	for _, session := range []*mcp.ClientSession{first, second} {
		status := callGatewayTool(t, ctx, session, "environment_mcp_status", map[string]any{
			"environment_id": environmentID,
			"mcp_id":         mcpID,
		})
		if status.IsError || !strings.Contains(toolText(t, status), "healthy") {
			t.Fatalf("owner status failed: %s", toolText(t, status))
		}
		info := callGatewayTool(t, ctx, session, "gateway_info", map[string]any{})
		if info.IsError || !strings.Contains(toolText(t, info), owner.Info().ID) {
			t.Fatalf("gateway_info missing shared owner %q: %s", owner.Info().ID, toolText(t, info))
		}
	}
	if connects != 1 {
		t.Fatalf("independent clients created %d owned sessions, want 1", connects)
	}
}

func runtimeOwnerTestService(t *testing.T) (*app.Service, string, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "marker.txt"), []byte("owner\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	workspace, err := service.Workspaces.Add(root, "owner-test")
	if err != nil {
		t.Fatal(err)
	}
	entry, err := service.MCPs.AddMCP("owned", "http://127.0.0.1:65534/mcp", false)
	if err != nil {
		t.Fatal(err)
	}
	environment, err := service.Environments.Create(workspace.ID, "owner-test", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.SetEnvironmentMCP(environment.ID, entry.ID, true); err != nil {
		t.Fatal(err)
	}
	return service, environment.ID, entry.ID
}
