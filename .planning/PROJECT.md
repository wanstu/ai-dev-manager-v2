# AI Dev Manager V2 — Project Baseline

## Product Mission

AI Dev Manager V2 is a local AI Coding Environment / development control plane for external Agents.

Its job is to give an Agent a reliable, inspectable, safe development environment over one stable Gateway: registered local roots, Environment-scoped context, file/runtime capabilities, Skills, external MCPs, verification, optional isolation, and later long-running task orchestration.

ADM is not primarily a Desktop app and is not primarily a CRUD manager. Human UI exists to manage the same Core that Agents use.

## Product Priority

1. Core/runtime behavior is real and reliable.
2. External Agents can use the capabilities through ADM without fallback tools.
3. Persistent lifecycle and optional isolation support real development work.
4. Agent/GSD orchestration builds only on validated runtime capabilities.
5. Desktop/CLI expose the completed Core; UI polish and packaging come last.

## Locked Architecture Rules

- Workspace is a registered local directory. Git is optional.
- Environment is a development context rooted at a directory. It must not intrinsically require Git/worktree/Docker/verifier.
- Missing optional capability blocks only the operation that needs it.
- One physical root has at most one active writer.
- Agent-facing operations use stable ADM identity and never trust arbitrary filesystem paths as authorization.
- Control/management and Runtime execution are separate responsibilities.
- One long-lived ADM Gateway should serve many Environments; users must not configure a new MCP server for every task.
- Skill and MCP must be runtime capabilities, not merely persisted records.
- GSD Skill installation is global/shared; project `.planning/` is project state and is not the Skill installation itself.
- Important planning decisions live under `.planning/`, not only in chat history.
- A Phase implements only its plan. New ideas outside the Phase go to backlog unless they invalidate a locked requirement or block the Phase exit criteria.
- Do not mark a capability complete because CRUD or metadata exists. Completion requires a real consumption path and acceptance test.
- No pre-stable migration/compatibility burden unless explicitly requested.

## Reality Audit — 2026-09-08

### Validated / usable

- Workspace registration and lifecycle for ordinary directories.
- Environment creation/lifecycle with arbitrary Workspace-contained roots.
- tree/read/search/write/exact-edit/single-file-delete.
- single-writer lease by physical root.
- allowlisted structured command execution.
- optional git status/diff/branch on a Git Environment root.
- persistent local state.
- one HTTP/stdio ADM MCP Gateway and Gateway lifecycle controls.
- explicit Global Memory and Environment-private Memory CRUD.
- CLI and Desktop management surfaces over the same state/services.
- Desktop real-state dogfood fixes: modal Environment detail, transient operation feedback, hidden Windows command consoles.
- Real Skill runtime (Phase 1): explicit discovery roots, real `SKILL.md` artifacts, stable IDs, explicit support roots, Environment-gated list/read, shared global installation, and real-host GSD Gateway acceptance.
- Installed GSD bootstrap (Phase 1): this Agent read `gsd-next` and its `gsd-core` smart-entry workflow through ADM, then real `gsd-tools` parsed and advanced the repository planning state.
- Structured verifier runtime (Phase 2): Environment-scoped test/lint/build/custom definitions, Gateway list/run tools, writer-gated execution through the existing allowlisted Runtime, bounded structured pass/fail/timeout results, cwd containment, zero-config development, and real Streamable HTTP non-Git `go test ./...` acceptance.
- Persistent Runtime owner (Phase 5, integrated): the existing HTTP/stdio Gateway owns live external MCP sessions, exposes one process-instance owner identity, reconciles from persisted desired Environment/catalog state after restart, rejects stale healthy observations, and closes owned resources on disable/removal/graceful shutdown without introducing a second daemon.
- Dev process lifecycle (Phase 6, integrated): the persistent Gateway can own allowlisted Environment-scoped `proc_` resources across Agent client exit, retain bounded stdout/stderr tails, report owned listening TCP ports on the Windows dogfood target, and deterministically stop process trees without exposing arbitrary PID control or persisting observed process state.
- Optional Git worktree isolation (Phase 7, integrated): one Git Workspace can create ADM-owned managed worktree Environments with generated `wt_` identity/branch/root, routed Runtime revalidation, source-checkout preservation and writer-exclusive safe destroy. Dirty/unpublished work is refused by default; force retains the branch. Ordinary non-Git Environment behavior remains unchanged.
- Agent Run lifecycle (Phase 8, locally verified; integration review pending): the persistent Gateway owns single-command asynchronous `run_` resources across Agent client exit, with stable list/status/cancel identity, writer-gated cancellation, succeeded/failed/canceled states, bounded command results, owner cleanup, and no restart resurrection or persisted Run observation.

### Scope and remaining work

#### External MCP — validated Phase 3 tracer

Phase 3 completed Streamable HTTP connection definitions, Environment gating, four-state health, activation-boundary secret resolution and real list/call acceptance. Phase 5 is integrated to local master with Gateway-owned persistent external MCP sessions and restart reconciliation. Broader transports/auth remain deferred.

#### Memory

Persistence and explicit scoped reads/writes are real. Automatic context injection is not currently required or implemented; do not describe Memory as automatically supplied to every Agent request.

#### Desktop

Desktop is a functional management shell, not the product completion gate. It must remain frozen except for blockers until the Core milestones below are complete.

### Not implemented

- cross-platform listening-port observation parity beyond the current Windows dogfood target;
- Planner/Executor/Reviewer orchestration;
- GSD phase execution/verified state advance in V2;
- parallel Agent/worktree orchestration;
- production installer/tray/autostart/updater/signing.

## R1 Dogfood — reviewed and integrated

A real external Agent completed a separate non-Git verifier-report CLI through ADM only, consuming real GSD Skill artifacts, structured red/green verification and private Memory. Restart persistence and operation-local negative acceptance passed. A reproducible Windows command-child cancellation blocker was fixed. R1 milestone review passed on 2026-09-07 and Phase 04 was fast-forwarded into master; post-integration tests/vet/diff-check passed. Phase 5 was later explicitly authorized, locally verified, and fast-forwarded into local master after integration review; Desktop expansion remains frozen. Evidence: `.planning/phases/04-external-agent-dogfood-gate/04-R1-REVIEW.md` and `.planning/phases/05-persistent-runtime-ownership/05-INTEGRATION-REVIEW.md`.

## Validated Requirements

### SKILL-RUNTIME — Phase 1

- **SKILL-RUN-01** ✅: Skills are discovered from explicit configured roots or explicit definitions, not by arbitrary disk scan.
- **SKILL-RUN-02** ✅: A Skill resolves to a real artifact/content source, not only a display name.
- **SKILL-RUN-03** ✅: Environment selection gates Agent access to a Skill.
- **SKILL-RUN-04** ✅: One global GSD Skill can be used by multiple Environments without copying it into each project.
- **SKILL-RUN-05** ✅: An Agent connected only through ADM can discover and read/use the GSD Skill for an enabled Environment; a disabled Environment cannot.

Evidence: `.planning/phases/01-skill-runtime/01-VERIFICATION.md` and `01-UAT.md`.

### VERIFY — Phase 2

- **VERIFY-01** ✅: Verification is an optional Runtime capability.
- **VERIFY-02** ✅: Environment/project development can declare structured test/lint/build/custom verifiers without requiring Git.
- **VERIFY-03** ✅: Verifier execution returns structured status, exit code, bounded output, timeout/failure identity, and is available through the Agent Gateway.
- **VERIFY-04** ✅: Absence of verifier configuration does not block normal Environment/file development.

Evidence: `.planning/phases/02-structured-verifier-runtime/02-VERIFICATION.md` (machine-verifiable; no human UAT required).

### MCP-RUNTIME — Phase 3

- **MCP-RUN-01** ✅: Configured MCP definitions represent actual connection information and transport.
- **MCP-RUN-02** ✅: Environment selection gates activation/access at runtime.
- **MCP-RUN-03** ✅: Health/status distinguishes configured, disabled, healthy, and error; configuration existence is not reported as healthy.
- **MCP-RUN-04** ✅: Real external MCP discovery/call works through ADM without bypassing Environment policy.
- **MCP-RUN-05** ✅: Secret values are resolved only at activation boundaries and are not exposed in normal status/log output.

Evidence: `.planning/phases/03-external-mcp-runtime-completion/03-VERIFICATION.md` and deterministic Phase 3 UAT artifacts.

### LIFECYCLE — Phase 5 (integrated)

- **LIFE-01** ✅: Long-lived Runtime/MCP/process ownership belongs to one persistent local control boundary, not to a short CLI invocation. The existing Gateway is that boundary; Phase 5 concretely owns external MCP sessions there.
- **LIFE-02** ✅: Desired state and observed state are distinct; restart rebuilds live MCP runtime state from persisted Environment/catalog configuration rather than serializing owner/session/healthy observation.
- **LIFE-03** ✅: Resources owned by the persistent boundary are cleaned up deterministically. Phase 5 proves external MCP session cleanup on disable/removal/owner close and owner-bound graceful Gateway shutdown; Phase 6 extends the same ownership boundary to long-running dev child processes and verifies process-tree cleanup on the Windows target.

Evidence: `.planning/phases/05-persistent-runtime-ownership/05-VERIFICATION.md` and independent dogfood evidence.

### PROCESS — Phase 6 (integrated)

- **PROC-01** ✅: Long-running dev processes can be start/list/status/stop by stable ADM `proc_` identity under the persistent Gateway; start/stop remain writer-gated and reuse the short-exec Runtime authority boundary.
- **PROC-02** ✅: stdout/stderr are retained as bounded owner-memory tails and are queryable from later clients, including after process exit while that owner remains alive.
- **PROC-03** ✅ on the current Windows dogfood target: listening TCP ports are reported only as facts for ADM-owned process PIDs; the Agent API does not accept arbitrary PIDs or expose a generic OS process/port manager. Non-Windows currently returns no port facts rather than broadening authority.

Evidence: `.planning/phases/06-dev-process-logs-ports/06-VERIFICATION.md`, `06-UAT.md`, independent dogfood evidence, and `06-INTEGRATION-REVIEW.md` (review passed and integrated to local master 2026-09-08; post-integration validation is recorded in `evidence/integration.json`).

### ISOLATION — Phase 7 (integrated)

- **ISO-01** ✅: Git worktree is an optional isolation capability around Environment. Ordinary non-Git Workspace/Environment creation and Gateway file development remain green; managed isolation failure is operation-local.
- **ISO-02** ✅: managed worktrees use ADM-generated `wt_` identities, branches and destinations beneath the ADM state-directory worktree root; managed roots are revalidated for Environment relation, owned-root containment, Git top-level, common-dir and branch identity before routed Runtime access.
- **ISO-03** ✅: destroy requires the matching writer, refuses dirty or locally advanced/unpublished work by default, requires explicit force for unsafe removal, and always retains the generated branch so committed work is not silently deleted with the worktree directory.

Evidence: `.planning/phases/07-optional-git-worktree-isolation/07-VERIFICATION.md`, `07-UAT.md`, `07-INTEGRATION-REVIEW.md`, `evidence/regression.json`, and post-integration validation on local master.

### AGENT RUN — Phase 8 (locally verified; integration review pending)

- **ARUN-01** ✅ locally: the persistent Gateway owns stable single-command asynchronous `run_` resources independently from the launching MCP client. Later clients can list/status/cancel by ADM identity; cancel is writer-gated; terminal states distinguish succeeded, failed and canceled; timeout/bounded result behavior reuses the existing Runtime authority; owner close/drop cleans active Runs; restart begins with no prior Run observation and `state.json` contains none.

Evidence: `.planning/phases/08-agent-run-lifecycle/08-VERIFICATION.md`, `08-UAT.md`, and `evidence/regression.json`. Integration review is still pending; do not treat Phase 8 as integrated until separately reviewed/authorized.

## Active Requirements

### AGENT / GSD

- **FLOW-01**: Planner/Executor/Reviewer use structured auditable contracts and distinguish review failure from orchestration failure.
- **GSD-01**: V2 can read the repository's real `.planning/PROJECT.md`, `STATE.md`, current CONTEXT and PLAN provenance.
- **GSD-02**: GSD execution uses an explicit operation allowlist and Runtime capabilities; no arbitrary hidden shell execution.
- **GSD-03**: State advances only after verified pass and only to a pre-existing unambiguous next plan.
- **PAR-01**: Parallel execution is added only after single-run lifecycle and optional worktree isolation are validated.

## Deferred / Frozen

Until the Core roadmap reaches the Human Manager milestone, do not spend feature time on:

- installer / MSI;
- tray;
- autostart;
- updater;
- signing;
- notifications;
- launch-wrapper polish;
- visual redesign;
- migration/compatibility with old V1/V2 development state.

A concrete blocker may override this only if it prevents executing the current Phase acceptance test.

## Bootstrap Note

This planning reset is based on the real GSD planning artifacts and validated sequencing found in the earlier ADM implementation under `D:\projects\.ai-dev-manager-worktrees\...\.planning`.

Phase 1 closed the bootstrap gap on 2026-09-06: V2 discovered the actual installed OpenCode GSD Skill suite, exposed `gsd-next` and its authorized `gsd-core` supporting workflow through the ADM Gateway, and used that GSD path to normalize and advance this repository's planning state to Phase 2.

---
*Last updated: 2026-09-08 after Phase 8 local verification; integration review pending*