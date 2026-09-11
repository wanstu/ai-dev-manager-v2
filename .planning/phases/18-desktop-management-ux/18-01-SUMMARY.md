# 18-01 Implementation / Validation Summary

Date: 2026-09-11
Plan: `18-01-PLAN.md`
Status: IMPLEMENTATION COMPLETE / NATIVE MANUAL ACCEPTANCE PENDING
Starting head: `3b26a92 docs: complete Phase 18 UI execution planning`

## Task 1 checkpoint — management shell and section routing

Implementation is complete enough for the first local code node; final browser/Wails acceptance remains pending under Task 3.

Changed/added code:

- `cmd/ai-dev-manager-desktop/frontend/index.html` — compact persistent connection header, grouped sidebar, persistent Management Environment context, ten route pages and Settings placement. Existing management controls/forms remain single-instance.
- `cmd/ai-dev-manager-desktop/frontend/navigation.js` — fixed route normalization/activation, hash/history handling, active navigation state and navigation guard integration. No adapter access or persistence.
- `cmd/ai-dev-manager-desktop/frontend/styles.css` — shell/sidebar/page/hidden/focus/responsive rules for the existing 1120x760 / 820x560 window envelope.
- `cmd/ai-dev-manager-desktop/frontend/app.js` — initializes the navigation helper, prevents page changes while a connection transition/editor/Environment detail is active, and restores focus after closing Environment detail.
- `cmd/ai-dev-manager-desktop/main_test.go` — embedded asset/shell marker coverage only; not treated as behavior proof.
- `tests/desktop-ui/navigation.test.cjs` — narrow route list/normalization/hash behavior tests.

Preserved boundaries:

- No Core, adapter, schema, persistence, Runtime authorization or build-pipeline change.
- Routing makes no management bridge call and does not mutate the selected Management Environment.
- Native editor dialogs and Environment detail remain outside hidden route pages.
- `connectionSelect`, `managementEnvironment` and every pre-existing action/form ID remain single-instance.
- Diagnostics remains an Environment-route shortcut rather than a new route/API.

Supporting checks run at this checkpoint:

- `node --check cmd/ai-dev-manager-desktop/frontend/navigation.js` — PASS.
- `node --check cmd/ai-dev-manager-desktop/frontend/app.js` — PASS.
- `node --test tests/desktop-ui/navigation.test.cjs` — PASS, 2/2 tests.
- HTML structural check — PASS: ten expected route pages, 150 IDs, zero duplicate IDs, overview is the only initially visible page, one connection selector and one Management Environment selector.
- `go test -count=1 ./cmd/ai-dev-manager-desktop` — PASS.

Not yet claimed:

- A01/A02/A06 real-browser behavior, focus/layout at 1120x760 and 820x560, back/forward with live DOM, and native Wails smoke are PENDING Task 3.
- Task 2 truthful dashboard/load-state/error isolation and stale-response protection are NOT STARTED.

## Task 2 checkpoint — truthful overview and scoped auxiliary reads

Implemented after Task 1 commit `6aa21bf`:

- Added `dashboard.js` as a presentation-only state model for unloaded/loading/success/stale/error and sanitized snapshot counts. Only a trusted successful snapshot establishes numeric zero; unloaded/error show `—`; loading/stale may display a clearly labeled prior successful snapshot.
- Removed the synthetic empty snapshot from `clearManagementData()`. Disconnected/unloaded management lists and counts now remain explicitly unavailable instead of looking like a successful empty installation.
- `GetSnapshot` is now the authoritative base read and renders before auxiliary reads. `ListSkillSources`, selected Environment inspection/Skill availability, Runtime verifiers/processes/runs and open Environment detail refresh cannot erase a valid base snapshot when they fail.
- Skill source, Environment context and the three Runtime lists expose local unavailable reasons. Runtime reads use independent settled results so one unavailable list does not clear the other valid lists.
- Added connection/Environment/detail generation guards. Environment switching clears old Runtime observation/output immediately and late reads from a previous connection/Environment/detail are discarded.
- Open Environment detail refresh failure is treated as an auxiliary failure and no longer bubbles into the refresh-button path that clears connection-bound state.
- Added `dashboard.test.cjs`, including unloaded-vs-zero, stale/loading retention and scope-generation rejection checks.

Task 2 checks:

- `node --check cmd/ai-dev-manager-desktop/frontend/app.js` — PASS.
- `node --check cmd/ai-dev-manager-desktop/frontend/dashboard.js` — PASS.
- `node --test tests/desktop-ui/navigation.test.cjs tests/desktop-ui/dashboard.test.cjs` — PASS, 6/6 tests.
- Focused `go test -count=1 ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/management ./internal/adminmcp` — PASS after moving the old `global_memory_count` embed marker assertion from `app.js` to its new `dashboard.js` owner.
- Structural data-flow check — PASS: zero duplicate DOM IDs, base snapshot no longer coupled to `ListSkillSources`, no synthetic empty snapshot on disconnect, auxiliary Environment/Runtime reads use settled results and scope generations.

## Task 3 checkpoint — browser/integration/artifact validation

Implementation and automated acceptance are complete. Native visible-window click-through remains pending, so 18-01 is not marked accepted and 18-02 must not start yet.

Commits:

- `6aa21bf feat(desktop): add management section navigation`
- `d243bb0 feat(desktop): add management overview and scoped loading states`
- `2cc2acd test(desktop): add real-browser phase18 smoke`

Real-browser production-asset evidence:

- `tests/desktop-ui/browser-smoke.cjs` reads the shipped `index.html` and absolute production CSS/JS assets, injects only a fake Wails adapter plus assertions, and runs in an installed Chromium browser without adding npm/Vite/Playwright/Puppeteer/jsdom dependencies.
- `node tests/desktop-ui/browser-smoke.cjs` — PASS at `1120x760`, `820x560`, and `1120x760 @ 125%`; 37 checks per scenario.
- The browser gate covers all ten routes, one visible/active page, navigation with zero adapter calls, explicit Global Memory read only, editor/detail route guards, failed Workspace-save draft preservation, focus restoration, long-path/no-horizontal-overflow behavior, Skill-source auxiliary failure isolation, Runtime-list partial failure isolation, profile-switch old-state clearing and absence of implicit probe/mutation calls.
- `node --check tests/desktop-ui/browser-smoke.cjs` — PASS.
- `node --test tests/desktop-ui/navigation.test.cjs tests/desktop-ui/dashboard.test.cjs` — PASS, 6/6 tests. Note: Node 22 on this Windows host treats `node --test tests/desktop-ui` as a module path instead of directory discovery, so the two test files are enumerated explicitly.

Repository/integration gates:

- Focused Desktop/management Go gate — PASS: `go test -count=1 ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/management ./internal/adminmcp`.
- Full repository gate — PASS via ADM async Run `run_2d1e9b0528fda01a`: `go test -count=1 ./...`, exit 0.
- Vet — PASS via ADM async Run `run_027f4d0395e337cc`: `go vet ./...`, exit 0.
- Wails build — PASS via ADM async Run `run_27752bef55f6d9c6`: `powershell -NoProfile -File scripts/build-desktop.ps1 -OutputName adm-desktop-phase18-01-windows-amd64.exe`.
- Exact artifact: `D:\projects\ai-dev-manager-v2\dist\adm-desktop-phase18-01-windows-amd64.exe`, 17,238,016 bytes, SHA-256 `81E297856EECAA6317FF5A5C0BF3084AFA755FB42FCB7F44F22E2086236F229C`.
- An older Desktop instance (`adm-desktop-v1.0.0-rc.local-runobs-windows-amd64.exe`) is currently running and owns the product's single-instance lock. Launching the exact Phase 18 artifact therefore exercised the expected second-instance handoff path and exited within 5 seconds with exit code 0. The active old Desktop/Gateway was deliberately not stopped merely for UI acceptance.

Acceptance status:

- A01/A03/A04/A08: automated browser/unit evidence PASS.
- A02/A05/A06: production-browser evidence PASS for routing/action reachability, profile clearing, modal/detail guard, failed-save draft retention, focus and minimum-size/scaling behavior. A05 late-scope rejection also has narrow generation tests; native manual observation remains part of the final acceptance boundary.
- A07: existing focused/full Go boundary tests remain green; browser mutation log confirms navigation does not acquire writers, probe MCPs, refresh sources or execute Runtime actions implicitly.
- A09: embedded assets, full Go/vet and exact Wails build PASS; exact artifact starts successfully into the existing single-instance handoff. **Visible new-artifact Wails window click-through/screenshots are NOT RUN** because the already-running user Desktop owns the single-instance lock and was not interrupted.

Remaining acceptance blocker:

- Run the exact Phase 18 artifact in a visible Wails window when the existing Desktop can be closed without disrupting active work. Click through overview and all ten routes, connection/Environment context, editor cancel/failure behavior, explicit Memory load, Runtime page and minimum-size layout. Record native evidence before accepting 18-01.
- Until that native check is recorded, keep `18-01` at implementation-complete/manual-acceptance-pending and do not execute 18-02.

Next action: native Wails click-through of the exact artifact when the single-instance lock is available; then record acceptance/closeout and only after that unblock 18-02. Do not push.
