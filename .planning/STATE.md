---
gsd_state_version: 1.0
milestone: V2
current_phase: 9
current_phase_name: Planner / Executor / Reviewer Contract
status: phase-review
stopped_at: Phase 09 locally verified on b73bb59; integration review pending; Phase 10 not started
last_updated: "2026-09-08T05:46:38Z"
last_activity: 2026-09-08
last_activity_desc: Phase 09 deterministic workflow contract locally verified; integration review pending
state_head: b73bb5905c36076bad11305a775498fbdf7f62f7
progress:
  total_phases: 13
  completed_phases: 8
  total_plans: 10
  completed_plans: 10
  percent: 62
---

# Project State

## Project Reference

See: `.planning/PROJECT.md` (rebaselined 2026-09-06)

**Core value:** Give external Agents one reliable, inspectable, safe local development control plane instead of disconnected CRUD/config surfaces.
**Current focus:** Phase 9 — Planner / Executor / Reviewer Contract

## Current Position

Phase: 9 — Planner / Executor / Reviewer Contract
Plan: 09-01 complete; locally verified
Status: Phase 9 locally verified on isolated feature worktree; integration review pending; Phase 10 not started
Last activity: 2026-09-08 — implementation `b73bb59` passed deterministic workflow, real Streamable HTTP, race, full Gateway, full Go, vet and diff-check gates

Progress: 8/13 phases integrated complete (62%); Phase 9 plan 1/1 complete and locally verified

## Phase 1 Completion Evidence

- Real host Skill root: `C:\Users\wanstu\.config\opencode\skills`
- Explicit support root: `C:\Users\wanstu\.config\opencode\gsd-core`
- `environment_skill_list` / `environment_skill_read` are Environment-gated Gateway capabilities.
- Enabled/disabled isolation, shared installation, path escape rejection, broken-artifact isolation, and real Streamable HTTP GSD acceptance pass.
- Final `go test ./...`, `go vet ./...`, and `git diff --check` passed.
- GSD `uat.classify-coverage` reported all four deliverables auto-covered; `phase uat-passed 1 --require-verification` returned `passed: true` with no blockers.
- GSD `query phase.complete 1` advanced the project to Phase 2 after the installed GSD successfully parsed the normalized ROADMAP/STATE.

## Phase 2 Completion Evidence

- Structured verifier definitions are Environment-scoped and persist with stable `vf_` identity.
- `environment_verifier_list` and writer-gated `environment_verifier_run` are real Agent Gateway capabilities backed by the existing allowlisted Runtime execution boundary.
- Pass/fail/timeout return bounded structured verifier results; cwd containment, executable authority, writer heartbeat, and zero-config behavior are covered by tests.
- Real Streamable HTTP acceptance runs a configured verifier and V2 `go test ./...` from a non-Git repository copy without recursive acceptance re-entry.
- Final implementation commit: `f7490f3` (`feat(verifier): add structured verifier runtime`).
- Canonical GSD re-verification on `f7490f3` passed 4/4 Phase goals and all 16 Nyquist points with no human UAT required.
- GSD `query phase.complete 02` completed 1/1 plans with no warnings and advanced the project to Phase 3.

## Phase 3 Completion Evidence

- MCP catalog definitions now carry Streamable HTTP transport and unresolved header references; legacy empty transport normalizes to `streamable-http`.
- Environment selection gates on-demand health, external tool listing, and external tool calls.
- Health distinguishes exactly `configured`, `disabled`, `healthy`, and `error`; configured state is never treated as health.
- Endpoint/header environment references resolve only at activation boundaries and resolved secret values are not exposed in normal health/error output.
- Real Streamable HTTP acceptance proves enable → healthy → list tools → call tool → disable → disabled → broken upstream → structured error, while unrelated Gateway tools remain available.
- CLI `mcp status --id MCP_ID --environment-id ENV_ID`, management, and Desktop adapter consume the same app-owned health semantics.
- Final local validation passed: `go test ./...`, `go vet ./...`, `git diff --check`, and focused MCP lifecycle acceptance.
- Per explicit user instruction, Phase 3 final verification/UAT used deterministic local evidence only and did not invoke any OpenCode/GSD/other LLM subagent. Local `gsd-tools` accepted `verification.status=passed`, 5/5 automated UAT coverage, and `phase complete 03` advanced the project to Phase 4 with zero warnings.

## Phase 4 Completion Evidence

- Standalone non-Git D:/projects/adm-verifier-report was developed only through ADM: failing structured verifier, exact implementation edit, passing verifier, actual red/green JSON consumed by the built CLI (exit 1/0).
- Real gsd-next and support workflow consumed through Environment gating; deterministic smart-entry and plan consistency used, no LLM agents.
- Wrong writer, disabled Skill, zero-verifier file access, local Git failure, private Memory isolation and owned Gateway restart persistence passed.
- Reproduced Windows timeout-child cleanup blocker twice; new parent/child cancellation test failed before the fix and passed afterward. Bounded OS process-tree cancellation fixes the unchanged HTTP acceptance; no timeout inflation or cleanup retry added.
- Final go test ./... and go vet ./... passed. Evidence: phases/04-external-agent-dogfood-gate/04-VERIFICATION.md.
- R1 milestone review passed and Phase 04 was fast-forwarded to master. Post-integration go test ./..., go vet ./..., and git diff --check passed.

## Phase 5 Completion Evidence

- Existing HTTP/stdio Gateway is the persistent runtime owner; owner ID/PID/start time and live external MCP sessions are process-instance state only.
- Canonical app activation resolution keeps persisted Environment/catalog desired state separate from resolved secret-bearing activation and observed owner/session health.
- Two independent Agent clients share one owner/session path; restart creates a fresh owner and rebuilds from persisted desired selection.
- Dead upstream returns structured `connection_refused`, evicts the owned session, and is never carried across restart as stale healthy state.
- Current Gateway stop uses owner-bound graceful HTTP shutdown and closes owned sessions; older/no-owner Gateway health retains the safe legacy termination fallback.
- Independent current-source dogfood used private ADM_V2_HOME plus loopback 43761/43762; shared 41137 was untouched. Full tests/vet/diff-check and five repeated cross-process restart tests passed.
- Evidence: `phases/05-persistent-runtime-ownership/05-VERIFICATION.md`, `05-UAT.md`, and `evidence/`.
- Integration review passed; `master` fast-forwarded `b36ef6a -> 79b2122`. Post-merge Gateway and remaining package tests, `go vet ./...`, and `git diff --check` passed. Evidence: `05-INTEGRATION-REVIEW.md`.

## Phase 6 Completion Evidence

- The existing persistent Gateway owner now owns long-running Environment-scoped dev processes by stable `proc_` identity; launching Agent/client exit does not terminate an HTTP-Gateway-owned process.
- Start/stop reuse writer authority, the Runtime executable allowlist, Environment-relative cwd containment and the existing OS process-tree cancellation path.
- Later clients can list/status processes and read bounded stdout/stderr tails; exited process logs remain owner-local and queryable until owner shutdown.
- Windows dogfood status reports listening TCP ports only for ADM-owned PIDs; the Agent API never accepts arbitrary PID control. Non-Windows currently returns no port facts rather than widening authority.
- Explicit stop, Environment removal and graceful Gateway shutdown deterministically stop owned process trees and release ports. Gateway restart creates a fresh owner and does not resurrect prior `proc_` observations.
- Independent current-source dogfood used a private ADM_V2_HOME, non-Git project, detached Gateway 127.0.0.1:36728 and separate probe processes. Real HTTP service, bounded logs, owned ports, explicit stop, owner cleanup, restart-empty state and no observed-state persistence all passed.
- Integration review gate on `75e919d`: Gateway tests 17.872s; real restart acceptance ×5 5.745s; focused tail/cross-client/writer-lease tests 0.993s. All 16 Go packages passed tests or compiled through bounded groups; vet and diff checks passed. Review fixed exact-capacity truncation reporting and made port assertions use decoded fields/platform expectations. Evidence: `phases/06-dev-process-logs-ports/06-INTEGRATION-REVIEW.md`, `06-VERIFICATION.md`, and `evidence/regression.json`.
- The user explicitly authorized local master integration on 2026-09-08. `master` fast-forwarded `0ad5488 -> 6dd87d6` with no conflicts. Post-integration Gateway tests passed in 15.004s; all remaining Go packages passed tests or compiled, and vet/diff-check passed. Evidence: `phases/06-dev-process-logs-ports/evidence/integration.json`. No push or automatic Phase 7 advance was performed.

## Phase 7 Completion Evidence

- Git worktree remains an optional isolation capability around Environment. Ordinary `environment_create` is still Workspace-contained and the existing non-Git Gateway development acceptance passes independently.
- Managed worktrees persist separate `wt_` metadata, use generated `adm/wt_...` branches and are created only under the ADM state-directory worktree root; Agent callers do not choose arbitrary destination paths or managed branch names.
- Real Git tests create two worktree Environments from one source Workspace and prove distinct roots/branches, cross-root file isolation, and unchanged source branch/HEAD/tracked content.
- `app.Service.Runtime` revalidates managed Environment relation, owned-root containment, Git top-level, common-dir and branch identity before routed access. Branch tamper and unmanaged out-of-Workspace roots fail before attempted file mutation.
- Destroy requires the matching writer. Dirty or locally advanced/unpublished work is refused without explicit force; clean and forced destroy both retain the generated branch, and committed work remains reachable after forced worktree removal.
- Real MCP Streamable HTTP acceptance creates/lists/destroys a managed worktree through the Agent Gateway, verifies generic `environment_remove` cannot bypass safety, and confirms a failed non-Git worktree request does not break ordinary Environment creation.
- Final local gate on implementation `4b11375`: named real-Git acceptance 7.820s; app tamper acceptance 1.179s; real HTTP Gateway acceptance 1.505s; non-Git Gateway regression 0.123s; focused packages pass; `go test ./...`, `go vet ./...`, and `git diff --check` pass. Evidence: `phases/07-optional-git-worktree-isolation/07-VERIFICATION.md`, `07-UAT.md`, and `evidence/regression.json`.
- Integration review passed; repeat real-Git isolation ×3 and real HTTP Gateway lifecycle ×3 passed. Following explicit user authorization, local `master` fast-forwarded `9ae56ad -> 43d413f`; post-integration `go test ./...`, `go vet ./...`, and `git diff --check` passed. Evidence: `07-INTEGRATION-REVIEW.md` and `evidence/integration.json`. No push was performed.

## Phase 8 Completion Evidence

- The existing persistent Gateway owner now owns single-command asynchronous Agent Runs by stable `run_` identity; the launching MCP client may disconnect while later clients list/status/cancel the same running Run.
- `run_start` requires the matching Environment writer and resolves the existing app Runtime before asynchronous execution, preserving executable allowlist, cwd containment, managed-worktree validation, bounded output, timeout and OS process-tree cancellation instead of adding a second command policy.
- Lifecycle distinguishes `running`, `succeeded`, `failed` and `canceled`. Tests prove exit 0, exit 7, timeout, explicit cancellation and a 64-byte result bound against a 4096-byte helper output.
- Wrong-writer cancel, forbidden executable and escaped cwd are rejected locally. Failed authority checks do not install a Run. Running Runs heartbeat the existing writer lease.
- Owner Environment drop and owner/Gateway close cancel and wait for active Run process trees. Real loopback helpers prove the serving port is released on cancel and graceful owner shutdown.
- Real Streamable HTTP subprocess acceptance proves client independence and restart semantics: a Run is left active, Gateway stop cleans it, restart yields a different owner, `run_list` is empty, the prior `run_` ID is invalid, and ordinary Environment `read` still works.
- Run identity/result/count observations are not persisted in `state.json`; restart does not infer, resume or resurrect prior Runs.
- Final local gate on implementation `a74b318`: focused Run tests 4.072s; focused race gate 9.746s; real HTTP restart acceptance ×3 2.966s; full Gateway suite 43.320s; final `go test ./...`, `go vet ./...`, and `git diff --check` passed. Evidence: `phases/08-agent-run-lifecycle/08-VERIFICATION.md`, `08-UAT.md`, and `evidence/regression.json`.
- Integration review passed with no blocking findings: cancel/drop/owner-close acceptance ×5 passed in 5.716s; real Streamable HTTP shutdown/restart acceptance ×5 passed in 5.372s; `git diff --check master...HEAD` passed and scope audit found no Phase 9-11 implementation. Evidence: `08-INTEGRATION-REVIEW.md`.
- Following explicit user continuation authorization, local `master` fast-forwarded `24f1df9 -> 4f73d4d`. Post-integration Gateway tests passed in 50.456s; final `go test ./...`, `go vet ./...`, and `git diff --check` passed. No push was performed.

## Phase 9 Local Verification Evidence

- `run_workflow_start` creates a workflow on the existing persistent Gateway-owned `run_` lifecycle; ordinary `run_list`, `run_status`, and `run_cancel` remain the observation/cancellation surface.
- Planner behavior is deterministic and auditable: explicit goal/ordered steps/verifier are materialized before Run installation, stable `step_01...` IDs are assigned, and every executor command is preflighted through existing Runtime authority.
- Executor runs sequentially on the owner-derived Run context and records per-step state, timing, exit code and bounded stdout/stderr. Unsafe authority failure installs no partial Run.
- Reviewer reuses `app.Service.RunVerifier` and retains structured verifier evidence. Normal verifier failure is `run=succeeded`, workflow outcome `rejected`, review `rejected`; it is not an orchestration failure.
- Executor non-zero and reviewer invocation/configuration errors are separately classified as `executor_step_failed` and `reviewer_error` with Run state `failed`.
- Workflow cancellation reuses Phase 8 writer-gated Run cancellation, writer heartbeat and owner process-tree cleanup. Workflow observations remain owner-local and are not persisted to `state.json`.
- Real Streamable HTTP acceptance proves accepted/rejected workflow audit is visible from a later client. Repeated acceptance ×5 passed in 0.914s; focused race gate passed in 12.293s; fresh full Gateway suite passed in 33.774s; `go test ./...`, `go vet ./...`, and `git diff --check` passed.
- Evidence: `phases/09-planner-executor-reviewer-contract/09-VERIFICATION.md`, `09-UAT.md`, and `evidence/regression.json`.

## Accumulated Context

### Decisions

- If @pjadm is unusable, @pj may be used to repair it; switch back to @pjadm after repair (user-authorized 2026-09-07). No fallback was needed in Phase 4.

- Skill runtime uses explicit discovery roots and real `SKILL.md` artifacts; no arbitrary disk scan.
- Real GSD requires an explicit sibling `gsd-core` support root; support-file access is limited to configured roots.
- Environment selection remains the Skill authorization gate.
- Git remains optional for Environment.
- Verifier precedes Agent/GSD orchestration.
- Structured verifier definitions stay Environment-scoped; verifier execution reuses the existing writer-gated, allowlisted `Runtime.Exec` boundary rather than adding a second execution primitive.
- Phase 8 Agent Run is deliberately one asynchronous Runtime command owned by the persistent Gateway; Phase 9 extends the same `run_` lifecycle with deterministic workflow audit rather than creating a second orchestration owner/persistence model.
- Phase 9 treats normal verifier failure as a review decision (`run=succeeded`, workflow/review `rejected`), while executor/reviewer infrastructure failures are distinct failed Run classifications. Phase 10 must preserve this distinction when deciding whether GSD state may advance.
- Desktop/package expansion remains frozen until the Core milestones reach human-manager parity.
- Do not invoke OpenCode/GSD/other LLM subagents or consume separate provider/model quota unless the user explicitly reverses this decision; use `@pjadm` local file/Git/Go/gsd-tools capabilities for continued work.

### Blockers/Concerns

Phase 9 is locally verified on `feat/planner-executor-reviewer-contract` at implementation `b73bb59` with no known blocker. Integration review is pending; local `master` is unchanged by Phase 9, push remains unauthorized, and Phase 10 has not started. No LLM subagents are permitted.

## Deferred Items

- Connector/runtime tool exposure drift: shared 41137 and installed connector expose an older subset than current source; final Phase 4 acceptance used a freshly built private Gateway. Refresh deployed Gateway/connector separately; do not claim shared deployment was upgraded.
- pjadm investigation-acceleration candidates (user proposal 2026-09-08): later evaluate whether to add evidence-first higher-level tools for HTTP endpoint resolution, symbol/reference/write tracing, response-field/data lineage, symbol-scoped Git history/diff, concrete test-DB metric explanation, API debug-SQL reverse mapping, and semantic-consistency detection. Prefer machine evidence chains with `confidence`, `evidence[]`, `uncertainties[]` (and alternatives where applicable) over opaque guesses. Keep deferred from active roadmap phases unless the capability becomes a concrete blocker or is explicitly prioritized.

- automatic Memory context composition
- full MCP server configuration model after the Phase 3 Streamable HTTP tracer is complete: validate server names (`[A-Za-z0-9._-]+`), optional description/comment, transports `stdio` / `sse` / `streamable-http` / `openapi`, authentication modes `none` / header token / OAuth, and explicit env/header configuration; use MCPHub's separation of transport/auth/config concerns as a design reference rather than copying its runtime model
- Desktop build UX: ordinary `go build ./cmd/ai-dev-manager-desktop` currently produces a binary that reports missing Wails build tags at runtime; make the supported release build path unmistakable or fail earlier without changing the current Phase 3 desktop freeze
- GitNexus tooling integration: `pjadm.exec` currently exposes only allowlisted executables and does not provide a directly invokable `gitnexus` executable; attempts to use GitNexus indirectly through Node/npx can start but `gitnexus analyze` has failed with `ERR_MODULE_NOT_FOUND: Cannot find package 'tree-sitter-swift'` from GitNexus's Swift ingestion module, preventing repository indexing and therefore impact analysis / `detect_changes`. Some repositories have also reported `gitnexus executable not found`, while `.gitnexus/run.cjs --help` has not produced a reliable usable path. Until fixed, do not claim GitNexus checks passed; fall back explicitly to source/call-site search, diff review, tests, and real runtime acceptance where required.
- installer/tray/autostart/updater/signing/notifications
- UI redesign
- migration/compatibility burden

## Session Continuity

Last session: 2026-09-08T05:46:38Z
Stopped at: Phase 09 locally verified on `b73bb59`; integration review pending; Phase 10 not started
Resume file: `.planning/phases/09-planner-executor-reviewer-contract/09-VERIFICATION.md`; perform independent Phase 9 integration review. Do not merge master, push, or start Phase 10 without the review decision/authorization.
