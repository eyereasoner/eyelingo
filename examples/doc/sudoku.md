# Sudoku

## What this example is about

This example solves a Sudoku puzzle and emits the completed grid only after checking that the solution is legal.

A Sudoku solution must place digits 1 through 9 so that each row, each column, and each 3×3 box contains every digit exactly once.

## How it works, in plain language

The input puzzle is an 81-character string. Digits are given clues, and dots or zeroes are empty cells. The program parses that string into a 9×9 grid, counts the clues, and compares the completed grid against the puzzle.

The solution grid is fixed in the example, so the emphasis is on verification and explanation rather than on search performance.

## What to notice in the output

The output shows both the original puzzle and the completed grid. This helps readers verify that the given clues were preserved and makes the transformation easy to follow.

## What the trust gate checks

The trust gate verifies that the puzzle has 81 cells, has the expected clue count, preserves every original clue, and that the completed grid satisfies all row, column, and 3×3-box constraints.

## Run it

From the repository root:

```sh
node examples/sudoku.js
```

## Files

- [JavaScript example](../sudoku.js)
- [Input data](../input/sudoku.json)
- [Expected output](../output/sudoku.md)
