# Dijkstra Risk Path

## What this example is about

This example chooses a route through a small delivery graph. Each route segment has a normal cost and a risk value. The question is not simply “what is shortest?” but “what is best after risk is priced in?”

## How it works, in plain language

The program gives every edge a risk-adjusted score: delivery cost plus a risk penalty. The risk penalty is the edge risk multiplied by the configured risk weight.

It then uses Dijkstra's shortest-path algorithm. Dijkstra expands the currently cheapest known partial route first. Once the destination is removed from the queue, the selected route is optimal for the weighted graph.

## What to notice in the output

The selected route goes through DepotB and LabD. A route through DepotC looks cheaper early, but it carries enough risk that it loses after the risk penalty is applied. The output separates raw cost, risk sum, and risk-adjusted score so the trade-off is visible.

## What the trust gate checks

The trust gate verifies that edge costs and risks are non-negative, at least one path reaches the goal, the selected path starts and ends correctly, the path appears in the enumerated simple paths, and no enumerated simple path has a lower risk-adjusted score.

## Run it

From the repository root:

```sh
node examples/dijkstra_risk_path.js
```

## Files

- [JavaScript example](../dijkstra_risk_path.js)
- [Input data](../input/dijkstra_risk_path.json)
- [Reference output](../output/dijkstra_risk_path.md)
