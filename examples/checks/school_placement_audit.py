# Independent Python checks for the school_placement_audit example.
from __future__ import annotations

import re

from .common import run_checks


def preference_rank(student: dict, school: str) -> int:
    try:
        return list(student["preferences"]).index(school)
    except ValueError:
        return 999


def by_student_and_school(data: dict) -> dict[tuple[str, str], dict]:
    return {(row["student"], row["school"]): row for row in data["distances"]}


def choose_school(data: dict, student: dict, audited: bool) -> dict:
    best = None
    penalty = int(data["policy"]["preferencePenaltyMeters"])
    for row in data["distances"]:
        if row["student"] != student["id"]:
            continue
        rank = preference_rank(student, row["school"])
        score = int(row["straightMeters"])
        if audited:
            score = int(row["walkingMeters"]) + rank * penalty
        candidate = {"student": student, "school": row["school"], "distance": row, "score": score, "rank": rank}
        key = (candidate["score"], candidate["rank"], candidate["school"])
        if best is None or key < (best["score"], best["rank"], best["school"]):
            best = candidate
    assert best is not None
    return best


def hidden_detour(assignment: dict) -> int:
    row = assignment["distance"]
    return int(row["walkingMeters"]) - int(row["straightMeters"])


def derive(data: dict):
    straight = {}
    audited = {}
    changes = []
    max_walk = int(data["policy"]["maxWalkingMeters"])
    for student in data["students"]:
        flawed = choose_school(data, student, audited=False)
        better = choose_school(data, student, audited=True)
        straight[student["id"]] = flawed
        audited[student["id"]] = better
        if flawed["school"] != better["school"] or int(flawed["distance"]["walkingMeters"]) > max_walk:
            changes.append({"student": student, "straight": flawed, "audited": better})
    largest = max(changes, key=lambda change: hidden_detour(change["straight"]))
    return straight, audited, changes, largest


def parse_answer(answer: str):
    result_match = re.search(r"audit result\s*:\s*(\w+)", answer)
    affected_match = re.search(r"children affected by straight-line rule\s*:\s*(.+)", answer)
    detour_match = re.search(r"largest hidden detour\s*:\s*([^,]+),\s*(\d+)\s*m", answer)
    assignment_match = re.search(r"recommended assignments\s*:\s*(.+)", answer)
    explanation_match = re.search(r"explanation requested\s*:\s*(\w+)", answer)

    affected = []
    if affected_match:
        affected = [part.strip() for part in affected_match.group(1).split(",") if part.strip()]

    assignments = {}
    if assignment_match:
        for part in assignment_match.group(1).split(";"):
            if "->" not in part:
                continue
            student, school = [piece.strip() for piece in part.split("->", 1)]
            assignments[student] = school

    return {
        "result": result_match.group(1) if result_match else None,
        "affected": affected,
        "detour_student": detour_match.group(1).strip() if detour_match else None,
        "detour_meters": int(detour_match.group(2)) if detour_match else None,
        "assignments": assignments,
        "explanation": explanation_match.group(1) if explanation_match else None,
    }


def run(ctx):
    data = ctx.load_input()
    parsed = parse_answer(ctx.answer)
    straight, audited, changes, largest = derive(data)
    schools = {school["name"] for school in data["schools"]}
    distance_pairs = by_student_and_school(data)

    expected_audited_by_name = {change["student"]["name"]: change["audited"]["school"] for change in changes}
    expected_audited_by_name.update({student["name"]: audited[student["id"]]["school"] for student in data["students"]})
    affected_names = sorted(change["student"]["name"] for change in changes)
    provisional_names = {student["name"]: straight[student["id"]]["school"] for student in data["students"]}
    audited_names = {student["name"]: audited[student["id"]]["school"] for student in data["students"]}
    over_limit = [
        change for change in changes
        if int(change["straight"]["distance"]["walkingMeters"]) > int(data["policy"]["maxWalkingMeters"])
    ]

    reason = ctx.reason.lower()

    checks = [
        (
            "the fixture has four students, four schools, and a complete 4 × 4 distance matrix",
            len(data["students"]) == 4 and len(data["schools"]) == 4 and len(distance_pairs) == 16,
        ),
        (
            "every student preference list covers every school exactly once",
            all(set(student["preferences"]) == schools and len(student["preferences"]) == len(schools) for student in data["students"]),
        ),
        (
            "the independent straight-line rule assigns Ada and Björn to Centrum",
            provisional_names["Ada"] == "Centrum" and provisional_names["Björn"] == "Centrum",
        ),
        (
            "the independent route-aware rule computes walking distance plus preference penalty",
            audited_names == {"Ada": "Lindholmen", "Björn": "Backa", "Clara": "Haga", "Davi": "Haga"},
        ),
        (
            "the reported recommended assignments match the Python audit",
            parsed["assignments"] == expected_audited_by_name,
        ),
        (
            "the reported affected children are exactly those whose provisional placement is flagged",
            parsed["affected"] == affected_names and affected_names == ["Ada", "Björn", "Davi"],
        ),
        (
            "the reported largest hidden detour matches the flawed straight-line placement",
            parsed["detour_student"] == largest["student"]["name"] and parsed["detour_meters"] == hidden_detour(largest["straight"]) == 3000,
        ),
        (
            "the failure result follows from at least one over-limit walking route and changed assignment",
            parsed["result"] == "fail" and len(over_limit) >= 2 and any(change["straight"]["school"] != change["audited"]["school"] for change in changes),
        ),
        (
            "the Reason text names the support-tool rule, walking-route recomputation, and inspectability requirement",
            "support-tool" in reason and "walking-route" in reason and "inspectable" in reason,
        ),
        (
            "the report explicitly requests an explanation for the affected placement decisions",
            parsed["explanation"] == "yes",
        ),
    ]
    return run_checks(checks)
