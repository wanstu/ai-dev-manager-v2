# Phase 11 Research — MCP Runtime Completion

## Current implementation audit

### Persisted model

`internal/model/types.go` currently stores MCP and Skill in the same `CatalogEntry`. MCP-specific fields are `Endpoint`, `Transport`, and `HeaderRefs`; Skill-specific fields live beside them. This makes desired MCP configuration difficult to validate as a transport-specific contract and encourages unrelated MCP/Skill changes to share one persistence type.

`internal/catalog/service.go` currently:

- defaults transport to `streamable-http`;
- accepts only HTTP(S) endpoints;
- rejects every transport other than `streamable-http`;
- stores environment-variable-like header references;
- provides global catalog/default-selection behavior.

Because compatibility is not a locked requirement, Phase 11 can introduce a typed MCP definition instead of preserving the generic transport fields indefinitely.

### Activation and transient health

`internal/app/mcp_health.go` already separates activation from persisted config in a useful way:

- Environment selection is checked first;
- missing catalog entries are reported as configured/unresolved;
- secret/header references resolve only at activation;
- on-demand probes are bounded;
- errors are classified into timeout, refused, auth and connection categories.

The current activation type is HTTP-specific (`Endpoint`, `Headers`) and needs to become transport-specific without leaking resolved secrets.

### Persistent Gateway owner

`internal/gateway/runtime_owner.go` owns MCP runtime state by `(environment_id, mcp_id)`. It already:

- reuses live sessions;
- verifies a session using `ListTools`;
- drops failed sessions;
- closes sessions on MCP disable/removal, Environment drop and Gateway shutdown;
- reconstructs desired runtime from persisted Environment selections after restart;
- keeps observed health owner-local rather than persisting it.

This is the correct owner boundary for Phase 11. Do not create a second MCP daemon.

### Missing runtime concepts

The current owner does not expose a rich observed record containing connection stage, timestamps or tool inventory. `ListTools` is both health check and inventory fetch, but there is no explicit refresh operation or auditable inventory status.

`CallTool` drops a session after failure. This is correct for future recovery, but any automatic immediate retry of the tool call would be unsafe because tool calls can mutate external systems. Phase 11 must explicitly forbid automatic replay.

## MCP protocol / SDK facts

Repository dependency: `github.com/modelcontextprotocol/go-sdk v1.7.0`.

The official Go SDK compatibility table states `v1.7.0+` supports MCP specification `2026-07-28` while retaining older supported protocol revisions. The current MCP specification deprecates legacy HTTP+SSE. The SDK also exposes `CommandTransport` for starting a local command and communicating over stdin/stdout.

Therefore Phase 11 does not need an SDK upgrade merely to speak the current protocol and should not add deprecated SSE just to claim transport breadth.

References:

- https://github.com/modelcontextprotocol/go-sdk
- https://blog.modelcontextprotocol.io/posts/2026-07-28/

## Recommended desired model

Introduce an MCP-specific persisted type rather than adding more fields to generic `CatalogEntry`.

Conceptual shape:

```text
MCPDefinition
  id
  name
  default_include_in_environment
  transport: streamable-http | stdio
  auth_mode: none | headers
  endpoint                 # streamable-http only
  header_refs              # streamable-http / headers auth only
  executable               # stdio only
  args                     # stdio only
  env_refs                 # stdio optional activation refs
```

Validation must be transport-local:

- HTTP requires valid http(s) endpoint and forbids stdio fields.
- stdio requires an executable and forbids HTTP endpoint/header fields.
- unsupported transport/auth is rejected before persistence/activation.
- secret refs store reference names/templates, not resolved values.

Stdio activation should use the selected Environment root as process cwd. The executable must pass the same executable allowlist authority as normal Runtime command execution so `mcp_add` cannot become an arbitrary command launcher.

## Recommended observed model

Owner-local runtime observation should include enough facts for diagnostics without becoming persisted desired state:

```text
MCPRuntimeObservation
  mcp_id
  environment_id
  state
  transport
  stage: activation | connect | discover | list_tools | call
  error_kind
  message
  last_checked_at
  last_success_at
  tool_inventory
  inventory_fetched_at
```

Tool inventory should contain public tool metadata already returned by MCP and remain bounded. No header/env secret values belong in this record.

## Refresh/reconnect semantics

- `environment_mcp_refresh` is explicit and safe: close/drop current observation, activate again, fetch real tool inventory.
- `environment_mcp_tools` may ensure a usable connection and fetch tools if inventory is absent/stale.
- failed safe discovery/list operations may reconnect on a subsequent request.
- failed `environment_mcp_call` drops the connection but never silently retries the same tool invocation.
- Environment disable/remove immediately invalidates the session/process and inventory.
- Gateway restart reconstructs from desired config; old observations are not treated as current truth.

## Auth boundary

Current environment-variable-backed HTTP headers are retained as the supported secret mechanism. Interactive OAuth requires a secure token/credential lifecycle that ADM does not currently have. Phase 11 must report unsupported auth explicitly rather than store refresh/access tokens insecurely or invent a half-OAuth implementation.

## Main risks

1. **stdio as exec bypass** — mitigate by requiring executable allowlist authority.
2. **secret leakage** — test all status/error/management output for resolved secret absence.
3. **unsafe retry** — never replay a tool call automatically.
4. **stale healthy/inventory** — drop observation on transport error, disable/remove and restart.
5. **optional-capability coupling** — one broken MCP must not block files/Skill/other MCP operations.
6. **model churn** — because this is pre-stable development, prefer a clean typed MCP model over compatibility fields.
