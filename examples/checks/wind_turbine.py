# Independent Python checks for the wind_turbine example.
from __future__ import annotations

import re

from .common import run_checks


def power_for_speed(v: float, cut_in: float, rated: float, cut_out: float, rated_power: float) -> tuple[str, float]:
    if v < cut_in or v >= cut_out:
        return "stopped", 0.0
    if v >= rated:
        return "rated", rated_power
    numerator = v ** 3 - cut_in ** 3
    denominator = rated ** 3 - cut_in ** 3
    return "partial", rated_power * numerator / denominator


def parse_int(answer: str, label: str) -> int | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9]+)", answer)
    return int(match.group(1)) if match else None


def parse_float(answer: str, label: str) -> float | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9.]+)", answer)
    return float(match.group(1)) if match else None


def run(ctx):
    data = ctx.load_input()
    cut_in = float(data["cutInMS"])
    rated = float(data["ratedMS"])
    cut_out = float(data["cutOutMS"])
    rated_power = float(data["ratedPowerMW"])
    hours = float(data["intervalMinutes"]) / 60.0
    speeds = [float(v) for v in data["windSpeedsMS"]]
    results = [power_for_speed(v, cut_in, rated, cut_out, rated_power) for v in speeds]
    classes = [kind for kind, _ in results]
    powers = [p for _, p in results]
    usable = sum(kind != "stopped" for kind in classes)
    rated_count = classes.count("rated")
    stopped = classes.count("stopped")
    energy = sum(powers) * hours

    expected = data["expected"]
    answer_line_ok = all(f"t{i + 1} {speeds[i]:.1f} m/s {classes[i]} {powers[i]:.3f} MW" in ctx.answer for i in range(len(speeds)))

    checks = [
        ("cut-in, rated, and cut-out thresholds are strictly ordered", cut_in < rated < cut_out),
        ("usable intervals are exactly the samples inside the operating envelope", usable == expected["usableIntervals"]),
        ("rated intervals are speeds at or above rated and below cut-out", rated_count == expected["ratedIntervals"]),
        ("stopped intervals are below cut-in or at/above cut-out", stopped == expected["stoppedIntervals"]),
        ("below-rated usable speeds follow the cubic normalized power curve", abs(powers[1] - 0.44008604702915216) <= 1e-9 and abs(powers[2] - 2.586496313329871) <= 1e-9),
        ("total interval energy is recomputed in MWh", abs(energy - 1.5710970600598373) <= 1e-9),
        ("reported usable count and total energy match recomputation", parse_int(ctx.answer, "usable intervals") == usable and abs((parse_float(ctx.answer, "total energy") or -1) - energy) <= 5e-4),
        ("the answer reports every sample classification and power", answer_line_ok),
    ]
    return run_checks(checks)
