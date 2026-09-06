# ADM V2 Development Status

This file is the project-level view of what is implemented, what is being worked on next, and what remains later.

`docs/PRODUCT_CONTRACT.md` defines product semantics. Phase documents describe delivered slices. This file tracks current implementation status and development order.

## Current milestone

First usable ADM V2: an Agent can develop software through MCP against an ordinary local directory without Git being a prerequisite.

The core implementation and management CLI parity exist. The remaining acceptance gate for this milestone is real dogfood through the V2 Gateway; concrete blockers found there take priority over further UX or abstraction work.

## Completed

### Git-independent development core

- Workspace registration for ordinary directories, including non-Git directories
- Environment creation rooted inside a Workspace
- tree / read / literal search
- write / exact edit / file delete
- single-writer lease per physical root
- explicitly allowlisted local command execution
- optional Git status / diff / branch capabilities
- persistent local state
- stdio and HTTP MCP Gateway

Tracked by `docs/PHASE_01_GIT_INDEPENDENT_CORE.md`.

### Development context

- Global MCP catalog: add / list / remove / set default
- Global Skill catalog: add / list / remove / set default
- Per-Environment MCP selections
- Per-Environment Skill selections
- Global Memory: list / read / write / delete
- Environment-private Memory: list / read / write / delete
- `environment_inspect`
- Gateway tools for the management capabilities above

Tracked by `docs/PHASE_02_DEVELOPMENT_CONTEXT.md`.

### Gateway lifecycle

- Foreground `gateway start`
- Detached `gateway start -d` / `--detach`
- `gateway status`
- `gateway stop`
- `gateway restart`

### Workspace lifecycle

- Workspace inspect by stable ID
- Workspace rename without moving the directory
- Workspace remove without deleting project files
- Removal is blocked while an Environment still references the Workspace

### Exec allowlist lifecycle

- Add allowed executable
- List allowed executables
- Remove allowed executable
- Capability and execution permission update immediately after removal

### Management CLI parity — delivered slices

- Global MCP catalog CLI: add / list / remove / set default
- Global Skill catalog CLI: add / list / remove / set default
- Global Memory CLI: list / read / write / delete
- Per-Environment MCP selection CLI: enable / disable
- Per-Environment Skill selection CLI: enable / disable
- Environment-private Memory CLI: list / read / write / delete by explicit Environment ID
- CLI help exposes MCP, Skill, both Memory scopes, and Environment selection management

Tracked by `docs/PHASE_03_MANAGEMENT_CLI.md`.

## In progress / next gate

### Dogfood the first milestone

Phase 03 management CLI parity is complete. Dogfood 01 exercised a dedicated non-Git Go project through the Agent Gateway and passed the core development loop.

Passed in Dogfood 01:

1. select/open a non-Git project
2. inspect tree and files
3. search code
4. create/edit/delete files
5. execute an allowlisted development command
6. observe an intentional test failure
7. make a follow-up code change from that failure
8. rerun tests successfully
9. verify Git failure stays local and unrelated file access still works

Still pending before the milestone gate is fully closed:

- restart the Gateway process, reconnect through MCP, and verify the same persisted Workspace/Environment/context remains available

See `docs/DOGFOOD_01_NON_GIT_LOOP.md`.

## Next

### Environment management UX

After CLI parity, review Environment lifecycle and management ergonomics as one coherent surface rather than adding isolated commands opportunistically.

Candidates to decide and implement from concrete UX needs include:

- Environment rename
- clearer inspect output for MCP / Skill / Memory state
- easier discovery of Workspace-to-Environment relationships
- consistent CLI naming and error messages across management commands

### Management boundary

Define a stable management service/API boundary for desktop management clients instead of coupling UI code directly to persistence internals.

This should reuse the same application services used by CLI and MCP rather than create a second product model.

### Desktop Management UI

Preferred direction: Wails desktop application, following the same general shape as CodexPro+.

Expected management areas:

- Gateway status/lifecycle
- Workspaces
- Environments and Writer state
- Exec allowlist
- MCP catalog and Environment selections
- Skill catalog and Environment selections
- Global and Environment-private Memory

Wails keeps a Go backend while allowing an HTML/CSS/JS UI and leaves room for Windows/macOS/Linux support. OS-specific integrations should remain isolated behind platform-specific files.

## Later / not part of the immediate slice

These are not current prerequisites for normal ADM V2 development:

- Git worktree lifecycle
- clone/copy/container/VM isolation
- verifier configuration
- task/session attribution
- process/log/port management
- parallel Agent orchestration
- GSD workflow automation
- generic isolation/orchestration framework
- V1 state migration or compatibility

They should be pulled forward only by a concrete product or dogfood requirement.

## Known development-process incident

A previous set of feature worktrees was accidentally based on prior feature commits instead of the intended integration branch, creating unintended stacked branches and cherry-pick conflicts.

The incident is also recorded in ADM Global Memory under `development.process.worktree-base-incident`.

For independent feature work, the worktree base must be checked against the intended integration branch before development starts. Stacked branches should exist only when the dependency is intentional.

## Status maintenance

When a feature changes from planned to implemented, update this file in the same development slice so project status does not depend on chat history or memory.
