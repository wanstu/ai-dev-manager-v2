package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPHealthState string

const (
	MCPHealthConfigured MCPHealthState = "configured"
	MCPHealthDisabled   MCPHealthState = "disabled"
	MCPHealthHealthy    MCPHealthState = "healthy"
	MCPHealthError      MCPHealthState = "error"
)

type MCPHealthStatus struct {
	MCPID     string         `json:"mcp_id"`
	State     MCPHealthState `json:"state"`
	ErrorKind string         `json:"error_kind,omitempty"`
	Message   string         `json:"message,omitempty"`
}

type MCPError struct {
	MCPID     string `json:"mcp_id"`
	ErrorKind string `json:"error_kind"`
	Message   string `json:"message"`
}

func (e MCPError) Error() string {
	return fmt.Sprintf("mcp_id=%s error_kind=%s message=%s", e.MCPID, e.ErrorKind, e.Message)
}

// MCPActivation contains resolved runtime-only connection material. It is never
// persisted or returned through ordinary management inspection/status output.
type MCPActivation struct {
	MCPID        string
	Transport    string
	Endpoint     string
	Headers      map[string]string
	Executable   string
	Args         []string
	Environment  map[string]string
	HealthPolicy model.MCPHealthPolicy
}

// ResolveMCPActivation keeps Environment selection as the authorization gate,
// then resolves secret/environment references only at this activation boundary.
func (s *Service) ResolveMCPActivation(environmentID, mcpID string) (*MCPActivation, MCPHealthStatus, error) {
	env, err := s.Environments.Get(environmentID)
	if err != nil {
		return nil, MCPHealthStatus{}, err
	}
	enabled := false
	for _, id := range env.EnabledMCPIDs {
		if id == mcpID {
			enabled = true
			break
		}
	}
	if !enabled {
		return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthDisabled}, nil
	}

	definition, err := s.MCPs.Get(mcpID)
	if err != nil {
		return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured, Message: "mcp definition is unresolved"}, nil
	}

	activation := &MCPActivation{
		MCPID:        mcpID,
		Transport:    definition.Transport,
		HealthPolicy: definition.HealthPolicy,
	}
	switch definition.Transport {
	case catalog.MCPTransportStreamableHTTP:
		if strings.TrimSpace(definition.Endpoint) == "" {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured, Message: "mcp has no configured endpoint"}, nil
		}
		if definition.AuthMode != catalog.MCPAuthNone && definition.AuthMode != catalog.MCPAuthHeaders {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthError, ErrorKind: "unsupported_auth", Message: "mcp auth mode is not supported by this runtime"}, nil
		}
		if hasUnresolvedEnvRef(definition.Endpoint) || hasUnresolvedStringMapRef(definition.HeaderRefs) {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured, Message: "mcp connection configuration is unresolved"}, nil
		}
		activation.Endpoint = os.ExpandEnv(definition.Endpoint)
		activation.Headers = resolveStringMap(definition.HeaderRefs)
	case catalog.MCPTransportStdio:
		if strings.TrimSpace(definition.Executable) == "" {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured, Message: "mcp has no configured executable"}, nil
		}
		if definition.AuthMode != catalog.MCPAuthNone {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthError, ErrorKind: "unsupported_auth", Message: "mcp auth mode is not supported by this transport"}, nil
		}
		if hasUnresolvedStringMapRef(definition.EnvRefs) {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured, Message: "mcp environment configuration is unresolved"}, nil
		}
		activation.Executable = definition.Executable
		activation.Args = append([]string(nil), definition.Args...)
		activation.Environment = resolveStringMap(definition.EnvRefs)
	default:
		return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthError, ErrorKind: "unsupported_transport", Message: "mcp transport is not supported by this runtime"}, nil
	}
	return activation, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured}, nil
}

// ConnectMCP activates either supported transport. Stdio command preparation
// reuses Runtime.PrepareCommand so executable allowlisting, Environment-root cwd,
// containment and OS cancellation policy stay under the existing authority.
func (s *Service) ConnectMCP(ctx context.Context, environmentID string, activation *MCPActivation, implementation *mcp.Implementation) (*mcp.ClientSession, error) {
	if activation == nil {
		return nil, fmt.Errorf("mcp activation is required")
	}
	if implementation == nil {
		implementation = &mcp.Implementation{Name: "adm-v2-mcp-client", Version: "dev"}
	}
	client := mcp.NewClient(implementation, nil)

	var transport mcp.Transport
	switch activation.Transport {
	case catalog.MCPTransportStreamableHTTP:
		httpTransport := &mcp.StreamableClientTransport{
			Endpoint:             activation.Endpoint,
			MaxRetries:           -1,
			DisableStandaloneSSE: true,
		}
		if len(activation.Headers) != 0 {
			httpTransport.HTTPClient = &http.Client{Transport: mcpHeaderRoundTripper{base: http.DefaultTransport, headers: activation.Headers}}
		}
		transport = httpTransport
	case catalog.MCPTransportStdio:
		rt, _, err := s.Runtime(environmentID)
		if err != nil {
			return nil, err
		}
		cmd, err := rt.PrepareCommand(ctx, activation.Executable, activation.Args, "")
		if err != nil {
			return nil, err
		}
		cmd.Env = append(os.Environ(), environmentEntries(activation.Environment)...)
		transport = &mcp.CommandTransport{Command: cmd}
	default:
		return nil, fmt.Errorf("unsupported mcp transport %q", activation.Transport)
	}
	return client.Connect(ctx, transport, nil)
}

// ProbeMCPHealth performs one bounded on-demand probe for management callers.
// Phase 11-02 adds owner-local protocol Ping monitoring/recovery; this explicit
// probe remains useful even when background health checks are disabled.
func (s *Service) ProbeMCPHealth(ctx context.Context, environmentID, mcpID string) (MCPHealthStatus, error) {
	activation, status, err := s.ResolveMCPActivation(environmentID, mcpID)
	if err != nil {
		return MCPHealthStatus{}, err
	}
	if activation == nil {
		return status, nil
	}

	probeSeconds := activation.HealthPolicy.ProbeTimeoutSeconds
	if probeSeconds <= 0 {
		probeSeconds = 10
	}
	probeCtx, cancel := context.WithTimeout(ctx, time.Duration(probeSeconds)*time.Second)
	defer cancel()

	session, err := s.ConnectMCP(probeCtx, environmentID, activation, &mcp.Implementation{Name: "adm-v2-health-probe", Version: "dev"})
	if err != nil {
		return mcpHealthErrorStatus(mcpID, classifyMCPError(err)), nil
	}
	defer session.Close()

	if err := session.Ping(probeCtx, &mcp.PingParams{}); err != nil {
		kind := classifyMCPError(err)
		if kind == "timeout" {
			kind = "ping_timeout"
		} else if kind == "connection_failed" {
			kind = "ping_protocol_failure"
		}
		return mcpHealthErrorStatus(mcpID, kind), nil
	}
	return MCPHealthStatus{MCPID: mcpID, State: MCPHealthHealthy}, nil
}

func resolveStringMap(refs map[string]string) map[string]string {
	if len(refs) == 0 {
		return nil
	}
	resolved := make(map[string]string, len(refs))
	for key, value := range refs {
		resolved[key] = os.ExpandEnv(value)
	}
	return resolved
}

func hasUnresolvedEnvRef(value string) bool {
	unresolved := false
	os.Expand(value, func(key string) string {
		if _, ok := os.LookupEnv(key); !ok {
			unresolved = true
		}
		return ""
	})
	return unresolved
}

func hasUnresolvedStringMapRef(refs map[string]string) bool {
	for _, value := range refs {
		if hasUnresolvedEnvRef(value) {
			return true
		}
	}
	return false
}

func environmentEntries(values map[string]string) []string {
	if len(values) == 0 {
		return nil
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	entries := make([]string, 0, len(keys))
	for _, key := range keys {
		entries = append(entries, key+"="+values[key])
	}
	return entries
}

func classifyMCPError(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return "timeout"
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return "timeout"
	}
	text := strings.ToLower(err.Error())
	switch {
	case strings.Contains(text, "not allowlisted"), strings.Contains(text, "not allowed"):
		return "executable_not_allowed"
	case strings.Contains(text, "deadline exceeded"), strings.Contains(text, "i/o timeout"), strings.Contains(text, "timed out"):
		return "timeout"
	case strings.Contains(text, "connection refused"), strings.Contains(text, "actively refused"):
		return "connection_refused"
	case strings.Contains(text, "401"), strings.Contains(text, "403"), strings.Contains(text, "unauthorized"), strings.Contains(text, "forbidden"):
		return "auth_failure"
	default:
		return "connection_failed"
	}
}

// ClassifyMCPError exposes the app-owned classifier to thin Gateway handlers.
func ClassifyMCPError(err error) string { return classifyMCPError(err) }

func mcpHealthErrorStatus(mcpID, kind string) MCPHealthStatus {
	return MCPHealthStatus{
		MCPID:     mcpID,
		State:     MCPHealthError,
		ErrorKind: kind,
		Message:   fmt.Sprintf("mcp health probe failed: %s", kind),
	}
}

type mcpHeaderRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

func (r mcpHeaderRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	clone.Header = request.Header.Clone()
	for key, value := range r.headers {
		clone.Header.Set(key, value)
	}
	return r.base.RoundTrip(clone)
}
