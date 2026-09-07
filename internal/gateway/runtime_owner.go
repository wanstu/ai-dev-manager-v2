package gateway

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"ai-dev-manager-v2/internal/app"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const runtimeOwnerReconcileTimeout = 10 * time.Second

var runtimeOwnerSequence atomic.Uint64

type ownedMCPSession interface {
	ListTools(context.Context, *mcp.ListToolsParams) (*mcp.ListToolsResult, error)
	CallTool(context.Context, *mcp.CallToolParams) (*mcp.CallToolResult, error)
	Close() error
}

type ownedMCPConnectFunc func(context.Context, string, string, map[string]string) (ownedMCPSession, error)

type runtimeOwnerKey struct {
	environmentID string
	mcpID         string
}

type runtimeOwnerInfo struct {
	ID               string    `json:"id"`
	PID              int       `json:"pid"`
	StartedAt        time.Time `json:"started_at"`
	OwnedMCPSessions int       `json:"owned_mcp_sessions"`
}

type runtimeOwner struct {
	service   *app.Service
	id        string
	pid       int
	startedAt time.Time
	connect   ownedMCPConnectFunc

	mu       sync.Mutex
	closed   bool
	sessions map[runtimeOwnerKey]ownedMCPSession
	observed map[runtimeOwnerKey]app.MCPHealthStatus
}

func newRuntimeOwner(service *app.Service) *runtimeOwner {
	startedAt := time.Now().UTC()
	owner := &runtimeOwner{
		service:   service,
		id:        fmt.Sprintf("owner_%d_%x_%x", os.Getpid(), startedAt.UnixNano(), runtimeOwnerSequence.Add(1)),
		pid:       os.Getpid(),
		startedAt: startedAt,
		sessions:  map[runtimeOwnerKey]ownedMCPSession{},
		observed:  map[runtimeOwnerKey]app.MCPHealthStatus{},
	}
	owner.connect = func(ctx context.Context, mcpID, endpoint string, headers map[string]string) (ownedMCPSession, error) {
		return connectExternalMCP(ctx, mcpID, endpoint, headers)
	}
	return owner
}

func (o *runtimeOwner) Info() runtimeOwnerInfo {
	o.mu.Lock()
	defer o.mu.Unlock()
	return runtimeOwnerInfo{
		ID:               o.id,
		PID:              o.pid,
		StartedAt:        o.startedAt,
		OwnedMCPSessions: len(o.sessions),
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
	params := &mcp.ListToolsParams{}
	var tools []*mcp.Tool
	for {
		page, err := session.ListTools(ctx, params)
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
	o.setObserved(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID}, app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy})
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
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: tool, Arguments: arguments})
	if err != nil {
		o.drop(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID})
		kind := runtimeMCPErrorKind(err)
		if kind == "connection_failed" {
			kind = "tool_call_failed"
		}
		return nil, &app.MCPError{MCPID: mcpID, ErrorKind: kind, Message: "external MCP tool call failed"}
	}
	o.setObserved(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID}, app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy})
	return result, nil
}

func (o *runtimeOwner) DropEnvironment(environmentID string) {
	if o == nil {
		return
	}
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
	sessions := make([]ownedMCPSession, 0, len(o.sessions))
	for _, session := range o.sessions {
		sessions = append(sessions, session)
	}
	o.sessions = map[runtimeOwnerKey]ownedMCPSession{}
	o.observed = map[runtimeOwnerKey]app.MCPHealthStatus{}
	o.mu.Unlock()

	var errs []error
	for _, session := range sessions {
		if err := session.Close(); err != nil {
			errs = append(errs, err)
		}
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
		o.setObserved(key, status)
		return nil, status, nil
	}

	if session := o.session(key); session != nil {
		if _, err := session.ListTools(ctx, nil); err == nil {
			status = app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy}
			o.setObserved(key, status)
			return session, status, nil
		}
		o.drop(key)
	}

	session, err := o.connect(ctx, mcpID, activation.Endpoint, activation.Headers)
	if err != nil {
		status = ownerMCPErrorStatus(mcpID, runtimeMCPErrorKind(err))
		o.setObserved(key, status)
		return nil, status, nil
	}
	if _, err := session.ListTools(ctx, nil); err != nil {
		_ = session.Close()
		kind := runtimeMCPErrorKind(err)
		if kind == "connection_failed" {
			kind = "tool_list_failed"
		}
		status = ownerMCPErrorStatus(mcpID, kind)
		o.setObserved(key, status)
		return nil, status, nil
	}
	current := o.installSession(key, session)
	if current == nil {
		_ = session.Close()
		status = ownerMCPErrorStatus(mcpID, "owner_closed")
		return nil, status, nil
	}
	if current != session {
		_ = session.Close()
		session = current
	}
	status = app.MCPHealthStatus{MCPID: mcpID, State: app.MCPHealthHealthy}
	o.setObserved(key, status)
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
