# Phase 14 Context — Evidence-first Investigation Toolkit

## Product boundary

ADM remains a local development control plane. Phase 14 may add investigation helpers only when they return concrete evidence and improve on plain text search. ADM must not infer tasks, plan work, advance `.planning` state, or choose orchestration policy.

## Phase 14 goal

Expose high-value code/runtime investigation helpers that an external Agent can call to gather bounded evidence. Every helper must return evidence, confidence, and uncertainties; no opaque AI guesses.

## First slice

14-01 implements endpoint evidence resolution. It takes one URL/path and optional HTTP method, searches the Environment root for route-like literals and nearby method evidence, and returns bounded file/line evidence with confidence and uncertainties.

## Non-goals

- No server startup, network request, browser probe, verifier run, command execution, or MCP tool call.
- No semantic code execution or framework-specific interpreter.
- No task planning or automatic fix suggestion.
- No `.planning` state advancement beyond normal plan/summary docs.

## Current development rule

After Phase 13 was merged back, active development continues directly on `master` in `D:\projects\ai-dev-manager-v2`. Old rebaseline ADM workspace/env metadata has been cleaned; the physical Git worktree is not deleted by ADM metadata cleanup.
