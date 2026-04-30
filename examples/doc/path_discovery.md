# Path Discovery

`path_discovery` is a Go translation/adaptation of Eyeling's `path-discovery.n3`.

The context is large graph routing. Airport labels and outbound facts live in JSON, while the Go logic derives candidate paths under stopover and route constraints.

## What it demonstrates

This example is mainly in the **Technology** category. Airport path discovery with stopover and routing constraints.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/path_discovery.json](../input/path_discovery.json)

Go translation: [../path_discovery.go](../path_discovery.go)

Expected Markdown output: [../output/path_discovery.md](../output/path_discovery.md)
