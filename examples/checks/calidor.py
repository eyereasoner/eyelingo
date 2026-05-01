# Independent Python checks for the calidor example.
from .common import run_fragment_checks

CHECKS = [
    ('the trusted precomputed HMAC signature verifies', ['Usage is permitted only for purpose "heatwave_response" and the envelope expires at 2026-07-18T21:00:00+00:00.']),
    ('payload hash matches the source envelope digest', ['reason.txt : The gateway keeps raw indoor heat, vulnerability, and prepaid-energy data local, derives a priority-support signal, and shares only a scoped heatwave-response envelope']),
    ('the shared insight strips raw heat, vulnerability, credit, and meter-trace terms', ['The gateway desensitizes local heat, vulnerability, and prepaid-energy stress into an expiring municipal support insight.']),
    ('the insight has a device scope, event scope, and expiry', ['question : Is the Calidor heat-response system allowed to use a narrow household support insight for heatwave response, and if so which support package should it recommend?']),
    ('heatwave-response use is authorized before the insight expires', ['expires at : 2026-07-18T21:00:00+00:00']),
    ('the heat alert is active', ['unsafe indoor heat : active - 31.4°C for 9 hours reaches the 30.0°C/6 hour threshold']),
    ('indoor heat is unsafe for long enough to trigger support', ['Tenant-screening reuse is blocked by a prohibition on odrl:distribute for purpose "tenant_screening".']),
    ('the active-need count reaches the priority threshold', ['energy constraint : active - €3.2 prepaid credit is at or below the €5.0 limit']),
    ('the recommended package is within budget and covers every required capability', ['pkg:VOUCHER  : reject, cost=€12, covers 2/4 required capabilities; within budget']),
    ('the lowest-cost eligible package is chosen', ['Selected package "Calidor Priority Cooling Bundle" covers bill_credit, cooling_kit, transport, welfare_check.']),
    ('the deletion duty is scheduled before expiry', ['vulnerability present : active - the local gateway sees heat-sensitive and mobility flags']),
    ('reuse for tenant screening is prohibited', ['Envelope HMAC-SHA-256: e635c7c1991742a5c36992fc0da32a7abc80b32aa5777a1142adaab55183681c']),
    ('all source-level policy, scope, signature, and package checks pass', ['pkg:DELUXE   : reject, cost=€28, covers all required capabilities; over budget by €8']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
