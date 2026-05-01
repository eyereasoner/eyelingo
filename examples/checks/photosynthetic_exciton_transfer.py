# Independent Python checks for the photosynthetic_exciton_transfer example.
from __future__ import annotations

from .common import run_checks


def efficient_delivery(complex_: dict, reaction_center: str) -> bool:
    return (
        complex_["excitonCoupling"] == "Strong"
        and complex_["stateDelocalization"] == "Present"
        and complex_["vibronicBridge"] == "Tuned"
        and complex_["energyLandscape"] == "DownhillToReactionCenter"
        and complex_["dephasing"] in {"Moderate", "Low"}
        and complex_["electronicCoherenceLifetime"] in {"Short", "Long"}
        and complex_["connectedTo"] == reaction_center
    )


def run(ctx):
    data = ctx.load_input()
    complexes = {item["role"]: item for item in data["complexes"]}
    tuned = complexes["positive"]
    detuned = complexes["contrast"]
    tuned_ok = efficient_delivery(tuned, data["reactionCenter"])
    detuned_ok = efficient_delivery(detuned, data["reactionCenter"])

    checks = [
        ("the tuned complex has strong excitonic coupling", tuned["excitonCoupling"] == "Strong"),
        ("the tuned complex has delocalized states", tuned["stateDelocalization"] == "Present"),
        ("the tuned complex has a tuned vibronic bridge", tuned["vibronicBridge"] == "Tuned"),
        ("the tuned complex has moderate dephasing and short-lived coherence", tuned["dephasing"] == "Moderate" and tuned["electronicCoherenceLifetime"] == "Short"),
        ("the tuned complex is connected downhill to the reaction center", tuned["energyLandscape"] == "DownhillToReactionCenter" and tuned["connectedTo"] == data["reactionCenter"]),
        ("efficient transfer is independently derived for the tuned complex", tuned_ok and "YES for the tuned antenna complex" in ctx.answer),
        ("the detuned complex lacks strong coupling, delocalization, and a vibronic bridge", detuned["excitonCoupling"] == "Weak" and detuned["stateDelocalization"] == "Absent" and detuned["vibronicBridge"] == "Absent"),
        ("the detuned complex has strong dephasing and a trapping mismatch", detuned["dephasing"] == "Strong" and detuned["energyLandscape"] == "TrappingMismatch"),
        ("efficient delivery is blocked for the detuned contrast complex", not detuned_ok and "NO for the detuned" in ctx.answer),
        ("the reaction-center connection alone is insufficient without the other conditions", detuned["connectedTo"] == data["reactionCenter"] and not detuned_ok),
    ]
    return run_checks(checks)
