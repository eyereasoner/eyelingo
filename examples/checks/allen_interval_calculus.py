# Independent Python checks for the allen_interval_calculus example.
from __future__ import annotations

import re
from datetime import datetime, timedelta, timezone

from .common import run_checks


def parse_time(value: str) -> datetime:
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def complete_interval(item: dict) -> tuple[datetime, datetime]:
    start = parse_time(item["start"])
    if "end" in item:
        end = parse_time(item["end"])
    else:
        end = start + timedelta(minutes=int(item["durationMinutes"]))
    return start, end


def relation(a: tuple[datetime, datetime], b: tuple[datetime, datetime]) -> str:
    a0, a1 = a
    b0, b1 = b
    if a1 < b0:
        return "before"
    if a0 > b1:
        return "after"
    if a1 == b0:
        return "meets"
    if a0 == b1:
        return "metBy"
    if a0 < b0 and a1 > b0 and a1 < b1:
        return "overlaps"
    if a0 > b0 and a0 < b1 and a1 > b1:
        return "overlappedBy"
    if a0 == b0 and a1 < b1:
        return "starts"
    if a0 == b0 and a1 > b1:
        return "startedBy"
    if a0 > b0 and a1 < b1:
        return "during"
    if a0 < b0 and a1 > b1:
        return "contains"
    if a0 > b0 and a1 == b1:
        return "finishes"
    if a0 < b0 and a1 == b1:
        return "finishedBy"
    if a0 == b0 and a1 == b1:
        return "equals"
    return "unknown"


def run(ctx):
    data = ctx.load_input()
    intervals = {item["name"]: complete_interval(item) for item in data["intervals"]}
    ordered = {(a, b): relation(intervals[a], intervals[b]) for a in intervals for b in intervals if a != b}
    required = {tuple(k.split("|")): v for k, v in data["expected"]["requiredRelations"].items()}
    completed = {item["name"]: intervals[item["name"]] for item in data["intervals"] if "durationMinutes" in item}
    invalid = [name for name, (start, end) in intervals.items() if end <= start]
    count_match = re.search(r"derived relations : (\d+) ordered interval pairs", ctx.answer)
    reported_count = int(count_match.group(1)) if count_match else None

    checks = [
        ("duration-based intervals are completed from start plus minutes", completed["I"][1].hour == 18 and completed["K"][1].hour == 14 and completed["K"][1].minute == 0),
        ("all completed intervals have strictly positive duration", not invalid),
        ("all ordered non-self interval pairs are classified", reported_count == len(data["intervals"]) * (len(data["intervals"]) - 1) == len(ordered)),
        ("every required Allen relation is recomputed from endpoints", all(ordered[pair] == rel for pair, rel in required.items())),
        ("A/B, A/C, and A/D demonstrate before, meets, and overlaps", ordered[("A", "B")] == "before" and ordered[("A", "C")] == "meets" and ordered[("A", "D")] == "overlaps"),
        ("converse relations are independently recovered", ordered[("B", "A")] == "after" and ordered[("C", "A")] == "metBy" and ordered[("D", "A")] == "overlappedBy"),
        ("start, finish, during, contains, and equals cases all occur", {"starts", "finishes", "during", "contains", "equals"}.issubset(set(ordered.values()))),
        ("the showcase text includes each required forward example", all(f"{a} {rel} {b}" in ctx.answer for (a, b), rel in required.items() if (a, b) in [("A","B"),("A","C"),("A","D"),("F","A"),("G","A"),("A","H"),("A","E"),("J","I"),("K","C")])),
    ]
    return run_checks(checks)
