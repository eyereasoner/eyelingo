# Dijkstra Risk Path

`dijkstra_risk_path` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on risk-adjusted path selection using weighted network edges. Its input fixture is organized around `caseName`, `question`, `start`, `goal`, `riskWeight`, `edges`, `expected`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Engineering** example. It demonstrates systems decisions, safety envelopes, route planning, and operational constraints in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: selected path : ClinicA -> DepotB -> LabD -> HubZ raw cost : 10.00 risk sum : 0.55

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - all route edges were loaded from JSON C2 OK - edge score is cost + 2.00 × risk C3 OK - Dijkstra reached HubZ from ClinicA

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/dijkstra_risk_path.json](../input/dijkstra_risk_path.json)

Go translation: [../dijkstra_risk_path.go](../dijkstra_risk_path.go)

Expected Markdown output: [../output/dijkstra_risk_path.md](../output/dijkstra_risk_path.md)
