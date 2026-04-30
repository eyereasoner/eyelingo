# Sudoku

`sudoku` is a Go translation/adaptation of Eyeling's `sudoku.n3`.

The context is grid-constraint solving. The puzzle is completed while preserving row, column, and box uniqueness, producing a verifiable solved board.

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
