# Deep Taxonomy 100000

`deep_taxonomy_100000` is a Go translation/adaptation of Eyeling's `deep-taxonomy-100000.n3`.

The context is large-scale materialization. A very deep taxonomy chain stresses rule application, termination, and counting, while the output keeps the benchmark result compact.

## How it works

A self-contained Go translation of deep-taxonomy-100000.n3 from the Eyeling
examples.

The original N3 file is a stress test for long rule chains. It contains
100,000 nearly identical taxonomy-step rules:

	{?X a :N42} => {?X a :N43, :I43, :J43}.

Starting from :ind a :N0, the chain must reach :N100000. The terminal rule
then derives :A2, and the success rule derives :test :is true.

This Go version translates the repeated rule family into a compact bit-set
forward chainer. A bit set is a memory-efficient way to record which classes
have been reached. The program is intentionally not a generic RDF/N3 reasoner;
it keeps this one deep classification derivation visible and deterministic
without embedding five megabytes of repetitive Go literals.

## What it demonstrates

This example is mainly in the **Technology** category. Large taxonomy materialization benchmark through a very deep class chain.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/deep_taxonomy_100000.json](../input/deep_taxonomy_100000.json)

Go translation: [../deep_taxonomy_100000.go](../deep_taxonomy_100000.go)

Expected Markdown output: [../output/deep_taxonomy_100000.md](../output/deep_taxonomy_100000.md)
