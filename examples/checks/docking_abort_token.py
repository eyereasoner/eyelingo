# Independent Python checks for the docking_abort_token example.
from __future__ import annotations

import re

from .common import run_checks


def run(ctx):
    data = ctx.load_input()
    media = data["media"]
    names = [m["name"] for m in media]
    copy_tasks = [(a, b) for a in names for b in names if a != b]
    serial = data["expected"]["serial"].split("->")
    parallel_source = data["expected"]["parallelSource"]
    reported_copy = re.search(r"possible copy tasks : (\d+)", ctx.answer)
    reported_copy_count = int(reported_copy.group(1)) if reported_copy else None

    checks = [
        ("all four media carry the same AbortBit variable", len(media) == 4 and data["variable"] == "AbortBit"),
        ("each medium has distinct zero and one states", all(m["zero"] != m["one"] for m in media)),
        ("the directed copy-task count is recomputed as n*(n-1)", reported_copy_count == len(copy_tasks) == 12),
        ("the expected serial witness uses known media in order", all(item in names for item in serial) and f"{serial[0]} -> {serial[1]} -> {serial[2]}" in ctx.answer),
        ("the serial witness is backed by two legal copy/measure edges", (serial[0], serial[1]) in copy_tasks and (serial[1], serial[2]) in copy_tasks),
        ("the expected parallel source can fan out to two other media", parallel_source in names and sum(1 for a, _b in copy_tasks if a == parallel_source) >= 2 and "parallel witness : flightPLC -> radioFrame and auditDisplay" in ctx.answer),
        ("the quantum seal is separate from the classical AbortBit variable", data["superinformationMedium"]["variable"] != data["variable"]),
        ("the answer blocks universal cloning and unrestricted audit fan-out for the seal", "cannot be universally cloned" in ctx.answer and "unrestricted audit fan-out" in ctx.answer),
    ]
    return run_checks(checks)
