# Independent Python checks for the ebike_motor_thermal_envelope example.
from __future__ import annotations

import math
import re

from .common import run_checks


def propagate(data: dict):
    lower = upper = data["InitialMotorC"] - data["AmbientC"]
    rows = []
    for i, assist in enumerate(data["AssistPlan"], 1):
        heat = data["HeatingEnvelopeByAssist"][assist]
        lower = lower * data["CoolingLower"] + heat["Lower"]
        upper = upper * data["CoolingUpper"] + heat["Upper"]
        rows.append({"step": i, "assist": assist, "lower_temp": data["AmbientC"] + lower, "upper_temp": data["AmbientC"] + upper})
    return rows


def parse_answer(answer: str):
    vals = {}
    for key, pattern in {
        "max_upper": r"maximum upper motor temperature : ([0-9.]+) C",
        "recovery_step": r"warning recovery : step (\d+) at",
        "hard_limit": r"hard limit : ([0-9.]+) C",
    }.items():
        m = re.search(pattern, answer)
        vals[key] = float(m.group(1)) if m and key != "recovery_step" else (int(m.group(1)) if m else None)
    return vals


def run(ctx):
    data = ctx.load_input()
    rows = propagate(data)
    reported = parse_answer(ctx.answer)
    max_upper = max(r["upper_temp"] for r in rows)
    first_warning = next(i for i, r in enumerate(rows, 1) if r["upper_temp"] > data["WarningLimitC"])
    recovery_step = next(i for i, r in enumerate(rows, 1) if i >= first_warning and r["upper_temp"] <= data["WarningLimitC"])
    hard_safe = all(r["upper_temp"] < data["HardLimitC"] for r in rows)
    exact_alpha = math.exp(-data["SamplePeriodSec"] / data["ThermalTimeConstantSec"])

    checks = [
        ("the cooling certificate brackets exp(-sample/tau)", data["CoolingLower"] <= exact_alpha <= data["CoolingUpper"]),
        ("the assist plan has twelve sampled thermal updates", len(rows) == 12),
        ("Turbo, Tour, Eco, and Coast heating envelopes are nonnegative intervals", all(v["Lower"] <= v["Upper"] and v["Lower"] >= 0 for v in data["HeatingEnvelopeByAssist"].values())),
        ("interval propagation recomputes the maximum upper motor temperature", abs(reported["max_upper"] - max_upper) < 5e-4),
        ("the upper envelope first exceeds the warning limit during Turbo", first_warning == 1 and rows[0]["assist"] == "Turbo"),
        ("the reported warning-recovery step matches the independently propagated envelope", reported["recovery_step"] == recovery_step == data["Expected"]["WarningRecoveryStep"]),
        ("all upper temperatures remain below the hard thermal limit", hard_safe and reported["hard_limit"] == data["HardLimitC"]),
        ("the final Coast samples cool monotonically", all(rows[i]["upper_temp"] < rows[i - 1]["upper_temp"] for i in range(8, 12))),
        ("the final decision matches the safety envelope", data["Expected"]["Decision"] in ctx.answer and hard_safe),
    ]
    return run_checks(checks)
