# Plan 11-03 Progress — Import Preview + Apply + Adapter Coverage

Date: 2026-09-09

This remains an implementation progress note, not a full 11-03 closeout. The preview/apply foundation now includes explicit parsing branches for the main planned source families.

## Scope completed so far

Implemented `mcp_import_preview` and `mcp_import_apply` for these shapes:

- OpenCode-style `mcp.servers`;
- Codex plugin / common `.mcp.json` wrapper `mcpServers`;
- Codex plugin direct top-level server map form;
- WorkBuddy explicit `mcpServers` import;
- CodeBuddy explicit `mcpServers` import;
- Claude Code project/plugin `mcpServers` import and single-project `projects.<path>.mcpServers` import;
- MCPHub `mcpServers` and hub-oriented `servers` map import.

Code added or touched so far:

- `internal/catalog/mcp_import.go`
- `internal/catalog/mcp_import_test.go`
- `internal/management/service.go`
- `internal/management/service_test.go`
- `internal/gateway/server.go`
- `internal/gateway/server_test.go`

## Preview behavior implemented

- Parses JSON plus JSONC line/block comments without corrupting URL strings such as `http://`.
- `format=auto` chooses a deterministic adapter when the source shape is unambiguous or has a source hint.
- Generic `mcpServers` remains a Codex/common shape in auto mode unless a source-specific hint is present.
- Sources with both OpenCode `mcp.servers` and common `mcpServers` return `ambiguous_format` instead of silently guessing.
- OpenCode remote servers normalize to ADM `streamable-http` definitions.
- OpenCode local command arrays normalize to ADM `stdio` executable/args.
- WorkBuddy/CodeBuddy/Codex/Claude/MCPHub remote URL or `streamableHttp`/HTTP aliases normalize to `streamable-http`.
- WorkBuddy/CodeBuddy/Codex/Claude/MCPHub `command`/`args` entries normalize to `stdio`.
- Claude Code full `projects` import is accepted only when there is one project scope; multiple scopes return `scope_selector_required` until an explicit selector surface is added.
- MCPHub hub-oriented `servers` map is accepted while hub-only fields become warnings.
- Deprecated SSE transport is rejected as a candidate-level error instead of being silently converted.
- Source `disabled`/`enabled` flags are preview warnings only and do not modify ADM Environment membership.
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
- Claude-style `${VAR:-default}` templates are preserved as references/templates without resolution.
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

- Importer CLI surface is still pending.
- Real activation after apply, once generated references are provisioned, is still pending.
- Update-by-name owner invalidation has Gateway coverage through drop-on-success plumbing but still needs a direct runtime-owner regression if required for closeout.
- WorkBuddy/CodeBuddy/Claude/MCPHub coverage uses representative inline fixtures; separate fixture files can still be added if the closeout standard requires file-backed fixtures.
- Formal 11-03 summary and closeout verification remain pending.
