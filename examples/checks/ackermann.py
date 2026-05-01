# Independent Python checks for the ackermann example.
from __future__ import annotations

import hashlib
import re
import sys

from .common import run_checks

try:
    sys.set_int_max_str_digits(0)
except AttributeError:
    pass


def summarize(n: int) -> tuple[int, str, str, str]:
    s = str(n)
    return len(s), s[:24], s[-24:], hashlib.sha256(s.encode("ascii")).hexdigest()


def ackermann_binary(x: int, y: int) -> int:
    # This example uses A(x,y)=T(x,y+3,2)-3, where the T rules are the
    # base-2 hyperoperation ladder used in the N3 source.
    t_y = y + 3
    if x == 0:
        value = t_y + 1
    elif x == 1:
        value = t_y + 2
    elif x == 2:
        value = 2 * t_y
    elif x == 3:
        value = 1 << t_y
    elif x == 4:
        if t_y == 3:
            value = 16
        elif t_y == 4:
            value = 65536
        elif t_y == 5:
            value = 1 << 65536
        else:
            raise ValueError("unsupported tetration query")
    elif x == 5 and t_y == 3:
        value = 65536
    else:
        raise ValueError("unsupported Ackermann query")
    return value - 3


def parse_small(answer: str) -> dict[str, int]:
    rows = {}
    for ident, value in re.findall(r"^(A\d+) ackermann\(\d+,\d+\) = (\d+)$", answer, flags=re.MULTILINE):
        rows[ident] = int(value)
    return rows


def parse_fingerprints(answer: str) -> dict[str, tuple[int, str, str, str]]:
    rows = {}
    for ident, digits, leading, trailing, sha in re.findall(r"^(A\d+) digits=(\d+) leading=([0-9]+) trailing=([0-9]+) sha256=([0-9a-f]+)$", answer, flags=re.MULTILINE):
        rows[ident] = (int(digits), leading, trailing, sha)
    return rows


def run(ctx):
    data = ctx.load_input()
    values = {item["ID"]: ackermann_binary(int(item["X"]), int(item["Y"])) for item in data}
    small = parse_small(ctx.answer)
    fingerprints = parse_fingerprints(ctx.answer)
    a7_summary = summarize(values["A7"])
    a10_summary = summarize(values["A10"])

    checks = [
        ("all twelve JSON Ackermann queries are recomputed", len(values) == 12),
        ("x=0 queries reduce to successor after the +3 binary offset", values["A0"] == 1 and values["A1"] == 7),
        ("x=1 queries reduce to addition after the +3 binary offset", values["A2"] == 4 and values["A3"] == 9),
        ("x=2 queries reduce to multiplication after the +3 binary offset", values["A4"] == 7 and values["A5"] == 21),
        ("x=3 queries reduce to exact base-2 exponentiation", values["A6"] == 125 and a7_summary == fingerprints.get("A7")),
        ("A(4,0) and A(4,1) match the first tetration cases", values["A8"] == 13 and values["A9"] == 65533),
        ("A(4,2) is held exactly as 2^65536-3 with the reported fingerprint", a10_summary == fingerprints.get("A10")),
        ("A(5,0) lands on the same value as A(4,1)", values["A11"] == values["A9"] == 65533),
        ("all non-huge reported exact values match recomputation", all(small[k] == values[k] for k in small)),
        ("the reported proof statistics match the query and memo structure", "primitive test queries : 12" in ctx.reason and "binary reductions : 12" in ctx.reason and "distinct ternary facts : 23" in ctx.reason),
    ]
    return run_checks(checks)
