# Independent Python checks for the path_discovery example.
from .common import run_fragment_checks

CHECKS = [
    ('source and destination airport labels are known', ['The N3 source defines a recursive :route relation over nepo:hasOutboundRouteTo facts. A route can use a direct edge when the current length is within the maximum, or extend through']),
    ('Ostend-Bruges has one outbound route in the full N3 graph, to Liège Airport', ['route 3 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Palma De Mallorca Airport -> Václav Havel Airport Prague']),
    ('the discovered route set matches the N3 query answer', ['The path discovery query finds 3 air routes with at most 2 stopovers.']),
    ('no direct or one-stop route exists under the same bound', ['route 2 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Heraklion International Nikos Kazantzakis Airport -> Václav Havel Airport Prague']),
    ('every discovered route has at most two stopovers', ['route 1 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Diagoras Airport -> Václav Havel Airport Prague']),
    ('every hop is backed by a nepo:hasOutboundRouteTo fact', ['source N3 outbound-route facts : 37505']),
    ('no route revisits an airport', ['Heraklion International Nikos Kazantzakis Airport (res:AIRPORT_1452)']),
    ('the Go translation loaded every airport label and outbound-route fact from the N3 source', ['source N3 airport labels : 7698']),
    ('route output is sorted deterministically by airport labels', ['translated full airport labels : 7698']),
]

def run(ctx):
    return run_fragment_checks(ctx, CHECKS)
