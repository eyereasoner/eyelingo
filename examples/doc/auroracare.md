# AuroraCare

`auroracare` is a Go translation/adaptation of Eyeling's `auroracare.n3`.

The context is health-data governance. Care, quality-improvement, and research uses are evaluated against purpose, consent, minimization, and policy conditions, giving a compact example of rule-based permit/deny reasoning.

## How it works

A self-contained Go translation of the Eyeling AuroraCare purpose-based
medical-data exchange example.

The source N3 file models a policy decision point (PDP). The same health data
can be allowed or denied depending on purpose, requester role, care-team
relationship, environment, privacy safeguards, and patient consent. This Go
version keeps those concrete policy checks visible as ordinary structs and
functions.

This is intentionally not a generic RDF, ODRL, DPV, EHDS, or N3 reasoner. It
translates the source facts and rules for this one example into explicit Go
data and deterministic checks, then emits a compact report.

## What it demonstrates

This example is mainly in the **Science** category. Health-data permit/deny scenarios across care, quality improvement, and research.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/auroracare.json](../input/auroracare.json)

Go translation: [../auroracare.go](../auroracare.go)

Expected Markdown output: [../output/auroracare.md](../output/auroracare.md)
