# Genetic Knapsack Selection  

## Insight  
final genome : 101000000101  
selected items : item01, item03, item10, item12  
weight : 50 / 50  
value : 101  
fitness : 999899  
generations evaluated : 5  
exhaustive optimum value : 104 at genome 001000011111  

## Explanation  
Each genome bit says whether the corresponding item is selected for the knapsack.  
Feasible candidates get fitness 1000000 minus value, so higher value means lower fitness; overweight candidates are penalized above every feasible candidate.  
At every generation, all single-bit mutants of the parent are compared with the parent, and the lowest-fitness candidate is selected deterministically.  
The run stops at 101000000101 because every one-bit neighbor is no better under the capacity 50 rule.  
For transparency, an exhaustive check also finds the global best feasible value 104; this example demonstrates a local mutation search, not a promise of global optimality.  
