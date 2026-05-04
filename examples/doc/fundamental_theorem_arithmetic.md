# Fundamental Theorem Arithmetic

## What this example is about

This example factors several integers into primes. It illustrates the core idea behind the Fundamental Theorem of Arithmetic: every integer greater than 1 can be written as a product of primes, and that prime multiset is unique up to order.

## How it works, in plain language

The program uses trial division. It repeatedly finds the smallest divisor of the current number, records that divisor, and divides it out. When no smaller divisor remains, the remaining number is prime.

The output also rewrites repeated factors using exponents, such as `3^2` instead of `3 * 3`.

## What to notice in the output

The primary sample `202692987` factors into `3 * 3 * 7 * 829 * 3881`. The output includes a step-by-step decomposition trace, then compares the same factors in reverse order and sorted order to show that order does not change the underlying factor multiset.

## What the trust gate checks

The trust gate verifies that the primary number is present, every factorization multiplies back to the original number, every factor is prime, the sample set is non-empty, all samples are integers greater than one, and prime samples remain single factors.

## Run it

From the repository root:

```sh
node examples/fundamental_theorem_arithmetic.js
```

## Files

- [JavaScript example](../fundamental_theorem_arithmetic.js)
- [Input data](../input/fundamental_theorem_arithmetic.json)
- [Reference output](../output/fundamental_theorem_arithmetic.md)
