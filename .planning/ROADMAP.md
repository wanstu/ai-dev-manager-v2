# Roadmap: AI Dev Manager V2

## Overview

AI Dev Manager V2 is being rebuilt dependency-first: first make external-Agent development context real (Skills, verifiers, MCP runtime), then add persistent lifecycle and optional isolation, then Agent/GSD orchestration, and only then bring Desktop to Core parity. Desktop/package polish stays frozen unless it blocks the active Phase acceptance.

Milestone grouping:

- **R1 — Agent-ready Development Context:** Phases 1-4
- **R2 — Persistent Development Runtime:** Phases 5-7
- **R3 — Agent Runs and GSD Execution:** Phases 8-11
- **R4 — Human Manager Parity:** Phases 12-13

The planning rebaseline completed before Phase 1 and is retained under `.planning/phases/00-rebaseline/`; it is bootstrap history, not a numbered GSD delivery phase.

## Phases

- [x] **Phase 1: Real Skill Runtime + GSD Bootstrap** - Make real global Skills discoverable and Environment-gated through ADM; bootstrap with the actual installed GSD suite. (completed 2026-09-06)
- [x] **Phase 2: Structured Verifier Runtime** - Make change → verify a first-class structured Agent capability. (completed 2026-09-07)
- [x] **Phase 3: External MCP Runtime Completion** - Complete Environment-scoped external MCP activation, health, and real connection semantics. (completed 2026-09-07)
- [ ] **Phase 4: External Agent Dogfood Gate** - Complete a real development loop using only ADM R1 capabilities.
- [ ] **Phase 5: Persistent Runtime Ownership** - Give long-lived runtime state one stable local owner with restart reconciliation.
- [ ] **Phase 6: Dev Process / Logs / Ports** - Add the first real long-running development-process lifecycle.
- [ ] **Phase 7: Optional Git Worktree Isolation** - Add safe managed worktree isolation without making Git an Environment prerequisite.
- [ ] **Phase 8: Agent Run Lifecycle** - Add persistent-owner Agent Run identity, status, and cancellation.
- [ ] **Phase 9: Planner / Executor / Reviewer Contract** - Add structured auditable orchestration on top of validated Runtime/Verifier capabilities.
- [ ] **Phase 10: GSD Phase Executor** - Execute real repository `.planning` phases through controlled ADM capabilities.
- [ ] **Phase 11: Parallel Runs / Worktrees** - Add isolated concurrent Agent lanes after single-run lifecycle is validated.
- [ ] **Phase 12: Desktop Core Parity** - Expose already-validated Core capabilities in the human Manager.
- [ ] **Phase 13: Distribution Only If Needed** - Add installer/tray/autostart/etc only when daily dogfood proves a concrete need.

## Phase Details

### Phase 1: Real Skill Runtime + GSD Bootstrap

**Milestone:** R1 — Agent-ready Development Context
**Goal**: Replace text-record Skill approximation with real discoverable Environment-scoped Skill artifacts and use the actual installed GSD suite as acceptance.
**Depends on**: Nothing (first delivery phase after planning rebaseline)
**Requirements**: SKILL-RUN-01, SKILL-RUN-02, SKILL-RUN-03, SKILL-RUN-04, SKILL-RUN-05
**Success Criteria** (what must be TRUE):

1. ADM discovers real `SKILL.md` artifacts only from explicitly configured roots.
2. An enabled Environment can list/read the actual installed `gsd-next` Skill and its explicitly authorized `gsd-core` supporting workflow through the ADM Gateway.
3. A disabled Environment cannot read the Skill; a second enabled Environment can use the same global installation without copying it.
4. Paths outside configured Skill artifact/support roots are rejected, and a broken Skill does not break unrelated Gateway tools.
5. The Agent consumes GSD through ADM and GSD can correctly read this repository's `.planning` state before Phase 2 begins.

**Plans**: 1 plan

Plans:

- [x] 01-01: Real Skill Vertical Slice

### Phase 2: Structured Verifier Runtime

**Milestone:** R1 — Agent-ready Development Context
**Goal**: Make change → verify a first-class Agent capability instead of an ad-hoc `exec` convention.
**Depends on**: Phase 1
**Requirements**: VERIFY-01, VERIFY-02, VERIFY-03, VERIFY-04
**Success Criteria** (what must be TRUE):

1. A non-Git Environment can run a configured verifier.
2. Pass/fail/timeout return verifier identity and structured bounded results.
3. Missing verifier configuration does not block normal Environment/file development.
4. V2 `go test ./...` can be invoked through the ADM verifier capability.

**Plans**: TBD

### Phase 3: External MCP Runtime Completion

**Milestone:** R1 — Agent-ready Development Context
**Goal**: Turn current HTTP proxying into a real Environment-scoped external MCP runtime lifecycle.
**Depends on**: Phase 2
**Requirements**: MCP-RUN-01, MCP-RUN-02, MCP-RUN-03, MCP-RUN-04, MCP-RUN-05
**Success Criteria** (what must be TRUE):

1. A real local external MCP initializes and lists tools through ADM.
2. Environment enable/disable gates activation and call access.
3. Health distinguishes configured/disabled/healthy/error rather than treating configuration as health.
4. Missing/bad configuration reports a structured error without secret leakage.
5. At least one real external MCP call succeeds end-to-end through the ADM Gateway.

**Plans**: 2 plans

Plans:

- [x] 03-01-PLAN.md — MCP runtime lifecycle tracer: Transport model, health probe, secret resolution, environment_mcp_status tool
- [x] 03-02-PLAN.md — Error enrichment, management/CLI health surface, comprehensive acceptance testing

### Phase 4: External Agent Dogfood Gate

**Milestone:** R1 — Agent-ready Development Context
**Goal**: Prove Skills, verifier, external MCP, files/exec, and scoped Memory as one real external-Agent development loop using only ADM.
**Depends on**: Phase 3
**Success Criteria** (what must be TRUE):

1. An external Agent completes a real code change without another filesystem/exec MCP.
2. GSD Skill is actually consumed through ADM.
3. Verification completes through ADM.
4. Any external MCP used by the task is accessed through Environment gating.
5. R1 closes before new Desktop feature work.

**Plans**: TBD

### Phase 5: Persistent Runtime Ownership

**Milestone:** R2 — Persistent Development Runtime
**Goal**: Give long-lived Runtime/MCP/process ownership one stable local control boundary.
**Depends on**: Phase 4
**Requirements**: LIFE-01, LIFE-02, LIFE-03
**Success Criteria** (what must be TRUE):

1. Independent invocations observe the same local owner.
2. Owned runtime/MCP state survives client exit.
3. Restart rebuilds desired runtime state and never reports dead sessions healthy.
4. Clean shutdown closes owned resources.

**Plans**: TBD

### Phase 6: Dev Process / Logs / Ports

**Milestone:** R2 — Persistent Development Runtime
**Goal**: Support the first long-running development-process vertical slice.
**Depends on**: Phase 5
**Requirements**: PROC-01, PROC-02, PROC-03
**Success Criteria** (what must be TRUE):

1. Start a real local dev server through ADM and return control immediately.
2. Query process logs/status from a later invocation.
3. Report the serving port as a fact.
4. Stop the process deterministically.

**Plans**: TBD

### Phase 7: Optional Git Worktree Isolation

**Milestone:** R2 — Persistent Development Runtime
**Goal**: Allow isolated parallel task roots for Git projects without redefining Environment around Git.
**Depends on**: Phase 6
**Requirements**: ISO-01, ISO-02, ISO-03
**Success Criteria** (what must be TRUE):

1. Existing non-Git Environment behavior remains green.
2. Two managed worktree Environments from one Git Workspace have isolated roots.
3. Main checkout is not switched or modified by worktree creation.
4. Unsafe destroy is refused unless explicit force policy is satisfied.
5. Missing/tampered managed worktree is detected before routed mutation.

**Plans**: TBD

### Phase 8: Agent Run Lifecycle

**Milestone:** R3 — Agent Runs and GSD Execution
**Goal**: Add stable Agent Run identity, lifecycle status, and cancellation under persistent ownership.
**Depends on**: Phase 7
**Requirements**: ARUN-01
**Success Criteria** (what must be TRUE):

1. A Run outlives the launching client.
2. A later client can list/status/cancel it by stable identity.
3. Owner restart does not resurrect stale observed Runs.

**Plans**: TBD

### Phase 9: Planner / Executor / Reviewer Contract

**Milestone:** R3 — Agent Runs and GSD Execution
**Goal**: Add structured auditable orchestration and verifier-backed review semantics.
**Depends on**: Phase 8
**Requirements**: FLOW-01
**Success Criteria** (what must be TRUE):

1. One real deterministic workflow plans, executes, verifies, and reviews.
2. Review failure is distinct from orchestration/runtime failure.
3. Plan/steps/review remain visible in Run status.

**Plans**: TBD

### Phase 10: GSD Phase Executor

**Milestone:** R3 — Agent Runs and GSD Execution
**Goal**: Execute repository GSD phases using controlled ADM Runtime capabilities and verifier/reviewer gates.
**Depends on**: Phase 9
**Requirements**: GSD-01, GSD-02, GSD-03
**Success Criteria** (what must be TRUE):

1. V2 reads real PROJECT/STATE/CONTEXT/PLAN provenance.
2. GSD execution uses an explicit operation allowlist with no hidden arbitrary shell path.
3. STATE advances only after verified pass to a pre-existing unambiguous next plan.

**Plans**: TBD

### Phase 11: Parallel Runs / Worktrees

**Milestone:** R3 — Agent Runs and GSD Execution
**Goal**: Add isolated concurrent Agent lanes after single-run lifecycle and optional worktree isolation are validated.
**Depends on**: Phase 10
**Requirements**: PAR-01
**Success Criteria** (what must be TRUE):

1. At least two isolated lanes execute concurrently without sharing a writable root.
2. Parent review/aggregation is auditable.
3. Cleanup is safe and never performs automatic integration decisions.

**Plans**: TBD

### Phase 12: Desktop Core Parity

**Milestone:** R4 — Human Manager Parity
**Goal**: Make Desktop a clear management surface for Core capabilities that already exist.
**Depends on**: Phase 11
**Success Criteria** (what must be TRUE):

1. Desktop exposes real Skills, MCP runtime/health, verifier, process/runtime, isolation, and Agent Run facts.
2. Every Desktop control maps to an already-validated Core operation.
3. No Desktop-only product state or Core semantics exist.

**Plans**: TBD

### Phase 13: Distribution Only If Needed

**Milestone:** R4 — Human Manager Parity
**Goal**: Add distribution/polish only when real daily use demonstrates a concrete need.
**Depends on**: Phase 12
**Success Criteria** (what must be TRUE):

1. Installer/tray/autostart/updater/signing/notifications are implemented only against a demonstrated user need.
2. Distribution work does not redefine Core behavior.

**Plans**: TBD

## Progress

| Phase | Milestone | Plans Complete | Status | Completed |
|-------|-----------|----------------|--------|-----------|
| 1. Real Skill Runtime + GSD Bootstrap | R1 | 1/1 | Complete    | 2026-09-06 |
| 2. Structured Verifier Runtime | R1 | 1/1 | Complete    | 2026-09-07 |
| 3. External MCP Runtime Completion | R1 | 2/2 | Complete    | 2026-09-07 |
| 4. External Agent Dogfood Gate | R1 | 0/TBD | Not started | - |
| 5. Persistent Runtime Ownership | R2 | 0/TBD | Not started | - |
| 6. Dev Process / Logs / Ports | R2 | 0/TBD | Not started | - |
| 7. Optional Git Worktree Isolation | R2 | 0/TBD | Not started | - |
| 8. Agent Run Lifecycle | R3 | 0/TBD | Not started | - |
| 9. Planner / Executor / Reviewer Contract | R3 | 0/TBD | Not started | - |
| 10. GSD Phase Executor | R3 | 0/TBD | Not started | - |
| 11. Parallel Runs / Worktrees | R3 | 0/TBD | Not started | - |
| 12. Desktop Core Parity | R4 | 0/TBD | Not started | - |
| 13. Distribution Only If Needed | R4 | 0/TBD | Not started | - |

## Execution Rules

1. Read `.planning/STATE.md`, `PROJECT.md`, current Phase CONTEXT and PLAN before changing code.
2. Implement only the active plan.
3. A new issue enters the active Phase only when it blocks that Phase success criteria or invalidates a locked product requirement.
4. Otherwise record it for later and continue.
5. Unit tests alone do not close a Phase whose success criteria require real Gateway/process/Desktop acceptance.
6. Update STATE and requirement truth at Phase completion.
7. Stop at milestone boundaries for review unless real GSD automation is deliberately enabled later.
