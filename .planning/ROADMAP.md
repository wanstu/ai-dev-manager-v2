# Roadmap: AI Dev Manager V2

## Overview

ADM V2 is a local AI development control plane. The roadmap is now organized around product capabilities an external Agent actually consumes: MCP, Skill, safe local Runtime, Environment capability control and diagnostics.

The 2026-09-08 rebaseline explicitly removes Agent/GSD orchestration from ADM's product scope. Phase 9 remains Git history but is superseded product direction; the Phase 10 GSD branch is abandoned and will not be merged.

Milestones:

- **R1 — Development Foundations:** Phases 1-4 (complete history)
- **R2 — Persistent Local Runtime:** Phases 5-8 (complete core history)
- **R3 — Core Boundary + MCP/Skill Completion:** Phases 10-13
- **R4 — Agent Investigation + Human Management:** Phases 14-16

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
- [ ] **Phase 10: Orchestration Boundary Cleanup** — remove ADM-owned workflow orchestration and restore the Core boundary.
- [ ] **Phase 11: MCP Runtime Completion** — make external MCP a complete first-class ADM capability.
- [ ] **Phase 12: Skill Runtime Completion** — make Skill discovery/access/availability a complete first-class ADM capability.
- [ ] **Phase 13: Environment Capability Diagnostics** — answer what is usable in an Environment, what is not, and why.
- [ ] **Phase 14: Evidence-first Investigation Toolkit** — add high-value code/runtime investigation helpers after Core completion.
- [ ] **Phase 15: Desktop Core Parity** — expose validated Core capabilities for human management.
- [ ] **Phase 16: Distribution Only If Needed** — installer/tray/autostart/etc only from demonstrated need.

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

- [ ] `10-01-PLAN.md` — remove workflow orchestration, preserve generic Run, prove MCP/Skill/Core regressions.

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

- [ ] `11-01-PLAN.md` — typed MCP desired configuration, health policy model, and real Streamable HTTP/stdio transport activation.
- [ ] `11-02-PLAN.md` — protocol Ping health monitor, configurable automatic reconnect, owner-local inventory/refresh and structured diagnostics.
- [ ] `11-03-PLAN.md` — preview/apply MCP JSON/JSONC import adapters for OpenCode, WorkBuddy/CodeBuddy, Codex plugin MCP JSON, Claude Code and MCPHub.

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

- [ ] `12-01-PLAN.md` — explicit Skill sources, source-aware stable identity and atomic refresh.
- [ ] `12-02-PLAN.md` — Environment Skill availability, bounded support inventory/read and diagnostics.

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

- [ ] `13-01-PLAN.md` — shared CapabilityFact model and resilient application-level Environment report.
- [ ] `13-02-PLAN.md` — Gateway-owner observation enrichment and canonical Agent-facing capability report.

### Phase 14: Evidence-first Investigation Toolkit

**Goal:** Add concrete Agent debugging/navigation helpers only where they improve development accuracy and speed beyond generic text search.

**Candidate slices:** endpoint resolution; symbol/reference/write tracing; response/data lineage; symbol-scoped Git history/diff; concrete test-data metric explanation; debug-SQL-to-code reverse mapping; semantic consistency checks.

**Rule:** each helper must return `evidence`, `confidence`, and `uncertainties`/alternatives where applicable. Do not add opaque AI guesses or orchestration policy.

**Plans:** choose slices from real dogfood blockers after Phases 11-13.

### Phase 15: Desktop Core Parity

**Goal:** Make validated MCP/Skill/Environment/runtime/diagnostic capabilities manageable by a human without creating Desktop-only semantics.

**Success Criteria:**

1. Desktop exposes the same application-level state/operations as Core for MCP, Skill, Environment capability facts, process/runtime, verifier and isolation.
2. No Desktop-only product state or authorization model exists.
3. UI work does not block or redefine Agent-facing Core behavior.

### Phase 16: Distribution Only If Needed

**Goal:** Add installer/tray/autostart/updater/signing/notifications only when daily use demonstrates a concrete need.

## Historical Phases 1-9

Phases 1-8 remain accepted implementation history. Phase 8's generic asynchronous single-command Run remains Core.

Phase 9 was technically verified and integrated, but the 2026-09-08 product rebaseline determined that Planner/Executor/Reviewer workflow semantics belong to an external Agent/GSD/orchestrator rather than ADM. Phase 10 removes that surface. Historical evidence is retained; it does not justify future dependencies on FLOW-01.

The abandoned `feat/gsd-phase-executor` branch is not a roadmap phase result and must not be merged.

## Progress

| Phase | Status |
|---|---|
| 1-8 | Complete / retained Core history |
| 9 | Historical experiment; integrated but superseded; cleanup required |
| 10 | Implementation in progress — orchestration boundary cleanup |
| 11 | Detailed planned — 3 plans, implementation not started |
| 12 | Detailed planned — 2 plans, implementation not started |
| 13 | Detailed planned — 2 plans, implementation not started |
| 14 | Deferred until Core completion — investigation toolkit |
| 15 | Deferred — Desktop Core parity |
| 16 | Conditional — distribution |

## Execution Rules

1. Read STATE, PROJECT, current CONTEXT and PLAN before feature code.
2. Product boundary wins over historical implementation.
3. ADM supplies capabilities; Agent/GSD supplies task orchestration.
4. Optional capability failures remain operation-local.
5. MCP/Skill completion requires real consumption and negative acceptance, not catalog CRUD.
6. New high-level investigation helpers require concrete evidence that generic Runtime/search is insufficient.
7. No automatic merge/push or next-phase transition without the existing review/authorization boundary.
