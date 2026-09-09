# Plan 11-02 Summary — MCP Health Monitor, Automatic Reconnect, Inventory + Diagnostics

Date: 2026-09-09

## Status

Plan 11-02 is implemented to a reviewable closeout checkpoint on local `master` after the Phase 11 typed MCP runtime work.

This plan adds an owner-local MCP runtime lifecycle: background health checks, fixed-interval reconnect, bounded inventory observation, explicit inspect/refresh diagnostics and safe recovery semantics without automatic tool-call replay.

## Implemented behavior

### Owner-local observation

`MCPRuntimeObservation` records live, owner-local facts for each `(environment_id, mcp_id)`:

- state and transport;
- failure stage and structured error kind;
- last check / last healthy / last success timestamps;
- consecutive failure count;
- next reconnect time;
- probe/reconnect in-flight flags;
- bounded public tool inventory and fetch timestamp.

The observation remains runtime-owner state only. It is not persisted as desired ADM state.

### Background monitor

`runtimeOwner.Monitor(ctx)` now runs for `RunHTTP` and `RunStdio`. It performs startup reconciliation and then calls `MonitorOnce(ctx)` on a bounded interval.

`MonitorOnce(ctx)` evaluates currently enabled Environment/MCP selections from desired state and prunes disabled/removed sessions.

### Health checks

For enabled MCPs with a live owner session:

- `health_policy.health_check_enabled=false` prevents background Ping;
- `health_policy.health_check_enabled=true` allows protocol Ping once `check_interval_seconds` has elapsed;
- each Ping is bounded by `probe_timeout_seconds`;
- Ping failure drops the session, marks the observation unhealthy, records failure stage/kind, clears inventory and optionally schedules reconnect.

### Automatic reconnect

For unhealthy/no-session MCPs:

- `auto_reconnect=false` prevents background reconnect;
- explicit safe status/list/refresh can still reconnect;
- `auto_reconnect=true` schedules reconnect at the fixed configured `reconnect_interval_seconds`;
- no exponential backoff, adaptive retry or hidden jitter is introduced;
- disable/drop prevents pending recovery from resurrecting old capability;
- successful reconnect restores healthy observation and refreshes inventory.

### Safe replay boundary

Background recovery only reconnects and refreshes public inventory. It never invokes arbitrary MCP tools and never replays a failed tool call. Later tool calls must be explicit caller actions.

### API diagnostics

`environment_mcp_inspect` exposes sanitized desired config plus current owner observation. Acceptance now verifies API-level fields including health policy, state, failure stage, failure count, next reconnect time, in-flight flags, inventory and inventory timestamps.

`environment_mcp_refresh` remains the explicit safe reconnect/refresh surface.

## Real transport acceptance

Plan 11-02 now has real background reconnect acceptance for both supported transports:

- Streamable HTTP: real SDK HTTP MCP server transitions healthy → broken → background reconnected healthy.
- Stdio: real helper MCP child process exits, monitor detects failure, fixed-interval reconnect starts a second helper and restores inventory.

Both tests verify that background recovery does not replay business tool calls.

## Verification snapshot

See `11-02-VERIFICATION.md` for commands and evidence. The closeout checkpoint includes split package test coverage, targeted race coverage, `go vet ./...` and `git diff --check`.

## Remaining Phase 11 work

Plan 11-02 closeout does not complete all of Phase 11. Plan 11-03 JSON/JSONC import adapters remain not implemented.

Desktop/UI parity should still be reviewed as part of broader product review, but the runtime owner, Gateway API and CLI/management foundations are now in place for MCP runtime completion.
