package gateway

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type fakeOwnedMCPSession struct {
	mu            sync.Mutex
	pingCalls     int
	listCalls     int
	callCalls     int
	closes        int
	pingFailAfter int
	pingErr       error
	failAfter     int
	listErr       error
}

func (s *fakeOwnedMCPSession) Ping(context.Context, *mcp.PingParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pingCalls++
	if s.pingFailAfter > 0 && s.pingCalls > s.pingFailAfter {
		if s.pingErr != nil {
			return s.pingErr
		}
		return errors.New("connection refused")
	}
	return nil
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

func (s *fakeOwnedMCPSession) pingCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pingCalls
}

func TestRuntimeOwnerReusesSessionAndClosesIt(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	fake := &fakeOwnedMCPSession{}
	connects := 0
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
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
	first.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
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
	second.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
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

func TestRuntimeOwnerRestartPersistsHealthPolicyButNotObservation(t *testing.T) {
	policy := model.MCPHealthPolicy{
		HealthCheckEnabled:       true,
		CheckIntervalSeconds:     3,
		ProbeTimeoutSeconds:      2,
		AutoReconnect:            true,
		ReconnectIntervalSeconds: 4,
	}
	service, environmentID, mcpID := runtimeOwnerPolicyTestService(t, policy)

	first := newRuntimeOwner(service)
	firstSession := &fakeOwnedMCPSession{}
	first.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		return firstSession, nil
	}
	if status, err := first.Status(context.Background(), environmentID, mcpID); err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("first owner status=%+v err=%v", status, err)
	}
	firstInspection, err := first.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if firstInspection.Observation.LastCheckAt.IsZero() || firstInspection.Observation.InventoryFetchedAt.IsZero() || len(firstInspection.Observation.Inventory) != 1 {
		t.Fatalf("first owner did not record live observation: %+v", firstInspection.Observation)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	second := newRuntimeOwner(service)
	defer second.Close()
	secondInspection, err := second.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if got := secondInspection.Definition.HealthPolicy; got != policy {
		t.Fatalf("health policy did not persist: got=%+v want=%+v", got, policy)
	}
	if secondInspection.Observation.LastCheckAt != (time.Time{}) ||
		secondInspection.Observation.InventoryFetchedAt != (time.Time{}) ||
		len(secondInspection.Observation.Inventory) != 0 ||
		secondInspection.Observation.State != app.MCPHealthConfigured {
		t.Fatalf("owner-local observation leaked across restart: %+v", secondInspection.Observation)
	}

	secondSession := &fakeOwnedMCPSession{}
	connects := 0
	second.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		connects++
		return secondSession, nil
	}
	if status, err := second.Status(context.Background(), environmentID, mcpID); err != nil || status.State != app.MCPHealthHealthy || connects != 1 {
		t.Fatalf("second owner did not rebuild from desired state: status=%+v err=%v connects=%d", status, err, connects)
	}
}

func TestRuntimeOwnerDeadSessionCannotRemainHealthy(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	defer owner.Close()
	first := &fakeOwnedMCPSession{pingFailAfter: 1, pingErr: errors.New("connection refused")}
	connects := 0
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
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
	inspection, err := owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.FailureStage != "ping" || inspection.Observation.ConsecutiveFailures != 1 || len(inspection.Observation.Inventory) != 0 {
		t.Fatalf("dead session observation=%+v", inspection.Observation)
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
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
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
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
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

func TestRuntimeOwnerBackgroundMonitorSkipsPingWhenHealthCheckDisabled(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerPolicyTestService(t, model.MCPHealthPolicy{
		HealthCheckEnabled:       false,
		CheckIntervalSeconds:     1,
		ProbeTimeoutSeconds:      1,
		AutoReconnect:            true,
		ReconnectIntervalSeconds: 1,
	})
	owner := newRuntimeOwner(service)
	defer owner.Close()

	fake := &fakeOwnedMCPSession{pingFailAfter: 1, pingErr: errors.New("connection refused")}
	connects := 0
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		connects++
		return fake, nil
	}

	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	inspection, err := owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if connects != 1 || fake.pingCount() != 1 || owner.Info().OwnedMCPSessions != 1 {
		t.Fatalf("health_check_enabled=false must not background Ping/reconnect: connects=%d pings=%d info=%+v", connects, fake.pingCount(), owner.Info())
	}
	if inspection.Observation.State != app.MCPHealthHealthy || inspection.Observation.NextReconnectAt != nil || inspection.Observation.FailureStage != "" {
		t.Fatalf("health_check_enabled=false changed observation unexpectedly: %+v", inspection.Observation)
	}
}

func TestRuntimeOwnerBackgroundMonitorDoesNotReconnectWhenDisabled(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerPolicyTestService(t, model.MCPHealthPolicy{
		HealthCheckEnabled:       true,
		CheckIntervalSeconds:     1,
		ProbeTimeoutSeconds:      1,
		AutoReconnect:            false,
		ReconnectIntervalSeconds: 1,
	})
	owner := newRuntimeOwner(service)
	defer owner.Close()

	first := &fakeOwnedMCPSession{pingFailAfter: 1, pingErr: errors.New("connection refused")}
	second := &fakeOwnedMCPSession{}
	connects := 0
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		connects++
		if connects == 1 {
			return first, nil
		}
		return second, nil
	}

	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	inspection, err := owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.State != app.MCPHealthError ||
		inspection.Observation.FailureStage != "ping" ||
		inspection.Observation.ErrorKind != "connection_refused" ||
		inspection.Observation.ConsecutiveFailures != 1 ||
		inspection.Observation.NextReconnectAt != nil ||
		len(inspection.Observation.Inventory) != 0 ||
		!inspection.Observation.InventoryFetchedAt.IsZero() ||
		inspection.Observation.LastCheckAt.IsZero() {
		t.Fatalf("monitor failure observation=%+v", inspection.Observation)
	}
	if owner.Info().OwnedMCPSessions != 0 || connects != 1 {
		t.Fatalf("auto_reconnect=false should not reconnect in background: connects=%d info=%+v", connects, owner.Info())
	}

	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	if connects != 1 {
		t.Fatalf("background monitor reconnected despite auto_reconnect=false: connects=%d", connects)
	}
	status, err = owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy || connects != 2 {
		t.Fatalf("explicit status should safely reconnect: status=%+v err=%v connects=%d", status, err, connects)
	}
}

func TestRuntimeOwnerBackgroundMonitorReconnectsOnFixedInterval(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerPolicyTestService(t, model.MCPHealthPolicy{
		HealthCheckEnabled:       true,
		CheckIntervalSeconds:     1,
		ProbeTimeoutSeconds:      1,
		AutoReconnect:            true,
		ReconnectIntervalSeconds: 1,
	})
	owner := newRuntimeOwner(service)
	defer owner.Close()

	first := &fakeOwnedMCPSession{pingFailAfter: 1, pingErr: errors.New("connection refused")}
	second := &fakeOwnedMCPSession{}
	connects := 0
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		connects++
		if connects == 1 {
			return first, nil
		}
		return second, nil
	}

	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	inspection, err := owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.State != app.MCPHealthError ||
		inspection.Observation.FailureStage != "ping" ||
		inspection.Observation.ErrorKind != "connection_refused" ||
		inspection.Observation.ConsecutiveFailures != 1 ||
		inspection.Observation.NextReconnectAt == nil ||
		inspection.Observation.ProbeInFlight ||
		inspection.Observation.ReconnectInFlight ||
		connects != 1 {
		t.Fatalf("expected scheduled reconnect after ping failure: connects=%d observation=%+v", connects, inspection.Observation)
	}
	owner.MonitorOnce(context.Background())
	if connects != 1 {
		t.Fatalf("reconnect ran before configured interval: connects=%d", connects)
	}

	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	inspection, err = owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if connects != 2 || owner.Info().OwnedMCPSessions != 1 || inspection.Observation.State != app.MCPHealthHealthy || inspection.Observation.NextReconnectAt != nil {
		t.Fatalf("background reconnect did not restore health: connects=%d info=%+v observation=%+v", connects, owner.Info(), inspection.Observation)
	}
	_, firstCalls, _ := first.counts()
	_, secondCalls, _ := second.counts()
	if firstCalls != 0 || secondCalls != 0 {
		t.Fatalf("background reconnect must not replay MCP tool calls: first_calls=%d second_calls=%d", firstCalls, secondCalls)
	}
}

func TestRuntimeOwnerBackgroundReconnectDoesNotResurrectDisabledMCP(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerPolicyTestService(t, model.MCPHealthPolicy{
		HealthCheckEnabled:       true,
		CheckIntervalSeconds:     1,
		ProbeTimeoutSeconds:      1,
		AutoReconnect:            true,
		ReconnectIntervalSeconds: 1,
	})
	owner := newRuntimeOwner(service)
	defer owner.Close()

	first := &fakeOwnedMCPSession{pingFailAfter: 1, pingErr: errors.New("connection refused")}
	connects := 0
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		connects++
		return first, nil
	}

	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	inspection, err := owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.NextReconnectAt == nil {
		t.Fatalf("expected pending reconnect before disable: %+v", inspection.Observation)
	}
	if _, err := service.SetEnvironmentMCP(environmentID, mcpID, false); err != nil {
		t.Fatal(err)
	}
	time.Sleep(1100 * time.Millisecond)
	owner.MonitorOnce(context.Background())
	status, err = owner.Status(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if status.State != app.MCPHealthDisabled || connects != 1 || owner.Info().OwnedMCPSessions != 0 {
		t.Fatalf("disabled MCP was resurrected: status=%+v connects=%d info=%+v", status, connects, owner.Info())
	}
}

func TestRuntimeOwnerInspectAndRefreshObservation(t *testing.T) {
	service, environmentID, mcpID := runtimeOwnerTestService(t)
	owner := newRuntimeOwner(service)
	defer owner.Close()

	var sessions []*fakeOwnedMCPSession
	owner.connect = func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error) {
		session := &fakeOwnedMCPSession{}
		sessions = append(sessions, session)
		return session, nil
	}

	status, err := owner.Status(context.Background(), environmentID, mcpID)
	if err != nil || status.State != app.MCPHealthHealthy {
		t.Fatalf("initial status=%+v err=%v", status, err)
	}
	inspection, err := owner.Inspect(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if inspection.Observation.State != app.MCPHealthHealthy ||
		inspection.Observation.Transport != "streamable-http" ||
		len(inspection.Observation.Inventory) != 1 ||
		inspection.Observation.Inventory[0].Name != "owned_fake" ||
		inspection.Observation.LastCheckAt.IsZero() ||
		inspection.Observation.InventoryFetchedAt.IsZero() {
		t.Fatalf("initial inspection = %+v", inspection)
	}

	refreshed, err := owner.Refresh(context.Background(), environmentID, mcpID)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 2 || refreshed.Observation.State != app.MCPHealthHealthy ||
		len(refreshed.Observation.Inventory) != 1 {
		t.Fatalf("refreshed inspection=%+v sessions=%d", refreshed, len(sessions))
	}
	_, _, closes := sessions[0].counts()
	if closes != 1 {
		t.Fatalf("stale session closes=%d, want 1", closes)
	}
}

func runtimeOwnerTestService(t *testing.T) (*app.Service, string, string) {
	t.Helper()
	return runtimeOwnerPolicyTestService(t, model.MCPHealthPolicy{})
}

func runtimeOwnerPolicyTestService(t *testing.T, policy model.MCPHealthPolicy) (*app.Service, string, string) {
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
	entry, err := service.MCPs.AddMCPConfig("owned", catalog.MCPConfig{
		Transport:    catalog.MCPTransportStreamableHTTP,
		Endpoint:     "http://127.0.0.1:65534/mcp",
		HealthPolicy: policy,
	})
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
