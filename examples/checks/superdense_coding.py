# Independent Python checks for the superdense_coding example.
from __future__ import annotations

import re
from collections import Counter

from .common import run_checks


Pair = tuple[bool, bool]


def relation(rows: list[list[bool]]) -> set[Pair]:
    # The JSON serializes each fact row as two mobit pairs: [a,b,c,d]
    # means relation contains (a,b) and (c,d). Duplicates collapse.
    out: set[Pair] = set()
    for row in rows:
        if len(row) != 4:
            raise ValueError("relation rows must contain four booleans")
        out.add((row[0], row[1]))
        out.add((row[2], row[3]))
    return out


def compose(first: set[Pair], second: set[Pair]) -> set[Pair]:
    return {(a, d) for (a, b) in first for (c, d) in second if b == c}


def derive(data: dict):
    states = {name: relation(rows) for name, rows in data["States"].items()}
    primitive = {name: relation(rows) for name, rows in data["Primitive"].items()}
    composed = {
        "kg": compose(primitive["g"], primitive["k"]),
        "gk": compose(primitive["k"], primitive["g"]),
    }
    lookup = lambda name: primitive.get(name, composed.get(name))
    alice = {int(k): v for k, v in data["AliceOps"].items()}
    # The Go translation names the four Bob tests by the relations that accept
    # each encoded support after the joint measurement.
    bob = {0: "gk", 1: "k", 2: "g", 3: "id"}
    counts = Counter()
    candidates = []
    encoded = {}
    for msg in range(4):
        alice_rel = lookup(alice[msg])
        encoded[msg] = set()
        for shared_a, shared_b in sorted(states["R"]):
            for move_a, move_b in sorted(alice_rel):
                if move_a != shared_a:
                    continue
                encoded[msg].add((move_b, shared_b))
                for decoded in range(4):
                    if (move_b, shared_b) in lookup(bob[decoded]):
                        counts[(msg, decoded)] += 1
                        candidates.append((msg, decoded, shared_a, move_b, shared_b))
    survivors = {msg: decoded for msg in range(4) for decoded in range(4) if counts[(msg, decoded)] % 2 == 1}
    return states, primitive, composed, counts, candidates, survivors, encoded


def parse_survivors(answer: str):
    return {int(a): int(b) for a, b in re.findall(r"^\s*(\d+) dqc:superdense-coding (\d+)", answer, flags=re.MULTILINE)}


def run(ctx):
    data = ctx.load_input()
    states, primitive, composed, counts, candidates, survivors, encoded = derive(data)
    reported = parse_survivors(ctx.answer)
    diagonal_counts = [counts[(i, i)] for i in range(4)]
    offdiag_even = all(counts[(i, j)] % 2 == 0 for i in range(4) for j in range(4) if i != j)

    checks = [
        ("shared entanglement R contains exactly |0,0) and |1,1)", states["R"] == {(False, False), (True, True)}),
        ("composition KG is recomputed by composing G then K", composed["kg"] == {(False, False), (False, True), (True, False)}),
        ("composition GK is recomputed by composing K then G", composed["gk"] == {(False, True), (True, False), (True, True)}),
        ("the raw superdense rule creates 24 candidate derivations", len(candidates) == 24),
        ("GF(2) cancellation leaves odd diagonal and even off-diagonal counts", diagonal_counts == [1, 1, 1, 1] and offdiag_even),
        ("reported surviving decoded messages match the parity survivors", reported == survivors == {0: 0, 1: 1, 2: 2, 3: 3}),
        ("the four Alice operations produce distinct encoded supports", len({tuple(sorted(v)) for v in encoded.values()}) == 4),
        ("the JSON relation facts match the primitive teaching-model relations", primitive["id"] == states["R"] and primitive["g"] == states["S"] and primitive["k"] == states["U"]),
    ]
    return run_checks(checks)
