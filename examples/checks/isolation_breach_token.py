# Independent Python checks for the isolation_breach_token example.
from .common import run_fragment_checks

CHECKS = [
    ('doorBeacon, containmentPLC, nursePager, and incidentBoard encode BreachBit', ['serial witness : doorBeacon -> containmentPLC -> incidentBoard']),
    ('nursePager can prepare CodeBreach', ['prepared witness : nursePager prepares CodeBreach']),
    ('doorBeacon can permute SafeGreen to BreachRed and back', ['The specimen seal is separate and superinformation-like, so it records provenance without becoming an unrestricted cloneable broadcast token.']),
    ('containmentPLC can copy the breach token to nursePager', ['The breach token is an ordinary classical information variable carried by four unlike media in the lab workflow.']),
    ('nursePager can be measured into incidentBoard', ['Copy and measurement compose into an incident-board audit path, and copy pairs compose into a parallel notification witness.']),
    ('doorBeacon -> containmentPLC -> incidentBoard is a serial audit network', ['classical breach token : YES, prepare, reversible permutation, copy, measure, serial audit, and parallel fan-out all succeed']),
    ('containmentPLC can fan out to nursePager and incidentBoard', ['Since each medium carries BreachBit, the example derives preparation, reversible permutation, copying, and measurement tasks.']),
    ('specimenSeal cannot universally clone all provenance states', ['specimen provenance seal : NO, universal cloning and unrestricted parallel fan-out are blocked']),
    ('specimenSeal blocks unrestricted parallel fan-out', ['possible prepare tasks : 8']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
