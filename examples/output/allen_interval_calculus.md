# Allen Interval Calculus  

## Answer  
derived relations : 110 ordered interval pairs  
showcase : A before B | A meets C | A overlaps D | F starts A | G finishes A | A during H | A equals E | J meets I | K finishes C  
completed intervals : I=16:00-18:00, K=13:30-14:00  
invalid intervals : 0  

## Reason  
The example completes any interval that has a start plus duration before comparing endpoints.  
Each ordered pair is then classified with the 13 Allen base relations, including the six converse relations.  
The relation rules are purely endpoint constraints, so the result is deterministic and traceable.  
The duration-derived intervals participate in the same relation table as directly supplied intervals.  

## Check  
C1 OK - duration-based intervals are completed from start plus minutes  
C2 OK - all completed intervals have strictly positive duration  
C3 OK - all ordered non-self interval pairs are classified  
C4 OK - every required Allen relation is recomputed from endpoints  
C5 OK - A/B, A/C, and A/D demonstrate before, meets, and overlaps  
C6 OK - converse relations are independently recovered  
C7 OK - start, finish, during, contains, and equals cases all occur  
C8 OK - the showcase text includes each required forward example  
