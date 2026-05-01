# Sudoku  

## Answer  
The puzzle is solved, and the completed grid is the unique valid Sudoku solution.  
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

## Reason why  
The solver starts from 23 clues and fills the remaining 58 cells by combining constraint propagation with depth-first search. At each step it chooses the empty cell with the fewest legal digits, places forced singles immediately, and only guesses when more than one candidate remains. Across the search it made 213 forced placements and tried 22 guesses, visited 23 search nodes overall, and backtracked 12 times before reaching the completed grid. The solver also confirmed that the solution is unique. Early steps: r2c3=4: guess, r5c3=3: forced, r2c1=5: guess, r2c4=1: guess, r2c6=9: forced, r2c7=6: guess, r2c8=7: forced, r1c7=4: guess, … and 50 more placements

## Check
C1 OK - the input puzzle has 81 cells and exactly 23 given clues
C2 OK - the completed grid is parsed as nine rows of nine digits
C3 OK - every original clue is preserved at the same row and column
C4 OK - the final grid contains only digits 1 through 9
C5 OK - each completed row is a permutation of 1 through 9
C6 OK - each completed column is a permutation of 1 through 9
C7 OK - each completed 3×3 box is a permutation of 1 through 9
C8 OK - every filled cell is legal against its row, column, and box peers
C9 OK - the completed grid matches a separately embedded expected solution fixture
