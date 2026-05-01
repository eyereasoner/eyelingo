# Independent Python checks for the calidor example.
from __future__ import annotations

import re
from datetime import datetime

from .common import run_checks


def dt(value: str) -> datetime:
    return datetime.fromisoformat(value)


REQUIRED_CAPABILITIES = {"bill_credit", "cooling_kit", "transport", "welfare_check"}


def active_needs(data: dict) -> dict[str, bool]:
    return {
        "heat_alert": data["CurrentAlertLevel"] >= data["AlertLevelAtLeast"],
        "unsafe_indoor_heat": data["CurrentIndoorTempC"] >= data["IndoorTempCAtLeast"] and data["HoursAtOrAboveThreshold"] >= data["HoursAtOrAboveThresholdAtLeast"],
        "vulnerability_present": len(data["VulnerabilityFlags"]) > 0,
        "energy_constraint": data["RemainingPrepaidCreditEur"] <= data["EnergyCreditEurAtMost"],
    }


def choose_package(data: dict):
    candidates = []
    for pkg in data["Packages"]:
        caps = set(pkg["Capabilities"])
        covers_all = REQUIRED_CAPABILITIES.issubset(caps)
        within_budget = pkg["CostEur"] <= data["MaxPackageCostEur"]
        if covers_all and within_budget:
            candidates.append(pkg)
    return min(candidates, key=lambda p: (p["CostEur"], p["PackageID"])) if candidates else None


def parse_answer(answer: str):
    vals = {}
    for key, pattern, cast in [
        ("active", r"Active need count: (\d+)/", int),
        ("pkg", r"Recommended package: ([^\n]+)", str),
        ("cost", r"Package cost: €(\d+)", int),
        ("hash", r"Payload SHA-256: ([0-9a-f]+)", str),
        ("hmac", r"Envelope HMAC-SHA-256: ([0-9a-f]+)", str),
    ]:
        m = re.search(pattern, answer)
        vals[key] = cast(m.group(1)) if m else None
    return vals


def run(ctx):
    data = ctx.load_input()
    needs = active_needs(data)
    active_count = sum(needs.values())
    selected = choose_package(data)
    reported = parse_answer(ctx.answer)
    insight = data["Insight"]
    policy = data["Policy"]
    signature = data["Signature"]
    allowed = (
        data["RequestAction"] == policy["PermissionAction"]
        and data["RequestPurpose"] == policy["PermissionPurpose"]
        and policy["PermissionTarget"] == insight["ID"]
        and dt(data["CityAuthAt"]) < dt(data["GatewayExpiresAt"])
        and active_count >= data["MinimumActiveNeedCount"]
        and selected is not None
    )
    delete_before_expiry = dt(data["CityDutyAt"]) < dt(data["GatewayExpiresAt"])
    serialized = insight["SerializedLowercase"].lower()
    local_terms = ["heat_sensitive_condition", "mobility_limitation", "prepaid"]

    checks = [
        ("four active heat-response needs are recomputed from the local signals", active_count == 4 and reported["active"] == 4),
        ("the insight threshold of three active needs is met", active_count >= data["MinimumActiveNeedCount"] == insight["ThresholdCount"]),
        ("the lowest-cost eligible package covering all required capabilities is selected", selected["Name"] == "Calidor Priority Cooling Bundle" and reported["pkg"] == selected["Name"]),
        ("the selected package fits inside the €20 budget", selected["CostEur"] == reported["cost"] == 18 and selected["CostEur"] <= data["MaxPackageCostEur"]),
        ("cheaper packages are rejected because they do not cover all capabilities", all(not REQUIRED_CAPABILITIES.issubset(set(p["Capabilities"])) for p in data["Packages"] if p["CostEur"] < selected["CostEur"])),
        ("the deluxe package is rejected because it is over budget", next(p for p in data["Packages"] if p["PackageID"] == "pkg:DELUXE")["CostEur"] > data["MaxPackageCostEur"]),
        ("policy permission authorizes heatwave-response use before expiry", allowed and "decision : ALLOWED" in ctx.answer),
        ("tenant-screening reuse is prohibited by the policy", policy["ProhibitionAction"] == "odrl:distribute" and policy["ProhibitionPurpose"] == "tenant_screening"),
        ("deletion duty is scheduled before envelope expiry", delete_before_expiry),
        ("vulnerability flags and local raw stress signals are omitted from the serialized insight", not any(term in serialized for term in local_terms)),
        ("reported signature metadata matches the trusted precomputed input", reported["hash"] == signature["PayloadHashSHA256"] and reported["hmac"] == signature["SignatureHMAC"] and signature["Algorithm"] == "HMAC-SHA256"),
        ("scope metadata is explicit for device, event, municipality, creation, and expiry", all(insight.get(k) for k in ["ScopeDevice", "ScopeEvent", "Municipality", "CreatedAt", "ExpiresAt"])),
    ]
    return run_checks(checks)
