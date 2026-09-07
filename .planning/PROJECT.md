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

## Reality Audit — 2026-09-07

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

### Scope and remaining work

#### External MCP — validated Phase 3 tracer

Phase 3 completed Streamable HTTP connection definitions, Environment gating, four-state health, activation-boundary secret resolution and real list/call acceptance. Broader transports/auth and persistent session ownership remain deferred; see Phase 3 verification.

#### Memory

Persistence and explicit scoped reads/writes are real. Automatic context injection is not currently required or implemented; do not describe Memory as automatically supplied to every Agent request.

#### Desktop

Desktop is a functional management shell, not the product completion gate. It must remain frozen except for blockers until the Core milestones below are complete.

### Not implemented

- persistent Runtime/process ownership beyond the current Gateway process lifecycle;
- dev-server/process/log/port lifecycle;
- Git worktree Environment lifecycle/isolation;
- Agent Run lifecycle/status/cancel;
- Planner/Executor/Reviewer orchestration;
- GSD phase execution/verified state advance in V2;
- parallel Agent/worktree orchestration;
- production installer/tray/autostart/updater/signing.

## R1 Dogfood — reviewed and integrated

A real external Agent completed a separate non-Git verifier-report CLI through ADM only, consuming real GSD Skill artifacts, structured red/green verification and private Memory. Restart persistence and operation-local negative acceptance passed. A reproducible Windows command-child cancellation blocker was fixed. R1 milestone review passed on 2026-09-07 and Phase 04 was fast-forwarded into master; post-integration tests/vet/diff-check passed. Phase 5/Desktop expansion is still not started automatically. Evidence: `.planning/phases/04-external-agent-dogfood-gate/04-R1-REVIEW.md`.

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

## Active Requirements

### LIFECYCLE

- **LIFE-01**: Long-lived Runtime/MCP/process ownership belongs to one persistent local control boundary, not to a short CLI invocation.
- **LIFE-02**: Desired state and observed state are distinct; restart rebuilds real runtime state rather than serializing in-memory sessions.
- **LIFE-03**: Owned child processes are cleaned up deterministically.

### PROCESS

- **PROC-01**: Long-running dev processes can be start/list/status/stop by stable identity.
- **PROC-02**: stdout/stderr logs are bounded and queryable.
- **PROC-03**: known/listening ports are reported as facts without turning ADM into a generic OS process manager.

### ISOLATION

- **ISO-01**: Git worktree is an optional isolation capability around Environment, never an Environment prerequisite.
- **ISO-02**: Managed worktrees are created only under an ADM-owned root and are revalidated before use.
- **ISO-03**: destroy defaults safe around dirty/unpushed work and never silently deletes branches/changes.

### AGENT / GSD

- **ARUN-01**: Agent Run has stable identity, lifecycle state, status, and cancellation owned by the persistent control boundary.
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
*Last updated: 2026-09-07 after R1 milestone review and master integration*