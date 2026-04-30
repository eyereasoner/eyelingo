# Goldbach 1000

`goldbach_1000` is a Go translation/adaptation of Eyeling's `goldbach-1000.n3`.

## Background

The strong Goldbach conjecture says every even integer greater than two can be written as the sum of two primes. The conjecture is not proven in full, but a bounded checker can verify it for a finite range. This example tests every even number from 4 through 1000 and records the prime-pair witnesses.

## What it demonstrates

**Category:** Mathematics. Bounded strong-Goldbach checker for every even integer from 4 through 1000.

## How the Go implementation works

The program first builds a prime table up to the configured even bound. It then scans every even number from 4 through that bound and searches for a pair of primes that sums to it.

Witness pairs are stored for sample evens, and checks fail if any even number lacks a decomposition.

## Files

Input JSON: [../input/goldbach_1000.json](../input/goldbach_1000.json)

Go implementation: [../goldbach_1000.go](../goldbach_1000.go)

Expected Markdown output: [../output/goldbach_1000.md](../output/goldbach_1000.md)
