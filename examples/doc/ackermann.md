# Ackermann

`ackermann` is a Go translation/adaptation of Eyeling's `ackermann.n3`.

## Background

The Ackermann function is a classic example of total computable recursion that grows much faster than primitive-recursive functions such as addition, multiplication, and exponentiation. Small inputs already produce very large integers, so it is a useful stress test for exact arithmetic and for implementations that must make recursion and memoization visible rather than approximate the result.

## What it demonstrates

**Category:** Mathematics. Exact Ackermann and hyperoperation facts, including very large integer results.

## How the Go implementation works

The Go code represents each requested Ackermann value as a typed query and reduces it through a three-argument helper. That helper switches between successor, addition, multiplication, exponentiation, and higher repeated-power cases, with memoization for repeated subproblems.

All arithmetic uses `math/big`, so the large integer results remain exact. Oversized values are summarized with digit counts and SHA-256 fingerprints before the report is written.

## Files

Input JSON: [../input/ackermann.json](../input/ackermann.json)

Go implementation: [../ackermann.go](../ackermann.go)

Expected Markdown output: [../output/ackermann.md](../output/ackermann.md)
