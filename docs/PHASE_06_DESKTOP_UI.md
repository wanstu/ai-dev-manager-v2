# Phase 06 — Desktop Management UI

## Requirement IDs

- ADM-DESKTOP-001
- ADM-MGMT-001
- ADM-DEV-001
- ADM-DEV-003

## Goal

Add a cross-platform desktop management shell using Wails while keeping ADM product semantics in the existing application/management layers.

The desktop app is a human management surface. It does not replace the Agent Gateway and must not become a prerequisite for normal ADM development.

## Delivery slices

### Slice A — Desktop shell and snapshot

Status: delivered.

- add a separate `ai-dev-manager-v2-desktop` entrypoint
- use Wails v2 and bind the existing `desktop.Adapter`
- reuse `store.DefaultPath()` and the existing app/management services
- embed plain HTML/CSS/JavaScript assets directly in the Go binary
- render counts and lightweight lists from `GetSnapshot`
- provide an explicit Refresh action
- show loading and error states
- no npm/Vite/Node build prerequisite
- no REST or running Gateway prerequisite

### Slice B — Workspace / Environment management

Status: delivered.

- add Workspace registration form
- add Workspace rename/remove actions through the desktop adapter
- add Environment create form with explicit Workspace selection and optional root
- add Environment rename/remove actions through the desktop adapter
- add Environment detail panel using `InspectEnvironment`
- surface backend lifecycle/safety errors instead of duplicating those rules in JavaScript
- make destructive action copy explicit that ADM records are removed without deleting project directories/files

### Slice C — Catalog / allowlist / Memory management

Add exec allowlist, MCP/Skill catalog and Environment selection management, plus explicit Global and Environment-private Memory panels. Memory values remain behind explicit Memory views.

### Slice D — Gateway lifecycle and desktop polish

Add Gateway status/lifecycle controls and platform-specific desktop polish only after the core management UI is usable.

## Acceptance tests for Slice A

- `go test ./...` still passes
- desktop command compiles on the current Windows development machine
- desktop command uses the normal ADM state path
- bound object is `desktop.Adapter`
- frontend calls `GetSnapshot` through Wails binding
- snapshot page renders Workspace / Environment / allowlist / MCP / Skill / Global Memory counts
- refresh reloads snapshot state
- frontend assets are embedded and require no Node tooling
- CLI/Gateway command remains unchanged

## Acceptance tests for Slice B

- desktop command still compiles on the current Windows development machine
- frontend binds Workspace add/rename/remove through `desktop.Adapter`
- frontend binds Environment create/inspect/rename/remove through `desktop.Adapter`
- Workspace removal safety remains enforced by the backend management boundary
- Environment removal does not delete project files and active Writer protection remains backend-enforced
- Environment detail comes from the shared sanitized inspect view and does not dump private Memory values
- `go test ./...`, `go vet ./...`, and `git diff --check` pass

## Explicit non-goals for Slice A

- Workspace/Environment mutation forms
- Memory value editor
- Gateway controls
- tray integration
- autostart
- single-instance handling
- notifications
- REST API
- Git/worktree lifecycle
- isolation/orchestration

## New prerequisites

- Wails v2 runtime/library for the separate desktop executable
- OS WebView runtime required by Wails on the target platform

These prerequisites apply only to the desktop executable. They are not prerequisites for CLI, MCP Gateway, Workspace, Environment, or normal Agent development.
