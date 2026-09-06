# ADM V2 Product Contract

Status: development baseline. Breaking changes are allowed until the user explicitly declares a compatibility boundary.

This file defines product semantics. Implementation plans and phases may refine how requirements are delivered, but must not silently change these requirements.

## Product goal

### ADM-GOAL-001 — Let an Agent develop software on this computer through MCP

ADM V2 is a local development MCP server. Its primary purpose is to let an Agent perform real software-development work against projects on this computer.

A usable ADM must let an Agent, through MCP, select/open a local project, inspect its files, search code, create/edit/delete files, execute permitted development commands, run tests/build/lint when available, inspect the results, and continue iterating until a real development task is complete.

Workspace, Environment, Memory, Skill, external MCP integration, Git integration, isolation, concurrency control, and other mechanisms exist only to support this goal. None of those abstractions may become an unrelated prerequisite that prevents ordinary local development.

The primary product test is therefore an end-to-end Agent development loop, not the existence of management abstractions.

## Core model

### ADM-CORE-001 — Workspace is a directory

A Workspace is a development directory explicitly added to ADM and identified by an ADM ID and a filesystem path.

Agents may operate only inside directories that have been added to ADM as Workspaces. ADM must not expose arbitrary filesystem access outside registered Workspace roots.

Acceptance:

- register an empty ordinary directory
- register a directory that is not a Git repository
- inspect/list it successfully

### ADM-CORE-002 — Environment is a development context

An Environment is a development context associated with one Workspace and rooted at a directory.

Minimum semantic fields:

- Environment ID
- Workspace ID
- human-readable name
- root path
- state/activity metadata
- writer lease state
- enabled global MCP IDs
- enabled global Skill IDs
- Environment-private memory

Environment is both a development root and a private AI development context. MCPs and Skills are defined once in the global catalog; an Environment stores only which global MCPs and Skills are enabled for that Environment. Environment-private memory belongs to that Environment and must not leak into other Environments.

Environment is not defined by Git branch, commit, worktree, Docker container, clone, or any other isolation mechanism.

Acceptance:

- create an Environment for a non-Git Workspace
- root defaults to the Workspace directory
- inspect/list survives daemon restart
- Environment MCP/Skill enable selections survive daemon restart
- enabling or disabling an MCP/Skill for one Environment does not affect another Environment or the global catalog
- Environment-private memory is not visible from another Environment by default

### ADM-CORE-003 — Tools are optional capabilities

Git, worktree, Docker, shell execution, verifier, MCP, language toolchains, and future integrations are optional capabilities.

Missing capability blocks only operations that require it.

Acceptance examples:

- no Git: read/search/write still work
- no verifier: Environment creation and file operations still work
- no shell execution: read/search still work
- no worktree: normal Environment development still works

### ADM-CORE-004 — Capability checks are operation-local

Each routed operation checks only the capability needed for that operation.

No workflow or Environment lifecycle may require unrelated capabilities as a global admission check.

### ADM-CORE-005 — Single writer per physical root

ADM must prevent concurrent writers from mutating the same physical development root through multiple Environments.

Read-only operations do not require the writer lease.

Mutation operations require the active writer owner.

A writer is a bounded lease, not a permanent lock. The lease has an explicit expiry time, successful owner activity renews it, and an expired lease must automatically stop blocking the physical root. A persisted writer without an explicit expiry is invalid and must be reacquired rather than inferred through compatibility rules. Force release remains a recovery fallback, not the normal lifecycle.

### ADM-CORE-006 — Core file development operations

The first usable V2 core must support, subject only to the relevant runtime policy:

- tree/list files
- read text files
- literal search
- write files
- exact edit
- delete a file

These operations must work in a plain non-Git directory.

### ADM-CORE-007 — Local command execution is allowlisted

Agents may execute local development commands only when the executable is explicitly allowed by ADM's command whitelist.

Command execution must run within the selected registered Workspace/Environment root and must not provide arbitrary command execution outside that scope.

Commands not permitted by the whitelist must be rejected clearly.

The allowlist is mutable management state. ADM must support removing a previously allowed executable, and that removal must affect subsequent capability checks and command execution without requiring an ADM restart. Removing a value that is not currently allowlisted must fail clearly rather than report a false success.

Acceptance:

- allow an executable and observe `shell.exec` when the Runtime otherwise supports execution
- remove the executable and observe the updated allowlist immediately
- after removing the last allowed executable, `shell.exec` is no longer advertised
- execution of the removed executable is rejected
- removing a non-allowlisted executable returns a clear error

### ADM-CORE-008 — Verifier is optional

Verifier execution is a Runtime capability.

A project without configured verifiers remains a valid development Workspace.

### ADM-CORE-009 — Git is a tool

Git features are optional tool operations.

Potential Git operations include status, diff, branch and later worktree support.

Git must not be part of the intrinsic Workspace or Environment data model.

### ADM-CORE-010 — Isolation is optional

Isolation is not part of the initial V2 Environment requirement.

Future isolation implementations may create/select another root directory and then create an Environment rooted there.

Possible future implementations include Git worktree, clone, copy, container, VM, or another mechanism. None is privileged in the core Environment model.

### ADM-CORE-011 — No pre-stable compatibility layer

During active V2 development, old development data may be discarded when a model changes.

Do not add migrations, legacy inference, compatibility fields, or fallback behavior unless the user explicitly introduces a compatibility requirement.

### ADM-CORE-012 — Global MCP/Skill catalog with Environment selection

All MCP and Skill definitions are managed in one global catalog.

Each global MCP and Skill definition includes a default-include-in-environment setting that controls whether newly created Environments enable that entry by default.

An Environment does not own duplicate/private MCP or Skill definitions. It stores only the IDs of global MCPs and Skills enabled for that Environment.

When a new Environment is created, its initial enabled MCP/Skill selections are copied from the current global default-include settings. Later changes to a global default must not silently rewrite existing Environment selections.

An Environment may then enable or disable any subset of the global catalog. Entries not selected for that Environment are not exposed in that Environment.

Changing an Environment's enabled selections must not modify the global MCP/Skill definitions, their global defaults, or another Environment's selections.

If a globally defined MCP or Skill is removed, affected Environment selections may become unresolved references; handling of unresolved selections must be explicit and must not silently recreate or copy definitions.

### ADM-CORE-013 — Global memory and Environment-private memory

ADM owns one global memory scope plus one private memory scope per Environment.

Global memory:

- persists independently from transient Agent sessions
- is available across Environments as shared durable context
- is not tied to any Workspace or Environment lifecycle
- is not deleted when an Environment is destroyed

Environment-private memory:

- persists independently from transient Agent sessions
- is available when working inside that Environment
- is isolated from other Environments by default
- is not promoted to global memory implicitly
- can be deleted with the Environment without deleting project files

When working inside an Environment, ADM may expose both Global Memory and that Environment's private Memory. Writes must target one scope explicitly; ADM must not silently copy, promote, merge, or synchronize memory between scopes.

### ADM-CORE-014 — Workspace lifecycle management is metadata-safe

ADM must support inspecting a Workspace by stable ID, changing its human-readable name, and removing its registration.

Renaming a Workspace changes only ADM metadata. It must not move or rename the Workspace directory.

Removing a Workspace removes only the ADM registration. It must not delete the Workspace directory or any project files. A Workspace that is still referenced by any Environment must not be removed; the Environment must be removed first.

Acceptance:

- inspect a registered Workspace by ID
- rename it while preserving its ID and filesystem path
- reject an empty replacement name
- reject removal while any Environment references the Workspace
- after Environment removal, remove the Workspace registration successfully
- Workspace removal leaves the directory and project files untouched

### ADM-CORE-015 — Human management surfaces expose persisted development context

ADM management state must not be usable only through the Agent-facing MCP Gateway. Human-facing management surfaces must expose the same persisted MCP catalog, Skill catalog, Memory scopes, and Environment selections needed to understand and administer the current development context.

The first human management surface is the CLI. It must expose management operations without changing the underlying product model or creating a second persistence path.

Acceptance:

- CLI can add/list/remove global MCP catalog entries and change their default-include setting
- CLI can add/list/remove global Skill catalog entries and change their default-include setting
- CLI can list/read/write/delete Global Memory
- CLI can enable/disable MCP and Skill IDs for one Environment without changing another Environment
- CLI can list/read/write/delete Environment-private Memory by explicit Environment ID
- CLI help makes these capabilities discoverable
- all operations reuse the same persisted state and application services used by MCP

### ADM-CORE-016 — Environment lifecycle management is metadata-safe

ADM must support changing an Environment's human-readable name by stable Environment ID.

Renaming an Environment changes only ADM metadata. It must not move or rename the Environment root directory, change the Workspace association, rewrite MCP/Skill selections, rewrite private memory, release an active writer, or touch project files.

Acceptance:

- rename an Environment by stable ID
- preserve Environment ID, Workspace ID, root path, creation time, MCP selections, Skill selections, private memory, and writer state
- reject an empty replacement name
- leave project directories and files untouched
- expose rename through both CLI and MCP management surfaces

## Agent-facing Gateway

### ADM-GW-001 — Gateway routes by Workspace/Environment identity

Agent tools operate against explicit stable IDs rather than assuming one current repository.

### ADM-GW-002 — Gateway must remain usable without Git

Gateway discovery, Environment lifecycle, read/search and permitted file mutation must work for plain directories.

Git-specific tools may return a clear unsupported-capability error when unavailable.

### ADM-GW-003 — Tool failure must be local

A failure or absence of one tool must not make unrelated Gateway tools unavailable.

### ADM-GW-004 — HTTP Gateway supports foreground and detached startup

`gateway start` defaults to foreground execution in the current terminal.

`gateway start -d` / `gateway start --detach` must start the same ADM V2 HTTP Gateway as a process detached from the current terminal, wait until its health endpoint reports ready, and then return control to the caller. The resulting process must remain manageable through the existing `gateway status` and `gateway stop` commands.

Acceptance:

- foreground `gateway start` remains blocking and Ctrl+C-stoppable
- `gateway start -d` and `gateway start --detach` return only after `/healthz` is ready
- closing the launching terminal does not stop the detached Gateway
- `gateway status` reports the detached Gateway PID/version
- `gateway stop` stops the detached Gateway
- detached startup does not introduce a separate daemon or generic process-manager prerequisite

## Development process contract

### ADM-DEV-001 — Requirement traceability

Every implementation phase must name the requirement IDs it changes or implements.

### ADM-DEV-002 — Explicit prerequisites

Every new prerequisite must be listed in the phase plan and requires user approval before implementation.

### ADM-DEV-003 — Negative acceptance tests

Tests must prove not only that features work, but that unrelated tools are not prerequisites.

### ADM-DEV-004 — Dogfood gate

Before expanding into speculative orchestration/isolation features, use V2 through its Agent Gateway to perform real development work on another project.

Real blockers discovered by dogfooding take priority over speculative features.

## Initial non-goals

The first V2 milestone intentionally does NOT include:

- Git worktree lifecycle
- clone-based Environment isolation
- branch management
- automatic Git synchronization
- parallel Agent workflow
- GSD automation engine
- Docker as a prerequisite
- UI / Windows Client
- migration from ADM V1 Environment data
- compatibility with V1 Environment store format
- generic process/log/port management

These may be added later only from concrete requirements.

## V2 first milestone acceptance gate

The first milestone is accepted only when an Agent can complete a real development loop through ADM MCP against an ordinary local project directory, including a directory with no `.git` directory.

The end-to-end gate is:

1. select/open a local project through ADM
2. inspect the project tree and read files
3. search for relevant code
4. create and edit files
5. delete a file when the task requires it
6. execute a permitted project command
7. run the project's tests/build/lint when those capabilities exist
8. inspect command results and make a follow-up code change
9. preserve the development context across an ADM restart
10. complete the task without Git being required

Supporting capabilities required by the same milestone include Workspace/Environment identity, global MCP/Skill definitions and Environment selections, global Memory, Environment-private Memory, and safe mutation/concurrency control. These supporting abstractions are not themselves the product acceptance goal; they are accepted only insofar as they support the Agent development loop above.
