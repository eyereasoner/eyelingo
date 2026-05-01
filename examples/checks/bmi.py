# Independent Python checks for the bmi example.
from __future__ import annotations

import math
import re

from .common import run_checks


def round_half_up(value: float, places: int) -> float:
    scale = 10 ** places
    return math.floor(value * scale + 0.5) / scale


def classify(bmi: float) -> str:
    if bmi < 18.5:
        return "Underweight"
    if bmi < 25.0:
        return "Normal"
    if bmi < 30.0:
        return "Overweight"
    if bmi < 35.0:
        return "Obesity I"
    if bmi < 40.0:
        return "Obesity II"
    return "Obesity III"


def parse_answer(answer: str):
    bmi_match = re.search(r"BMI = ([0-9]+(?:\.[0-9]+)?)", answer)
    category_match = re.search(r"Category = ([^\n]+)", answer)
    band_match = re.search(r"range is about ([0-9]+(?:\.[0-9]+)?)–([0-9]+(?:\.[0-9]+)?) kg", answer)
    return (
        float(bmi_match.group(1)) if bmi_match else None,
        category_match.group(1).strip() if category_match else None,
        float(band_match.group(1)) if band_match else None,
        float(band_match.group(2)) if band_match else None,
    )


def run(ctx):
    data = ctx.load_input()
    unit = data["UnitSystem"]
    weight = float(data["Weight"])
    height = float(data["Height"])

    if unit == "metric":
        weight_kg = weight
        height_m = height / 100.0
    elif unit == "us":
        weight_kg = weight * 0.45359237
        height_m = height * 0.0254
    else:
        weight_kg = float("nan")
        height_m = float("nan")

    height_sq = height_m * height_m
    bmi = weight_kg / height_sq
    bmi_one_decimal = round_half_up(bmi, 1)
    bmi_two_decimal = round_half_up(bmi, 2)
    healthy_min = round_half_up(18.5 * height_sq, 1)
    healthy_max = round_half_up(24.9 * height_sq, 1)
    actual_bmi, actual_category, actual_min, actual_max = parse_answer(ctx.answer)

    boundary_order = [
        classify(18.49),
        classify(18.50),
        classify(24.99),
        classify(25.00),
        classify(29.99),
        classify(30.00),
        classify(35.00),
        classify(40.00),
    ]

    checks = [
        ("input units are normalized to positive SI kg and m", weight_kg > 0 and height_m > 0 and unit in {"metric", "us"}),
        ("height squared is recomputed independently from the normalized height", abs(height_sq - 3.1684) < 1e-12),
        ("reported BMI matches independent kg/m² computation rounded to one decimal", actual_bmi == bmi_one_decimal),
        ("the unrounded BMI independently rounds to the expected two-decimal value", abs(bmi_two_decimal - 22.72) < 1e-12),
        ("reported category is the independent WHO category for the computed BMI", actual_category == classify(bmi)),
        ("WHO category boundaries are half-open at 18.5, 25, 30, 35, and 40", boundary_order == ["Underweight", "Normal", "Normal", "Overweight", "Overweight", "Obesity I", "Obesity II", "Obesity III"]),
        ("reported healthy-weight lower bound equals BMI 18.5 at the same height", actual_min == healthy_min),
        ("reported healthy-weight upper bound equals BMI 24.9 at the same height", actual_max == healthy_max),
        ("the explanation mentions the same formula that the Python check recomputed", "weight in kilograms divided by height in meters squared" in ctx.reason),
    ]
    return run_checks(checks)
