# 8-Queens

## What this example is about

This example solves the classic 8-Queens puzzle: place eight queens on an 8×8 chessboard so that no two queens attack each other.

A queen attacks along its row, column, and diagonals. The challenge is to place one queen in each row without sharing columns or diagonals.

## How it works, in plain language

The solver places queens row by row. It uses bit masks to remember which columns and diagonals are already occupied. That makes it fast to list only the safe columns for the next row.

The example prints only the first board, but it keeps counting after that so the total number of solutions is complete.

## What to notice in the output

The printed board is one valid arrangement. The final line reports 92 total solutions for the 8-Queens puzzle, showing that the program did not stop after the first solution.

## What the trust gate checks

The trust gate verifies that at least one solution is found, the reported total matches the enumerated solutions, every solution is valid, and the solution list contains no duplicates. These checks catch mistakes in diagonal handling, counting, or early termination.

## Run it

From the repository root:

```sh
node examples/queens.js
```

## Files

- [JavaScript example](../queens.js)
- [Input data](../input/queens.json)
- [Reference output](../output/queens.md)
