# Cold Chain Release

`cold_chain_release` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on biologic lot release and allocation under cold-chain and transit constraints. Its input fixture is organized around `CaseName`, `Question`, `Product`, `Lots`, `RecalledSources`, `Events`, `Clinics`, `Transit`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Science** example. It demonstrates scientific measurement, evidence handling, and domain checks in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Release decision : 3 of 6 candidate lots can ship. Safe release lots : BIO-17A, BIO-20A, BIO-21A Quarantined lots : BIO-17B, BIO-18A, BIO-19A

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - every final lot has a computed ancestry closure. C2 OK - recalled Seed-Beta propagates to BIO-18A and not to unrelated safe lots. C3 OK - custody ledgers are hash-linked or explicitly rejected.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/cold_chain_release.json](../input/cold_chain_release.json)

Go translation: [../cold_chain_release.go](../cold_chain_release.go)

Expected Markdown output: [../output/cold_chain_release.md](../output/cold_chain_release.md)
