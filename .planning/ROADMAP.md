# Roadmap: AI Dev Manager V2

## Overview

ADM V2 is a local AI development control plane. The roadmap is now organized around the capabilities required for a usable 1.0 release candidate: MCP, Skill, safe local Runtime, Environment capability diagnostics, and a human-manageable Desktop surface over the same Core.

The 2026-09-08 rebaseline explicitly removes Agent/GSD orchestration from ADM's product scope. Phase 9 remains Git history but is superseded product direction; the Phase 10 GSD branch is abandoned and will not be merged.

The 2026-09-09 RC-first replan keeps Phase 14/15/17 as important usability and polish work, but makes Desktop Core Parity the next release-critical path. Evidence helpers, GitNexus provider integration, temporary resource lifecycle and distribution polish should not block a 1.0 RC unless dogfood proves they are necessary for basic daily use.

Milestones:

- **R1 — Development Foundations:** Phases 1-4 (complete history)
- **R2 — Persistent Local Runtime:** Phases 5-8 (complete core history)
- **R3 — Core Boundary + MCP/Skill Completion:** Phases 10-13
- **R4 — 1.0 RC Readiness:** Phase 16 Desktop Core Parity first, then only RC blockers
- **R5 — Post-RC Usability + Polish:** remaining Phase 14 evidence helpers, Phase 15 lifecycle, Phase 17 distribution as needed

## Phases

- [x] **Phase 1: Real Skill Foundation + Bootstrap** — real `SKILL.md` discovery/read, support roots and Environment gating. (2026-09-06)
- [x] **Phase 2: Structured Verifier Runtime** — structured optional test/lint/build/custom verification. (2026-09-07)
- [x] **Phase 3: External MCP Foundation** — Streamable HTTP vertical slice, health, secrets, Environment gating and real tool call. (2026-09-07)
- [x] **Phase 4: External Agent Dogfood Gate** — prove real development through ADM-only capabilities. (2026-09-07)
- [x] **Phase 5: Persistent Runtime Ownership** — Gateway-owned live resources and restart reconciliation. (2026-09-07)
- [x] **Phase 6: Dev Process / Logs / Ports** — long-running development process lifecycle. (2026-09-08)
- [x] **Phase 7: Optional Git Worktree Isolation** — safe optional isolated roots. (2026-09-08)
- [x] **Phase 8: Generic Async Run Lifecycle** — stable single-command `run_` start/list/status/cancel. (2026-09-08)
- [x] **Phase 9: Planner / Executor / Reviewer Experiment** — integrated historically, but superseded as out-of-scope orchestration by the 2026-09-08 rebaseline.
- [x] **Phase 10: Orchestration Boundary Cleanup** — removed ADM-owned workflow orchestration and restored the Core boundary. (2026-09-08)
- [x] **Phase 11: MCP Runtime Completion** — complete first-class external MCP runtime, health/recovery, diagnostics and import adapters. (2026-09-09)
- [x] **Phase 12: Skill Runtime Completion** — complete source-aware Skill refresh, availability diagnostics and bounded support inventory/read. (2026-09-09)
- [x] **Phase 13: Environment Capability Diagnostics** — canonical static and Gateway-owner-enriched capability report. (2026-09-09)
- [ ] **Phase 14: Evidence-first Investigation Toolkit** — current resumed phase; 14-01 endpoint evidence is complete, next slice is 14-02 optional code intelligence provider + GitNexus integration boundary.
- [ ] **Phase 15: Temporary Resource Lifecycle** — important cleanup/usability work for temporary Env/MCP/Skill/provider resources; post-RC unless dogfood shows it blocks basic use.
- [x] **Phase 16: Desktop Core Parity + 1.0 RC Readiness** — complete; Desktop RC path, packaging, CI, tray/autostart, docs and Release automation are closed.
- [ ] **Phase 17: Distribution Only If Needed** — installer/tray/autostart/updater/signing/notifications only after the Desktop RC is useful.

## Phase Details

### Phase 10: Orchestration Boundary Cleanup

**Goal:** Remove product behavior that belongs to the Agent/GSD layer while retaining useful generic Runtime primitives.

**Requirements:** BOUNDARY-01, BOUNDARY-02, BOUNDARY-03

**What this means to the user:** ADM stops pretending to plan/review projects. It remains the safe local capability server that an Agent/GSD calls.

**Success Criteria:**

1. `run_start/list/status/cancel` remains a generic asynchronous single-command Runtime capability.
2. `run_workflow_start` and Planner/Executor/Reviewer domain/status code are removed from the Agent-facing ADM Core.
3. No `gsd_phase_*` or `.planning` state-advance implementation is merged.
4. Existing files, exec, verifier, process, MCP, Skill, Git/worktree and ordinary Run behavior remain green.
5. Product Contract and docs describe ADM as infrastructure/control plane, not an orchestrator.

**Plans:** 1 plan

- [x] `10-01-PLAN.md` — remove workflow orchestration, preserve generic Run, prove MCP/Skill/Core regressions. (2026-09-08)

### Phase 11: MCP Runtime Completion

**Goal:** Turn the current MCP foundation into a complete, reliable, importable, self-monitoring and diagnosable external MCP runtime.

**Requirements:** MCP-COMP-01..07

**What this means to the user:** configure or import MCPs once, enable them for an Environment, and ADM reliably connects, exposes tools, checks health, can automatically reconnect unhealthy servers when the user enables that policy, refreshes inventory and explains failures.

**Success Criteria:**

1. MCP definitions have one explicit desired configuration model with supported transport/auth fields separated from observed health/session state.
2. Every advertised transport has a real activation/lifecycle implementation; unsupported transport is rejected locally without affecting other Environment capabilities.
3. Secret/auth values resolve only at activation and do not leak through normal diagnostics or import previews.
4. Agent can inspect/refresh the actual tool inventory and Environment disable revokes access immediately.
5. Enabled MCPs can use configurable protocol health checks with bounded timeout and configurable check interval. Automatic reconnect defaults to off; when enabled it uses the configured fixed reconnect interval with no exponential/adaptive backoff. Health/recovery state is owner-local and stale sessions are never reported healthy.
6. Connection/init/ping/tool-discovery/call/reconnect failures return structured actionable diagnostics. A failed tool call is never automatically replayed.
7. ADM can preview then atomically import one or many definitions from supported JSON/JSONC source formats: OpenCode, WorkBuddy/CodeBuddy, Codex plugin `.mcp.json`, Claude Code, and MCPHub. Name conflicts default to error. Phase 11 import writes only global MCP definitions, never existing Environment selections, and converts literal credential-bearing values into secret/environment-reference requirements instead of persisting them.
8. Real acceptance covers both supported transports, automatic unhealthy→reconnect recovery, and representative single/batch imports from each supported source adapter.

**Non-goals:** ADM does not decide which MCP tool an Agent should call as part of a task plan; importer compatibility does not make external config formats part of ADM's internal model.

**Plans:** 3 plans

- [x] `11-01-PLAN.md` — typed MCP desired configuration, health policy model, and real Streamable HTTP/stdio transport activation. (2026-09-08)
- [x] `11-02-PLAN.md` — protocol Ping health monitor, configurable automatic reconnect, owner-local inventory/refresh and structured diagnostics. (2026-09-08)
- [x] `11-03-PLAN.md` — preview/apply MCP JSON/JSONC import adapters for OpenCode, WorkBuddy/CodeBuddy, Codex plugin MCP JSON, Claude Code and MCPHub. (2026-09-09)

### Phase 12: Skill Runtime Completion

**Goal:** Turn real `SKILL.md` discovery/read into a reliable Skill availability/runtime context capability.

**Requirements:** SKILL-COMP-01..04

**What this means to the user:** ADM can tell an Agent exactly which Skills are installed/enabled/usable, read the real Skill/support files safely, refresh changes, and explain broken Skills.

**Success Criteria:**

1. Explicit refresh updates discovered real Skill artifacts without arbitrary host scanning.
2. Source path/artifact/support-file facts are inspectable and bounded.
3. Broken/missing support references make that Skill unavailable/diagnosable without breaking unrelated Skills.
4. Environment enabled/disabled state remains the authorization gate.
5. Agent gets clear availability states/reasons and can read the content it is authorized to consume.
6. ADM does not create a second Skill interpreter/execution engine; the consuming Agent follows the Skill instructions.

**Plans:** 2 plans

- [x] `12-01-PLAN.md` — explicit Skill sources, source-aware stable identity and atomic refresh. (2026-09-09)
- [x] `12-02-PLAN.md` — Environment Skill availability, bounded support inventory/read and diagnostics. (2026-09-09)

### Phase 13: Environment Capability Diagnostics

**Goal:** Give one authoritative view of what an Agent can actually use in one Environment.

**Requirements:** CAP-01, CAP-02

**What this means to the user:** instead of probing ten tools, ask ADM why a capability works or does not work.

**Success Criteria:**

1. One Environment inspection returns file/exec/verifier/process/Git/isolation/MCP/Skill availability facts.
2. Unavailable capabilities include structured reasons such as disabled, unconfigured, missing executable, broken Skill, unsupported transport, auth error or runtime error.
3. Evidence identifies the relevant ADM definition/resource without leaking secrets/private Memory values.
4. A broken optional capability never marks the whole Environment unusable.

**Plans:** 2 plans

- [x] `13-01-PLAN.md` — shared CapabilityFact model and resilient application-level Environment report. (2026-09-09)
- [x] `13-02-PLAN.md` — Gateway-owner observation enrichment and canonical Agent-facing capability report. (2026-09-09)

### Phase 14: Evidence-first Investigation Toolkit

**Goal:** Add concrete Agent debugging/navigation helpers only where they improve development accuracy and speed beyond generic text search.

**Status:** Current resumed phase after Phase 16 closeout. 14-01 is complete; continue with 14-02 optional code intelligence provider + GitNexus integration boundary.

**Candidate slices:** endpoint resolution; optional code intelligence provider integration; GitNexus integration; symbol/reference/write tracing; impact/blast-radius evidence; response/data lineage; symbol-scoped Git history/diff; concrete test-data metric explanation; debug-SQL-to-code reverse mapping; semantic consistency checks.

**Provider direction:** GitNexus and similar code graph tools may be integrated as optional providers through existing MCP/runtime authorization. ADM should normalize provider evidence and preserve fallback static heuristics instead of embedding a full code graph engine.

**Plans:**

- [x] `14-01-PLAN.md` — endpoint evidence resolution for URL/path + optional HTTP method. (2026-09-09)
- [ ] `14-02-PLAN.md` — optional code intelligence provider + GitNexus integration boundary. Next planned slice after Phase 16 closeout.
- [ ] Additional evidence-first slices after RC or explicit dogfood blocker.

### Phase 15: Temporary Resource Lifecycle

**Goal:** Give ADM-owned temporary Env/MCP/Skill/provider resources explicit ownership, TTL, attachment and cleanup semantics.

**Status:** Planned but deferred behind Phase 16 unless unmanaged temporary resources become a day-to-day RC blocker.

**What this means to the user:** MCP/Gateway-created temporary resources can be listed, preview-cleaned and safely removed, while CLI/UI-created resources remain durable by default.

**Success Criteria:**

1. Distinguish durable vs temporary resources in persisted metadata without changing existing durable defaults.
2. Track creator surface, owner/session/run identity, last-used time, expiration/retention policy and cleanup eligibility.
3. Provide preview/dry-run before destructive cleanup.
4. Refuse cleanup when active writer/process/run, dirty/unpublished worktree or ambiguous ownership makes deletion unsafe.
5. Expose the same lifecycle semantics through Core surfaces first; Desktop consumes them later.

**Plans:**

- [ ] `15-01-PLAN.md` — temporary resource metadata and safe cleanup preview/execute. Planned, deferred behind Phase 16 unless it becomes an RC blocker.

### Phase 16: Desktop Core Parity + 1.0 RC Readiness

**Goal:** Make validated Core capabilities manageable by a human without creating Desktop-only semantics, then cut a local 1.0 RC candidate when the Desktop can support daily use.

**Why this moves ahead:** Phases 10-13 established the foundation: safe local runtime, MCP, Skill, capability diagnostics and Core boundaries. Phase 14/15/17 mostly improve usability or polish. Desktop Core Parity is the missing product surface that lets the validated Core become a usable RC.

**Success Criteria:**

1. Desktop exposes the same application-level state/operations as Core for Workspace, Environment, capability facts, MCP, Skill, process/runtime, verifier and isolation.
2. No Desktop-only product state or authorization model exists.
3. Desktop safely displays private/sensitive state: private Memory values and secret-backed MCP values are not leaked.
4. Desktop can manage daily Core workflows: inspect status, see capability reasons, configure MCP/Skill, run allowed operations, and view process/run/verifier results.
5. Desktop has an explicit RC gate: full tests/vet/diff, manual smoke checklist, known-limitations note, and no release-blocking Core regressions.
6. 1.0 RC is local/no-push unless explicitly requested.

**Plans:**

- [x] `16-01-PLAN.md` — Desktop Core parity inventory, GitHub CI baseline and RC gate. (2026-09-09)
- [x] `16-02-PLAN.md` — Desktop MCP/Skill visual management UI for the first RC user journey. (2026-09-10)
- [x] `16-03-PLAN.md` — local management-plane convergence complete through 16-03D; 16-03E authenticated remote enablement is explicitly deferred post-RC. (2026-09-10)
- [x] RC1 readiness slice — reserved-port blocker fixed, local bind preflight added, Desktop verifier/process/run visibility added, and `v1.0.0-rc.1` local candidate cut. Superseded after dogfood exposed Desktop build and Gateway target-selection blockers. (2026-09-10)
- [x] RC2 blocker-fix slice — Desktop release/CI build now uses Wails, scripts/build-rc.ps1 produces versioned artifacts/checksums, Gateway lifecycle honors the selected ADM Base URL, and real GUI dogfood fixed the broken MCP/Skill CSS plus long-list interaction at `a3478f2`. Exact RC2 artifacts passed local smoke. (2026-09-10)
- [x] RC3 import-fix slice — generic top-level `mcpServers` now auto-detects as `generic-mcpservers`, literal env/header values remain reference-only, and Desktop preview shows generated reference requirements. (2026-09-10)
- [x] RC4/RC5 closeout slice — `adm` / `adm-desktop` naming, unified `dist/` packaging, fitted tray icon, Windows canonical-path CI fixes, tray event loop OS-thread fix and tag-triggered Release automation. Phase 16 closed by user acceptance. (2026-09-10)

### Phase 17: Distribution Only If Needed

**Goal:** Add installer/tray/autostart/updater/signing/notifications only when daily use demonstrates a concrete need after the Desktop RC is useful.

**Status:** Post-RC polish by default. Do not block 1.0 RC on installer or platform packaging unless manual dogfood proves it is necessary.

## Historical Phases 1-9

Phases 1-8 remain accepted implementation history. Phase 8's generic asynchronous single-command Run remains Core.

Phase 9 was technically verified and integrated, but the 2026-09-08 product rebaseline determined that Planner/Executor/Reviewer workflow semantics belong to an external Agent/GSD/orchestrator rather than ADM. Phase 10 removes that surface. Historical evidence is retained; it does not justify future dependencies on FLOW-01.

The abandoned `feat/gsd-phase-executor` branch is not a roadmap phase result and must not be merged.

## Progress

| Phase | Status |
|---|---|
| 1-8 | Complete / retained Core history |
| 9 | Historical experiment; integrated but superseded; workflow surface removed by Phase 10 |
| 10 | Complete / integrated — orchestration boundary cleanup |
| 11 | Complete — MCP runtime, health/recovery/diagnostics and JSON/JSONC import adapters |
| 12 | Complete — Skill source refresh, source/artifact identity, availability diagnostics and support inventory/read |
| 13 | Complete — static CapabilityFact report and Gateway-owner observation enrichment |
| 14 | Current / resumed after Phase 16 — continue from 14-02 GitNexus/provider boundary |
| 15 | Planned but deferred — temporary resource lifecycle is post-RC unless blocking |
| 16 | Complete — Desktop RC path closed through tray/autostart, packaging/docs, canonical path CI fixes, tray event fix and tag-triggered Release automation |
| 17 | Conditional post-RC — distribution |

## Execution Rules

1. Read STATE, PROJECT, current CONTEXT and PLAN before feature code.
2. Product boundary wins over historical implementation.
3. ADM supplies capabilities; Agent/GSD supplies task orchestration.
4. Optional capability failures remain operation-local.
5. MCP/Skill completion requires real consumption and negative acceptance, not catalog CRUD.
6. New high-level investigation helpers require concrete evidence that generic Runtime/search is insufficient.
7. During active development, continue on `master`; make a commit at each clear node. Do not push unless explicitly requested.
8. RC-first rule: prioritize Phase 16 until the Desktop can support a local 1.0 RC. Phase 14/15/17 work should not preempt Phase 16 unless it is a proven RC blocker.
