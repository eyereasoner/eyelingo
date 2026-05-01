# Independent Python checks for the bayes_therapy example.
from __future__ import annotations

import math
import re

from .common import run_checks


def parse_recommendation(answer: str) -> tuple[str | None, float | None]:
    match = re.search(r"recommended therapy is ([^(]+) \(utility = ([0-9.\-]+)\)", answer)
    if not match:
        return None, None
    return match.group(1).strip(), float(match.group(2))


def parse_posteriors(answer: str) -> dict[str, tuple[float, float]]:
    rows = {}
    pattern = r"^\s*([A-Za-z0-9_]+)\s+posterior = ([0-9.]+)\s+\(unnormalized = ([0-9.]+)\)"
    for disease, posterior, unnormalized in re.findall(pattern, answer, flags=re.MULTILINE):
        rows[disease] = (float(posterior), float(unnormalized))
    return rows


def parse_therapies(answer: str) -> dict[str, tuple[float, float, float]]:
    rows = {}
    pattern = r"^\s*([A-Za-z0-9_]+)\s+expectedSuccess = ([0-9.]+)\s+adverse = ([0-9.]+)\s+utility = ([0-9.\-]+)"
    for therapy, success, adverse, utility in re.findall(pattern, answer, flags=re.MULTILINE):
        rows[therapy] = (float(success), float(adverse), float(utility))
    return rows


def evidence_total_from_reason(reason: str) -> float | None:
    match = re.search(r"Evidence total \(normalizing constant\) = ([0-9]+(?:\.[0-9]+)?)", reason)
    return float(match.group(1)) if match else None


def recompute_posteriors(data: dict):
    diseases = [d["Name"] for d in data["Diseases"]]
    priors = {d["Name"]: float(d["Prior"]) for d in data["Diseases"]}
    conditionals = data["ProbGiven"]
    raw = {}
    factors = {}
    for disease in diseases:
        likelihood = priors[disease]
        disease_factors = []
        for item in data["Evidence"]:
            p = float(conditionals[disease][item["Symptom"]])
            factor = p if item["Present"] else 1.0 - p
            disease_factors.append((item["Symptom"], item["Present"], factor))
            likelihood *= factor
        raw[disease] = likelihood
        factors[disease] = disease_factors
    total = sum(raw.values())
    posteriors = {disease: raw[disease] / total for disease in diseases}
    return diseases, priors, raw, posteriors, total, factors


def recompute_therapies(data: dict, diseases: list[str], posteriors: dict[str, float]):
    benefit = float(data["BenefitWeight"])
    harm = float(data["HarmWeight"])
    results = {}
    for therapy in data["Therapies"]:
        success_by_disease = [float(value) for value in therapy["SuccessByDisease"]]
        expected_success = sum(posteriors[disease] * success_by_disease[index] for index, disease in enumerate(diseases))
        adverse = float(therapy["Adverse"])
        utility = benefit * expected_success - harm * adverse
        results[therapy["Name"]] = (expected_success, adverse, utility)
    return results


def run(ctx):
    data = ctx.load_input()
    diseases, priors, raw, posteriors, total, factors = recompute_posteriors(data)
    therapy_results = recompute_therapies(data, diseases, posteriors)
    reported_posteriors = parse_posteriors(ctx.answer)
    reported_therapies = parse_therapies(ctx.answer)
    reported_best, reported_best_utility = parse_recommendation(ctx.answer)
    reported_total = evidence_total_from_reason(ctx.reason)
    best = max(therapy_results, key=lambda name: therapy_results[name][2])

    absent_factors_used = all(
        math.isclose(factor, 1.0 - float(data["ProbGiven"][disease][symptom]), rel_tol=0, abs_tol=1e-15)
        for disease, disease_factors in factors.items()
        for symptom, present, factor in disease_factors
        if not present
    )
    disease_count = len(diseases)
    therapies_have_aligned_success_vectors = all(len(t["SuccessByDisease"]) == disease_count for t in data["Therapies"])

    checks = [
        ("all priors are valid probabilities and the model has nonzero prior mass", all(0 <= p <= 1 for p in priors.values()) and sum(priors.values()) > 0),
        ("all evidence symptoms are available for every disease and absent evidence uses complement factors", absent_factors_used and all(item["Symptom"] in data["ProbGiven"][d] for d in diseases for item in data["Evidence"])),
        ("the Bayesian normalizing constant is recomputed independently", reported_total is not None and abs(reported_total - total) < 5e-9),
        ("the reported posterior table contains one row for each disease", set(reported_posteriors) == set(diseases)),
        ("each reported unnormalized disease likelihood matches the Python recomputation", all(abs(reported_posteriors[d][1] - raw[d]) < 5e-8 for d in reported_posteriors)),
        ("each reported posterior equals likelihood divided by the evidence total", all(abs(reported_posteriors[d][0] - posteriors[d]) < 5e-6 for d in reported_posteriors)),
        ("therapy success vectors align with diseases and all therapy probabilities are in [0, 1]", therapies_have_aligned_success_vectors and all(0 <= float(t["Adverse"]) <= 1 and all(0 <= float(p) <= 1 for p in t["SuccessByDisease"]) for t in data["Therapies"])),
        ("each expected therapy success is recomputed as posterior-weighted disease success", set(reported_therapies) == set(therapy_results) and all(abs(reported_therapies[t][0] - therapy_results[t][0]) < 5e-6 for t in reported_therapies)),
        ("each therapy utility applies benefitWeight × expectedSuccess − harmWeight × adverse", all(abs(reported_therapies[t][2] - therapy_results[t][2]) < 5e-6 for t in reported_therapies)),
        ("the reported recommendation is the maximum-utility therapy", reported_best == best and reported_best_utility is not None and abs(reported_best_utility - therapy_results[best][2]) < 5e-6),
    ]
    return run_checks(checks)
