# Sudoku

`sudoku` is a Go translation/adaptation of Eyeling's `sudoku.n3`.

The context is grid-constraint solving. The puzzle is completed while preserving row, column, and box uniqueness, producing a verifiable solved board.

## How it works

Standalone Go translation of the Eyeling Sudoku example:
  - examples/sudoku.n3 supplies the default puzzle and report/check structure.
  - examples/builtin/sudoku.js supplies the solving and validation logic.

The program reads an 81-cell Sudoku puzzle, solves it with constraint propagation
plus depth-first search, and prints an N3-style explanation report with answer,
reasoning summary, consistency checks, and Go-specific audit details.

Usage:

	go run sudoku.go
	go run sudoku.go -puzzle "<81 cells using digits, 0, ., or _>"

Puzzle input accepts digits 1-9 for givens and 0, '.', or '_' for blanks.
Whitespace and common board separators such as '|', '+', and '-' are ignored.

## What it demonstrates

This example is mainly in the **Mathematics** category. Sudoku constraint solving with a unique completed grid.

The JSON file contains the example-specific facts, data, or parameters. The Go file makes the translated N3 rules, calculations, or search procedure explicit. The Markdown output records the result in ARC style so the answer, reasoning, checks, and implementation audit can be reviewed separately.

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

`Go audit details` separates implementation evidence from the domain conclusion: source scenario names, input sizes, thresholds, counters, precision choices, rule counts, or platform details.

## Files

Input JSON: [../input/sudoku.json](../input/sudoku.json)

Go translation: [../sudoku.go](../sudoku.go)

Expected Markdown output: [../output/sudoku.md](../output/sudoku.md)
