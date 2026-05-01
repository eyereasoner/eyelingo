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
C1 OK - the Python checker loaded the normalized 8-Queens input
C2 OK - the printed solution gives one column for each row
C3 OK - all printed columns are within the board
C4 OK - the printed solution uses each column exactly once
C5 OK - no pair of printed queens shares a diagonal
C6 OK - the rendered board contains exactly eight rows and eight queens
C7 OK - an independent Python bit-mask search counts 92 total solutions
C8 OK - the reported total matches the independent Python solution count
