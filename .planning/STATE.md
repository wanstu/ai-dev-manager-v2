---
gsd_state_version: 1.0
milestone: V2
current_phase: 17
current_phase_name: Distribution Only If Needed
status: v1.0.0-release-ready
stopped_at: ADM v1.0.0 release gate passed locally; closeout ready for tag and publish
last_updated: "2026-09-11T12:23:00Z"
last_activity: 2026-09-11
last_activity_desc: Prepared ADM v1.0.0 release closeout after local gates, artifact build and live async Run observability validation
state_head: da5d2ea
progress:
  total_phases: 17
  completed_phases: 17
  total_plans: 27
  completed_plans: 27
  percent: 100
---

# Project State

## Project Reference

See `.planning/PROJECT.md`, `.planning/ROADMAP.md`, and `.planning/rebaseline/2026-09-08-core-boundary.md`.

**Core value:** Give external Agents one reliable, inspectable and safe local development control plane for MCP, Skill and project Runtime capabilities.

**Explicit boundary:** ADM supplies capabilities, lifecycle, authorization and diagnostics. Agent/GSD supplies task planning/orchestration.

## Current Position

Phase: 17 — Distribution Only If Needed
status: phase-17-standby
Base master: `9fc4c4b`
Active development branch: `master`
Phase 11 importer implementation: `ea0d85af99f4591a431ba22dd8f1df036c76cac6`
Phase 12 source-aware Skill implementation: `4074d3e8e5349ee717bf63ad027391d14579cecc`
Phase 12 Skill availability implementation: `8f46f6e019afc8722708e3a4478de06faa4b241a`
Phase 13 static capability report implementation: `b0eb0c74e6f43844b3a754feaa2d34c93bb82da6`
Phase 13 Gateway-owner capability report implementation: `46c86f3`
Phase 14 endpoint evidence resolver implementation: `0e40394`
Phase 14/15/16/17 priority planning: `9c4c539`
Phase 16 GitHub CI/RC baseline: `6d51e6d`
Phase 16 RC baseline closeout: `5a96fe4`
Phase 16 Desktop management adapter APIs: `374fa12`
Phase 16 Desktop MCP/Skill visual management UI: `74be256`
Phase 16 Agent/Admin MCP surface split: `2c6ab1c`
Phase 16 Desktop ADM connection profiles: `cd191ff`
Phase 16 Desktop Admin MCP convergence: `0dfe01d`
Phase 16 CLI Admin MCP convergence: `ecf0516`
Phase 16 Windows default Gateway port fix: `65396ef`
Phase 16 local Gateway bind preflight: `28bf1ac`
Phase 16 Desktop Runtime visibility: `19808f7`
Phase 16 RC1 dogfood blocker fixes: `8dbe817`
Phase 16 MCP/Skill Desktop UI dogfood polish: `a3478f2`
Phase 16 generic mcpServers import fix: `c52d584`
Phase 16 Windows path canonicalization: `11c49ee`
Phase 16 Desktop tray/autostart polish: `53300d8`
Phase 16 saved connections/modal editors/icons: `41a161a`
Phase 16 packaging/docs/tray sizing polish: `abafcc0`
Phase 16 Windows canonical path CI fix: `73c4481`
Phase 16 tray event loop fix: `5c8574e`
Phase 16 release automation: `559aa3d`
Phase 15 lifecycle cleanup implementation: `48b0858`
Phase 15 closeout validation: `04d1bba`
Phase 17 local distribution standby validation: `9fc4c4b`
Run observability mitigation: `6b793d3`
Phase 17 run-observability artifact refresh: `da5d2ea`
ADM v1.0.0 pre-release validation head: `da5d2ea`

The prior `feat/gsd-phase-executor` branch is abandoned and must not merge. Its planning/provenance/state-advance implementation is not ADM product scope.

## Retained Core History

### Phases 1-4 — Development foundations

- real `SKILL.md` discovery/read from explicit roots with support-root containment and Environment gating;
- structured optional verifier runtime;
- real external MCP Streamable HTTP tool discovery/call with Environment gating, four-state health and activation-boundary secret resolution;
- external Agent dogfood through ADM-only local development capabilities.

### Phases 5-8 — Persistent local Runtime

- persistent Gateway owner and external MCP session restart reconciliation;
- long-running dev process lifecycle, bounded logs and owned-port facts;
- optional managed Git worktree isolation;
- generic single-command asynchronous `run_` start/list/status/cancel lifecycle.

These capabilities remain ADM Core.

## Superseded Work

### Phase 9 — Planner / Executor / Reviewer

Phase 9 was technically implemented, reviewed and integrated to local master, but the product rebaseline determined that its workflow orchestration semantics belong to the Agent/GSD/orchestrator layer.

The Git/evidence history remains for auditability. The Agent-facing workflow surface and workflow-specific domain code are cleanup scope under BOUNDARY-01. Future ADM capabilities must not depend on FLOW-01.

### Abandoned Phase 10 GSD executor branch

`feat/gsd-phase-executor` is not to be merged. `gsd_phase_inspect`, `gsd_phase_start`, `gsd_phase_advance`, `.planning` provenance interpretation and STATE advancement are explicitly outside ADM Core.

## Active Requirements

### Phase 10 — Boundary cleanup (complete)

- BOUNDARY-01 ✅: Planner/Executor/Reviewer workflow surface/domain is removed from ADM Core while generic single-command Runs remain.
- BOUNDARY-02 ✅: no GSD `.planning` interpretation/state-advance API was merged or introduced.
- BOUNDARY-03 ✅: files/exec, verifier, process, MCP, Skill, managed worktree, ordinary Environment and generic Run regression gates passed.

### Phase 11 — MCP runtime completion (complete)

- 11-01 ✅: typed desired MCP configuration, activation-only references and real Streamable HTTP/stdio owner-bound transport runtime.
- 11-02 ✅: owner-local health/recovery observation, configurable probe/reconnect policy, inventory/inspect/refresh, stable-ID update invalidation and safe no-replay semantics; full test/vet/race/diff gates passed at `4c4f4dc`.
- 11-03 ✅: external JSON/JSONC preview/apply import adapters, batch atomicity, conflict policy, credential-reference conversion, Gateway/management/CLI surfaces and runtime-owner definition-fingerprint/generation reconciliation; full `go test -count=1 ./...` plus vet/race/diff gates passed at `ea0d85a`.

### Phase 12 — Skill runtime completion (complete)

- 12-01 ✅: persisted Skill sources, source/artifact stable identity, atomic per-source refresh, safe failed-refresh preservation, unresolved selection preservation and Gateway/management/CLI source surfaces; full test/vet/race/diff gates passed at `4074d3e`.
- 12-02 ✅: Environment-specific Skill availability, bounded artifact/support inventory, structured Skill read diagnostics and broken-Skill isolation; full test/vet/race/diff gates passed at `8f46f6e`.

### Phase 13 — Environment capability diagnostics (complete)

- 13-01 ✅: shared `CapabilityReport`/`CapabilityFact`/`CapabilityEvidence` model and resilient side-effect-free application-level Environment report with static file/exec/verifier/Git/isolation/MCP/Skill/process/run facts; full test/vet/race/diff gates passed at `b0eb0c7`.
- 13-02 ✅: Gateway-owner observation enrichment and canonical Agent-facing `environment_capability_report`, plus CLI/management wrappers; full test/vet/race/diff gates passed at `46c86f3`.

### Phase 14 — Evidence-first Investigation Toolkit (slice complete / later helpers deferred)

- 14-01 ✅: endpoint evidence resolution for URL/path plus optional HTTP method; returns bounded static route evidence, confidence and uncertainties; full test/vet/diff gates passed at `0e40394`.
- 14-02 ✅: optional code intelligence provider + GitNexus integration boundary uses existing Environment MCP/runtime authorization, passive Gateway-owner observations, explicit provenance/freshness/uncertainty and static fallback; focused app/Gateway plus full test/vet/diff gates passed before the later lifecycle closeout.
- Additional evidence-first slices remain deferred by default and do not block Phase 15.

### Phase 15 — Temporary Resource Lifecycle (complete)

- 15-01 ✅: explicit durable/temporary retention metadata, cleanup inspect/dry-run/execute, mark/promote, Skill/MCP safe cleanup and retention lifecycle surfaces landed at `fc8303d`.
- 15-02 ✅: temporary Environment state-only cleanup, Gateway-owner managed worktree cleanup through existing destroy safety, and bounded Environment capability inspection to avoid unbounded disabled catalog facts landed at `48b0858`; rebuilt Gateway live dogfood confirmed `skill.catalog` summary with `suppressed_disabled_fact_count` instead of per-disabled-Skill facts.

### Phase 16 — Desktop Core Parity + 1.0 RC Readiness (complete)

- 16-01 ✅: Desktop parity matrix, RC gate and GitHub Actions CI build baseline completed at `6d51e6d`; closeout recorded at `5a96fe4`.
- 16-02 ✅: Desktop MCP/Skill visual management UI completed through adapter API commit `374fa12` and UI commit `74be256`; full repository test/vet/build/diff gates passed. Human Wails click-through continues as RC1 dogfood.
- 16-03A ✅: Agent/Admin MCP surface split completed at `2c6ab1c`. `/mcp` now excludes Admin-only management tools; `/admin/mcp` exposes the privileged management superset over the same Core/runtime owner.
- 16-03B ✅: Desktop connection profile + configurable health/liveness completed at `cd191ff`; Base URL derives health/Agent/Admin MCP endpoints and local bootstrap is loopback-only.
- 16-03C ✅: Production Desktop normal management moved to Admin MCP at `0dfe01d`; disconnected/stopped ADM clears the management backend and never silently falls back to writable local state.
- 16-03D ✅: Normal CLI `workspace`/`environment`/`exec`/`mcp`/`skill`/`memory` management converged on the shared Admin MCP client at `ecf0516`; `--adm-url` / `ADM_V2_URL` select the management target and connection failure never falls back to writable local state. `gateway`/`doctor`/`state` remain explicit local bootstrap/offline/recovery paths.
- 16-03E ⏸️: Authenticated/TLS-safe non-loopback remote Admin MCP is intentionally deferred post-RC.
- RC1 ⚠️: dogfood exposed that the published Desktop artifact was incorrectly produced by raw `go build` and therefore failed Wails build-tag validation at launch; it also exposed that Gateway lifecycle did not honor `--adm-url` / `ADM_V2_URL`. RC1 is superseded by RC2.
- RC2 ⚠️: Wails build/CI, Gateway target convergence and screenshot-driven MCP/Skill UI polish were fixed, but dogfood later exposed generic top-level `mcpServers` auto import as unnecessarily ambiguous; RC2 is superseded by RC3.
- RC3 ✅: canonical `generic-mcpservers` auto import landed at `c52d584`; literal imported env/header values remain reference-only, Desktop preview shows generated reference requirements, Admin MCP acceptance passed, exact Wails Desktop launch smoke passed, and full `go test -count=1 ./...` / vet / diff gates passed.
- 16-04 ✅ local code/gates: Windows short/long path canonicalization landed at `11c49ee`; Desktop UI/tray/autostart polish landed at `53300d8`; full local tests/vet/diff/Wails build passed. Tray menu and real login acceptance remain manual.
- 16-05 ✅ local code/gates: saved ADM connection profiles, modal child editors and ADM application/tray/window icon pipeline landed at `41a161a`; full tests/vet/diff/Wails build plus hidden/single-instance smoke passed. Visual/modal and tray/autostart click-through remain manual.
- 16-06 DONE local code/gates: ime-lock-v2 tray lifecycle alignment, fitted tray icon sizing, simplified `adm` / `adm-desktop` user-facing names, unified `dist/` packaging, tag-aware GitHub artifacts and README/docs split landed at `abafcc0`; full tests/vet/diff/default Wails build and RC packaging smoke passed.
- 16-07 DONE closeout: Windows CI short/long path test compatibility landed at `73c4481`; tray left/right click dispatch was fixed by keeping the Win32 tray message loop on one OS thread at `5c8574e`; tag-triggered GitHub Release automation landed at `559aa3d`; user accepted the tray behavior and closed Phase 16.

### Phase 17 — Distribution Only If Needed (post-RC)

- Installer/updater/signing/notifications remain post-RC. Tray/autostart were pulled into Phase 16 by dogfood requirements and are implemented locally; interactive acceptance remains open.

### Next Core priorities

1. Address the dogfood timeout candidate in `.planning/dogfood/2026-09-11-pjadm-timeouts.md`: long verification should use an explicit async/inspectable path or return an actionable synchronous diagnostic instead of disappearing into outer tool timeout.
2. Keep Phase 16 closed; do not reopen Desktop/CI/release work unless a new concrete blocker appears.
3. Preserve Phase 15 cleanup safety in future changes: active writers/processes/runs, dirty or unpublished managed worktrees, unknown ownership and insufficient evidence must continue to block cleanup.
4. Phase 17 is on conditional standby after local artifact refresh; do not implement installer/updater/signing/notifications unless future manual dogfood proves a concrete blocker.
5. Future development should be dogfood-driven: continue deferred evidence-first helpers only when generic Runtime/search is insufficient for a real task.

## Product Decisions

- MCP and Skill are core product capabilities and must stay visible in the main roadmap.
- A foundation vertical slice is not the same as completion. Phase 3 does not mean MCP is finished; Phase 1 does not mean Skill is finished.
- `run_` is retained only as generic asynchronous Runtime lifecycle; it must not grow task semantics.
- Worktree remains an optional isolation primitive. ADM does not orchestrate parallel Agents or choose integration policy.
- Verifier returns structured evidence; the external Agent/orchestrator decides what that evidence means for its task/phase.
- No LLM/GSD/OpenCode subagent quota is used unless the user explicitly reverses the existing instruction.
- No automatic push at review boundaries.
- During active development, continue on `master`; make a commit at each clear node.
- Phase 11 transport scope is Streamable HTTP + stdio; legacy HTTP+SSE is not added. Stdio MCP executable launch must still obey ADM's executable allowlist.
- Phase 11 health-check interval/probe timeout are configurable per MCP. Automatic reconnect defaults to off; when enabled it retries at one fixed configured interval with no exponential/adaptive backoff. Health/recovery observation is owner-local, and recovery never replays failed tool calls.
- Phase 11 supports preview + atomic single/batch JSON/JSONC import adapters for OpenCode, WorkBuddy/CodeBuddy, Codex plugin MCP JSON, Claude Code and supported MCPHub shapes. Name conflicts default to error. Import writes only global MCP definitions, never existing Environment selections, and converts literal credential-bearing values into secret/environment-reference requirements instead of persisting them.
- Phase 11 does not persist interactive OAuth tokens until an approved secure credential lifecycle exists; supported HTTP auth is none or explicit secret-backed headers.
- Phase 12 uses persisted Skill sources, source/artifact-based stable identity and explicit atomic refresh; ADM does not interpret Skill instructions.
- Phase 12 availability is Environment-specific; a broken Skill is local to that Skill and does not poison unrelated Skill, MCP, file, verifier or Runtime operations.
- Phase 13 capability inspection is side-effect-free by default and aggregates per-capability facts instead of probing/executing optional tools.
- Phase 13 Gateway-owner enrichment reads existing owner-local observations only; it does not reconnect, Ping, refresh inventory, call MCP tools, run verifiers, start processes, or acquire writer leases.
- Phase 14 investigation helpers must return concrete evidence, confidence and uncertainties. 14-01 endpoint investigation is static literal/dynamic-candidate search only and does not execute code or call endpoints.
- GitNexus is an optional evidence provider through existing Environment MCP/runtime boundaries; ADM does not embed a mandatory code graph engine, and passive provider inspection must remain side-effect-free with static fallback.
- Phase 15 temporary resource lifecycle/retention is complete; CLI/UI-created resources remain durable by default, and temporary cleanup must remain explicit, inspectable and conservative.
- Phase 16 is RC-first. GitHub Actions CI auto build is RC infrastructure, not post-RC polish.
- Desktop must remain a management surface over Core, with no Desktop-only state or authorization model.
- Normal management uses the separate Admin MCP management plane. Agent/Admin MCP privilege separation is implemented; production Desktop and normal CLI management have converged on the shared Admin MCP client. Direct writable state-file access is not a peer normal mode; retain it only as explicit offline/bootstrap/recovery behavior, preferably read-only until cross-process locking and service-stopped safety are designed.
- Desktop connection configuration should support explicit scheme/host/domain/port health checks. Non-loopback remote Admin MCP must remain disabled until authentication/TLS/Host-boundary semantics are defined.
- Phase 16 first UI priority is MCP/Skill visual management: configure/import MCPs, manage Skill sources, enable/disable both for one Environment and see health/availability reasons from Core diagnostics.
- Phase 17 distribution polish should not block local 1.0 RC unless daily use proves it is necessary.

## Deferred

- automatic Memory context composition;
- additional evidence-first investigation slices after 14-02;
- broad Desktop feature expansion beyond RC blockers;
- installer/updater/signing/notifications;
- migration/compatibility burden.

## Session Continuity

Phase 16 remains closed and must not be reopened for broad Desktop/CI/release work. Phase 15 is complete: 15-01 landed at `fc8303d`; 15-02 landed at `48b0858` with temporary Environment state-only cleanup, Gateway-owner managed worktree cleanup and bounded Environment capability inspection.

Current dogfood note: a rebuilt Gateway on `127.0.0.1:43138` live-validated bounded `environment inspect` output. The report now emits `mcp.catalog` and `skill.catalog` summary facts plus selected MCP details; it suppresses unselected disabled Skill detail facts and showed `catalog_count=207`, `selected_count=0`, `suppressed_disabled_fact_count=207` for this Environment.

Latest dogfood mitigation: repeated long synchronous `pjadm.exec` outer timeouts are tracked in `.planning/dogfood/2026-09-11-pjadm-timeouts.md`. Async `run_start` remains the intended path for long verification, and code at `6b793d3` makes `run_status` expose bounded stdout/stderr snapshots while a command is still running, including truncation flags when `max_output_bytes` is reached.

Live dogfood status: the currently connected ChatGPT `pjadm` Gateway was rechecked after the run-observability mitigation and `run_status` returned stdout while a command was still `running`. This confirms the active Gateway path now includes the `6b793d3` fix.

Release closeout: `v1.0.0` local release gates passed, Windows CLI/Desktop artifacts were built with checksums, and `.planning/releases/v1.0.0.md` records the release decision and retained post-1.0 limits. Next action is to commit this closeout record, tag `v1.0.0`, and push master plus the tag as requested.