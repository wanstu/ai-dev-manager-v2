# 15-02 Validation — Temporary Environment Cleanup and Bounded Inspection

Date: 2026-09-11
Status: working tree validated; live `pjadm` Gateway verification requires rebuild/restart.

## Scope validated

- Temporary ordinary Environment cleanup can remove only the ADM Environment state record.
- Workspace records and ordinary project directories are preserved during ordinary Environment cleanup.
- Temporary ADM-managed worktree cleanup is delegated through existing isolation destroy safety instead of duplicating deletion logic.
- Dirty managed worktrees are blocked from retention cleanup.
- Runtime-owner evidence is required for Environment and MCP cleanup execution paths that depend on Gateway-owned activity observations.
- Environment capability reporting now bounds default output by returning catalog summary facts plus selected MCP/Skill detail facts only.

## Commands run

```text
go test -count=1 ./internal/app -run TestCapabilityReport
go test -count=1 ./internal/gateway -run TestGatewayCapabilityReport
go test -count=1 ./internal/app -run TestResourceRetention
go test -count=1 ./internal/gateway -run TestGatewayResourceRetention
go test -count=1 ./internal/adminmcp ./internal/management
go test -count=1 ./cmd/...
go vet ./internal/app ./internal/gateway ./internal/adminmcp ./internal/management ./cmd/...
git diff --check
```

## Dogfood note

The bounded capability inspection source fix is validated by unit tests, but the already-running `pjadm` Gateway still serves the old code until it is rebuilt/restarted. Live `environment_inspect` should not be used as the acceptance signal before restart because it still emits the old unbounded disabled Skill facts.

## Remaining before closeout

- Rebuild/restart the local ADM Gateway used by `pjadm`.
- Re-run live `environment_inspect` once after restart and confirm the capability report contains catalog summary facts instead of unselected disabled Skill/MCP detail facts.
- Commit the 15-02 code and planning updates only after the live dogfood check passes.
