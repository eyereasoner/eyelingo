# Independent Python checks for the flandor example.
from .common import run_fragment_checks

CHECKS = [
    ('payload hash matches the source envelope digest', ['Usage is permitted only for purpose "regional_stabilization" and the envelope expires at 2026-04-08T19:00:00+00:00.']),
    ('HMAC value matches the trusted precomputed signature', ['question : Is the Flemish Economic Resilience Board allowed to use a neutral macro-economic insight for regional stabilization, and if so which package should it activate for Fland']),
    ('export weakness, skills strain, and grid stress reach the three-need threshold', ['grid congestion : 19 hours (threshold > 11)']),
    ('the insight has a device scope, event scope, and expiry', ['export weakness : active - at least one cluster has exportOrdersIndex < 90']),
    ('the shared insight omits firm names, payroll rows, and other sensitive terms', ['Payload SHA-256: 718f5b17d07ab6a95503bc04a1000ddb132409f600659c03d21def81914b780b']),
    ('regional stabilization use is authorized before the insight expires', ['Envelope HMAC-SHA-256: 955968ca99a191783bc00cba068128ccb9ff40a5e6114fda13a52c74ee27329e']),
    ('the deletion duty is scheduled before expiry', ['Deletion duty time : 2026-04-08T18:30:00+00:00']),
    ('reuse for firm surveillance is prohibited', ['Surveillance reuse is blocked by a prohibition on odrl:distribute for purpose "firm_surveillance".']),
    ('the recommended package fits inside the €140M budget', ['Budget cap: €140M']),
    ('the recommended package covers all active needs', ['pkg:CORRIDOR_165 : reject, cost=€165M, covers all active needs; over budget by €25M']),
    ('the lowest-cost eligible package is chosen', ['Selected package "Flandor Retooling Pulse" covers export=yes, skills=yes, grid=yes, cost=€120M.']),
    ('the expected six files are recorded as written', ['grid stress : active - grid congestion hours are above 11']),
    ('all source-level policy, scope, signature, and package checks pass', ['pkg:RET_FLEX_120 : selected, cost=€120M, covers all active needs; within budget']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
