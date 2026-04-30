# Bayes Diagnosis

`bayes_diagnosis` is a Go translation/adaptation of Eyeling's `bayes-diagnosis.n3`.

The context is medical diagnostic support under uncertainty. Symptoms and test evidence are combined into posterior probabilities, then ranked so the answer can be read as a transparent Bayesian decision aid rather than as a black-box classifier.

## What it demonstrates

This example is mainly in the **Science** category. Bayesian posterior ranking of possible diseases from symptoms and test evidence.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/bayes_diagnosis.json](../input/bayes_diagnosis.json)

Go translation: [../bayes_diagnosis.go](../bayes_diagnosis.go)

Expected Markdown output: [../output/bayes_diagnosis.md](../output/bayes_diagnosis.md)
