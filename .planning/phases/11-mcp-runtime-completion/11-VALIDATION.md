# Phase 11 Validation — MCP Runtime Completion

## Validation principle

MCP completion is accepted only through real runtime consumption. Catalog CRUD or unit-only transport validation is insufficient.

## Required proof matrix

### P1 — typed desired configuration

Prove valid `streamable-http` and `stdio` definitions persist with only their allowed fields. Invalid endpoint, missing stdio executable, mixed transport fields, unsupported transport and unsupported auth are rejected before activation.

Negative requirement: MCP configuration must not persist resolved secret values.

### P2 — Streamable HTTP remains real

A real Streamable HTTP MCP server must initialize/negotiate through the repository's current official Go SDK, list tools and execute a tool through ADM. Existing Environment enable/disable gating remains authoritative.

### P3 — real stdio MCP

Start a real helper MCP as a child process through the Gateway owner using `stdio`, list at least one tool and call it successfully from a later Agent client.

Negative requirements:

- an executable not in ADM's allowlist cannot be activated as stdio MCP;
- escaped/arbitrary cwd is not accepted; stdio process cwd is the selected Environment root;
- disabling/removing the MCP or closing the Gateway terminates the owned child process.

### P4 — secrets stay secret

Use a sentinel secret supplied through an environment reference. Prove activation succeeds while the sentinel never appears in:

- persisted ADM state;
- MCP definition/status output;
- runtime observation;
- normal error text/log output.

Missing secret ref must return a structured configuration/activation error without exposing unrelated environment variables.

### P5 — real tool inventory and refresh

Prove tool inventory is fetched from the actual server, includes a fetch timestamp/count/public tool metadata, and is invalidated/refreshed explicitly.

Change the helper MCP inventory between refreshes and prove `environment_mcp_refresh` exposes the new inventory rather than stale owner state.

Environment disable must immediately remove access and observed inventory.

### P6 — lifecycle and safe reconnect

Prove:

- broken connection becomes structured error, never stale healthy;
- subsequent safe status/list/refresh can establish a new connection after server recovery;
- a failed tool call invalidates the bad connection;
- ADM does not automatically replay that tool call.

Use a helper tool that records invocation count so automatic replay would be machine-detectable.

### P7 — restart semantics

With an enabled MCP, stop/restart the real Gateway owner. Desired config/Environment selection persists, old observed session/inventory does not, and a later request rebuilds or reconciles a healthy runtime from desired state.

### P8 — optional failure is local

With one broken MCP and one healthy MCP in the same Environment, prove:

- healthy MCP remains callable;
- file read/search remains usable;
- Skill list/read remains usable;
- verifier/file/runtime capabilities are not globally marked failed.

## Automated gates

- focused catalog/config tests;
- focused HTTP/stdio MCP runtime tests;
- race tests for owner session/refresh/drop paths;
- real Streamable HTTP acceptance;
- real stdio child-process acceptance;
- Gateway restart acceptance;
- `go test ./...` or equivalent complete package-group coverage if the tool transport cannot observe one long aggregate call;
- `go vet ./...`;
- `git diff --check`.

## Review focus

Integration review must explicitly inspect:

- no secret values in persisted/returned structures;
- no automatic tool-call replay;
- no stdio exec-policy bypass;
- no legacy SSE implementation;
- no Planner/GSD orchestration reintroduction;
- no MCP failure promoted into a global Environment prerequisite.
