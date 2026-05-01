# Independent Python checks for the complex_matrix_stability example.
from __future__ import annotations

import math
import re

from .common import run_checks


def modulus(z: dict) -> float:
    return math.hypot(float(z["re"]), float(z["im"]))


def classify(radius: float) -> str:
    if radius < 1.0:
        return "damped"
    if abs(radius - 1.0) <= 1e-12:
        return "marginally stable"
    return "unstable"


def parse_rows(answer: str) -> dict[str, tuple[float, str]]:
    rows = {}
    for name, radius, label in re.findall(r"^([A-Za-z0-9_]+) : spectral radius ([0-9.]+) -> ([^\n]+)$", answer, flags=re.MULTILINE):
        rows[name] = (float(radius), label.strip())
    return rows


def run(ctx):
    data = ctx.load_input()
    rows = parse_rows(ctx.answer)
    radii = {
        matrix["name"]: max(modulus(z) for z in matrix["diagonal"])
        for matrix in data["matrices"]
    }
    z = complex(float(data["sampleProduct"]["z"]["re"]), float(data["sampleProduct"]["z"]["im"]))
    w = complex(float(data["sampleProduct"]["w"]["re"]), float(data["sampleProduct"]["w"]["im"]))
    scale = float(data["scale"])

    checks = [
        ("diagonal entries are used as the eigenvalues", set(radii) == {m["name"] for m in data["matrices"]}),
        ("A_unstable has independently recomputed spectral radius 2", abs(radii["A_unstable"] - 2.0) <= 1e-12 and classify(radii["A_unstable"]) == "unstable"),
        ("A_stable has spectral radius exactly 1 and is marginal", abs(radii["A_stable"] - 1.0) <= 1e-12 and classify(radii["A_stable"]) == "marginally stable"),
        ("A_damped has spectral radius 0 and is damped", radii["A_damped"] == 0.0 and classify(radii["A_damped"]) == "damped"),
        ("reported matrix classes and radii match recomputation", set(rows) == set(radii) and all(abs(rows[name][0] - radii[name]) <= 1e-12 and rows[name][1] == classify(radii[name]) for name in rows)),
        ("squared modulus of z*w equals product of squared moduli", abs(abs(z * w) ** 2 - (abs(z) ** 2 * abs(w) ** 2)) <= 1e-12),
        ("scaling a matrix by 2 multiplies spectral-radius-squared by 4", abs((scale * radii["A_unstable"]) ** 2 - (scale ** 2) * radii["A_unstable"] ** 2) <= 1e-12),
    ]
    return run_checks(checks)
