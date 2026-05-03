# Docking Abort Token

`docking_abort_token` is a Go translation/adaptation of Eyeling's `act-docking-abort-token.n3`.

## Background

A docking abort is a safety-critical decision where evidence must be relayed quickly but not altered or copied beyond the allowed trust boundary. Audit tokens capture the provenance and authorization of the abort signal. This example distinguishes ordinary propagation of safety evidence from restricted operations that would weaken accountability or allow uncontrolled fan-out.

## What it demonstrates

**Category:** Engineering. Docking abort audit-token flow and safety-system copy restrictions.

## How the Go implementation works

The program models abort-token media, possible propagation steps, and restricted seal behavior as data. It derives allowed classical token transfers separately from impossible unrestricted-copy tasks for the protected seal.

The checks confirm that the audit token can move where required while unsafe copying or fan-out is rejected.

## Files

Input JSON: [../input/docking_abort_token.json](../input/docking_abort_token.json)

Go implementation: [../docking_abort_token.go](../docking_abort_token.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/docking_abort_token.md](../output/docking_abort_token.md)
