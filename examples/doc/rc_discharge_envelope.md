# RC Discharge Envelope

## What this example is about

This example models a capacitor discharging through a resistor. Instead of relying on an exact symbolic value at every step, it uses a small numerical interval that safely contains the decay factor.

The goal is to certify when the voltage envelope drops below a tolerance.

## How it works, in plain language

At each sample, the capacitor voltage is multiplied by a decay factor. The physical factor is represented symbolically as `exp(-1/4)`, but the program uses a finite interval around that value.

Because the interval is between 0 and 1, repeated multiplication shrinks the voltage envelope. The upper bound is the safety-relevant value: if the upper bound is below the tolerance, then every compatible exact trajectory is also below it.

## What to notice in the output

The output identifies step 13 as the first below-tolerance witness. It also prints the time and the upper-envelope voltage at that step, making the safety claim auditable.

## What the trust gate checks

The trust gate verifies that the decay interval contracts, the previous step is not yet below tolerance, the settled step is below the threshold, and the witness occurs before the configured maximum step.

## Run it

From the repository root:

```sh
node examples/rc_discharge_envelope.js
```

## Files

- [JavaScript example](../rc_discharge_envelope.js)
- [Input data](../input/rc_discharge_envelope.json)
- [Reference output](../output/rc_discharge_envelope.md)
