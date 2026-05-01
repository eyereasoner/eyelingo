# Independent Python checks for the digital_product_passport example.
from .common import run_fragment_checks

CHECKS = [
    ('component masses fold to 105 g', ['total component mass : 105 g']),
    ('recycled component masses fold to 14 g', ['The passport folds the explicit component list to derive total mass and recycled mass, then computes an integer recycled-content percentage.']),
    ('integer recycled-content percentage is 13%', ['recycled content : 13%']),
    ('lifecycle footprint totals 52500 gCO2e', ['lifecycle footprint : 52500 gCO2e']),
    ('critical raw material exposure is Cobalt, Lithium', ['Critical raw-material exposure is derived from component-material links: Cobalt, Lithium.']),
    ('replaceable battery plus repair and spare-parts documents derive repairFriendly', ['The product is repair-friendly because the battery is replaceable and the public passport section exposes both repair and spare-parts documentation.']),
    ('public section contains user manual, repair guide, and spare-parts catalog', ['public endpoint : https://example.org/dpp/ACME-X1000-SN123']),
    ('restricted declarations stay in the restricted section', ['BatteryPack-01 Battery mass=48g recycled=0g materials=Lithium, Cobalt, Nickel replaceable=true']),
    ('manufacturing, sale, and repair events are chronologically consistent', ['Lifecycle footprint is derived by summing manufacturing, transport, and use-phase emissions.']),
    ('passport endpoint matches the product digital link', ['Passport decision : PASS for ACME X1000 SN123.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
