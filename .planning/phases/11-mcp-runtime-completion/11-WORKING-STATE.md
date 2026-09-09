# Phase 11 Working State — 2026-09-09

This file records the live Phase 11 repository state so new sessions do not mistake older planning snapshots or chat-only context for the current code checkpoint.

## Git snapshot

- Branch: `master`.
- Last completed checkpoint before Plan 11-03 work: `31fed37 docs(11): close mcp monitor evidence`.
- Local master was `ahead 2` relative to `origin/master` after the 11-02 closeout checkpoint.
- Current worktree contains the first Plan 11-03 import preview implementation slice.

Prior Phase 11 checkpoints:

- `3b39c40 feat(11): add typed MCP runtime and monitor`
- `b47c370 test(11): add stdio MCP reconnect acceptance`
- `31fed37 docs(11): close mcp monitor evidence`

## Plan 11-01 — Typed MCP Configuration + HTTP/Stdio Runtime

Status: **implemented and documented for review**.

Evidence:

- `.planning/phases/11-mcp-runtime-completion/11-01-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-01-VERIFICATION.md`

Implemented behavior includes typed MCP definitions, HTTP/stdio activation/runtime, CLI and management typed add surfaces, Gateway typed configuration, secret/reference boundary enforcement, stdio allowlist authority, real HTTP/stdio acceptance, Ping health probing and owner context-boundary fixes.

Residual review item:

- Desktop/UI parity has not been fully reviewed beyond adapter type changes.

## Plan 11-02 — Health Monitor, Automatic Reconnect, Inventory + Diagnostics

Status: **implemented and documented for review**.

Evidence:

- `.planning/phases/11-mcp-runtime-completion/11-02-SUMMARY.md`
- `.planning/phases/11-mcp-runtime-completion/11-02-VERIFICATION.md`

Implemented behavior includes owner-local observation, inspect/refresh, Ping-based monitor, fixed-interval auto reconnect, disabled/drop cleanup, no tool replay, desired-vs-observed restart boundary, real Streamable HTTP and real stdio background reconnect acceptance, Gateway API-level inspect assertions and targeted race coverage.

## Plan 11-03 — JSON / JSONC Import Adapters

Status: **in progress; import preview foundation implemented in the current worktree**.

Current slice includes:

- source-neutral catalog preview boundary in `internal/catalog/mcp_import.go`;
- JSON + JSONC parser that strips comments without corrupting URL strings;
- `format=auto` deterministic detection with `ambiguous_format` failure when more than one adapter matches;
- OpenCode-style `mcp.servers` preview support;
- Codex plugin/common `.mcp.json` `mcpServers` wrapper preview support;
- Codex plugin direct top-level server map preview support;
- local command-array / command+args normalization to ADM `stdio` config;
- remote URL normalization to ADM `streamable-http` config;
- deprecated SSE rejection;
- source disabled flag warning without changing Environment selections;
- `{env:NAME}` and `Bearer {env:NAME}` normalization to `${NAME}` templates;
- literal credential-bearing header/env conversion to generated reference requirements without echoing secret values;
- management preview surface `MCPImportPreview`;
- Gateway tool `mcp_import_preview`;
- read-only behavior: preview does not persist MCP definitions.

Current evidence files:

- `.planning/phases/11-mcp-runtime-completion/11-03-PROGRESS.md`
- `.planning/phases/11-mcp-runtime-completion/11-03-VERIFICATION.md`

Remaining 11-03 work:

- `mcp_import_apply` and atomic selected-batch apply;
- conflict detection, skip and update-by-name policies;
- WorkBuddy/CodeBuddy adapter fixtures;
- Claude Code adapter fixtures;
- MCPHub adapter fixtures;
- CLI import surface;
- real activation after apply once generated references are provisioned;
- formal 11-03 closeout summary/verification.

## Verification snapshot

11-01/11-02 verification is recorded in their respective verification files and remains green at the latest checkpoints.

Current 11-03 preview slice has passed focused tests for catalog importer preview, management preview boundary and Gateway `mcp_import_preview` read-only/sanitized behavior. Run the broader split-package and hygiene commands before committing this slice.

## Immediate next action

Complete verification for the current 11-03 preview slice, commit it as a small checkpoint, then continue with `mcp_import_apply` and conflict/atomicity semantics.

Do not restart 11-01/11-02 from older assumptions. Do not automatically merge or push.
