# Dijkstra Risk Path  

## Answer  
selected path : ClinicA -> DepotB -> LabD -> HubZ  
raw cost : 10.00  
risk sum : 0.55  
risk-adjusted score : 11.10  
edges in selected path : 3  

## Reason why  
Each edge contributes its delivery cost plus the configured risk penalty.  
Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph.  
The DepotC shortcut has lower early cost but carries enough risk to lose under the configured risk weight.  
The selected route balances cost and risk through DepotB and LabD.  

## Check  
C1 OK - all route edges were loaded from JSON  
C2 OK - edge score is cost + 2.00 × risk  
C3 OK - Dijkstra reached HubZ from ClinicA  
C4 OK - selected path is ClinicA -> DepotB -> LabD -> HubZ  
C5 OK - selected total score is 11.10  
C6 OK - higher-risk shortcut through DepotC is rejected  
