# ADM V2 Development Status

This file is the project-level view of what is implemented, what is being worked on next, and what remains later.

`docs/PRODUCT_CONTRACT.md` defines product semantics. Phase documents describe delivered slices. This file tracks current implementation status and development order.

## Current milestone

First usable ADM V2: an Agent can develop software through MCP against an ordinary local directory without Git being a prerequisite.

The first usable ADM V2 milestone is accepted: the non-Git development loop passed through the Agent Gateway, and persisted context survived an actual temporary Gateway process restart and fresh MCP reconnect.

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

### Environment lifecycle management

- Environment rename by stable ID
- Rename changes only the human-readable name
- Root, Workspace association, MCP/Skill selections, private Memory, active Writer, and project files remain unchanged
- Environment list returns lightweight summaries with private Memory count but not values
- Environment inspect exposes Workspace relation, capabilities, resolved MCP/Skill entries, unresolved selection IDs, and private Memory count
- CLI and MCP share the same application-level Environment management view

Tracked by `docs/PHASE_04_ENVIRONMENT_MANAGEMENT_UX.md`.

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

### First milestone dogfood

- Completed a real non-Git development loop through the Agent Gateway
- Observed a failing test, made a follow-up code change, and reran successfully
- Verified Git failure remains operation-local
- Verified persisted context across an actual temporary Gateway process restart and fresh MCP reconnect

Tracked by `docs/DOGFOOD_01_NON_GIT_LOOP.md` and `TestHTTPGatewayPersistsContextAcrossProcessRestart`.

## In progress / next slice

### Management boundary

Phase 05 is active. Slice A is delivered: a read-only management snapshot now aggregates Workspaces, lightweight Environment summaries, exec allowlist, MCP catalog, Skill catalog, and Global Memory count through the existing application services without exposing Memory values.

Next Slice B adds explicit management mutation methods for already-defined product operations. They must delegate to existing services rather than duplicate validation, persistence, or product semantics.

## Next

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
