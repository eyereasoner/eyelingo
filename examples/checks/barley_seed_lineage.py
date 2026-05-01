# Independent Python checks for the barley_seed_lineage example.
from .common import run_fragment_checks

CHECKS = [
    ('no-design laws are loaded', ['The main lineage satisfies the constructor-theory style CAN side: digital heredity under no-design laws, repair support, a protected dormant seed stage, germination resources, prop']),
    ('mainLine can copy its digitally instantiated genome', ['mainLine CAN : genome-copy, protected-dormancy, germination, propagule-production, accurate-self-reproduction, lineage-closure, adaptive-persistence']),
    ('mainLine has repair, protected dormancy, and greenhouse support', ['analogLine lacks a digital hereditary medium, fragileLine lacks repair, coatlessLine lacks the protected dormant compartment, and staticLine lacks heritable variation.']),
    ('mainLine is evolvable under the saline selection environment', ['Only mainLine closes the life cycle and adaptively persists in the saline selection environment.']),
    ('analogLine is blocked by non-digital heredity', ['blocked contrast lineages : analogLine, fragileLine, coatlessLine, staticLine']),
    ('fragileLine is blocked by missing repair', ["The contrast lineages are deliberately near misses so the CAN'T side is explicit."]),
    ('coatlessLine is blocked by missing dormancy protection', ["mainLine CAN'T : none of the modeled blockers apply"]),
    ('staticLine is blocked by missing heritable variation', ['evolvable lineage : mainLine']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
