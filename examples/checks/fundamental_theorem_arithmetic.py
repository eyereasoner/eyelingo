# Independent Python checks for the fundamental_theorem_arithmetic example.
from __future__ import annotations

import math
import re
from collections import Counter

from .common import run_checks


def is_prime(n: int) -> bool:
    if n < 2:
        return False
    if n % 2 == 0:
        return n == 2
    d = 3
    while d * d <= n:
        if n % d == 0:
            return False
        d += 2
    return True


def factor(n: int) -> list[int]:
    out = []
    while n % 2 == 0:
        out.append(2)
        n //= 2
    d = 3
    while d * d <= n:
        while n % d == 0:
            out.append(d)
            n //= d
        d += 2
    if n > 1:
        out.append(n)
    return out


def format_power_form(factors: list[int]) -> str:
    counts = Counter(factors)
    parts = []
    for p in sorted(counts):
        e = counts[p]
        parts.append(str(p) if e == 1 else f"{p}^{e}")
    return " * ".join(parts)


def parse_metric(answer: str, label: str) -> int | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9]+)", answer)
    return int(match.group(1)) if match else None


def parse_sample_rows(answer: str) -> dict[int, str]:
    rows = {}
    for n, factors in re.findall(r"^\s+([0-9]+)\s+=\s+(.+)$", answer, flags=re.MULTILINE):
        rows[int(n)] = factors.strip()
    return rows


def run(ctx):
    numbers = [int(n) for n in ctx.load_input()]
    factorizations = {n: factor(n) for n in numbers}
    primary = 202692987
    rows = parse_sample_rows(ctx.answer)
    distinct_primes = sorted({p for factors in factorizations.values() for p in factors})
    total_factors = sum(len(factors) for factors in factorizations.values())
    trace_factors = [int(x) for x in re.findall(r"(?:=|\*)\s+([0-9]+)\s+\*", ctx.reason)]

    checks = [
        ("the primary source number is factored independently", factorizations[primary] == [3, 3, 7, 829, 3881]),
        ("the reported primary prime-power form matches grouped exponents", f"primary prime-power form : {format_power_form(factorizations[primary])}" in ctx.answer),
        ("multiplying every factor list reconstructs its source integer", all(math.prod(factors) == n for n, factors in factorizations.items())),
        ("every distinct factor found by trial division is prime", all(is_prime(p) for p in distinct_primes)),
        ("the report includes one sample factorization row for every JSON number", set(rows) == set(numbers)),
        ("every reported sample row matches the independently formatted factorization", all(rows[n] == format_power_form(factorizations[n]) for n in rows)),
        ("smallest-first and largest-first traversals sort to the same primary multiset", "source sorted comparison : 3 * 3 * 7 * 829 * 3881" in ctx.reason),
        ("reported sample count and largest sample match JSON", parse_metric(ctx.answer, "sample count") == len(numbers) and parse_metric(ctx.answer, "largest sample") == max(numbers)),
        ("reported total factor multiplicity and distinct-prime count match recomputation", parse_metric(ctx.answer, "total prime factors counted with multiplicity") == total_factors and parse_metric(ctx.answer, "distinct primes seen across samples") == len(distinct_primes)),
        ("the ten-digit prime remains unfactored after trial division", factorizations[9999999967] == [9999999967] and is_prime(9999999967)),
    ]
    return run_checks(checks)
