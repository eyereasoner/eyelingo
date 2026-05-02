# Flandor

`flandor` is a Go translation/adaptation of Eyeling's `flandor.n3`.

## Background

Regional economic planning often needs information from exporters, training providers, and grid operators, but those actors may not be able to reveal sensitive microdata. A minimal signed insight can aggregate enough pressure indicators to justify a public response while keeping firm-level details local. This example uses that pattern for a Flanders retooling-priority decision.

## What it demonstrates

**Category:** Engineering. Regional retooling priority calculation for a Flanders scenario.

## How the Go implementation works

The program models regional retooling pressure with exporters, training capacity, grid signals, intervention packages, and a signed decision envelope. It aggregates only the narrow indicators needed for the policy decision, then checks authorization, scope, expiry, and threshold rules.

The chosen priority is reported together with the private-signal minimization checks that allowed it.

## Files

Input JSON: [../input/flandor.json](../input/flandor.json)

Go implementation: [../flandor.go](../flandor.go)

Python check: [../checks/flandor.py](../checks/flandor.py)

Expected Markdown output: [../output/flandor.md](../output/flandor.md)
