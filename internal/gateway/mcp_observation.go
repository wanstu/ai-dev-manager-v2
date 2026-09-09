package gateway

import (
	"context"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const maxObservedMCPTools = 200

type MCPRuntimeObservation struct {
	EnvironmentID       string             `json:"environment_id"`
	MCPID               string             `json:"mcp_id"`
	State               app.MCPHealthState `json:"state"`
	Transport           string             `json:"transport,omitempty"`
	FailureStage        string             `json:"failure_stage,omitempty"`
	ErrorKind           string             `json:"error_kind,omitempty"`
	Message             string             `json:"message,omitempty"`
	LastCheckAt         time.Time          `json:"last_check_at,omitempty"`
	LastHealthyAt       time.Time          `json:"last_healthy_at,omitempty"`
	LastSuccessAt       time.Time          `json:"last_success_at,omitempty"`
	ConsecutiveFailures int                `json:"consecutive_failures"`
	NextReconnectAt     *time.Time         `json:"next_reconnect_at,omitempty"`
	ProbeInFlight       bool               `json:"probe_in_flight"`
	ReconnectInFlight   bool               `json:"reconnect_in_flight"`
	Inventory           []*mcp.Tool        `json:"inventory"`
	InventoryFetchedAt  time.Time          `json:"inventory_fetched_at,omitempty"`
}

type MCPRuntimeInspection struct {
	Definition  model.MCPDefinition   `json:"definition"`
	Observation MCPRuntimeObservation `json:"observation"`
}

func (o *runtimeOwner) Inspect(ctx context.Context, environmentID, mcpID string) (MCPRuntimeInspection, error) {
	definition, err := o.service.MCPs.Get(mcpID)
	if err != nil {
		return MCPRuntimeInspection{}, err
	}
	activation, status, err := o.service.ResolveMCPActivation(environmentID, mcpID)
	if err != nil {
		return MCPRuntimeInspection{}, err
	}
	transport := definition.Transport
	if activation != nil {
		transport = activation.Transport
	}
	key := runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID}
	observation, ok := o.observation(key)
	if !ok {
		observation = MCPRuntimeObservation{
			EnvironmentID: environmentID,
			MCPID:         mcpID,
			State:         status.State,
			Transport:     transport,
			ErrorKind:     status.ErrorKind,
			Message:       status.Message,
			Inventory:     []*mcp.Tool{},
		}
	}
	return MCPRuntimeInspection{Definition: definition, Observation: observation}, nil
}

func (o *runtimeOwner) Refresh(ctx context.Context, environmentID, mcpID string) (MCPRuntimeInspection, error) {
	o.drop(runtimeOwnerKey{environmentID: environmentID, mcpID: mcpID})
	if _, status, err := o.ensureHealthySession(ctx, environmentID, mcpID); err != nil {
		return MCPRuntimeInspection{}, err
	} else if status.State != app.MCPHealthHealthy {
		return o.Inspect(ctx, environmentID, mcpID)
	}
	return o.Inspect(ctx, environmentID, mcpID)
}

func mcpProbeContext(ctx context.Context, seconds int64) (context.Context, context.CancelFunc) {
	if seconds <= 0 {
		seconds = 10
	}
	return context.WithTimeout(ctx, time.Duration(seconds)*time.Second)
}

func pingMCPErrorKind(err error) string {
	kind := runtimeMCPErrorKind(err)
	if kind == "timeout" {
		return "ping_timeout"
	}
	if kind == "connection_failed" {
		return "ping_protocol_failure"
	}
	return kind
}

func discoverOwnedMCPTools(ctx context.Context, session ownedMCPSession) ([]*mcp.Tool, error) {
	params := &mcp.ListToolsParams{}
	var tools []*mcp.Tool
	for {
		page, err := session.ListTools(ctx, params)
		if err != nil {
			return nil, err
		}
		tools = append(tools, page.Tools...)
		if page.NextCursor == "" || len(tools) >= maxObservedMCPTools {
			if len(tools) > maxObservedMCPTools {
				tools = tools[:maxObservedMCPTools]
			}
			return tools, nil
		}
		params.Cursor = page.NextCursor
	}
}

func (o *runtimeOwner) observation(key runtimeOwnerKey) (MCPRuntimeObservation, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	value, ok := o.observations[key]
	if !ok {
		return MCPRuntimeObservation{}, false
	}
	value.Inventory = append([]*mcp.Tool(nil), value.Inventory...)
	return value, true
}

func (o *runtimeOwner) observe(key runtimeOwnerKey, status app.MCPHealthStatus, transport, stage string, inventory []*mcp.Tool, successful bool) {
	now := time.Now().UTC()
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return
	}
	current := o.observations[key]
	current.EnvironmentID = key.environmentID
	current.MCPID = key.mcpID
	current.State = status.State
	if transport != "" {
		current.Transport = transport
	}
	current.FailureStage = stage
	current.ErrorKind = status.ErrorKind
	current.Message = status.Message
	current.LastCheckAt = now
	current.ProbeInFlight = false
	current.ReconnectInFlight = false
	current.NextReconnectAt = nil
	switch status.State {
	case app.MCPHealthHealthy:
		current.ConsecutiveFailures = 0
		current.FailureStage = ""
		current.LastHealthyAt = now
		if successful {
			current.LastSuccessAt = now
		}
	case app.MCPHealthError:
		current.ConsecutiveFailures++
		current.Inventory = []*mcp.Tool{}
		current.InventoryFetchedAt = time.Time{}
	default:
		current.ConsecutiveFailures = 0
	}
	if inventory != nil {
		if len(inventory) > maxObservedMCPTools {
			inventory = inventory[:maxObservedMCPTools]
		}
		current.Inventory = append([]*mcp.Tool(nil), inventory...)
		current.InventoryFetchedAt = now
	}
	if current.Inventory == nil {
		current.Inventory = []*mcp.Tool{}
	}
	o.observations[key] = current
}
