# Flandor

`flandor` is a Go translation/adaptation of Eyeling's `flandor.n3`.

The context is regional industrial policy. Exporters, training capacity, intervention packages, and payload-signature checks are combined to select a retooling priority.

## How it works

A self-contained Go translation of the Eyeling Flandor insight-economy
example.

The source N3 file models a regional policy decision. Exporters, training
actors, and grid operators keep their sensitive details local. The region only
receives a narrow macro-economic insight: enough aggregated pressure is active
to justify a temporary retooling response for Flanders.

This program is intentionally not a generic RDF, ODRL, crypto, or N3 reasoner.
It translates the concrete facts and rules for this example into ordinary Go
structs and checks. That keeps the main idea visible: private micro-signals can
become a minimal, signed, expiring decision object for a public policy board.

## What it demonstrates

This example is mainly in the **Engineering** category. Regional retooling priority calculation for a Flanders scenario.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/flandor.json](../input/flandor.json)

Go translation: [../flandor.go](../flandor.go)

Expected Markdown output: [../output/flandor.md](../output/flandor.md)
