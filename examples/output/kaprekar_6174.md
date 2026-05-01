# Kaprekar 6174  

## Answer  
Kaprekar chains that end at 6174 are emitted as :kaprekar facts.  
total emitted : 9990  
omitted 0000 basin : 10  
maximum steps to 6174 : 7  

Selected facts, shown with four-digit padding for readability:  
  0001 :kaprekar (0999 8991 8082 8532 6174)  
  3524 :kaprekar (3087 8352 6174)  
  6174 :kaprekar (6174)  
  9831 :kaprekar (8442 5994 5355 1998 8082 8532 6174)  
  9998 :kaprekar (0999 8991 8082 8532 6174)  

## Reason why  
Each start is read as four digits, so 1 is treated as 0001.  
The digits are sorted once, then the optimized identity computes the  
same result as descending-number minus ascending-number.  
The search is bounded to seven steps, matching the N3 source: any  
four-digit start that reaches 6174 does so within that bound.  

Step-count distribution for emitted starts:  
  0 step(s) : 1 start(s)  
  1 step(s) : 383 start(s)  
  2 step(s) : 576 start(s)  
  3 step(s) : 2400 start(s)  
  4 step(s) : 1272 start(s)  
  5 step(s) : 1518 start(s)  
  6 step(s) : 1656 start(s)  
  7 step(s) : 2184 start(s)  

Examples omitted because they fall to 0000:  
  0000 -> (0000)  
  1111 -> (0000)  
  2222 -> (0000)  
  9999 -> (0000)  

## Check  
C1 OK - all four-digit starts from 0000 through 9999 are considered  
C2 OK - the optimized Kaprekar step equals direct descending-minus-ascending digit sorting  
C3 OK - the classic 3524 chain is recomputed independently  
C4 OK - 0001 is treated as padded four-digit input  
C5 OK - 0000 and the nine nonzero repdigits fall to the zero basin  
C6 OK - every emitted chain reaches 6174 within the configured bound  
C7 OK - the recomputed maximum step count is seven  
C8 OK - reported emitted and omitted counts match recomputation  
C9 OK - the step-count distribution in the explanation matches recomputation  
C10 OK - selected reported chains match the recomputed chains  
