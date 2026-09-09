# Plan 11-02 Verification — Interim Monitor/Reconnect Evidence

Date: 2026-09-09

This verification file records the current green evidence for the first Plan 11-02 monitor/reconnect implementation slice. It is not a full Plan 11-02 closeout.

## Focused monitor tests

```text
go test ./internal/gateway -run TestRuntimeOwnerBackground -count=1 -v
```

Observed passing tests:

```text
TestRuntimeOwnerBackgroundMonitorSkipsPingWhenHealthCheckDisabled
TestRuntimeOwnerBackgroundMonitorDoesNotReconnectWhenDisabled
TestRuntimeOwnerBackgroundMonitorReconnectsOnFixedInterval
TestRuntimeOwnerBackgroundReconnectDoesNotResurrectDisabledMCP
```

These cover:

- `health_check_enabled=false` prevents background Ping/reconnect;
- `auto_reconnect=false` prevents background reconnect;
- explicit safe `Status` can still reconnect after background reconnect is disabled;
- `auto_reconnect=true` reconnects after the fixed configured interval;
- failure observation records stage/kind/failure count and clears stale inventory;
- reconnect does not replay arbitrary MCP tool calls;
- disabling the MCP before the reconnect interval prevents background resurrection.

## Restart / desired-vs-observed boundary

```text
go test ./internal/gateway -run TestRuntimeOwnerRestartPersistsHealthPolicyButNotObservation -count=1 -v
```

Result: passed.

This proves:

- desired `MCPHealthPolicy` survives owner restart;
- owner-local observation timestamps/inventory do not persist into the next owner;
- the new owner can rebuild a healthy session from persisted desired state.

## Real Streamable HTTP background reconnect acceptance

```text
go test ./internal/gateway -run TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance -count=1 -v
```

Result:

```text
=== RUN   TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance
--- PASS: TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance (2.24s)
PASS
ok  ai-dev-manager-v2/internal/gateway  2.269s
```

This test uses a real SDK Streamable HTTP MCP server. It proves:

- owner starts healthy and can discover real tool inventory;
- background Ping detects a temporarily broken upstream;
- stale inventory/session are cleared;
- reconnect does not run before the configured fixed interval;
- once the upstream recovers, background reconnect restores healthy state and inventory;
- a later explicit tool call succeeds after recovery.

## Gateway regression split

```text
go test ./internal/gateway -run TestRuntimeOwnerBackground|TestRuntimeOwnerRestart|TestRuntimeOwnerInspect|TestRuntimeOwnerRealHTTPBackgroundReconnectAcceptance -count=1 -v
```

Result: passed.

```text
go test ./internal/gateway -run TestMCP|TestStdio|TestRuntimeOwner|TestGateway|TestHTTPGateway -count=1
```

Result: passed.

Previously recorded gateway split groups also passed:

```text
go test ./internal/gateway -run TestGateway|TestHTTPGateway|TestRealHost|TestStopHTTP|TestSameADM|TestMCP|TestStdio|TestRuntimeOwner -count=1
go test ./internal/gateway -run TestVerifier|TestProcess|TestRun|TestLegacy|TestAllowed|TestOwner|TestWait|TestInspect -count=1
```

The full `internal/gateway` package in one plugin call can exceed the plugin transport observation window, so gateway verification is recorded through split test groups.

## Related package regression

```text
go test ./internal/app ./internal/catalog ./internal/management -count=1
```

Result:

```text
ok  ai-dev-manager-v2/internal/app         26.267s
ok  ai-dev-manager-v2/internal/catalog      1.535s
ok  ai-dev-manager-v2/internal/management   1.774s
```

```text
go test ./cmd/ai-dev-manager ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/environment -count=1
```

Result: passed.

```text
go test ./internal/isolation ./internal/runtime ./internal/skill ./internal/verifier -count=1
```

Result: passed.

```text
go test ./internal/identity ./internal/memory ./internal/model ./internal/store ./internal/workspace -count=1
```

Result: compile/no-test packages passed.

## Vet and diff hygiene

```text
go vet ./...
```

Result: passed.

```text
git diff --check
```

Result: passed. Only existing LF-to-CRLF working-copy warnings were emitted; no whitespace errors were reported.

## Remaining verification gaps before 11-02 closeout

- real stdio unhealthy-to-background-reconnect acceptance still needs stronger end-to-end coverage;
- Gateway API-level `environment_mcp_inspect` assertions should explicitly cover next reconnect timing, in-flight flags, failure stage, failure count and inventory timestamps;
- deeper race/fake-clock style coverage is still incomplete;
- monolithic `go test ./... -count=1` still times out through the plugin transport and should be rerun directly in a local terminal if required for closeout.
