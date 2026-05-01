# 8-Queens  

## Answer  
Solving 8-Queens...  
Printing at most 1 solution(s).  

Solution 1:  
Q . . . . . . .  
. . . . Q . . .  
. . . . . . . Q  
. . . . . Q . .  
. . Q . . . . .  
. . . . . . Q .  
. Q . . . . . .  
. . . Q . . . .  
As column positions by row: [1, 5, 8, 6, 3, 7, 2, 4]  

Total solutions for 8-Queens: 92  

## Reason why  
The solver places one queen per row on a 8x8 board.  
At each row it uses bit masks for occupied columns and both diagonal directions to enumerate only safe candidate columns.  
Counting continues after the printed solution limit, so the total solution count remains complete.  

## Check  
C1 OK - search reached depth 8  
C2 OK - first solution places one queen in each row  
C3 OK - first solution columns are unique  
C4 OK - no pair of queens in the first solution shares a diagonal  
C5 OK - counted 92 solutions for the normalized 8-Queens input  
