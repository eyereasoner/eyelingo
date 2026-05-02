# Eco Route Insight

`eco_route_insight` is a compact insight-economy example based on a local eco-routing demonstrator: raw logistics facts are used locally to produce a small, signed, purpose-bound insight.

## Background

A delivery tablet can know a shipment payload, a current route, and possible alternatives without exporting all of that data. The useful result is smaller: whether to show an eco banner, which route to suggest, what saving is expected, and which audience may consume the suggestion.

This example keeps the data small and deterministic. It demonstrates the pattern: ship the decision, not the data.

## What it demonstrates

**Category:** Technology. Local eco-route insight with a signed envelope and independent verification.

## How the Go implementation works

The Go implementation loads one shipment, one current route, two route alternatives, a fuel-index threshold, and a signature key. It computes:

- the current route fuel index as `distanceKm × payload tonnes × gradientFactor`;
- each alternative's fuel index, saving, and ETA delay;
- the best eligible lower-fuel alternative;
- a compact envelope containing only the audience, allowed use, expiry, route suggestion, fuel indices, saving, and signature metadata.

The Go example does not emit a `## Check` section. During testing, the Python check independently recomputes the fuel indices, route choice, envelope expiry, canonical payload digest, and HMAC signature.

## Files

Input JSON: [../input/eco_route_insight.json](../input/eco_route_insight.json)

Go implementation: [../eco_route_insight.go](../eco_route_insight.go)

Expected Markdown output: [../output/eco_route_insight.md](../output/eco_route_insight.md)
