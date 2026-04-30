# Sudoku

`sudoku` is a Go translation/adaptation of Eyeling's `sudoku.n3`.

## Background

Sudoku is a grid-constraint puzzle where each row, column, and 3 by 3 box must contain the digits 1 through 9 exactly once. A solver alternates between eliminating impossible candidates and branching when several values remain possible. The solved grid is valid only if it preserves the original givens and satisfies all uniqueness constraints.

## What it demonstrates

**Category:** Mathematics. Sudoku constraint solving with a unique completed grid.

## How the Go implementation works

The implementation reads an 81-cell puzzle, normalizes blanks and givens, then solves it with constraint propagation plus depth-first search. Candidate sets are reduced by row, column, and box constraints before the search branches.

The report prints the solved grid and checks givens, uniqueness constraints, and solver status.

## Files

Input JSON: [../input/sudoku.json](../input/sudoku.json)

Go implementation: [../sudoku.go](../sudoku.go)

Expected Markdown output: [../output/sudoku.md](../output/sudoku.md)
