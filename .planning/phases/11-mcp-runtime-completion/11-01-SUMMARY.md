# Plan 11-01 Summary — Typed MCP Configuration + HTTP/Stdio Runtime

## Status

Plan 11-01 is implemented in the current uncommitted working tree and has split-package green verification. It should still receive a human/integration review before commit because the working tree also contains adjacent Phase 11-02 observation scaffolding.

## Implemented scope

- Added a dedicated persisted `MCPDefinition` model separate from Skill `CatalogEntry`.
- Added `MCPHealthPolicy` desired config fields for explicit health/reconnect policy.
- Added a dedicated MCP catalog service with transport-specific validation.
- Supported `streamable-http` and `stdio` MCP transports.
- Rejected mixed transport fields, unsupported transports and unsupported auth modes.
- Preserved Skill catalog behavior behind its existing generic catalog service.
- Resolved MCP activation at Environment boundary through `ResolveMCPActivation`.
- Resolved HTTP endpoints/header references and stdio env references only at activation.
- Kept resolved secret values out of persisted desired state.
- Added Streamable HTTP SDK connection and stdio command transport activation through ADM runtime authority.
- Routed stdio executable startup through the existing ADM Runtime command preparation and executable allowlist.
- Ensured stdio MCP child process cleanup on Environment disable and owner close.
- Added Gateway typed `mcp_add` support.
- Added CLI typed `mcp add` support for HTTP, stdio, header/env refs and health policy flags.
- Added management typed `MCPAddConfig` support.
- Added secret-boundary validation so HTTP headers and credential-bearing stdio env values cannot be persisted as ordinary literals.
- Added acceptance coverage for real Streamable HTTP list/call and real stdio list/call.
- Added no-owner Gateway compatibility behavior so direct list/call consumption is not rejected only because Ping is unsupported or protocol-incompatible.
- Added runtime-owner context boundary so owner-owned upstream MCP operations use the owner lifecycle context instead of the outer Gateway request context.

## Explicit non-goals / not complete here

- Phase 11-02 periodic background health monitor is not complete.
- Phase 11-02 automatic reconnect scheduler is not complete.
- Phase 11-03 import adapters are not started.
- Desktop/UI parity has not been fully reviewed beyond adapter compatibility changes.
- No push or merge is implied by this summary.

## Notes for review

The secret/reference boundary intentionally distinguishes HTTP `HeaderRefs` from stdio `EnvRefs`:

- HTTP headers are always treated as sensitive configuration and must use environment references.
- Stdio env keys that appear credential-bearing, such as token/secret/password/auth/api-key keys, must use environment references.
- Plain stdio runtime env keys may remain literal so helper/runtime configuration such as flags or root directories can still work.

This preserves the Phase 11 requirement that credential-bearing literals must not become ordinary persisted MCP config while avoiding false positives for non-secret stdio environment values.
