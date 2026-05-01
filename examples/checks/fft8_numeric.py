# Independent Python checks for the fft8_numeric example.
from __future__ import annotations

import cmath
import math
import re

from .common import run_checks


def dft(samples: list[float]) -> list[complex]:
    n = len(samples)
    return [sum(samples[t] * cmath.exp(-2j * math.pi * k * t / n) for t in range(n)) for k in range(n)]


def parse_dominant(answer: str) -> dict[int, tuple[float, float]]:
    rows = {}
    for k, magnitude, phase in re.findall(r"k=([0-9]+) magnitude=([0-9.]+) phase=([-0-9.]+)", answer):
        rows[int(k)] = (float(magnitude), float(phase))
    return rows


def parse_float(answer: str, label: str) -> float | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*([0-9.]+)", answer)
    return float(match.group(1)) if match else None


def run(ctx):
    data = ctx.load_input()
    samples = [float(x) for x in data["samples"]]
    spectrum = dft(samples)
    magnitudes = [abs(z) for z in spectrum]
    max_mag = max(magnitudes)
    tol = float(data["expected"]["tolerance"])
    dominant = [k for k, mag in enumerate(magnitudes) if abs(mag - max_mag) <= tol]
    reported = parse_dominant(ctx.answer)
    time_energy = sum(x * x for x in samples)
    freq_energy = sum(abs(z) ** 2 for z in spectrum) / len(samples)

    checks = [
        ("the input contains exactly eight samples", len(samples) == 8),
        ("dominant bins are recomputed as k=1 and k=7", dominant == data["expected"]["dominantBins"]),
        ("the DC component cancels to zero within tolerance", abs(spectrum[0]) <= tol),
        ("the two sine bins have magnitude four", all(abs(magnitudes[k] - 4.0) <= 1e-9 for k in dominant)),
        ("reported dominant magnitudes and phases match recomputation", set(reported) == set(dominant) and all(abs(reported[k][0] - magnitudes[k]) <= 1e-6 and abs(reported[k][1] - cmath.phase(spectrum[k])) <= 1e-6 for k in reported)),
        ("Parseval energy is preserved under the unnormalized DFT convention", abs(time_energy - freq_energy) <= 1e-9),
        ("reported time and frequency-domain energies match recomputation", abs((parse_float(ctx.answer, "time-domain energy") or -1) - time_energy) <= 1e-6 and abs((parse_float(ctx.answer, "frequency-domain energy / 8") or -1) - freq_energy) <= 1e-6),
        ("real samples produce conjugate-symmetric bins", all(abs(spectrum[k] - spectrum[-k].conjugate()) <= 1e-9 for k in range(1, len(samples)))),
    ]
    return run_checks(checks)
