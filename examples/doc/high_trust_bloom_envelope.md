# High Trust RDF Bloom Envelope

`high_trust_bloom_envelope` is a Go translation/adaptation of Eyeling's `high-trust-rdf-bloom-envelope.n3`.

## Background

A Bloom filter is a compact probabilistic set-membership structure: it can say an item is definitely absent or maybe present, but a maybe-present result can be a false positive. High-trust use therefore requires exact confirmation against a canonical source. This example treats the Bloom filter as a prefilter and separately certifies the expected false-positive workload.

## What it demonstrates

**Category:** Technology. Bloom-envelope acceptance using canonical graph, index, and false-positive checks.

## How the Go implementation works

The implementation treats the Bloom filter as a prefilter rather than as final proof. Maybe-positive membership hits are confirmed against the canonical graph, while decimal bounds estimate the expected false-positive workload.

The envelope is accepted only when canonical membership, index evidence, certificate precision, and trust thresholds all pass.

## Files

Input JSON: [../input/high_trust_bloom_envelope.json](../input/high_trust_bloom_envelope.json)

Go implementation: [../high_trust_bloom_envelope.go](../high_trust_bloom_envelope.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/high_trust_bloom_envelope.md](../output/high_trust_bloom_envelope.md)
