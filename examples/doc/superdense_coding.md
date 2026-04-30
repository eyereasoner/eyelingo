# Superdense Coding

`superdense_coding` is a Go translation/adaptation of Eyeling's `superdense-coding.n3`.

## Background

Superdense coding is a quantum-information protocol where a shared entangled pair lets one party transmit two classical bits by sending one qubit after applying one of four operations. The receiver decodes by measuring the joint state. This example represents the relevant parity and relation facts explicitly rather than simulating a physical quantum device.

## What it demonstrates

**Category:** Science. Quantum-information parity facts for superdense coding.

## How the Go implementation works

The Go code represents mobit relations, compositions, candidate messages, and parity cancellation explicitly. It generates the four two-bit message cases, applies the relation for Alice's half, and keeps only decoded answers with odd parity.

Checks verify that each message decodes to the expected classical bit pair.

## Files

Input JSON: [../input/superdense_coding.json](../input/superdense_coding.json)

Go implementation: [../superdense_coding.go](../superdense_coding.go)

Expected Markdown output: [../output/superdense_coding.md](../output/superdense_coding.md)
