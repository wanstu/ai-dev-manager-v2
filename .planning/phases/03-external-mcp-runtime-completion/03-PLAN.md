---
phase: 03
slug: external-mcp-runtime-completion
plans:
  - id: "03-01"
    title: "Tracer — Core Health Lifecycle End-to-End"
    wave: 1
    depends_on: []
    autonomous: true
    files_modified:
      - internal/model/types.go
      - internal/catalog/service.go
      - internal/catalog/service_test.go
      - internal/app/mcp_health.go
      - internal/app/mcp_health_test.go
      - internal/gateway/server.go
      - internal/gateway/mcp_acceptance_test.go
    requirements: [MCP-RUN-01, MCP-RUN-02, MCP-RUN-03, MCP-RUN-04, MCP-RUN-05]
  - id: "03-02"
    title: "Surface Enrichment — CLI, Management, Desktop"
    wave: 2
    depends_on: ["03-01"]
    autonomous: true
    files_modified:
      - internal/management/service.go
      - internal/desktop/adapter.go
      - cmd/ai-dev-manager/main.go
      - cmd/ai-dev-manager/main_test.go
    requirements: [MCP-RUN-01, MCP-RUN-02, MCP-RUN-03, MCP-RUN-04, MCP-RUN-05]
---

# Phase 03 Plan — External MCP Runtime Completion

## Goal

Turn current HTTP proxying into a real Environment-scoped external MCP runtime lifecycle.

## Tracer Strategy

Lead with one production-quality end-to-end tracer slice (Plan 03-01) that proves the full health lifecycle works:

> Add MCP → enable for Environment → probe healthy → list tools → call tool → disable → probe disabled → break endpoint → probe error

Then expand with management surface enrichment (Plan 03-02).

---

## Threat Model

**ASVS Level 1** (local-only development gateway, no network exposure).

| Threat | Severity | Mitigation |
|--------|----------|------------|
| Secret endpoint values leaked in tool output / error messages / inspect | HIGH | Secrets resolved only at connect boundary via `os.ExpandEnv`; resolved values never stored, returned, or included in error text; `connectExternalMCP` errors sanitize the endpoint before returning. |
| Health probe used as SSRF scanner against internal network | LOW | Gateway is loopback-only (`allowedGatewayHost`); ASVS L1 accepts local-only risk. |
| Catalog entry without Transport defaults to unsupported transport | MEDIUM | `Transport` defaults to `"streamable-http"` when empty; backward-compatible with existing entries. No migration needed. |
| `connectExternalMCP` error message includes resolved endpoint URL containing secrets | HIGH | `connectExternalMCP` is refactored to never include the endpoint URL in error messages; callers receive `error_kind` + `mcp_id` only. |

---

## Specification

### What this phase delivers

1. `CatalogEntry` gains `Transport` and `HeaderRefs` fields for transport-specific MCP definitions (MCP-RUN-01).
2. `ProbeMCPHealth` returns 4-state health (configured/disabled/healthy/error) from on-demand connect+list-tools evidence (MCP-RUN-03).
3. Gateway tools (`environment_mcp_tools`, `environment_mcp_call`) return structured `MCPError` responses with `error_kind` on failure (MCP-RUN-03/04).
4. `environment_mcp_status` Gateway tool explicitly probes health without listing tools (MCP-RUN-03).
5. Secret values resolved only at connect boundary via `os.ExpandEnv` on `HeaderRefs`; never exposed in output (MCP-RUN-05).
6. CLI, Management, and Desktop surfaces expose health status (MCP-RUN-03).

### Non-goals

- Persistent connection pooling or session ownership (Phase 5, LIFE-01..03).
- Periodic/background health monitoring (Phase 5+).
- stdio/command-based MCP transport (ADM-CORE-012: first transport is HTTP).
- Desktop MCP health UI beyond basic inspection surface.
- Migration or compatibility burden for existing catalog entries.

---

## Plan 03-01: Tracer — Core Health Lifecycle End-to-End (Wave 1)

### Task 03-01-01: Model + Catalog Enrichment

**wave:** 1
**depends_on:** []

<read_first>
- internal/model/types.go (current CatalogEntry struct at line 46-55)
- internal/catalog/service.go (current AddMCP at line 35-45, add at line 93-118)
- internal/catalog/service_test.go (existing test patterns)
</read_first>

<acceptance_criteria>
- `CatalogEntry` has `Transport string` field defaulting to `"streamable-http"` when empty.
- `CatalogEntry` has `HeaderRefs map[string]string` field.
- `AddMCPConfig` method accepts `MCPConfig` struct with Transport, HeaderRefs, Endpoint, DefaultInclude.
- Existing `AddMCP` wraps `AddMCPConfig` with default Transport `"streamable-http"` and nil HeaderRefs.
- Backward-compatible: existing catalog entries without Transport are treated as `"streamable-http"`.
- Catalog test `TestAddMCPConfigSetsTransportAndHeaders` passes.
- Catalog test `TestAddMCPDefaultsTransportToStreamableHTTP` passes.
- `go test ./internal/catalog/...` passes.
</acceptance_criteria>

<verify>
Run `go test ./internal/catalog/...` — all tests pass including new AddMCPConfig tests.
</verify>

<action>
Add Transport and HeaderRefs fields to CatalogEntry in internal/model/types.go. Add MCPConfig struct and AddMCPConfig method to internal/catalog/service.go. Refactor existing AddMCP to wrap AddMCPConfig with defaults. Add tests to internal/catalog/service_test.go for both the new method and backward-compat default behavior.
</action>

<artifacts>
- internal/model/types.go — Transport, HeaderRefs fields on CatalogEntry
- internal/catalog/service.go — MCPConfig struct, AddMCPConfig method, refactored AddMCP
- internal/catalog/service_test.go — TestAddMCPConfigSetsTransportAndHeaders, TestAddMCPDefaultsTransportToStreamableHTTP
</artifacts>

---

### Task 03-01-02: Health Probe Module

**wave:** 1
**depends_on:** ["03-01-01"]

<read_first>
- internal/model/types.go (CatalogEntry with Transport/HeaderRefs from Task 01)
- internal/catalog/service.go (List, Get methods for catalog lookup)
- internal/app/service.go (InspectEnvironment at line 98-130, resolveCatalogSelections at line 138-153)
- internal/gateway/server.go (connectExternalMCP at line 571-578, enabledMCPEndpoint at line 553-568)
</read_first>

<acceptance_criteria>
- NEW FILE `internal/app/mcp_health.go` exists with:
  - `MCPHealthState` type (string alias) with constants: `MCPHealthConfigured`, `MCPHealthDisabled`, `MCPHealthHealthy`, `MCPHealthError`.
  - `MCPHealthStatus` struct with `MCPID`, `State`, `ErrorKind`, `Message` fields.
  - `MCPError` struct with `MCPID`, `ErrorKind`, `Message` fields implementing `error`.
  - `ProbeMCPHealth(ctx, environmentID, mcpID)` method on `*Service`.
  - `resolveEndpoint(endpoint, headerRefs)` function using `os.ExpandEnv`.
  - `resolveHeaders(headerRefs)` function using `os.ExpandEnv`.
- `ProbeMCPHealth` returns:
  - `"disabled"` when MCP is not enabled for the Environment.
  - `"configured"` when catalog entry exists but endpoint is empty or not yet probed.
  - `"healthy"` when connect + ListTools succeeds within bounded timeout.
  - `"error"` with structured `ErrorKind` when connect or ListTools fails.
- `resolveEndpoint` expands `${VAR}` references in endpoint URL via `os.ExpandEnv`.
- `resolveHeaders` expands `${VAR}` references in header values via `os.ExpandEnv`.
- `MCPError` implements `error` interface and formats as `mcp_id=<id> error_kind=<kind> message=<msg>`.
- NEW FILE `internal/app/mcp_health_test.go` exists with tests for all 4 states + secret resolution + non-exposure.
- `go test ./internal/app/...` passes.
</acceptance_criteria>

<verify>
Run `go test ./internal/app/...` — all tests pass including health probe tests for all 4 states.
</verify>

<action>
Create internal/app/mcp_health.go with MCPHealthState constants, MCPHealthStatus struct, MCPError type, ProbeMCPHealth method, resolveEndpoint function, and resolveHeaders function. ProbeMCPHealth checks enablement via InspectEnvironment, checks catalog entry existence, resolves secrets at connect boundary, attempts connect+ListTools with bounded timeout (10s). Create internal/app/mcp_health_test.go with unit tests covering all 4 health states, env-var expansion in resolveEndpoint, and secret non-exposure in MCPError messages.
</action>

<artifacts>
- internal/app/mcp_health.go (NEW) — MCPHealthState, MCPHealthStatus, MCPError, ProbeMCPHealth, resolveEndpoint, resolveHeaders
- internal/app/mcp_health_test.go (NEW) — TestProbeMCPHealthConfigured, TestProbeMCPHealthDisabled, TestProbeMCPHealthHealthy, TestProbeMCPHealthError, TestResolveEndpointExpandsEnvVars, TestSecretsNotExposedInErrorMessages
</artifacts>

---

### Task 03-01-03: Gateway Enrichment + Acceptance Test

**wave:** 1
**depends_on:** ["03-01-01", "03-01-02"]

<read_first>
- internal/app/mcp_health.go (ProbeMCPHealth, MCPHealthStatus, MCPError from Task 02)
- internal/gateway/server.go (environment_mcp_tools at line 358-383, environment_mcp_call at line 384-400, connectExternalMCP at line 571-578, enabledMCPEndpoint at line 553-568)
- internal/gateway/server_test.go (TestGatewayUsesOnlyEnabledConfiguredExternalMCPAndSkillRuntime at line 572-750, connectInMemory at line 1028, toolText at line 1062)
</read_first>

<acceptance_criteria>
- `environment_mcp_status` tool added to Gateway: accepts `EnvironmentMCPStatusInput`, calls `ProbeMCPHealth`, returns structured `MCPHealthStatus`. No writer required.
- `environment_mcp_tools` handler returns structured `MCPError` with `error_kind` when connect fails (instead of raw `fmt.Errorf`).
- `environment_mcp_call` handler returns structured `MCPError` with `error_kind` when connect fails.
- `connectExternalMCP` no longer includes the endpoint URL in error messages (prevents secret leakage).
- `connectExternalMCP` accepts resolved headers parameter and sets them via `HTTPClient` round-tripper on `StreamableClientTransport`.
- NEW FILE `internal/gateway/mcp_acceptance_test.go` exists with full lifecycle test.
- Acceptance test proves: add MCP → enable → probe healthy → list tools → call tool → disable → probe disabled → break endpoint → probe error.
- Acceptance test proves: structured errors on broken MCP; unrelated Gateway tools survive MCP error.
- `go test ./internal/gateway/...` passes.
- `go test ./...` passes.
</acceptance_criteria>

<verify>
Run `go test ./internal/gateway/...` — all tests pass including new acceptance test. Run `go test ./...` — full suite green.
</verify>

<action>
Add `EnvironmentMCPStatusInput` struct and `environment_mcp_status` tool to internal/gateway/server.go. Enrich `environment_mcp_tools` and `environment_mcp_call` error paths to return `MCPError`-formatted errors. Refactor `connectExternalMCP` to accept resolved headers and avoid including endpoint URL in error messages. Create internal/gateway/mcp_acceptance_test.go with TestMCPHealthLifecycleEndToEnd that exercises the full health lifecycle through the in-memory Gateway.
</action>

<artifacts>
- internal/gateway/server.go — EnvironmentMCPStatusInput, environment_mcp_status tool, enriched error handlers, refactored connectExternalMCP
- internal/gateway/mcp_acceptance_test.go (NEW) — TestMCPHealthLifecycleEndToEnd
</artifacts>

---

## Plan 03-02: Surface Enrichment — CLI, Management, Desktop (Wave 2)

### Task 03-02-01: Management + Desktop + CLI Passthrough

**wave:** 2
**depends_on:** ["03-01"]

<read_first>
- internal/app/mcp_health.go (ProbeMCPHealth, MCPHealthStatus from Plan 03-01)
- internal/management/service.go (current MCP methods at line 100-110)
- internal/desktop/adapter.go (current MCP methods at line 162-181)
- cmd/ai-dev-manager/main.go (runCatalog at line 505-612, runEnvironmentSelection at line 353-389)
- cmd/ai-dev-manager/main_test.go (existing CLI test patterns)
</read_first>

<acceptance_criteria>
- `management/service.go` gains `MCPHealth(ctx, environmentID, mcpID) (app.MCPHealthStatus, error)` method.
- `desktop/adapter.go` gains `ProbeMCPHealth(environmentID, mcpID) (app.MCPHealthStatus, error)` method.
- CLI `mcp status --id MCP_ID --environment-id ENV_ID` command added, prints JSON health status.
- CLI `mcp status -h` prints usage help.
- `go test ./internal/management/...` passes (or existing tests continue passing).
- `go test ./cmd/ai-dev-manager/...` passes.
</acceptance_criteria>

<verify>
Run `go test ./...` — full suite green after management/desktop/CLI changes.
</verify>

<action>
Add MCPHealth method to internal/management/service.go that delegates to app.ProbeMCPHealth. Add ProbeMCPHealth method to internal/desktop/adapter.go that delegates to management.MCPHealth. Add "mcp status" case to runCatalog in cmd/ai-dev-manager/main.go with --id and --environment-id flags. Add CLI test for mcp status --help in cmd/ai-dev-manager/main_test.go.
</action>

<artifacts>
- internal/management/service.go — MCPHealth method
- internal/desktop/adapter.go — ProbeMCPHealth method
- cmd/ai-dev-manager/main.go — mcp status CLI command
- cmd/ai-dev-manager/main_test.go — test for mcp status help output
</artifacts>

---

### Task 03-02-02: Final Validation Suite

**wave:** 2
**depends_on:** ["03-02-01"]

<read_first>
- All files modified in Plan 03-01 and Plan 03-02
- internal/gateway/mcp_acceptance_test.go (full lifecycle test from Plan 03-01)
- internal/app/mcp_health_test.go (health probe unit tests from Plan 03-01)
- internal/gateway/server_test.go (existing tests that must remain green)
</read_first>

<acceptance_criteria>
- `go test ./...` passes with zero failures.
- `go vet ./...` passes.
- `git diff --check` is clean.
- All MCP-RUN-01 through MCP-RUN-05 acceptance tests pass.
- Structured errors return `error_kind` and `mcp_id` without secret leakage.
- Health states correctly distinguish configured/disabled/healthy/error.
- At least one real external MCP call succeeds end-to-end through the Gateway.
- No resolved secret values appear in any test output, error message, or inspection result.
- Empty endpoint returns "configured" state (not "error").
- Empty tool name in mcp_call returns structured error with `error_kind: "missing_tool_name"`.
- Disabled Environment MCP access returns structured error with `error_kind: "not_enabled"`.
- Broken MCP endpoint returns structured error with `error_kind: "connection_refused"` or `"timeout"`.
</acceptance_criteria>

<verify>
Run `go test ./...` — full suite green. Run `go vet ./...` — no warnings. Run `git diff --check` — clean.
</verify>

<action>
Run full validation suite: go test ./..., go vet ./..., git diff --check. Verify all acceptance criteria are met. If any test fails, fix the underlying issue. This task is primarily a validation gate — no new code unless a pre-existing test breaks.
</action>

<artifacts>
- No new files; this task validates existing artifacts from Plans 03-01 and 03-02.
</artifacts>

---

## Verification Criteria

1. `go test ./...` passes.
2. `go vet ./...` passes.
3. `git diff --check` is clean.
4. All `MCP-RUN-01` through `MCP-RUN-05` acceptance tests pass.
5. Real Streamable HTTP acceptance proves: add MCP → enable → probe healthy → list tools → call tool → disable → probe disabled → break endpoint → probe error.
6. Structured errors carry `error_kind` and `mcp_id` without secret leakage.
7. PROJECT/ROADMAP/STATE/verification/UAT/SUMMARY are updated only after implementation evidence exists.

---

## must_haves

### truths

- "D-01: CatalogEntry has Transport field defaulting to streamable-http when empty; no stdio transport added in Phase 3"
- "D-01: CatalogEntry has HeaderRefs map for env-var-referenced auth headers"
- "D-01: AddMCPConfig persists Transport and HeaderRefs; AddMCP wraps it with defaults"
- "D-04: ProbeMCPHealth returns 'configured' for unprobed MCP with empty endpoint"
- "D-04: ProbeMCPHealth returns 'disabled' for MCP not enabled in Environment"
- "D-04: ProbeMCPHealth returns 'healthy' when connect+ListTools succeeds"
- "D-04: ProbeMCPHealth returns 'error' with structured error_kind when connect or ListTools fails"
- "D-03: resolveEndpoint expands ${VAR} references via os.ExpandEnv at connect boundary only"
- "D-03: resolveHeaders expands ${VAR} references via os.ExpandEnv at connect boundary only"
- "D-07: MCPError implements error interface with mcp_id, error_kind, message fields"
- "D-06: environment_mcp_status tool returns structured MCPHealthStatus without requiring writer"
- "D-06: environment_mcp_tools returns MCPError with error_kind on connection failure"
- "D-06: environment_mcp_call returns MCPError with error_kind on connection failure"
- "D-03/D-07: connectExternalMCP never includes resolved endpoint URL in error messages"
- "D-03: connectExternalMCP sets resolved headers via HTTPClient on StreamableClientTransport"
- "D-05: CLI mcp status --id ID --environment-id ENV_ID returns JSON health status"
- "D-05: management.MCPHealth delegates to app.ProbeMCPHealth"
- "D-05: desktop.ProbeMCPHealth delegates to management.MCPHealth"
- "D-02: Concurrent MCP health probes create independent per-call sessions with no shared mutable state"
- "D-04: ProbeMCPHealth returns exactly one of four states per invocation; no partial health state exists"
- "D-07: Empty tool name in environment_mcp_call returns structured error with error_kind missing_tool_name"
- "D-04: Empty or missing endpoint returns 'configured' state with diagnostic message"
- "D-04: Health is queried per individual MCP ID, not as an ordered collection"
- "Tool listing order from external MCPs is preserved as returned by the remote server"
- { statement: "D-02: os.ExpandEnv is stateless per call; concurrent probes resolve independently without shared mutable state", verification: backstop }

### prohibitions

- { requirement_id: "MCP-RUN-05", category: "privacy", status: "unresolved", verification: null, statement: "Secret endpoint values (resolved from env vars) must never appear in Gateway tool output, CLI output, error messages, inspection results, or log output", check_kind: null, check_target: null, check_rule: null, check_violation_fixture: null, check_clean_fixture: null }
- { requirement_id: "MCP-RUN-05", category: "privacy", status: "unresolved", verification: null, statement: "The resolved endpoint URL containing actual secret values must never be persisted to catalog state, returned from AddMCP/AddMCPConfig, or included in Snapshot output", check_kind: null, check_target: null, check_rule: null, check_violation_fixture: null, check_clean_fixture: null }
- { requirement_id: "MCP-RUN-03", category: "safety", status: "dismissed", verification: null, reason: "Health probe SSRF risk is mitigated by loopback-only gateway (allowedGatewayHost); ASVS L1 accepts local-only risk", statement: "Health probe must not be weaponized as SSRF scanner against internal network" }

### flagged_assumptions

- MCP-RUN-02 / unclassified: "Environment selection gating mechanism (MCP-RUN-02) is already proven working by existing test TestGatewayUsesOnlyEnabledConfiguredExternalMCPAndSkillRuntime. Phase 3 enriches it with health semantics but the gating mechanism itself requires no changes. This assumption should be verified during Phase 4 dogfood."
- MCP-RUN-05 prohibition (privacy): "External MCP server error responses may echo back the endpoint URL or auth headers. ADM wraps external errors but cannot control what the remote server includes in its error message body. This is an accepted limitation for Phase 3 — a future Phase may add explicit error-body sanitization."

---

## Artifacts this phase produces

| Artifact | Type | Location |
|----------|------|----------|
| `MCPHealthState` type | Go type | internal/app/mcp_health.go |
| `MCPHealthStatus` struct | Go type | internal/app/mcp_health.go |
| `MCPError` struct | Go type | internal/app/mcp_health.go |
| `MCPConfig` struct | Go type | internal/catalog/service.go |
| `ProbeMCPHealth` method | Go method | internal/app/mcp_health.go |
| `resolveEndpoint` function | Go function | internal/app/mcp_health.go |
| `resolveHeaders` function | Go function | internal/app/mcp_health.go |
| `AddMCPConfig` method | Go method | internal/catalog/service.go |
| `MCPHealth` method | Go method | internal/management/service.go |
| `ProbeMCPHealth` adapter method | Go method | internal/desktop/adapter.go |
| `environment_mcp_status` tool | Gateway tool | internal/gateway/server.go |
| `mcp status` CLI command | CLI command | cmd/ai-dev-manager/main.go |
| Transport field on CatalogEntry | Model field | internal/model/types.go |
| HeaderRefs field on CatalogEntry | Model field | internal/model/types.go |
| Health probe unit tests | Tests | internal/app/mcp_health_test.go |
| Catalog AddMCPConfig tests | Tests | internal/catalog/service_test.go |
| Gateway MCP acceptance test | Tests | internal/gateway/mcp_acceptance_test.go |
| CLI mcp status test | Tests | cmd/ai-dev-manager/main_test.go |
