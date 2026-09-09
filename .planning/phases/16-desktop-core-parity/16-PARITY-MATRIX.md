# Phase 16 Desktop Core Parity Matrix

## Priority legend

- **P0 / RC blocker:** unsafe or unusable without this.
- **P1 / RC target:** strongly expected for a useful Desktop RC; may use documented CLI/Gateway fallback only if implementation risk is too high.
- **P2 / Post-RC:** useful, but not required for local 1.0 RC.
- **N/A:** intentionally not a Desktop concern.

## Current inventory

| Core area | Core surface exists | Current Desktop exposure | Gap | RC priority |
|---|---:|---|---|---|
| Gateway lifecycle | Yes | Status/start/stop via adapter/frontend | Needs smoke confirmation only | P0 |
| Workspace list/add/rename/remove/inspect | Yes | Exposed | Needs smoke confirmation only | P0 |
| Environment list/create/rename/remove/inspect | Yes | Exposed | Needs smoke confirmation only | P0 |
| CapabilityReport / CapabilityFact | Yes | Included in `InspectEnvironment`, but not clearly rendered in frontend | Add visible structured unavailable/degraded capability facts | P0 |
| MCP simple catalog add/default/remove | Yes | Exposed for endpoint-style add/default/remove | Lacks typed config fields and import preview/apply | P1 |
| MCP Environment selection | Yes | Exposed | Needs health/status display in Environment detail | P1 |
| MCP health/status/probe | Yes | Adapter exposes `ProbeMCPHealth`; frontend usage appears limited | Add clear health/status action/result display | P1 |
| MCP import adapters | Yes | Management has preview/apply; Desktop adapter/frontend do not expose | Add Desktop import preview/apply or document CLI fallback for RC | P1 |
| MCP tool inventory/refresh | Yes through Gateway owner surfaces | Not clearly exposed in Desktop | Post-RC unless needed during smoke | P2 |
| Skill source add/list/refresh/remove | Yes | Management has source APIs; Desktop exposes basic root add but not source list/refresh/remove | Add source refresh/list/status or document CLI fallback for RC | P1 |
| Skill Environment selection | Yes | Exposed | Needs availability display in Environment detail | P1 |
| Skill availability/support inventory/read diagnostics | Yes | Management has APIs; Desktop adapter/frontend do not clearly expose | Add availability state/reason view; support file browsing may be post-RC | P1 |
| Exec allowlist | Yes | Add/remove/list exposed | Needs smoke confirmation only | P0 |
| Exec command execution | Yes | Not exposed | CLI/Gateway fallback acceptable for RC; Desktop execution UI can be post-RC | P2 |
| File tree/read/search/write/edit/delete | Yes | Not exposed as file manager | CLI/Gateway fallback acceptable for RC; do not build a broad editor for RC | P2 |
| Verifier list/run/results | Yes | Not exposed | Add list/run/result surface or documented CLI/Gateway fallback | P1 |
| Process lifecycle/logs/ports | Yes | Not exposed | Add process list/log/status/stop surface or documented CLI/Gateway fallback | P1 |
| Generic run lifecycle | Yes | Not exposed | Add run list/status/cancel surface or documented CLI/Gateway fallback | P1 |
| Git status/diff/branch | Yes | Not exposed | Useful but optional; CLI fallback acceptable | P2 |
| Managed worktree isolation | Yes | Not exposed | Optional; post-RC unless blocking daily use | P2 |
| Global Memory | Yes | List/read/write/delete exposed with explicit load/read path | Ensure snapshot does not leak values | P0 |
| Environment-private Memory | Yes | List/read/write/delete exposed after explicit Environment detail load | Ensure inspect/snapshot expose count only by default | P0 |
| 14-01 endpoint investigation | Yes | Gateway only | Post-RC; not required for Desktop RC | P2 |
| GitNexus/provider integration | Planned only | Not exposed | Post-RC unless dogfood proves blocker | P2 |
| Temporary resource lifecycle | Planned only | Not exposed | Post-RC unless unmanaged resources block daily use | P2 |
| CI auto build | Added in Phase 16 | GitHub Actions workflow | Must pass on GitHub before RC declaration | P0 |

## Initial RC implementation targets

1. CI baseline and local smoke build.
2. Capability report display in Desktop Environment detail.
3. MCP status/import minimum Desktop path or documented fallback.
4. Skill availability/source refresh minimum Desktop path or documented fallback.
5. Verifier/process/run visibility minimum Desktop path or documented fallback.
6. RC smoke checklist and known limitations.

## Keep out of RC unless proven blocker

- Full file manager/editor.
- Full Git/worktree UI.
- GitNexus/provider integration.
- Temporary resource cleanup implementation.
- Installer/tray/autostart/updater/signing.
