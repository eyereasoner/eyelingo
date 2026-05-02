# AuroraCare

`auroracare` is a Go translation/adaptation of Eyeling's `auroracare.n3`.

## Background

Health-data access decisions depend on more than identity. A request may be allowed for treatment but denied for research, or allowed only when consent, role, minimization, purpose, and safeguards match the policy. This example frames those governance rules as explicit permit/deny logic so each decision can be traced back to conditions that passed or failed.

## What it demonstrates

**Category:** Science. Health-data permit/deny scenarios across care, quality improvement, and research.

## How the Go implementation works

The Go code builds a policy-decision fixture with subjects, requesters, policies, and access scenarios. Each scenario is evaluated by matching prohibitions first, then testing permission conditions such as purpose, role, care-team relationship, environment, safeguards, and consent.

Results are stored per scenario, and the checks verify that permitted and denied cases line up with the policy constraints rather than with the raw data category alone.

## Files

Input JSON: [../input/auroracare.json](../input/auroracare.json)

Go implementation: [../auroracare.go](../auroracare.go)

Python check: [../checks/auroracare.py](../checks/auroracare.py)

Expected Markdown output: [../output/auroracare.md](../output/auroracare.md)
