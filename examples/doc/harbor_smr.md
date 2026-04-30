# HarborSMR Insight Dispatch

`harbor_smr` is a Go translation/adaptation of Eyeling's `harborsmr.n3`.

The context is port-energy dispatch. A small modular reactor insight is used only if scope, expiry, privacy, safety margin, and dispatch thresholds all pass.

## How it works

A self-contained Go translation inspired by Eyeling's `examples/harborsmr.n3`.

The scenario models a port hydrogen hub that asks to use a narrow,
permissioned SMR flexible-export insight for one four-hour electrolyzer
dispatch window. The key point is data minimization: the hub receives only a
bounded decision object, while raw reactor telemetry stays local to the SMR
operator.

This is intentionally not a generic RDF/N3/ODRL reasoner. The concrete facts
and rules are represented as Go structs and ordinary functions so the policy,
safety, and minimization checks stay visible and directly runnable.

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
