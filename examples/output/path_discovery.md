# Path Discovery  

## Answer  
The path discovery query finds 3 air routes with at most 2 stopovers.  
from : Ostend-Bruges International Airport  
to : Václav Havel Airport Prague  
max stopovers : 2  

Discovered routes:  
route 1 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Diagoras Airport -> Václav Havel Airport Prague  
route 2 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Heraklion International Nikos Kazantzakis Airport -> Václav Havel Airport Prague  
route 3 (2 stopovers): Ostend-Bruges International Airport -> Liège Airport -> Palma De Mallorca Airport -> Václav Havel Airport Prague  

## Reason why  
The N3 source defines a recursive :route relation over nepo:hasOutboundRouteTo facts. A route can use a direct edge when the current length is within the maximum, or extend through a non-visited intermediate airport and recurse with length+1. The final log:collectAllIn query collects the labels of each airport in every route from the source to the destination.  
source N3 airport labels : 7698  
source N3 outbound-route facts : 37505  
translated full airport labels : 7698  
translated full outbound-route facts : 37505  
airport terms appearing in outbound facts : 3334  
frontier airports expanded : 9  
bounded search outbound facts touched : 338  
source outbound candidates : 1  
Liège outbound candidates : 7  
direct routes : 0  
one-stop routes : 0  
two-stopover routes : 3  
search recursive calls : 333  
search edge tests : 338  
search depth-limit leaves : 321  
Second-hop candidates from Liège:  
Ajaccio-Napoléon Bonaparte Airport (res:AIRPORT_1324)  
Al Massira Airport (res:AIRPORT_1064)  
Bastia-Poretta Airport (res:AIRPORT_1321)  
Diagoras Airport (res:AIRPORT_1472)  
Heraklion International Nikos Kazantzakis Airport (res:AIRPORT_1452)  
Lille-Lesquin Airport (res:AIRPORT_1399)  
Palma De Mallorca Airport (res:AIRPORT_3998)  

## Check  
C1 OK - source and destination airport labels are known  
C2 OK - Ostend-Bruges has one outbound route in the full graph, to Liège Airport  
C3 OK - bounded DFS independently finds exactly three two-stopover routes  
C4 OK - reported route labels match the independently discovered route set  
C5 OK - no direct or one-stop route exists under the same bound  
C6 OK - every discovered hop is backed by an outbound-route fact  
C7 OK - no discovered route revisits an airport  
C8 OK - the translated graph size matches the full source counts  
C9 OK - the second-hop candidates from Liège are independently recovered  
C10 OK - route output is sorted deterministically by airport labels  
