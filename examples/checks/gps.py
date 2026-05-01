# Independent Python checks for the gps example.
from __future__ import annotations

import re
from collections import defaultdict

from .common import run_checks


def enumerate_paths(edges: list[dict], start: str, goal: str) -> list[list[dict]]:
    by_from = defaultdict(list)
    for edge in edges:
        by_from[edge["From"]].append(edge)
    paths = []

    def dfs(node: str, path: list[dict], seen: set[str]):
        if node == goal:
            paths.append(path[:])
            return
        for edge in by_from[node]:
            if edge["To"] in seen:
                continue
            path.append(edge)
            seen.add(edge["To"])
            dfs(edge["To"], path, seen)
            seen.remove(edge["To"])
            path.pop()

    dfs(start, [], {start})
    return paths


def route_metrics(path: list[dict]) -> tuple[float, float, float, float]:
    duration = sum(float(e["Duration"]) for e in path)
    cost = sum(float(e["Cost"]) for e in path)
    belief = 1.0
    comfort = 1.0
    for edge in path:
        belief *= float(edge["Belief"])
        comfort *= float(edge["Comfort"])
    return duration, cost, belief, comfort


def label_for_path(path: list[dict]) -> str:
    nodes = [path[0]["From"]] + [edge["To"] for edge in path]
    return " → ".join(nodes)


def parse_goal(question: str) -> str:
    match = re.search(r"to ([A-Za-zÀ-ÿ]+)", question)
    if not match:
        raise ValueError(question)
    return match.group(1).rstrip("?")


def run(ctx):
    data = ctx.load_input()
    start = data["Traveller"]["Location"]
    goal = parse_goal(data["Question"])
    paths = enumerate_paths(data["Edges"], start, goal)
    metrics_by_label = {label_for_path(path): route_metrics(path) for path in paths}
    direct = data["Routes"]["routeDirect"]["Label"]
    alt = data["Routes"]["routeViaKortrijk"]["Label"]
    best_label = min(metrics_by_label, key=lambda label: metrics_by_label[label][0])
    d_duration, d_cost, d_belief, d_comfort = metrics_by_label[direct]
    a_duration, a_cost, a_belief, a_comfort = metrics_by_label[alt]

    checks = [
        ("the direct Gent-Brugge-Oostende route is derived from edges", direct in metrics_by_label),
        ("the Kortrijk alternative route is derived from edges", alt in metrics_by_label),
        ("exactly two simple routes connect the traveller to the destination", len(metrics_by_label) == 2),
        ("route duration and cost are additive over edges", d_duration == 2400.0 and abs(d_cost - 0.010) <= 1e-12 and a_duration == 4100.0 and abs(a_cost - 0.018) <= 1e-12),
        ("route belief and comfort are multiplicative over edges", abs(d_belief - 0.9408) <= 1e-12 and abs(d_comfort - 0.99) <= 1e-12 and abs(a_belief - 0.903168) <= 1e-12 and abs(a_comfort - 0.9801) <= 1e-12),
        ("the recommended route is faster than the alternative", d_duration < a_duration and best_label == direct),
        ("the recommended route is cheaper than the alternative", d_cost < a_cost),
        ("the recommended route has higher belief and comfort scores", d_belief > a_belief and d_comfort > a_comfort),
        ("the answer names the independently chosen direct route", direct in ctx.answer and "Take the direct route" in ctx.answer),
    ]
    return run_checks(checks)
