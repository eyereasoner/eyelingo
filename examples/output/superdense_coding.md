# Superdense Coding  

## Answer  
Superdense-coding facts that survive GF(2) parity cancellation:  
  0 dqc:superdense-coding 0  
  1 dqc:superdense-coding 1  
  2 dqc:superdense-coding 2  
  3 dqc:superdense-coding 3  

## Reason  
Alice and Bob start with |R) = |0,0) + |1,1).  
Alice chooses one relation for the first mobit; Bob applies one joint test.  
The N3 example keeps only answers with odd derivation count, so duplicate  
modal paths cancel just like addition over GF(2).  

Alice operations:  
  message 0 -> ID -> encoded support {(0,0), (1,1)}  
  message 1 -> G  -> encoded support {(0,1), (1,0)}  
  message 2 -> K  -> encoded support {(0,0), (0,1), (1,1)}  
  message 3 -> KG -> encoded support {(0,0), (0,1), (1,0)}  

Raw candidate counts before parity filtering:  
  encoded 0 -> decoded counts {0:1, 1:2, 2:0, 3:2}  
  encoded 1 -> decoded counts {0:2, 1:1, 2:2, 3:0}  
  encoded 2 -> decoded counts {0:2, 1:2, 2:1, 3:2}  
  encoded 3 -> decoded counts {0:2, 1:2, 2:2, 3:1}  

Surviving explanations:  
  0 -> 0 because count=1 is odd; all competing counts are even.  
  1 -> 1 because count=1 is odd; all competing counts are even.  
  2 -> 2 because count=1 is odd; all competing counts are even.  
  3 -> 3 because count=1 is odd; all competing counts are even.  

## Check  
C1 OK - shared entanglement R contains exactly |0,0) and |1,1)  
C2 OK - composition KG is recomputed by composing G then K  
C3 OK - composition GK is recomputed by composing K then G  
C4 OK - the raw superdense rule creates 24 candidate derivations  
C5 OK - GF(2) cancellation leaves odd diagonal and even off-diagonal counts  
C6 OK - reported surviving decoded messages match the parity survivors  
C7 OK - the four Alice operations produce distinct encoded supports  
C8 OK - the JSON relation facts match the primitive teaching-model relations  
