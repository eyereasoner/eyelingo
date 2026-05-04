# EV Roadtrip Planner

## What this example is about

This example is a small electric-vehicle route planner. It starts with a car in Brussels and searches for acceptable ways to reach Cologne while tracking battery state, route actions, duration, cost, reliability, and comfort.

## How it works, in plain language

The input describes possible actions such as driving, charging, buying a pass, or taking a shuttle. The planner composes these actions into candidate plans. For each plan, it adds duration and cost, and multiplies belief and comfort values across the chosen actions.

A plan is acceptable only if it satisfies the configured thresholds for reliability, cost, and duration. Among acceptable plans, the program selects the fastest one.

## What to notice in the output

The selected route uses a shuttle from Aachen to Cologne. That avoids an extra charge stop while staying inside the reliability, cost, and duration limits. The output also lists the top acceptable plans so the reader can compare near alternatives.

## What the trust gate checks

The trust gate checks that the planner finds the expected fastest acceptable plan, that the selected final state is Cologne with low battery and no pass, that the number of acceptable plans is stable, and that key score values match the fixture.

## Run it

From the repository root:

```sh
node examples/ev_roundtrip_planner.js
```

## Files

- [JavaScript example](../ev_roundtrip_planner.js)
- [Input data](../input/ev_roundtrip_planner.json)
- [Expected output](../output/ev_roundtrip_planner.md)
