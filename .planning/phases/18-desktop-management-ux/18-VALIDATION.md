# Phase 18 — Validation and Acceptance

Date: 2026-09-11
Current scope: full Phase 18 delivery (18-01 -> 18-02 -> 18-03 -> 18-04); 18-01 is accepted and 18-02 is the next executable plan
Current result: 18-01 A01-A09 PASS. 18-02 B01-B07 PASS and B08 automated/Wails-build PASS; exact visible native UI acceptance for the 18-02 artifact is pending because an existing user Desktop single-instance window was not closed. 18-04 remains the mandatory final integrated gate.

## Acceptance matrix

IDs below are local test/checklist identifiers; they are not new product requirements.

| ID | Requirement basis | Scenario / expected behavior | Evidence required |
|---|---|---|---|
| A01 | ADM-DESKTOP-001 | Default/unknown hash resolves to overview; all ten menu routes, metric shortcuts and back/forward work; one visible page and aria-current item. | Route logic tests plus real browser against production assets. |
| A02 | ADM-DESKTOP-001 / ADM-CORE-015 | Every existing management action family in the design's migration table is reachable after movement; no duplicate DOM IDs, detached handlers or hidden modal ancestors. Navigation alone makes zero adapter calls. | DOM/interaction assertion and action log from fake Wails bridge; Wails click-through. |
| A03 | ADM-MGMT-001 / ADM-CORE-017 | Disconnected/loading/error show unavailable counters; only successful empty snapshot shows zero; populated counts match source; snapshot/overview contains no private Memory values. | Dashboard state tests, deliberately empty/malformed/error responses and browser checks. |
| A04 | ADM-CORE-003 / 004 / 020 | Fail ListSkillSources, Environment availability or one Runtime list independently; a valid base snapshot and unrelated management areas remain usable; show local reason/retry and do not report overall success. | Controlled rejected bridge promises with per-area assertions. |
| A05 | ADM-MGMT-001 / ADM-CORE-013 / 017 | Switch profile A -> B and Environment A -> B while reads are pending. Prior response never populates the new context. Old Memory, output, detail and observations clear; request serialization remains intact. | Deferred-promise tests and browser profile/Environment sequence with distinct sentinel values; no sensitive production data. |
| A06 | ADM-DESKTOP-001 | Keyboard-only navigation, visible focus, hidden sections not focusable; editor validation/save failure retains draft; Escape/cancel restores visible opener; modal/hash/connection-transition guard prevents hidden drafts. Long labels/paths fit at 1120x760 and 820x560 plus 125% scaling. | Real browser focus/layout checks and actual Wails screenshots/click notes. |
| A07 | ADM-CORE-003 / 005 / 007 / 012 | No Git, verifier, selected MCP/Skill or writer does not block shell/global management; only the operation requiring the capability is unavailable. Existing Runtime authorization is preserved; routing never acquires writer, enables capabilities or adds executables. | Non-Git disposable fixture plus fake bridge mutation log; existing Go boundary tests. |
| A08 | ADM-CORE-013 / 017, PROC-01..03, ARUN-01 | Visiting overview/Memory/Runtime does not load Memory values, probe MCPs, refresh Skill sources, execute verifier, start/stop/cancel or fetch arbitrary logs. Explicit existing buttons still perform their intended operation once. | Positive and negative adapter call assertions, scoped Memory and Runtime interaction checks. |
| A09 | ADM-DESKTOP-001 / ADM-MGMT-001 | New JS assets are actually embedded/loaded; production still uses NewClientAdapter, Admin MCP and existing profile/tray/build behavior; built Wails app launches and routes successfully. | Focused/full Go tests, vet, JS gate, existing Wails build, tested artifact path/hash and live launch evidence. |

## Plan-level acceptance index

The A-cases above establish the shared shell/state safety contract in 18-01. Successor plans add focused cases without replacing that contract:

| Plan | Local acceptance IDs | Focus | Final carry-forward rule |
|---|---|---|---|
| 18-01 | A01-A09 | routing, truthful dashboard/load state, optional-failure isolation, scope transitions, focus/layout and production Wails path | Must pass before 18-02; affected cases rerun after shared shell/state changes. |
| 18-02 | B01-B08 | Workspace/Environment lifecycle presentation, stable-target detail/actions and existing Runtime lists/output | Must pass before 18-03; carry only evidence unchanged by later shared-state edits. |
| 18-03 | C01-C08 | MCP/Skill global-vs-Environment hierarchy, explicit Memory scopes, existing system controls and diagnostics | Must pass before 18-04; Memory/scope evidence is invalidated by later related changes. |
| 18-04 | D01-D10 | integrated cross-section journeys, degradation, keyboard/scaling and exact final Wails artifact | Mandatory final Phase 18 gate; cannot be waived because A/B/C passed separately. |

The detailed positive and negative cases live in each numbered plan. `18-04-SUMMARY.md` now records final automated evidence and the exact tested source/artifact; `18-CLOSEOUT.md` remains pending until D10 exact visible native acceptance passes.

## Required fixtures / test isolation

Use the real frontend assets with a fake `window.go.desktop.Adapter` installed before deferred scripts run. Test fixtures remain under `tests/desktop-ui/`, outside the embedded production frontend; do not duplicate the production HTML/JS or replace browser behavior with a handmade DOM model.

Fixtures should cover:

- no selected profile; connected empty installation; populated snapshot;
- two different connection profiles and two Environment IDs with distinguishable harmless data;
- ordinary non-Git Environment with no MCP/Skill/verifier selection;
- failed sources/availability/Runtime reads with a successful base snapshot;
- delayed old-scope read resolution, delayed profile transition and failed explicit save;
- long Windows paths/IDs/names and a long MCP/Skill list;
- a Memory sentinel absent until the explicit read operation, then cleared on scope change;
- bridge call counters identifying all mutations, probes, source refreshes and value/log reads.

For native smoke, use a dedicated test ADM connection and disposable fixture records. Do not stop the Gateway serving the active development session or mutate unrelated real user catalogs/Memory merely to test UI wiring. Read-only native checks can use an existing authorized connection; any management mutation evidence should identify the disposable resource.

## Execution gate order

1. Execute 18-01 and pass A01-A09. Run its focused route/dashboard tests, affected Go gates, real-browser production-asset checks and its named Wails smoke before closing the plan.
2. Execute 18-02 only after 18-01 acceptance. Run the focused project/runtime tests and B01-B08 browser/native checks defined there; preserve the shared A contract and rerun any A case affected by shared shell/state changes.
3. Execute 18-03 only after 18-02 acceptance. Run focused MCP/Skill/Memory/system tests and C01-C08 checks; explicitly recheck A/B scope/focus behavior affected by shared coordinator, modal or reset changes.
4. Execute 18-04 against the final integrated code. D01-D10 are mandatory, including real browser production assets, final full repository tests/vet and build/launch of the exact named Wails artifact. Evidence-backed fixes invalidate and rerun the affected checks.
5. At every code/document commit node, run `git diff --check`; after explicit staging run `git diff --cached --check`, inspect staged names and keep the commit scoped to the current task/plan.

New JS tests must exercise behavior, state transitions and deferred results rather than mirror implementation strings. Long tests/builds use ADM asynchronous run_start/status with bounded output; record run IDs and terminal exit codes. An outer tool timeout does not prove command failure and must not cause blind repeated launches.

The full suite, vet and Wails build need not be repeated for documentation-only updates after an unchanged verified code node, but 18-04 still performs its final integrated gates. A raw `go build` launch is not Wails artifact acceptance.

## Evidence / closeout record

Implementation session must record:

| Item | Current planning status |
|---|---|
| Source baseline / scope checkpoint | `0173259` / `d95b088` |
| Detailed planning commits | `93a74d9` (design/18-01), `b6da7fc` (18-02), `f52b872` (18-03), `3b26a92` (18-04 + synchronized planning) |
| Implementation commits | `6aa21bf` navigation shell; `d243bb0` truthful overview/scoped states; `2cc2acd` production-browser smoke harness |
| A01-A09 / B01-B08 / C01-C08 | A01-A09 PASS for 18-01. 18-02 B01-B07 PASS and B08 automated PASS from helper/browser/focused Go gates; exact Wails build PASS and second-instance launch exited 0, but exact visible native 18-02 UI checks remain pending under the active user Desktop single-instance window. 18-03 C01-C07 PASS and C08 automated PASS from helper/browser/focused Go/Wails gates; exact visible native 18-03 UI checks also remain pending under the same active single-instance window. |
| D01-D10 integrated acceptance | D01-D09 PASS from final integrated helper/browser/full Go/vet evidence; D10 automated build and second-instance launch PASS but exact visible native artifact acceptance PENDING due active Desktop single-instance window |
| Focused/full Go tests and vet for new UI | PASS. Focused desktop/management gate PASS; full `go test -count=1 ./...` via `run_2d1e9b0528fda01a` exit 0; `go vet ./...` via `run_027f4d0395e337cc` exit 0. |
| JS/browser tests | PASS. Route/dashboard tests 6/6; production-asset Chromium smoke 37 checks each at 1120x760, 820x560 and 1120x760@125%. |
| Wails build / artifact SHA-256 | PASS via `run_27752bef55f6d9c6`; `dist/adm-desktop-phase18-01-windows-amd64.exe`, 17,238,016 bytes, SHA-256 `81E297856EECAA6317FF5A5C0BF3084AFA755FB42FCB7F44F22E2086236F229C`. |
| Real Wails native acceptance | PASS. Exact artifact opened visibly in its own Wails/WebView2 window after the old Desktop parent had exited naturally. UI Automation invoked all ten routes plus diagnostics, verified route-heading focus and persistent context, opened/closed Workspace modal with focus return, and exercised the native 820x560 minimum. Evidence: `evidence/18-01-native-ui-automation.json` plus three PNG screenshots. |
| Known limitations | Native acceptance intentionally did not display or record real Memory values; explicit Memory value-load behavior is covered by the production-browser fake bridge. Synthetic OS `Alt+Left` was not accepted as native history evidence; browser back/forward is covered by the production-browser gate. No implementation defect is currently known for 18-01. |
| Next action | Commit 18-01 closeout/evidence, calibrate 18-02 against the delivered shell, then execute 18-02 Task 1 under the existing Phase 18 boundaries. Do not push. |

Write actual results to each plan's `18-0N-SUMMARY.md` during implementation and aggregate final evidence in `18-04-SUMMARY.md` / `18-CLOSEOUT.md`. If native GUI access is unavailable at an intermediate slice, preserve pending evidence accurately; required final native evidence cannot be silently waived to mark Phase 18 complete.

## Planning-only checks

This session reviews the source/requirement references, route and file boundaries, task dependencies, negative acceptance, continuation instructions and consistency of STATE/PROJECT/ROADMAP/PHASE-MAP. Final documentation integrity and staged diff checks are recorded at the planning commit boundary. No application test/build result is inferred from a documentation-only change.

Planning integrity evidence (2026-09-11): PASS — all 8 Phase 18 planning documents are present/readable; four numbered Phase 18 plans exist; their 21 unique declared requirement IDs all resolve to `docs/PRODUCT_CONTRACT.md`; all 10 routes are preserved and A01-A09 / B01-B08 / C01-C08 / D01-D10 define 35 local acceptance cases. The numbered plan inventory for phases 01-18 is 33 when phase-00 bootstrap is excluded. Planning closeout landed at `3b26a92`; current execution status is 18-01 implementation/automated validation complete with native visible-window acceptance pending, while 18-02-04 remain blocked/sequential. The scope checkpoint is `d95b088`; detailed plan commits are `93a74d9`, `b6da7fc`, `f52b872` and `3b26a92`.


## 2026-09-12 user feedback polish evidence

- Management Context route visibility and standalone Diagnostics route covered by production Chromium smoke: 131 checks x 3.
- ADM connection layout flattening and unique `gatewayState` id covered by production Chromium smoke and embedded asset marker tests.
- Focused Go PASS after marker updates: `./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/management`.
- Final-named artifact rebuild could not overwrite the currently running final exe; polish-named artifact build PASS via `run_401fb4c5bf51b1d3`, SHA-256 `4D3ACB64B781447A49543D4C68AC0989884C01763ABB237B7F91123FFA74DEB8`.
