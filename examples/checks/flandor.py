# Independent Python checks for the flandor example.
from __future__ import annotations

import re
from datetime import datetime

from .common import run_checks


def dt(value: str) -> datetime:
    return datetime.fromisoformat(value)


def active_needs(data: dict) -> dict[str, bool]:
    return {
        "export": any(c["ExportOrdersIndex"] < 90 for c in data["Clusters"]),
        "skills": data["Labour"]["TechVacancyRate"] > 3.9,
        "grid": data["Grid"]["CongestionHours"] > 11,
    }


def covers(pkg: dict, needs: dict[str, bool]) -> bool:
    return (
        (not needs["export"] or pkg["CoversExportWeakness"])
        and (not needs["skills"] or pkg["CoversSkillsStrain"])
        and (not needs["grid"] or pkg["CoversGridStress"])
    )


def choose_package(data: dict, needs: dict[str, bool]):
    candidates = [p for p in data["Packages"] if covers(p, needs) and p["CostMEUR"] <= data["Budget"]["MaxMEUR"]]
    return min(candidates, key=lambda p: (p["CostMEUR"], p["PackageID"])) if candidates else None


def parse_answer(answer: str):
    vals = {}
    for key, pattern, cast in [
        ("active", r"Active need count: (\d+)/", int),
        ("pkg", r"Recommended package: ([^\n]+)", str),
        ("cost", r"Package cost: €(\d+)M", int),
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
    selected = choose_package(data, needs)
    reported = parse_answer(ctx.answer)
    insight = data["Insight"]
    policy = data["Policy"]
    signature = data["Signature"]
    allowed = (
        data["RequestAction"] == policy["PermissionAction"]
        and data["RequestPurpose"] == policy["PermissionPurpose"]
        and policy["PermissionTarget"] == insight["ID"]
        and dt(data["BoardAuthAt"]) < dt(insight["ExpiresAt"])
        and active_count >= insight["ThresholdScore"]
        and selected is not None
    )
    export_weak = [c["Name"] for c in data["Clusters"] if c["ExportOrdersIndex"] < 90]
    serialized = insight["SerializedLowercase"].lower()

    checks = [
        ("export weakness, skills strain, and grid stress are all active", needs == {"export": True, "skills": True, "grid": True} and active_count == 3),
        ("active need count meets the insight threshold", reported["active"] == active_count == insight["ThresholdScore"]),
        ("the lowest-cost package covering all active needs is selected", selected["Name"] == "Flandor Retooling Pulse" and reported["pkg"] == selected["Name"]),
        ("the selected package fits inside the €140M budget", selected["CostMEUR"] == reported["cost"] == 120 and selected["CostMEUR"] <= data["Budget"]["MaxMEUR"]),
        ("cheaper packages are rejected because each covers only one active need", all(not covers(p, needs) for p in data["Packages"] if p["CostMEUR"] < selected["CostMEUR"])),
        ("the full corridor package covers all needs but is over budget", covers(data["Packages"][-1], needs) and data["Packages"][-1]["CostMEUR"] > data["Budget"]["MaxMEUR"]),
        ("policy permission authorizes regional-stabilization use before expiry", allowed and "decision : ALLOWED" in ctx.answer),
        ("firm-surveillance redistribution is prohibited", policy["ProhibitionAction"] == "odrl:distribute" and policy["ProhibitionPurpose"] == "firm_surveillance"),
        ("deletion duty is scheduled before envelope expiry", dt(data["BoardDutyAt"]) < dt(insight["ExpiresAt"])),
        ("shared insight omits firm names and payroll rows", not data["Signals"]["ContainsFirmNames"] and not data["Signals"]["ContainsPayrollRows"] and "firm" not in serialized and "payroll" not in serialized),
        ("reported signature metadata matches the trusted precomputed input", reported["hash"] == signature["DisplayPayloadSHA256"] and reported["hmac"] == signature["SignatureHMAC"] and signature["Algorithm"] == "HMAC-SHA256"),
        ("the expected six files and one audit entry are recorded", data["FilesWritten"] == data["ExpectedFilesWritten"] == 6 and data["AuditEntries"] == 1),
        ("export-weak cluster names are independently identified", export_weak == ["Antwerp chemicals", "Ghent manufacturing"]),
    ]
    return run_checks(checks)
