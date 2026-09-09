---
gsd_state_version: 1.0
milestone: V2
current_phase: 14
current_phase_name: Evidence-first Investigation Toolkit
status: phase-14-02-planned-gitnexus-provider-next
stopped_at: Phase 14 Plan 14-02 optional code intelligence provider and GitNexus integration boundary planned; implementation not started
last_updated: "2026-09-09T15:52:00Z"
last_activity: 2026-09-09
last_activity_desc: Planned GitNexus as optional code intelligence provider, prioritized remaining Phase 14 investigation slices, promoted temporary Env/MCP/Skill/provider lifecycle to Phase 15, and moved Desktop Core parity after lifecycle semantics; no feature code was changed
state_head: ad974ad
progress:
  total_phases: 17
  completed_phases: 12
  total_plans: 21
  completed_plans: 19
  percent: 76
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, `.planning/rebaseline/2026-09-08-core-boundary.md`, and `.planning/phases/14-evidence-first-investigation-toolkit/14-PRIORITY-MAP.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 14 — Evidence-first Investigation Toolkit
Status: 14-01 endpoint evidence resolution is complete; 14-02 optional code intelligence provider + GitNexus integration boundary is planned but not implemented
Base master: `703593f`
Active development branch: `master`
Phase 11 importer implementation: `ea0d85af99f4591a431ba22dd8f1df036c76cac6`
Phase 12 source-aware Skill implementation: `4074d3e8e5349ee717bf63ad027391d14579cecc`
Phase 12 Skill availability implementation: `8f46f6e019afc8722708e3a4478de06faa4b241a`
Phase 13 static capability report implementation: `b0eb0c74e6f43844b3a754feaa2d34c93bb82da6`
Phase 13 Gateway-owner capability report implementation: `46c86f3`
Phase 14 endpoint evidence resolver implementation: `0e40394`
Phase 14 endpoint evidence resolver closeout: `ad974ad`

The prior `feat/gsd-phase-executor` branch is abandoned and must not merge. Its planning/provenance/state-advance implementation is not ADM product scope.

## Retained Core History

### Phases 1-4 — Development foundations

- real `SKILL.md` discovery/read from explicit roots with support-root containment and Environment gating;
- structured optional verifier runtime;
- real external MCP Streamable HTTP tool discovery/call with Environment gating, four-state health and activation-boundary secret resolution;
- external Agent dogfood through ADM-only local development capabilities.

### Phases 5-8 — Persistent local Runtime

- persistent Gateway owner and external MCP session restart reconciliation;
- long-running dev process lifecycle, bounded logs and owned-port facts;
- optional managed Git worktree isolation;
- generic single-command asynchronous `run_` start/list/status/cancel lifecycle.

These capabilities remain ADM Core.

## Superseded Work

### Phase 9 — Planner / Executor / Reviewer

Phase 9 was technically implemented, reviewed and integrated to local master, but the product rebaseline determined that its workflow orchestration semantics belong to the Agent/GSD/orchestrator layer.

The Git/evidence history remains for auditability. The Agent-facing workflow surface and workflow-specific domain code are cleanup scope under BOUNDARY-01. Future ADM capabilities must not depend on FLOW-01.

### Abandoned Phase 10 GSD executor branch

`feat/gsd-phase-executor` is not to be merged. `gsd_phase_inspect`, `gsd_phase_start`, `gsd_phase_advance`, `.planning` provenance interpretation and STATE advancement are explicitly outside ADM Core.

## Active Requirements

### Phase 10 — Boundary cleanup (complete)

- BOUNDARY-01 ✅: Planner/Executor/Reviewer workflow surface/domain is removed from ADM Core while generic single-command Runs remain.
- BOUNDARY-02 ✅: no GSD `.planning` interpretation/state-advance API was merged or introduced.
- BOUNDARY-03 ✅: files/exec, verifier, process, MCP, Skill, managed worktree, ordinary Environment and generic Run regression gates passed.

### Phase 11 — MCP runtime completion (complete)

- 11-01 ✅: typed desired MCP configuration, activation-only references and real Streamable HTTP/stdio owner-bound transport runtime.
- 11-02 ✅: owner-local health/recovery observation, configurable probe/reconnect policy, inventory/inspect/refresh, stable-ID update invalidation and safe no-replay semantics; full test/vet/race/diff gates passed at `4c4f4dc`.
- 11-03 ✅: external JSON/JSONC preview/apply import adapters, batch atomicity, conflict policy, credential-reference conversion, Gateway/management/CLI surfaces and runtime-owner definition-fingerprint/generation reconciliation; full `go test -count=1 ./...` plus vet/race/diff gates passed at `ea0d85a`.

### Phase 12 — Skill runtime completion (complete)

- 12-01 ✅: persisted Skill sources, source/artifact stable identity, atomic per-source refresh, safe failed-refresh preservation, unresolved selection preservation and Gateway/management/CLI source surfaces; full test/vet/race/diff gates passed at `4074d3e`.
- 12-02 ✅: Environment-specific Skill availability, bounded artifact/support inventory, structured Skill read diagnostics and broken-Skill isolation; full test/vet/race/diff gates passed at `8f46f6e`.

### Phase 13 — Environment capability diagnostics (complete)

- 13-01 ✅: shared `CapabilityReport`/`CapabilityFact`/`CapabilityEvidence` model and resilient side-effect-free application-level Environment report with static file/exec/verifier/Git/isolation/MCP/Skill/process/run facts; full test/vet/race/diff gates passed at `b0eb0c7`.
- 13-02 ✅: Gateway-owner observation enrichment and canonical Agent-facing `environment_capability_report`, plus CLI/management wrappers; full test/vet/race/diff gates passed at `46c86f3`.

### Phase 14 — Evidence-first Investigation Toolkit (in progress)

- 14-01 ✅: endpoint evidence resolution for URL/path plus optional HTTP method; returns bounded static route evidence, confidence and uncertainties; full test/vet/diff gates passed at `0e40394`.
- 14-02 📝: optional code intelligence provider + GitNexus integration boundary is planned in `14-02-PLAN.md`; implementation has not started.
- Next implementation target ⏳: provider-neutral investigation layer with GitNexus as optional external provider, then symbol/reference/write tracing.

### Phase 15 — Temporary Resource Lifecycle (planned)

- 15-01 📝: temporary Env/MCP/Skill/provider resource lifecycle and cleanup preview/execute semantics are planned for after the initial Phase 14 provider/symbol work.
- CLI/UI-created resources stay durable by default; Gateway/MCP-created resources may become temporary only with explicit metadata and conservative cleanup rules.

### Phase 16 — Desktop Core Parity (planned)

- Desktop should expose the same Core state and lifecycle semantics after temporary resource lifecycle is stable.
- No Desktop-only state or cleanup semantics.

### Phase 17 — Distribution Only If Needed (conditional)

- Installer/tray/autostart/updater/signing/notifications remain conditional on demonstrated daily-use need.

### Next Core priorities

1. Implement Phase 14 Plan 14-02: provider-neutral investigation layer and optional GitNexus integration boundary.
2. Implement 14-03 symbol/reference/write tracing using GitNexus/provider evidence when available and static fallback otherwise.
3. Implement Phase 15 temporary resource lifecycle before Desktop Core parity.
4. Move Desktop Core parity to Phase 16 after lifecycle semantics stabilize.

## Product Decisions

- MCP and Skill are core product capabilities and must stay visible in the main roadmap.
- A foundation vertical slice is not the same as completion. Phase 3 does not mean MCP is finished; Phase 1 does not mean Skill is finished.
- `run_` is retained only as generic asynchronous Runtime lifecycle; it must not grow task semantics.
- Worktree remains an optional isolation primitive. ADM does not orchestrate parallel Agents or choose integration policy.
- Verifier returns structured evidence; the external Agent/orchestrator decides what that evidence means for its task/phase.
- No LLM/GSD/OpenCode subagent quota is used unless the user explicitly reverses the existing instruction.
- No automatic push at review boundaries.
- During active development, continue directly on `master`; commit at each clear node.
- Phase 11 transport scope is Streamable HTTP + stdio; legacy HTTP+SSE is not added. Stdio MCP executable launch must still obey ADM's executable allowlist.
- Phase 11 health-check interval/probe timeout are configurable per MCP. Automatic reconnect defaults to off; when enabled it retries at one fixed configured interval with no exponential/adaptive backoff. Health/recovery observation is owner-local, and recovery never replays failed tool calls.
- Phase 11 supports preview + atomic single/batch JSON/JSONC import adapters for OpenCode, WorkBuddy/CodeBuddy, Codex plugin MCP JSON, Claude Code and supported MCPHub shapes. Name conflicts default to error. Import writes only global MCP definitions and never modifies existing Environment selections. Literal credentials are converted into generated secret/environment-reference requirements and only references are persisted.
- Phase 11 does not persist interactive OAuth tokens until an approved secure credential lifecycle exists; supported HTTP auth is none or explicit secret-backed headers.
- Phase 12 uses persisted Skill sources, source/artifact-based stable identity and explicit atomic refresh; ADM does not interpret Skill instructions.
- Phase 12 availability is Environment-specific; a broken Skill is local to that Skill and does not poison unrelated Skill, MCP, file, verifier or Runtime operations.
- Phase 13 capability inspection is side-effect-free by default and aggregates per-capability facts instead of probing/executing optional tools.
- Phase 13 Gateway-owner enrichment reads existing owner-local observations only; it does not reconnect, Ping, refresh inventory, call MCP tools, run verifiers, start processes, or acquire writer leases.
- Phase 14 investigation helpers must return concrete evidence, confidence and uncertainties. GitNexus/code-graph support is optional provider integration, not a hard dependency or ADM-owned graph engine.
- Phase 15 temporary resource lifecycle/retention is now planned as an explicit Core phase before Desktop. CLI/UI-created resources are durable by default; Gateway/MCP-created temporary resource semantics need explicit ownership, TTL, attachment and safe cleanup rules before implementation.
- Desktop Core parity moves after temporary lifecycle so Desktop exposes stable Core semantics instead of inventing its own cleanup model.

## Deferred

- automatic Memory context composition;
- Phase 14 follow-up slices after provider integration: symbol/reference/write tracing, impact/blast-radius evidence, data lineage, symbol-scoped Git history, test-data metric explanation, debug-SQL reverse mapping and semantic consistency;
- Desktop feature expansion beyond blockers until validated Core parity and temporary lifecycle semantics are stable;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: planning-only update on `master`: 14-02 optional code intelligence provider + GitNexus integration boundary is planned; Phase 15 temporary resource lifecycle is planned; Desktop Core parity is moved to Phase 16; Distribution is moved to Phase 17. No feature implementation was started after 14-01.

Next action: implement 14-02 on `master`, then commit at the next clear implementation node. Do not resume or merge `feat/gsd-phase-executor`; do not push unless explicitly requested.
