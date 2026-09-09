# Plan 11-03 Progress — Import Preview Foundation

Date: 2026-09-09

This is an initial implementation slice for Plan 11-03. It establishes the source-neutral preview boundary and first supported adapters. It is not a full 11-03 closeout.

## Scope completed in this slice

Implemented `mcp_import_preview` preview-only behavior for the first supported JSON/JSONC shapes:

- OpenCode-style `mcp.servers`;
- Codex plugin / common `.mcp.json` wrapper `mcpServers`;
- Codex plugin direct top-level server map form.

Code added or touched:

- `internal/catalog/mcp_import.go`
- `internal/catalog/mcp_import_test.go`
- `internal/management/service.go`
- `internal/management/service_test.go`
- `internal/gateway/server.go`
- `internal/gateway/server_test.go`

## Behavior implemented

- Preview parses JSON plus JSONC line/block comments without corrupting URL strings such as `http://`.
- `format=auto` chooses a deterministic adapter only when the source shape is unambiguous.
- Ambiguous sources return `ambiguous_format` instead of silently guessing.
- OpenCode remote servers normalize to ADM `streamable-http` definitions.
- OpenCode local command arrays normalize to ADM `stdio` executable/args.
- Codex/common remote `url` servers normalize to `streamable-http`.
- Codex/common `command`/`args` servers normalize to `stdio`.
- Deprecated SSE transport is rejected as a candidate-level error instead of being silently converted.
- Source disabled flags are preview warnings only and do not modify ADM Environment membership.
- Import preview defaults `default_include=false` unless the caller explicitly sets the preview option.
- Preview does not persist MCP definitions.

## Secret/reference behavior

- Recognized source env references like `{env:NAME}` normalize to `${NAME}` templates without resolving the environment.
- Template strings such as `Bearer {env:REMOTE_TOKEN}` normalize to `Bearer ${REMOTE_TOKEN}`.
- Literal credential-bearing header/env values are converted to generated reference requirements.
- Authorization literals preserve safe auth-scheme context, for example `Bearer plain-secret-token` becomes `Bearer ${REMOTE_AUTHORIZATION}`.
- Preview output records only generated reference names, candidate names and field paths; it does not echo literal credential values.

## Surface exposed

- `management.Service.MCPImportPreview(...)`
- Gateway tool `mcp_import_preview`

Both surfaces are preview-only and read-only in this slice.

## Remaining 11-03 work

- `mcp_import_apply` is not implemented yet.
- WorkBuddy/CodeBuddy, Claude Code and MCPHub adapters still need committed fixtures and explicit parsing branches.
- Conflict detection and atomic selected-batch apply are still pending.
- Update-by-name semantics and owner session invalidation are still pending.
- Importer CLI surface is still pending.
- Real activation after apply, once references are provisioned, is still pending.
