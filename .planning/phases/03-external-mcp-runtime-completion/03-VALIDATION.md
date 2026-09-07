---
phase: 3
slug: external-mcp-runtime-completion
status: draft
nyquist_compliant: true
wave_0_complete: false
created: 2026-09-07
---

# Phase 3 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | go test (stdlib) |
| **Config file** | go.mod (module root) |
| **Quick run command** | `go test ./internal/catalog/... ./internal/app/... ./internal/gateway/...` |
| **Full suite command** | `go test ./...` |
| **Estimated runtime** | ~15 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./internal/catalog/... ./internal/app/... ./internal/gateway/...`
- **After every plan wave:** Run `go test ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 15 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Test Type | Automated Command | Status |
|---------|------|------|-------------|-----------|-------------------|--------|
| 03-01-01 | 03-01 | 1 | MCP-RUN-01 | unit | `go test ./internal/catalog/...` | ⬜ pending |
| 03-01-02 | 03-01 | 1 | MCP-RUN-03, MCP-RUN-05 | unit | `go test ./internal/app/...` | ⬜ pending |
| 03-01-03 | 03-01 | 1 | MCP-RUN-03, MCP-RUN-04 | integration | `go test ./internal/gateway/...` | ⬜ pending |
| 03-02-01 | 03-02 | 2 | MCP-RUN-03 | integration | `go test ./internal/management/... ./cmd/ai-dev-manager/...` | ⬜ pending |
| 03-02-02 | 03-02 | 2 | MCP-RUN-01..05 | validation | `go test ./... && go vet ./...` | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/app/mcp_health.go` — new file, health probe logic
- [ ] `internal/app/mcp_health_test.go` — new file, health probe unit tests
- [ ] `internal/gateway/mcp_acceptance_test.go` — new file, real HTTP acceptance tests

*Existing test infrastructure (go test) covers all phase requirements. No framework install needed.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Real local MCP initializes and lists tools | MCP-RUN-04 | Requires running external MCP process | Start test fixture MCP server, run `environment_mcp_tools` |
| Secret not visible in `mcp list` output | MCP-RUN-05 | CLI output inspection | Configure env-var endpoint, run `mcp list`, verify no resolved secret |

*All other phase behaviors have automated verification.*

---

## Validation Sign-Off

- [x] All tasks have `<automated>` verify commands
- [x] Sampling continuity: no 3 consecutive tasks without automated verify
- [x] Wave 0 covers all MISSING references
- [x] No watch-mode flags
- [x] Feedback latency < 15s
- [x] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
