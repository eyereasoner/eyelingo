# Go Capture Scenario (Weiqi)

## Answer
The move W at (1,0) is legal and captures 1 opponent stone(s).
- Captured stones at: (1,1)
- Board after the move:

. W . . . 
W . W . . 
. W . . . 
. . . . . 
. . . . . 


## Reason why
In Go, a stone lives if it belongs to a group with at least one liberty (empty adjacent intersection).
When a player places a stone, any opponent group that loses its last liberty is immediately captured and removed.
Before the move, the Black stone at (1,1) had only one liberty at (1,0) – the other three were occupied by White.
By playing White at (1,0), that last liberty is taken, so the Black stone is captured and removed.
The White move itself is safe because the newly placed stone still has liberties (e.g., (0,0), (2,0), ...).

## Check
- C1 OK - Black stone had exactly 1 liberty before the move (count=1).
- C2 OK - After capture, the Black stone is gone (0 liberties).
- C3 OK - The White move is not suicide (it has remaining liberties).
- C4 OK - At least one stone was captured.
- C5 OK - Board size is 5x5.

## Go audit details
- platform : go1.26.2 linux/amd64
- board size : 5x5
- move : W (1,0)
- initial board:
. W . . . 
. B W . . 
. W . . . 
. . . . . 
. . . . . 

- result : legal move: captures 1 opponent stone(s)
- captured stones : 1
- checks passed : 5/5
- recommendation consistent : yes
