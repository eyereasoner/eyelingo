# Allen Interval Calculus

`allen_interval_calculus` is a Go translation/adaptation of Eyeling's `allen-interval-calculus.n3`.

## Background

Allen interval calculus represents time with intervals rather than instants. Two intervals can stand in relations such as before, meets, overlaps, during, starts, finishes, equals, or the corresponding inverse relations. Those relations let a planner reason about temporal ordering even when each event has a start and end time instead of a single timestamp.

## What it demonstrates

**Category:** Mathematics. Allen temporal interval relation closure over completed and explicit intervals.

## How the Go implementation works

The implementation parses each interval into `time.Time` values, completes its duration, and classifies every ordered interval pair. The relation function encodes Allen's base cases directly, so each pair receives one relation such as before, meets, overlaps, starts, during, finishes, equals, or the inverse relation.

The report highlights required relations from the input and counts invalid or missing cases through explicit checks.

## Files

Input JSON: [../input/allen_interval_calculus.json](../input/allen_interval_calculus.json)

Go implementation: [../allen_interval_calculus.go](../allen_interval_calculus.go)

Expected Markdown output: [../output/allen_interval_calculus.md](../output/allen_interval_calculus.md)
