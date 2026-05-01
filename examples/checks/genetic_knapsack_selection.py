# Independent Python checks for the genetic_knapsack_selection example.
from __future__ import annotations

import re

from .common import run_checks


def parse_answer(answer: str):
    def int_field(label: str) -> int | None:
        match = re.search(rf"{re.escape(label)}\s*:\s*(\d+)", answer)
        return int(match.group(1)) if match else None

    genome_match = re.search(r"final genome\s*:\s*([01]+)", answer)
    items_match = re.search(r"selected items\s*:\s*(.+)", answer)
    weight_match = re.search(r"weight\s*:\s*(\d+)\s*/\s*(\d+)", answer)
    return {
        "genome": genome_match.group(1) if genome_match else "",
        "items": [part.strip() for part in items_match.group(1).split(",")] if items_match else [],
        "weight": int(weight_match.group(1)) if weight_match else None,
        "capacity": int(weight_match.group(2)) if weight_match else None,
        "value": int_field("value"),
        "fitness": int_field("fitness"),
        "generations": int_field("generations evaluated"),
    }


def evaluate(genome: str, items: list[dict], capacity: int):
    weight = sum(int(item["Weight"]) for bit, item in zip(genome, items) if bit == "1")
    value = sum(int(item["Value"]) for bit, item in zip(genome, items) if bit == "1")
    fitness = 1_000_000 - value if weight <= capacity else 2_000_000 + (weight - capacity)
    selected = [item["Name"] for bit, item in zip(genome, items) if bit == "1"]
    return {"genome": genome, "weight": weight, "value": value, "fitness": fitness, "items": selected}


def mutants(genome: str) -> list[str]:
    out = []
    for index, bit in enumerate(genome):
        flipped = "1" if bit == "0" else "0"
        out.append(genome[:index] + flipped + genome[index + 1:])
    return out


def select_best(genomes: list[str], items: list[dict], capacity: int):
    candidates = [evaluate(genome, items, capacity) for genome in genomes]
    return min(candidates, key=lambda candidate: (candidate["fitness"], candidate["genome"]))


def simulate(data: dict):
    genome = data["StartGenome"]
    history = []
    for generation in range(int(data["MaxGenerations"]) + 1):
        parent = evaluate(genome, data["Items"], int(data["Capacity"]))
        best = select_best([genome] + mutants(genome), data["Items"], int(data["Capacity"]))
        history.append((generation, parent, best))
        if best["genome"] == genome:
            break
        genome = best["genome"]
    return evaluate(genome, data["Items"], int(data["Capacity"])), history


def run(ctx):
    data = ctx.load_input()
    parsed = parse_answer(ctx.answer)
    final, history = simulate(data)
    final_neighbors = [evaluate(genome, data["Items"], int(data["Capacity"])) for genome in mutants(final["genome"])]
    best_neighbor = min(final_neighbors, key=lambda candidate: (candidate["fitness"], candidate["genome"]))
    global_feasible_count = sum(
        1
        for value in range(1 << len(data["Items"]))
        if evaluate(format(value, f"0{len(data['Items'])}b"), data["Items"], int(data["Capacity"]))["weight"] <= int(data["Capacity"])
    )

    checks = [
        ("the input has one item per genome bit and a valid binary start genome", len(data["Items"]) == len(data["StartGenome"]) and set(data["StartGenome"]) <= {"0", "1"}),
        ("Python independently simulates the deterministic single-bit local search to the reported final genome", parsed["genome"] == final["genome"]),
        ("reported selected items are exactly the one-bits in the final genome", parsed["items"] == final["items"]),
        ("reported weight, value, and fitness match independent genome evaluation", parsed["weight"] == final["weight"] and parsed["value"] == final["value"] and parsed["fitness"] == final["fitness"]),
        ("the final candidate is feasible and matches the expected fixture totals", final["weight"] <= int(data["Capacity"]) and final["genome"] == data["Expected"]["FinalGenome"] and final["weight"] == data["Expected"]["FinalWeight"] and final["value"] == data["Expected"]["FinalValue"]),
        ("no one-bit neighbor has a lower fitness than the final candidate", best_neighbor["fitness"] >= final["fitness"]),
        ("the reported generation count matches the independent simulation history length", parsed["generations"] == len(history)),
        ("the fixture has many feasible genomes, so the check validates the local-search rule rather than a text fragment", global_feasible_count > 1000),
    ]
    return run_checks(checks)
