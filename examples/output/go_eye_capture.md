# Go Eye Capture Scenario (Weiqi)

## Answer
The Black group had only a single eye and no outside liberties.
White plays inside the eye at (2,2) and captures the entire group.
Initial board:
W W W W W 
W B B B W 
W B . B W 
W B B B W 
W W W W W 
Move: W at (2,2)
Board after the move:
W W W W W 
W . . . W 
W . W . W 
W . . . W 
W W W W W 
Number of captured Black stones: 8

## Reason why
In Go, a group lives only if it has at least two eyes, or can make them
eventually. A group with only one eye and no other liberties is dead
because the opponent can safely play inside that eye.
When the opponent fills the last liberty, the entire group is removed
from the board as captured stones. The played stone itself remains alive
because after the capture it inherits the newly vacated liberties.

In this scenario, the Black enclosure is completely sealed by White
stones. The only empty intersection inside is the eye at (2,2).
White’s move inside the eye is not suicide because the Black stones are
immediately captured, freeing the space around the newly placed White stone.

## Check
C1 OK - Before the move, the Black group had exactly 1 liberty (the eye).
C2 OK - Initial Black group size is 8 (8 expected).
C3 OK - The killing move was legal.
C4 OK - All Black stones were captured.
C5 OK - After capture, the White move is not a suicide (it has liberties).

## Go audit details
platform : go1.26.2 linux/amd64
board size : 5x5
initial board:
W W W W W 
W B B B W 
W B . B W 
W B B B W 
W W W W W 
killing move : W (2,2)
result : legal move: captures 8 opponent stone(s)
captured stones : 8
captured positions : (1,1), (2,1), (1,2), (3,1), (1,3), (3,2), (2,3), (3,3)
checks passed : 5/5
recommendation consistent : yes
