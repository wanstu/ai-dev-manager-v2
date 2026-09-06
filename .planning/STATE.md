# AI Dev Manager V2 — Planning State

## Current Position

Milestone: R1 — Agent-ready Development Context
Phase: 01 — Real Skill Runtime + GSD Bootstrap
Plan: 01-01 — Real Skill Vertical Slice
Status: Ready

## Current Objective

Make one real globally installed GSD Skill discoverable and consumable through ADM for an enabled Environment, with disabled-Environment isolation proven end-to-end.

## Completed Planning Phases

- Phase 00 — Rebaseline and Planning Authority: completed 2026-09-06. Established PROJECT/ROADMAP/STATE, corrected Skill/MCP completion claims, froze peripheral Desktop/package work, and created the Phase 01 executable plan.

## What Triggered the Reset

Real Desktop dogfood exposed that several capabilities described as complete were only management/configuration surfaces:

- Skill was catalog CRUD + Environment IDs + an instructions field, not a real shared Skill runtime.
- MCP was initially catalog CRUD + Environment IDs; a minimal real HTTP proxy was later added, but lifecycle/health/runtime semantics remain partial.
- Desktop work was consuming attention while verifier, persistent process lifecycle, isolation, Agent Run and GSD runtime remained absent.
- Ad-hoc reaction to the latest bug replaced Phase planning and caused roadmap drift.

## Current Truth

### Validated

Workspace/Environment basics, file development tools, writer lease, allowlisted exec, optional basic Git, persistent state, Gateway transport/lifecycle, scoped Memory CRUD, CLI/Desktop management shell.

### Partial

Skill runtime, external MCP runtime, Desktop completeness.

### Absent

Verifier, persistent Runtime ownership, process/log/port lifecycle, worktree isolation, Agent Run, Planner/Executor/Reviewer, V2 GSD executor, parallel orchestration.

See `.planning/PROJECT.md` for the detailed audit.

## Decisions Recorded During Reset

1. The earlier ADM implementation contains a real GSD planning system and validated engineering sequence. It is evidence/reference, not code automatically copied into V2.
2. V2 keeps its Git-independent Environment semantics; optional worktree isolation comes later.
3. Skill runtime is the first Core gap because it is required to truthfully consume the global GSD Skill through ADM.
4. Verifier comes before Agent/GSD orchestration.
5. Persistent ownership comes before long-running processes and Agent Runs.
6. Desktop is frozen except for blockers until Core parity milestone.
7. Packaging/polish is frozen.
8. Do not call Phase 01 “using GSD through ADM” until the real GSD Skill acceptance passes.

## Immediate Next

Execute only Phase 01 Plan 01-01:

- inspect the actual GSD Skill installation on this machine;
- define the minimum real Skill artifact/source model;
- implement Environment-gated Skill discovery/read through ADM;
- prove enabled/disabled isolation with the real GSD Skill;
- use the Skill through ADM before advancing to Phase 02.

## Backlog — Do Not Interrupt Current Phase

- automatic Memory context composition;
- additional MCP transports beyond the first real need;
- generic Docker capability;
- debugger/DAP;
- LAN/remote Gateway auth/exposure;
- installer/tray/autostart/updater/signing/notifications;
- UI redesign;
- migration from prior ADM implementations.

## Bootstrap GSD Status

The repository now uses GSD-compatible planning artifacts as a manual bootstrap, informed by the real earlier ADM `.planning` system.

However, V2 cannot yet claim that this Agent consumed the GSD Skill through V2. That becomes true only after Phase 01 acceptance.