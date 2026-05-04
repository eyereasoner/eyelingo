# Path Discovery

## What this example is about

This example searches an airport graph for routes from Ostend-Bruges International Airport to Václav Havel Airport Prague with at most two stopovers.

The source data is large, but the query is bounded so the explanation can show the relevant slice of the search.

## How it works, in plain language

The graph contains airports and outbound route facts. The program starts at the source airport and recursively follows outbound routes. It avoids revisiting airports already in the current path, and it stops expanding when the route would exceed the maximum number of stopovers.

A route is accepted only if it reaches the destination within the stopover limit.

## What to notice in the output

The output finds three routes, all with two stopovers. Each route starts at Ostend-Bruges, goes through Liège, then through one of three second-hop airports before reaching Prague.

The diagnostic counts show how a large graph can still produce an inspectable bounded result: number of labels, outbound facts, frontier airports expanded, recursive calls, and edge tests.

## What the trust gate checks

The trust gate verifies the expected number of routes, source and destination, stopover bound, route simplicity, and diagnostic search counts. This catches errors in recursion, depth limits, cycle prevention, or graph translation.

## Run it

From the repository root:

```sh
node examples/path_discovery.js
```

## Files

- [JavaScript example](../path_discovery.js)
- [Input data](../input/path_discovery.json)
- [Expected output](../output/path_discovery.md)
