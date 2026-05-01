# Independent Python checks for the bayes_diagnosis example.
from __future__ import annotations

import math
import re

from .common import run_checks


def parse_distribution(answer: str) -> dict[str, tuple[float, float]]:
    rows = {}
    for disease, posterior, unnorm in re.findall(r"^\s*([A-Za-z0-9_]+)\s+posterior = ([0-9.]+)\s+\(unnormalized = ([0-9.]+)\)", answer, flags=re.MULTILINE):
        rows[disease] = (float(posterior), float(unnorm))
    return rows


def run(ctx):
    data = ctx.load_input()
    priors = {d["Name"]: float(d["Prior"]) for d in data["Diseases"]}
    conditionals = data["ProbGiven"]
    evidence = data["Evidence"]
    raw = {}
    for disease, prior in priors.items():
        likelihood = prior
        for item in evidence:
            p = float(conditionals[disease][item["Symptom"]])
            likelihood *= p if item["Present"] else (1.0 - p)
        raw[disease] = likelihood
    total = sum(raw.values())
    posterior = {d: raw[d] / total for d in raw}
    reported = parse_distribution(ctx.answer)
    best = max(posterior, key=posterior.get)
    total_match = re.search(r"Evidence total \(normalizing constant\) = ([0-9]+(?:\.[0-9]+)?)", ctx.reason)
    reported_total = float(total_match.group(1)) if total_match else None

    checks = [
        ("all priors are probabilities and the prior mass is less than one", all(0 <= p <= 1 for p in priors.values()) and sum(priors.values()) < 1),
        ("every conditional probability is in [0, 1]", all(0 <= float(p) <= 1 for probs in conditionals.values() for p in probs.values())),
        ("all evidence symptoms are available for every disease", all(item["Symptom"] in conditionals[d] for d in priors for item in evidence)),
        ("the absent Sneezing evidence uses the complement likelihood", math.isclose(raw["COVID19"], 0.05 * 0.7 * 0.65 * 0.4 * (1 - 0.15) * 0.2, rel_tol=0, abs_tol=1e-15)),
        ("the Bayesian normalizing constant is recomputed independently", reported_total is not None and abs(reported_total - total) < 5e-9),
        ("the reported distribution contains one posterior for each disease", set(reported) == set(priors)),
        ("each reported unnormalized likelihood matches the Python recomputation", all(abs(reported[d][1] - raw[d]) < 5e-8 for d in reported)),
        ("each reported posterior matches likelihood divided by evidence total", all(abs(reported[d][0] - posterior[d]) < 5e-6 for d in reported)),
        ("the reported posteriors sum to one after rounding", abs(sum(v[0] for v in reported.values()) - 1.0) < 2e-6),
        ("COVID19 is independently selected as the maximum-posterior disease", best == "COVID19" and "most likely disease is COVID19" in ctx.answer),
    ]
    return run_checks(checks)
