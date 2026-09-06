# Phase 04 — Environment management UX

## Requirement IDs

This phase implements or refines:

- ADM-CORE-002
- ADM-CORE-016
- ADM-DEV-001
- ADM-DEV-003

## Goal

Make Environment lifecycle management as explicit and safe as Workspace lifecycle management without changing the Environment model or introducing isolation semantics.

## Delivery slices

### Slice A — Environment rename

Status: delivered.

- rename an Environment by stable ID
- rename changes only the human-readable name
- preserve Workspace ID, root, selections, private memory, writer state, and creation time
- reject an empty replacement name
- do not move, rename, or delete project files
- expose the operation through CLI and MCP

### Slice B — Environment inspection ergonomics

Review the existing `environment list` / `environment inspect` output after Slice A and improve only concrete discoverability gaps. Do not add a second management model.

## Acceptance tests

- rename preserves stable identity and root
- rename preserves MCP/Skill selections and private memory
- rename preserves an active writer lease
- empty names fail clearly
- project files are untouched
- CLI help exposes rename
- MCP discovery exposes `environment_rename`
- existing non-Git development tests continue to pass

## Explicit non-goals

- Git worktree lifecycle
- branch management
- isolation
- task/session attribution
- orchestration
- desktop UI
- V1 compatibility or migration

## New prerequisites

None.
