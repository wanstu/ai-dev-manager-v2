---
gsd_state_version: 1.0
milestone: V2
current_phase: 11
current_phase_name: MCP Runtime Completion
status: in-progress-uncommitted
stopped_at: Phase 10 is integrated on local master at 703593f; Phase 11 plans are committed; 11-01 implementation is active in the uncommitted working tree, 11-02 has partial observation/refresh scaffolding, and 11-03 has not started
last_updated: "2026-09-09T07:03:53Z"
last_activity: 2026-09-09
last_activity_desc: Reconciled planning state with local master HEAD and the live uncommitted Phase 11 implementation
state_head: 703593f
progress:
  total_phases: 16
  completed_phases: 10
  total_plans: 18
  completed_plans: 12
  percent: 63
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, `.planning/rebaseline/2026-09-08-core-boundary.md`, and `.planning/phases/11-mcp-runtime-completion/11-WORKING-STATE.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 11 — MCP Runtime Completion
Status: Phase 10 is complete on local `master`; Phase 11 planning is committed and Phase 11 implementation is active but uncommitted and not yet verified green.
Local master HEAD: `703593f` (`docs(10): record integration review`)
Relative to `origin/master`: local master is ahead by 21 commits.
Working tree: 18 changed entries — 16 tracked modifications plus 2 untracked Phase 11 source files; tracked shortstat is 689 insertions / 317 deletions.

The previous state saying local master remained at `eca6cc6`, Phase 10 integration was pending, and Phase 11 had not started is superseded by the live Git/worktree state recorded on 2026-09-09.

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

### Phase 10 — Boundary cleanup (complete on local master)

- BOUNDARY-01 ✅: Planner/Executor/Reviewer workflow surface/domain is removed from ADM Core while generic single-command Runs remain.
- BOUNDARY-02 ✅: no GSD `.planning` interpretation/state-advance API was merged or introduced.
- BOUNDARY-03 ✅: files/exec, verifier, process, MCP, Skill, managed worktree, ordinary Environment and generic Run regression gates passed during Phase 10 closeout.

### Phase 11 — MCP Runtime Completion (active, uncommitted)

#### Plan 11-01 — Typed MCP Configuration + HTTP/Stdio Runtime

Status: in progress, late implementation stage; not closed.

Current working tree contains the core typed MCP model/catalog split, Streamable HTTP + stdio activation/connect paths, stdio allowlist authority path, Gateway transport-specific configuration work, Ping-based explicit health probing and acceptance work.

Closure gaps still include:

- CLI `mcp add` remains the old HTTP-oriented `--name + --endpoint` surface;
- Desktop/shared product-surface parity is not established;
- secret/reference persistence boundaries require hardening so credential-bearing literals cannot become ordinary persisted/returned config;
- no 11-01 summary/verification/evidence has been recorded;
- current worktree has not been established as green by the 2026-09-09 state reconciliation.

#### Plan 11-02 — Health Monitor, Automatic Reconnect, Inventory + Diagnostics

Status: partially started, not complete.

Current working tree includes owner-local MCP runtime observation, inspect/refresh paths, Ping/inventory plumbing and reconnect-related observation fields.

Required background behavior is still missing: no periodic `HealthCheckEnabled` monitor scheduling, no `AutoReconnect` scheduler, and no fixed-interval use of `NextReconnectAt` / `ReconnectInFlight` was found in the live source snapshot.

#### Plan 11-03 — JSON / JSONC Import Adapters

Status: not started in code.

No `mcp_import_preview` implementation was found in the repository snapshot. The committed 11-03 plan remains planning-only.

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

1. **Secret boundary:** `HeaderRefs` / `EnvRefs` are named references, but the current typed config normalization does not by itself prove that credential-bearing literals cannot be persisted or returned. Phase 11-01 must harden this before closure.
2. **Product surface parity:** Gateway typed configuration is ahead of CLI/Desktop surfaces.
3. **Background health/recovery:** observation fields and explicit refresh exist, but periodic Ping scheduling and fixed-interval automatic reconnect are not implemented yet.
4. **Import:** Phase 11-03 remains unimplemented.
5. **Verification:** the current uncommitted Phase 11 worktree must not be called green until focused/full tests and Phase evidence are recorded.
6. **Multi-session continuity:** all sessions must use this STATE plus `11-WORKING-STATE.md` rather than older chat-only assumptions.

## Deferred

- automatic Memory context composition;
- evidence-first investigation helpers until MCP/Skill/capability Core is complete (endpoint resolution, symbol/reference/write tracing, data lineage, symbol-scoped Git history, test-data metric explanation, debug-SQL reverse mapping, semantic consistency);
- Desktop feature expansion beyond blockers until validated Core parity phase;
- installer/tray/autostart/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Stopped at: local `master` HEAD `703593f` already contains Phase 10 integration review and is 21 commits ahead of `origin/master`. Phase 11 plans are committed. The current working tree contains active, uncommitted Phase 11 implementation: 11-01 is in late implementation but has product-surface/secret-boundary/verification gaps; 11-02 has observation/refresh/inventory scaffolding but not periodic monitor/automatic reconnect; 11-03 is not implemented.

Detailed live snapshot: `.planning/phases/11-mcp-runtime-completion/11-WORKING-STATE.md`.

Next action: close Plan 11-01 before advancing 11-02/11-03. Harden secret/reference persistence boundaries, finish required typed configuration surfaces, add/finish focused negative and transport acceptance tests, establish a green verification result, and write 11-01 summary/verification evidence. Do not automatically merge or push.
