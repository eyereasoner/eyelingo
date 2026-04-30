# Delfour

`delfour` is a Go translation/adaptation of Eyeling's `delfour.n3`.

The context is privacy-preserving retail insight. A household-shopping recommendation is only allowed when the insight, signature, minimization, and policy checks all line up.

## How it works

A self-contained Go translation of delfour.js, which itself is a concrete
JavaScript translation of examples/delfour.n3.

The scenario models a privacy-preserving retail interaction. A phone turns a
sensitive household condition into a narrow, expiring shopping insight; a
store scanner checks the envelope and policy; and, if authorized, suggests a
lower-sugar product.

This is intentionally not a generic RDF/N3 reasoner. The concrete N3 facts
and rules are represented as Go structs and ordinary functions so the data
flow is clear and directly runnable.

## What it demonstrates

This example is mainly in the **Technology** category. Privacy-preserving retail insight and recommendation policy.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/delfour.json](../input/delfour.json)

Go translation: [../delfour.go](../delfour.go)

Expected Markdown output: [../output/delfour.md](../output/delfour.md)
