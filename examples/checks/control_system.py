# Independent Python checks for the control_system example.
from __future__ import annotations

import math
import re

from .common import run_checks


def measurement10(pair):
    a, b = pair
    if a < b:
        return math.sqrt(b - a), "lessThan"
    return a, "notLessThan"


def parse_answer(answer: str) -> dict[str, float]:
    values = {}
    patterns = {
        "actuator1": r"actuator1 control1 = (-?[0-9.]+)",
        "actuator2": r"actuator2 control1 = (-?[0-9.]+)",
        "input1_m10": r"input1 measurement10 = (-?[0-9.]+)",
        "disturbance2_m10": r"disturbance2 measurement10 = (-?[0-9.]+)",
    }
    for key, pattern in patterns.items():
        m = re.search(pattern, answer)
        values[key] = float(m.group(1)) if m else None
    return values


def close_reported(value, expected, places=6):
    return value is not None and abs(value - expected) <= 0.5 * 10 ** (-places)


def run(ctx):
    data = ctx.load_input()
    input1_m10, input_branch = measurement10(data["input1"]["measurement1"])
    disturb2_m10, disturb_branch = measurement10(data["disturbance2"]["measurement1"])
    product = input1_m10 * 19.6
    compensation = math.log10(data["disturbance1"]["measurement3"])
    actuator1 = product - compensation
    error = data["output2"]["target2"] - data["output2"]["measurement4"]
    differential = data["state3"]["observation3"] - data["output2"]["measurement4"]
    proportional = 5.8 * error
    nonlinear_factor = 7.3 / error
    nonlinear = nonlinear_factor * differential
    actuator2 = proportional + nonlinear
    reported = parse_answer(ctx.answer)

    checks = [
        ("input1 measurement10 recomputes through the lessThan square-root branch", input_branch == "lessThan" and close_reported(reported["input1_m10"], input1_m10)),
        ("disturbance2 measurement10 recomputes through the notLessThan branch", disturb_branch == "notLessThan" and close_reported(reported["disturbance2_m10"], disturb2_m10)),
        ("feedforward guard is true before actuator1 arithmetic is applied", data["input2"]["measurement2"] is True),
        ("actuator1 control recomputes product minus log10 compensation", close_reported(reported["actuator1"], actuator1)),
        ("target-minus-measurement error is recomputed as 5", error == 5),
        ("state/output differential error is recomputed as -2", differential == -2),
        ("nonlinear feedback term uses 7.3/error times the differential", math.isclose(nonlinear, -2.92, rel_tol=0, abs_tol=1e-12)),
        ("actuator2 control recomputes proportional plus nonlinear feedback", close_reported(reported["actuator2"], actuator2)),
    ]
    return run_checks(checks)
