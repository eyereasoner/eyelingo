# Independent Python checks for the doctor_advice_work_conflict example.
from __future__ import annotations

import re

from .common import run_checks


def closure(job: str, subclass_of: dict[str, str]) -> set[str]:
    out = {job}
    while job in subclass_of:
        job = subclass_of[job]
        out.add(job)
    return out


def evaluate(data: dict) -> dict[str, dict[str, str]]:
    sick = data["Person"]["Condition"] == "Flu"
    doctor_jobs = set(data["DoctorCanDoJobs"])
    policies = data["ConflictPolicies"]
    results = {}
    for req in data["Requests"]:
        raw = set()
        job_family = closure(req["Job"], data["SubclassOf"])
        if req["Job"] in doctor_jobs:
            raw.add("Permit")
        if sick and "Work" in job_family and policies["SickWorkDefault"] == "Deny":
            raw.add("Deny")
        status = "BothPermitDeny" if raw == {"Permit", "Deny"} else next(iter(raw), "None")
        if status == "BothPermitDeny" and req["Location"] == "Home" and policies["HomeProgrammingWork"] == "Permit":
            effective = "Permit"
        elif status == "BothPermitDeny" and req["Location"] == "Office" and policies["SickOfficeWorkDefault"] == "Deny":
            effective = "Deny"
        else:
            effective = status
        results[req["ID"]] = {"raw": "+".join(sorted(raw)), "status": status, "effective": effective}
    return results


def parse_rows(answer: str) -> dict[str, dict[str, str]]:
    rows = {}
    for req_id, raw, status, effective in re.findall(r"^(Request_[A-Za-z0-9_]+) : raw=([A-Za-z+]+) status=([A-Za-z]+) effective=([A-Za-z]+)", answer, flags=re.MULTILINE):
        rows[req_id] = {"raw": raw, "status": status, "effective": effective}
    return rows


def run(ctx):
    data = ctx.load_input()
    computed = evaluate(data)
    reported = parse_rows(ctx.answer)
    home = next(req["ID"] for req in data["Requests"] if req["Location"] == "Home")
    office = next(req["ID"] for req in data["Requests"] if req["Location"] == "Office")
    overall = "RemoteWorkOnly" if computed[home]["effective"] == "Permit" and computed[office]["effective"] == "Deny" else "Other"

    checks = [
        ("Flu classifies Jos as sick for the policy conflict", data["Person"]["Condition"] == "Flu" and "Jos has Flu" in ctx.reason),
        ("ProgrammingWork is closed upward to Work", "Work" in closure("ProgrammingWork", data["SubclassOf"])),
        ("doctor advice contributes Permit for every ProgrammingWork request", all("Permit" in computed[req["ID"]]["raw"] for req in data["Requests"])),
        ("sick-work default contributes Deny for both requests", all("Deny" in computed[req["ID"]]["raw"] for req in data["Requests"])),
        ("the home request keeps the raw Permit+Deny conflict before resolution", reported.get(home) == computed[home] and computed[home]["status"] == "BothPermitDeny"),
        ("the office request keeps the raw Permit+Deny conflict before resolution", reported.get(office) == computed[office] and computed[office]["status"] == "BothPermitDeny"),
        ("conflict resolution permits sick programming work at Home", computed[home]["effective"] == "Permit"),
        ("conflict resolution denies Office work", computed[office]["effective"] == "Deny"),
        ("the combined recommendation recomputes to RemoteWorkOnly", overall == data["Expected"]["OverallDecision"] and f"overall decision for {data['Person']['Name']} : {overall}" in ctx.answer),
    ]
    return run_checks(checks)
