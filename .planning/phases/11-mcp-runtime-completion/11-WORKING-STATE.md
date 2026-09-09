# Phase 11 Working State — 2026-09-09

This file records the current uncommitted Phase 11 implementation state so new sessions do not mistake older planning snapshots for the live repository state.

## Git snapshot

- Branch: `master`
- HEAD: `703593f docs(10): record integration review`
- Local master is `21` commits ahead of `origin/master`.
- Working tree remains uncommitted Phase 11 implementation state.
- Current tracked diff shortstat observed during this update: 19 files changed, 1769 insertions, 394 deletions.
- New/untracked Phase 11 files include:
  - `.planning/phases/11-mcp-runtime-completion/11-01-SUMMARY.md`
  - `.planning/phases/11-mcp-runtime-completion/11-01-VERIFICATION.md`
  - `.planning/phases/11-mcp-runtime-completion/11-02-PROGRESS.md`
  - `.planning/phases/11-mcp-runtime-completion/11-02-VERIFICATION.md`
  - `.planning/phases/11-mcp-runtime-completion/11-WORKING-STATE.md`
  - `internal/catalog/mcp_service.go`
  - `internal/gateway/mcp_observation.go`

Phase 10 is already integrated into local `master`. Phase 11 planning is committed, while the Phase 11 implementation described below is still uncommitted working-tree state.

## Plan 11-01 — Typed MCP Configuration + HTTP/Stdio Runtime

Status: **implementation substantially complete in the working tree and documented for review**.

Observed implementation now includes:

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

Evidence files now exist:

- `.planning/phases/11-mcp-runtime-completion/11-01-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-01-VERIFICATION.md`

Known residual 11-01 review item:

- Desktop/UI parity has not been fully reviewed beyond the current adapter type changes.

## Plan 11-02 — Health Monitor, Automatic Reconnect, Inventory + Diagnostics

Status: **in progress; owner-local monitor/reconnect core loop and real HTTP background reconnect acceptance are implemented, but full stdio/race/closeout evidence is not complete**.

Observed implementation now includes:

- owner-local `MCPRuntimeObservation`;
- health/error/inventory timestamps and bounded tool inventory fields;
- explicit inspect/refresh paths;
- Ping-based probing helpers;
- inventory discovery and refresh plumbing;
- observation fields for `NextReconnectAt`, `ProbeInFlight`, and `ReconnectInFlight`;
- `runtimeOwner.Monitor(ctx)` long-running background loop;
- `runtimeOwner.MonitorOnce(ctx)` deterministic single-iteration monitor entrypoint for tests;
- Gateway `RunHTTP` / `RunStdio` now start `owner.Monitor(ownerCtx)` rather than one-shot reconciliation;
- `health_check_enabled=false` prevents background Ping/reconnect and leaves explicit safe status/list/refresh as recovery paths;
- background health check for enabled MCPs with `health_policy.health_check_enabled=true` and an existing live session;
- Ping failure drops the unusable session, records failure stage/error kind/failure count, clears stale inventory and optionally schedules reconnect;
- fixed-interval reconnect scheduling for unhealthy/no-session MCPs with `health_policy.auto_reconnect=true`;
- `auto_reconnect=false` behavior: background monitor does not reconnect, while explicit safe `Status`/`ListTools`/`Refresh` paths can still reconnect;
- disable/drop cleanup prevents pending background reconnect from resurrecting disabled MCPs;
- background reconnect never invokes arbitrary MCP tools and does not replay failed tool calls;
- restart proof that desired health policy persists while owner-local observation/inventory/timestamps do not;
- real Streamable HTTP background reconnect acceptance: healthy -> broken -> background-reconnected healthy with inventory restored.

Current 11-02 tests cover:

- background monitor skips Ping when `health_check_enabled=false`;
- background monitor does not reconnect when `auto_reconnect=false`;
- background monitor reconnects on configured fixed interval when `auto_reconnect=true`;
- failure observation details: failure stage, error kind, consecutive failures, cleared inventory, in-flight flags and next reconnect time;
- background reconnect does not resurrect a disabled MCP;
- background reconnect does not replay a failed tool call;
- restart persists `MCPHealthPolicy` but does not persist owner-local observation/inventory/timestamps;
- real HTTP background reconnect acceptance through a real SDK Streamable HTTP MCP server;
- existing Gateway restart/reconciliation, HTTP, stdio, verifier/process/run regressions remain green in split runs.

Evidence files now exist:

- `.planning/phases/11-mcp-runtime-completion/11-02-PROGRESS.md`
- `.planning/phases/11-mcp-runtime-completion/11-02-VERIFICATION.md`

Remaining 11-02 gaps:

- real stdio unhealthy-to-background-reconnect acceptance still needs stronger end-to-end evidence;
- Gateway API-level `environment_mcp_inspect` assertions should explicitly cover next reconnect timing, in-flight flags, failure stage, failure count and inventory timestamps;
- deeper race/fake-clock coverage is not complete;
- no formal 11-02 closeout summary exists yet.

Therefore 11-02 must still be described as active/in-progress, not complete.

## Plan 11-03 — JSON / JSONC Import Adapters

Status: **not implemented**.

Repository search found no `mcp_import_preview` implementation. The committed 11-03 plan remains planning-only at this snapshot.

## Verification snapshot

Current split-package verification completed successfully:

```text
go test ./internal/app ./internal/catalog ./internal/management -count=1
ok  ai-dev-manager-v2/internal/app         26.267s
ok  ai-dev-manager-v2/internal/catalog      1.535s
ok  ai-dev-manager-v2/internal/management   1.774s
```

Gateway verification was split to avoid one long plugin transport call:

```text
go test ./internal/gateway -run TestRuntimeOwnerBackground|TestRuntimeOwnerRestart|TestRuntimeOwnerInspect|TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance -count=1 -v
go test ./internal/gateway -run TestMCP|TestStdio|TestRuntimeOwner|TestGateway|TestHTTPGateway -count=1
go test ./internal/gateway -run TestGateway|TestHTTPGateway|TestRealHost|TestStopHTTP|TestSameADM|TestMCP|TestStdio|TestRuntimeOwner -count=1
go test ./internal/gateway -run TestVerifier|TestProcess|TestRun|TestLegacy|TestAllowed|TestOwner|TestWait|TestInspect -count=1
```

Additional packages completed successfully:

```text
go test ./cmd/ai-dev-manager ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/environment -count=1
go test ./internal/isolation ./internal/runtime ./internal/skill ./internal/verifier -count=1
go test ./internal/identity ./internal/memory ./internal/model ./internal/store ./internal/workspace -count=1
```

Static/diff checks:

```text
go vet ./...
git diff --check
```

Both passed. `git diff --check` emitted only existing LF-to-CRLF working-copy warnings and no whitespace errors.

A monolithic `go test ./... -count=1` was attempted through the plugin but timed out at the tool-call layer. The split-package runs above cover the listed packages and should be used as the current verification evidence unless a local terminal can run the monolithic command directly.

## Immediate next action

Do **not** start 11-03 and do not treat 11-02 as complete. Continue Plan 11-02 by expanding real stdio reconnect and Gateway API diagnostics evidence:

1. add real stdio unhealthy-to-background-reconnect acceptance;
2. add Gateway API-level `environment_mcp_inspect` assertions for `next_reconnect_at`, `probe_in_flight`, `reconnect_in_flight`, failure stage, failure count and inventory timestamps;
3. consider fake-clock/race coverage for monitor/reconnect/drop paths;
4. rerun split-package verification and `git diff --check`;
5. write formal 11-02 summary/verification only when those acceptance gaps are closed.

No automatic merge or push is implied by this working-state record.
