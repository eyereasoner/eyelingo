# Sudoku  

## Insight  
The puzzle is solved, and the completed grid is a valid Sudoku solution.  
case : sudoku  
default puzzle : classic  

Puzzle  
1 . . | . . 7 | . 9 .  
. 3 . | . 2 . | . . 8  
. . 9 | 6 . . | 5 . .  

. . 5 | 3 . . | 9 . .  
. 1 . | . 8 . | . . 2  
6 . . | . . 4 | . . .  

3 . . | . . . | . 1 .  
. 4 . | . . . | . . 7  
. . 7 | . . . | 3 . .  

Completed grid  
1 6 2 | 8 5 7 | 4 9 3  
5 3 4 | 1 2 9 | 6 7 8  
7 8 9 | 6 4 3 | 5 2 1  

4 7 5 | 3 1 2 | 9 8 6  
9 1 3 | 5 8 6 | 7 4 2  
6 2 8 | 7 9 4 | 1 3 5  

3 5 6 | 4 7 8 | 2 1 9  
2 4 1 | 9 3 5 | 8 6 7  
8 9 7 | 2 6 1 | 3 5 4  

## Explanation  
The input contains 23 given clues and 58 empty cells.  
The trust gate checks that every clue is preserved, each row contains digits 1 through 9, each column contains digits 1 through 9, and each 3×3 box contains digits 1 through 9.  
Only after those constraints hold does the example emit the completed grid.  
