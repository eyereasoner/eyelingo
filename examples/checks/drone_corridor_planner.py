# Independent Python checks for the drone_corridor_planner example.
from __future__ import annotations

import math
import re

from .common import run_checks


FIELDS = ("Location", "Battery", "Permit")


def state_key(state: dict) -> tuple[str, str, str]:
    return tuple(state[f] for f in FIELDS)


def matches(pattern: dict, actual: dict) -> bool:
    return all(pattern[f] == "*" or pattern[f] == actual[f] for f in FIELDS)


def apply_state(pattern: dict, current: dict) -> dict:
    return {f: current[f] if pattern[f] == "*" else pattern[f] for f in FIELDS}


def search(data: dict) -> list[dict]:
    plans = []
    def walk(state, fuel, path, duration, cost, belief, comfort, seen):
        if state["Location"] == data["GoalLocation"]:
            if belief > data["Thresholds"]["MinBelief"] and cost < data["Thresholds"]["MaxCost"]:
                plans.append({"actions": path[:], "end": state.copy(), "duration": duration, "cost": cost, "belief": belief, "comfort": comfort, "fuel_left": fuel})
            return
        if fuel == 0:
            return
        for action in data["Actions"]:
            if not matches(action["From"], state):
                continue
            nxt = apply_state(action["To"], state)
            key = state_key(nxt)
            if key in seen:
                continue
            walk(nxt, fuel - 1, path + [action["Name"]], duration + action["DurationSec"], cost + action["Cost"], belief * action["Belief"], comfort * action["Comfort"], seen | {key})
    start = data["Start"]
    walk(start, data["Fuel"], [], 0, 0.0, 1.0, 1.0, {state_key(start)})
    plans.sort(key=lambda p: (round(p["cost"], 12), p["duration"], -p["belief"], ",".join(p["actions"])))
    return plans


def parse_answer(answer: str) -> dict:
    def grab(pattern, cast=str):
        m = re.search(pattern, answer)
        return cast(m.group(1)) if m else None
    plan = grab(r"selected plan : ([^\n]+)")
    return {
        "actions": [p.strip() for p in plan.split("->")] if plan else [],
        "duration": grab(r"duration : (\d+) s", int),
        "cost": grab(r"cost : ([0-9.]+)", float),
        "belief": grab(r"belief : ([0-9.]+)", float),
        "comfort": grab(r"comfort : ([0-9.]+)", float),
        "survivors": grab(r"surviving plans : (\d+)", int),
    }


def close(a, b, eps=5e-7):
    return a is not None and abs(a - b) <= eps


def run(ctx):
    data = ctx.load_input()
    plans = search(data)
    best = plans[0]
    reported = parse_answer(ctx.answer)
    next_cheapest = plans[1]
    checks = [
        ("14 corridor actions are loaded from JSON", len(data["Actions"]) == 14),
        ("bounded search independently finds 17 surviving plans", len(plans) == data["Expected"]["SurvivingPlans"] == reported["survivors"]),
        ("the selected path is the lowest-cost survivor", reported["actions"] == best["actions"] == ["fly_gent_brugge", "public_coastline_brugge_oostende"]),
        ("the selected path starts with the expected first action", best["actions"][0] == data["Expected"]["SelectedFirstAction"]),
        ("duration, cost, belief, and comfort are recomputed along the selected path", reported["duration"] == best["duration"] and close(reported["cost"], best["cost"], 5e-4) and close(reported["belief"], best["belief"]) and close(reported["comfort"], best["comfort"])),
        ("the selected end state reaches Oostende with low battery and no permit", best["end"] == {"Location": "Oostende", "Battery": "low", "Permit": "none"}),
        ("the selected belief and cost satisfy the thresholds", best["belief"] > data["Thresholds"]["MinBelief"] and best["cost"] < data["Thresholds"]["MaxCost"]),
        ("state-cycle pruning keeps every selected-path state unique", len(best["actions"]) + 1 == len({state_key(data["Start"]), state_key(best["end"]), ("Brugge", "mid", "none")})),
        ("the next cheapest survivor costs 0.014 as stated", math.isclose(next_cheapest["cost"], 0.014, rel_tol=0, abs_tol=1e-12) and "next cheapest" in ctx.reason),
    ]
    return run_checks(checks)
