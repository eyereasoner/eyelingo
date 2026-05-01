# Dining Philosophers  

## Answer  
The Chandy-Misra dining-philosophers trace completes without conflict.  
philosophers : 5  
forks : 5  
rounds : 9  
meals : 15  
everyone ate 3 times : yes  

Meal trace:  
  round 1 cycle 1 : P1#1 uses F51,F12; P3#1 uses F23,F34  
  round 2 cycle 1 : P2#1 uses F12,F23; P4#1 uses F34,F45  
  round 3 cycle 1 : P5#1 uses F45,F51  
  round 4 cycle 2 : P1#2 uses F51,F12; P3#2 uses F23,F34  
  round 5 cycle 2 : P2#2 uses F12,F23; P4#2 uses F34,F45  
  round 6 cycle 2 : P5#2 uses F45,F51  
  round 7 cycle 3 : P1#3 uses F51,F12; P3#3 uses F23,F34  
  round 8 cycle 3 : P2#3 uses F12,F23; P4#3 uses F34,F45  
  round 9 cycle 3 : P5#3 uses F45,F51  

## Reason why  
Each round has three phases. Hungry philosophers first request adjacent  
forks they do not hold. A holder sends a requested fork only when it is  
Dirty, and the receiver gets that fork Clean. After the meal phase, all  
forks are marked Dirty for the next round.  
The Go code uses goroutines inside each phase, then applies state changes  
between phases so the output remains deterministic.  

Round derivation:  
  round 1 hungry : P1,P3  
    requests : P3 asks P2 for F23  
    transfers : P2 sends F23 to P3  
    kept forks : F12,F34,F45,F51  
    meals : P1#1,P3#1  
  round 2 hungry : P2,P4  
    requests : P2 asks P1 for F12; P2 asks P3 for F23; P4 asks P3 for F34  
    transfers : P1 sends F12 to P2; P3 sends F23 to P2; P3 sends F34 to P4  
    kept forks : F45,F51  
    meals : P2#1,P4#1  
  round 3 hungry : P5  
    requests : P5 asks P4 for F45; P5 asks P1 for F51  
    transfers : P4 sends F45 to P5; P1 sends F51 to P5  
    kept forks : F12,F23,F34  
    meals : P5#1  
  round 4 hungry : P1,P3  
    requests : P1 asks P5 for F51; P1 asks P2 for F12; P3 asks P2 for F23; P3 asks P4 for F34  
    transfers : P5 sends F51 to P1; P2 sends F12 to P1; P2 sends F23 to P3; P4 sends F34 to P3  
    kept forks : F45  
    meals : P1#2,P3#2  
  round 5 hungry : P2,P4  
    requests : P2 asks P1 for F12; P2 asks P3 for F23; P4 asks P3 for F34; P4 asks P5 for F45  
    transfers : P1 sends F12 to P2; P3 sends F23 to P2; P3 sends F34 to P4; P5 sends F45 to P4  
    kept forks : F51  
    meals : P2#2,P4#2  
  round 6 hungry : P5  
    requests : P5 asks P4 for F45; P5 asks P1 for F51  
    transfers : P4 sends F45 to P5; P1 sends F51 to P5  
    kept forks : F12,F23,F34  
    meals : P5#2  
  round 7 hungry : P1,P3  
    requests : P1 asks P5 for F51; P1 asks P2 for F12; P3 asks P2 for F23; P3 asks P4 for F34  
    transfers : P5 sends F51 to P1; P2 sends F12 to P1; P2 sends F23 to P3; P4 sends F34 to P3  
    kept forks : F45  
    meals : P1#3,P3#3  
  round 8 hungry : P2,P4  
    requests : P2 asks P1 for F12; P2 asks P3 for F23; P4 asks P3 for F34; P4 asks P5 for F45  
    transfers : P1 sends F12 to P2; P3 sends F23 to P2; P3 sends F34 to P4; P5 sends F45 to P4  
    kept forks : F51  
    meals : P2#3,P4#3  
  round 9 hungry : P5  
    requests : P5 asks P4 for F45; P5 asks P1 for F51  
    transfers : P4 sends F45 to P5; P1 sends F51 to P5  
    kept forks : F12,F23,F34  
    meals : P5#3

## Check  
C1 OK - the JSON schedule contains nine deterministic rounds  
C2 OK - the Chandy-Misra trace yields exactly 15 meals  
C3 OK - each philosopher eats exactly three times  
C4 OK - no two meals in the same round use the same fork  
C5 OK - dirty-fork transfer simulation reproduces the reported meal pattern  
C6 OK - the first round transfers only F23 to P3 and lets P1/P3 eat  
C7 OK - the P5-only rounds derive one meal each  
C8 OK - all forks end dirty after the final phase  
C9 OK - the final fork holders match the independent simulation  
C10 OK - request and transfer counts are nontrivial and internally consistent
