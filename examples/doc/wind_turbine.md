# Wind Turbine Envelope

## What this example is about

This example classifies wind-speed samples for a turbine and estimates energy production over a short sequence of intervals.

A turbine has practical operating thresholds: too little wind produces no useful power, moderate wind produces partial power, strong but safe wind reaches rated power, and very high wind stops the turbine for safety.

## How it works, in plain language

For each wind-speed sample, the program compares the speed with cut-in, rated, and cut-out thresholds. Below cut-in and at or above cut-out, power is zero. Between cut-in and rated speed, power follows a cubic curve. Between rated speed and cut-out, power is capped at rated power.

Energy is then accumulated by multiplying each interval's power by the interval duration.

## What to notice in the output

The output lists every interval with speed, status, and power. It also reports four usable intervals and total energy of 1.571 MWh. The interval-by-interval listing makes the aggregate result easier to audit.

## What the trust gate checks

The trust gate verifies the expected counts of usable, rated, and stopped intervals, and checks that total energy is stable after rounding. This catches threshold-boundary and duration-conversion errors.

## Run it

From the repository root:

```sh
node examples/wind_turbine.js
```

## Files

- [JavaScript example](../wind_turbine.js)
- [Input data](../input/wind_turbine.json)
- [Expected output](../output/wind_turbine.md)
