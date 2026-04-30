# Calidor

`calidor` is a Go translation/adaptation of Eyeling's `calidor.n3`.

The context is municipal response planning during heat or cooling stress. Needs, packages, budgets, signatures, and policy constraints are combined to decide whether an intervention bundle is admissible.

## How it works

A self-contained Go translation of the Eyeling Calidor heat-response
insight-economy example.

The source N3 file models an urgent city service without making the city
collect raw household sensor, health, or prepaid-energy data. A local gateway
turns those private signals into a narrow insight: this household needs
priority cooling support during a heat-alert window. The city may use that
signed and expiring envelope for heatwave response, but not for unrelated
purposes such as tenant screening.

This program is intentionally not a generic RDF, ODRL, crypto, or N3 reasoner.
It translates the concrete facts and rules for this one example into ordinary
Go structs and checks. That keeps the privacy pattern visible: raw signals stay
local, while a minimal decision object is shared for a specific public-service
purpose.

## What it demonstrates

This example is mainly in the **Engineering** category. Municipal cooling intervention bundle chosen from active needs and budget constraints.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/calidor.json](../input/calidor.json)

Go translation: [../calidor.go](../calidor.go)

Expected Markdown output: [../output/calidor.md](../output/calidor.md)
