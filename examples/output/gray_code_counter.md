# Gray Code Counter  

## Answer  
bits : 4  
states visited : 16  
unique states : 16  
sequence prefix : 0000, 0001, 0011, 0010, 0110, 0111, 0101, 0100  
wrap transition : 1000 -> 0000  
maximum adjacent Hamming distance : 1  

## Reason why  
The counter maps each integer n to n xor (n >> 1), which is the reflected binary Gray-code construction.  
For 4 bits, the first 16 integers cover the full state space without duplicates.  
The Hamming-distance check compares each state with the next state, including the final wraparound transition.  
A valid cyclic Gray counter therefore changes exactly one bit at every step.  

## Check  
C1 OK - 16 states were generated for a 4-bit counter  
C2 OK - all generated states are unique  
C3 OK - each adjacent transition flips exactly one bit  
C4 OK - the final state wraps to the first with one bit flip  
C5 OK - first four states match the reflected binary Gray-code prefix  
C6 OK - the numeric generator is n xor (n >> 1)  

## Go audit details  
platform : go1.26.2 linux/amd64  
case : gray_code_counter  
question : Does a 4-bit Gray counter visit all states while flipping exactly one bit at each step?  
bits : 4  
requested steps : 16  
unique states : 16  
adjacent transitions checked : 16  
wrap hamming distance : 1  
checks passed : 6/6  
