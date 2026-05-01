# Independent Python checks for the barley_seed_lineage example.
from __future__ import annotations

import re

from .common import run_checks


def blockers(world: dict, lineage: dict) -> list[str]:
    out = []
    if not world["noDesignLaws"]:
        out.append("design-laws")
    if not lineage["digitalHeredity"]:
        out.append("digital-heredity")
    if not lineage["repair"]:
        out.append("repair")
    if not lineage["dormancyProtection"]:
        out.append("protected-dormancy")
    if not set(["warmth", "moisture", "light"]).issubset(set(world["greenhouse"])):
        out.append("germination-resources")
    if not lineage["heritableVariation"]:
        out.append("heritable-variation")
    if lineage["variant"] != world["selectionFavours"]:
        out.append("adaptive-variant")
    return out


def run(ctx):
    data = ctx.load_input()
    computed_blockers = {line["name"]: blockers(data["world"], line) for line in data["lineages"]}
    evolvable = [name for name, blocks in computed_blockers.items() if not blocks]
    blocked = [name for name, blocks in computed_blockers.items() if blocks]
    answer_blocked = re.search(r"blocked contrast lineages : ([^\n]+)", ctx.answer)
    reported_blocked = [x.strip() for x in answer_blocked.group(1).split(",")] if answer_blocked else []

    checks = [
        ("no-design laws and greenhouse resources are available", data["world"]["noDesignLaws"] and set(data["world"]["greenhouse"]) == {"warmth", "moisture", "light"}),
        ("mainLine satisfies every CAN condition", computed_blockers["mainLine"] == []),
        ("mainLine is the unique evolvable lineage", evolvable == [data["expected"]["evolvable"]] and "evolvable lineage : mainLine" in ctx.answer),
        ("analogLine is blocked by missing digital heredity", computed_blockers["analogLine"] == ["digital-heredity"]),
        ("fragileLine is blocked by missing repair", computed_blockers["fragileLine"] == ["repair"]),
        ("coatlessLine is blocked by missing dormancy protection", computed_blockers["coatlessLine"] == ["protected-dormancy"]),
        ("staticLine is blocked by missing heritable variation", computed_blockers["staticLine"] == ["heritable-variation"]),
        ("reported blocked lineages match independent blocker analysis", reported_blocked == data["expected"]["blocked"] == blocked),
        ("adaptive persistence follows from the selected salt-tolerant variant", next(l for l in data["lineages"] if l["name"] == "mainLine")["variant"] == data["world"]["selectionFavours"]),
    ]
    return run_checks(checks)
