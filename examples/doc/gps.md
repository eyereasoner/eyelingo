# GPS — Goal Driven Route Planning

## What this example is about

This example compares two simple routes from Gent to Oostende. It is named “GPS” because it demonstrates goal-driven route planning: start somewhere, reach a destination, and compare candidate routes.

## How it works, in plain language

Each route has four scores: duration, cost, belief, and comfort. Duration and cost are added along the route. Belief and comfort are multiplied, because each route segment affects the overall reliability and comfort.

The planner compares the direct route through Brugge with an alternative route through Kortrijk and Brugge.

## What to notice in the output

The direct route is faster, cheaper, more reliable, and slightly more comfortable. Because it dominates the alternative on all four measures, the explanation is easy to audit.

## What the trust gate checks

The trust gate verifies that both expected routes are found and that the direct route has the expected duration, cost, belief, comfort, and recommendation. This catches mistakes in route enumeration and score aggregation.

## Run it

From the repository root:

```sh
node examples/gps.js
```

## Files

- [JavaScript example](../gps.js)
- [Input data](../input/gps.json)
- [Expected output](../output/gps.md)
