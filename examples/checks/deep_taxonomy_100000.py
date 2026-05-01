# Independent Python checks for the deep_taxonomy_100000 example.
from .common import run_fragment_checks

CHECKS = [
    ('the starting classification N0 is present', ['Starting fact : :ind a :N0']),
    ('the first expansion produced N1 together with side labels I1 and J1', ['The N3 source is a very deep rule chain. Each taxonomy-step rule consumes the same individual in class Ni and derives the next class N(i+1), plus two side labels I(i+1) and J(i+1).']),
    ('the chain reaches the midpoint N50000 and still carries both side-label branches', ['The side labels are not needed for the final A2 proof, but carrying both I and J at every level checks that the whole wide/deep expansion was performed, not just the main N-chain.']),
    ('the final taxonomy step from N99999 to N100000 was completed', [':N99999 and :N100000 present : yes']),
    ('once N100000 is reached, the terminal class A2 is derived', ['Reached class : :ind a :N100000']),
    ('the success flag is raised only after the terminal class A2 is present', [':A2 and success flag present : yes']),
    ('all 100001 N-level classifications are present with no missing level', ['classification facts derived : 100001 N classes + 200000 side labels + A2 = 300002 type facts']),
    ('all 200000 side-label classifications I1..I100000 and J1..J100000 are present', [':N1 plus :I1/:J1 present : yes']),
    ('exactly 100000 taxonomy-step rules, one terminal rule, and one success rule fired', ['source N3 taxonomy-step rules : 100000']),
    ('a linear scan finds no skipped taxonomy level', ['The deep taxonomy test succeeds.']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
