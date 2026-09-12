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

### Post-1.0 priority update — 2026-09-11

The stable 1.0 line is green at `v1.0.1`. Post-1.0 work now prioritizes human management usability before deeper Agent automation. Desktop UX reorganization is the next phase because daily operation has become the most visible bottleneck, as long as it stays a management surface over the shared Core and does not introduce Desktop-only state, authorization or persistence.

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
- Long-lived ADM-owned capabilities may expose explicit health/recovery policies. Recovery is capability-specific, configurable and observable; it must never become blind replay of potentially mutating operations. MCP is the first implementation.
- Worktree is an optional isolation primitive; ADM does not decide parallel Agent policy or automatic integration.
- Important planning decisions live under `.planning/`, not only chat history.
- No capability is complete because CRUD/metadata exists; it requires a real consumption path and acceptance evidence.
- No pre-stable migration/compatibility burden unless explicitly requested.

## Current Delivery Status — 2026-09-11

- Phases 10-13 are complete: boundary cleanup, MCP runtime/import/health, Skill source/availability and capability diagnostics. The earlier foundation-gap audit below is historical.
- Phase 14 is complete through 14-02; further investigation helpers remain conditional. Phase 15 retention/cleanup is complete and must be reused by future workflow UX.
- Phase 16 production Desktop/normal CLI management use the shared Admin MCP client. Desktop may open while disconnected, but management cannot silently use writable local state as a fallback.
- Phases 16-17 are closed through the repository-recorded green v1.0.1 release.
- Phase 18 is complete. Desktop Management UX Reorganization is accepted through production-browser, full Go/vet, Wails build and visible native Wails/WebView2 evidence. See `.planning/phases/18-desktop-management-ux/18-CLOSEOUT.md`, `18-04-SUMMARY.md`, and `18-VALIDATION.md`.
- Phase 19 is complete. Bounded metadata-only Workspace discovery and Environment tree digest now flow through Core, Agent/Admin MCP, normal CLI and explicit Desktop navigation/root handoff. Automated B01-B08/B10 acceptance is green; B09 native GUI click-through is explicitly pending because no native GUI-control tool was available, while the uniquely named Wails artifact build passed. See `.planning/phases/19-workspace-discovery/19-CLOSEOUT.md` and `19-02-SUMMARY.md`.
- Transient navigation/filter/focus state is presentation state; existing connection/autostart preferences remain the prior explicit exception. Phases 18-19 create no new state file, schema or authority.

## Historical Reality Audit — 2026-09-08

This dated audit is retained for provenance. Its MCP/Skill/capability gaps were closed by Phases 11-13; it is not the current backlog.

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

### Historical superseded orchestration

Phase 9 introduced `run_workflow_start` and Planner / Executor / Reviewer workflow semantics. The 2026-09-08 rebaseline classified this orchestration layer as outside ADM's product responsibility, and Phase 10 removed the Agent-facing workflow tool plus workflow runtime/domain state from Core while retaining generic single-command `run_` lifecycle.

The GSD Phase Executor implementation remains only on the abandoned `feat/gsd-phase-executor` branch and was not merged into Core.

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

### CORE-CLEANUP

- **BOUNDARY-01** ✅ — ADM-owned Planner/Executor/Reviewer workflow surface is removed while generic asynchronous `run_` Runtime behavior remains.
- **BOUNDARY-02** ✅ — no GSD `.planning` interpretation/state-advance API was merged or introduced.
- **BOUNDARY-03** ✅ — removal passed generic Run, verifier, MCP, Skill, process, file/non-Git and managed-worktree regression gates.

## Completed Delivery Requirements — Phase 18

The active delivery is a refinement of ADM-DESKTOP-001, ADM-MGMT-001, ADM-CORE-015/017/020 and the existing operation-local capability, catalog and Memory boundaries.

- Menu shell, stable section routing and a read-only dashboard over the existing sanitized management snapshot.
- Current connection/Environment context remains clear across pages; disconnected, empty, partial-failure and stale data are distinguished.
- Existing feature actions, modal behavior, Runtime authority and explicit Memory reads remain intact.
- No new prerequisites, Core APIs, persistence, frontend build pipeline or task orchestration.
- 18-01 through 18-04 are complete and accepted. Phase 18 closeout records final D01-D10 evidence and an explicit D10 artifact-path note.

## Completed Delivery — Phase 19

Context and detailed plans: `.planning/phases/19-workspace-discovery/19-CONTEXT.md`, `19-01-PLAN.md` and `19-02-PLAN.md`. Phase 19 is complete through bounded metadata-only Workspace discovery, Environment tree digest, Agent/Admin MCP exposure, normal CLI Admin-MCP usage and explicit Desktop navigation/root handoff. Full Go/vet and production-browser acceptance passed on final implementation commit `8a5c99a`; uniquely named CLI and Wails artifacts were built and hashed. Native Wails/WebView2 click-through remains explicitly pending under B09 because this session had no native GUI-control tool. See `19-01-SUMMARY.md`, `19-02-SUMMARY.md` and `19-CLOSEOUT.md`. Requirements: ADM-GOAL-001, ADM-CORE-001/002/003/004/005/006, ADM-MGMT-001 and ADM-DESKTOP-001. No new prerequisites, automatic scans, persistence or task orchestration were introduced. Global MCP probes and Skill structural availability remain independent of Environment selection.

## Active Planning — Phase 20

Detailed planning is open under `.planning/phases/20-agent-context-bundle/20-CONTEXT.md`, `20-01-PLAN.md` and `20-02-PLAN.md`. The contract composes existing Environment authority, bounded Phase-19 tree digest, Phase-13 capability facts and passive Gateway-owner observations into one explicit read-only Environment context bundle. ADM keeps stable `environment_id` routing and does not add hidden session selection state. 20-01 Core is complete: the shared bounded bundle, passive no-Git capability projection, safe MCP/Skill/verifier summaries, factual guidance and deterministic omission/output caps are validated by full Go/vet evidence. MCP initialize instructions remain 20-02 scope and may contain only static non-project-specific usage guidance; dynamic context will require an explicit bundle call. No MCP probing/execution, Skill content injection, Memory value injection, persistence or task orchestration occurs while building the bundle. 20-02 is the next executable plan and has not started.

## Core Delivery Requirements — Completed in Phases 11-13

These requirement definitions remain the contract for regression and traceability. Current delivery evidence is indexed in STATE.

### MCP COMPLETION

- **MCP-COMP-01** — one explicit MCP definition model represents supported transports/configuration without mixing desired config and observed health.
- **MCP-COMP-02** — supported MCP transports have real lifecycle implementations and local diagnostics; unsupported transport fails only that MCP.
- **MCP-COMP-03** — auth/secret configuration is explicit, resolved only at activation, and never leaked through status/errors/logs.
- **MCP-COMP-04** — tool inventory can be refreshed/inspected and Environment enable/disable revokes access immediately.
- **MCP-COMP-05** — connection/session/reconnect failures return actionable structured diagnostics and never become stale healthy state.
- **MCP-COMP-06** — enabled MCPs have an explicit configurable health/recovery policy: periodic protocol health checks, bounded probe timeout, user-configurable check interval, automatic reconnect defaulting to off, and fixed configurable reconnect interval when enabled. Background recovery is owner-local, has no exponential/adaptive backoff in Phase 11, and never replays a failed tool call.
- **MCP-COMP-07** — ADM can preview and atomically import one or many MCP definitions from supported external JSON/JSONC formats (OpenCode, WorkBuddy/CodeBuddy, Codex plugin MCP JSON, Claude Code, and MCPHub) into one canonical MCPDefinition model. Name conflicts default to error; Phase 11 import writes global MCP definitions only, never changes existing Environment selections, and converts literal credential-bearing values into secret/environment-reference requirements instead of persisting literals.

### SKILL COMPLETION

- **SKILL-COMP-01** — explicit refresh discovers current real Skill artifacts and reports source facts without arbitrary disk scanning.
- **SKILL-COMP-02** — Skill artifact/support-file availability and broken references are diagnosable per Skill without breaking unrelated Skills or Environment capabilities.
- **SKILL-COMP-03** — Environment selection remains the authorization boundary and availability status explains enabled/disabled/missing/broken states.
- **SKILL-COMP-04** — ADM exposes Skill content/support facts for the consuming Agent without inventing an ADM-specific Skill execution/orchestration engine.

### CAPABILITY DIAGNOSTICS

- **CAP-01** — one Environment capability inspection surface reports usable/unavailable capabilities with reasons and evidence.
- **CAP-02** — optional capability failures are operation-local and never silently become global Environment prerequisites.

## Deferred / Frozen

Keep the following deferred/frozen unless explicitly opened through the current roadmap:

- Planner/Executor/Reviewer orchestration;
- GSD state automation;
- parallel Agent policy/parent aggregation;
- automatic Git merge/rebase/push;
- automatic Memory context composition;
- installer/MSI, updater, signing, notifications unless post-1.0 dogfood proves need;
- broad visual redesign outside the approved Phase 18 Desktop management UX reorganization;
- migration/compatibility work;
- broad Desktop feature expansion except blockers.

## Rebaseline Reference

See `.planning/rebaseline/2026-09-08-core-boundary.md` for the decision record and existing implementation audit.

---
*Last updated: 2026-09-12 for Phase 20-01 Core completion. The 2026-09-08 audit is historical; Phases 18-19 are complete, Phase 20-01 is complete, and 20-02 is ready but not started.*
