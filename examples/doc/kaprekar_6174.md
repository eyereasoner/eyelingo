# Kaprekar 6174

`kaprekar_6174` is a Go translation/adaptation of Eyeling's `kaprekar-6174.n3`.

## Background

Kaprekar's routine for four-digit numbers sorts the digits descending and ascending, subtracts the smaller number from the larger, and repeats. For many four-digit starts with at least two distinct digits, the chain reaches the fixed point 6174. Leading zeros matter because the routine always treats the value as a four-digit string.

## What it demonstrates

**Category:** Mathematics. Kaprekar chains and basin facts ending at 6174.

## How the Go implementation works

The implementation runs Kaprekar's four-digit routine with leading zero handling. Each step sorts the digits high-to-low and low-to-high, subtracts, and appends the result to the chain until it reaches 6174 or the bounded limit.

The checks verify convergence, the seven-step bound, trace consistency, and omission of equal-digit starts that fall into 0000.

## Files

Input JSON: [../input/kaprekar_6174.json](../input/kaprekar_6174.json)

Go implementation: [../kaprekar_6174.go](../kaprekar_6174.go)

Expected Markdown output: [../output/kaprekar_6174.md](../output/kaprekar_6174.md)
