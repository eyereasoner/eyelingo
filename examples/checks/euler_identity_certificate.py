# Independent Python checks for the euler_identity_certificate example.
from __future__ import annotations

import math
import re

from .common import run_checks


def exp_i_taylor(angle: float, terms: int) -> complex:
    z = 1j * angle
    term = 1 + 0j
    total = 1 + 0j
    for k in range(1, terms):
        term *= z / k
        total += term
    return total


def parse_float(answer: str, label: str) -> float | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([-+0-9.eE]+)", answer)
    return float(match.group(1)) if match else None


def run(ctx):
    data = ctx.load_input()
    angle = float(data["angle"])
    terms = int(data["terms"])
    tolerance = float(data["tolerance"])
    value = exp_i_taylor(angle, terms)
    residual = abs(value + 1)
    reported_real = parse_float(ctx.answer, "computed real part of exp(iπ)")
    reported_imag = parse_float(ctx.answer, "computed imaginary part of exp(iπ)")
    reported_residual = parse_float(ctx.answer, "residual magnitude")

    checks = [
        ("the Taylor expansion uses the terms count from JSON", terms == 28 and f"terms used : {terms}" in ctx.answer),
        ("the input angle is pi to floating precision", abs(angle - math.pi) <= 1e-15),
        ("the finite Taylor real part is close to -1", abs(value.real + 1.0) <= tolerance),
        ("the finite Taylor imaginary part is close to zero", abs(value.imag) <= tolerance),
        ("the residual |exp(iπ)+1| is below the configured tolerance", residual < tolerance and data["expected"]["residualBelowTolerance"] is True),
        ("reported real, imaginary, and residual values match recomputation", reported_real is not None and reported_imag is not None and reported_residual is not None and abs(reported_real - value.real) <= 1e-12 and abs(reported_imag - value.imag) <= 1e-12 and abs(reported_residual - residual) <= 5e-16),
        ("the explanation explicitly treats the result as a finite certificate", "finite" in ctx.reason and "not claimed" in ctx.reason and "tolerance" in ctx.reason),
    ]
    return run_checks(checks)
