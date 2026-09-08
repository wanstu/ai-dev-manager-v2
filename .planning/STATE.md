---
gsd_state_version: 1.0
milestone: V2
current_phase: 10
current_phase_name: Orchestration Boundary Cleanup
status: planned
stopped_at: Core-boundary rebaseline drafted; Phase 10 cleanup plan is next; Phase 10 GSD branch abandoned
last_updated: "2026-09-08T08:24:00Z"
last_activity: 2026-09-08
last_activity_desc: Rebased ADM around MCP, Skill, Environment capability control and local Runtime; orchestration removed from product scope
state_head: eca6cc6
progress:
  total_phases: 16
  completed_phases: 8
  total_plans: 10
  completed_plans: 10
  percent: 50
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, and `.planning/rebaseline/2026-09-08-core-boundary.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 10 — Orchestration Boundary Cleanup
Status: planned after 2026-09-08 core-boundary rebaseline
Base master: `eca6cc6`
Active rebaseline branch: `rebaseline/core-mcp-skill-runtime`

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

### Phase 10 — Boundary cleanup

- BOUNDARY-01: remove Planner/Executor/Reviewer workflow surface/domain from ADM Core while retaining generic single-command Runs.
- BOUNDARY-02: do not merge or introduce GSD `.planning` interpretation/state-advance APIs.
- BOUNDARY-03: cleanup must not regress files, exec, verifier, process, MCP, Skill, Git/worktree, Environment or generic Run behavior.

### Next Core priorities

1. MCP-COMP-01..05 — complete MCP configuration, supported transports/auth, inventory refresh, lifecycle/reconnect and actionable diagnostics.
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

## Deferred

- automatic Memory context composition;
- evidence-first investigation helpers until MCP/Skill/capability Core is complete (endpoint resolution, symbol/reference/write tracing, data lineage, symbol-scoped Git history, test-data metric explanation, debug-SQL reverse mapping, semantic consistency);
- Desktop feature expansion beyond blockers until validated Core parity phase;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: rebaseline documents and Product Contract are being corrected on `rebaseline/core-mcp-skill-runtime`.

Next action: validate rebaseline consistency, commit the planning/product-boundary change, then write Phase 10 Orchestration Boundary Cleanup CONTEXT/PLAN. Do not resume or merge `feat/gsd-phase-executor`.
