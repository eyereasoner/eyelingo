# Independent Python checks for the digital_product_passport example.
from __future__ import annotations

import re
from datetime import date

from .common import run_checks


def answer_field(answer: str, label: str) -> str | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*(.+)", answer)
    return match.group(1).strip() if match else None


def parse_int_field(answer: str, label: str, suffix: str = "") -> int | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*(\d+){re.escape(suffix)}", answer)
    return int(match.group(1)) if match else None


def document_types(documents: list[dict], section: str) -> set[str]:
    return {doc["DocType"] for doc in documents if doc["Section"] == section}


def declares(documents: list[dict], section: str, claim: str) -> bool:
    return any(doc["Section"] == section and claim in doc.get("Declares", []) for doc in documents)


def lifecycle_order_ok(events: list[dict]) -> bool:
    order = {"ManufacturingEvent": 0, "SaleEvent": 1, "RepairEvent": 2}
    last_day = None
    last_order = -1
    for event in events:
        day = date.fromisoformat(event["OnDate"])
        if last_day is not None and day < last_day:
            return False
        event_order = order.get(event["Type"], last_order)
        if event_order < last_order:
            return False
        last_day = day
        last_order = event_order
    return True


def run(ctx):
    data = ctx.load_input()
    total_mass = sum(int(component["MassG"]) for component in data["Components"])
    recycled_mass = sum(int(component["RecycledMassG"]) for component in data["Components"])
    recycled_pct = recycled_mass * 100 // total_mass if total_mass else 0
    footprint = data["Footprint"]
    lifecycle = int(footprint["ManufacturingGCO2e"]) + int(footprint["TransportGCO2e"]) + int(footprint["UsePhaseGCO2e"])
    material_is_critical = {material["ID"]: bool(material["CriticalRawMaterial"]) for material in data["Materials"]}
    critical_materials = sorted({
        material
        for component in data["Components"]
        for material in component["ContainsMaterial"]
        if material_is_critical.get(material, False)
    })
    public_section = data["AccessPolicy"]["PublicSection"]
    restricted_section = data["AccessPolicy"]["RestrictedSection"]
    public_docs = document_types(data["Documents"], public_section)
    restricted_docs = document_types(data["Documents"], restricted_section)
    battery_replaceable = any(component["Type"].lower() == "battery" and component["Replaceable"] for component in data["Components"])
    repair_friendly = (
        battery_replaceable
        and "RepairGuide" in public_docs
        and "SparePartsCatalog" in public_docs
        and declares(data["Documents"], public_section, data["AccessPolicy"]["RequiredPublicClaim"])
    )
    restricted_doc_types = set(data["AccessPolicy"]["RestrictedDocTypes"])
    public_doc_types = set(data["AccessPolicy"]["PublicDocTypes"])
    reported_critical = answer_field(ctx.answer, "critical raw materials")
    reported_decision = answer_field(ctx.answer, "Passport decision")

    checks = [
        ("component masses are folded from the JSON component list", parse_int_field(ctx.answer, "total component mass", " g") == total_mass == data["Expected"]["TotalMassG"]),
        ("recycled mass and integer recycled-content percentage are recomputed independently", recycled_mass == data["Expected"]["RecycledMassG"] and parse_int_field(ctx.answer, "recycled content", "%") == recycled_pct == data["Expected"]["RecycledContentPct"]),
        ("manufacturing, transport, and use-phase footprint values sum to the reported lifecycle footprint", parse_int_field(ctx.answer, "lifecycle footprint", " gCO2e") == lifecycle == data["Expected"]["LifecycleGCO2e"]),
        ("critical raw-material exposure is derived by joining component materials to material declarations", reported_critical == ", ".join(critical_materials) and critical_materials == ["Cobalt", "Lithium"]),
        ("repairFriendly is derived from a replaceable battery plus public repair and spare-parts documents", repair_friendly and answer_field(ctx.answer, "circularity hint") == data["Expected"]["CircularityHint"]),
        ("all required public document types are present only in the public document section", public_doc_types <= public_docs and not (public_doc_types & restricted_docs)),
        ("all restricted declaration document types stay in the restricted section", restricted_doc_types <= restricted_docs and not (restricted_doc_types & public_docs)),
        ("lifecycle events are chronological and follow manufacturing before sale before repair", lifecycle_order_ok(data["Lifecycle"])),
        ("the passport endpoint equals the product digital link and is reported publicly", data["Passport"]["PublicEndpoint"] == data["Product"]["DigitalLink"] == answer_field(ctx.answer, "public endpoint")),
        ("PASS is reported only because every independent passport check succeeds", reported_decision == f"PASS for {data['Product']['Model']} {data['Product']['SerialNumber']}."),
    ]
    return run_checks(checks)
