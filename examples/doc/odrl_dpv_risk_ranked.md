# Ranked DPV Risk Report

`odrl_dpv_risk_ranked` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on oDRL/DPV clause risk ranking by severity and risk class. Its input fixture is organized around `Consumer`, `Agreement`, `Process`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Technology** example. It demonstrates data representation, interoperability, policies, and computational artifacts in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: Agreement: Example Agreement Profile: Example consumer profile score=100 (risk:HighRisk, risk:HighSeverity) clause C1

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - 4 risk rows were derived. C2 OK - ranked output is in descending score order. C3 OK - score range is 70 to 100.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/odrl_dpv_risk_ranked.json](../input/odrl_dpv_risk_ranked.json)

Go translation: [../odrl_dpv_risk_ranked.go](../odrl_dpv_risk_ranked.go)

Expected Markdown output: [../output/odrl_dpv_risk_ranked.md](../output/odrl_dpv_risk_ranked.md)
