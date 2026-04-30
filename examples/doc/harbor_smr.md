# HarborSMR Insight Dispatch

`harbor_smr` is a Go translation/adaptation of Eyeling's `harborsmr.n3`.

The context is port-energy dispatch. A small modular reactor insight is used only if scope, expiry, privacy, safety margin, and dispatch thresholds all pass.

## What it demonstrates

This example is mainly in the **Engineering** category. Port electrolysis dispatch decision with safety margin and policy checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/harbor_smr.json](../input/harbor_smr.json)

Go translation: [../harbor_smr.go](../harbor_smr.go)

Expected Markdown output: [../output/harbor_smr.md](../output/harbor_smr.md)
