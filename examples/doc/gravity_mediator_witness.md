# Gravity Mediator Witness

`gravity_mediator_witness` is a Go translation/adaptation of Eyeling's `act-gravity-mediator-witness.n3`.

The context is a physics witness scenario. The rules distinguish mediator-only entanglement evidence from classical alternatives, so the conclusion depends on which transformations are possible.

## How it works

Inspired by Eyeling's `examples/act-gravity-mediator-witness.n3`.

## What it demonstrates

This example is mainly in the **Science** category. Mediator-only entanglement witness contrasting non-classical and purely classical gravitational mediators.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/gravity_mediator_witness.json](../input/gravity_mediator_witness.json)

Go translation: [../gravity_mediator_witness.go](../gravity_mediator_witness.go)

Expected Markdown output: [../output/gravity_mediator_witness.md](../output/gravity_mediator_witness.md)
