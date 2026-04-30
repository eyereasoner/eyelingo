# FFT32 Numeric

`fft32_numeric` is a Go translation/adaptation of Eyeling's `fft32-numeric.n3`.

The context is signal processing. The example evaluates complete 32-point Fourier spectra for multiple waveform fixtures, then checks the dominant frequency bins and spectral invariants.

## What it demonstrates

This example is mainly in the **Technology** category. Thirty-two-point Fourier transform over several sampled waveforms with dominant-bin, flat-spectrum, conjugate-symmetry, and energy checks.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/fft32_numeric.json](../input/fft32_numeric.json)

Go translation: [../fft32_numeric.go](../fft32_numeric.go)

Expected Markdown output: [../output/fft32_numeric.md](../output/fft32_numeric.md)
