# Parcel Locker

`parcellocker` is a Go translation/adaptation of Eyeling's `parcellocker.n3`.

The context is delegated authorization. A one-time parcel pickup request is permitted only when identity, token, time, locker, and delegation constraints are all satisfied.

## What it demonstrates

This example is mainly in the **Technology** category. Delegated parcel pickup-token authorization.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/parcellocker.json](../input/parcellocker.json)

Go translation: [../parcellocker.go](../parcellocker.go)

Expected Markdown output: [../output/parcellocker.md](../output/parcellocker.md)
