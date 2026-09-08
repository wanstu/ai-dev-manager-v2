package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"ai-dev-manager-v2/internal/catalog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const mcpHealthProbeTimeout = 10 * time.Second

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

// MCPActivation is resolved only at a runtime activation boundary. Endpoint
// and header values may contain secrets and must never be persisted or exposed
// through ordinary inspection/status output.
type MCPActivation struct {
	MCPID      string
	Transport  string
	Endpoint   string
	Headers    map[string]string
	Executable string
	Args       []string
	Env        map[string]string
}

// ResolveMCPActivation separates persisted desired configuration from observed
// runtime state. A nil activation means the MCP is disabled, incomplete, or
// otherwise not activatable; status explains that desired/configuration state.
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

	entry, err := s.MCPs.Get(mcpID)
	if err != nil {
		return nil, MCPHealthStatus{
			MCPID:   mcpID,
			State:   MCPHealthConfigured,
			Message: "mcp catalog entry is unresolved",
		}, nil
	}
	if hasUnresolvedEnvRef(entry.Endpoint) || hasUnresolvedMapRef(entry.HeaderRefs) || hasUnresolvedMapRef(entry.EnvRefs) {
		return nil, MCPHealthStatus{
			MCPID:     mcpID,
			State:     MCPHealthConfigured,
			ErrorKind: "unresolved_secret_reference",
			Message:   "mcp connection configuration has an unresolved environment reference",
		}, nil
	}
	activation := &MCPActivation{
		MCPID:      mcpID,
		Transport:  entry.Transport,
		Endpoint:   os.ExpandEnv(entry.Endpoint),
		Headers:    resolveMap(entry.HeaderRefs),
		Executable: entry.Executable,
		Args:       append([]string(nil), entry.Args...),
		Env:        resolveMap(entry.EnvRefs),
	}
	if entry.Transport == catalog.MCPTransportStdio {
		rt, _, err := s.Runtime(environmentID)
		if err != nil {
			return nil, MCPHealthStatus{}, err
		}
		_, err = rt.Command(context.Background(), entry.Executable, entry.Args, activation.Env)
		if err != nil {
			return nil, MCPHealthStatus{MCPID: mcpID, State: MCPHealthError, ErrorKind: "executable_not_allowed", Message: "stdio MCP executable is unavailable under Environment authority"}, nil
		}
	}
	return activation, MCPHealthStatus{MCPID: mcpID, State: MCPHealthConfigured}, nil
}

// ProbeMCPHealth performs one bounded on-demand probe for management callers.
// The long-lived Gateway uses the same activation resolution but owns sessions
// in its runtime owner instead of persisting or reusing this transient probe.
func (s *Service) ProbeMCPHealth(ctx context.Context, environmentID, mcpID string) (MCPHealthStatus, error) {
	activation, status, err := s.ResolveMCPActivation(environmentID, mcpID)
	if err != nil {
		return MCPHealthStatus{}, err
	}
	if activation == nil {
		return status, nil
	}

	definition, _ := s.MCPs.Get(mcpID)
	timeout := mcpHealthProbeTimeout
	if definition.HealthPolicy.ProbeTimeoutSeconds > 0 {
		timeout = time.Duration(definition.HealthPolicy.ProbeTimeoutSeconds) * time.Second
	}
	probeCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client := mcp.NewClient(&mcp.Implementation{Name: "adm-v2-health-probe", Version: "dev"}, nil)
	transport, err := s.mcpTransport(probeCtx, environmentID, activation)
	if err != nil {
		return mcpHealthErrorStatus(mcpID, "activation_failed"), nil
	}
	session, err := client.Connect(probeCtx, transport, nil)
	if err != nil {
		return mcpHealthErrorStatus(mcpID, classifyMCPError(err)), nil
	}
	defer session.Close()

	if _, err := session.ListTools(probeCtx, nil); err != nil {
		kind := classifyMCPError(err)
		if kind == "connection_failed" {
			kind = "tool_list_failed"
		}
		return mcpHealthErrorStatus(mcpID, kind), nil
	}
	return MCPHealthStatus{MCPID: mcpID, State: MCPHealthHealthy}, nil
}

func resolveMap(refs map[string]string) map[string]string {
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

func hasUnresolvedMapRef(refs map[string]string) bool {
	for _, value := range refs {
		if hasUnresolvedEnvRef(value) {
			return true
		}
	}
	return false
}

func (s *Service) mcpTransport(ctx context.Context, environmentID string, activation *MCPActivation) (mcp.Transport, error) {
	switch activation.Transport {
	case catalog.MCPTransportStreamableHTTP:
		transport := &mcp.StreamableClientTransport{Endpoint: activation.Endpoint, MaxRetries: -1, DisableStandaloneSSE: true}
		if len(activation.Headers) != 0 {
			transport.HTTPClient = &http.Client{Transport: mcpHeaderRoundTripper{base: http.DefaultTransport, headers: activation.Headers}}
		}
		return transport, nil
	case catalog.MCPTransportStdio:
		rt, _, err := s.Runtime(environmentID)
		if err != nil {
			return nil, err
		}
		cmd, err := rt.Command(ctx, activation.Executable, activation.Args, activation.Env)
		if err != nil {
			return nil, err
		}
		return &mcp.CommandTransport{Command: cmd}, nil
	default:
		return nil, fmt.Errorf("unsupported mcp transport %q", activation.Transport)
	}
}

// MCPCommand returns the allowlisted command used for one resolved stdio
// activation. It is intentionally runtime-only and contains resolved values.
func (s *Service) MCPCommand(ctx context.Context, environmentID string, activation *MCPActivation) (*exec.Cmd, error) {
	if activation == nil || activation.Transport != catalog.MCPTransportStdio {
		return nil, fmt.Errorf("stdio MCP activation is required")
	}
	rt, _, err := s.Runtime(environmentID)
	if err != nil {
		return nil, err
	}
	return rt.Command(ctx, activation.Executable, activation.Args, activation.Env)
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
// The canonical implementation remains classifyMCPError in this package.
func ClassifyMCPError(err error) string {
	return classifyMCPError(err)
}

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
