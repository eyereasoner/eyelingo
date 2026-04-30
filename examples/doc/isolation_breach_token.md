# Isolation Breach Token

`isolation_breach_token` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on isolation-breach audit-token flow with cloning and fan-out restrictions. Its input fixture is organized around `caseName`, `question`, `variable`, `media`, `superinformationMedium`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: classical breach token : YES, prepare, reversible permutation, copy, measure, serial audit, and parallel fan-out all succeed specimen provenance seal : NO, universal cloning and unrestricted parallel fan-out are blocked prepared witness : nursePager prepares CodeBreach

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - doorBeacon, containmentPLC, nursePager, and incidentBoard encode BreachBit C2 OK - nursePager can prepare CodeBreach C3 OK - doorBeacon can permute SafeGreen to BreachRed and back

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/isolation_breach_token.json](../input/isolation_breach_token.json)

Go translation: [../isolation_breach_token.go](../isolation_breach_token.go)

Expected Markdown output: [../output/isolation_breach_token.md](../output/isolation_breach_token.md)
