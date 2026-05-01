# Allen Interval Calculus  

## Answer  
derived relations : 110 ordered interval pairs  
showcase : A before B | A meets C | A overlaps D | F starts A | G finishes A | A during H | A equals E | J meets I | K finishes C  
completed intervals : I=16:00-18:00, K=13:30-14:00  
invalid intervals : 0  

## Reason why  
The example completes any interval that has a start plus duration before comparing endpoints.  
Each ordered pair is then classified with the 13 Allen base relations, including the six converse relations.  
The relation rules are purely endpoint constraints, so the result is deterministic and auditable.  
The duration-derived intervals participate in the same relation table as directly supplied intervals.  

## Check  
C1 OK - 11 intervals were loaded and completed when duration was present  
C2 OK - A before B was derived  
C3 OK - A meets C and C metBy A were derived  
C4 OK - A overlaps D and D overlappedBy A were derived  
C5 OK - F starts A and G finishes A were derived  
C6 OK - A during H and H contains A were derived  
C7 OK - A equals E was derived  
C8 OK - duration completion produced I ending at 18:00 and K ending at 14:00  
C9 OK - no invalid intervals were detected  
