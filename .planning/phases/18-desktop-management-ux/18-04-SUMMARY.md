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
