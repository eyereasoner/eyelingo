# Docking Abort Token

`docking_abort_token` is a Go translation/adaptation of Eyeling's `act-docking-abort-token.n3`.

The context is safety-token handling for an automated docking abort. The example distinguishes allowed propagation of an audit token from unsafe copying or fan-out behavior.

## How it works

Inspired by Eyeling's `examples/act-docking-abort.n3`.

This example derives possible classical abort-token tasks and impossible
quantum-seal tasks from a small constructor-theory style fixture.

## What it demonstrates

This example is mainly in the **Engineering** category. Docking abort audit-token flow and safety-system copy restrictions.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/docking_abort_token.json](../input/docking_abort_token.json)

Go translation: [../docking_abort_token.go](../docking_abort_token.go)

Expected Markdown output: [../output/docking_abort_token.md](../output/docking_abort_token.md)
