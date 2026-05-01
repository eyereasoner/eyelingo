#!/usr/bin/env python3
from __future__ import annotations

import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def strip_go_tail(text: str) -> str:
    for marker in ("\n## Check", "\n## Go audit details"):
        pos = text.find(marker)
        if pos >= 0:
            text = text[:pos]
    return text.rstrip()


def main(argv: list[str]) -> int:
    if len(argv) != 3:
        print("usage: build_output.py EXAMPLE GO_STDOUT_FILE", file=sys.stderr)
        return 2
    name = argv[1]
    raw_path = Path(argv[2])
    prefix = strip_go_tail(raw_path.read_text(encoding="utf-8"))
    prefix_path = raw_path.with_suffix(raw_path.suffix + ".prefix")
    prefix_path.write_text(prefix + "\n", encoding="utf-8")
    proc = subprocess.run(
        [sys.executable, str(ROOT / "tools" / "run_check.py"), name, str(prefix_path)],
        cwd=str(ROOT),
        text=True,
        capture_output=True,
    )
    if proc.returncode != 0:
        if proc.stdout:
            print(proc.stdout, end="")
        if proc.stderr:
            print(proc.stderr, file=sys.stderr, end="")
        return proc.returncode
    print(prefix)
    print()
    print(proc.stdout.rstrip())
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
