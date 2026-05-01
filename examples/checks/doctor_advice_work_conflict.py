# Independent Python checks for the doctor_advice_work_conflict example.
from .common import run_fragment_checks

CHECKS = [
    ('Jos is classified as Sick from condition Flu', ['Jos has Flu, so the health context classifies the agent as Sick.']),
    ('home programming request keeps both Permit and Deny before resolution', ['Since Home is permitted and Office is denied, the combined recommendation is RemoteWorkOnly.']),
    ('home programming conflict resolves to Permit', ['The conflict resolver permits sick ProgrammingWork at Home but denies Office work because of colleague-infection risk.']),
    ('office programming request keeps both Permit and Deny before resolution', ["The doctor's statement permits ProgrammingWork, but the general sick-work policy also denies work, so the raw closure deliberately keeps the conflict visible."]),
    ('office programming conflict resolves to Deny', ['Request_Jos_Prog_Home : raw=Deny+Permit status=BothPermitDeny effective=Permit']),
    ('overall work decision is RemoteWorkOnly', ['overall decision for Jos : RemoteWorkOnly']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
