# Independent Python checks for the french_cities example.
from __future__ import annotations

from collections import defaultdict, deque

from .common import run_checks


def reachable(graph: dict[str, list[str]], start: str, goal: str) -> list[str] | None:
    q = deque([(start, [start])])
    seen = {start}
    while q:
        node, path = q.popleft()
        if node == goal:
            return path
        for nxt in graph.get(node, []):
            if nxt not in seen:
                seen.add(nxt)
                q.append((nxt, path + [nxt]))
    return None


def run(ctx):
    data = ctx.load_input()
    graph = defaultdict(list)
    for edge in data["Edges"]:
        graph[edge["From"]].append(edge["To"])
    labels = data["Labels"]
    goal = "nantes"
    paths = {city: reachable(graph, city, goal) for city in labels}
    reaching = {city: path for city, path in paths.items() if path and city != goal}
    reported_names = {name.strip() for name in ctx.answer.split("reach Nantes:", 1)[1].split(".", 1)[0].split(",")} if "reach Nantes:" in ctx.answer else set()
    expected_names = {labels[city] for city in reaching}

    checks = [
        ("Angers has a direct one-way edge to Nantes", paths["angers"] == ["angers", "nantes"]),
        ("Le Mans reaches Nantes through Angers", paths["lemans"] == ["lemans", "angers", "nantes"]),
        ("Chartres reaches Nantes through Le Mans and Angers", paths["chartres"] == ["chartres", "lemans", "angers", "nantes"]),
        ("Paris reaches Nantes through Chartres, Le Mans, and Angers", paths["paris"] == ["paris", "chartres", "lemans", "angers", "nantes"]),
        ("exactly four non-destination cities can reach Nantes", len(reaching) == 4),
        ("the reported city set matches the transitive closure", reported_names == expected_names),
        ("cities without any chain to Nantes are rejected", all(paths[c] is None for c in ["amiens", "orleans", "blois", "bourges", "tours"])),
    ]
    return run_checks(checks)
