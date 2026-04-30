# 8-Queens

`queens` is a Go translation/adaptation of Eyeling's `queens.n3`.

The context is constraint satisfaction. The eight queens are placed so that rows, columns, and diagonals remain conflict-free, giving a compact combinatorial proof object.

## How it works

A small, standalone Go implementation of the N-Queens puzzle.

The program counts every way to place N queens on an N×N chessboard so that
no two queens attack each other. It uses a compact bit-mask backtracking
search: each recursive level represents one board row, and three masks track
columns plus the two diagonal directions that are already under attack.

Run with the default 8-Queens puzzle:

	go run queens.go

Or choose the board size and number of example boards to print:

	go run queens.go 14 1

The final “Go audit details” section is intentionally diagnostic. It exposes
the normalized command-line settings and search counters so the translation is
easy to inspect, compare, and regression-test.

## What it demonstrates

This example is mainly in the **Mathematics** category. 8-Queens constraint satisfaction with a valid board solution.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/queens.json](../input/queens.json)

Go translation: [../queens.go](../queens.go)

Expected Markdown output: [../output/queens.md](../output/queens.md)
