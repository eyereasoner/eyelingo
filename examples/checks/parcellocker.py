# Independent Python checks for the parcellocker example.
from __future__ import annotations

import re

from .common import run_checks


def evaluate(data: dict, request: dict):
    auth = data["Authorization"]
    parcel = data["Parcel"]
    checks = {
        "C1": request["Requester"] == auth["Delegate"],
        "C2": request["Parcel"] == auth["Parcel"],
        "C3": request["Locker"] == auth["Locker"],
        "C4": request["Action"] == auth["Action"],
        "C5": request["Purpose"] == auth["Purpose"],
        "C6": auth["State"] == "Active" and not request["UsedOnce"],
        "C7": auth["Reuse"] == "SingleUse",
        "C8": parcel["Status"] == "ReadyForPickup",
        "C9": not (request["Action"] == "ViewBillingDetails" or request["Purpose"] == "BillingAccess") or auth["BillingAccess"] != "None",
        "C10": not (request["Action"] == "RedirectParcel" or request["Purpose"] == "RedirectDelivery") or auth["RedirectAllowed"] != "No",
    }
    decision = "PERMIT" if all(checks.values()) else "DENY"
    return decision, checks


def parse_decisions(answer: str):
    rows = {}
    for key, decision in re.findall(r"^\s+([a-z-]+)\s+: (PERMIT|DENY)", answer, flags=re.MULTILINE):
        rows[key] = decision
    return rows


def run(ctx):
    data = ctx.load_input()
    computed = {r["Key"]: evaluate(data, r) for r in data["Requests"]}
    decisions = {k: v[0] for k, v in computed.items()}
    reported = parse_decisions(ctx.answer)
    pickup_checks = computed["pickup"][1]
    guardrail_denials = sum(1 for key, decision in decisions.items() if key != "pickup" and decision == "DENY")

    checks = [
        ("the source pickup request satisfies all ten authorization conditions", decisions["pickup"] == "PERMIT" and sum(pickup_checks.values()) == 10),
        ("all reported request decisions match independent policy evaluation", reported == decisions),
        ("billing access is denied by the privacy guardrail", decisions["billing"] == "DENY" and computed["billing"][1]["C9"] is False),
        ("redirect is denied by the parcel-redirection guardrail", decisions["redirect"] == "DENY" and computed["redirect"][1]["C10"] is False),
        ("wrong-person use is denied because requester must match the delegate", decisions["wrong-person"] == "DENY" and computed["wrong-person"][1]["C1"] is False),
        ("wrong-locker use is denied because locker must match the token", decisions["wrong-locker"] == "DENY" and computed["wrong-locker"][1]["C3"] is False),
        ("single-use reuse is denied after the token is already consumed", decisions["reuse"] == "DENY" and computed["reuse"][1]["C6"] is False),
        ("guardrail denials recompute to five out of five", guardrail_denials == 5 and "guardrail denials : 5/5" in ctx.answer),
        ("the release text matches parcel owner, delegate, parcel, locker, and site", "Noah may collect parcel123 for Maya from locker B17 at Station West" in ctx.answer),
    ]
    return run_checks(checks)
