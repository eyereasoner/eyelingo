# Independent Python checks for the drone_corridor_planner example.
from .common import run_fragment_checks

CHECKS = [
    ('14 corridor actions were loaded from JSON', ['The planner treats each corridor description as a state transition over location, battery, and permit state.']),
    ('bounded search found 17 plans meeting belief and cost thresholds', ['Plans are retained only when belief is greater than 0.94 and cost is less than 0.03.']),
    ('lowest-cost selected plan starts with fly_gent_brugge', ['The selected plan is the lowest-cost surviving plan; the next cheapest starts with fly_gent_brugge and costs 0.014.']),
    ('selected plan reaches Oostende', ['selected plan : fly_gent_brugge -> public_coastline_brugge_oostende']),
    ('selected belief 0.960300 is above 0.94', ['belief : 0.960300']),
    ('selected cost 0.012 is below 0.03', ['cost : 0.012']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
