# Phase 11 Context — MCP Runtime Completion

## Goal

Turn ADM's current external-MCP foundation into a complete, reliable and diagnosable runtime capability that external Agents can use without bypass tools.

## Requirements

- MCP-COMP-01 — explicit desired MCP definition model separated from observed runtime state.
- MCP-COMP-02 — every advertised transport has a real implementation and lifecycle.
- MCP-COMP-03 — auth/secret configuration is explicit and secret values never leak through normal status/errors/logs.
- MCP-COMP-04 — real tool inventory is inspectable/refreshable and Environment disable revokes access immediately.
- MCP-COMP-05 — connection/discovery/call/reconnect failures have actionable structured diagnostics and stale healthy state is impossible.

## Existing foundation

ADM already has:

- global MCP catalog entries;
- Environment enable/disable selection;
- Streamable HTTP activation;
- environment-variable-backed header references;
- four-state configured/disabled/healthy/error health;
- real `tools/list` and `tools/call` through the Gateway;
- Gateway-owned external MCP sessions with restart reconciliation and cleanup.

These are a strong vertical slice, but the persisted MCP definition still shares a generic `CatalogEntry` with Skill fields, transport support is restricted to Streamable HTTP, observed health is relatively coarse, and there is no explicit tool-inventory refresh/status contract.

## Locked decisions

1. MCP is a Core ADM capability. This Phase does not add task planning or tool-selection policy.
2. Supported transports for this Phase are:
   - `streamable-http` for remote/local HTTP MCP;
   - `stdio` for local command MCP.
3. Legacy HTTP+SSE is not added. The current MCP specification deprecates the legacy HTTP+SSE transport.
4. The repository's current official Go MCP SDK `v1.7.0` already supports protocol revision `2026-07-28`; dependency replacement is not required merely for protocol-version support.
5. Persisted desired MCP configuration and Gateway-owner observed runtime status are different models. Session IDs, health observations and tool inventories are never persisted as desired config.
6. Stdio MCP executable launch must obey ADM's executable allowlist. Adding an MCP definition must not create a second arbitrary-process authority path.
7. A stdio MCP is activated per Environment selection under the Gateway owner. Its process working directory is the Environment root; the definition remains global/shareable.
8. Secret values are references in persisted config and resolve only at activation. Status can expose that a reference exists/missing, but never its resolved value.
9. Phase 11 supported HTTP auth modes are `none` and explicit secret-backed headers. Interactive OAuth is not silently approximated; it remains unsupported until ADM has an approved secure credential/token lifecycle.
10. `environment_mcp_refresh` explicitly drops/re-establishes runtime observation and refreshes the actual tool inventory.
11. On connection/list failure, ADM may reconnect for a subsequent safe inspection/list operation. On tool-call failure, ADM drops the bad runtime session but does not automatically replay the tool call because a failed response does not prove the remote side had no side effect.
12. Disabling/removing an MCP or dropping an Environment immediately closes its owned session/process and invalidates observed inventory.
13. One broken MCP affects only that MCP. Files, Skills, verifiers and other MCPs remain usable.

## Non-goals

- no Planner/Executor/Reviewer logic;
- no automatic choice of which MCP tool an Agent should call;
- no legacy SSE transport implementation;
- no insecure OAuth token persistence;
- no generic arbitrary-process launch outside the existing executable authority;
- no Desktop redesign in this Phase.
