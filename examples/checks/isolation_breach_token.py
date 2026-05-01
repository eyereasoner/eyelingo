# Independent Python checks for the isolation_breach_token example.
from __future__ import annotations

import re

from .common import run_checks


def run(ctx):
    data = ctx.load_input()
    media = data["media"]
    names = [m["name"] for m in media]
    states = [(m["zero"], m["one"]) for m in media]
    prepare_tasks = [(m["name"], state) for m in media for state in (m["zero"], m["one"])]
    copy_tasks = [(a, b) for a in names for b in names if a != b]
    serial = data["expected"]["serial"].split("->")
    reported_prepare = re.search(r"possible prepare tasks : (\d+)", ctx.answer)
    reported_prepare_count = int(reported_prepare.group(1)) if reported_prepare else None
    prepared_state = data["expected"]["preparedState"]

    checks = [
        ("all classical lab media carry the BreachBit variable", len(media) == 4 and data["variable"] == "BreachBit"),
        ("each medium has distinguishable safe/breach states", all(zero != one for zero, one in states)),
        ("prepare-task count is independently recomputed from all states", reported_prepare_count == len(prepare_tasks) == 8),
        ("the expected prepared breach state belongs to nursePager", any(m["name"] == "nursePager" and m["one"] == prepared_state for m in media) and f"nursePager prepares {prepared_state}" in ctx.answer),
        ("the expected serial audit path is backed by legal directed edges", all(item in names for item in serial) and all(edge in copy_tasks for edge in zip(serial, serial[1:]))),
        ("containmentPLC has at least two legal fan-out targets", sum(1 for a, _b in copy_tasks if a == "containmentPLC") >= 2),
        ("the specimen seal has a separate non-classical provenance variable", data["superinformationMedium"]["variable"] != data["variable"] and len(data["superinformationMedium"]["states"]) >= 3),
        ("the answer blocks universal cloning and unrestricted parallel fan-out for the specimen seal", "universal cloning" in ctx.answer and "unrestricted parallel fan-out" in ctx.answer),
        ("all three expected witnesses are reported in the answer", "classical breach token : YES" in ctx.answer and "specimen provenance seal : NO" in ctx.answer and "serial witness" in ctx.answer),
    ]
    return run_checks(checks)
