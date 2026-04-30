# Flandor

`flandor` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on regional retooling priority calculation for a Flanders scenario. Its input fixture is organized around `CaseName`, `Question`, `ExpectedFilesWritten`, `RequestPurpose`, `RequestAction`, `HubCreatedAt`, `HubExpiresAt`, `BoardAuthAt`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Name: Flandor Region: Flanders Metric: regional_retooling_priority

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - payload hash matches the source envelope digest C2 OK - HMAC value matches the trusted precomputed signature C3 OK - export weakness, skills strain, and grid stress reach the three-need threshold

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/flandor.json](../input/flandor.json)

Go translation: [../flandor.go](../flandor.go)

Expected Markdown output: [../output/flandor.md](../output/flandor.md)
