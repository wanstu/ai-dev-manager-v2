---
gsd_state_version: 1.0
milestone: V2
current_phase: 10
current_phase_name: Orchestration Boundary Cleanup
status: locally-verified-review-pending
stopped_at: Phase 10 orchestration cleanup implementation and verification passed on rebaseline/core-mcp-skill-runtime; integration review pending; Phase 11 implementation not started
last_updated: "2026-09-08T11:26:00Z"
last_activity: 2026-09-08
last_activity_desc: Removed ADM-owned Planner/Executor/Reviewer workflow semantics while retaining generic single-command Run lifecycle; full Core regression gates passed
state_head: 3620248
progress:
  total_phases: 16
  completed_phases: 9
  total_plans: 18
  completed_plans: 11
  percent: 56
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, and `.planning/rebaseline/2026-09-08-core-boundary.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 10 — Orchestration Boundary Cleanup
Status: Phase 10 implementation and local verification passed; integration review pending; Phases 11-13 remain planned only and not started
Base master: `eca6cc6`
Active rebaseline branch: `rebaseline/core-mcp-skill-runtime`
Phase 10 implementation: `36202489c13111cba2621b6d5a373707b7fde3f2`

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

### Next Core priorities

1. MCP-COMP-01..07 — complete MCP configuration, supported transports/auth, inventory refresh, configurable health/automatic reconnect, actionable diagnostics and single/batch external JSON/JSONC import.
2. SKILL-COMP-01..04 — complete Skill refresh/source/support availability and diagnostics without a Skill execution engine.
3. CAP-01..02 — one authoritative Environment capability availability/diagnostic view.

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
- Phase 13 capability inspection is side-effect-free by default and aggregates per-capability facts instead of probing/executing optional tools.

## Deferred

- automatic Memory context composition;
- evidence-first investigation helpers until MCP/Skill/capability Core is complete (endpoint resolution, symbol/reference/write tracing, data lineage, symbol-scoped Git history, test-data metric explanation, debug-SQL reverse mapping, semantic consistency);
- Desktop feature expansion beyond blockers until validated Core parity phase;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: Phase 10 implementation commit `3620248` is complete and local verification passed; closeout evidence/planning truth is being finalized before integration review. Phase 11 MCP, Phase 12 Skill and Phase 13 capability-diagnostics remain planned only.

Next action: complete Phase 10 integration review against the rebaseline base/master boundary. Do not start Phase 11 implementation until that review/integration boundary is resolved. Do not resume or merge `feat/gsd-phase-executor`.
