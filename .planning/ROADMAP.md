# AI Dev Manager V2 — Roadmap Reset

This roadmap replaces the old sequential Phase 01–07 development order as the active implementation plan. Old `docs/PHASE_*` files remain historical evidence until they are reconciled; they are not the current execution queue.

The roadmap is dependency-driven. A newly discovered issue enters the current Phase only when it blocks that Phase exit criteria or invalidates a locked product requirement.

## Milestone R0 — Planning Rebaseline

**Goal:** Restore truthful project state and a GSD-compatible planning boundary before more feature work.

### Phase 00 — Rebaseline and Planning Authority

**Scope:**
- create `.planning/PROJECT.md`, `ROADMAP.md`, `STATE.md`;
- classify existing V2 capabilities as validated / partial / absent;
- explicitly stop calling Skill/MCP CRUD “complete runtime integration”;
- make `.planning/` the active development queue while `PRODUCT_CONTRACT.md` remains product-semantics authority;
- preserve existing working code; no feature implementation in this Phase.

**Exit criteria:**
1. Current capability truth is recorded in-repo.
2. The next Phase is explicit and has a plan.
3. Deferred Desktop/package work is frozen.
4. Git working state contains planning/docs changes only for this Phase.

---

## Milestone R1 — Agent-ready Development Context

**Goal:** Complete the capabilities an external Agent needs before adding lifecycle/orchestration abstractions: real Skills, structured verification, and real external MCP runtime semantics.

### Phase 01 — Real Skill Runtime + GSD Bootstrap

**Goal:** Replace the current text-record Skill approximation with a real discoverable Environment-scoped Skill capability, and use the actual GSD Skill as the acceptance case.

**Scope:**
- explicit Skill roots / definitions; no arbitrary disk scan;
- discover actual Skill artifacts (initial contract centered on `SKILL.md`);
- stable Skill identity and source/path metadata;
- Environment enable/disable gates runtime access;
- Agent-facing list/read capability through ADM;
- configure one real global GSD Skill without copying it into the project;
- use that GSD Skill through ADM to guide the next Phase.

**Non-goals:**
- Agent Run orchestration;
- GSD automatic phase executor;
- UI redesign;
- generic plugin marketplace.

**Exit criteria:**
1. One global GSD Skill is discovered from its real installation path.
2. Enabled Environment can list/read the actual GSD Skill through ADM.
3. Disabled Environment cannot access it.
4. A second Environment can use the same global GSD Skill without copying files.
5. The next Phase planning is performed using the Skill through ADM, not by pretending the current instructions field is equivalent.
6. Full Go tests/vet and focused real Gateway acceptance pass.

### Phase 02 — Structured Verifier Runtime

**Goal:** Make “change code → verify it” a first-class Agent capability instead of ad-hoc `exec` conventions.

**Scope:**
- optional structured verifier definitions (test/lint/build/custom);
- Environment-root-aware execution;
- timeout/output bounds and stable result schema;
- `run_verifier` / `run_verifiers` Agent Gateway surface;
- verifier absence remains operation-local.

**Exit criteria:**
1. A non-Git Environment can run a configured verifier.
2. Pass/fail/timeout identify the verifier and structured result.
3. No verifier configuration does not block file/exec development.
4. Real V2 `go test ./...` can be invoked through ADM verifier capability.

### Phase 03 — External MCP Runtime Completion

**Goal:** Turn current HTTP proxying into a real Environment-scoped external MCP runtime lifecycle.

**Scope:**
- structured MCP definition with transport-specific connection data;
- keep Streamable HTTP; add stdio only as needed for the real fixture/GSD ecosystem;
- activation/probe and health state;
- Environment enable/disable gates runtime access;
- secret/env references resolved only at activation boundary;
- no false healthy state;
- retain direct tool discovery/call through ADM.

**Exit criteria:**
1. A real local external MCP initializes and lists tools through ADM.
2. Environment A enabled / Environment B disabled produces isolated access.
3. Disabled means no activation/call.
4. Missing/bad configuration reports structured error without secret leakage.
5. At least one real call succeeds end-to-end through the ADM Gateway.

### Phase 04 — External Agent Dogfood Gate

**Goal:** Prove the R1 capabilities as one development loop using only ADM.

**Scope:**
- use one real registered project and Environment;
- consume GSD Skill through ADM;
- inspect/search/edit/delete files;
- run allowlisted command where needed;
- run structured verifier;
- exercise one configured external MCP if the target task benefits from it;
- explicitly read/write scoped Memory where useful;
- record blockers, but only blockers to this gate are fixed in Phase 04.

**Exit criteria:**
1. External Agent completes a real code change without falling back to another filesystem/exec MCP.
2. GSD Skill is actually consumed through ADM.
3. Verifier completes through ADM.
4. Any external MCP used by the task is accessed through Environment gating.
5. R1 closes before any new Desktop feature work.

---

## Milestone R2 — Persistent Development Runtime

**Goal:** Give long-lived development activity a stable local ownership boundary, then add only the lifecycle capabilities real development needs.

### Phase 05 — Persistent Runtime Ownership

**Goal:** Long-lived Runtime/MCP/process ownership survives CLI/Desktop calls and is rebuilt correctly after restart.

**Scope:**
- one persistent local owner boundary built around the existing Gateway/control process rather than a second competing daemon;
- desired vs observed runtime state;
- Runtime/MCP session ownership and cleanup;
- restart reconciliation;
- no serialization of live Go sessions/listeners/process objects.

**Exit criteria:**
1. Start/status/stop from independent invocations observe the same owner.
2. Owned MCP/runtime state survives client exit.
3. Restart rebuilds desired runtime state and never reports dead sessions healthy.
4. Clean shutdown closes owned resources.

### Phase 06 — Dev Process / Logs / Ports

**Goal:** Support the first long-running development-process vertical slice.

**Scope:**
- start/list/status/stop stable process identities;
- bounded stdout/stderr capture;
- process exit facts;
- known/listening port facts where discoverable;
- Environment/writer/policy boundaries remain authoritative.

**Exit criteria:**
1. Start a real local dev server through ADM and return control immediately.
2. Query logs/status from a later invocation.
3. Report the serving port as a fact.
4. Stop it deterministically.
5. Gateway/control restart behavior is explicit and tested.

### Phase 07 — Optional Git Worktree Isolation

**Goal:** Allow isolated parallel task roots for Git projects without redefining Environment around Git.

**Scope:**
- managed worktree create/list/destroy capability;
- Environment may be created on a managed worktree root;
- safe dirty/unpushed destruction guard;
- branch/base facts;
- no automatic merge/rebase/commit/push;
- non-Git Environment path remains fully supported.

**Exit criteria:**
1. Non-Git Environment tests remain green.
2. Two managed worktree Environments from one Git Workspace have isolated roots.
3. Main checkout is not switched or modified by worktree creation.
4. Unsafe destroy is refused unless explicit force policy is satisfied.
5. Missing/tampered managed worktree is detected before routed mutation.

---

## Milestone R3 — Agent Runs and GSD Execution

**Goal:** Add orchestration only after the development substrate and persistent ownership are real.

### Phase 08 — Agent Run Lifecycle

**Scope:** stable Run ID, start/list/status/cancel, persistent-owner execution, explicit restart semantics.

**Exit criteria:** a Run outlives the launching client, is observable/cancellable later, and stale runs are not resurrected after owner restart.

### Phase 09 — Planner / Executor / Reviewer Contract

**Scope:** structured Plan/Step/Review, verifier-backed workflow, review-fail vs orchestration-error separation, auditable run trace.

**Exit criteria:** one real deterministic workflow plans, executes, verifies, reviews, and exposes all states through Run status.

### Phase 10 — GSD Phase Executor

**Scope:** load real `.planning/PROJECT.md`, `STATE.md`, Phase CONTEXT/PLAN; execute only audited operation specs/capabilities; advance STATE only after review pass to a pre-existing unambiguous next plan.

**Exit criteria:** same safety properties previously validated in the earlier ADM GSD implementation are reproduced in V2 against V2 Runtime contracts.

### Phase 11 — Parallel Runs / Worktrees

**Scope:** parent/child run aggregation, optional managed worktree lanes, verifier aggregation, safe cleanup, no automatic integration.

**Exit criteria:** at least two isolated lanes execute concurrently without sharing a writable root; review/cleanup state is auditable.

---

## Milestone R4 — Human Manager Parity

**Goal:** Make Desktop a clear management surface for the Core that now exists. Do not invent Core semantics in the UI.

### Phase 12 — Desktop Core Parity

**Scope:** expose real Skill roots/status, MCP runtime/health, verifiers, Runtime/process/log/port status, Environment isolation facts, Agent Run status. Reuse existing application services.

**Exit criteria:** every Desktop control maps to an already-validated Core operation; no Desktop-only product state exists.

### Phase 13 — Distribution Only If Needed

Installer, shortcuts, tray, autostart, signing, updater, notifications and other distribution/polish work are considered only after real daily dogfood proves a specific need.

---

## Execution Rules

1. Read `.planning/STATE.md` before changing code.
2. Read `PROJECT.md` locked decisions and the current Phase CONTEXT/PLAN.
3. Implement only the current plan.
4. A new issue enters the current Phase only if it blocks the exit criteria or contradicts a locked requirement.
5. Otherwise record it under STATE backlog and continue.
6. No Phase is complete from unit tests alone when the exit criteria require a real Gateway/Desktop/process acceptance.
7. Update STATE and requirement status at Phase completion.
8. Stop at milestone boundaries for review unless `workflow.auto_advance` is deliberately enabled after the real GSD Skill bootstrap.
