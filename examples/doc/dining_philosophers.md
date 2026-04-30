# Dining Philosophers

`dining_philosophers` is a Go translation/adaptation of Eyeling's `dining-philosophers.n3`.

The context is concurrent resource use. A Chandy-Misra style trace is replayed so the output can check fairness, ownership, and conflict freedom between philosophers and forks.

## How it works

A self-contained Go translation of dining-philosophers.n3 from the Eyeling
examples.

The N3 source describes a concrete Chandy-Misra dining-philosophers trace.
Five philosophers sit in a ring and share five forks. A hungry philosopher
asks for any adjacent fork they do not hold. The current holder sends the fork
only when it is Dirty. A transferred fork arrives Clean, and after a meal the
trace marks every fork Dirty again.

This Go version keeps the same nine-round schedule from the N3 file:
P1/P3 eat, then P2/P4 eat, then P5 eats, repeated three times. Goroutines are
used inside each round to model philosophers making requests and checking
whether they can eat. State updates are still applied phase by phase so the
run is deterministic and easy to compare with the source trace.

## What it demonstrates

This example is mainly in the **Mathematics** category. Chandy-Misra dining-philosophers trace with concurrency conflict checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/dining_philosophers.json](../input/dining_philosophers.json)

Go translation: [../dining_philosophers.go](../dining_philosophers.go)

Expected Markdown output: [../output/dining_philosophers.md](../output/dining_philosophers.md)
