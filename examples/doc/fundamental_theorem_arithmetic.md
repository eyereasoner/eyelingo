# Fundamental Theorem Arithmetic

`fundamental_theorem_arithmetic` is a Go translation/adaptation of Eyeling's `fundamental-theorem-arithmetic.n3`.

## Background

The fundamental theorem of arithmetic says every integer greater than one has a unique factorization into primes, apart from the order of the factors. A prime-power representation groups repeated factors by exponent. This example computes and checks that decomposition using exact integer arithmetic.

## What it demonstrates

**Category:** Mathematics. Prime factorization and prime-power representation.

## How the Go implementation works

The Go code factors each integer by repeatedly taking the smallest divisor, then groups equal primes into prime powers. It verifies the product reconstruction and tests each reported factor for primality.

To audit uniqueness up to order, the program compares smallest-first and largest-first traversals of the same factorization.

## Files

Input JSON: [../input/fundamental_theorem_arithmetic.json](../input/fundamental_theorem_arithmetic.json)

Go implementation: [../fundamental_theorem_arithmetic.go](../fundamental_theorem_arithmetic.go)

Expected Markdown output: [../output/fundamental_theorem_arithmetic.md](../output/fundamental_theorem_arithmetic.md)
