# Post-1.0 Phase Map

Date: 2026-09-11
Base: `c898586` (`v1.0.1` green release hotfix)

This document sets the high-level post-1.0 phase sequence. It intentionally does not define execution-level `PLAN.md` files for every phase. Each phase gets a detailed plan only when work on that phase starts.

## Post-1.0 planning rules

1. User-facing management UX is the first priority after the green 1.0 line.
2. Desktop remains a management surface over the shared Core; it must not introduce Desktop-only state, authorization or persistence.
3. CLI remains a management/runtime surface over the same Core; CLI and Desktop can be separated logically before any repository split.
4. ADM supplies capabilities, lifecycle, authorization and diagnostics. It must not become a task planner/orchestrator.
5. Long-running work should prefer async start/status/cancel patterns over blocking request/response calls.
6. Large workspace navigation must be bounded and explainable; ADM should help Agents map human language to likely project roots without scanning arbitrary host paths.
7. Temporary resources must remain explicit, inspectable and conservative to clean up.
8. Distribution polish remains conditional unless daily use proves it is blocking.

## Phase sequence

| Phase | Name | Priority | Status | Purpose |
|---|---|---:|---|---|
| 18 | Desktop Management UX Reorganization | P1 | 18-01/02/03/04 detailed | Make ADM easier for the human operator to manage through a menu-based Desktop shell and validated section workflows. |
| 19 | Workspace Discovery + Project Navigation | P1 | Planned | Let ADM summarize large workspace directories and suggest likely project roots such as `projects/p2`. |
| 20 | Agent Context Bundle + Capability Injection | P1 | Planned | Give Agents a compact Environment context bundle: root, tree digest, enabled MCP/Skill summary, verifier/run guidance and capability reasons. |
| 21 | Async Verifier + Long Operation Observability | P1 | Planned | Move heavy verifier/test workflows toward async observable lifecycle instead of long blocking calls. |
| 22 | Temporary Task Environments + Safe Cleanup Workflow | P1/P2 | Planned | Make task-scoped temporary Environments first-class and safely cleanable/promotable. |
| 23 | CLI Agent UX + MCP/Skill Provisioning | P2 | Planned | Improve CLI flows for MCP/Skill import, enablement, refresh, diagnostics and Agent-friendly setup. |
| 24 | Desktop/CLI Surface Boundary Split | P2/P3 | Planned | Separate Desktop and CLI surfaces logically while preserving one shared Core/state model. |
| 25 | Distribution Polish If Needed | P3 conditional | Standby | Installer, updater, signing, notifications and deeper packaging polish only after concrete dogfood need. |
| 26 | Evidence-first Investigation Expansion | Conditional | Standby | Add deeper code/debug/data evidence helpers only when generic Runtime/search is insufficient. |

## Phase 18 — Desktop Management UX Reorganization

**Goal:** Reorganize Desktop from a dense management panel into a navigable application.

**Scope:** left/top menu shell, dashboard, section routing, clearer current ADM/Gateway status, and first pass at separating Workspace, Environment, Runs, MCP, Skill, Memory, Exec Allowlist, Gateway, Diagnostics and Settings views.

**Non-goals:** new Core model, new persistence, new Desktop-only authorization, installer/updater/signing, major backend expansion.

**Completion direction:** the same management capabilities become easier to find and operate without reducing existing functionality.

**Active planning:** `.planning/phases/18-desktop-management-ux/18-CONTEXT.md` records the scope review; `18-UI-REFACTOR.md` defines the source-grounded UI design; `18-01-PLAN.md` through `18-04-PLAN.md` define the sequential delivery and mandatory integrated closeout; `18-VALIDATION.md` is the shared acceptance index. Implementation has not started.

**Requirements:** ADM-DESKTOP-001, ADM-MGMT-001, ADM-CORE-015/017/020, preserving optional-capability, writer, catalog/Memory and Runtime boundaries. **New prerequisites:** none.

**Delivery sequence:** 18-01 establishes the compact connection/context header, grouped menu, ten real routes, truthful dashboard and shared scope/load-state contract. 18-02 refines Workspace/Environment and existing Runtime views. 18-03 refines MCP/Skill, explicit Global/Environment Memory, existing system controls and diagnostics. 18-04 is mandatory integrated acceptance across final routes/context/focus/scaling and the exact Wails artifact; only evidence-backed regressions are fixed there. Distinguish unloaded from zero, isolate auxiliary read failures, preserve dialogs and discard old-scope responses throughout.

## Phase 19 — Workspace Discovery + Project Navigation

**Goal:** Make large workspace roots understandable when a human says things like “look at p2”.

**Scope:** bounded workspace scan, project candidate summaries, git/module/package markers, suggested Environment root, tree digest and Desktop/CLI/Agent-facing views.

**Non-goals:** arbitrary full-disk indexing, semantic search engine, automatic code editing, hidden background scans.

**Completion direction:** ADM can summarize `projects/p1`, `projects/p2`, `projects/p3` style directories and help select/create the right Environment quickly.

## Phase 20 — Agent Context Bundle + Capability Injection

**Goal:** Reduce Agent setup friction by providing one compact context object for the selected Environment.

**Scope:** Environment root, bounded tree digest, enabled MCPs and tool inventory summary, enabled Skills and support summary, verifier availability, run/verifier usage guidance, unavailable capability reasons.

**Non-goals:** prompt orchestration, choosing the task plan, auto-calling MCP tools, interpreting Skills into ADM workflows.

**Completion direction:** an Agent can start work with one concise ADM context instead of manually probing many surfaces.

## Phase 21 — Async Verifier + Long Operation Observability

**Goal:** Make heavy verification flows observable and resilient across client/tool timeouts.

**Scope:** async verifier start/list/status/cancel, bounded live output, terminal result retention, clearer sync diagnostic when a request is too large/long for blocking execution.

**Non-goals:** replacing generic `run_`, CI orchestration, automatic retry of failed mutating commands.

**Completion direction:** long tests/builds/verifiers can be monitored without losing the result to an outer client timeout.

## Phase 22 — Temporary Task Environments + Safe Cleanup Workflow

**Goal:** Let Agents create task-scoped work areas without polluting durable project state.

**Scope:** temporary Environment creation, TTL/owner/run attachment, cleanup preview/execute, promote-to-durable, dirty/unpublished/active-writer/run blockers, Desktop visibility.

**Non-goals:** unsafe deletion of host files, automatic merge/push, deleting ambiguous non-ADM roots.

**Completion direction:** temporary work can be safely created, inspected, promoted or cleaned up with conservative evidence.

## Phase 23 — CLI Agent UX + MCP/Skill Provisioning

**Goal:** Make CLI setup and automation smoother for humans and Agents.

**Scope:** clearer MCP/Skill import/preview/apply flows, enable/disable commands, refresh/status/diagnostics summaries, Environment capability/context CLI output, script-friendly JSON modes where missing.

**Non-goals:** a second Admin API, CLI-only model, bypassing Environment selection.

**Completion direction:** common MCP/Skill/Environment setup tasks are easy from CLI and produce Agent-usable output.

## Phase 24 — Desktop/CLI Surface Boundary Split

**Goal:** Keep Desktop and CLI independently optimizable without forking product semantics.

**Scope:** package/module cleanup, shared management adapter boundaries, release/build naming clarity, surface-specific tests, frontend/backend separation inside the existing repository.

**Non-goals:** immediate multi-repo split, second state file, duplicate business logic.

**Completion direction:** Desktop and CLI can evolve separately while still sharing one Core contract.

## Phase 25 — Distribution Polish If Needed

**Goal:** Add distribution features only when daily use proves they matter.

**Scope candidates:** installer, updater, signing, notifications, deeper autostart polish, release UX.

**Non-goals:** speculative packaging work before a real blocker.

**Completion direction:** open only with concrete post-1.0 dogfood evidence.

## Phase 26 — Evidence-first Investigation Expansion

**Goal:** Add deeper investigation helpers only when real development tasks show generic file/search/runtime tools are insufficient.

**Scope candidates:** symbol/reference/write tracing, impact evidence, data lineage, debug-SQL-to-code mapping, test metric explanation, symbol-scoped history.

**Non-goals:** mandatory code graph engine, speculative provider lock-in, ADM task planning.

**Completion direction:** each slice must return evidence, confidence and uncertainty, with fallback behavior when optional providers are unavailable.

## Review notes — 2026-09-11

The sequence is retained after source/contract review. Phase 19's bounded tree digest informs 20. Async verifier (21) has no hard dependency on 19/20 and may be explicitly reprioritized if the existing pjadm timeout blocker recurs. Phase 22 reuses the completed Phase 15 lifecycle/cleanup primitives; task-scoped UX must not create ADM task orchestration. Local frontend extraction required by 18 does not wait for the wider Desktop/CLI boundary cleanup in 24. Phases 25/26 remain conditional.

## Immediate next phase

Phase 18 planning is ready. Execute `.planning/phases/18-desktop-management-ux/18-01-PLAN.md` Task 1 after Git/context checks and writer acquisition. Then proceed only through the recorded gates to 18-02, 18-03 and mandatory 18-04. The detailed successors do not authorize premature implementation before predecessor acceptance. Phases 19-26 remain phase-level.
