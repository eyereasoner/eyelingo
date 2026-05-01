# Independent Python checks for the deep_taxonomy_100000 example.
from __future__ import annotations

import re

from .common import run_checks


def grab_int(text: str, pattern: str):
    m = re.search(pattern, text)
    return int(m.group(1)) if m else None


def run(ctx):
    data = ctx.load_input()
    depth = int(data["TaxonomyDepth"])
    total_rules = depth + 2 + 7
    type_facts = (depth + 1) + (2 * depth) + 1
    checks = [
        ("the requested taxonomy depth is 100000", depth == 100000),
        ("the terminal class N100000 is reported as reached", f":ind a :N{depth}" in ctx.answer),
        ("the terminal A2 class and success flag are reported", ":ind a :A2" in ctx.answer and ":test :is true" in ctx.answer),
        ("taxonomy-step rule count matches the JSON depth", grab_int(ctx.reason, r"source N3 taxonomy-step rules : (\d+)") == depth),
        ("total source-rule count recomputes from start, step, terminal, success, and report rules", grab_int(ctx.reason, r"source N3 total rules counted : (\d+)") == total_rules),
        ("agenda pops match one pop per N-class fact in the main chain", grab_int(ctx.reason, r"agenda pops : (\d+)") == depth + 1),
        ("taxonomy-step applications match the depth", grab_int(ctx.reason, r"taxonomy-step rule applications : (\d+)") == depth),
        ("terminal and success rules fire exactly once", grab_int(ctx.reason, r"terminal rule applications : (\d+)") == 1 and grab_int(ctx.reason, r"success rule applications : (\d+)") == 1),
        ("classification fact total accounts for N, I/J side labels, and A2", str(type_facts) in ctx.reason and "300002 type facts" in ctx.reason),
        ("midpoint and endpoint checkpoints are present", ":N50000 plus :I50000/:J50000 present : yes" in ctx.answer and ":N99999 and :N100000 present : yes" in ctx.answer),
    ]
    return run_checks(checks)
