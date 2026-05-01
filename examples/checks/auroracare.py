# Independent Python checks for the auroracare example.
from __future__ import annotations

import re

from .common import run_checks


def policy_matches(policy: dict, scenario: dict, requester: dict, subject: dict) -> bool:
    purpose = scenario["Purpose"]
    cats = set(scenario["Categories"])
    if policy["Kind"] == "permission":
        if purpose not in (policy["AllowedPurposes"] or []):
            return False
        if policy["RequiredRole"] and scenario["Role"] != policy["RequiredRole"]:
            return False
        if policy["RequiredEnvironment"] and scenario["Environment"] != policy["RequiredEnvironment"]:
            return False
        if policy["RequiredTOM"] and scenario["TOM"] != policy["RequiredTOM"]:
            return False
        if policy["AllowAnyCategories"] and not (cats & set(policy["AllowAnyCategories"])):
            return False
        if policy["RequireAllCategories"] and not set(policy["RequireAllCategories"]).issubset(cats):
            return False
        if policy["UID"] == "urn:policy:primary-care-001" and requester["LinkedTo"] != scenario["SubjectID"]:
            return False
        if policy["UID"] == "urn:policy:research-aurora-diabetes" and purpose not in subject["ConsentAllow"]:
            return False
        return True
    if policy["Kind"] == "prohibition":
        return purpose in (policy["ProhibitedPurposes"] or [])
    return False


def evaluate(data: dict):
    decisions = {}
    traces = {}
    for scenario in data["Scenarios"]:
        requester = data["Requesters"][scenario["RequesterID"]]
        subject = data["Subjects"][scenario["SubjectID"]]
        if scenario["Purpose"] in subject["ConsentDeny"]:
            decisions[scenario["Key"]] = ("DENY", None, "subject_opted_out")
            continue
        prohibited = [p for p in data["Policies"] if p["Kind"] == "prohibition" and policy_matches(p, scenario, requester, subject)]
        if prohibited:
            decisions[scenario["Key"]] = ("DENY", prohibited[0]["UID"], "prohibition")
            continue
        permitted = [p for p in data["Policies"] if p["Kind"] == "permission" and policy_matches(p, scenario, requester, subject)]
        if permitted:
            decisions[scenario["Key"]] = ("PERMIT", permitted[0]["UID"], "permission")
        else:
            decisions[scenario["Key"]] = ("DENY", None, "no_policy")
    return decisions


def parse_rows(answer: str):
    rows = {}
    for key, label, decision, source in re.findall(r"^\s+([A-Z]) – (.*?) : (PERMIT|DENY) \((.*?)\)$", answer, flags=re.MULTILINE):
        rows[key] = (decision, source)
    return rows


def run(ctx):
    data = ctx.load_input()
    computed = evaluate(data)
    reported = parse_rows(ctx.answer)
    permit_count = sum(1 for decision, _uid, _why in computed.values() if decision == "PERMIT")
    deny_count = sum(1 for decision, _uid, _why in computed.values() if decision == "DENY")

    checks = [
        ("all seven scenario decisions are recomputed", len(computed) == 7 and len(reported) == 7),
        ("reported decisions match the independent PDP evaluation", all(reported[k][0] == decision for k, (decision, _uid, _why) in computed.items())),
        ("primary-care access requires clinician role and care-team link", computed["A"][:2] == ("PERMIT", "urn:policy:primary-care-001") and computed["E"][:2] == ("PERMIT", "urn:policy:primary-care-001")),
        ("quality improvement is allowed only with both required categories in the secure environment", computed["B"][:2] == ("PERMIT", "urn:policy:qi-2025-aurora") and computed["C"][0] == "DENY"),
        ("insurance management is denied by the matching prohibition", computed["D"][:2] == ("DENY", "urn:policy:deny-insurance")),
        ("research requires opt-in, anonymisation, and the secure environment", computed["F"][:2] == ("PERMIT", "urn:policy:research-aurora-diabetes")),
        ("AI training is denied because the subject opted out", computed["G"] == ("DENY", None, "subject_opted_out")),
        ("permit and deny counts match the report", permit_count == 4 and deny_count == 3 and "permit count : 4" in ctx.answer and "deny count : 3" in ctx.answer),
    ]
    return run_checks(checks)
