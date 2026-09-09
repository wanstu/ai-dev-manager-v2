# Phase 16 Validation — Desktop Core Parity + RC Readiness

## Phase-level gates

Phase 16 must prove Desktop is a safe management surface over existing Core and that the repository has an automated CI/build path for RC candidates. It must not introduce Desktop-only semantics or bypass Core authorization.

Required validation categories:

1. Core parity inventory is complete enough for daily local use.
2. Sensitive state remains protected: secret values and private Memory values are not displayed by default.
3. Writer-gated operations still require writer ownership and do not bypass Environment lease rules.
4. Runtime operations still obey executable allowlist and Environment containment.
5. Optional capability failures are displayed as local/degraded/unavailable facts, not as whole-app failure.
6. Desktop calls existing application state/services or Gateway-compatible surfaces; no duplicated product state.
7. CI runs tests/vet/build automatically and uploads short-retention artifacts without publishing releases.
8. RC gate has a manual smoke checklist and automated regression gate.

## CI baseline gates

For the GitHub Actions workflow:

1. Trigger on `push` to `master`, `pull_request`, and `workflow_dispatch`.
2. Use read-only repository permissions.
3. Use Go version information from `go.mod`.
4. Run Windows full `go test -count=1 ./...`.
5. Run `go vet ./...`.
6. Build Windows CLI and Desktop smoke binaries.
7. Build CLI artifacts on Linux/macOS/Windows where ordinary `go build` supports them.
8. Upload artifacts with short retention.
9. Do not deploy, publish releases, push commits or require secrets.

## Automated gates for implementation slices

At minimum:

- focused Desktop tests for changed UI/adapter behavior;
- relevant app/gateway/management tests if service boundaries change;
- `go test -count=1 ./...`;
- `go vet ./...`;
- `git diff --check`.

## RC smoke checklist

Before a local 1.0 RC is declared, manually verify:

1. Desktop launches and can read existing ADM state.
2. Workspace/Environment list and inspect work.
3. Capability report is visible with unavailable/degraded reasons.
4. MCP definitions/selection/health status are manageable or clearly linked to CLI/Gateway fallback if not yet surfaced.
5. Skill sources/selection/availability are manageable or clearly linked to CLI/Gateway fallback if not yet surfaced.
6. Process/run/verifier results are visible enough for daily use.
7. Sensitive values are not exposed in normal Desktop views.
8. CI artifacts are produced and downloadable for smoke review.
9. Known limitations are documented.

## Non-blocking by default

The following remain post-RC unless dogfood proves they block daily use:

- GitNexus/provider integration;
- additional Phase 14 investigation helpers;
- temporary resource lifecycle cleanup;
- installer/tray/autostart/updater/signing/notifications;
- automatic GitHub Release publishing.
