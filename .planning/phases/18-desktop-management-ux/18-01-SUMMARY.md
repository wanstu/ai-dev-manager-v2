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

Next action after the Task 1 commit: begin Task 2 by separating management data load state from synthetic empty data and isolating GetSnapshot from auxiliary reads, then add dashboard state tests. Do not start 18-02.
