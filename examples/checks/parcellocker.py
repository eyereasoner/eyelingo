# Independent Python checks for the parcellocker example.
from .common import run_fragment_checks

CHECKS = [
    ('the source pickup request satisfies all ten authorization conditions', ['pickup        : PERMIT (source request)']),
    ('the same token is denied for billing, redirect, wrong person, wrong locker, and reuse', ['reason : Permit: the requester, parcel, locker, action, purpose, active state, single-use limit, parcel readiness, and privacy restrictions all match.']),
    ('every request matches its expected PERMIT or DENY result', ['request : requester=noah parcel=parcel123 locker=lockerB17 action=ViewBillingDetails purpose=BillingAccess']),
    ('single-use state permits the first pickup but rejects reuse', ['reason : Deny: C4 requested action must be parcel collection; C5 requested purpose must be pickup only; C10 parcel redirection must stay blocked.']),
    ('billing details stay hidden and parcel redirection remains blocked', ['reason : Deny: C4 requested action must be parcel collection; C5 requested purpose must be pickup only; C9 billing details must stay hidden.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
