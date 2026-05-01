# Independent Python checks for the harbor_smr example.
from .common import run_fragment_checks

CHECKS = [
    ('reserve margin 24 MW exceeds threshold 19 MW', ['PERMIT - North Quay Hydrogen Hub may use https://example.org/insight/harborsmr to run PEM_electrolyzer_train_2 at 16 MW from 2026-06-18T14:00:00Z to 2026-06-18T18:00:00Z.']),
    ('cooling margin 18% exceeds threshold 14%', ['The SMR operator exposes a bounded 18 MW flexible-export insight for day_ahead_balancing, not raw reactor telemetry.']),
    ('no planned outage blocks the balancing window', ['The requested 16 MW electrolysis dispatch fits inside that window, safety margins clear the thresholds, no outage is planned, and the policy permits use only for electrolysis_dispa']),
    ('requested 16 MW fits inside the 18 MW flexible-export insight', ['The approved dispatch is 64 MWh over the four-hour window, scoped to port-hydrogen-hub and PEM_electrolyzer_train_2.']),
    ('serialized insight omits sensitive telemetry terms', ['## Answer']),
    ('aggregate flags confirm raw reactor telemetry stays local', ['## Answer']),
    ('policy permits use for electrolysis dispatch before the insight expires', ['## Answer']),
    ('policy prohibits redistribution for market resale', ['## Answer']),
    ('scope is explicit: device, event, start, and expiry', ['## Answer']),
    ('dispatch plan matches the requested load, power, and insight window', ['## Answer']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
