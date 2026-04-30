# Path Discovery

`path_discovery` is a Go translation/adaptation of Eyeling's `path-discovery.n3`.

The context is large graph routing. Airport labels and outbound facts live in JSON, while the Go logic derives candidate paths under stopover and route constraints.

## How it works

A self-contained Go translation of path-discovery.n3 from the Eyeling
examples.

The original N3 file contains a large Neptune air-routes graph and a bounded
recursive rule:

	(?source ?destination () 0 2) :route ?airports

This Go version translates the query, the recursive no-revisit route rule,
and the complete airport/flight data loaded from examples/input/path_discovery.json:
7,698 airport labels and 37,505 nepo:hasOutboundRouteTo facts. The bounded
query still touches only the 338 outbound candidates reachable from
Ostend-Bruges within the two-stopover bound, but the full graph is loaded
and checked.

This is intentionally not a generic RDF/N3 reasoner. The concrete route rules
are represented as ordinary Go functions, and the concrete airport labels and
flight facts are loaded from JSON input so the derivation remains visible and
directly runnable.

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
