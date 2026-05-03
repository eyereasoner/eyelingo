# Deep Taxonomy 100000

`deep_taxonomy_100000` is a Go translation/adaptation of Eyeling's `deep-taxonomy-100000.n3`.

## Background

A taxonomy is a hierarchy of classes where membership can propagate through subclass links. If A is a subclass of B and B is a subclass of C, then an instance of A is also an instance of C. Very deep taxonomies stress a materialization engine because it must apply transitive closure many times while still terminating and counting the derived facts accurately.

## What it demonstrates

**Category:** Technology. Large taxonomy materialization benchmark through a very deep class chain.

## How the Go implementation works

The Go code turns the long taxonomy chain into a compact forward-chaining loop backed by a bit set. Starting from the initial class, each reached class activates the next one, and the terminal condition derives the success flag.

This avoids storing thousands of repeated literals in code while still exercising the same termination, counting, and reachability behavior.

## Files

Input JSON: [../input/deep_taxonomy_100000.json](../input/deep_taxonomy_100000.json)

Go implementation: [../deep_taxonomy_100000.go](../deep_taxonomy_100000.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/deep_taxonomy_100000.md](../output/deep_taxonomy_100000.md)
