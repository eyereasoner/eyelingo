# Fibonacci Example (Big)

## What this example is about

This example computes a very large Fibonacci number exactly. The target index is 10000, which is far too large for ordinary fixed-size integer arithmetic.

## How it works, in plain language

The Fibonacci sequence starts with `F(0)=0` and `F(1)=1`. Every later value is the sum of the two previous values. The program repeats that rule up to the target index.

The important implementation detail is arbitrary-precision integer arithmetic. In JavaScript, that means using `BigInt`, so the result does not overflow into an approximate floating-point number.

## What to notice in the output

The output prints the exact value of `F(10000)`. The number is long, but the example is intentionally simple: it demonstrates that a small trusted derivation can still produce a large exact result.

## What the trust gate checks

The trust gate verifies the configured target, sample recurrence checks, the two base cases, and that the computed target value is emitted as an integer string. The input no longer carries the final Fibonacci answer; the program derives it.

## Run it

From the repository root:

```sh
node examples/fibonacci.js
```

## Files

- [JavaScript example](../fibonacci.js)
- [Input data](../input/fibonacci.json)
- [Reference output](../output/fibonacci.md)
