# Gray Code Counter  

## Insight  
bits : 4  
states visited : 16  
unique states : 16  
sequence prefix : 0000, 0001, 0011, 0010, 0110, 0111, 0101, 0100  
wrap transition : 1000 -> 0000  
maximum adjacent Hamming distance : 1  

## Explanation  
The counter maps each integer n to n xor (n >> 1), which is the reflected binary Gray-code construction.  
For 4 bits, the first 16 integers cover the full state space without duplicates.  
The Hamming-distance check compares each state with the next state, including the final wraparound transition.  
A valid cyclic Gray counter therefore changes exactly one bit at every step.  
