# Dijkstra Risk Path  

## Answer  
selected path : ClinicA -> DepotB -> LabD -> HubZ  
raw cost : 10.00  
risk sum : 0.55  
risk-adjusted score : 11.10  
edges in selected path : 3  

## Reason  
Each edge contributes its delivery cost plus the configured risk penalty.  
Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph.  
The DepotC shortcut has lower early cost but carries enough risk to lose under the configured risk weight.  
The selected route balances cost and risk through DepotB and LabD.  

## Check  
C1 OK - the graph fixture has the requested start, goal, risk weight, and eight directed edges  
C2 OK - every edge weight is independently computed as cost + riskWeight × risk  
C3 OK - the reported path is made only of directed edges present in the input graph  
C4 OK - reported raw cost, risk sum, score, and edge count match the parsed path  
C5 OK - an independent Dijkstra-style search selects ClinicA -> DepotB -> LabD -> HubZ  
C6 OK - the independent shortest-path score matches the fixture expectation  
C7 OK - exhaustive simple-path enumeration finds no lower risk-adjusted score  
C8 OK - the best route through DepotC is independently more expensive than the selected route  
