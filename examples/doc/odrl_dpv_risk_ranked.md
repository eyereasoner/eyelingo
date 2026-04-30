# Ranked DPV Risk Report

`odrl_dpv_risk_ranked` is a Go translation/adaptation of Eyeling's `odrl-dpv-risk-ranked.n3`.

The context is data-policy risk analysis. ODRL-style permissions and prohibitions are combined with DPV-like risk labels to produce ranked mitigation-oriented findings.

## How it works

A self-contained Go translation of `examples/odrl-dpv-risk-ranked.n3`.

The original N3 example models a Terms-of-Service style agreement using
policy and privacy-risk vocabulary terms. It links machine-readable
permissions and prohibitions to human-readable clauses, derives risks from
missing or weak safeguards, scores those risks, and emits a ranked report.

This program is intentionally not a generic RDF, ODRL, DPV, or N3 reasoner.
Instead, it translates the concrete facts and rules from the fixture into
typed Go data structures and explicit inference functions. That keeps the
rule mechanics visible while preserving the deterministic ranked output style
of the N3 `log:outputString` section.

## What it demonstrates

This example is mainly in the **Technology** category. ODRL/DPV clause risk ranking by severity and risk class.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/odrl_dpv_risk_ranked.json](../input/odrl_dpv_risk_ranked.json)

Go translation: [../odrl_dpv_risk_ranked.go](../odrl_dpv_risk_ranked.go)

Expected Markdown output: [../output/odrl_dpv_risk_ranked.md](../output/odrl_dpv_risk_ranked.md)
