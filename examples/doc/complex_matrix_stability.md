# Complex Matrix Stability

`complex_matrix_stability` is a Go translation/adaptation of Eyeling's `complex-matrix-stability.n3`.

The context is discrete-time system stability. Complex diagonal entries are reduced to spectral radii, and the result classifies whether the matrix dynamics remain bounded.

## What it demonstrates

This example is mainly in the **Engineering** category. Discrete-time stability classification using spectral radii of diagonal complex matrices.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/complex_matrix_stability.json](../input/complex_matrix_stability.json)

Go translation: [../complex_matrix_stability.go](../complex_matrix_stability.go)

Expected Markdown output: [../output/complex_matrix_stability.md](../output/complex_matrix_stability.md)
