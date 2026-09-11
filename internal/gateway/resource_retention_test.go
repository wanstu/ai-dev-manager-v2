package gateway

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"ai-dev-manager-v2/internal/app"
	"ai-dev-manager-v2/internal/catalog"
	"ai-dev-manager-v2/internal/model"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGatewayResourceRetentionCleanupDryRunUsesRuntimeOwnerObservations(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	expiredAt := time.Now().UTC().Add(-time.Hour)
	mcpEntry, err := service.MCPs.AddMCPConfigWithRetention("temporary-unused-mcp", catalog.MCPConfig{
		Transport: catalog.MCPTransportStreamableHTTP,
		AuthMode:  catalog.MCPAuthNone,
		Endpoint:  "http://127.0.0.1:65534/mcp",
	}, model.ResourceRetention{
		Persistence:    model.PersistenceTemporary,
		CreatorSurface: "gateway",
		OwnerID:        "owner-runtime",
		ExpiresAt:      &expiredAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	ctx := context.Background()
	session := connectInMemory(t, ctx, newServerForSurface(service, owner, serverSurfaceAdmin))
	defer session.Close()

	dryRun, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "resource_retention_cleanup",
		Arguments: map[string]any{},
	})
	if err != nil || dryRun.IsError {
		t.Fatalf("resource_retention_cleanup dry-run failed: err=%v result=%+v", err, dryRun)
	}
	text := toolText(t, dryRun)
	for _, required := range []string{"\"dry_run\":true", "\"would_remove\"", mcpEntry.ID, model.RetentionResourceMCP} {
		if !strings.Contains(strings.ReplaceAll(text, " ", ""), strings.ReplaceAll(required, " ", "")) {
			t.Fatalf("dry-run cleanup output missing %q:\n%s", required, text)
		}
	}
	if strings.Contains(text, app.RetentionRuntimeObservationBlocker) {
		t.Fatalf("Gateway-owner dry run should resolve runtime observation blocker: %s", text)
	}
}

func TestGatewayResourceRetentionCleanupExecuteRemovesRuntimeSafeMCP(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	expiredAt := time.Now().UTC().Add(-time.Hour)
	mcpEntry, err := service.MCPs.AddMCPConfigWithRetention("temporary-runtime-safe-mcp", catalog.MCPConfig{
		Transport: catalog.MCPTransportStreamableHTTP,
		AuthMode:  catalog.MCPAuthNone,
		Endpoint:  "http://127.0.0.1:65534/mcp",
	}, model.ResourceRetention{
		Persistence:    model.PersistenceTemporary,
		CreatorSurface: "gateway",
		OwnerID:        "owner-runtime",
		ExpiresAt:      &expiredAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	ctx := context.Background()
	session := connectInMemory(t, ctx, newServerForSurface(service, owner, serverSurfaceAdmin))
	defer session.Close()

	executed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "resource_retention_cleanup",
		Arguments: map[string]any{"execute": true},
	})
	if err != nil || executed.IsError {
		t.Fatalf("resource_retention_cleanup execute failed: err=%v result=%+v text=%s", err, executed, toolText(t, executed))
	}
	text := toolText(t, executed)
	for _, required := range []string{"\"dry_run\":false", "\"removed\"", mcpEntry.ID, model.RetentionResourceMCP} {
		if !strings.Contains(strings.ReplaceAll(text, " ", ""), strings.ReplaceAll(required, " ", "")) {
			t.Fatalf("execute cleanup output missing %q:\n%s", required, text)
		}
	}
	if _, err := service.MCPs.Get(mcpEntry.ID); err == nil {
		t.Fatalf("execute did not remove runtime-safe temporary MCP %s", mcpEntry.ID)
	}
}

func TestGatewayResourceRetentionMarkTemporaryAndPromoteTools(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	expiredAt := time.Now().UTC().Add(-time.Hour)
	mcpEntry, err := service.MCPs.AddMCP("retention-tool-mcp", "http://127.0.0.1:65534/mcp", false)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	session := connectInMemory(t, ctx, NewAdmin(service))
	defer session.Close()
	marked, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "resource_retention_mark_temporary",
		Arguments: map[string]any{
			"kind":       model.RetentionResourceMCP,
			"id":         mcpEntry.ID,
			"owner_id":   "owner-tool",
			"expires_at": expiredAt.Format(time.RFC3339Nano),
		},
	})
	if err != nil || marked.IsError {
		t.Fatalf("resource_retention_mark_temporary failed: err=%v result=%+v", err, marked)
	}
	if text := toolText(t, marked); !strings.Contains(text, model.RetentionUpdateActionMarkTemporary) || !strings.Contains(text, "owner-tool") {
		t.Fatalf("mark temporary output missing retention facts: %s", text)
	}
	storedMarked, err := service.MCPs.Get(mcpEntry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if storedMarked.Retention.Persistence != model.PersistenceTemporary || storedMarked.Retention.OwnerID != "owner-tool" {
		t.Fatalf("mark temporary did not persist: %+v", storedMarked.Retention)
	}

	promoted, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "resource_retention_promote",
		Arguments: map[string]any{"kind": model.RetentionResourceMCP, "id": mcpEntry.ID},
	})
	if err != nil || promoted.IsError {
		t.Fatalf("resource_retention_promote failed: err=%v result=%+v", err, promoted)
	}
	if text := toolText(t, promoted); !strings.Contains(text, model.RetentionUpdateActionPromoteDurable) || !strings.Contains(text, model.PersistenceDurable) {
		t.Fatalf("promote output missing durable facts: %s", text)
	}
	storedPromoted, err := service.MCPs.Get(mcpEntry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if storedPromoted.Retention.Persistence != model.PersistenceDurable || storedPromoted.Retention.OwnerID != "" || storedPromoted.Retention.ExpiresAt != nil {
		t.Fatalf("promote did not clear temporary fields: %+v", storedPromoted.Retention)
	}
}

func TestGatewayResourceRetentionCleanupExecuteBlocksActiveMCPRuntime(t *testing.T) {
	service := app.New(filepath.Join(t.TempDir(), "state.json"))
	expiredAt := time.Now().UTC().Add(-time.Hour)
	mcpEntry, err := service.MCPs.AddMCPConfigWithRetention("temporary-inflight-mcp", catalog.MCPConfig{
		Transport: catalog.MCPTransportStreamableHTTP,
		AuthMode:  catalog.MCPAuthNone,
		Endpoint:  "http://127.0.0.1:65534/mcp",
	}, model.ResourceRetention{
		Persistence:    model.PersistenceTemporary,
		CreatorSurface: "gateway",
		OwnerID:        "owner-runtime",
		ExpiresAt:      &expiredAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	owner := newRuntimeOwner(service)
	defer owner.Close()
	owner.mu.Lock()
	owner.inFlight[runtimeOwnerKey{environmentID: "env_active", mcpID: mcpEntry.ID}] = true
	owner.mu.Unlock()

	ctx := context.Background()
	session := connectInMemory(t, ctx, newServerForSurface(service, owner, serverSurfaceAdmin))
	defer session.Close()

	executed, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "resource_retention_cleanup",
		Arguments: map[string]any{"execute": true},
	})
	if err != nil || executed.IsError {
		t.Fatalf("resource_retention_cleanup execute with active runtime failed unexpectedly: err=%v result=%+v", err, executed)
	}
	text := toolText(t, executed)
	if strings.Contains(text, "\"removed\"") && strings.Contains(text, mcpEntry.ID) {
		t.Fatalf("active MCP runtime should not be removed:\n%s", text)
	}
	if !strings.Contains(text, "mcp_operation_in_flight") {
		t.Fatalf("active MCP runtime blocker missing:\n%s", text)
	}
	if _, err := service.MCPs.Get(mcpEntry.ID); err != nil {
		t.Fatalf("active runtime cleanup removed MCP: %v", err)
	}
}
