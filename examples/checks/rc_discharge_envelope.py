# Independent Python checks for the rc_discharge_envelope example.
from __future__ import annotations

import re

from .common import run_checks


def parse_int(answer: str, label: str) -> int | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9]+)", answer)
    return int(match.group(1)) if match else None


def parse_float(answer: str, label: str) -> float | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9.]+)", answer)
    return float(match.group(1)) if match else None


def run(ctx):
    data = ctx.load_input()
    lower = float(data["decayLower"])
    upper = float(data["decayUpper"])
    initial = float(data["initialVoltage"])
    tolerance = float(data["tolerance"])
    sample_period = float(data["samplePeriod"])
    max_step = int(data["maxStep"])
    upper_voltage = {step: initial * (upper ** step) for step in range(max_step + 1)}
    lower_voltage = {step: initial * (lower ** step) for step in range(max_step + 1)}
    first = next(step for step in range(max_step + 1) if upper_voltage[step] < tolerance)

    reported_step = parse_int(ctx.answer, "first below tolerance step")
    reported_time = parse_float(ctx.answer, "first below tolerance time")
    reported_voltage = parse_float(ctx.answer, f"upper voltage at step {first}")

    checks = [
        ("the decay interval is nonempty, positive, and below one", 0 < lower <= upper < 1),
        ("the certified upper voltage envelope decreases at every sample", all(upper_voltage[s + 1] < upper_voltage[s] for s in range(max_step))),
        ("the lower and upper envelopes bracket every compatible decay", all(lower_voltage[s] <= upper_voltage[s] for s in range(max_step + 1))),
        ("step 12 remains above the voltage tolerance", upper_voltage[first - 1] >= tolerance),
        ("step 13 is the first certified below-tolerance sample", first == data["expected"]["firstSettledStep"]),
        ("reported first step, time, and upper voltage match recomputation", reported_step == first and abs((reported_time or -1) - first * sample_period) <= 5e-4 and abs((reported_voltage or -1) - upper_voltage[first]) <= 5e-4),
        ("the report uses the JSON double interval rather than only the exact symbol", f"[{lower:.10f}, {upper:.10f}]" in ctx.answer and data["exactDecaySymbol"] in ctx.answer),
    ]
    return run_checks(checks)
