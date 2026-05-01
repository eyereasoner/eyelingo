# Independent Python checks for the high_trust_bloom_envelope example.
from __future__ import annotations

import math
import re

from .common import run_checks


def parse_answer(answer: str):
    def grab(pattern, cast=float):
        m = re.search(pattern, answer)
        return cast(m.group(1)) if m else None
    fp = re.search(r"false-positive envelope : ([0-9.]+) \.\. ([0-9.]+)", answer)
    return {
        "lambda": grab(r"lambda : ([0-9.]+)"),
        "fp_low": float(fp.group(1)) if fp else None,
        "fp_high": float(fp.group(2)) if fp else None,
        "extra": grab(r"expected extra exact lookups upper : ([0-9.]+)"),
    }


def run(ctx):
    data = ctx.load_input()
    a = data["Artifact"]
    reported = parse_answer(ctx.answer)
    lam = a["HashFunctions"] * a["CanonicalTripleCount"] / a["BloomBits"]
    fp_low = (1.0 - a["ExpMinusLambdaUpper"]) ** a["HashFunctions"]
    fp_high = (1.0 - a["ExpMinusLambdaLower"]) ** a["HashFunctions"]
    extra_upper = fp_high * a["NegativeLookupsPerBatch"]
    exact_exp = math.exp(-lam)

    checks = [
        ("numeric Bloom parameters are positive", all(a[k] > 0 for k in ["CanonicalTripleCount", "BloomBits", "HashFunctions", "NegativeLookupsPerBatch"])),
        ("canonical graph and SPO index triple counts agree", a["CanonicalTripleCount"] == a["SPOIndexTripleCount"]),
        ("lambda recomputes as k*n/m", abs(reported["lambda"] - lam) < 1e-12 and abs(lam - a["CertifiedLambda"]) < 1e-12),
        ("the exp(-lambda) decimal certificate brackets the exact Python value", a["ExpMinusLambdaLower"] <= exact_exp <= a["ExpMinusLambdaUpper"]),
        ("false-positive envelope is recomputed from the certificate", abs(reported["fp_low"] - fp_low) < 5e-10 and abs(reported["fp_high"] - fp_high) < 5e-10),
        ("false-positive upper bound stays below the policy budget", fp_high < a["FPRateBudget"]),
        ("expected extra exact lookups stay below budget", abs(reported["extra"] - extra_upper) < 5e-3 and extra_upper < a["ExtraExactLookupsBudget"]),
        ("maybe-positive answers must be confirmed against the canonical graph", data["Policies"]["MaybePositivePolicy"] == "ConfirmAgainstCanonicalGraph" and "maybe-positive policy : ConfirmAgainstCanonicalGraph" in ctx.answer),
        ("definite negatives may return absent without exact lookup", data["Policies"]["DefiniteNegativePolicy"] == "ReturnAbsent" and "definite-negative policy : ReturnAbsent" in ctx.answer),
        ("the deployment decision matches the recomputed envelope", data["Expected"]["Decision"] in ctx.answer and fp_high < a["FPRateBudget"] and extra_upper < a["ExtraExactLookupsBudget"]),
    ]
    return run_checks(checks)
