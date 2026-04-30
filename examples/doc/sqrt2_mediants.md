# Integer-First Sqrt2 Mediants

`sqrt2_mediants` is a Go translation/adaptation of Eyeling's `integer-first-sqrt2-mediants.n3`.

## Background

The square root of two is irrational, so it cannot be represented exactly as a fraction. Rational approximations can still be certified by squaring: if p²/q² is below 2, p/q is a lower bound; if it is above 2, p/q is an upper bound. Integer comparisons avoid floating-point rounding while narrowing the bracket.

## What it demonstrates

**Category:** Mathematics. Rational lower/upper bounds for sqrt(2) certified by integer square comparisons.

## How the Go implementation works

The Go code builds rational lower and upper bounds for `sqrt(2)` using integer arithmetic. It generates continued-fraction convergents, compares squared integers instead of floats, and keeps the certified bracket.

Checks confirm that the lower bound squares below 2, the upper bound squares above 2, and the interval width meets the target.

## Files

Input JSON: [../input/sqrt2_mediants.json](../input/sqrt2_mediants.json)

Go implementation: [../sqrt2_mediants.go](../sqrt2_mediants.go)

Expected Markdown output: [../output/sqrt2_mediants.md](../output/sqrt2_mediants.md)
