# Independent Python checks for the odrl_dpv_risk_ranked example.
from .common import run_fragment_checks

CHECKS = [
    ('4 risk rows were derived', ['Risk: portability is restricted because exporting user data is prohibited. Clause C4: Users are not permitted to export their data.']),
    ('ranked output is in descending score order', ['Rows are sorted by descending score so the highest-risk clauses are reviewed first.']),
    ('score range is 70 to 100', ['score=100 (risk:HighRisk, risk:HighSeverity) clause C1']),
    ('high=3 moderate=1 low=0 risk levels were derived', ['Risk: account/data removal is permitted without notice safeguards (no notice constraint and no duty to inform). Clause C1: Provider may remove the user account (and associated data']),
    ('5 mitigation measures were generated', ['Each triggered rule derives a risk row with a normalized score, a source clause, and one or more mitigation measures.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
