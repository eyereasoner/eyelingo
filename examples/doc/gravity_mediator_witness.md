# Gravity Mediator Witness

`gravity_mediator_witness` is a Go translation/adaptation of Eyeling's `act-gravity-mediator-witness.n3`.

The example models a mediator-only entanglement witness between two quantum sensors. The positive run uses a gravitational mediator under locality and interoperability; the contrast run uses a purely classical mediator model.

## What it demonstrates

This is mainly a **Science** example. The output separates what the witness can support from what the classical contrast model cannot support under the same mediator-only setup.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, counters, thresholds, derived facts, or platform details.

## Files

Input JSON: [../input/gravity_mediator_witness.json](../input/gravity_mediator_witness.json)

Go translation: [../gravity_mediator_witness.go](../gravity_mediator_witness.go)

Expected Markdown output: [../output/gravity_mediator_witness.md](../output/gravity_mediator_witness.md)
