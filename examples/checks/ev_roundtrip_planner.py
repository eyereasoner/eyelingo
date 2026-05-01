# Independent Python checks for the ev_roundtrip_planner example.
from __future__ import annotations

import re

from .common import run_checks


FIELDS = ("At", "Battery", "Pass")


def key(state: dict) -> tuple[str, str, str]:
    return tuple(state[f] for f in FIELDS)


def matches(pattern: dict, state: dict) -> bool:
    return all(pattern[f] == "*" or pattern[f] == state[f] for f in FIELDS)


def apply(pattern: dict, state: dict) -> dict:
    return {f: state[f] if pattern[f] == "*" else pattern[f] for f in FIELDS}


def matches_goal(state: dict, goal: dict) -> bool:
    return matches(goal, state)


def search(data: dict):
    plans = []
    stats = {"max_depth": 0}
    def walk(state, path, duration, cost, belief, comfort, fuel_left, seen):
        stats["max_depth"] = max(stats["max_depth"], len(path))
        if matches_goal(state, data["Goal"]):
            th = data["Thresholds"]
            if belief > th["MinBelief"] and cost < th["MaxCost"] and duration < th["MaxDuration"]:
                plans.append({"actions": path[:], "state": state.copy(), "duration": duration, "cost": cost, "belief": belief, "comfort": comfort, "fuel": fuel_left})
            return
        if fuel_left == 0:
            return
        for action in data["Actions"]:
            if not matches(action["From"], state):
                continue
            nxt = apply(action["To"], state)
            nxt_key = key(nxt)
            if nxt_key in seen and nxt_key != key(state):
                continue
            walk(nxt, path + [action["Name"]], duration + action["Duration"], cost + action["Cost"], belief * action["Belief"], comfort * action["Comfort"], fuel_left - 1, seen | {nxt_key})
    start = {f: data["Vehicle"][f] for f in FIELDS}
    walk(start, [], 0.0, 0.0, 1.0, 1.0, data["FuelSteps"], {key(start)})
    plans.sort(key=lambda p: (p["duration"], p["cost"], "/".join(p["actions"])))
    return plans, stats


def parse_answer(answer: str):
    def grab(pattern, cast=str):
        m = re.search(pattern, answer)
        return cast(m.group(1)) if m else None
    plan = grab(r"Select plan : ([^.]+)\.")
    return {
        "actions": [p.strip() for p in plan.split("->")] if plan else [],
        "duration": grab(r"duration : ([0-9.]+) minutes", float),
        "cost": grab(r"cost : ([0-9.]+)", float),
        "belief": grab(r"belief : ([0-9.]+)", float),
        "comfort": grab(r"comfort : ([0-9.]+)", float),
        "acceptable": grab(r"acceptable plans : (\d+)", int),
        "fuel": grab(r"fuel remaining : (\d+) of", int),
    }


def close(a, b, eps=5e-7):
    return a is not None and abs(a - b) <= eps


def run(ctx):
    data = ctx.load_input()
    plans, stats = search(data)
    best = plans[0]
    reported = parse_answer(ctx.answer)
    checks = [
        ("bounded search finds eight acceptable Brussels-to-Cologne plans", len(plans) == 8 == reported["acceptable"]),
        ("the selected plan is the fastest acceptable candidate", reported["actions"] == best["actions"] == ["drive_bru_liege", "drive_liege_aachen", "shuttle_aachen_cologne"]),
        ("selected duration and fuel remaining are recomputed", close(reported["duration"], best["duration"]) and reported["fuel"] == best["fuel"]),
        ("selected cost, belief, and comfort are recomputed by summing/multiplying actions", close(reported["cost"], best["cost"], 5e-4) and close(reported["belief"], best["belief"]) and close(reported["comfort"], best["comfort"])),
        ("the selected final state satisfies the wildcard goal", matches_goal(best["state"], data["Goal"]) and best["state"] == {"At": "Cologne", "Battery": "low", "Pass": "none"}),
        ("the selected plan satisfies reliability, cost, and duration thresholds", best["belief"] > data["Thresholds"]["MinBelief"] and best["cost"] < data["Thresholds"]["MaxCost"] and best["duration"] < data["Thresholds"]["MaxDuration"]),
        ("the last mile uses the high-belief Aachen-Cologne shuttle", best["actions"][-1] == "shuttle_aachen_cologne"),
        ("search depth stays within the fuel-step bound", stats["max_depth"] <= data["FuelSteps"]),
        ("the top two acceptable plans are ordered by duration then cost", plans[0]["duration"] <= plans[1]["duration"] and plans[0]["cost"] < plans[1]["cost"]),
    ]
    return run_checks(checks)
