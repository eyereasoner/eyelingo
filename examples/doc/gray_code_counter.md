# Gray Code Counter

## What this example is about

This example generates a 4-bit Gray code. A Gray code is a sequence of binary states where neighboring states differ by exactly one bit.

This is useful in hardware, sensors, and encoders because changing one bit at a time reduces ambiguity during transitions.

## How it works, in plain language

The program counts from 0 to 15 and maps each integer `n` to `n xor (n >> 1)`. That formula creates the reflected binary Gray-code sequence.

After generating the sequence, it checks each neighboring pair, including the transition from the final state back to the first state.

## What to notice in the output

The output reports 16 visited states and 16 unique states, so the sequence covers the whole 4-bit state space. The maximum adjacent Hamming distance is 1, meaning no transition changes more than one bit.

## What the trust gate checks

The trust gate verifies the bit width, sequence length, uniqueness, cyclic wraparound behavior, and one-bit-transition property. These checks catch duplicate states and incorrect binary transformations.

## Run it

From the repository root:

```sh
node examples/gray_code_counter.js
```

## Files

- [JavaScript example](../gray_code_counter.js)
- [Input data](../input/gray_code_counter.json)
- [Reference output](../output/gray_code_counter.md)
