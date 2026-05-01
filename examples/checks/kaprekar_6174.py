# Independent Python checks for the kaprekar_6174 example.
from __future__ import annotations

import re
from collections import Counter

from .common import run_checks


def step(n: int) -> int:
    digits = f"{n:04d}"
    asc = int("".join(sorted(digits)))
    desc = int("".join(sorted(digits, reverse=True)))
    return desc - asc


def chain(n: int, max_steps: int, target: int, zero: int) -> list[int]:
    if n in {target, zero}:
        return [n]
    out = []
    current = n
    for _ in range(max_steps):
        current = step(current)
        out.append(current)
        if current in {target, zero}:
            return out
    return out


def parse_int(answer: str, label: str) -> int | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9]+)", answer)
    return int(match.group(1)) if match else None


def parse_distribution(reason: str) -> dict[int, int]:
    return {int(step): int(count) for step, count in re.findall(r"^\s+([0-9]+) step\(s\) : ([0-9]+) start\(s\)", reason, flags=re.MULTILINE)}


def parse_selected(answer: str) -> dict[int, list[int]]:
    rows = {}
    for start, values in re.findall(r"^\s+([0-9]{4}) :kaprekar \(([^)]*)\)", answer, flags=re.MULTILINE):
        rows[int(start)] = [int(v) for v in values.split()]
    return rows


def run(ctx):
    data = ctx.load_input()
    start_count = int(data["StartCount"])
    target = int(data["TargetConstant"])
    zero = int(data["ZeroBasin"])
    max_steps = int(data["MaxKaprekarSteps"])

    emitted = {}
    omitted = {}
    for n in range(start_count):
        ch = chain(n, max_steps, target, zero)
        if ch[-1] == target:
            emitted[n] = ch
        elif ch[-1] == zero:
            omitted[n] = ch
    distances = {n: (0 if n == target else len(ch)) for n, ch in emitted.items()}
    distribution = Counter(distances.values())
    selected = parse_selected(ctx.answer)

    identity_ok = all(step(n) == int("".join(sorted(f"{n:04d}", reverse=True))) - int("".join(sorted(f"{n:04d}"))) for n in range(start_count))

    checks = [
        ("all four-digit starts from 0000 through 9999 are considered", start_count == 10000 and len(emitted) + len(omitted) == start_count),
        ("the optimized Kaprekar step equals direct descending-minus-ascending digit sorting", identity_ok),
        ("the classic 3524 chain is recomputed independently", emitted[3524] == [3087, 8352, 6174]),
        ("0001 is treated as padded four-digit input", emitted[1] == [999, 8991, 8082, 8532, 6174]),
        ("0000 and the nine nonzero repdigits fall to the zero basin", set(omitted) == {0, 1111, 2222, 3333, 4444, 5555, 6666, 7777, 8888, 9999}),
        ("every emitted chain reaches 6174 within the configured bound", all(ch[-1] == target and len(ch) <= max_steps + 1 for ch in emitted.values())),
        ("the recomputed maximum step count is seven", max(distances.values()) == max_steps),
        ("reported emitted and omitted counts match recomputation", parse_int(ctx.answer, "total emitted") == len(emitted) and parse_int(ctx.answer, "omitted 0000 basin") == len(omitted)),
        ("the step-count distribution in the explanation matches recomputation", parse_distribution(ctx.reason) == dict(sorted(distribution.items()))),
        ("selected reported chains match the recomputed chains", all(selected[n] == emitted[n] for n in selected if n in emitted)),
    ]
    return run_checks(checks)
