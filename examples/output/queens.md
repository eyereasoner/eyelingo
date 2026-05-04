# 8-Queens  

## Insight  
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

## Explanation  
The solver places one queen per row on a 8x8 board.  
At each row it uses bit masks for occupied columns and both diagonal directions to enumerate only safe candidate columns.  
Counting continues after the printed solution limit, so the total solution count remains complete.  
