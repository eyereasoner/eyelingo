# Path Discovery

## Answer
The path discovery query finds 3 air routes with at most 2 stopovers.
- from : Ostend-Bruges International Airport
- to : Václav Havel Airport Prague
- max stopovers : 2

- Discovered routes:
 - route 1 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Diagoras Airport -> Václav Havel Airport Prague
 - route 2 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Heraklion International Nikos Kazantzakis Airport -> Václav Havel Airport Prague
 - route 3 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Palma De Mallorca Airport -> Václav Havel Airport Prague

## Reason why
- The N3 source defines a recursive :route relation over nepo:hasOutboundRouteTo facts. A route can use a direct edge when the current length is within the maximum, or extend through a non-visited intermediate airport and recurse with length+1. The final log:collectAllIn query collects the labels of each airport in every route from the source to the destination.
- source N3 airport labels : 7698
- source N3 outbound-route facts : 37505
- translated full airport labels : 7698
- translated full outbound-route facts : 37505
- airport terms appearing in outbound facts : 3334
- frontier airports expanded : 9
- bounded search outbound facts touched : 338
- source outbound candidates : 1
- Liège outbound candidates : 7
- direct routes : 0
- one-stop routes : 0
- two-stopover routes : 3
- search recursive calls : 333
- search edge tests : 338
- search depth-limit leaves : 321
- Second-hop candidates from Liège:
 - Ajaccio-Napoléon Bonaparte Airport (res:AIRPORT_1324)
 - Al Massira Airport (res:AIRPORT_1064)
 - Bastia-Poretta Airport (res:AIRPORT_1321)
 - Diagoras Airport (res:AIRPORT_1472)
 - Heraklion International Nikos Kazantzakis Airport (res:AIRPORT_1452)
 - Lille-Lesquin Airport (res:AIRPORT_1399)
 - Palma De Mallorca Airport (res:AIRPORT_3998)

## Check
- C1 OK - source and destination airport labels are known.
- C2 OK - Ostend-Bruges has one outbound route in the full N3 graph, to Liège Airport.
- C3 OK - the discovered route set matches the N3 query answer.
- C4 OK - no direct or one-stop route exists under the same bound.
- C5 OK - every discovered route has at most two stopovers.
- C6 OK - every hop is backed by a nepo:hasOutboundRouteTo fact.
- C7 OK - no route revisits an airport.
- C8 OK - the Go translation loaded every airport label and outbound-route fact from the N3 source.
- C9 OK - route output is sorted deterministically by airport labels.

## Go audit details
- platform : go1.26.2 linux/amd64
- question : Find routes from Ostend-Bruges International Airport to Václav Havel Airport Prague with at most 2 stopovers.
- source airport : Ostend-Bruges International Airport (res:AIRPORT_310)
- destination airport : Václav Havel Airport Prague (res:AIRPORT_1587)
- source graph airport labels : 7698
- source graph outbound facts : 37505
- translated full airport labels : 7698
- translated full outbound-route facts : 37505
- airport terms appearing in outbound facts : 3334
- bounded search outbound facts touched : 338
- max stopovers : 2
- max hops : 3
- routes discovered : 3
- mandatory first hop : Liège Airport (res:AIRPORT_309)
- expanded airports:
 - Ostend-Bruges International Airport (res:AIRPORT_310)
 - Liège Airport (res:AIRPORT_309)
 - Ajaccio-Napoléon Bonaparte Airport (res:AIRPORT_1324)
 - Al Massira Airport (res:AIRPORT_1064)
 - Bastia-Poretta Airport (res:AIRPORT_1321)
 - Diagoras Airport (res:AIRPORT_1472)
 - Heraklion International Nikos Kazantzakis Airport (res:AIRPORT_1452)
 - Lille-Lesquin Airport (res:AIRPORT_1399)
 - Palma De Mallorca Airport (res:AIRPORT_3998)
- route 1 terms : res:AIRPORT_310 -> res:AIRPORT_309 -> res:AIRPORT_1472 -> res:AIRPORT_1587
- route 1 labels : Ostend-Bruges International Airport -> Liège Airport -> Diagoras Airport -> Václav Havel Airport Prague
- route 1 hops : 3
- route 1 stopovers : 2
- route 2 terms : res:AIRPORT_310 -> res:AIRPORT_309 -> res:AIRPORT_1452 -> res:AIRPORT_1587
- route 2 labels : Ostend-Bruges International Airport -> Liège Airport -> Heraklion International Nikos Kazantzakis Airport -> Václav Havel Airport Prague
- route 2 hops : 3
- route 2 stopovers : 2
- route 3 terms : res:AIRPORT_310 -> res:AIRPORT_309 -> res:AIRPORT_3998 -> res:AIRPORT_1587
- route 3 labels : Ostend-Bruges International Airport -> Liège Airport -> Palma De Mallorca Airport -> Václav Havel Airport Prague
- route 3 hops : 3
- route 3 stopovers : 2
- search recursive calls : 333
- search edge tests : 338
- search edges extended : 332
- search revisit prunes : 6
- search depth-limit leaves : 321
- search dead ends : 0
- search routes emitted : 3
- search max depth : 3
- checks passed : 9/9
- all checks pass : yes
