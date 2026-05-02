# Independent Python checks for the eco_route_insight example.
from __future__ import annotations

import base64
import hashlib
import hmac
import json
import re
from datetime import datetime, timedelta, timezone

from .common import run_checks


def parse_bool(value: str | None) -> bool | None:
    if value is None:
        return None
    value = value.strip().lower()
    if value == "yes":
        return True
    if value == "no":
        return False
    return None


def parse_answer(answer: str) -> dict:
    fields = {}
    for line in answer.splitlines():
        if " : " in line:
            key, value = line.split(" : ", 1)
            fields[key.strip()] = value.strip()
    return fields


def fuel_index(route: dict, payload_tonnes: float) -> float:
    return float(route["distanceKm"]) * payload_tonnes * float(route["gradientFactor"])


def expiry(issued_at: str, ttl_hours: int) -> str:
    parsed = datetime.fromisoformat(issued_at.replace("Z", "+00:00"))
    return (parsed + timedelta(hours=int(ttl_hours))).astimezone(timezone.utc).isoformat().replace("+00:00", "Z")


def choose_alternative(data: dict, current_fuel: float) -> dict:
    current_eta = int(data["currentRoute"]["etaMinutes"])
    max_delay = int(data["policy"]["maxEtaDelayMinutes"])
    threshold = float(data["policy"]["fuelIndexThreshold"])
    payload_tonnes = float(data["shipment"]["payloadKg"]) / 1000.0
    scored = []
    for route in data["alternativeRoutes"]:
        fi = fuel_index(route, payload_tonnes)
        saving = current_fuel - fi
        delay = int(route["etaMinutes"]) - current_eta
        eligible = saving > 0 and delay <= max_delay
        if data["policy"].get("requireAlternativeBelowThreshold"):
            eligible = eligible and fi <= threshold
        scored.append({"route": route, "fuelIndex": fi, "saving": saving, "delay": delay, "eligible": eligible})
    scored.sort(key=lambda row: (not row["eligible"], -row["saving"], row["route"]["id"]))
    return scored[0]


def go_number(value: float):
    rounded = round(value, 2)
    if float(rounded).is_integer():
        return int(rounded)
    return rounded


def canonical_envelope(data: dict, current_fuel: float, selected: dict) -> dict:
    issued_at = data["scenario"]["issuedAt"]
    exp = expiry(issued_at, data["scenario"]["ttlHours"])
    issue = current_fuel > float(data["policy"]["fuelIndexThreshold"]) and selected["eligible"]
    return {
        "audience": data["scenario"]["depot"],
        "allowedUse": data["policy"]["allowedUse"],
        "issuedAt": issued_at,
        "expiry": exp,
        "keyId": data["signing"]["keyId"],
        "assertions": {
            "showEcoBanner": issue,
            "suggestedRoute": selected["route"]["id"],
            "currentFuelIndex": go_number(current_fuel),
            "suggestedFuelIndex": go_number(selected["fuelIndex"]),
            "estimatedSaving": go_number(selected["saving"]),
            "rawDataExported": False,
        },
    }


def stable_json(value: dict) -> str:
    return json.dumps(value, separators=(",", ":"), ensure_ascii=False)


def signature(secret: str, canonical: str) -> str:
    digest = hmac.new(secret.encode("utf-8"), canonical.encode("utf-8"), hashlib.sha256).digest()
    return base64.urlsafe_b64encode(digest).decode("ascii").rstrip("=")


def extract_reason_numbers(reason: str) -> dict:
    current = re.search(r"gives\s+([0-9]+(?:\.[0-9]+)?)\s+×\s+([0-9]+(?:\.[0-9]+)?)\s+×\s+([0-9]+(?:\.[0-9]+)?)\s+=\s+([0-9]+(?:\.[0-9]+)?)", reason)
    selected = re.search(r"alternative\s+(\S+)\s+gives\s+([0-9]+(?:\.[0-9]+)?)\s+×\s+([0-9]+(?:\.[0-9]+)?)\s+×\s+([0-9]+(?:\.[0-9]+)?)\s+=\s+([0-9]+(?:\.[0-9]+)?),\s+saving\s+([0-9]+(?:\.[0-9]+)?)", reason)
    return {
        "current": tuple(float(x) for x in current.groups()) if current else None,
        "selected_route": selected.group(1) if selected else None,
        "selected": tuple(float(x) for x in selected.groups()[1:]) if selected else None,
    }


def run(ctx):
    data = ctx.load_input()
    answer = parse_answer(ctx.answer)
    payload_tonnes = float(data["shipment"]["payloadKg"]) / 1000.0
    current_fuel = fuel_index(data["currentRoute"], payload_tonnes)
    selected = choose_alternative(data, current_fuel)
    env = canonical_envelope(data, current_fuel, selected)
    canonical = stable_json(env)
    digest = hashlib.sha256(canonical.encode("utf-8")).hexdigest()
    sig = signature(data["signing"]["secret"], canonical)
    assertions = env["assertions"]
    reason_numbers = extract_reason_numbers(ctx.reason)
    forbidden = data["dataMinimization"]["forbiddenEnvelopeTerms"]

    checks = [
        (
            "the current fuel index is recomputed from distance, payload tonnes, and gradient",
            round(current_fuel, 2) == 120.75 and answer.get("current fuel index") == "120.75",
        ),
        (
            "the policy threshold triggers a local eco banner",
            current_fuel > float(data["policy"]["fuelIndexThreshold"]) and parse_bool(answer.get("show eco banner")) is True,
        ),
        (
            "the selected alternative is the best eligible lower-fuel route",
            selected["route"]["id"] == "alt-low-fuel" and selected["eligible"] and answer.get("suggested route") == "alt-low-fuel",
        ),
        (
            "the alternative fuel index and saving are independently recomputed",
            round(selected["fuelIndex"], 2) == 99.75 and round(selected["saving"], 2) == 21.00 and answer.get("suggested fuel index") == "99.75" and answer.get("estimated saving") == "21.00",
        ),
        (
            "the selected route stays within the allowed ETA delay",
            selected["delay"] == 2 and selected["delay"] <= int(data["policy"]["maxEtaDelayMinutes"]),
        ),
        (
            "the envelope audience, allowed use, expiry, and raw-data flag match the policy",
            answer.get("audience") == env["audience"] and answer.get("allowed use") == env["allowedUse"] and answer.get("expires at") == env["expiry"] and parse_bool(answer.get("raw data exported")) is False,
        ),
        (
            "the canonical envelope omits forbidden raw logistics terms",
            not any(term in canonical for term in forbidden),
        ),
        (
            "the payload digest is SHA-256 over the independently rebuilt envelope",
            answer.get("payload digest") == digest,
        ),
        (
            "the signature is the expected base64url HMAC-SHA256 value",
            answer.get("signature algorithm") == data["policy"]["signatureAlgorithm"] and answer.get("signature key") == data["signing"]["keyId"] and answer.get("signature") == sig,
        ),
        (
            "the Reason text reports the same arithmetic and the insight pattern",
            reason_numbers["current"] == (42.0, 2.5, 1.15, 120.75) and reason_numbers["selected_route"] == "alt-low-fuel" and reason_numbers["selected"] == (38.0, 2.5, 1.05, 99.75, 21.0) and "ship the decision, not the data" in ctx.reason,
        ),
    ]
    return run_checks(checks)
