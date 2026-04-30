# Path Discovery

`path_discovery` is a Go translation/adaptation of Eyeling's `path-discovery.n3`.

## Background

Airport or transport routing can be represented as a directed graph with outbound connections from each node. A bounded route search explores candidate paths while avoiding cycles and respecting stopover limits. This example emphasizes controlled graph exploration over a larger fixture than the small city-route examples.

## What it demonstrates

**Category:** Technology. Airport path discovery with stopover and routing constraints.

## How the Go implementation works

The implementation loads the airport labels and outbound-route facts from JSON, then runs a bounded no-revisit route search from the configured origin. The search tracks stopover count and path state so it can avoid cycles while still exploring the large graph.

The report summarizes the matching route candidates and audits the graph size, touched candidates, and bound checks.

## Files

Input JSON: [../input/path_discovery.json](../input/path_discovery.json)

Go implementation: [../path_discovery.go](../path_discovery.go)

Expected Markdown output: [../output/path_discovery.md](../output/path_discovery.md)
