---
gsd_state_version: 1.0
milestone: V2
current_phase: 14
current_phase_name: Evidence-first Investigation Toolkit
status: phase-14-01-complete-next-slice-pending
stopped_at: Phase 14 Plan 14-01 endpoint evidence resolution completed and locally verified at 0e40394; additional Phase 14 slices pending explicit plan selection
last_updated: "2026-09-09T15:38:00Z"
last_activity: 2026-09-09
last_activity_desc: Merged Phase 13 back to master, cleaned obsolete rebaseline ADM workspace/env metadata, and completed Phase 14 Plan 14-01 endpoint evidence resolver with Gateway investigate_endpoint; full test/vet/diff gates passed
state_head: 0e40394
progress:
  total_phases: 16
  completed_phases: 12
  total_plans: 19
  completed_plans: 19
  percent: 79
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, and `.planning/rebaseline/2026-09-08-core-boundary.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 14 — Evidence-first Investigation Toolkit
Status: Phase 14 Plan 14-01 endpoint evidence resolution is complete and locally verified; additional Phase 14 slices are pending explicit plan selection
Base master: `703593f`
Active development branch: `master`
Phase 11 importer implementation: `ea0d85af99f4591a431ba22dd8f1df036c76cac6`
Phase 12 source-aware Skill implementation: `4074d3e8e5349ee717bf63ad027391d14579cecc`
Phase 12 Skill availability implementation: `8f46f6e019afc8722708e3a4478de06faa4b241a`
Phase 13 static capability report implementation: `b0eb0c74e6f43844b3a754feaa2d34c93bb82da6`
Phase 13 Gateway-owner capability report implementation: `46c86f3`
Phase 14 endpoint evidence resolver implementation: `0e40394`

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
- Next slice ⏳: choose only with an explicit plan; candidates include symbol/reference/write tracing, data lineage, symbol-scoped Git history/diff, debug-SQL reverse mapping, test-data metric explanation and semantic consistency.

### Next Core priorities

1. Continue Phase 14 only by adding one explicit evidence-first slice at a time.
2. Desktop Core parity only after validated Core capability semantics are stable.
3. Temporary resource lifecycle/retention design remains deferred and must be explicitly planned before implementation.

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
- Phase 13 Gateway-owner enrichment reads existing owner-local observations only; it does not reconnect, Ping, refresh inventory, call MCP tools, run verifiers, start processes, or acquire writer leases.
- Phase 14 investigation helpers must return concrete evidence, confidence and uncertainties. 14-01 endpoint investigation is static literal/dynamic-candidate search only and does not execute code or call endpoints.
- Active ADM V2 development now continues directly on `master`; commit at each clear node and do not push unless explicitly requested.
- Temporary resource lifecycle/retention is captured for later design in `.planning/rebaseline/2026-09-09-temporary-resource-lifecycle.md`; CLI/UI-created resources are durable by default, while MCP/Gateway-created temporary resource semantics need explicit ownership, TTL, attachment and safe cleanup rules before implementation.

## Deferred

- automatic Memory context composition;
- temporary resource lifecycle/retention policy for MCP/Gateway-created Env/MCP/Skill/resources, including ownership, TTL, last-used tracking, attachment semantics, dry-run cleanup and managed-worktree safety;
- additional evidence-first investigation slices after 14-01 (symbol/reference/write tracing, data lineage, symbol-scoped Git history, test-data metric explanation, debug-SQL reverse mapping, semantic consistency);
- Desktop feature expansion beyond blockers until validated Core parity phase;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: Phase 14 Plan 14-01 implementation `0e40394` is complete and locally verified on `master`. Phase 13 has been merged back to master with local merge commit `50e6b75`. Obsolete rebaseline ADM workspace/env metadata has been removed; the physical Git worktree directory was not deleted. Temporary resource lifecycle/retention remains a future design item, not implemented.

Next action: choose the next Phase 14 evidence-first slice with an explicit plan, or pause for review. Continue directly on `master`; do not resume or merge `feat/gsd-phase-executor`; do not push unless explicitly requested.
