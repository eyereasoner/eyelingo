# FFT8 Numeric

`fft8_numeric` is a Go translation/adaptation of Eyeling's `fft8-numeric.n3`.

The context is signal processing. Eight samples from a sine wave are transformed into frequency bins, and the checks verify dominant bins, conjugate symmetry, zero DC, and energy preservation.

## What it demonstrates

This example is mainly in the **Technology** category. Eight-point Fourier transform over a sampled sine wave with conjugate-bin and energy checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/fft8_numeric.json](../input/fft8_numeric.json)

Go translation: [../fft8_numeric.go](../fft8_numeric.go)

Expected Markdown output: [../output/fft8_numeric.md](../output/fft8_numeric.md)
