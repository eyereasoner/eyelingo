# Genetic Knapsack Selection  

## Answer  
final genome : 101000000101  
selected items : item01, item03, item10, item12  
weight : 50 / 50  
value : 101  
fitness : 999899  
generations evaluated : 5  

## Reason why  
Each genome bit says whether the corresponding item is selected for the knapsack.  
Feasible candidates get fitness 1000000 minus value, so higher value means lower fitness; overweight candidates are penalized above every feasible candidate.  
At every generation, all single-bit mutants of the parent are compared with the parent, and the lowest-fitness candidate is selected deterministically.  
The run stops at 101000000101 because every one-bit neighbor is no better under the capacity 50 rule.  

## Check  
C1 OK - 12 items align with a 12-bit genome  
C2 OK - final weight 50 is within capacity 50  
C3 OK - final genome is 101000000101  
C4 OK - final value is 101 using item01, item03, item10, item12  
C5 OK - no single-bit neighbor improves the final candidate  
