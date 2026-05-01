# Independent Python checks for the photosynthetic_exciton_transfer example.
from .common import run_fragment_checks

CHECKS = [
    ('the tuned complex can sample exciton pathways coherently', ['The tuned complex combines strong excitonic coupling, delocalization, a tuned vibronic bridge, moderate dephasing, short-lived coherence, and a downhill route to the reaction cente']),
    ('the tuned complex can use vibronically assisted transfer', ['YES for the tuned antenna complex.']),
    ('short-lived quantum assistance is enough in the tuned downhill regime', ['The detuned contrast complex has weak coupling, absent delocalization, no vibronic bridge, strong dephasing, and a trapping mismatch, so the same efficient delivery task is blocked']),
    ('efficient exciton transfer is possible in the tuned complex', ['NO for the detuned, strongly decohered contrast complex.']),
    ('the tuned complex can deliver excitation to the reaction center', ['## Answer']),
    ('the detuned complex cannot sample pathways coherently', ['## Answer']),
    ('the detuned complex cannot use vibronically assisted transfer', ['## Answer']),
    ('the detuned complex cannot achieve directed reaction-center transfer', ['## Answer']),
    ('the detuned complex cannot achieve efficient exciton transfer', ['## Answer']),
    ('the detuned complex cannot deliver excitation efficiently to the reaction center', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
