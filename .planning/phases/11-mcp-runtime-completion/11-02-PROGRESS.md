# Plan 11-02 Progress — Owner Monitor/Reconnection Core

Date: 2026-09-09

This is an interim progress note, not a Plan 11-02 closeout. Plan 11-02 remains active until the remaining restart, diagnostics and race/fake-clock evidence is complete.

## Scope completed so far

Implemented the owner-local background monitor/reconnect loop under the persistent Gateway runtime owner and added real background reconnect acceptance for both supported transports.

Code touched in this Plan 11-02 slice:

- `internal/gateway/runtime_owner.go`
- `internal/gateway/runtime_owner_test.go`
- `internal/gateway/server.go`
- `internal/gateway/mcp_acceptance_test.go`
- `.planning/phases/11-mcp-runtime-completion/11-WORKING-STATE.md`
- `.planning/phases/11-mcp-runtime-completion/11-02-PROGRESS.md`
- `.planning/phases/11-mcp-runtime-completion/11-02-VERIFICATION.md`

## Runtime behavior now covered

- `RunHTTP` and `RunStdio` start `owner.Monitor(ownerCtx)` instead of a one-shot `owner.Reconcile(ownerCtx)`.
- `Monitor(ctx)` performs startup reconciliation and then repeatedly runs `MonitorOnce(ctx)`.
- `MonitorOnce(ctx)` evaluates enabled Environment/MCP selections against their desired `MCPHealthPolicy`.
- `health_check_enabled=false` prevents background Ping, leaving explicit safe status/list/refresh operations as the way to reconnect or refresh.
- `health_check_enabled=true` allows background Ping checks for an already-owned live session once `check_interval_seconds` has elapsed.
- Ping failure drops the unusable session, marks the observation unhealthy, records failure stage/kind, clears stale inventory and optionally schedules fixed-interval reconnect.
- `auto_reconnect=false` prevents background reconnect; explicit safe status/list/refresh operations may still reconnect.
- `auto_reconnect=true` schedules fixed-interval reconnect attempts using `reconnect_interval_seconds`.
- Background reconnect does not call arbitrary MCP tools and does not replay a failed tool call.
- Disable/drop closes runtime state and prevents pending background reconnect from resurrecting a disabled MCP.
- Restart inspection preserves desired health policy while owner-local observation/inventory/timestamps start empty for the new owner.
- A real Streamable HTTP MCP server can transition healthy -> broken -> background-reconnected healthy with inventory restored and no attached Agent client required during recovery.
- A real stdio MCP helper process can transition healthy -> stopped child/session -> background-reconnected healthy with a newly started helper and no automatic business tool replay.

## New/expanded tests

Added or expanded owner-level tests covering:

- `health_check_enabled=false` does not background Ping or reconnect;
- background monitor does not reconnect when `auto_reconnect=false`;
- background monitor reconnects after the configured fixed interval when `auto_reconnect=true`;
- observation fields after failure: `failure_stage`, `error_kind`, `consecutive_failures`, cleared inventory, no in-flight flags and scheduled `next_reconnect_at`;
- background reconnect does not resurrect a disabled MCP;
- background reconnect does not replay a failed tool call;
- restart preserves desired `MCPHealthPolicy` but not owner-local observation state.

Added real acceptance coverage:

- `TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance` uses a real SDK Streamable HTTP MCP server and proves background Ping failure detection plus fixed-interval reconnect after the upstream recovers.
- `TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance` uses the real test binary as a stdio MCP helper, stops the helper through an explicit test tool, proves the monitor marks the session unhealthy, waits for the fixed reconnect interval, starts a second helper in the background and verifies the reconnect did not replay `stdio_echo`.

Existing runtime owner tests continue to cover:

- session reuse and close;
- restart rebuilding from persisted desired state;
- dead session cannot remain healthy;
- disable evicts a live session;
- shared owner across independent Agent sessions;
- inspect/refresh observation behavior;
- real HTTP Gateway restart reconciliation.

## Important non-goals still remaining

This slice still does not complete all of 11-02. Remaining work includes:

- deeper race/fake-clock coverage for monitor/reconnect/drop paths;
- broader Gateway API-level `environment_mcp_inspect` assertions, especially around serialized in-flight fields and reconnect timing;
- formal 11-02 summary and closeout verification after the remaining review gates.

## Safety note

Background recovery is limited to connection/session recovery and bounded inventory discovery. It never invokes arbitrary MCP tools and does not replay a tool call that previously failed.
