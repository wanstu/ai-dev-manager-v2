# Plan 11-03 Verification — Import Preview + Apply Foundation

Date: 2026-09-09

This verification records the current Plan 11-03 import preview/apply implementation slice. It is not a full Plan 11-03 closeout.

## Focused catalog importer tests

```text
go test ./internal/catalog -run TestApplyMCPImport|TestPreviewMCPImport -count=1 -v
```

Result: passed.

Covered preview cases:

- OpenCode JSONC `mcp.servers` with line comments and URL strings;
- OpenCode local command-array normalization to stdio executable/args;
- OpenCode remote URL normalization to Streamable HTTP;
- source disabled flag reported as a warning only;
- source `{env:NAME}` reference normalization;
- template `Bearer {env:NAME}` normalization to `Bearer ${NAME}`;
- literal credential-bearing env/header conversion to generated reference requirements;
- preview does not leak literal credential values;
- Codex/common `mcpServers` wrapper;
- Codex/common direct top-level server map;
- `format=auto` deterministic detection;
- ambiguous source returns `ambiguous_format`;
- legacy SSE transport is rejected as unsupported.

Covered apply cases:

- default conflict policy is `error` and does not persist partial batch state;
- `conflict_policy=skip` imports only non-conflicting candidates;
- `conflict_policy=update_by_name` preserves existing ADM MCP ID while replacing desired config by name;
- `selected_names` imports an explicit subset;
- invalid selected batch is atomic and writes none of the selected candidates;
- literal credential-bearing headers are converted to generated reference templates before persistence;
- apply result and persisted MCP definitions do not leak literal credential values.

## Management boundary tests

```text
go test ./internal/management -run TestManagementMCPImport -count=1 -v
```

Result: passed.

Covered cases:

- shared management preview boundary delegates to the source-neutral importer;
- shared management apply boundary delegates through the canonical MCP service;
- literal bearer token is converted to `Bearer ${REMOTE_AUTHORIZATION}`;
- reference requirements contain only generated reference names and field paths;
- preview output does not contain literal tokens;
- preview does not persist MCP definitions;
- apply persists sanitized generated references only.

## Gateway tests

```text
go test ./internal/gateway -run TestGatewayMCPImportPreview|TestGatewayMCPImportApply|TestGatewayDevelopsPlainDirectoryWithoutGit -count=1 -v
```

Result: passed.

Covered cases:

- `mcp_import_preview` and `mcp_import_apply` are present in the Gateway tool list;
- Gateway preview returns generated reference requirements;
- Gateway preview does not echo literal credentials;
- Gateway preview does not persist MCP definitions;
- Gateway apply persists selected imported MCP definitions;
- Gateway apply does not modify Environment selections;
- Gateway apply default conflict policy rejects atomically;
- Gateway apply `update_by_name` replaces desired config by name;
- existing plain-directory Gateway behavior remains green.

## Package regression and hygiene

```text
go test ./internal/catalog ./internal/management -count=1
go test ./internal/gateway -run TestGateway|TestMCP|TestStdio|TestRuntimeOwner -count=1
go vet ./internal/catalog ./internal/management ./internal/gateway
git diff --check
```

Result: passed in this slice. `git diff --check` emitted only LF-to-CRLF working-copy warnings and no whitespace errors.

## Remaining verification gaps before 11-03 closeout

- WorkBuddy/CodeBuddy, Claude Code and MCPHub fixtures are not covered yet.
- CLI import surface is not covered yet.
- Real runtime activation after apply with provisioned references is not covered yet.
- Owner drop after Gateway `update_by_name` is wired but could use a dedicated runtime-owner regression if required before closeout.
