# Independent Python checks for the gps example.
from .common import run_fragment_checks

CHECKS = [
    ('the direct Gent → Brugge → Oostende route was derived', ['The direct route (Gent → Brugge → Oostende) takes 2400.0 seconds at cost 0.01, with belief 0.9408 and comfort 0.99. The alternative (Gent → Kortrijk → Brugge → Oostende) takes 4100']),
    ('the alternative Gent → Kortrijk → Brugge → Oostende route was derived', ['Recommended route: Gent → Brugge → Oostende']),
    ('the recommended route is faster than the alternative', ['So the direct route is faster, cheaper, more reliable, and slightly more comfortable.']),
    ('the recommended route is cheaper than the alternative', ['Take the direct route via Brugge.']),
    ('the recommended route has higher belief and comfort scores', ['From Gent to Oostende, the planner found two routes in this small map.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
