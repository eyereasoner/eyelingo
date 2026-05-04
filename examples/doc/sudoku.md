# Sudoku

## What this example is about

This example solves a Sudoku puzzle and emits the completed grid only after checking that the solution is legal.

A Sudoku solution must place digits 1 through 9 so that each row, each column, and each 3×3 box contains every digit exactly once.

## How it works, in plain language

The input puzzle is an 81-character string. Digits are given clues, and dots or zeroes are empty cells. The program parses that string into a 9×9 grid, counts the clues, and then fills the empty cells with backtracking search. It chooses the open cell with the fewest legal candidates first, which keeps the search small and deterministic.

## What to notice in the output

The output shows both the original puzzle and the completed grid. This helps readers verify that the given clues were preserved and makes the transformation easy to follow.

## What the trust gate checks

The trust gate verifies that the puzzle has 81 cells, the clues do not conflict, the solver finds exactly one solution within the search limit, every original clue is preserved, and the completed grid satisfies all row, column, and 3×3-box constraints.

## Run it

From the repository root:

```sh
node examples/sudoku.js
```

## Files

- [JavaScript example](../sudoku.js)
- [Input data](../input/sudoku.json)
- [Reference output](../output/sudoku.md)
