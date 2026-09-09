# Plan 11-01 Verification — 2026-09-09

## Verification status

Current result: **split-package green**.

A monolithic `go test ./... -count=1` was attempted through the ChatGPT plugin transport but timed out at the tool-call layer before a final result could be observed. The package-group runs below cover the repository package list returned by `go list ./...` and are the current machine-observed verification evidence.

## Package list

`go list ./...` returned:

```text
ai-dev-manager-v2/cmd/ai-dev-manager
ai-dev-manager-v2/cmd/ai-dev-manager-desktop
ai-dev-manager-v2/internal/app
ai-dev-manager-v2/internal/catalog
ai-dev-manager-v2/internal/desktop
ai-dev-manager-v2/internal/environment
ai-dev-manager-v2/internal/gateway
ai-dev-manager-v2/internal/identity
ai-dev-manager-v2/internal/isolation
ai-dev-manager-v2/internal/management
ai-dev-manager-v2/internal/memory
ai-dev-manager-v2/internal/model
ai-dev-manager-v2/internal/runtime
ai-dev-manager-v2/internal/skill
ai-dev-manager-v2/internal/store
ai-dev-manager-v2/internal/verifier
ai-dev-manager-v2/internal/workspace
```

## Passing commands

```text
go test ./internal/app ./internal/catalog ./internal/management ./internal/gateway -count=1
ok  ai-dev-manager-v2/internal/app         7.751s
ok  ai-dev-manager-v2/internal/catalog     0.315s
ok  ai-dev-manager-v2/internal/management  0.414s
ok  ai-dev-manager-v2/internal/gateway     57.415s
```

```text
go test ./cmd/ai-dev-manager ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/environment -count=1
ok  ai-dev-manager-v2/cmd/ai-dev-manager          0.753s
ok  ai-dev-manager-v2/cmd/ai-dev-manager-desktop  0.056s
ok  ai-dev-manager-v2/internal/desktop            0.260s
ok  ai-dev-manager-v2/internal/environment        0.216s
```

```text
go test ./internal/isolation ./internal/runtime ./internal/skill ./internal/verifier -count=1
ok  ai-dev-manager-v2/internal/isolation  8.049s
ok  ai-dev-manager-v2/internal/runtime    0.780s
ok  ai-dev-manager-v2/internal/skill      0.077s
ok  ai-dev-manager-v2/internal/verifier   0.043s
```

```text
go test ./internal/identity ./internal/memory ./internal/model ./internal/store ./internal/workspace -count=1
?   ai-dev-manager-v2/internal/identity   [no test files]
?   ai-dev-manager-v2/internal/memory     [no test files]
?   ai-dev-manager-v2/internal/model      [no test files]
?   ai-dev-manager-v2/internal/store      [no test files]
?   ai-dev-manager-v2/internal/workspace  [no test files]
```

```text
git diff --check
```

`git diff --check` exited 0 with only existing LF-to-CRLF working-copy warnings and no whitespace errors.

## Focused tests exercised during repair

Additional focused runs were used while closing regressions:

```text
go test ./internal/gateway -run TestRuntimeOwnerRealHTTPRestartReconciliation -count=1 -v
PASS
```

```text
go test ./internal/gateway -run TestStdioMCPAcceptanceAndProcessCleanup -count=1 -v
PASS
```

```text
go test ./internal/gateway -run TestGatewayMCPAddRejectsLiteralSecretRefs|TestGatewayUsesOnlyEnabledConfiguredExternalMCPAndSkillRuntime -count=1 -v
PASS
```

```text
go test ./internal/management -run TestManagementMCPAddConfigUsesTypedModelAndSecretBoundary -count=1 -v
PASS
```

## Requirements covered

- P1 typed desired configuration: covered by catalog tests and CLI/management tests.
- P2 Streamable HTTP runtime: covered by Gateway external MCP list/call acceptance.
- P3 real stdio MCP: covered by stdio acceptance and process cleanup test.
- P4 secrets stay secret: covered by catalog, CLI, management and Gateway literal-secret negative tests plus activation persistence test.
- P8 restart semantics for owner desired-state reconciliation: partially covered by `TestRuntimeOwnerRealHTTPRestartReconciliation` for owner identity, desired selection persistence, stale observation reset and dead upstream behavior.
- P9 optional failure stays local: covered by Gateway tests proving broken MCP does not break unrelated Gateway/file operations.

## Remaining validation outside 11-01

- P5 inventory refresh is partly implemented and belongs to 11-02 observation/refresh closure.
- P6 periodic health check and automatic reconnect are not implemented yet.
- P7 automatic recovery/no replay needs the 11-02 reconnect lifecycle before full validation.
- P10 external import adapters belong to 11-03 and are not implemented yet.

## Review caveat

This evidence makes the current working tree suitable for Plan 11-01 review, not for declaring all Phase 11 complete. Before commit, a local terminal may optionally rerun the monolithic command:

```text
go test ./... -count=1
```

The plugin transport has already shown that this command can exceed the observable tool-call window even when the equivalent package groups pass.
