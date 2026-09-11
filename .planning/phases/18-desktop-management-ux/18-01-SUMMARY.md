# 18-01 Implementation / Validation Summary

Date: 2026-09-11
Plan: `18-01-PLAN.md`
Status: IN PROGRESS
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

Still pending Task 3:

- Production-asset browser interaction for A01-A08, including live focus/hidden-page/back-forward and deferred-response behavior.
- Full repository test/vet gates, exact Wails build and native launch/click-through.
- 18-01 cannot be marked complete until required native/browser evidence exists; if this environment cannot provide GUI interaction, preserve manual acceptance as pending and do not begin 18-02.

Next action after the Task 2 commit: execute Task 3 automated/integration gates, determine what real-browser/native evidence can be produced in this environment, and write the final 18-01 evidence accurately. Do not start 18-02.
