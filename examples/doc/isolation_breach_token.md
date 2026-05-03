# Isolation Breach Token

`isolation_breach_token` is a Go translation/adaptation of Eyeling's `act-isolation-breach-token.n3`.

## Background

Isolation-breach handling requires evidence to travel across systems while preserving provenance and safety boundaries. A token can prove that a breach signal is authorized, but copying or fanning it out without controls may create unsafe or misleading evidence. This example models which propagation steps are permitted and which are blocked.

## What it demonstrates

**Category:** Engineering. Isolation-breach audit-token flow with cloning and fan-out restrictions.

## How the Go implementation works

The Go code derives permitted breach-token propagation tasks for classical evidence and separates them from restricted provenance-seal operations. It records where copying or fan-out would violate the protected-token rules.

Checks compare the allowed and blocked task sets with the expected safety boundary.

## Files

Input JSON: [../input/isolation_breach_token.json](../input/isolation_breach_token.json)

Go implementation: [../isolation_breach_token.go](../isolation_breach_token.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/isolation_breach_token.md](../output/isolation_breach_token.md)
