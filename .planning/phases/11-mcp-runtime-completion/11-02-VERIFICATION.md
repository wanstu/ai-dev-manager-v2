# Plan 11-02 Verification — MCP Monitor/Reconnect Closeout Evidence

Date: 2026-09-09

This file records the closeout verification for Plan 11-02: MCP health monitor, automatic reconnect, inventory and diagnostics.

## Focused monitor and reconnect tests

```text
go test ./internal/gateway -run TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance|TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance|TestRuntimeOwnerBackground|TestRuntimeOwnerRestart|TestRuntimeOwnerInspect -count=1 -v
```

Result: passed.

Observed passing tests include:

- `TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance`
- `TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance`
- `TestRuntimeOwnerBackgroundMonitorSkipsPingWhenHealthCheckDisabled`
- `TestRuntimeOwnerBackgroundMonitorDoesNotReconnectWhenDisabled`
- `TestRuntimeOwnerBackgroundMonitorReconnectsOnFixedInterval`
- `TestRuntimeOwnerBackgroundReconnectDoesNotResurrectDisabledMCP`
- `TestRuntimeOwnerRestartPersistsHealthPolicyButNotObservation`
- `TestRuntimeOwnerInspectAndRefreshObservation`

These cover:

- `health_check_enabled=false` prevents background Ping/reconnect;
- `health_check_enabled=true` runs background Ping after the configured interval;
- `auto_reconnect=false` prevents background reconnect while explicit safe status/list/refresh can still reconnect;
- `auto_reconnect=true` reconnects after the fixed configured interval;
- no exponential/adaptive backoff is used in Phase 11;
- failure observation records stage/kind/failure count and clears stale inventory;
- pending reconnect is cancelled by disabled desired state;
- desired `MCPHealthPolicy` persists across owner restart;
- owner-local observation/inventory/timestamps do not persist across owner restart.

## Gateway API-level inspect evidence

`TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance` verifies `environment_mcp_inspect` through a real Gateway tool call, not only direct owner internals.

It asserts structured JSON fields for:

- `definition.transport`;
- `definition.health_policy.health_check_enabled`;
- `definition.health_policy.check_interval_seconds`;
- `definition.health_policy.auto_reconnect`;
- `definition.health_policy.reconnect_interval_seconds`;
- `observation.state`;
- `observation.failure_stage`;
- `observation.consecutive_failures`;
- `observation.next_reconnect_at`;
- `observation.probe_in_flight`;
- `observation.reconnect_in_flight`;
- `observation.inventory`;
- `observation.inventory_fetched_at`.

The test proves healthy inventory is visible, broken inventory is cleared, and recovered inventory is restored.

## Real Streamable HTTP background reconnect acceptance

```text
go test ./internal/gateway -run TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance -count=1 -v
```

Result: passed.

This test uses a real SDK Streamable HTTP MCP server and proves:

- owner starts healthy and discovers real tool inventory;
- background Ping detects a temporarily broken upstream;
- stale inventory/session are cleared;
- reconnect does not run before the configured fixed interval;
- after upstream recovery, background reconnect restores healthy state and inventory;
- a later explicit tool call succeeds after recovery.

## Real stdio background reconnect acceptance

```text
go test ./internal/gateway -run TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance -count=1 -v
```

Result: passed.

This test uses a real stdio MCP helper child process and proves:

- stdio MCP starts only after executable allowlist authority is satisfied;
- initial status/list/call work with the real child process;
- a test-only `stdio_stop` tool terminates the helper;
- background Ping detects the stopped child/session;
- stale inventory/session are cleared;
- reconnect does not run before the fixed interval;
- background reconnect starts a second helper and restores health/inventory;
- reconnect does not replay `stdio_echo`;
- a later explicit `stdio_echo` call runs exactly once after recovery.

## Race coverage

```text
go test -race ./internal/gateway -run TestRuntimeOwnerBackground|TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance|TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance -count=1
```

Result:

```text
ok  ai-dev-manager-v2/internal/gateway  14.074s
```

## Gateway regression split

```text
go test ./internal/gateway -run TestMCP|TestStdio|TestRuntimeOwner|TestGateway|TestHTTPGateway -count=1
```

Result: passed.

The gateway package can be long-running in a single plugin call, so Gateway verification is recorded through focused and split test groups.

## Related package regression

```text
go test ./internal/app ./internal/catalog ./internal/management -count=1
```

Result: passed.

Previously recorded split package runs after the 11-01 checkpoint passed:

```text
go test ./cmd/ai-dev-manager ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/environment ./internal/isolation ./internal/runtime ./internal/skill ./internal/verifier -count=1

go test ./internal/identity ./internal/memory ./internal/model ./internal/store ./internal/workspace -count=1
```

## Vet and diff hygiene

```text
go vet ./...
```

Result: passed.

```text
git diff --check
```

Result: passed. Only existing LF-to-CRLF working-copy warnings were emitted; no whitespace errors were reported.

## Monolithic test note

A monolithic `go test ./... -count=1` has timed out through the plugin transport in prior attempts. The package coverage above is the current recorded closeout evidence. A local terminal monolithic run may still be useful before push/release, but the plugin-observed package split is green.
