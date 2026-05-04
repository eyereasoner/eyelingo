# Goldbach 1000

## What this example is about

This example checks a bounded version of Goldbach's conjecture. The full conjecture says every even integer greater than 2 is the sum of two primes. This example only checks the finite range from 4 through 1000.

## How it works, in plain language

The program first builds a cache of prime numbers up to the configured maximum. Then, for each even number, it searches for a prime `P` no larger than half of that number such that the remainder is also prime.

When it finds such a pair, that pair is a witness for that even number.

## What to notice in the output

The output says all 499 even integers in the range have witnesses and prints a few samples, such as `1000 = 3 + 997`. This is a bounded computational result, not a proof of the infinite conjecture.

## What the trust gate checks

The trust gate verifies that the expected count of even numbers is checked, that every even number has a valid prime-pair witness, and that the configured upper bound is stable.

## Run it

From the repository root:

```sh
node examples/goldbach_1000.js
```

## Files

- [JavaScript example](../goldbach_1000.js)
- [Input data](../input/goldbach_1000.json)
- [Expected output](../output/goldbach_1000.md)
