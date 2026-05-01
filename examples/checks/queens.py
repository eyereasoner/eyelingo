# Independent Python checks for the queens example.
from __future__ import annotations

import re

from .common import run_checks


def parse_columns(answer: str) -> list[int]:
    match = re.search(r"As column positions by row: \[([^\]]+)\]", answer)
    if not match:
        return []
    return [int(part.strip()) for part in match.group(1).split(",")]


def parse_total(answer: str) -> int | None:
    match = re.search(r"Total solutions for \d+-Queens: (\d+)", answer)
    return int(match.group(1)) if match else None


def count_n_queens(n: int) -> int:
    total = 0

    def search(row: int, cols: int, diag1: int, diag2: int) -> None:
        nonlocal total
        if row == n:
            total += 1
            return
        available = ((1 << n) - 1) & ~(cols | diag1 | diag2)
        while available:
            bit = available & -available
            available -= bit
            search(row + 1, cols | bit, (diag1 | bit) << 1, (diag2 | bit) >> 1)

    search(0, 0, 0, 0)
    return total


def run(ctx):
    data = ctx.load_input()
    n = int(data["N"])
    max_print = int(data["MaxPrint"])
    columns = parse_columns(ctx.answer)
    total = parse_total(ctx.answer)
    zero_based = [c - 1 for c in columns]
    diagonal_pairs = [
        (r1, r2)
        for r1 in range(len(zero_based))
        for r2 in range(r1 + 1, len(zero_based))
        if abs(zero_based[r1] - zero_based[r2]) == abs(r1 - r2)
    ]
    rendered_rows = [line for line in ctx.answer.splitlines() if set(line.strip().split()) <= {"Q", "."} and "Q" in line]
    independent_total = count_n_queens(n)

    checks = [
        ("the Python checker loaded the normalized 8-Queens input", n == 8 and max_print == 1),
        ("the printed solution gives one column for each row", len(columns) == n),
        ("all printed columns are within the board", all(1 <= c <= n for c in columns)),
        ("the printed solution uses each column exactly once", sorted(columns) == list(range(1, n + 1))),
        ("no pair of printed queens shares a diagonal", diagonal_pairs == []),
        ("the rendered board contains exactly eight rows and eight queens", len(rendered_rows) == n and sum(row.split().count("Q") for row in rendered_rows) == n),
        ("an independent Python bit-mask search counts 92 total solutions", independent_total == 92),
        ("the reported total matches the independent Python solution count", total == independent_total),
    ]
    return run_checks(checks)
