# Genetic Knapsack Selection  

## Answer  
final genome : 101000000101  
selected items : item01, item03, item10, item12  
weight : 50 / 50  
value : 101  
fitness : 999899  
generations evaluated : 5  

## Reason  
Each genome bit says whether the corresponding item is selected for the knapsack.  
Feasible candidates get fitness 1000000 minus value, so higher value means lower fitness; overweight candidates are penalized above every feasible candidate.  
At every generation, all single-bit mutants of the parent are compared with the parent, and the lowest-fitness candidate is selected deterministically.  
The run stops at 101000000101 because every one-bit neighbor is no better under the capacity 50 rule.  

## Check  
C1 OK - the input has one item per genome bit and a valid binary start genome  
C2 OK - the checker independently simulates the deterministic single-bit local search to the reported final genome  
C3 OK - reported selected items are exactly the one-bits in the final genome  
C4 OK - reported weight, value, and fitness match independent genome evaluation  
C5 OK - the final candidate is feasible and matches the expected fixture totals  
C6 OK - no one-bit neighbor has a lower fitness than the final candidate  
C7 OK - the reported generation count matches the independent simulation history length  
C8 OK - the fixture has many feasible genomes, so the check validates the local-search rule rather than a text fragment  
