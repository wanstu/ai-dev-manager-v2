# ai-dev-manager-v2

ADM V2 is a clean redevelopment of AI Development Manager. V1 is a reference implementation and a source of proven components, but V2 does not preserve V1's assumption that a development directory must first qualify as a Git/worktree-capable repository.

## Core rule

**A Workspace is a directory.**

A directory does not need Git, GSD, worktree, Docker, a verifier, or a language toolchain in order to be registered and developed through ADM.

Those integrations are tools. Missing tools block only the operations that need them.

```text
registered local directory
        ↓
     Workspace
        ↓
    Environment
        ↓
      Runtime
        ↓
operation-local capabilities
```

Product semantics are defined in `docs/PRODUCT_CONTRACT.md`. Development rules are defined in `AGENTS.md`.

## Implemented vertical slice

The current V2 code supports:

- registering an existing ordinary local directory as a Workspace
- creating a persistent Environment whose root defaults to the Workspace directory
- creating Environments for directories with no `.git`
- persistent Environment writer leases
- single writer protection per physical root
- `tree`
- `read`
- literal `search`
- `write`
- exact `edit`
- safe single-file `delete`
- explicitly allowlisted local command execution
- optional `git_status`, `git_diff`, and `git_branch`
- stdio MCP Gateway routing by explicit `environment_id`
- persistent state across process restarts
- global MCP catalog with per-Environment enable selections
- global Skill catalog with per-Environment enable selections
- global durable Memory
- Environment-private Memory
- Gateway bootstrap tools for Workspace registration and executable allowlisting

Git capabilities are detected per Environment root. On a non-Git directory, Git operations fail locally as unsupported while file development continues normally.

GSD and Git worktree lifecycle are intentionally not part of this first V2 core.

## CLI

Build:

```powershell
go build -o ai-dev-manager-v2.exe ./cmd/ai-dev-manager
```

Register any existing directory, including an empty or non-Git directory:

```powershell
.\ai-dev-manager-v2.exe workspace add D:\projects\some-directory
.\ai-dev-manager-v2.exe workspace list
```

Create a development Environment:

```powershell
.\ai-dev-manager-v2.exe env create --workspace <ws_id> --name task-a
.\ai-dev-manager-v2.exe env inspect <env_id>
```

Acquire one writer for the physical root:

```powershell
.\ai-dev-manager-v2.exe env writer acquire --owner session-a <env_id>
```

Allow a development executable explicitly:

```powershell
.\ai-dev-manager-v2.exe exec allow go
.\ai-dev-manager-v2.exe exec list
```

Start the Agent-facing MCP server over stdio:

```powershell
.\ai-dev-manager-v2.exe gateway stdio
```

Or expose the same Gateway as Streamable HTTP on a loopback address. The MCP endpoint is `/mcp`:

```powershell
.\ai-dev-manager-v2.exe gateway http --listen 127.0.0.1:41137
# endpoint: http://127.0.0.1:41137/mcp
```

V2 intentionally restricts this HTTP listener to loopback during the bootstrap phase. Docker/non-loopback exposure can be added later as an explicit capability instead of widening the local security boundary by default.

The state location can be inspected with:

```powershell
.\ai-dev-manager-v2.exe state path
```

Set `ADM_V2_HOME` to use an alternate state directory during development or tests.

## MCP tools in the current slice

Discovery/lifecycle:

- `workspace_list`
- `workspace_add`
- `environment_list`
- `environment_create`
- `environment_inspect`
- `environment_writer_acquire`
- `environment_writer_release`
- `exec_allow`
- `exec_allow_list`

Development context:

- `mcp_list`, `mcp_add`, `mcp_remove`, `mcp_set_default`
- `environment_mcp_set`
- `skill_list`, `skill_add`, `skill_remove`, `skill_set_default`
- `environment_skill_set`
- `memory_global_list`, `memory_global_read`, `memory_global_write`, `memory_global_delete`
- `memory_environment_list`, `memory_environment_read`, `memory_environment_write`, `memory_environment_delete`

Development:

- `tree`
- `read`
- `search`
- `write`
- `edit`
- `delete`
- `exec`

Optional Git tools:

- `git_status`
- `git_diff`
- `git_branch`

Mutation tools require a matching `writer_owner`. Read-only operations do not require a writer.

## Verified behavior

`go test ./...` includes negative acceptance tests proving that:

- a directory with no `.git` can be registered
- an Environment can be created for it
- file development works through the core service
- the same development loop works through an in-memory MCP client/server session
- Git capability is absent for the plain directory
- calling the Git tool fails only that tool
- command execution is absent until an executable is explicitly allowlisted
- two Environments sharing one physical root cannot hold two writers concurrently
- Environment/writer state survives a service restart

## Current non-goals

The current slice deliberately does not implement:

- Git worktree lifecycle
- branch-per-task Environment semantics
- GSD workflow automation
- automatic isolation
- parallel Agent orchestration
- Docker requirements
- V1 state compatibility/migrations

These can be added later as optional capabilities after concrete requirements or dogfood blockers justify them.
