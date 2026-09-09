package gateway

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"ai-dev-manager-v2/internal/app"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	runtimeOwnerReconcileTimeout = 10 * time.Second
	runtimeOwnerMonitorInterval  = time.Second
)

var runtimeOwnerSequence atomic.Uint64

type ownedMCPSession interface {
	Ping(context.Context, *mcp.PingParams) error
	ListTools(context.Context, *mcp.ListToolsParams) (*mcp.ListToolsResult, error)
	CallTool(context.Context, *mcp.CallToolParams) (*mcp.CallToolResult, error)
	Close() error
}

type ownedMCPConnectFunc func(context.Context, string, *app.MCPActivation) (ownedMCPSession, error)

type runtimeOwnerKey struct {
	environmentID string
	mcpID         string
}

type runtimeOwnerInfo struct {
	ID                string    `json:"id"`
	PID               int       `json:"pid"`
	StartedAt         time.Time `json:"started_at"`
	OwnedMCPSessions  int       `json:"owned_mcp_sessions"`
	OwnedDevProcesses int       `json:"owned_dev_processes"`
	OwnedAgentRuns    int       `json:"owned_agent_runs"`
}

type runtimeOwner struct {
	service   *app.Service
	id        string
	pid       int
	startedAt time.Time
	connect   ownedMCPConnectFunc
	ctx       context.Context
	cancel    context.CancelFunc

	mu           sync.Mutex
	closed       bool
	sessions     map[runtimeOwnerKey]ownedMCPSession
	observed     map[runtimeOwnerKey]app.MCPHealthStatus
	observations map[runtimeOwnerKey]MCPRuntimeObservation
	processes    map[string]*ownedDevProcess
	runs         map[string]*ownedAgentRun
}

func newRuntimeOwner(service *app.Service) *runtimeOwner {
	startedAt := time.Now().UTC()
	ownerCtx, cancel := context.WithCancel(context.Background())
	owner := &runtimeOwner{
		service:      service,
		id:           fmt.Sprintf("owner_%d_%x_%x", os.Getpid(), startedAt.UnixNano(), runtimeOwnerSequence.Add(1)),
		pid:          os.Getpid(),
		startedAt:    startedAt,
		ctx:          ownerCtx,
		cancel:       cancel,
		sessions:     map[runtimeOwnerKey]ownedMCPSession{},
		observed:     map[runtimeOwnerKey]app.MCPHealthStatus{},
		observations: map[runtimeOwnerKey]MCPRuntimeObservation{},
		processes:    map[string]*ownedDevProcess{},
		runs:         map[string]*ownedAgentRun{},
	}
	owner.connect = func(ctx context.Context, environmentID string, activation *app.MCPActivation) (ownedMCPSession, error) {
		return service.ConnectMCP(ctx, environmentID, activation, &mcp.Implementation{Name: serverName + "-proxy", Version: serverVersion})
	}
	return owner
}

func (o *runtimeOwner) Info() runtimeOwnerInfo {
	o.mu.Lock()
	defer o.mu.Unlock()
	return runtimeOwnerInfo{
		ID:                o.id,
		PID:               o.pid,
		StartedAt:         o.startedAt,
		OwnedMCPSessions:  len(o.sessions),
		OwnedDevProcesses: o.runningDevProcessCountLocked(),
		OwnedAgentRuns:    o.runningAgentRunCountLocked(),
	}
}

func (o *runtimeOwner) Monitor(ctx context.Context) {
	if o == nil || o.service == nil {
		return
	}
	o.Reconcile(ctx)
	ticker := time.NewTicker(runtimeOwnerMonitorInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-o.operationContext(ctx).Done():
			return
		case <-ticker.C:
			o.MonitorOnce(ctx)
		}
	}
}

func (o *runtimeOwner) Reconcile(ctx context.Context) {
	if o == nil || o.service == nil {
		return
	}
	environments, err := o.service.Environments.List()
	if err != nil {
		return
	}
	desired := make(map[runtimeOwnerKey]struct{})
	for _, env := range environments {
		for _, mcpID := range env.EnabledMCPIDs {
			key := runtimeOwnerKey{environmentID: env.ID, mcpID: mcpID}
			desired[key] = struct{}{}
			probeCtx, cancel := context.WithTimeout(ctx, runtimeOwnerReconcileTimeout)
			_, _ = o.Status(probeCtx, env.ID, mcpID)
			cancel()
		}
	}
	o.pruneToDesired(desired)
}

func (o *runtimeOwner) MonitorOnce(ctx context.Context) {
	if o == nil || o.service == nil {
		return
	}
	environments, err := o.service.Environments.List()
	if err != nil {
		return
	}
	desired := make(map[runtimeOwnerKey]struct{})
	for _, env := range environments {
		for _, mcpID := range env.EnabledMCPIDs {
			key := runtimeOwnerKey{environmentID: env.ID, mcpID: mcpID}
			desired[key] = struct{}{}
			o.monitorMCP(ctx, key)
		}
	}
	o.pruneToDesired(desired)
}

func (o *runtimeOwner) monitorMCP(ctx context.Context, key runtimeOwnerKey) {
	activation, status, err := o.service.ResolveMCPActivation(key.environmentID, key.mcpID)
	if err != nil {
		return
	}
	if activation == nil {
		o.drop(key)
		o.observe(key, status, "", "activation", nil, false)
		return
	}
	if session := o.session(key); session != nil {
		if activation.HealthPolicy.HealthCheckEnabled && o.probeDue(key, activation.HealthPolicy.CheckIntervalSeconds) {
			o.monitorPing(ctx, key, activation, session)
		}
		return
	}
	if activation.HealthPolicy.AutoReconnect && o.reconnectDue(key) {
		o.monitorReconnect(ctx, key, activation)
	}
}

func (o *runtimeOwner) monitorPing(ctx context.Context, key runtimeOwnerKey, activation *app.MCPActivation, session ownedMCPSession) {
	if !o.markProbeInFlight(key) {
		return
	}
	probeCtx, cancel := mcpProbeContext(o.operationContext(ctx), activation.HealthPolicy.ProbeTimeoutSeconds)
	defer cancel()
	if err := session.Ping(probeCtx, &mcp.PingParams{}); err != nil {
		o.drop(key)
		status := ownerMCPErrorStatus(key.mcpID, pingMCPErrorKind(err))
		o.observe(key, status, activation.Transport, "ping", nil, false)
		o.scheduleReconnect(key, activation.HealthPolicy.AutoReconnect, activation.HealthPolicy.ReconnectIntervalSeconds)
		return
	}
	status := app.MCPHealthStatus{MCPID: key.mcpID, State: app.MCPHealthHealthy}
	o.observe(key, status, activation.Transport, "", nil, true)
}

func (o *runtimeOwner) monitorReconnect(ctx context.Context, key runtimeOwnerKey, activation *app.MCPActivation) {
	if !o.markReconnectInFlight(key) {
		return
	}
	_, status, err := o.ensureHealthySession(ctx, key.environmentID, key.mcpID)
	if err != nil {
		status = ownerMCPErrorStatus(key.mcpID, runtimeMCPErrorKind(err))
		o.observe(key, status, activation.Transport, "reconnect", nil, false)
		o.scheduleReconnect(key, true, activation.HealthPolicy.ReconnectIntervalSeconds)
		return
	}
	if status.State != app.MCPHealthHealthy {
		o.scheduleReconnect(key, true, activation.HealthPolicy.ReconnectIntervalSeconds)
	}
}

func (o *runtimeOwner) Status(ctx context.Context, environmentID, mcpID string) (app.MCPHealthStatus, error) {
	_, status, err := o.ensureHealthySession(ctx, environmentID, mcpID)
	return status, err
}

func (o *runtimeOwner) ListTools(ctx context.Context, environmentID, mcpID string) ([]*mcp.Tool, error) {
	session, status, err := o.ensureHealthySession(ctx, environmentID, mcpID)
	if err != nil {
		return nil, err
	}
	if status.State != app.MCPHealthHealthy || session == nil {
		return nil, statusAsMCPError(status)
	}
	operationCtx := o.operationContext(ctx)
	params := &mcp.ListToolsParams{}
	var tools []*mcp.Tool
	for {
		page, err := session.ListTools(operationCtx, params)
		if err != nil {
			o.drop(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID})
			kind := runtimeMCPErrorKind(err)
			if kind == "connection_failed" {
				kind = "tool_list_failed"
			}
			return nil, &app.MCPError{MCPID: mcpID, ErrorKind: kind, Message: "external MCP tool listing failed"}
		}
		tools = append(tools, page.Tools...)
		if page.NextCursor == "" {
			break
		}
		params.Cursor = page.NextCursor
	}
	key := runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID}
	status = app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy}
	o.setObserved(key, status)
	o.observe(key, status, "", "", tools, true)
	return tools, nil
}

func (o *runtimeOwner) CallTool(ctx context.Context, environmentID, mcpID, tool string, arguments map[string]any) (*mcp.CallToolResult, error) {
	session, status, err := o.ensureHealthySession(ctx, environmentID, mcpID)
	if err != nil {
		return nil, err
	}
	if status.State != app.MCPHealthHealthy || session == nil {
		return nil, statusAsMCPError(status)
	}
	operationCtx := o.operationContext(ctx)
	result, err := session.CallTool(operationCtx, &mcp.CallToolParams{Name: tool, Arguments: arguments})
	if err != nil {
		o.drop(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID})
		kind := runtimeMCPErrorKind(err)
		if kind == "connection_failed" {
			kind = "tool_call_failed"
		}
		return nil, &app.MCPError{MCPID: mcpID, ErrorKind: kind, Message: "external MCP tool call failed"}
	}
	key := runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID}
	status = app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy}
	o.setObserved(key, status)
	o.observe(key, status, "", "", nil, true)
	return result, nil
}

func (o *runtimeOwner) DropEnvironment(environmentID string) {
	if o == nil {
		return
	}
	o.dropDevProcessesForEnvironment(environmentID)
	o.dropAgentRunsForEnvironment(environmentID)
	o.dropMatching(func(key runtimeOwnerKey) bool { return key.environmentID == environmentID })
}

func (o *runtimeOwner) DropMCP(mcpID string) {
	if o == nil {
		return
	}
	o.dropMatching(func(key runtimeOwnerKey) bool { return key.mcpID == mcpID })
}

func (o *runtimeOwner) Drop(environmentID, mcpID string) {
	if o == nil {
		return
	}
	o.drop(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID})
}

func (o *runtimeOwner) Close() error {
	if o == nil {
		return nil
	}
	o.mu.Lock()
	if o.closed {
		o.mu.Unlock()
		return nil
	}
	o.closed = true
	if o.cancel != nil {
		o.cancel()
	}
	sessions := make([]ownedMCPSession, 0, len(o.sessions))
	for _, session := range o.sessions {
		sessions = append(sessions, session)
	}
	processes := make([]*ownedDevProcess, 0, len(o.processes))
	for _, process := range o.processes {
		processes = append(processes, process)
	}
	runs := make([]*ownedAgentRun, 0, len(o.runs))
	for _, run := range o.runs {
		runs = append(runs, run)
	}
	o.sessions = map[runtimeOwnerKey]ownedMCPSession{}
	o.observed = map[runtimeOwnerKey]app.MCPHealthStatus{}
	o.observations = map[runtimeOwnerKey]MCPRuntimeObservation{}
	o.processes = map[string]*ownedDevProcess{}
	o.runs = map[string]*ownedAgentRun{}
	o.mu.Unlock()

	var errs []error
	for _, session := range sessions {
		if err := session.Close(); err != nil && !isBenignOwnedMCPCloseError(err) {
			errs = append(errs, err)
		}
	}
	if err := o.closeDevProcesses(processes); err != nil {
		errs = append(errs, err)
	}
	if err := o.closeAgentRuns(runs); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func (o *runtimeOwner) ensureHealthySession(ctx context.Context, environmentID, mcpID string) (ownedMCPSession, app.MCPHealthStatus, error) {
	if o == nil || o.service == nil {
		return nil, app.MCPHealthStatus{}, fmt.Errorf("runtime owner is not initialized")
	}
	key := runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID}
	activation, status, err := o.service.ResolveMCPActivation(environmentID, mcpID)
	if err != nil {
		return nil, app.MCPHealthStatus{}, err
	}
	if activation == nil {
		o.drop(key)
		o.observe(key, status, "", "activation", nil, false)
		return nil, status, nil
	}

	baseCtx := o.ctx
	if baseCtx == nil {
		baseCtx = ctx
	}
	probeCtx, cancelProbe := mcpProbeContext(baseCtx, activation.HealthPolicy.ProbeTimeoutSeconds)
	defer cancelProbe()
	if session := o.session(key); session != nil {
		if err := session.Ping(probeCtx, &mcp.PingParams{}); err == nil {
			status = app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy}
			o.observe(key, status, activation.Transport, "", nil, true)
			return session, status, nil
		} else if pingMCPErrorKind(err) == "ping_protocol_failure" {
			o.drop(key)
		} else {
			o.drop(key)
			status = ownerMCPErrorStatus(mcpID, pingMCPErrorKind(err))
			o.observe(key, status, activation.Transport, "ping", nil, false)
			return nil, status, nil
		}
	}

	session, err := o.connect(baseCtx, environmentID, activation)
	if err != nil {
		status = ownerMCPErrorStatus(mcpID, runtimeMCPErrorKind(err))
		o.observe(key, status, activation.Transport, "connect", nil, false)
		return nil, status, nil
	}
	if err := session.Ping(probeCtx, &mcp.PingParams{}); err != nil && pingMCPErrorKind(err) != "ping_protocol_failure" {
		_ = session.Close()
		status = ownerMCPErrorStatus(mcpID, pingMCPErrorKind(err))
		o.observe(key, status, activation.Transport, "ping", nil, false)
		return nil, status, nil
	}
	inventory, err := discoverOwnedMCPTools(probeCtx, session)
	if err != nil {
		_ = session.Close()
		kind := runtimeMCPErrorKind(err)
		if kind == "connection_failed" {
			kind = "tool_list_failed"
		}
		status = ownerMCPErrorStatus(mcpID, kind)
		o.observe(key, status, activation.Transport, "list_tools", nil, false)
		return nil, status, nil
	}
	current := o.installSession(key, session)
	if current == nil {
		_ = session.Close()
		status = ownerMCPErrorStatus(mcpID, "owner_closed")
		o.observe(key, status, activation.Transport, "connect", nil, false)
		return nil, status, nil
	}
	if current != session {
		_ = session.Close()
		session = current
	}
	status = app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy}
	o.observe(key, status, activation.Transport, "", inventory, true)
	return session, status, nil
}

func (o *runtimeOwner) session(key runtimeOwnerKey) ownedMCPSession {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	return o.sessions[key]
}

func (o *runtimeOwner) installSession(key runtimeOwnerKey, session ownedMCPSession) ownedMCPSession {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	if current := o.sessions[key]; current != nil {
		return current
	}
	o.sessions[key] = session
	return session
}

func (o *runtimeOwner) setObserved(key runtimeOwnerKey, status app.MCPHealthStatus) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return
	}
	o.observed[key] = status
}

func (o *runtimeOwner) drop(key runtimeOwnerKey) {
	var session ownedMCPSession
	o.mu.Lock()
	if !o.closed {
		session = o.sessions[key]
		delete(o.sessions, key)
		delete(o.observed, key)
		delete(o.observations, key)
	}
	o.mu.Unlock()
	if session != nil {
		_ = session.Close()
	}
}

func (o *runtimeOwner) dropMatching(match func(runtimeOwnerKey) bool) {
	var sessions []ownedMCPSession
	o.mu.Lock()
	if !o.closed {
		for key, session := range o.sessions {
			if !match(key) {
				continue
			}
			sessions = append(sessions, session)
			delete(o.sessions, key)
			delete(o.observed, key)
			delete(o.observations, key)
		}
	}
	o.mu.Unlock()
	for _, session := range sessions {
		_ = session.Close()
	}
}

func (o *runtimeOwner) pruneToDesired(desired map[runtimeOwnerKey]struct{}) {
	o.dropMatching(func(key runtimeOwnerKey) bool {
		_, ok := desired[key]
		return !ok
	})
}

func (o *runtimeOwner) observedStatuses() []app.MCPHealthStatus {
	o.mu.Lock()
	defer o.mu.Unlock()
	statuses := make([]app.MCPHealthStatus, 0, len(o.observed))
	for _, status := range o.observed {
		statuses = append(statuses, status)
	}
	sort.Slice(statuses, func(i, j int) bool { return statuses[i].MCPID < statuses[j].MCPID })
	return statuses
}

func (o *runtimeOwner) probeDue(key runtimeOwnerKey, intervalSeconds int64) bool {
	if intervalSeconds <= 0 {
		intervalSeconds = 30
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return false
	}
	observation := o.observations[key]
	if observation.ProbeInFlight || observation.ReconnectInFlight {
		return false
	}
	if observation.LastCheckAt.IsZero() {
		return true
	}
	return time.Since(observation.LastCheckAt) >= time.Duration(intervalSeconds)*time.Second
}

func (o *runtimeOwner) reconnectDue(key runtimeOwnerKey) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return false
	}
	observation := o.observations[key]
	if observation.ProbeInFlight || observation.ReconnectInFlight {
		return false
	}
	return observation.NextReconnectAt == nil || !time.Now().UTC().Before(*observation.NextReconnectAt)
}

func (o *runtimeOwner) markProbeInFlight(key runtimeOwnerKey) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return false
	}
	observation := o.observations[key]
	if observation.ProbeInFlight || observation.ReconnectInFlight {
		return false
	}
	observation.EnvironmentID = key.environmentID
	observation.MCPID = key.mcpID
	observation.ProbeInFlight = true
	o.observations[key] = observation
	return true
}

func (o *runtimeOwner) markReconnectInFlight(key runtimeOwnerKey) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return false
	}
	observation := o.observations[key]
	if observation.ProbeInFlight || observation.ReconnectInFlight {
		return false
	}
	observation.EnvironmentID = key.environmentID
	observation.MCPID = key.mcpID
	observation.ReconnectInFlight = true
	o.observations[key] = observation
	return true
}

func (o *runtimeOwner) scheduleReconnect(key runtimeOwnerKey, enabled bool, intervalSeconds int64) {
	if !enabled {
		return
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 30
	}
	next := time.Now().UTC().Add(time.Duration(intervalSeconds) * time.Second)
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return
	}
	observation := o.observations[key]
	observation.EnvironmentID = key.environmentID
	observation.MCPID = key.mcpID
	observation.NextReconnectAt = &next
	observation.ProbeInFlight = false
	observation.ReconnectInFlight = false
	o.observations[key] = observation
}

func (o *runtimeOwner) operationContext(fallback context.Context) context.Context {
	if o != nil && o.ctx != nil {
		return o.ctx
	}
	return fallback
}

func isBenignOwnedMCPCloseError(err error) bool {
	if err == nil {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "canceling cmd: invalid argument")
}

func ownerMCPErrorStatus(mcpID, kind string) app.MCPHealthStatus {
	if kind == "" {
		kind = "connection_failed"
	}
	return app.MCPHealthStatus{
		MCPID:     mcpID,
		State:     app.MCPHealthError,
		ErrorKind: kind,
		Message:   fmt.Sprintf("mcp health probe failed: %s", kind),
	}
}

func runtimeMCPErrorKind(err error) string {
	if err == nil {
		return ""
	}
	var mcpErr *app.MCPError
	if errors.As(err, &mcpErr) && mcpErr.ErrorKind != "" {
		return mcpErr.ErrorKind
	}
	return app.ClassifyMCPError(err)
}

func statusAsMCPError(status app.MCPHealthStatus) error {
	kind := status.ErrorKind
	if kind == "" {
		switch status.State {
		case app.MCPHealthDisabled:
			kind = "not_enabled"
		case app.MCPHealthConfigured:
			kind = "configured"
		default:
			kind = string(status.State)
		}
	}
	message := status.Message
	if message == "" {
		message = fmt.Sprintf("external MCP is %s", status.State)
	}
	return &app.MCPError{MCPID: status.MCPID, ErrorKind: kind, Message: message}
}
