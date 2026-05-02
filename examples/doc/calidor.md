# Calidor

`calidor` is a Go translation/adaptation of Eyeling's `calidor.n3`.

## Background

Heat-response planning often needs sensitive household signals, such as health risk, energy access, or indoor conditions, but public services should not receive more private data than necessary. A privacy-preserving workflow can replace raw signals with a signed, purpose-limited insight that authorizes a narrow intervention. This example uses that pattern for cooling support during a heat-alert window.

## What it demonstrates

**Category:** Engineering. Municipal cooling intervention bundle chosen from active needs and budget constraints.

## How the Go implementation works

The implementation represents needs, service packages, signatures, budget limits, and authorization rules as typed structs. It filters active needs, matches them to package capabilities, checks whether the envelope is signed, current, minimized, and purpose-bound, then chooses the lowest eligible package.

The audit output records why the package was accepted and which requirements were covered.

## Files

Input JSON: [../input/calidor.json](../input/calidor.json)

Go implementation: [../calidor.go](../calidor.go)

Python check: [../checks/calidor.py](../checks/calidor.py)

Expected Markdown output: [../output/calidor.md](../output/calidor.md)
