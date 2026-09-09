# Plan 11-03 Progress — Import Preview + Apply Foundation

Date: 2026-09-09

This remains an implementation progress note, not a full 11-03 closeout. The preview foundation is now joined by the first atomic apply slice for common Codex/OpenCode-style JSON/JSONC shapes.

## Scope completed so far

Implemented `mcp_import_preview` and the first `mcp_import_apply` behavior for these supported shapes:

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

## Preview behavior implemented

- Parses JSON plus JSONC line/block comments without corrupting URL strings such as `http://`.
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

## Apply behavior implemented

- `catalog.MCPService.ApplyMCPImport(...)` reparses and revalidates the same source through the preview pipeline before writing.
- `selected_names` can import a chosen subset of candidates.
- Default `conflict_policy` is `error`.
- `conflict_policy=skip` skips conflicting names and imports only non-conflicting selected candidates.
- `conflict_policy=update_by_name` preserves the existing ADM MCP ID while replacing desired config by name.
- Selected batch apply is atomic: one invalid selected candidate or default conflict error prevents all selected writes.
- Apply writes only global MCP definitions and never modifies existing Environment selections.
- Gateway `mcp_import_apply` drops affected owner-local MCP runtime state after successful import/update so stale sessions/observations are not retained.

## Secret/reference behavior

- Recognized source env references like `{env:NAME}` normalize to `${NAME}` templates without resolving the environment.
- Template strings such as `Bearer {env:REMOTE_TOKEN}` normalize to `Bearer ${REMOTE_TOKEN}`.
- Literal credential-bearing header/env values are converted to generated reference requirements.
- Authorization literals preserve safe auth-scheme context, for example `Bearer plain-secret-token` becomes `Bearer ${REMOTE_AUTHORIZATION}`.
- Preview/apply output records only generated reference names, candidate names and field paths; it does not echo literal credential values.
- Apply persists only the generated reference/template value, not the literal credential.

## Surfaces exposed

- `management.Service.MCPImportPreview(...)`
- `management.Service.MCPImportApply(...)`
- Gateway tool `mcp_import_preview`
- Gateway tool `mcp_import_apply`

## Remaining 11-03 work

- WorkBuddy/CodeBuddy, Claude Code and MCPHub adapters still need committed fixtures and explicit parsing branches.
- Importer CLI surface is still pending.
- Real activation after apply, once references are provisioned, is still pending.
- Update-by-name owner invalidation has Gateway coverage through drop-on-success plumbing but still needs a direct runtime-owner regression if required for closeout.
- Formal 11-03 summary and closeout verification remain pending.
