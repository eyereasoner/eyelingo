# Fibonacci Example (Big)

`fibonacci` is a Go translation/adaptation of Eyeling's `fibonacci.n3`.

## Background

The Fibonacci sequence starts with two base values and then adds the previous two terms to get the next term. Although the recurrence is simple, large indices quickly require arbitrary-precision integers. This makes the example a useful check that the implementation computes exact integer results rather than overflowing fixed-width machine integers.

## What it demonstrates

**Category:** Mathematics. Exact computation of a large Fibonacci number.

## How the Go implementation works

The Go code computes the requested Fibonacci value with arbitrary-precision integers. It uses an iterative dynamic-programming loop so the result is deterministic, exact, and able to handle very large indices.

Large values are summarized in the output with size and fingerprint information instead of flooding the report.

## Files

Input JSON: [../input/fibonacci.json](../input/fibonacci.json)

Go implementation: [../fibonacci.go](../fibonacci.go)

Expected Markdown output: [../output/fibonacci.md](../output/fibonacci.md)
