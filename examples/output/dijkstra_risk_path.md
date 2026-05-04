# Dijkstra Risk Path  

## Insight  
selected path : ClinicA -> DepotB -> LabD -> HubZ  
raw cost : 10.00  
risk sum : 0.55  
risk-adjusted score : 11.10  
edges in selected path : 3  

## Explanation  
Each edge contributes its delivery cost plus the configured risk penalty.  
Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph.  
The DepotC shortcut has lower early cost but carries enough risk to lose under the configured risk weight.  
The selected route balances cost and risk through DepotB and LabD.  
