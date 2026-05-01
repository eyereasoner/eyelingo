#!/usr/bin/env python3
from __future__ import annotations

import importlib
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT))

from examples.checks.common import CheckContext


def main(argv: list[str]) -> int:
    if len(argv) != 3:
        print("usage: run_check.py EXAMPLE PREFIX_FILE", file=sys.stderr)
        return 2
    name = argv[1]
    prefix_path = Path(argv[2])
    prefix = prefix_path.read_text(encoding="utf-8")
    module = importlib.import_module(f"examples.checks.{name}")
    ok, lines = module.run(CheckContext(ROOT, name, prefix))
    print("## Check")
    for line in lines:
        print(line)
    return 0 if ok else 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
