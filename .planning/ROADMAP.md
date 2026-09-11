# Roadmap: AI Dev Manager V2

## Overview

ADM V2 is a local AI development control plane. After the green `v1.0.1` release line, the roadmap is organized around post-1.0 management UX and Agent usability: a navigable Desktop surface, large-workspace discovery, compact Agent context, async long-operation observability, temporary task Environments, CLI provisioning and conditional polish.

The 2026-09-08 rebaseline explicitly removes Agent/GSD orchestration from ADM's product scope. Phase 9 remains Git history but is superseded product direction; the Phase 10 GSD branch is abandoned and will not be merged.

The 2026-09-09 RC-first replan is complete. The 2026-09-11 post-1.0 phase map deliberately puts Desktop management UX first because daily operation is now the most visible bottleneck; later Agent-facing and distribution work remains bounded by dogfood evidence.

Milestones:

- **R1 — Development Foundations:** Phases 1-4 (complete history)
- **R2 — Persistent Local Runtime:** Phases 5-8 (complete core history)
- **R3 — Core Boundary + MCP/Skill Completion:** Phases 10-13
- **R4 — 1.0 RC Readiness:** Phase 16 Desktop Core Parity first, then only RC blockers
- **R5 — Stable 1.0 line:** Phase 17 closeout plus the `v1.0.1` green release hotfix
- **R6 — Post-1.0 Management + Agent UX:** Phases 18-24, starting with Desktop management UX
- **R7 — Conditional Post-1.0 Polish:** Phases 25-26 only when dogfood proves need

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
- [x] **Phase 14: Evidence-first Investigation Toolkit** — complete through 14-02 optional code intelligence provider + GitNexus integration boundary; additional helper slices remain deferred until dogfood proves need.
- [x] **Phase 15: Temporary Resource Lifecycle** — complete through temporary Env/MCP/Skill/provider retention metadata, safe cleanup and bounded capability inspection dogfood fix.
- [x] **Phase 16: Desktop Core Parity + 1.0 RC Readiness** — complete; Desktop RC path, packaging, CI, tray/autostart, docs and Release automation are closed.
- [x] **Phase 17: Distribution Only If Needed** — complete through `v1.0.1` green release hotfix; installer/updater/signing/notifications remain conditional.
- [ ] **Phase 18: Desktop Management UX Reorganization** — 18-01 implementation and automated/browser/Wails-build gates complete; visible exact-artifact native Wails acceptance pending. 18-02/03/04 remain sequential successors.
- [ ] **Phase 19: Workspace Discovery + Project Navigation** — bounded large-workspace project candidate discovery and tree digest.
- [ ] **Phase 20: Agent Context Bundle + Capability Injection** — compact Environment context for Agents: tree/capabilities/MCP/Skill/verifier/run guidance.
- [ ] **Phase 21: Async Verifier + Long Operation Observability** — async verifier lifecycle and better long-operation diagnostics.
- [ ] **Phase 22: Temporary Task Environments + Safe Cleanup Workflow** — task-scoped temporary Environment workflow, cleanup and promotion.
- [ ] **Phase 23: CLI Agent UX + MCP/Skill Provisioning** — clearer CLI setup/import/enable/diagnostics for MCP, Skill and Agent use.
- [ ] **Phase 24: Desktop/CLI Surface Boundary Split** — logical surface separation without splitting the Core model.
- [ ] **Phase 25: Distribution Polish If Needed** — installer/updater/signing/notifications only if post-1.0 dogfood proves need.
- [ ] **Phase 26: Evidence-first Investigation Expansion** — optional deeper investigation helpers only with concrete dogfood blockers.

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

**Status:** Complete through 14-02. Additional evidence-first helper slices are deferred by default and should only be opened when dogfood shows generic Runtime/search is insufficient.

**Candidate slices:** endpoint resolution; optional code intelligence provider integration; GitNexus integration; symbol/reference/write tracing; impact/blast-radius evidence; response/data lineage; symbol-scoped Git history/diff; concrete test-data metric explanation; debug-SQL-to-code reverse mapping; semantic consistency checks.

**Provider direction:** GitNexus and similar code graph tools may be integrated as optional providers through existing MCP/runtime authorization. ADM should normalize provider evidence and preserve fallback static heuristics instead of embedding a full code graph engine.

**Plans:**

- [x] `14-01-PLAN.md` — endpoint evidence resolution for URL/path + optional HTTP method. (2026-09-09)
- [x] `14-02-PLAN.md` — optional code intelligence provider + GitNexus integration boundary. (2026-09-11)
- [ ] Additional evidence-first slices after explicit dogfood blocker.

### Phase 15: Temporary Resource Lifecycle

**Goal:** Give ADM-owned temporary Env/MCP/Skill/provider resources explicit ownership, TTL, attachment and cleanup semantics.

**Status:** Complete. Phase 15 shipped explicit retention metadata, conservative cleanup inspect/dry-run/execute surfaces, temporary Skill/MCP cleanup, temporary Environment state-only cleanup, managed worktree cleanup through existing destroy safety, and bounded capability inspection to fix the `pjadm` dogfood timeout/bloat path.

**What this means to the user:** MCP/Gateway-created temporary resources can be listed, preview-cleaned and safely removed, while CLI/UI-created resources remain durable by default.

**Success Criteria:**

1. Distinguish durable vs temporary resources in persisted metadata without changing existing durable defaults.
2. Track creator surface, owner/session/run identity, last-used time, expiration/retention policy and cleanup eligibility.
3. Provide preview/dry-run before destructive cleanup.
4. Refuse cleanup when active writer/process/run, dirty/unpublished worktree or ambiguous ownership makes deletion unsafe.
5. Expose the same lifecycle semantics through Core surfaces first; Desktop consumes them later.

**Plans:**

- [x] `15-01-PLAN.md` — temporary resource metadata and safe cleanup preview/execute. (2026-09-11)
- [x] `15-02-PLAN.md` — temporary Environment cleanup, managed worktree cleanup safety and bounded capability inspection dogfood fix. (2026-09-11)

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

**Status:** Complete for the stable 1.0 line. `v1.0.0` was tagged but its GitHub Actions release run failed on a Windows shutdown race in a nested Gateway acceptance test. The shutdown race was fixed and `v1.0.1` completed the GitHub Actions release workflow successfully. Installer, updater, signing and notifications remain post-1.0 work only if new manual dogfood proves a concrete daily-use blocker.

**Plans / release records:**

- [x] `17-01-PLAN.md` — conditional distribution decision gate and local dogfood artifact refresh using existing build paths; no new distribution implementation required. (2026-09-11)
- [x] `.planning/releases/v1.0.0.md` — stable release closeout, local gate evidence, Windows artifact checksums and tag/publish decision. Superseded by `v1.0.1` because the remote release workflow failed. (2026-09-11)
- [x] `.planning/releases/v1.0.1.md` — release hotfix record for the Windows shutdown race and green GitHub Release. (2026-09-11)

## Post-1.0 Phase Map

See `.planning/post-1.0/PHASE-MAP.md` for the high-level phase sequence. Detailed `PLAN.md` files are intentionally deferred until a phase is actively started.

### Phase 18: Desktop Management UX Reorganization

**Goal:** Reorganize Desktop from a dense management panel into a navigable management application.

**Priority:** P1 / 18-01 implementation complete with native visible-window acceptance pending; 18-02/03/04 are detailed successors and remain blocked by predecessor gates.

**Requirements:** ADM-DESKTOP-001, ADM-MGMT-001, ADM-CORE-015/017/020; preserve ADM-CORE-003/004/005/007/012/013, metadata-safe lifecycle and existing PROC/ARUN boundaries.

**Direction:** compact current connection/Environment header, grouped menu, ten routes to existing panels, truthful dashboard and diagnostics through existing Environment inspection.

**Non-goals:** feature-page redesign in 18-01, new Core APIs/schema/persistence, frontend framework/build pipeline, task orchestration, discovery/context bundle, async verifier, new cleanup workflow, CLI split and distribution work.

**New prerequisites:** none.

**Success criteria:**

1. All existing management areas/actions remain reachable through menu/shortcuts; navigation preserves context, filters and modal drafts.
2. Dashboard uses sanitized existing snapshot data; unloaded/error is distinct from zero and configured inventory is distinct from Runtime health.
3. Auxiliary read failures stay local; profile/Environment switches clear old-scope data and reject late old-scope renders.
4. Optional Git/verifier/MCP/Skill/writer availability does not become a global shell prerequisite; Memory values and mutating/probing operations remain explicit.
5. Real browser and Wails interaction evidence validates routing, focus, minimum-size layout and feature preservation in addition to existing Go/binding gates.

**Plans:** 4 detailed sequential plans

- [ ] `.planning/phases/18-desktop-management-ux/18-01-PLAN.md` — menu shell, dashboard and section routing; implementation/automated acceptance complete, native visible-window acceptance pending before plan acceptance.
- [ ] `.planning/phases/18-desktop-management-ux/18-02-PLAN.md` — Workspace/Environment management readability and existing Runtime subviews/output; waits for accepted 18-01.
- [ ] `.planning/phases/18-desktop-management-ux/18-03-PLAN.md` — MCP/Skill hierarchy, explicit Memory scopes, existing system controls and diagnostics; waits for accepted 18-02.
- [ ] `.planning/phases/18-desktop-management-ux/18-04-PLAN.md` — mandatory integrated Desktop acceptance/closeout and evidence-backed fixes; waits for accepted 18-03.

**Design / validation:** `.planning/phases/18-desktop-management-ux/18-CONTEXT.md`, `18-UI-REFACTOR.md`, `18-VALIDATION.md` and `18-PLANNING-LOG.md`. Detailed successor plans are deliberately written early for continuation safety but must be calibrated against predecessor implementation evidence.

### Phase 19: Workspace Discovery + Project Navigation

**Goal:** Help ADM and Agents understand large workspace roots such as `projects/p1`, `projects/p2`, `projects/p3`.

**Priority:** P1.

**Direction:** bounded project candidate discovery, tree digest, likely root summaries and suggested Environment roots. No arbitrary full-disk indexing.

### Phase 20: Agent Context Bundle + Capability Injection

**Goal:** Give Agents one compact Environment context bundle instead of forcing them to probe many tools.

**Priority:** P1.

**Direction:** root, bounded tree digest, enabled MCP/Skill summary, verifier/run guidance and unavailable capability reasons. No ADM task orchestration.

### Phase 21: Async Verifier + Long Operation Observability

**Goal:** Make heavy verification flows observable and resilient across client/tool timeouts.

**Priority:** P1.

**Direction:** async verifier start/list/status/cancel, bounded live output and clearer sync diagnostics for long operations.

### Phase 22: Temporary Task Environments + Safe Cleanup Workflow

**Goal:** Let Agents create task-scoped temporary work areas without polluting durable project state.

**Priority:** P1/P2.

**Direction:** temporary Environment creation, TTL/owner/run attachment, cleanup preview/execute, promote-to-durable and conservative blockers for dirty/unpublished/active work.

### Phase 23: CLI Agent UX + MCP/Skill Provisioning

**Goal:** Make CLI setup and automation smoother for MCP, Skill, Environment and Agent workflows.

**Priority:** P2.

**Direction:** clearer import/preview/apply, enable/disable, refresh/status/diagnostics and script-friendly context output.

### Phase 24: Desktop/CLI Surface Boundary Split

**Goal:** Keep Desktop and CLI independently optimizable without forking product semantics.

**Priority:** P2/P3.

**Direction:** logical surface/package split inside the repository first; no immediate multi-repo split or second state model.

### Phase 25: Distribution Polish If Needed

**Goal:** Add installer/updater/signing/notifications only when daily use proves they matter.

**Priority:** P3 conditional.

**Direction:** keep on standby until concrete dogfood evidence exists.

### Phase 26: Evidence-first Investigation Expansion

**Goal:** Add deeper investigation helpers only when real tasks show generic Runtime/search is insufficient.

**Priority:** Conditional.

**Direction:** each slice must return evidence, confidence and uncertainty with safe fallback behavior.

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
| 14 | Complete through optional code intelligence provider + GitNexus boundary; later helpers deferred until dogfood need |
| 15 | Complete — temporary resource lifecycle, cleanup safety and bounded capability inspection |
| 16 | Complete — Desktop RC path closed through tray/autostart, packaging/docs, canonical path CI fixes, tray event fix and tag-triggered Release automation |
| 17 | Complete — stable 1.0 line closed through `v1.0.1`; distribution implementation remains conditional |
| 18 | 18-01 implementation + automated/browser/Wails-build gates complete; native visible-window acceptance pending; 18-02 blocked |
| 19 | Planned — Workspace Discovery + Project Navigation |
| 20 | Planned — Agent Context Bundle + Capability Injection |
| 21 | Planned — Async Verifier + Long Operation Observability |
| 22 | Planned — Temporary Task Environments + Safe Cleanup Workflow |
| 23 | Planned — CLI Agent UX + MCP/Skill Provisioning |
| 24 | Planned — Desktop/CLI Surface Boundary Split |
| 25 | Standby — Distribution Polish If Needed |
| 26 | Standby — Evidence-first Investigation Expansion |

## Execution Rules

1. Read STATE, PROJECT, current CONTEXT and PLAN before feature code.
2. Product boundary wins over historical implementation.
3. ADM supplies capabilities; Agent/GSD supplies task orchestration.
4. Optional capability failures remain operation-local.
5. MCP/Skill completion requires real consumption and negative acceptance, not catalog CRUD.
6. New high-level investigation helpers require concrete evidence that generic Runtime/search is insufficient.
7. During active development, continue on `master`; make a commit at each clear node. Do not push unless explicitly requested.
8. Post-1.0 rule: Phase 18 Desktop UX is first because human management is the current daily-use bottleneck; later phases should stay phase-level until execution starts.
