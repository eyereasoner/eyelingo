# High Trust RDF Bloom Envelope

`high_trust_bloom_envelope` is a Go translation/adaptation of Eyeling's `high-trust-rdf-bloom-envelope.n3`.

The context is trustworthy membership checking. Canonical graph data, Bloom-index evidence, a decimal certificate, and false-positive bounds are all combined before accepting an envelope.

## What it demonstrates

This example is mainly in the **Technology** category. Bloom-envelope acceptance using canonical graph, index, and false-positive checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/high_trust_bloom_envelope.json](../input/high_trust_bloom_envelope.json)

Go translation: [../high_trust_bloom_envelope.go](../high_trust_bloom_envelope.go)

Expected Markdown output: [../output/high_trust_bloom_envelope.md](../output/high_trust_bloom_envelope.md)
