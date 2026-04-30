# Gradient Descent Step

`gradient_descent_step` is a Go translation/adaptation of Eyeling's `gd-step-certified.n3`.

The context is certified numerical optimization. A single gradient-descent step for a quadratic objective is computed and checked against the expected decrease conditions.

## What it demonstrates

This example is mainly in the **Mathematics** category. Certified single gradient-descent step for a quadratic objective.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/gradient_descent_step.json](../input/gradient_descent_step.json)

Go translation: [../gradient_descent_step.go](../gradient_descent_step.go)

Expected Markdown output: [../output/gradient_descent_step.md](../output/gradient_descent_step.md)
