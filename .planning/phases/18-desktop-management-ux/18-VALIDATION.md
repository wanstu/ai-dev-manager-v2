# Phase 18 — Validation and Acceptance

Date: 2026-09-11
Current scope: 18-01
Current result: planning only; every implementation/browser/Wails gate below is NOT RUN.

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

1. Run focused tests while each task changes its own behavior. New JS tests must exercise route/view-state logic and deferred results, not mirror strings from the implementation.
2. Exercise browser behaviors A01-A08 against production assets, including hidden/focus behavior that Node logic tests cannot prove.
3. At the final code node, run:
   - `node --test tests/desktop-ui/navigation.test.cjs tests/desktop-ui/dashboard.test.cjs`
   - `go test -count=1 ./cmd/ai-dev-manager-desktop ./internal/desktop ./internal/management ./internal/adminmcp`
   - `go test -count=1 ./...`
   - `go vet ./...`
   - `powershell -NoProfile -File scripts/build-desktop.ps1 -OutputName adm-desktop-phase18-01-windows-amd64.exe`
4. Launch that exact Wails artifact and record a short click-through: disconnected start; connected overview; each section; Workspace/Environment modal; MCP import/editor failure; Skill filters; explicit Memory read; Runtime view; profile switch; Settings/tray availability; minimum-size layout.
5. Check `git diff --check` before staging and `git diff --cached --check` before each commit. Inspect staged names; stage only the current plan's files.

Long tests/builds use ADM asynchronous run_start/status with bounded output. Record run IDs and terminal exit codes; an outer tool timeout does not prove command failure and must not cause blind repeated launches.

Full-suite, vet and Wails build need not be repeated for documentation-only updates after the verified code node. Broaden testing only for a concrete remaining risk. A raw `go build` launch is not Wails artifact acceptance.

## Evidence / closeout record

Implementation session must record:

| Item | Current planning status |
|---|---|
| Source baseline / scope checkpoint | 0173259 / d95b088 |
| Implementation commits | None |
| A01-A08 behavior checks | NOT RUN |
| Focused/full Go tests and vet for new UI | NOT RUN |
| New JS test files/gate | Planned; files not created |
| Wails build / artifact SHA-256 | NOT RUN |
| Real Wails manual acceptance | NOT RUN |
| Known limitations | Current source review only; UI has not been changed |
| Next action | Start 18-01 Task 1 after Git/context/writer preflight |

Write actual results to `18-01-SUMMARY.md` during implementation and update this table. If native GUI access is unavailable, save the implemented code and automated evidence with manual acceptance explicitly pending; do not mark the plan/phase fully complete.

## Planning-only checks

This session reviews the source/requirement references, route and file boundaries, task dependencies, negative acceptance, continuation instructions and consistency of STATE/PROJECT/ROADMAP/PHASE-MAP. Final documentation integrity and staged diff checks are recorded at the planning commit boundary. No application test/build result is inferred from a documentation-only change.

Planning integrity evidence (2026-09-11): PASS — all 8 planning documents are present/readable; 9 declared execution-plan requirement IDs resolve to PRODUCT_CONTRACT; all 10 routes and 9 acceptance cases are documented; the numbered plan inventory is 30; current status is 18-01 ready / implementation not started. Source/planning file references were checked. Worktree `git diff --check` passed; staged diff checks are required immediately before the final planning commit. The scope checkpoint is `d95b088`.
