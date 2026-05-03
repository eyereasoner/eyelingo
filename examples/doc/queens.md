# 8-Queens

`queens` is a Go translation/adaptation of Eyeling's `queens.n3`.

## Background

The N-Queens problem asks for queen placements on an N by N chessboard so that no two queens share a row, column, or diagonal. It is a standard constraint-satisfaction problem because each placement restricts later choices. The eight-queen case is small enough to solve quickly but rich enough to show pruning and backtracking.

## What it demonstrates

**Category:** Mathematics. 8-Queens constraint satisfaction with a valid board solution.

## How the Go implementation works

The implementation solves N-Queens with recursive bit-mask backtracking. Each recursion level places one row, while three masks track occupied columns and both diagonal directions.

The solver counts all valid boards, captures a limited number of examples for the report, and records search counters for auditability.

## Files

Input JSON: [../input/queens.json](../input/queens.json)

Go implementation: [../queens.go](../queens.go)

Go check: [../../checks/main.go](../../checks/main.go)

Expected Markdown output: [../output/queens.md](../output/queens.md)
