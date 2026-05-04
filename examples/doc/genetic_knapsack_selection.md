# Genetic Knapsack Selection

## What this example is about

This example uses a deterministic mutation search for a knapsack problem. A knapsack problem asks which items to select when every item has weight and value, and the bag has a maximum capacity.

The example borrows vocabulary from genetic algorithms, but it is intentionally small and reproducible: no randomness is used.

## How it works, in plain language

A genome is a string of bits. A `1` means the corresponding item is selected; a `0` means it is not selected. At each generation, the program creates every one-bit mutation of the current genome and compares those candidates with the parent.

Feasible candidates are better when they have higher value. Overweight candidates receive a penalty so they always lose to feasible candidates. The search stops when no one-bit change improves the result.

## What to notice in the output

The final genome fills the capacity exactly: weight 50 out of 50, with value 101. The word “fitness” may look backwards because this example treats lower fitness numbers as better; feasible candidates use `1000000 - value`.

## What the trust gate checks

The trust gate verifies the final genome, selected items, weight, value, generation count, capacity rule, and local-optimum condition: every one-bit neighbor is no better than the final genome.

## Run it

From the repository root:

```sh
node examples/genetic_knapsack_selection.js
```

## Files

- [JavaScript example](../genetic_knapsack_selection.js)
- [Input data](../input/genetic_knapsack_selection.json)
- [Expected output](../output/genetic_knapsack_selection.md)
