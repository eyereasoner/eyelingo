# Ranked DPV Risk Report

`odrl_dpv_risk_ranked` is a Go translation/adaptation of Eyeling's `odrl-dpv-risk-ranked.n3`.

## Background

ODRL expresses permissions, prohibitions, duties, and constraints over data use, while DPV-style vocabularies describe privacy concepts and risks. A risk report can combine policy clauses with safeguards and severity labels to rank the most important findings. This example makes that ranking explicit so mitigation priorities are visible.

## What it demonstrates

**Category:** Technology. ODRL/DPV clause risk ranking by severity and risk class.

## How the Go implementation works

The implementation stores clauses, permissions, prohibitions, safeguards, and privacy-risk labels as typed data. It derives risk findings from missing or weak safeguards, assigns severity scores, and sorts the findings into a ranked report.

The checks verify the expected top risks, severity ordering, and link between policy clauses and mitigation categories.

## Files

Input JSON: [../input/odrl_dpv_risk_ranked.json](../input/odrl_dpv_risk_ranked.json)

Go implementation: [../odrl_dpv_risk_ranked.go](../odrl_dpv_risk_ranked.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/odrl_dpv_risk_ranked.md](../output/odrl_dpv_risk_ranked.md)
