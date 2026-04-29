# Go Ko Rule Scenario (Weiqi)

## Answer
- Move 1: W at (1,1) is legal move: no captures.
- Board after White's capture:
. W . 
B W W 
. W . 

- Move 2: B at (1,0) is legal.
- Board after recapture (should not happen here):
. W . 
B W W 
. W . 


## Reason why
- In Go, a ko occurs when:
1. A player captures exactly one stone.
2. After the capture, the capturing stone is left with exactly one liberty.
3. If the opponent immediately recaptures, the board would return to its
   exact state before the first capture.
The ko rule forbids such an immediate recapture to prevent infinite loops.

- In this scenario:
- The Black stone at (1,0) had only one liberty (1,1).
- White captured it by playing at (1,1), removing the Black stone.
- The White stone at (1,1) now had only one liberty (1,0).
- If Black were allowed to play at (1,0), it would capture the White stone
  and the board would look identical to how it was before White's capture.
- Therefore, Black's immediate recapture is illegal under the ko rule.

## Check
- C1 OK - Black stone had exactly 1 liberty before capture (count=3).
- C2 OK - White's capture move was legal.
- C3 FAIL - No capture occurred.
- C4 OK - After capture, White stone at (1,1) has exactly 1 liberty (count=4).
- C5 FAIL - Ko detection did not work.

## Go audit details
- platform : go1.26.2 linux/amd64
- initial board:
. W . 
B . W 
. W . 
- move 1 : W (1,1) -> legal move: no captures
- move 2 attempt: B (1,0) -> occupied by B
- ko detection : false
- checks passed : 1/5
- recommendation consistent : no
