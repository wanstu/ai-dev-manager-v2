# Phase 07 — Desktop dogfood and release readiness

## Requirement IDs

- ADM-DESKTOP-001
- ADM-DEV-001
- ADM-DEV-003

## Goal

Turn the Phase 06 desktop manager from a development-only `go build` target into a reproducible Wails release build, then use that build for concrete desktop dogfood before adding speculative desktop features.

## Delivery slices

### Slice A — Reproducible Wails release build

Status: delivered.

- keep `wails.json` beside the desktop Go `main` in `cmd/ai-dev-manager-desktop`
- build from that Wails project root so binding generation sees the correct Go package
- provide `scripts/build-desktop.ps1` as the root-level release build entrypoint
- pin the invoked Wails CLI to v2.15.0
- ignore generated Wails bindings and build output
- verify the official Wails build produces `build/bin/ai-dev-manager-v2-desktop.exe`
- smoke-launch the built desktop executable against a temporary ADM state directory and verify the process remains alive

### Slice B — Desktop dogfood

Use the built desktop executable against real ADM state for representative management tasks. Record and fix concrete usability or lifecycle blockers before pulling forward tray, autostart, single-instance handling, notifications, installers, or other desktop polish.

## Acceptance tests for Slice A

- running Wails from the repository root without the corrected project layout is known to fail with `no Go files`
- `cmd/ai-dev-manager-desktop/wails.json` is colocated with the desktop Go main
- no root `wails.json` remains
- `scripts/build-desktop.ps1` runs Wails from the desktop command directory
- the script completes successfully on the current Windows development machine
- the release executable is created under `cmd/ai-dev-manager-desktop/build/bin`
- the release executable starts and remains alive during a short smoke run using temporary ADM state
- generated `frontend/wailsjs` and `build` output are ignored
- `go test ./...`, `go vet ./...`, and `git diff --check` pass

## Explicit non-goals for Slice A

- tray integration
- autostart
- single-instance handling
- notifications
- installer design
- code signing
- automatic updates
- REST API
- Git/worktree lifecycle
- isolation/orchestration

## New prerequisites

None beyond the Wails v2 desktop dependency already introduced in Phase 06.
