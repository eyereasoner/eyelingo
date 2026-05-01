# Independent Python checks for the french_cities example.
from .common import run_fragment_checks

CHECKS = [
    ('Angers has a direct one-way connection to Nantes', ['The original example says that every :oneway link is also a :path, and that :path is transitive. So once Angers can reach Nantes directly, longer routes can be built by chaining ea']),
    ('Le Mans reaches Nantes by chaining Le Mans → Angers → Nantes', ['Four cities in this small network can reach Nantes: Paris, Chartres, Le Mans, Angers.']),
    ('Chartres reaches Nantes by chaining Chartres → Le Mans → Angers → Nantes', ['## Answer']),
    ('Paris reaches Nantes by chaining Paris → Chartres → Le Mans → Angers → Nantes', ['## Answer']),
    ('cities without a chain to Nantes are rejected by fail-loud fuse rules', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
