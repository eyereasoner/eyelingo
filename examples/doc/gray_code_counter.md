# Gray Code Counter

`gray_code_counter` is a Go translation/adaptation of Eyeling's `gray-code-counter.n3`.

## Background

A Gray code orders binary words so consecutive values differ in exactly one bit. This is useful in hardware and sensing because a one-step transition cannot momentarily flip several bits at once. The example builds the sequence and checks the one-bit Hamming-distance property between neighboring codes.

## What it demonstrates

**Category:** Technology. n-bit Gray-code sequence with one-bit transition checks.

## How the Go implementation works

The Go code generates a reflected binary Gray sequence for the configured bit width. It compares each adjacent pair, including the wrap-around pair, by counting bit differences.

The checks verify sequence length, uniqueness, valid bit width, and the one-bit transition property.

## Files

Input JSON: [../input/gray_code_counter.json](../input/gray_code_counter.json)

Go implementation: [../gray_code_counter.go](../gray_code_counter.go)

Python check: [../checks/gray_code_counter.py](../checks/gray_code_counter.py)

Expected Markdown output: [../output/gray_code_counter.md](../output/gray_code_counter.md)
