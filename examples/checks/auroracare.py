# Independent Python checks for the auroracare example.
from .common import run_fragment_checks

CHECKS = [
    ('all seven scenarios match the PERMIT/DENY outcomes encoded in the N3 example', ['deny count : 3']),
    ('primary-care access requires a clinician role and a care-team link', ["reason : Permitted: a clinician in the patient's care team matched the primary-care policy."]),
    ('quality improvement is allowed only when both lab results and patient summary are requested in the secure environment', ['reason : Permitted: the quality-improvement policy matched the secure environment and required data categories.']),
    ('insurance management is denied by a matching prohibition', ['D – Insurance management : DENY (urn:policy:deny-insurance)']),
    ('research is allowed only with patient opt-in and anonymisation in the secure environment', ['reason : Permitted: the subject opted in, the dataset is anonymised, and the research policy matched.']),
    ('AI training is denied because the subject opted out', ['reason : Denied: the subject opted out of data use for AI training.']),
    ('four scenarios are permitted and three are denied', ['For each AuroraCare scenario, should the PDP permit or deny the requested use of health data, and why?']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
