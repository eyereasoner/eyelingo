# Independent Python checks for the alarm_bit_interoperability example.
from __future__ import annotations

import re

from .common import run_checks


def run(ctx):
    data = ctx.load_input()
    media = data["ClassicalMedia"]
    variable = media[0]["Variable"] if media else None
    directed = [(a["Name"], b["Name"]) for a in media for b in media if a["Name"] != b["Name"] and a["Variable"] == b["Variable"]]
    reported_tasks = re.findall(r"copy task : ([A-Za-z0-9_]+) -> ([A-Za-z0-9_]+) for ([A-Za-z0-9_]+)", ctx.answer)
    blocked_match = re.search(r"blocked tasks : ([^\n]+)", ctx.answer)
    blocked = {x.strip() for x in blocked_match.group(1).split(",")} if blocked_match else set()
    answer_can = "classical alarm-bit interoperability : YES" in ctx.answer
    answer_cannot = "universal cloning of the superinformation token : NO" in ctx.answer

    checks = [
        ("all classical media encode the same abstract variable", len(media) >= 2 and all(m["Variable"] == variable for m in media)),
        ("the directed copy-task count is recomputed from the media graph", len(directed) == data["ExpectedCopyTasks"]),
        ("the report lists exactly the expected directed copy tasks", {(a, b) for a, b, _ in reported_tasks} == set(directed)),
        ("each classical medium has distinguishable zero and one states", all(m["ZeroState"] != m["OneState"] for m in media)),
        ("the superinformation contrast has more than two named states", len(data["Superinformation"]["States"]) > 2),
        ("all expected impossible tasks are reported as blocked", blocked == set(data["ExpectedImpossible"])),
        ("the reported CAN/CANNOT decisions match the expected polarities", answer_can == (data["ExpectedCanDecision"] == "YES") and answer_cannot == (data["ExpectedCantDecision"] == "NO")),
    ]
    return run_checks(checks)
