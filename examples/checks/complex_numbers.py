# Independent Python checks for the complex_numbers example.
from __future__ import annotations

import cmath
import math
import re

from .common import run_checks


def to_complex(pair: dict) -> complex:
    return complex(float(pair["Re"]), float(pair["Im"]))


def parse_answer_values(answer: str) -> dict[str, complex]:
    values = {}
    pattern = r"^(C[0-9]+).*?=\s+([-0-9.]+)\s+([+-])\s+([0-9.]+)i\s*$"
    for ident, real, sign, imag in re.findall(pattern, answer, flags=re.MULTILINE):
        values[ident] = complex(float(real), float(imag) * (1 if sign == "+" else -1))
    return values


def close(a: complex, b: complex, tol: float = 5e-12) -> bool:
    return abs(a.real - b.real) <= tol and abs(a.imag - b.imag) <= tol


def run(ctx):
    data = ctx.load_input()
    expected = {}
    recomputed = {}
    polar_angles = {}
    for item in data["Exponents"]:
        ident = item["ID"]
        base = to_complex(item["Base"])
        power = to_complex(item["Power"])
        expected[ident] = to_complex(item["Expected"])
        recomputed[ident] = base ** power
        polar_angles[ident] = cmath.phase(base)
    for item in data["Inverses"]:
        ident = item["ID"]
        z = to_complex(item["Input"])
        expected[ident] = to_complex(item["Expected"])
        recomputed[ident] = cmath.asin(z) if item["Operation"] == "asin" else cmath.acos(z)

    reported = parse_answer_values(ctx.answer)
    asin_value = recomputed["C5"]
    acos_value = recomputed["C6"]

    checks = [
        ("principal polar angles for -1, e, and i match the N3 dial cases", abs(polar_angles["C1"] - math.pi) <= 1e-12 and abs(polar_angles["C2"]) <= 1e-12 and abs(polar_angles["C3"] - math.pi / 2) <= 1e-12),
        ("all four complex exponentiation answers match independent complex arithmetic", all(close(recomputed[item["ID"]], expected[item["ID"]]) for item in data["Exponents"])),
        ("i^i and e^(-pi/2) recompute to the same real value", abs(recomputed["C3"] - recomputed["C4"]) <= 1e-12 and abs(recomputed["C3"].imag) <= 1e-12),
        ("asin(2) and acos(2) match independent inverse-trig recomputation", close(recomputed["C5"], expected["C5"]) and close(recomputed["C6"], expected["C6"])),
        ("sin(asin(2)) and cos(acos(2)) round-trip to 2+0i", abs(cmath.sin(asin_value) - 2) <= 1e-12 and abs(cmath.cos(acos_value) - 2) <= 1e-12),
        ("asin(2) + acos(2) equals pi/2 with cancelling imaginary parts", abs((asin_value + acos_value).real - math.pi / 2) <= 1e-12 and abs((asin_value + acos_value).imag) <= 1e-12),
        ("all six reported complex outputs match recomputation to displayed precision", set(reported) == set(recomputed) and all(abs(reported[k] - recomputed[k]) <= 5e-12 for k in reported)),
        ("the report contains four exponentiation and two inverse-trig queries", len(data["Exponents"]) == 4 and len(data["Inverses"]) == 2 and "primitive test facts : 6" in ctx.reason),
        ("all recomputed outputs are finite real/imaginary pairs", all(math.isfinite(z.real) and math.isfinite(z.imag) for z in recomputed.values())),
    ]
    return run_checks(checks)
