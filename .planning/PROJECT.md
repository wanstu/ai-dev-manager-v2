# AI Dev Manager V2 — Project Baseline

## Product Mission

AI Dev Manager V2 is a local AI development control plane for external Agents.

Its job is to expose one reliable, inspectable and safe local development boundary over a stable Gateway: Workspace/Environment roots, files, commands, processes, verification, external MCPs, Skills, optional Git/worktree capabilities, scoped Memory and capability diagnostics.

ADM is infrastructure for Agents. It is not the task planner or project orchestrator. An Agent, GSD, or another orchestrator decides what task to do, how to decompose it and when a phase is complete; ADM supplies the controlled local capabilities needed to perform that work.

## Product Priority

1. External MCP Runtime is complete, reliable and diagnosable.
2. Skill Runtime is complete, reliable and diagnosable.
3. Environment capability routing and safety are explicit and consistent.
4. Local files/exec/process/verifier/Git/isolation primitives remain reliable and useful to any external Agent.
5. Evidence-first diagnostics improve Agent accuracy without becoming an orchestration engine.
6. CLI/Desktop expose the same validated Core; packaging/polish come last.

## Locked Architecture Rules

- Workspace is a registered local directory. Git is optional.
- Environment is a development context rooted at a directory. It must not intrinsically require Git, worktree, Docker, verifier, MCP or Skill.
- Missing optional capability blocks only the operation that needs it.
- One physical root has at most one active writer.
- Agent-facing operations use stable ADM identity and Environment authority; arbitrary host paths are not authorization.
- Control/management and Runtime execution are separate responsibilities.
- One long-lived ADM Gateway serves many Environments; users do not configure a new ADM MCP server for every task.
- MCP and Skill are first-class runtime capabilities, not merely persisted catalog records.
- Skill execution semantics belong to the consuming Agent/Skill; ADM safely discovers, exposes and diagnoses Skill artifacts.
- Task planning/orchestration belongs to the Agent/GSD/orchestrator, not ADM.
- ADM must not interpret natural-language plans into commands, advance GSD `.planning` state, choose the next task/phase, or implement Planner/Executor/Reviewer role policy.
- Generic asynchronous Runtime resources are allowed when needed for client-independent lifecycle, but they must remain task-semantic-neutral.
- Worktree is an optional isolation primitive; ADM does not decide parallel Agent policy or automatic integration.
- Important planning decisions live under `.planning/`, not only chat history.
- No capability is complete because CRUD/metadata exists; it requires a real consumption path and acceptance evidence.
- No pre-stable migration/compatibility burden unless explicitly requested.

## Reality Audit — 2026-09-08

### Validated / usable core

- Workspace registration/lifecycle for ordinary directories.
- Environment creation/lifecycle with Workspace-contained roots.
- tree/read/search/write/exact-edit/single-file-delete.
- single-writer lease by physical root.
- allowlisted structured command execution.
- structured verifier definitions and execution.
- persistent Gateway Runtime owner.
- Gateway-owned long-running dev processes with bounded logs and Windows owned-port facts.
- optional Git status/diff/branch.
- optional managed Git worktree isolation with safe destroy/revalidation.
- persistent local ADM desired state.
- Global Memory and Environment-private Memory explicit CRUD.
- CLI/Desktop management shells over the same application state/services.
- real Skill foundation: explicit discovery roots, real `SKILL.md`, support roots, Environment gating, shared installation.
- external MCP foundation: Streamable HTTP connection definition, Environment gating, four-state health, secret resolution, real tool list/call, Gateway-owned persistent sessions/restart reconciliation.
- generic asynchronous single-command `run_` lifecycle: stable start/list/status/cancel across client disconnect; owner-local observation and cleanup.

### Implemented but no longer product-authoritative

Phase 9 introduced `run_workflow_start` and Planner / Executor / Reviewer workflow semantics. The implementation is integrated in Git history, but the 2026-09-08 rebaseline classifies this orchestration layer as outside ADM's product responsibility. It is scheduled for cleanup before further Core feature expansion.

The Phase 10 GSD Phase Executor implementation exists only on the abandoned `feat/gsd-phase-executor` branch and must not be merged.

### Core gaps

#### External MCP Runtime

The current implementation proves one strong Streamable HTTP vertical slice and persistent session ownership, but MCP is not considered complete. Remaining product work includes a coherent configuration/runtime model, supported transport coverage, auth/secret handling, lifecycle/reconnect behavior, tool inventory refresh and actionable diagnostics.

#### Skill Runtime

The current implementation proves real artifact discovery/read and support-root containment, but Skill is not considered complete. Remaining product work includes refresh/source facts, broken-artifact isolation, support/dependency visibility, Environment availability diagnostics and clear explanations of why a Skill is or is not usable.

#### Environment capability diagnostics

ADM can expose many capabilities separately, but it does not yet provide one authoritative Environment view that answers: what can this Agent use here, what is disabled/unconfigured/broken, and what concrete evidence explains that result.

#### Investigation helpers

Evidence-first helpers such as endpoint resolution, symbol/reference/write tracing, data lineage, symbol-scoped Git history, test-data metric explanation and debug-SQL reverse mapping are candidates after MCP/Skill core completion. They must return evidence/confidence/uncertainty rather than opaque guesses.

## Validated Requirements

### SKILL FOUNDATION

- **SKILL-RUN-01** ✅ — discovery is limited to explicit configured roots/definitions.
- **SKILL-RUN-02** ✅ — a Skill resolves to a real artifact/content source.
- **SKILL-RUN-03** ✅ — Environment selection gates access.
- **SKILL-RUN-04** ✅ — one global installation can serve multiple Environments.
- **SKILL-RUN-05** ✅ — an Agent using ADM can discover/read an enabled real Skill; disabled access is rejected.

These are foundation requirements, not the final Skill Runtime completion gate.

### VERIFIER

- **VERIFY-01..04** ✅ — verifier is optional, Environment-scoped, structured, bounded and does not become a global development prerequisite.

### MCP FOUNDATION

- **MCP-RUN-01** ✅ — definitions carry actual connection information.
- **MCP-RUN-02** ✅ — Environment selection gates runtime access.
- **MCP-RUN-03** ✅ — status distinguishes configured/disabled/healthy/error.
- **MCP-RUN-04** ✅ — real external tool discovery/call works through ADM.
- **MCP-RUN-05** ✅ — secrets resolve at activation boundaries and are not exposed normally.

These are foundation requirements, not the final MCP Runtime completion gate.

### RUNTIME / PROCESS / ISOLATION

- **LIFE-01..03** ✅ — one persistent local owner, desired-vs-observed separation and deterministic owned-resource cleanup.
- **PROC-01..03** ✅ on the current Windows dogfood target — stable dev process lifecycle, bounded logs and owned-port facts.
- **ISO-01..03** ✅ — optional managed worktree isolation without making Git an Environment prerequisite.
- **ARUN-01** ✅ retained — stable single-command asynchronous Runtime lifecycle. This requirement is task-semantic-neutral and remains ADM Core.

### Historical mis-scoped requirement

- **FLOW-01** — retired from the ADM product contract by the 2026-09-08 rebaseline. The implementation is historical/cleanup scope, not a capability future ADM features should depend on.
- **GSD-01..03 / PAR-01** — removed from the ADM roadmap. GSD state advancement and parallel Agent orchestration belong above ADM.

## Active Requirements

### CORE-CLEANUP

- **BOUNDARY-01** — remove ADM-owned Planner/Executor/Reviewer workflow surface while retaining generic asynchronous `run_` Runtime behavior.
- **BOUNDARY-02** — no GSD `.planning` interpretation/state-advance API is merged or introduced.
- **BOUNDARY-03** — removal must not regress command Run, verifier, MCP, Skill, process, file or Environment behavior.

### MCP COMPLETION

- **MCP-COMP-01** — one explicit MCP definition model represents supported transports/configuration without mixing desired config and observed health.
- **MCP-COMP-02** — supported MCP transports have real lifecycle implementations and local diagnostics; unsupported transport fails only that MCP.
- **MCP-COMP-03** — auth/secret configuration is explicit, resolved only at activation, and never leaked through status/errors/logs.
- **MCP-COMP-04** — tool inventory can be refreshed/inspected and Environment enable/disable revokes access immediately.
- **MCP-COMP-05** — connection/session/reconnect failures return actionable structured diagnostics and never become stale healthy state.

### SKILL COMPLETION

- **SKILL-COMP-01** — explicit refresh discovers current real Skill artifacts and reports source facts without arbitrary disk scanning.
- **SKILL-COMP-02** — Skill artifact/support-file availability and broken references are diagnosable per Skill without breaking unrelated Skills or Environment capabilities.
- **SKILL-COMP-03** — Environment selection remains the authorization boundary and availability status explains enabled/disabled/missing/broken states.
- **SKILL-COMP-04** — ADM exposes Skill content/support facts for the consuming Agent without inventing an ADM-specific Skill execution/orchestration engine.

### CAPABILITY DIAGNOSTICS

- **CAP-01** — one Environment capability inspection surface reports usable/unavailable capabilities with reasons and evidence.
- **CAP-02** — optional capability failures are operation-local and never silently become global Environment prerequisites.

## Deferred / Frozen

Until MCP/Skill/capability Core is complete, defer:

- Planner/Executor/Reviewer orchestration;
- GSD state automation;
- parallel Agent policy/parent aggregation;
- automatic Git merge/rebase/push;
- automatic Memory context composition;
- installer/MSI, tray, autostart, updater, signing, notifications;
- visual redesign;
- migration/compatibility work;
- broad Desktop feature expansion except blockers.

## Rebaseline Reference

See `.planning/rebaseline/2026-09-08-core-boundary.md` for the decision record and existing implementation audit.

---
*Last updated: 2026-09-08 after core-boundary rebaseline; Phase 10 GSD branch abandoned and orchestration cleanup is next.*
