# Independent Python checks for the odrl_dpv_risk_ranked example.
from __future__ import annotations

import re

from .common import run_checks


def notice_days(permission: dict) -> int | None:
    for duty in permission.get("Duties") or []:
        if duty.get("Action") == "odrl:inform":
            for constraint in duty.get("Constraints") or []:
                if constraint.get("LeftOperand") == "tosl:noticeDays":
                    return constraint["RightOperand"]["Int"]
    return None


def has_consent_constraint(permission: dict) -> bool:
    return any(c.get("LeftOperand") in {"dpv:Consent", "tosl:explicitConsent"} for c in (permission.get("Constraints") or []))


def risk_rows(data: dict):
    rows = []
    permissions = data["Agreement"]["Policy"]["Permissions"]
    prohibitions = data["Agreement"]["Policy"]["Prohibitions"]
    clauses = data["Agreement"]["Clauses"]
    needs = data["Consumer"]["Needs"]

    delete = permissions[":PermDeleteAccount"]
    if delete["Action"] == "tosl:removeAccount" and notice_days(delete) is None:
        rows.append({"clause": clauses[delete["ClauseID"]]["ID"], "score": 100, "level": "HighRisk", "mitigations": 2})

    share = permissions[":PermShareData"]
    if share["Action"] == "tosl:shareData" and not has_consent_constraint(share):
        rows.append({"clause": clauses[share["ClauseID"]]["ID"], "score": 97, "level": "HighRisk", "mitigations": 1})

    change = permissions[":PermChangeTerms"]
    got_notice = notice_days(change)
    required_notice = needs[":Need_ChangeOnlyWithPriorNotice"]["MinNoticeDays"]
    if got_notice is not None and got_notice < required_notice:
        rows.append({"clause": clauses[change["ClauseID"]]["ID"], "score": 85, "level": "HighRisk", "mitigations": 1})

    export = prohibitions[":ProhibitExportData"]
    if export["Action"] == "tosl:exportData":
        rows.append({"clause": clauses[export["ClauseID"]]["ID"], "score": 70, "level": "ModerateRisk", "mitigations": 1})

    rows.sort(key=lambda r: (-r["score"], r["clause"]))
    return rows


def parse_rows(answer: str):
    return [{"score": int(score), "level": level, "clause": clause} for score, level, clause in re.findall(r"score=(\d+) \(risk:([A-Za-z]+), risk:[A-Za-z]+\) clause (C\d)", answer)]


def run(ctx):
    data = ctx.load_input()
    rows = risk_rows(data)
    reported = parse_rows(ctx.answer)
    scores = [r["score"] for r in rows]
    mitigation_count = len(re.findall(r"^mitigation for clause", ctx.answer, flags=re.MULTILINE))

    checks = [
        ("four risk rows are derived from the policy/profile conflict scan", len(rows) == 4 and len(reported) == 4),
        ("reported rows match independently computed clauses and scores", [(r["clause"], r["score"]) for r in reported] == [(r["clause"], r["score"]) for r in rows]),
        ("ranked output is in descending score order", scores == sorted(scores, reverse=True)),
        ("account/data removal without notice safeguards is the highest risk", rows[0]["clause"] == "C1" and rows[0]["score"] == 100),
        ("user-data sharing without explicit consent is scored as high risk", any(r["clause"] == "C3" and r["score"] == 97 for r in rows)),
        ("three-day terms-change notice is below the fourteen-day consumer requirement", notice_days(data["Agreement"]["Policy"]["Permissions"][":PermChangeTerms"]) == 3 and data["Consumer"]["Needs"][":Need_ChangeOnlyWithPriorNotice"]["MinNoticeDays"] == 14),
        ("data-export prohibition creates the portability risk row", any(r["clause"] == "C4" and r["score"] == 70 and r["level"] == "ModerateRisk" for r in rows)),
        ("risk level counts recompute to high=3, moderate=1, low=0", sum(r["level"] == "HighRisk" for r in rows) == 3 and sum(r["level"] == "ModerateRisk" for r in rows) == 1),
        ("five mitigation measures are generated", mitigation_count == sum(r["mitigations"] for r in rows) == 5),
    ]
    return run_checks(checks)
