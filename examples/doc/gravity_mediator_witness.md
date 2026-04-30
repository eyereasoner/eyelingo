# Gravity Mediator Witness

`gravity_mediator_witness` is a Go translation/adaptation of Eyeling's `act-gravity-mediator-witness.n3`.

## Background

In quantum-information thought experiments, entanglement between two systems can reveal something about the mediator that couples them. If gravity could mediate entanglement, that would suggest non-classical behavior in the mediator; a purely classical channel should not create the same witness. This example contrasts those cases using explicit CAN/CAN'T-style evidence.

## What it demonstrates

**Category:** Science. Mediator-only entanglement witness contrasting non-classical and purely classical gravitational mediators.

## How the Go implementation works

The implementation evaluates each experimental run from its mediator model, coupling mode, observed outcome, and control/probe status. It identifies the positive witness run and the contrast run, then derives whether mediator-only entanglement evidence is present.

Checks confirm that the positive case satisfies all witness conditions and that the classical contrast is blocked for a specific reason.

## Files

Input JSON: [../input/gravity_mediator_witness.json](../input/gravity_mediator_witness.json)

Go implementation: [../gravity_mediator_witness.go](../gravity_mediator_witness.go)

Expected Markdown output: [../output/gravity_mediator_witness.md](../output/gravity_mediator_witness.md)
