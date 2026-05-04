# Complex Matrix Stability

## What this example is about

This example classifies three small complex-valued matrices for a discrete-time system. In everyday terms, it asks whether repeated application of a matrix makes a signal grow, stay on the edge of stability, or shrink away.

The matrices are diagonal, which keeps the example focused on the stability rule rather than on matrix algebra machinery.

## How it works, in plain language

For a diagonal matrix, the important values are already on the diagonal. Each diagonal entry is a complex number. The program computes the modulus, or distance from zero, of each complex number.

The largest modulus is called the spectral radius. A radius greater than 1 means repeated steps can grow, so the system is unstable. A radius equal to 1 is marginally stable. A radius below 1 is damped.

## What to notice in the output

The output labels the three fixtures as unstable, marginally stable, and damped. This makes the stability boundary visible: the difference between radius 1 and radius 2 is not just numerical; it changes the long-term behavior.

## What the trust gate checks

The trust gate verifies that all three matrices are classified, that the expected radius/class pairs are stable, and that a basic complex-number identity about product moduli holds. That last check protects the small complex arithmetic helper used by the example.

## Run it

From the repository root:

```sh
node examples/complex_matrix_stability.js
```

## Files

- [JavaScript example](../complex_matrix_stability.js)
- [Input data](../input/complex_matrix_stability.json)
- [Expected output](../output/complex_matrix_stability.md)
