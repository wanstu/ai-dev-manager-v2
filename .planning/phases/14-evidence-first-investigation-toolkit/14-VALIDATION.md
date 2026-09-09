# Phase 14 Validation

Phase 14 helpers must prove investigation is evidence-only, bounded and side-effect-free.

## Required gates for every Phase 14 slice

- focused app tests for the slice;
- focused Gateway tests for the slice when an Agent-facing tool is added;
- `go test -count=1 ./...`;
- `go vet ./...`;
- `git diff --check`.

## Safety checks for every Phase 14 slice

- no writer lease acquisition unless the explicit operation is already writer-gated;
- no process start, verifier run, shell command, endpoint call, browser probe, MCP tool call, or provider indexing during passive inspection;
- no private Memory values in evidence;
- evidence must be bounded by file/match/byte limits;
- every result must include confidence and uncertainties when evidence is partial, stale or heuristic;
- provider failures must be local and must not break built-in fallback investigation.

## 14-01 gates

14-01 endpoint evidence resolution proved:

- exact route literal evidence;
- conservative dynamic route candidates;
- optional same-line HTTP method evidence;
- bounded results;
- private Memory non-leakage;
- Gateway `investigate_endpoint` output;
- full test/vet/diff gates at `0e40394`.

## 14-02 planned gates

14-02 optional code intelligence provider + GitNexus integration must prove:

- GitNexus/provider absence reports `unconfigured` or `disabled`, not failure;
- provider failure reports structured unavailable/degraded state while static fallback remains available;
- passive provider inspection does not start GitNexus, run indexing, call MCP tools, call endpoints or mutate files;
- provider evidence includes provenance, freshness and uncertainty;
- GitNexus remains optional and can be integrated through existing MCP/runtime authorization rather than a new privileged path;
- no private Memory or secret values leak through provider diagnostics.
