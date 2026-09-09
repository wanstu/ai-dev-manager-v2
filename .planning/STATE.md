---
gsd_state_version: 1.0
milestone: V2
current_phase: 11
current_phase_name: MCP Runtime Completion
status: in-progress
stopped_at: Phase 11-01 typed MCP runtime and Phase 11-02 monitor/reconnect are checkpointed on local master; Phase 11-03 JSON/JSONC import adapters remain next
last_updated: "2026-09-09T09:45:00Z"
last_activity: 2026-09-09
last_activity_desc: Closed the Phase 11-02 MCP monitor/reconnect checkpoint with HTTP and stdio background reconnect acceptance
state_head: b47c370
progress:
  total_phases: 16
  completed_phases: 10
  total_plans: 18
  completed_plans: 14
  percent: 78
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, `.planning/rebaseline/2026-09-08-core-boundary.md`, and `.planning/phases/11-mcp-runtime-completion/11-WORKING-STATE.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 11 — MCP Runtime Completion
Status: Plan 11-01 and Plan 11-02 are implemented and checkpointed on local `master`; Plan 11-03 import adapters remain not implemented.
Local master checkpoint before this state update: `b47c370` (`test(11): add stdio MCP reconnect acceptance`)
Relative to `origin/master`: local master is ahead by 1 commit at the previous checkpoint.
Working tree after this state update contains only 11-02 closeout documentation and strengthened inspect assertions intended for the next small commit.

The prior state saying Phase 11-01 was active in an uncommitted working tree and 11-02 had only observation/refresh scaffolding is superseded by the local master checkpoints:

- `3b39c40 feat(11): add typed MCP runtime and monitor`
- `b47c370 test(11): add stdio MCP reconnect acceptance`

The prior `feat/gsd-phase-executor` branch remains abandoned and must not merge. Its planning/provenance/state-advance implementation is not ADM product scope.

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

Status: complete on local master.

- BOUNDARY-01 ✅: Planner/Executor/Reviewer workflow surface/domain is removed from ADM Core while generic single-command Runs remain.
- BOUNDARY-02 ✅: no GSD `.planning` interpretation/state-advance API was merged or introduced.
- BOUNDARY-03 ✅: files/exec, verifier, process, MCP, Skill, managed worktree, ordinary Environment and generic Run regression gates passed during Phase 10 closeout.

### Phase 11 — MCP Runtime Completion

#### Plan 11-01 — Typed MCP Configuration + HTTP/Stdio Runtime

Status: implemented and documented for review.

Checkpoint includes typed MCP model/catalog separation, Streamable HTTP + stdio activation/connect paths, stdio allowlist authority, CLI and management typed `mcp add` surfaces, Gateway typed configuration, persisted secret/reference boundary, Ping-based explicit health probing and real transport acceptance.

Evidence:

- `.planning/phases/11-mcp-runtime-completion/11-01-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-01-VERIFICATION.md`

Residual review item: Desktop/UI parity should still be reviewed beyond adapter type changes.

#### Plan 11-02 — Health Monitor, Automatic Reconnect, Inventory + Diagnostics

Status: implemented and documented for review.

Checkpoint includes owner-local runtime observation, `environment_mcp_inspect`, `environment_mcp_refresh`, background `runtimeOwner.Monitor`, configurable Ping health checks, fixed-interval auto reconnect, real Streamable HTTP and stdio background reconnect acceptance, desired-vs-observed restart boundary and no automatic tool-call replay.

Evidence:

- `.planning/phases/11-mcp-runtime-completion/11-02-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-02-VERIFICATION.md`

#### Plan 11-03 — JSON / JSONC Import Adapters

Status: not started in code.

No `mcp_import_preview` implementation was found before the 11-01/11-02 checkpoints. The committed 11-03 plan remains the next Phase 11 implementation target.

### Next Core priorities after Phase 11

1. SKILL-COMP-01..04 — complete Skill refresh/source/support availability and diagnostics without a Skill execution engine.
2. CAP-01..02 — one authoritative Environment capability availability/diagnostic view.

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

## Current Risks / Closure Gates

1. **Import:** Phase 11-03 remains unimplemented and is the main remaining Phase 11 scope.
2. **Desktop/UI parity:** review typed MCP configuration display/edit expectations before product-facing release.
3. **Monolithic verification:** plugin-observed split package tests are green, but a local terminal `go test ./... -count=1` is still useful before push/release because prior monolithic plugin calls timed out at the tool transport layer.
4. **Multi-session continuity:** all sessions must use this STATE plus `11-WORKING-STATE.md` rather than older chat-only assumptions.

## Deferred

- automatic Memory context composition;
- evidence-first investigation helpers until MCP/Skill/capability Core is complete;
- Desktop feature expansion beyond blockers until validated Core parity phase;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: Phase 11-01 and 11-02 have local master checkpoints and review evidence. Phase 11-03 JSON/JSONC import adapters remain the next implementation target. Do not restart 11-01/11-02 from old assumptions, and do not merge/push automatically.
