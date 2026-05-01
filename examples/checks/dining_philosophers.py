# Independent Python checks for the dining_philosophers example.
from __future__ import annotations

import re
from collections import Counter

from .common import run_checks


PHILOSOPHERS = ["P1", "P2", "P3", "P4", "P5"]
FORKS = ["F12", "F23", "F34", "F45", "F51"]
LEFT_FORK = {"P1": "F51", "P2": "F12", "P3": "F23", "P4": "F34", "P5": "F45"}
RIGHT_FORK = {"P1": "F12", "P2": "F23", "P3": "F34", "P4": "F45", "P5": "F51"}
INITIAL_STATE = {
    "F12": ["P1", "Dirty"],
    "F23": ["P2", "Dirty"],
    "F34": ["P3", "Dirty"],
    "F45": ["P4", "Dirty"],
    "F51": ["P1", "Dirty"],
}


def holder_for(fork: str, state: dict[str, list[str]]) -> str:
    return state[fork][0]


def derive_trace(schedule: list[dict]) -> tuple[list[dict], Counter, dict[str, list[str]]]:
    state = {k: v[:] for k, v in INITIAL_STATE.items()}
    meal_counts: Counter[str] = Counter()
    traces = []
    for round_ in schedule:
        hungry = round_["Hungry"]
        requests = []
        for p in sorted(hungry, key=PHILOSOPHERS.index):
            for side, fork in (("left", LEFT_FORK[p]), ("right", RIGHT_FORK[p])):
                holder, clean = state[fork]
                if holder != p:
                    requests.append((p, holder, fork, side, clean == "Dirty"))
        transfers = []
        for p, holder, fork, _side, dirty in requests:
            if dirty:
                transfers.append((holder, p, fork))
        for holder, p, fork in transfers:
            state[fork] = [p, "Clean"]
        meals = []
        forks_used = []
        for p in sorted(hungry, key=PHILOSOPHERS.index):
            lf, rf = LEFT_FORK[p], RIGHT_FORK[p]
            if holder_for(lf, state) == p and holder_for(rf, state) == p:
                meal_counts[p] += 1
                meals.append((p, meal_counts[p], lf, rf))
                forks_used.extend([lf, rf])
        for fork in FORKS:
            state[fork][1] = "Dirty"
        traces.append({"round": round_, "requests": requests, "transfers": transfers, "meals": meals, "forks_used": forks_used, "end": {k: v[:] for k, v in state.items()}})
    return traces, meal_counts, state


def parse_meal_trace(answer: str) -> list[tuple[int, list[str]]]:
    rows = []
    for round_no, meal_text in re.findall(r"round (\d+) cycle \d+ : ([^\n]+)", answer):
        philosophers = re.findall(r"(P\d)#\d+ uses", meal_text)
        rows.append((int(round_no), philosophers))
    return rows


def run(ctx):
    schedule = ctx.load_input()
    traces, meal_counts, final_state = derive_trace(schedule)
    reported_rows = parse_meal_trace(ctx.answer)
    meal_total = sum(len(t["meals"]) for t in traces)
    request_total = sum(len(t["requests"]) for t in traces)
    transfer_total = sum(len(t["transfers"]) for t in traces)
    no_shared_forks = all(len(t["forks_used"]) == len(set(t["forks_used"])) for t in traces)
    final_holders = {fork: holder for fork, (holder, _clean) in final_state.items()}

    checks = [
        ("the JSON schedule contains nine deterministic rounds", len(schedule) == 9 and [r["Number"] for r in schedule] == list(range(1, 10))),
        ("the Chandy-Misra trace yields exactly 15 meals", meal_total == 15 and "meals : 15" in ctx.answer),
        ("each philosopher eats exactly three times", all(meal_counts[p] == 3 for p in PHILOSOPHERS) and "everyone ate 3 times : yes" in ctx.answer),
        ("no two meals in the same round use the same fork", no_shared_forks),
        ("dirty-fork transfer simulation reproduces the reported meal pattern", [ps for _r, ps in reported_rows] == [[m[0] for m in t["meals"]] for t in traces]),
        ("the first round transfers only F23 to P3 and lets P1/P3 eat", traces[0]["transfers"] == [("P2", "P3", "F23")] and [m[0] for m in traces[0]["meals"]] == ["P1", "P3"]),
        ("the P5-only rounds derive one meal each", all([m[0] for m in traces[i]["meals"]] == ["P5"] for i in [2, 5, 8])),
        ("all forks end dirty after the final phase", all(clean == "Dirty" for _holder, clean in final_state.values())),
        ("the final fork holders match the independent simulation", final_holders == {"F12": "P2", "F23": "P2", "F34": "P4", "F45": "P5", "F51": "P5"}),
        ("request and transfer counts are nontrivial and internally consistent", request_total == 26 and transfer_total == 26),
    ]
    return run_checks(checks)
