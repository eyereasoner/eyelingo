# Independent Python checks for the delfour example.
from .common import run_fragment_checks

CHECKS = [
    ('signature verifies', ['signature alg : HMAC-SHA256']),
    ('payload hash matches', ['expires at : 2025-10-05T22:33:48.907185+00:00']),
    ('minimization strips sensitive terms', ["reason.txt : Household requires low-sugar guidance (diabetes in POD). A neutral Insight is scoped to device 'self-scanner', event 'pick_up_scanner', retailer 'Delfour', and expires"]),
    ('scope complete', ['scope : self-scanner @ pick_up_scanner']),
    ('authorization allowed', ['The scanner is allowed to use a neutral shopping insight and recommends Low-Sugar Tea Biscuits instead of Classic Tea Biscuits.']),
    ('high-sugar banner', ['banner headline : Track sugar per serving while you scan']),
    ('alternative lowers sugar', ['suggested alternative: Low-Sugar Tea Biscuits']),
    ('duty timing consistent', ['The phone desensitizes a diabetes-related household condition into a scoped low-sugar need, wraps it in an expiring Insight + Policy envelope, signs it, and the scanner consumes th']),
    ('marketing prohibited', ['threshold : 10.0']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
