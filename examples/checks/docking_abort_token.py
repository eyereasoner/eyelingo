# Independent Python checks for the docking_abort_token example.
from .common import run_fragment_checks

CHECKS = [
    ('four classical media encode AbortBit', ['All four classical media encode the same abstract AbortBit variable, so the model treats them as interoperable information media.']),
    ('each classical medium can distinguish and locally permute the abort bit', ['Local permutation and local cloning are allowed on each classical medium, while copying and measuring are allowed between media that carry the same variable.']),
    ('abortLamp can copy the token to flightPLC', ['classical abort token : YES, it can be permuted, copied, measured, and composed into audit networks']),
    ('radioFrame can be measured into auditDisplay', ['parallel witness : flightPLC -> radioFrame and auditDisplay']),
    ('a serial abortLamp -> flightPLC -> auditDisplay network is possible', ['serial witness : abortLamp -> flightPLC -> auditDisplay']),
    ('a parallel flightPLC fan-out to radioFrame and auditDisplay is possible', ['Those primitive tasks compose into a serial audit path and a parallel fan-out witness.']),
    ('quantumSeal cannot universally clone all seal states', ['quantum authenticity seal : NO, it cannot be universally cloned or used as unrestricted audit fan-out']),
    ('quantumSeal cannot be used for unrestricted audit fan-out', ['The quantum seal is modeled as a superinformation medium, so universal cloning and unrestricted audit fan-out are explicitly blocked.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
