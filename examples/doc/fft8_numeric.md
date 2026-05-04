# FFT8 Numeric

## What this example is about

This example shows a tiny discrete Fourier transform over eight samples. It is a compact way to see how a time-domain signal can be described by frequency bins.

The input samples represent one sine wave cycle measured at eight evenly spaced points.

## How it works, in plain language

The program projects the eight sample values onto eight complex frequency patterns. Each projection is a frequency bin. A large magnitude in a bin means the original signal contains that frequency.

For a real sine wave, energy appears in a positive-frequency bin and the matching negative-frequency bin. With eight samples, those show up as two dominant bins.

## What to notice in the output

The output prints the sample vector, the dominant bins, and an energy check. The time-domain energy and the scaled frequency-domain energy match, which is a numerical form of Parseval's identity.

## What the trust gate checks

The trust gate verifies that the dominant bins really have maximum magnitude, the DC component cancels within tolerance, the sine bins have the expected magnitude for this transform, Parseval energy is preserved, and the real signal gives conjugate symmetry. These checks catch mistakes in complex arithmetic, bin indexing, or normalization.

## Run it

From the repository root:

```sh
node examples/fft8_numeric.js
```

## Files

- [JavaScript example](../fft8_numeric.js)
- [Input data](../input/fft8_numeric.json)
- [Reference output](../output/fft8_numeric.md)
