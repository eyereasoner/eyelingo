# Bayes Therapy Decision Support

`bayes_therapy` is a Go translation/adaptation of Eyeling's `bayes-therapy-decision-support.n3`.

The context is treatment selection under probabilistic evidence. Candidate therapies are scored by posterior-weighted utility, so the example combines Bayesian inference with an explicit decision criterion.

## How it works

A self-contained Go translation of examples/bayes-therapy.n3 from the Eyeling
example suite, in ARC style.

The original N3 program extends a Naive‑Bayes diagnostic model with a
decision‑theoretic layer that scores five therapies by expected utility.
It then recommends the therapy with the highest utility.

This is intentionally not a generic N3 reasoner. The concrete N3 facts and
rules are represented as ordinary Go data and functions so the probabilistic
inference and decision logic are easy to read and directly runnable.

## What it demonstrates

This example is mainly in the **Science** category. Posterior-weighted utility selection of the best therapy.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/bayes_therapy.json](../input/bayes_therapy.json)

Go translation: [../bayes_therapy.go](../bayes_therapy.go)

Expected Markdown output: [../output/bayes_therapy.md](../output/bayes_therapy.md)
