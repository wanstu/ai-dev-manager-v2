# Phase 05 — Management boundary

## Requirement IDs

- ADM-MGMT-001
- ADM-CORE-015
- ADM-CORE-017
- ADM-DEV-001
- ADM-DEV-003

## Goal

Create one application-backed boundary for desktop/human management clients so future Wails or HTTP adapters do not read `state.json` directly and do not recreate ADM product semantics.

## Delivery slices

### Slice A — Read-only management snapshot

Status: delivered.

- aggregate registered Workspaces
- aggregate lightweight Environment summaries
- include exec allowlist
- include MCP catalog
- include Skill catalog
- include Global Memory entry count, not values
- reuse current application services and persisted store
- no Git/worktree/Gateway prerequisite

### Slice B — Explicit management commands

Add mutation methods only for concrete desktop-management operations already defined by the product contract. These methods must delegate to existing application services rather than duplicate validation or persistence.

### Slice C — Desktop adapter

Expose the management boundary to the desktop client through the smallest adapter required by the chosen Wails integration. Do not introduce REST unless a concrete client needs an HTTP boundary.

## Acceptance tests

- snapshot reflects Workspace and Environment state
- Environment private Memory values are absent
- Global Memory values are absent
- catalog and allowlist state are present
- persisted changes appear in a later snapshot
- snapshot works for a non-Git Workspace
- no second state/cache file is introduced

## Explicit non-goals

- UI implementation in Slice A
- REST API without a client requirement
- Git worktree lifecycle
- isolation/orchestration
- direct state.json access from UI code
- V1 compatibility or migration

## New prerequisites

None.
