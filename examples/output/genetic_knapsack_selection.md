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

## Go audit details  
platform : go1.26.2 linux/amd64  
case : genetic_knapsack_selection  
question : Which feasible knapsack genome is reached by deterministic single-bit mutation and selection?  
capacity : 50  
items : 12  
max generations : 64  
generation 0 : parent=000000000000 fitness=1000000 weight=0 value=0 best=000000000100 bestFitness=999972 bestWeight=14 bestValue=28  
generation 1 : parent=000000000100 fitness=999972 weight=14 value=28 best=000000000101 bestFitness=999946 bestWeight=27 bestValue=54  
generation 2 : parent=000000000101 fitness=999946 weight=27 value=54 best=100000000101 bestFitness=999922 bestWeight=39 bestValue=78  
generation 3 : parent=100000000101 fitness=999922 weight=39 value=78 best=101000000101 bestFitness=999899 bestWeight=50 bestValue=101  
generation 4 : parent=101000000101 fitness=999899 weight=50 value=101 best=101000000101 bestFitness=999899 bestWeight=50 bestValue=101  
checks passed : 5/5  
