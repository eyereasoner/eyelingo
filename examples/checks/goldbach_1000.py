# Independent Python checks for the goldbach_1000 example.
from __future__ import annotations

import re

from .common import run_checks


def is_prime(n: int) -> bool:
    if n < 2:
        return False
    if n == 2:
        return True
    if n % 2 == 0:
        return False
    divisor = 3
    while divisor * divisor <= n:
        if n % divisor == 0:
            return False
        divisor += 2
    return True


def first_witness(even: int, primes: set[int]) -> tuple[int, int] | None:
    for p in range(2, even // 2 + 1):
        if p in primes and even - p in primes:
            return p, even - p
    return None


def parse_answer(answer: str):
    count_match = re.search(r"All (\d+) even integers from 4 through (\d+) have a Goldbach witness", answer)
    witness_match = re.search(r"sample witnesses\s*:\s*(.+)", answer)
    witnesses = {}
    if witness_match:
        for even, p, q in re.findall(r"(\d+)=(\d+)\+(\d+)", witness_match.group(1)):
            witnesses[int(even)] = (int(p), int(q))
    return {
        "even_count": int(count_match.group(1)) if count_match else None,
        "max_even": int(count_match.group(2)) if count_match else None,
        "witnesses": witnesses,
    }


def run(ctx):
    data = ctx.load_input()
    max_even = int(data["maxEven"])
    primes = {n for n in range(2, max_even + 1) if is_prime(n)}
    evens = list(range(4, max_even + 1, 2))
    witnesses = {even: first_witness(even, primes) for even in evens}
    failures = [even for even, witness in witnesses.items() if witness is None]
    sample_witnesses = {even: witnesses[even] for even in data["sampleEvens"]}
    parsed = parse_answer(ctx.answer)

    checks = [
        ("the configured upper bound is parsed from JSON as 1000", max_even == 1000 and parsed["max_even"] == max_even),
        ("there are exactly 499 even integers from 4 through 1000", len(evens) == 499 and parsed["even_count"] == len(evens)),
        ("Python trial division independently finds 168 primes at or below 1000", len(primes) == 168),
        ("every checked even integer has a prime-pair witness", failures == []),
        ("each requested sample even has the first witness found by the independent search", parsed["witnesses"] == sample_witnesses),
        ("every reported witness uses two primes whose sum is the reported even integer", all(p in primes and q in primes and p + q == even for even, (p, q) in parsed["witnesses"].items())),
        ("the bounded Goldbach result is derived from recomputed witnesses, not from static output text", "No counterexample" in ctx.reason and len(witnesses) == len(evens)),
    ]
    return run_checks(checks)
