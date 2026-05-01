# Independent Python checks for the path_discovery example.
from __future__ import annotations

import re
from collections import defaultdict

from .common import run_checks


def discover(data: dict, max_stopovers: int = 2):
    graph = defaultdict(list)
    for edge in data["Edges"]:
        graph[edge["From"]].append(edge["To"])
    for node in graph:
        graph[node].sort(key=lambda x: data["Labels"].get(x, x))
    max_edges = max_stopovers + 1
    routes = []
    calls = 0
    edge_tests = 0

    def dfs(node, path):
        nonlocal calls, edge_tests
        calls += 1
        if node == data["DestinationID"]:
            routes.append(path[:])
            return
        if len(path) - 1 == max_edges:
            return
        for nxt in graph.get(node, []):
            edge_tests += 1
            if nxt in path:
                continue
            dfs(nxt, path + [nxt])

    dfs(data["SourceID"], [data["SourceID"]])
    routes.sort(key=lambda r: [data["Labels"][x] for x in r])
    return routes, graph, calls, edge_tests


def parse_routes(answer: str):
    rows = []
    for _idx, labels in re.findall(r"route (\d+) \(2 stopovers\): ([^\n]+)", answer):
        rows.append([p.strip() for p in labels.split("->")])
    return rows


def run(ctx):
    data = ctx.load_input()
    routes, graph, calls, edge_tests = discover(data)
    labels = data["Labels"]
    label_routes = [[labels[x] for x in route] for route in routes]
    reported = parse_routes(ctx.answer)
    second_hop_labels = sorted(labels[x] for x in graph[graph[data["SourceID"]][0]])
    expected_dest = labels[data["DestinationID"]]

    checks = [
        ("source and destination airport labels are known", labels[data["SourceID"]] == "Ostend-Bruges International Airport" and expected_dest == "Václav Havel Airport Prague"),
        ("Ostend-Bruges has one outbound route in the full graph, to Liège Airport", len(graph[data["SourceID"]]) == 1 and labels[graph[data["SourceID"]][0]] == "Liège Airport"),
        ("bounded DFS independently finds exactly three two-stopover routes", len(routes) == 3 and all(len(r) == 4 for r in routes)),
        ("reported route labels match the independently discovered route set", reported == label_routes),
        ("no direct or one-stop route exists under the same bound", all(len(r) != 2 and len(r) != 3 for r in routes)),
        ("every discovered hop is backed by an outbound-route fact", all(route[i + 1] in graph[route[i]] for route in routes for i in range(len(route) - 1))),
        ("no discovered route revisits an airport", all(len(route) == len(set(route)) for route in routes)),
        ("the translated graph size matches the full source counts", len(labels) == 7698 and sum(len(v) for v in graph.values()) == 37505),
        ("the second-hop candidates from Liège are independently recovered", second_hop_labels == ["Ajaccio-Napoléon Bonaparte Airport", "Al Massira Airport", "Bastia-Poretta Airport", "Diagoras Airport", "Heraklion International Nikos Kazantzakis Airport", "Lille-Lesquin Airport", "Palma De Mallorca Airport"]),
        ("route output is sorted deterministically by airport labels", label_routes == sorted(label_routes)),
    ]
    return run_checks(checks)
