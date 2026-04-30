# HarborSMR Insight Dispatch

`harbor_smr` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on port electrolysis dispatch decision with safety margin and policy checks. Its input fixture is organized around `CaseName`, `Question`, `RequestPurpose`, `RequestAction`, `HubAuthAt`, `Hub`, `Aggregate`, `Insight`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: PERMIT - North Quay Hydrogen Hub may use https://example.org/insight/harborsmr to run PEM_electrolyzer_train_2 at 16 MW from 2026-06-18T14:00:00Z to 2026-06-18T18:00:00Z.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - reserve margin 24 MW exceeds threshold 19 MW C2 OK - cooling margin 18% exceeds threshold 14% C3 OK - no planned outage blocks the balancing window

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/harbor_smr.json](../input/harbor_smr.json)

Go translation: [../harbor_smr.go](../harbor_smr.go)

Expected Markdown output: [../output/harbor_smr.md](../output/harbor_smr.md)
