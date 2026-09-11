# Dogfood Note — pjadm Long Synchronous Command Timeouts

Date: 2026-09-11
Status: open blocker candidate

## Observation

During Phase 17 follow-on validation, long synchronous `pjadm.exec` calls repeatedly timed out at the ChatGPT/tool boundary even when the ADM command `timeout_ms` was set much higher.

Observed examples:

- `go test -count=1 ./...` via synchronous `pjadm.exec` timed out before returning package-level results.
- `go test -count=1 ./internal/gateway` via synchronous `pjadm.exec` also timed out before returning results.
- Smaller package groups returned successfully, which suggests the issue is the synchronous request/response path or outer MCP/tool timeout rather than a simple test compile failure.

## Impact

This hurts daily dogfood because an Agent cannot reliably run natural long verification commands and get a clear terminal result. The workaround is to manually split commands into smaller groups or use asynchronous run lifecycle surfaces, but that workaround is easy to forget and makes failures ambiguous.

## Current workaround

- Prefer focused tests or small package groups for `pjadm.exec`.
- Prefer `run_start` + `run_status`/`run_list` for long-running commands when available.
- Keep `max_output_bytes` bounded to avoid returning huge payloads.

## Async path check

A `run_start` attempt for `go test -count=1 ./internal/gateway` returned a stable run id immediately, so async start avoids losing the ChatGPT turn to a synchronous outer timeout. However, while the run was still `running`, `run_status`/`run_list` returned only lifecycle metadata and no bounded live log/progress excerpt. The run was manually canceled to avoid leaving a long-running validation blocker active.

This means the async path is a useful workaround but not yet a complete dogfood answer for long verification, because the Agent still lacks an ergonomic running-output view and terminal result handoff inside normal chat cadence.

## Candidate fix direction

Do not increase global timeouts blindly. Instead, make the long-command path explicit and inspectable:

1. Document or surface a clear threshold where synchronous `exec` is not appropriate.
2. Prefer an async run flow for long commands, with stable run IDs and bounded status/log excerpts.
3. Consider having synchronous `exec` fail early with an actionable diagnostic when requested timeout exceeds the likely client/tool deadline.
4. Preserve the existing safety model: writer lease required, allowlisted executable only, bounded output, no arbitrary shell escape.

## Acceptance idea

A future dogfood slice is successful when an Agent can run repository-wide verification without losing the result to outer tool timeout, either by a first-class async path or by an explicit diagnostic that tells the Agent to use that path.
