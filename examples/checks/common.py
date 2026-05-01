from __future__ import annotations

import json
from pathlib import Path
from typing import Callable, Iterable


class CheckContext:
    def __init__(self, root: Path, name: str, prefix: str):
        self.root = Path(root)
        self.name = name
        self.prefix = prefix.rstrip() + "\n"
        self.input_path = self.root / "examples" / "input" / f"{name}.json"
        self._input_loaded = False
        self._input = None

    def load_input(self):
        if not self._input_loaded:
            with self.input_path.open("r", encoding="utf-8") as f:
                self._input = json.load(f)
            self._input_loaded = True
        return self._input

    def section(self, heading: str) -> str:
        marker = f"## {heading}"
        if marker not in self.prefix:
            return ""
        tail = self.prefix.split(marker, 1)[1]
        next_heading = tail.find("\n## ")
        if next_heading >= 0:
            tail = tail[:next_heading]
        return tail.strip()

    @property
    def answer(self) -> str:
        return self.section("Answer")

    @property
    def reason(self) -> str:
        return self.section("Reason why")

    def contains_all(self, fragments):
        if isinstance(fragments, str):
            fragments = [fragments]
        missing = [fragment for fragment in fragments if fragment not in self.prefix]
        return not missing, missing


def check_line(ok: bool, index: int, description: str) -> str:
    return f"C{index} {'OK' if ok else 'FAIL'} - {description}"


def run_checks(specs: Iterable[tuple[str, bool]]) -> tuple[bool, list[str]]:
    lines = []
    all_ok = True
    for index, (description, ok) in enumerate(specs, 1):
        ok = bool(ok)
        lines.append(check_line(ok, index, description))
        all_ok = all_ok and ok
    return all_ok, lines


def run_fragment_checks(ctx: CheckContext, specs):
    # Legacy helper used by examples that have not yet been upgraded to a
    # domain-specific Python verifier.
    ctx.load_input()
    lines = []
    all_ok = True
    for index, (description, fragments) in enumerate(specs, 1):
        ok, _missing = ctx.contains_all(fragments)
        lines.append(check_line(ok, index, description))
        all_ok = all_ok and ok
    return all_ok, lines
