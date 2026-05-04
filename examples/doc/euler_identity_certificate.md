# Euler Identity Certificate

## What this example is about

This example gives a finite numerical witness for Euler's identity, the famous relationship behind `exp(iπ) + 1 = 0`.

The program does not claim to prove the identity symbolically. Instead, it computes a reproducible approximation and checks that the residual is smaller than an explicit tolerance.

## How it works, in plain language

The input says how many Taylor-series terms to use and what tolerance is acceptable. The program approximates `exp(iπ)` using complex arithmetic. It then adds 1 and measures how far the result is from zero.

Because computers use finite arithmetic, the answer is near zero rather than exactly zero. That is why the residual and tolerance are part of the explanation.

## What to notice in the output

The output shows that the computed real part is effectively `-1`, the imaginary part is effectively `0`, and the residual magnitude is tiny. The `within tolerance : true` line is the finite certificate.

## What the trust gate checks

The trust gate verifies the configured term count, that the residual is below tolerance, and that the real and imaginary parts remain stable at the chosen precision. This keeps the example honest about numerical approximation.

## Run it

From the repository root:

```sh
node examples/euler_identity_certificate.js
```

## Files

- [JavaScript example](../euler_identity_certificate.js)
- [Input data](../input/euler_identity_certificate.json)
- [Expected output](../output/euler_identity_certificate.md)
