# Phase 16 RC Gate — Local 1.0 Release Candidate

## RC definition

A 1.0 RC is a local release candidate suitable for daily dogfood. It is not a public packaged release and does not require installer, tray, autostart, updater, signing or notifications.

## Required automated gate

Before declaring a local 1.0 RC:

1. GitHub Actions CI must exist and run on `master`, pull requests and manual dispatch.
2. CI must run full Windows `go test -count=1 ./...`.
3. CI must run `go vet ./...`.
4. CI must build the Windows CLI binary.
5. CI must build a Windows Desktop smoke binary.
6. CI should build CLI artifacts for Windows/Linux/macOS when supported by normal `go build`.
7. CI must upload short-retention artifacts for smoke review.
8. CI must not publish releases, push commits, require secrets or deploy anything.

## Required local gate

Before declaring RC locally:

1. `go test -count=1 ./...` passes.
2. `go vet ./...` passes.
3. `git diff --check` passes, allowing only platform line-ending warnings if there are no whitespace errors.
4. CLI smoke build passes.
5. Desktop smoke build passes.
6. `git status` is clean after the RC commit.
7. No push unless explicitly requested.

## Desktop smoke checklist

Before RC:

1. Desktop launches.
2. It can read existing ADM state.
3. Workspace and Environment list/inspect work.
4. Capability report is visible, including unavailable/degraded reasons.
5. MCP definitions, selection and health/status are manageable or clearly documented as CLI/Gateway fallback.
6. Skill sources, selection and availability are manageable or clearly documented as CLI/Gateway fallback.
7. Process/run/verifier state is visible enough for daily use.
8. Private Memory values and secret-backed MCP values are not displayed in normal views.
9. Known limitations are documented.

## RC blockers

A gap blocks RC only if it prevents safe daily use of validated Core capabilities through Desktop or makes the CI/artifact path unreliable.

Default RC blockers:

- Desktop cannot launch.
- Desktop cannot load ADM state.
- Desktop leaks private Memory values or secret-backed MCP values.
- Desktop bypasses writer lease, executable allowlist or Environment selection rules.
- Capability diagnostics are absent from Desktop with no clear CLI/Gateway fallback.
- CI cannot run test/vet/build or cannot produce any usable binary artifact.

## Non-blocking by default

The following are useful but not RC blockers unless dogfood proves otherwise:

- GitNexus/provider integration.
- Additional Phase 14 investigation helpers.
- Temporary resource lifecycle cleanup.
- Installer/tray/autostart/updater/signing/notifications.
- Automatic release publishing.

## Known limitations policy

A known limitation is acceptable for RC when:

1. the Core capability remains available through CLI/Gateway;
2. Desktop clearly avoids unsafe or misleading behavior;
3. the limitation is documented in the RC notes;
4. there is a post-RC follow-up path.
