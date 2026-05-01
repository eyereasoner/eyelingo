# Independent Python checks for the dijkstra_risk_path example.
from .common import run_fragment_checks

CHECKS = [
    ('all route edges were loaded from JSON', ["Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph."]),
    ('edge score is cost + 2.00 × risk', ['Each edge contributes its delivery cost plus the configured risk penalty.']),
    ('Dijkstra reached HubZ from ClinicA', ['selected path : ClinicA -> DepotB -> LabD -> HubZ']),
    ('selected path is ClinicA -> DepotB -> LabD -> HubZ', ['The selected route balances cost and risk through DepotB and LabD.']),
    ('selected total score is 11.10', ['risk-adjusted score : 11.10']),
    ('higher-risk shortcut through DepotC is rejected', ['The DepotC shortcut has lower early cost but carries enough risk to lose under the configured risk weight.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
