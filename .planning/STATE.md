---
gsd_state_version: 1.0
milestone: V2
current_phase: 13
current_phase_name: Environment Capability Diagnostics
status: phase-13-01-complete-phase-13-02-pending
stopped_at: Phase 13 Plan 13-01 resilient CapabilityFact model and static Environment report completed and locally verified at b0eb0c7; Phase 13 Plan 13-02 Gateway-owner observation enrichment not started
last_updated: "2026-09-09T14:45:00Z"
last_activity: 2026-09-09
last_activity_desc: Completed shared CapabilityReport/CapabilityFact model, side-effect-free app-level Environment capability report, static file/exec/verifier/Git/isolation/MCP/Skill/process/run facts, optional failure isolation, and capability_report exposure through environment inspect; full repository test plus vet/race/diff gates passed
state_head: b0eb0c7
progress:
  total_phases: 16
  completed_phases: 11
  total_plans: 18
  completed_plans: 17
  percent: 72
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, and `.planning/rebaseline/2026-09-08-core-boundary.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 13 — Environment Capability Diagnostics
Status: Plan 13-01 is complete and locally verified; Plan 13-02 Gateway-owner observation enrichment remains pending and has not started
Base master: `703593f`
Active rebaseline branch: `rebaseline/core-mcp-skill-runtime`
Phase 11 importer implementation: `ea0d85af99f4591a431ba22dd8f1df036c76cac6`
Phase 12 source-aware Skill implementation: `4074d3e8e5349ee717bf63ad027391d14579cecc`
Phase 12 Skill availability implementation: `8f46f6e019afc8722708e3a4478de06faa4b241a`
Phase 13 static capability report implementation: `b0eb0c74e6f43844b3a754feaa2d34c93bb82da6`

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

### Phase 10 — Boundary cleanup (locally verified)

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

### Phase 13 — Environment capability diagnostics (in progress)

- 13-01 ✅: shared `CapabilityReport`/`CapabilityFact`/`CapabilityEvidence` model and resilient side-effect-free application-level Environment report with static file/exec/verifier/Git/isolation/MCP/Skill/process/run facts; full test/vet/race/diff gates passed at `b0eb0c7`.
- 13-02 ⏳: Gateway-owner observation enrichment and canonical Agent-facing capability report remain pending. Do not start unless explicitly authorized.

### Next Core priorities

1. CAP-02 finalization through 13-02 — enrich the same capability model with honest Gateway-owner observations.
2. Evidence-first investigation helpers only after MCP/Skill/capability Core is complete.
3. Desktop Core parity only after validated Core capability semantics are stable.

## Product Decisions

- MCP and Skill are core product capabilities and must stay visible in the main roadmap.
- A foundation vertical slice is not the same as completion. Phase 3 does not mean MCP is finished; Phase 1 does not mean Skill is finished.
- `run_` is retained only as generic asynchronous Runtime lifecycle; it must not grow task semantics.
- Worktree remains an optional isolation primitive. ADM does not orchestrate parallel Agents or choose integration policy.
- Verifier returns structured evidence; the external Agent/orchestrator decides what that evidence means for its task/phase.
- No LLM/GSD/OpenCode subagent quota is used unless the user explicitly reverses the existing instruction.
- No automatic merge/push at review boundaries.
- Phase 11 transport scope is Streamable HTTP + stdio; legacy HTTP+SSE is not added. Stdio MCP executable launch must still obey ADM's executable allowlist.
- Phase 11 health-check interval/probe timeout are configurable per MCP. Automatic reconnect defaults to off; when enabled it retries at one fixed configured interval with no exponential/adaptive backoff. Health/recovery observation is owner-local, and recovery never replays failed tool calls.
- Phase 11 supports preview + atomic single/batch JSON/JSONC import adapters for OpenCode, WorkBuddy/CodeBuddy, Codex plugin MCP JSON, Claude Code and supported MCPHub shapes. Name conflicts default to error. Import writes only global MCP definitions and never modifies existing Environment selections. Literal credentials are converted into generated secret/environment-reference requirements and only references are persisted.
- Phase 11 does not persist interactive OAuth tokens until an approved secure credential lifecycle exists; supported HTTP auth is none or explicit secret-backed headers.
- Phase 12 uses persisted Skill sources, source/artifact-based stable identity and explicit atomic refresh; ADM does not interpret Skill instructions.
- Phase 12 availability is Environment-specific; a broken Skill is local to that Skill and does not poison unrelated Skill, MCP, file, verifier or Runtime operations.
- Phase 13 capability inspection is side-effect-free by default and aggregates per-capability facts instead of probing/executing optional tools.
- Temporary resource lifecycle/retention is captured for later design in `.planning/rebaseline/2026-09-09-temporary-resource-lifecycle.md`; CLI/UI-created resources are durable by default, while MCP/Gateway-created temporary resource semantics need explicit ownership, TTL, attachment and safe cleanup rules before implementation.

## Deferred

- automatic Memory context composition;
- temporary resource lifecycle/retention policy for MCP/Gateway-created Env/MCP/Skill/resources, including ownership, TTL, last-used tracking, attachment semantics, dry-run cleanup and managed-worktree safety;
- evidence-first investigation helpers until MCP/Skill/capability Core is complete (endpoint resolution, symbol/reference/write tracing, data lineage, symbol-scoped Git history, test-data metric explanation, debug-SQL reverse mapping, semantic consistency);
- Desktop feature expansion beyond blockers until validated Core parity phase;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: Phase 13 Plan 13-01 implementation `b0eb0c7` is complete and locally verified on `rebaseline/core-mcp-skill-runtime`; Plan 13-02 is pending and has not started. Temporary resource lifecycle/retention has been recorded as a future design item, not implemented.

Next action: review/integrate the 13-01 closeout under the existing boundary, then begin Phase 13 Plan 13-02 only when authorized. Do not start Phase 14 early and do not resume or merge `feat/gsd-phase-executor`.
