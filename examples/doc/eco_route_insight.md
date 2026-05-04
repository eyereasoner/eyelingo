# Eco Route Insight

## What this example is about

This example turns logistics data into a compact, privacy-preserving route insight. The scenario is a depot that may show an eco banner when the current route uses more fuel than a configured threshold and an acceptable lower-fuel alternative exists.

## How it works, in plain language

The program computes a fuel index from distance, payload, and gradient factor. It compares the current route with alternative routes. An alternative is eligible only if it saves fuel, stays within the allowed ETA delay, and satisfies the policy threshold rule.

Instead of exporting raw payload, GPS trace, driver behavior, or telemetry, the program creates a signed envelope with only the audience, allowed use, expiry, suggested route, compact fuel-index values, and a yes/no banner decision.

## What to notice in the output

The output shows an issued eco insight, the current and suggested fuel indices, the estimated saving, the expiry, and the signature metadata. The line `raw data exported : no` is central: the example is about shipping a decision, not the raw data behind the decision.

## What the trust gate checks

The trust gate verifies that the current route crosses the fuel threshold, the selected alternative is eligible, the ETA delay is acceptable, no other eligible alternative saves more fuel, forbidden raw-data terms are absent from the envelope, and the digest and signature match the canonical envelope.

## Run it

From the repository root:

```sh
node examples/eco_route_insight.js
```

## Files

- [JavaScript example](../eco_route_insight.js)
- [Input data](../input/eco_route_insight.json)
- [Reference output](../output/eco_route_insight.md)
