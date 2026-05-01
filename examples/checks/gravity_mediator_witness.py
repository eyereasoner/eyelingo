# Independent Python checks for the gravity_mediator_witness example.
from __future__ import annotations

from .common import run_checks


def nonclassical_witness(run: dict) -> bool:
    return (
        run["couplingMode"] == "Gravitational"
        and "Locality" in run["assumes"]
        and "Interoperability" in run["assumes"]
        and run["directCouplingStatus"] == "NoDirectCoupling"
        and run["observed"] == "EntanglementWitnessPassed"
        and run["probeStatus"] == "LocalProbeReadoutPresent"
        and run["controlStatus"] == "CopyLikeControlPresent"
    )


def run(ctx):
    data = ctx.load_input()
    by_role = {item["role"]: item for item in data["runs"]}
    positive = by_role["positive"]
    contrast = by_role["contrast"]
    positive_yes = nonclassical_witness(positive)
    contrast_yes = nonclassical_witness(contrast)

    checks = [
        ("the positive run assumes locality", "Locality" in positive["assumes"]),
        ("the positive run assumes interoperability", "Interoperability" in positive["assumes"]),
        ("direct coupling is excluded in the positive run", positive["directCouplingStatus"] == "NoDirectCoupling"),
        ("the positive run is gravitational and mediator-only", positive["couplingMode"] == "Gravitational" and positive["mediator"] == "gravityMediator"),
        ("the entanglement witness is observed in the positive run", positive["observed"] == "EntanglementWitnessPassed"),
        ("local readout and copy-like control interfaces are both present", positive["probeStatus"] == "LocalProbeReadoutPresent" and positive["controlStatus"] == "CopyLikeControlPresent"),
        ("the non-classical mediator conclusion is independently derived for the positive run", positive_yes and "YES for the mediator-only witness run" in ctx.answer),
        ("the purely classical contrast has the same mediator-only setup but no entanglement witness", contrast["mediatorModel"] == "PurelyClassical" and contrast["directCouplingStatus"] == "NoDirectCoupling" and contrast["observed"] == "NoEntanglementWitness"),
        ("the contrast run cannot support the same witness conclusion", not contrast_yes and "NO for a purely classical mediator model" in ctx.answer),
    ]
    return run_checks(checks)
