# Independent Python checks for the dijkstra_risk_path example.
from __future__ import annotations

import heapq
import math
import re
from collections import defaultdict

from .common import run_checks


def parse_answer(answer: str):
    def number(label: str) -> float | None:
        match = re.search(rf"{re.escape(label)}\s*:\s*([0-9]+(?:\.[0-9]+)?)", answer)
        return float(match.group(1)) if match else None

    path_match = re.search(r"selected path\s*:\s*(.+)", answer)
    edge_match = re.search(r"edges in selected path\s*:\s*(\d+)", answer)
    path = [part.strip() for part in path_match.group(1).split("->")] if path_match else []
    return {
        "path": path,
        "raw_cost": number("raw cost"),
        "risk_sum": number("risk sum"),
        "score": number("risk-adjusted score"),
        "edges": int(edge_match.group(1)) if edge_match else None,
    }


def edge_score(edge: dict, risk_weight: float) -> float:
    return float(edge["cost"]) + risk_weight * float(edge["risk"])


def path_totals(edges: list[dict], path: list[str], risk_weight: float):
    by_pair = {(edge["from"], edge["to"]): edge for edge in edges}
    raw = 0.0
    risk = 0.0
    score = 0.0
    for source, target in zip(path, path[1:]):
        edge = by_pair.get((source, target))
        if edge is None:
            return None
        raw += float(edge["cost"])
        risk += float(edge["risk"])
        score += edge_score(edge, risk_weight)
    return raw, risk, score


def dijkstra(data: dict):
    graph = defaultdict(list)
    for edge in data["edges"]:
        graph[edge["from"]].append(edge)
    risk_weight = float(data["riskWeight"])
    heap = [(0.0, data["start"], (data["start"],))]
    best = {}
    while heap:
        score, node, path = heapq.heappop(heap)
        if node in best and best[node] <= score:
            continue
        best[node] = score
        if node == data["goal"]:
            return list(path), score, best
        for edge in graph[node]:
            heapq.heappush(heap, (score + edge_score(edge, risk_weight), edge["to"], path + (edge["to"],)))
    return [], math.inf, best


def enumerate_simple_paths(data: dict):
    graph = defaultdict(list)
    for edge in data["edges"]:
        graph[edge["from"]].append(edge["to"])
    out = []

    def search(node: str, path: list[str]) -> None:
        if node == data["goal"]:
            out.append(path[:])
            return
        for nxt in graph[node]:
            if nxt not in path:
                path.append(nxt)
                search(nxt, path)
                path.pop()

    search(data["start"], [data["start"]])
    return out


def run(ctx):
    data = ctx.load_input()
    parsed = parse_answer(ctx.answer)
    risk_weight = float(data["riskWeight"])
    totals = path_totals(data["edges"], parsed["path"], risk_weight)
    best_path, best_score, visited = dijkstra(data)
    all_paths = enumerate_simple_paths(data)
    all_scores = [path_totals(data["edges"], path, risk_weight)[2] for path in all_paths]
    depot_c_paths = [path for path in all_paths if "DepotC" in path]
    depot_c_best = min(path_totals(data["edges"], path, risk_weight)[2] for path in depot_c_paths)

    checks = [
        ("the graph fixture has the requested start, goal, risk weight, and eight directed edges", data["start"] == "ClinicA" and data["goal"] == "HubZ" and risk_weight == 2.0 and len(data["edges"]) == 8),
        ("every edge weight is independently computed as cost + riskWeight × risk", all(edge_score(edge, risk_weight) > float(edge["cost"]) for edge in data["edges"])),
        ("the reported path is made only of directed edges present in the input graph", totals is not None and parsed["path"][0] == data["start"] and parsed["path"][-1] == data["goal"]),
        ("reported raw cost, risk sum, score, and edge count match the parsed path", totals is not None and abs(parsed["raw_cost"] - totals[0]) < 5e-9 and abs(parsed["risk_sum"] - totals[1]) < 5e-9 and abs(parsed["score"] - totals[2]) < 5e-9 and parsed["edges"] == len(parsed["path"]) - 1),
        ("an independent Python Dijkstra run selects ClinicA -> DepotB -> LabD -> HubZ", best_path == ["ClinicA", "DepotB", "LabD", "HubZ"]),
        ("the independent shortest-path score matches the fixture expectation", abs(best_score - float(data["expected"]["score"])) < 1e-9),
        ("exhaustive simple-path enumeration finds no lower risk-adjusted score", all(score >= best_score - 1e-12 for score in all_scores)),
        ("the best route through DepotC is independently more expensive than the selected route", depot_c_best > best_score and "DepotC" not in parsed["path"]),
    ]
    return run_checks(checks)
