# Phase 11 Working State — 2026-09-09

This file records the live Phase 11 repository state so new sessions do not mistake older planning snapshots or chat-only context for the current code checkpoint.

## Git snapshot

- Branch: `master`
- Last implementation checkpoint before this closeout doc update: `b47c370 test(11): add stdio MCP reconnect acceptance`.
- Local master was `ahead 1` relative to `origin/master` at that checkpoint.
- Current closeout worktree contains only strengthened Gateway API inspect assertions and planning evidence updates for Plan 11-02.

Prior large uncommitted Phase 11 work was checkpointed in:

- `3b39c40 feat(11): add typed MCP runtime and monitor`
- `b47c370 test(11): add stdio MCP reconnect acceptance`

## Plan 11-01 — Typed MCP Configuration + HTTP/Stdio Runtime

Status: **implemented and documented for review**.

Evidence:

- `.planning/phases/11-mcp-runtime-completion/11-01-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-01-VERIFICATION.md`

Implemented behavior includes:

- dedicated `MCPDefinition` and `MCPHealthPolicy` model fields;
- dedicated MCP catalog service with `streamable-http` and `stdio` validation;
- transport-specific activation through `ResolveMCPActivation`;
- Streamable HTTP and stdio connection paths through the official MCP SDK;
- stdio execution routed through existing ADM Runtime command preparation/allowlist authority;
- Gateway-side transport-specific MCP configuration support;
- CLI `mcp add` typed HTTP/stdio configuration surface, including header/env refs and health policy flags;
- management `MCPAddConfig` typed configuration surface;
- persisted secret/reference boundary: HTTP `HeaderRefs` must use environment references, and credential-bearing stdio `EnvRefs` keys must use environment references instead of literal values;
- Gateway `mcp_add` negative coverage proving literal credential-bearing refs are rejected and not persisted;
- stdio acceptance coverage, including executable allowlist enforcement and owner cleanup;
- Ping-based explicit health probing;
- compatibility fix for no-owner Gateway external MCP tools/call path so it does not hard-gate consumption on Ping when direct list/call succeeds;
- runtime-owner context-boundary fix so owner-owned upstream MCP operations use owner lifecycle context rather than the outer Gateway request context.

Residual review item:

- Desktop/UI parity has not been fully reviewed beyond adapter type changes.

## Plan 11-02 — Health Monitor, Automatic Reconnect, Inventory + Diagnostics

Status: **implemented and documented for review**.

Evidence:

- `.planning/phases/11-mcp-runtime-completion/11-02-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-02-VERIFICATION.md`

Implemented behavior includes:

- owner-local `MCPRuntimeObservation`;
- health/error/inventory timestamps and bounded tool inventory fields;
- explicit `environment_mcp_inspect` and `environment_mcp_refresh` paths;
- Ping-based probing helpers;
- bounded inventory discovery and refresh plumbing;
- `runtimeOwner.Monitor(ctx)` long-running background loop;
- `runtimeOwner.MonitorOnce(ctx)` deterministic single-iteration monitor entrypoint for tests;
- Gateway `RunHTTP` / `RunStdio` start `owner.Monitor(ownerCtx)` rather than one-shot reconciliation;
- background health check for enabled MCPs with `health_policy.health_check_enabled=true` and an existing live session;
- `health_check_enabled=false` prevents background Ping and reconnect;
- fixed-interval reconnect scheduling for unhealthy/no-session MCPs with `health_policy.auto_reconnect=true`;
- `auto_reconnect=false` behavior: background monitor does not reconnect, while explicit safe `Status`/`ListTools`/`Refresh` paths can still reconnect;
- disable/drop cleanup prevents pending background reconnect from resurrecting disabled MCPs;
- background reconnect never invokes arbitrary MCP tools and does not replay failed tool calls;
- desired `MCPHealthPolicy` persists across owner restart while owner-local observation/inventory/timestamps do not;
- real Streamable HTTP healthy → broken → background-reconnected healthy acceptance;
- real stdio child healthy → stopped → background-reconnected healthy acceptance;
- Gateway API-level structured `environment_mcp_inspect` assertions for health policy, state, failure stage/count, reconnect timing, in-flight flags and inventory timestamps;
- targeted `-race` coverage for monitor/reconnect paths.

## Plan 11-03 — JSON / JSONC Import Adapters

Status: **not implemented**.

Repository search before the 11-01/11-02 checkpoints found no `mcp_import_preview` implementation. The committed 11-03 plan remains the next implementation target.

## Verification snapshot

Current closeout evidence includes:

```text
go test ./internal/gateway -run TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance|TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance|TestRuntimeOwnerBackground|TestRuntimeOwnerRestart|TestRuntimeOwnerInspect -count=1 -v

go test -race ./internal/gateway -run TestRuntimeOwnerBackground|TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance|TestRuntimeOwnerRealStdioBackgroundReconnectAcceptance -count=1

go test ./internal/gateway -run TestMCP|TestStdio|TestRuntimeOwner|TestGateway|TestHTTPGateway -count=1

go test ./internal/app ./internal/catalog ./internal/management -count=1

go vet ./...

git diff --check
```

All passed in the plugin-observed split runs. `git diff --check` emitted only LF-to-CRLF working-copy warnings and no whitespace errors.

A monolithic `go test ./... -count=1` has timed out through the plugin transport in prior attempts. A local terminal monolithic run remains useful before push/release if required, but current package split evidence is green.

## Immediate next action

Commit this closeout doc/API-inspect assertion update, then continue Phase 11 with Plan 11-03 JSON/JSONC import adapters.

Do not restart 11-01/11-02 from older assumptions. Do not automatically merge or push.
