# Kaprekar 6174

## What this example is about

This example explores Kaprekar's routine for four-digit numbers. The routine sorts a number's digits descending and ascending, subtracts the smaller arrangement from the larger one, and repeats.

For many four-digit starts, this process reaches 6174, known as Kaprekar's constant.

## How it works, in plain language

Every start is treated as four digits, so `1` becomes `0001`. The program applies the Kaprekar step repeatedly, tracking chains that reach 6174 within the configured bound. Starts with all identical digits, such as `1111`, fall to `0000` instead and are omitted from the emitted Kaprekar facts.

## What to notice in the output

The output reports 9,990 emitted starts, 10 omitted starts, and a maximum of 7 steps to reach 6174. It also prints selected chains and a distribution showing how many starts take 0, 1, 2, ..., 7 steps.

## What the trust gate checks

The trust gate verifies the emitted and omitted counts, the maximum step count, selected known chains, and the step-count distribution. This protects the digit-padding, subtraction, and bounded-search logic.

## Run it

From the repository root:

```sh
node examples/kaprekar_6174.js
```

## Files

- [JavaScript example](../kaprekar_6174.js)
- [Input data](../input/kaprekar_6174.json)
- [Expected output](../output/kaprekar_6174.md)
