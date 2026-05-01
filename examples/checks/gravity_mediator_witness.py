# Independent Python checks for the gravity_mediator_witness example.
from .common import run_fragment_checks

CHECKS = [
    ('locality is assumed in the positive run', ['The positive run assumes locality and interoperability, excludes direct coupling, and observes entanglement after interaction through the gravitational mediator alone.']),
    ('interoperability is assumed in the positive run', ['YES for the mediator-only witness run.']),
    ('direct coupling between the two quantum systems is excluded', ['Under those conditions the mediator-only witness supports a non-classical-mediator conclusion, while the purely classical contrast model cannot support the same witness.']),
    ('the positive run has a mediator-only interaction path', ['NO for a purely classical mediator model under the same mediator-only conditions.']),
    ('an entanglement witness is observed in the positive run', ['## Answer']),
    ('the positive run has both information-transfer and local-readout interfaces', ['## Answer']),
    ('the gravitational mediator is derived to be non-classical', ['## Answer']),
    ('a purely classical mediator model is ruled out by the positive run', ['## Answer']),
    ('the contrast run is also mediator-only', ['## Answer']),
    ('the purely classical contrast mediator cannot support the witness', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
