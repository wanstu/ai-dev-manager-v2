# Phase 14 Validation

14-01 must prove endpoint investigation is evidence-only and side-effect-free.

## Required gates

- focused app tests for exact and dynamic route evidence;
- focused Gateway test for `investigate_endpoint` exposure;
- full repository test;
- vet;
- diff check.

## Safety checks

- no writer lease required;
- no command execution;
- no network request;
- no verifier/process/run/MCP call;
- no private Memory values in report output.
