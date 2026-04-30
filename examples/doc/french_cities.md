# French Cities

`french_cities` is a Go translation/adaptation of Eyeling's `french-cities.n3`.

## Background

A route graph represents cities as nodes and direct connections as edges. Reachability asks whether one node can be reached from another by following edges, optionally through intermediate nodes. This example keeps the graph small so the derived paths and transitive connections can be inspected directly.

## What it demonstrates

**Category:** Technology. Reachability over a small French city route graph.

## How the Go implementation works

The implementation builds a directed road graph from city links and derives reachability by expanding paths until no new destinations appear. It keeps the discovered predecessor/path information so reachable cities can be explained rather than merely listed.

Checks confirm the expected destinations and guard against unreachable-city false positives.

## Files

Input JSON: [../input/french_cities.json](../input/french_cities.json)

Go implementation: [../french_cities.go](../french_cities.go)

Expected Markdown output: [../output/french_cities.md](../output/french_cities.md)
