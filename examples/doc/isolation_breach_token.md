# Isolation Breach Token

`isolation_breach_token` is a Go translation/adaptation of Eyeling's `act-isolation-breach-token.n3`.

The context is safety evidence in an isolation-breach scenario. The token flow records what may be copied or relayed and what must remain restricted.

## How it works

Inspired by Eyeling's `examples/act-isolation-breach.n3`.

This example keeps the can/can't split: breach-token tasks are possible for
classical media, while the provenance seal refuses unrestricted fan-out.

## What it demonstrates

This example is mainly in the **Engineering** category. Isolation-breach audit-token flow with cloning and fan-out restrictions.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/isolation_breach_token.json](../input/isolation_breach_token.json)

Go translation: [../isolation_breach_token.go](../isolation_breach_token.go)

Expected Markdown output: [../output/isolation_breach_token.md](../output/isolation_breach_token.md)
