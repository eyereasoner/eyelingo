# Independent Python checks for the fft32_numeric example.
from __future__ import annotations

import cmath
import math
import re

from .common import run_checks


def waveform_samples(spec: dict, n: int) -> list[float]:
    kind = spec["kind"]
    if kind == "alternating":
        return [1.0 if i % 2 == 0 else -1.0 for i in range(n)]
    if kind == "constant":
        return [1.0 for _ in range(n)]
    if kind == "cosine":
        b = int(spec["bin"])
        return [math.cos(2 * math.pi * b * i / n) for i in range(n)]
    if kind == "sine":
        b = int(spec["bin"])
        return [math.sin(2 * math.pi * b * i / n) for i in range(n)]
    if kind == "impulse":
        return [1.0] + [0.0 for _ in range(n - 1)]
    if kind == "ramp":
        return [float(i) for i in range(n)]
    raise ValueError(kind)


def dft(samples: list[float]) -> list[complex]:
    n = len(samples)
    return [sum(samples[t] * cmath.exp(-2j * math.pi * k * t / n) for t in range(n)) for k in range(n)]


def dominant_bins(spectrum: list[complex], tol: float) -> list[int]:
    mags = [abs(z) for z in spectrum]
    peak = max(mags)
    return [k for k, mag in enumerate(mags) if abs(mag - peak) <= tol]


def parse_answer_rows(answer: str) -> dict[str, str]:
    rows = {}
    for line in answer.splitlines():
        if " : " in line and any(token in line for token in ("k=", "all 32 bins")):
            name, rest = line.split(" : ", 1)
            rows[name.strip()] = rest.strip()
    return rows


def run(ctx):
    data = ctx.load_input()
    n = int(data["length"])
    tol = float(data["tolerance"])
    rows = parse_answer_rows(ctx.answer)
    spectra = {}
    samples_by_name = {}
    for spec in data["waveforms"]:
        name = spec["name"]
        samples = waveform_samples(spec, n)
        samples_by_name[name] = samples
        spectra[name] = dft(samples)

    dominant_ok = True
    magnitude_ok = True
    energy_ok = True
    conjugate_ok = True
    report_ok = True
    for spec in data["waveforms"]:
        name = spec["name"]
        spectrum = spectra[name]
        samples = samples_by_name[name]
        mags = [abs(z) for z in spectrum]
        dom = dominant_bins(spectrum, tol)
        dominant_ok = dominant_ok and dom == spec["expectedDominantBins"]
        magnitude_ok = magnitude_ok and all(abs(mags[k] - float(spec["expectedDominantMagnitude"])) <= 1e-8 for k in dom)
        energy_ok = energy_ok and abs(sum(x * x for x in samples) - sum(abs(z) ** 2 for z in spectrum) / n) <= 1e-8
        conjugate_ok = conjugate_ok and all(abs(spectrum[k] - spectrum[-k].conjugate()) <= 1e-8 for k in range(1, n))
        
        if spec.get("expectedFlatSpectrum"):
            report_ok = report_ok and name in rows and "all 32 bins magnitude" in rows[name]
        else:
            report_ok = report_ok and name in rows and all(f"k={k}" in rows[name] for k in spec["expectedDominantBins"])

    impulse_mags = [abs(z) for z in spectra["impulse"]]

    checks = [
        ("each generated waveform has exactly 32 samples", all(len(samples) == n for samples in samples_by_name.values())),
        ("dominant bins match the fixture expectations", dominant_ok),
        ("dominant magnitudes match the configured certificates", magnitude_ok),
        ("Parseval energy is preserved for every waveform", energy_ok),
        ("real-valued waveforms produce conjugate-symmetric spectra", conjugate_ok),
        ("the impulse waveform has a flat unit-magnitude spectrum", all(abs(mag - 1.0) <= 1e-10 for mag in impulse_mags)),
        ("the answer reports every expected dominant bin by waveform name", report_ok),
        ("the configured operation certificate matches six full 32-bin direct DFTs", int(data["expectedOperations"]) == len(data["waveforms"]) * n * n),
    ]
    return run_checks(checks)
