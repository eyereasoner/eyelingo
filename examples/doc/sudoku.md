# Sudoku

`sudoku` translates a selected N3-style reasoning scenario into a compact Go example. It focuses on sudoku constraint solving with a unique completed grid. Its input fixture is organized around `Puzzle`, `Name`.

The example keeps the reasoning deliberately visible: the JSON file supplies the facts or parameters, the Go file encodes the translated rules and calculations, and the Markdown output records the result in ARC style.

## What it demonstrates

This is mainly a **Mathematics** example. It demonstrates exact computation, formal constraints, certificates, and algorithmic invariants in a form that can be read as code, data, and expected output.

In plain words, the answer section highlights: The puzzle is solved, and the completed grid is the unique valid Sudoku solution. case : sudoku default puzzle : classic

## How to read the output

`Answer` gives the computed conclusion or selected result.

`Reason why` explains the rule path, calculation path, or decision chain that led to the answer.

`Check` records invariants that should hold if the translation is faithful and the computation is consistent.

For this example, the checks include: C1 OK - every given clue is preserved in the final grid. C2 OK - the final grid contains only digits 1 through 9, with no blanks left. C3 OK - each row contains every digit exactly once.

`Go audit details` separates implementation evidence from the domain conclusion: input sizes, thresholds, counters, source scenario names, precision choices, or platform details.

## Files

Input JSON: [../input/sudoku.json](../input/sudoku.json)

Go translation: [../sudoku.go](../sudoku.go)

Expected Markdown output: [../output/sudoku.md](../output/sudoku.md)
