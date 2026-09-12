# 18-04 Summary — Integrated Desktop acceptance progress

Date: 2026-09-12
Status: integrated automated/browser/full Go/vet/Wails-build verification complete; exact visible native artifact acceptance pending because an existing user Desktop instance owns the single-instance window.

## Tested source

Tested source head before this documentation checkpoint: `dd4821c docs: record Phase 18-03 verification`.

18-04 performed final integrated acceptance against the Phase 18 Desktop code produced by 18-01/18-02/18-03. No new feature code was added during 18-04. No evidence-backed UI regression was found in automated/browser/Go/vet/build checks, so Task 3 has no code fix.

## Commands and evidence

### JavaScript / browser

Initial planned command note:

- `node --test tests/desktop-ui` — NOT USED as the final test command because this Node/Windows invocation treated the directory path as a module and returned `MODULE_NOT_FOUND` for `tests\\desktop-ui`.
- Actual discovered test set:
  - `tests/desktop-ui/dashboard.test.cjs`
  - `tests/desktop-ui/navigation.test.cjs`
  - `tests/desktop-ui/project-pages.test.cjs`
  - `tests/desktop-ui/runtime-view.test.cjs`
  - `tests/desktop-ui/skill-bulk.test.cjs`

Final JS/helper command:

- `node --test tests/desktop-ui/dashboard.test.cjs tests/desktop-ui/navigation.test.cjs tests/desktop-ui/project-pages.test.cjs tests/desktop-ui/runtime-view.test.cjs tests/desktop-ui/skill-bulk.test.cjs` — PASS, 20/20.

Final production-asset browser smoke:

- `node tests/desktop-ui/browser-smoke.cjs` — PASS, 122 checks at each size/scale:
  - 1120x760 @ 1.0
  - 820x560 @ 1.0
  - 1120x760 @ 1.25

The browser fixture covers the integrated Desktop shell and management surfaces: connection/profile startup, routing, Workspace/Environment lists, detail Summary/Diagnostics, MCP editor/import/probe freshness, Skill/Sources subviews and filtered bulk actions, explicit Memory scopes, Runtime subviews/output, system grouping, modal guards, stale response rejection and negative no-implicit-action checks.

### Go and vet

- `go test -count=1 ./...` — PASS via ADM async Run `run_3f517e281969f9c4`, exit 0.
  - Slowest observed package: `internal/gateway` at `103.843s`.
- `go vet ./...` — PASS via ADM async Run `run_a9b2cf42319f474e`, exit 0.

### Wails final artifact

- Exact final Wails build — PASS via ADM async Run `run_635cc447e2592061`:
  - `powershell -NoProfile -File scripts/build-desktop.ps1 -OutputName adm-desktop-phase18-final-windows-amd64.exe`
- Artifact:
  - `D:\projects\ai-dev-manager-v2\dist\adm-desktop-phase18-final-windows-amd64.exe`
  - size: `17,356,800` bytes
  - SHA-256: `78881B25AAB6C744532AC65CDF884006AB50040A0755482C80EC3965CF0D233F`
- Non-disruptive second-instance launch — PASS/PENDING distinction:
  - Started PID `36828` from the exact final artifact path.
  - The process exited within 5 seconds with exit code `0` and no window title because an existing user Desktop instance already owned the single-instance window.
  - This proves non-disruptive launch/hand-off behavior, not visible exact-artifact UI acceptance.

## D01-D10 status

- D01 PASS — connection/profile/overview and partial-load behavior covered by integrated browser fixture and full Go tests.
- D02 PASS — all ten routes, metric shortcuts, Workspace-filtered Environment list and Diagnostics entry covered by browser fixture. Routing/tab/filter changes do not trigger implicit probe/mutation/value reads.
- D03 PASS — ordinary non-Git project/Environment assumptions and metadata-safe lifecycle are covered by existing Core tests and integrated Desktop fixture call boundaries. No new Git prerequisite was introduced.
- D04 PASS — MCP editor/import/preview/apply/probe and Skill/Sources/default/current-Environment interactions covered by browser fixture. Stale preview cannot apply; unobserved is not healthy; no route-triggered scan/read/probe occurs.
- D05 PASS — Global and Environment-private Memory scopes covered by browser fixture. Closing Environment detail cannot retarget the Memory page; Environment-private reads are explicit and current-scope-bound.
- D06 PASS — Runtime Verifiers/Processes/Runs subviews and bounded output identity covered by browser fixture plus full Go tests. Tabs do not auto-run actions; output remains current-owner/resource-scoped.
- D07 PASS — modal/focus/draft/transition guards and delayed-response rejection covered by browser fixture.
- D08 PASS — Diagnostics renders only existing InspectEnvironment facts/unresolved IDs and honest unknown/source/reason/freshness. Diagnostics switching performs no adapter call.
- D09 PASS — browser fixture ran at 1120x760, 820x560 and 125% scale with long path/list data and visible focus/modal checks.
- D10 AUTOMATED PASS / EXACT VISIBLE NATIVE PENDING — final Wails build and second-instance launch passed, but the exact final artifact was not visibly exercised because the active user Desktop single-instance window is already running. I did not close or kill that window.

## Phase 18 completion decision

Phase 18 is not marked complete in this checkpoint because D10 requires visible native acceptance of the exact final Wails artifact. Automated/browser/full Go/vet/build evidence is green, but native single-instance ownership prevents honestly claiming that `adm-desktop-phase18-final-windows-amd64.exe` was the visible tested window.

Next action: when the active Desktop window can be closed or replaced safely, launch `D:\projects\ai-dev-manager-v2\dist\adm-desktop-phase18-final-windows-amd64.exe` visibly and exercise the D10 native route/dialog/profile/tray checks. If that passes, create `18-CLOSEOUT.md` and mark Phase 18 complete. If it fails, record the exact failure and fix only the demonstrated regression.


## User feedback polish: context and diagnostics

Status: committed; see latest local Git log for the exact commit hash.

The active Desktop screenshot showed three issues: the ADM connection page retained the global Management Context block where Environment selection is irrelevant, the grouped ADM connection layout was visually broken by old flex styling, and the sidebar `环境诊断` entry still routed through the Environments/detail flow instead of being a standalone diagnostics page.

Fixes made:

- Management Context is now visible only on Environment-relevant routes: Environments, Runtime, MCP, Skills, Memory and Diagnostics. Overview, ADM connection, Exec allowlist and Settings hide it.
- `环境诊断` is now a real `#/diagnostics` route. Without a selected Environment it shows an explicit no-fan-out message; with a selected Management Environment it renders the existing `InspectEnvironment` payload and does not open the Environment detail modal.
- The ADM connection page layout is flattened to a single-column card flow for profile, endpoints and local lifecycle. The duplicate `gatewayState` DOM id was removed so the header status remains the single authoritative status badge.

Evidence after this polish:

- JS syntax PASS: `navigation.js`, `app.js`.
- Helper regression PASS: 20/20.
- Production Chromium smoke PASS: 131 checks x 3 viewports/scales.
- Focused Go PASS: `./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/management`.
- Wails build with the original final output name compiled but could not copy over `dist/adm-desktop-phase18-final-windows-amd64.exe` because that exact binary is currently running and locks the file. Run: `run_a0342c4e84b10bab`; result: command_failed at Copy-Item only.
- Wails polish artifact PASS: `run_401fb4c5bf51b1d3`; `dist/adm-desktop-phase18-polish-windows-amd64.exe`; 17,362,432 bytes; SHA-256 `4D3ACB64B781447A49543D4C68AC0989884C01763ABB237B7F91123FFA74DEB8`.
- Polish artifact second-instance smoke: PID 36676 exited 0 because the running final Desktop instance still owns the single-instance window.

D10 remains PENDING for exact visible native acceptance. The next concrete action is to close the current running Desktop window/process that owns `dist/adm-desktop-phase18-final-windows-amd64.exe`, rebuild/copy the final-named artifact, launch it, and perform the visible native route/layout checks against that exact binary.
