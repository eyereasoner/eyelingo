# Independent Python checks for the lldm example.
from __future__ import annotations

import math
import re

from .common import run_checks


def line_through(p1, p2):
    x1, y1 = p1
    x2, y2 = p2
    slope = (y2 - y1) / (x2 - x1)
    intercept = y1 - slope * x1
    return slope, intercept


def intersection(m1, b1, m2, b2):
    x = (b2 - b1) / (m1 - m2)
    y = m1 * x + b1
    return x, y


def dist(a, b):
    return math.hypot(a[0] - b[0], a[1] - b[1])


def parse_answer(answer: str):
    vals = {}
    for key, pattern in {
        "sl1": r"SL1 = (-?[0-9.]+)",
        "sl3": r"SL3 = SL4 = (-?[0-9.]+)",
        "p5x": r"p5\s+= \((-?[0-9.]+),",
        "p5y": r"p5\s+= \(-?[0-9.]+, (-?[0-9.]+)\)",
        "p6x": r"p6\s+= \((-?[0-9.]+),",
        "p6y": r"p6\s+= \(-?[0-9.]+, (-?[0-9.]+)\)",
        "d53": r"d53 = ([0-9.]+) cm",
        "d64": r"d64 = ([0-9.]+) cm",
        "dcm": r"dCm = (-?[0-9.]+) cm",
    }.items():
        m = re.search(pattern, answer)
        vals[key] = float(m.group(1)) if m else None
    return vals


def close(a, b, eps=5e-4):
    return a is not None and abs(a - b) <= eps


def run(ctx):
    data = ctx.load_input()
    p1 = (data["P1x"], data["P1y"])
    p2 = (data["P2x"], data["P2y"])
    p3 = (data["P3x"], data["P3y"])
    p4 = (data["P4x"], data["P4y"])
    sl1, b1 = line_through(p1, p2)
    sl3 = -1.0 / sl1
    b3 = p3[1] - sl3 * p3[0]
    b4 = p4[1] - sl3 * p4[0]
    p5 = intersection(sl1, b1, sl3, b3)
    p6 = intersection(sl1, b1, sl3, b4)
    d53 = dist(p5, p3)
    d64 = dist(p6, p4)
    dcm = d53 - d64
    reported = parse_answer(ctx.answer)

    checks = [
        ("L1 slope is recomputed from p1 and p2", close(reported["sl1"], sl1, 5e-6)),
        ("L3/L4 slopes are perpendicular to L1", close(reported["sl3"], sl3, 5e-6) and abs(sl1 * sl3 + 1) < 1e-12),
        ("p5 is the analytic intersection of L1 and the perpendicular through p3", close(reported["p5x"], p5[0]) and close(reported["p5y"], p5[1])),
        ("p6 is the analytic intersection of L1 and the perpendicular through p4", close(reported["p6x"], p6[0]) and close(reported["p6y"], p6[1])),
        ("d53 recomputes as the Euclidean distance from p5 to p3", close(reported["d53"], d53, 5e-6)),
        ("d64 recomputes as the Euclidean distance from p6 to p4", close(reported["d64"], d64, 5e-6)),
        ("dCm recomputes as d53 minus d64", close(reported["dcm"], dcm, 5e-6)),
        ("the discrepancy is finite and negative for this geometry", math.isfinite(dcm) and dcm < 0),
        ("the alarm condition follows from |dCm| > 1.25 cm", abs(dcm) > 1.25 and "LLD Alarm = TRUE" in ctx.answer),
    ]
    return run_checks(checks)
