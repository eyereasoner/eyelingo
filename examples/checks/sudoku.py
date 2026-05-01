# Independent Python checks for the sudoku example.
from __future__ import annotations

import re

from .common import run_checks

DIGITS = set("123456789")


def parse_completed_grid(answer: str) -> list[list[int]]:
    if "Completed grid" not in answer:
        return []
    tail = answer.split("Completed grid", 1)[1]
    rows = []
    for line in tail.splitlines():
        line = line.strip()
        if not line:
            continue
        if not re.match(r"^[1-9 .|]+$", line):
            break
        digits = [int(ch) for ch in line if ch.isdigit()]
        if len(digits) == 9:
            rows.append(digits)
        if len(rows) == 9:
            break
    return rows


def givens_from_puzzle(puzzle: str):
    return [(i // 9, i % 9, int(ch)) for i, ch in enumerate(puzzle) if ch != "0"]


def legal(grid: list[list[int]], row: int, col: int, value: int) -> bool:
    if any(grid[row][c] == value for c in range(9) if c != col):
        return False
    if any(grid[r][col] == value for r in range(9) if r != row):
        return False
    box_r = row // 3 * 3
    box_c = col // 3 * 3
    for r in range(box_r, box_r + 3):
        for c in range(box_c, box_c + 3):
            if (r, c) != (row, col) and grid[r][c] == value:
                return False
    return True


EXPECTED_SOLUTION = [
    [1, 6, 2, 8, 5, 7, 4, 9, 3],
    [5, 3, 4, 1, 2, 9, 6, 7, 8],
    [7, 8, 9, 6, 4, 3, 5, 2, 1],
    [4, 7, 5, 3, 1, 2, 9, 8, 6],
    [9, 1, 3, 5, 8, 6, 7, 4, 2],
    [6, 2, 8, 7, 9, 4, 1, 3, 5],
    [3, 5, 6, 4, 7, 8, 2, 1, 9],
    [2, 4, 1, 9, 3, 5, 8, 6, 7],
    [8, 9, 7, 2, 6, 1, 3, 5, 4],
]


def run(ctx):
    data = ctx.load_input()
    puzzle = data["Puzzle"]
    grid = parse_completed_grid(ctx.answer)
    givens = givens_from_puzzle(puzzle)
    rows = [set(map(str, row)) for row in grid]
    columns = [set(str(grid[r][c]) for r in range(9)) for c in range(9)] if len(grid) == 9 else []
    boxes = []
    if len(grid) == 9:
        for br in range(0, 9, 3):
            for bc in range(0, 9, 3):
                boxes.append(set(str(grid[r][c]) for r in range(br, br + 3) for c in range(bc, bc + 3)))

    clue_count = len(givens)
    preserved = len(grid) == 9 and all(grid[r][c] == value for r, c, value in givens)
    full_digits = len(grid) == 9 and all(1 <= cell <= 9 for row in grid for cell in row)
    every_placement_legal = len(grid) == 9 and all(legal(grid, r, c, grid[r][c]) for r in range(9) for c in range(9))

    checks = [
        ("the input puzzle has 81 cells and exactly 23 given clues", len(puzzle) == 81 and clue_count == 23),
        ("the completed grid is parsed as nine rows of nine digits", len(grid) == 9 and all(len(row) == 9 for row in grid)),
        ("every original clue is preserved at the same row and column", preserved),
        ("the final grid contains only digits 1 through 9", full_digits),
        ("each completed row is a permutation of 1 through 9", len(rows) == 9 and all(row == DIGITS for row in rows)),
        ("each completed column is a permutation of 1 through 9", len(columns) == 9 and all(col == DIGITS for col in columns)),
        ("each completed 3×3 box is a permutation of 1 through 9", len(boxes) == 9 and all(box == DIGITS for box in boxes)),
        ("every filled cell is legal against its row, column, and box peers", every_placement_legal),
        ("the completed grid matches a separately embedded expected solution fixture", grid == EXPECTED_SOLUTION),
    ]
    return run_checks(checks)
