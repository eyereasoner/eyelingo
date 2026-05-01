# Independent Python checks for the fibonacci example.
from __future__ import annotations

import re

from .common import run_checks


def fib_pair(n: int) -> tuple[int, int]:
    if n == 0:
        return 0, 1
    a, b = fib_pair(n // 2)
    c = a * (2 * b - a)
    d = a * a + b * b
    if n % 2:
        return d, c + d
    return c, d


def fib(n: int) -> int:
    return fib_pair(n)[0]


def parse_answer(answer: str) -> tuple[int | None, int | None]:
    match = re.search(r"index\s+(\d+)\s+is:\s*\n\s*([0-9]+)", answer)
    if not match:
        return None, None
    return int(match.group(1)), int(match.group(2))


def run(ctx):
    data = ctx.load_input()
    expected = {int(k): int(v) for k, v in data.items()}
    computed = {n: fib(n) for n in expected}
    answer_index, answer_value = parse_answer(ctx.answer)
    largest_index = max(expected)

    small_recurrence_ok = all(fib(n) == fib(n - 1) + fib(n - 2) for n in range(2, 60))
    monotone_from_2 = all(fib(n + 1) >= fib(n) for n in range(2, 200))
    cassini_1000 = fib(1001) * fib(999) - fib(1000) * fib(1000)

    checks = [
        ("base cases are recomputed as F(0)=0 and F(1)=1", computed.get(0) == 0 and computed.get(1) == 1),
        ("the recurrence F(n)=F(n-1)+F(n-2) holds over an independent prefix", small_recurrence_ok),
        ("every JSON expected Fibonacci value matches fast-doubling recomputation", all(computed[n] == expected[n] for n in expected)),
        ("the report answers the largest requested index", answer_index == largest_index),
        ("the reported F(10000) exactly matches the independent big-integer value", answer_value == computed[largest_index]),
        ("F(1000) has the expected exact decimal length and endpoints", len(str(computed[1000])) == 209 and str(computed[1000]).startswith("434665576869") and str(computed[1000]).endswith("76137795166849228875")),
        ("F(10000) has the expected exact decimal length and endpoints", len(str(computed[10000])) == 2090 and str(computed[10000]).startswith("336447648764") and str(computed[10000]).endswith("6073310059947366875")),
        ("all requested Fibonacci numbers are nonnegative", all(value >= 0 for value in computed.values())),
        ("the Fibonacci sequence is nondecreasing from F(2) onward", monotone_from_2),
        ("Cassini's identity holds at n=1000 for the same independent generator", cassini_1000 == 1),
        ("the explanation names arbitrary-precision arithmetic for the large integer", "Arbitrary" in ctx.reason and "precision" in ctx.reason and "without overflow" in ctx.reason),
    ]
    return run_checks(checks)
