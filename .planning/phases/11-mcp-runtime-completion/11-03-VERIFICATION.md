# Plan 11-03 Verification — Import Preview Foundation

Date: 2026-09-09

This verification records the initial Plan 11-03 import preview implementation slice. It is not a full Plan 11-03 closeout.

## Focused catalog importer tests

```text
go test ./internal/catalog -run TestPreviewMCPImport -count=1 -v
```

Result: passed.

Covered cases:

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

## Management boundary test

```text
go test ./internal/management -run TestManagementMCPImportPreview -count=1 -v
```

Result: passed.

Covered cases:

- shared management preview boundary delegates to the source-neutral importer;
- literal bearer token is converted to `Bearer ${REMOTE_AUTHORIZATION}`;
- reference requirement contains only generated reference name and field path;
- preview output does not contain the literal token;
- preview does not persist any MCP definitions.

## Gateway preview test

```text
go test ./internal/gateway -run TestGatewayDevelopsPlainDirectoryWithoutGit|TestGatewayMCPImportPreviewIsSanitizedAndReadOnly|TestGatewayMCPAddRejectsLiteralSecretRefs -count=1 -v
```

Result: passed.

Covered cases:

- `mcp_import_preview` is present in the Gateway tool list;
- Gateway preview returns generated reference requirements;
- Gateway preview does not echo literal credentials;
- Gateway preview does not persist MCP definitions;
- existing `mcp_add` literal-secret rejection remains green.

## Hygiene

This slice should be followed by:

```text
go test ./internal/catalog ./internal/management ./internal/gateway -count=1
go vet ./...
git diff --check
```

before commit.

## Remaining verification gaps before 11-03 closeout

- `mcp_import_apply` and atomic selected batch apply are not covered yet.
- Conflict detection/update-by-name behavior is not covered yet.
- WorkBuddy/CodeBuddy, Claude Code and MCPHub fixtures are not covered yet.
- CLI import surface is not covered yet.
- Real runtime activation after apply with provisioned references is not covered yet.
