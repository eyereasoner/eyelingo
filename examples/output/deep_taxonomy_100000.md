# Deep Taxonomy 100000  

## Answer  
The deep taxonomy test succeeds.  
Starting fact : :ind a :N0  
Reached class : :ind a :N100000  
Terminal class : :ind a :A2  
Success flag : :test :is true  

Proof checkpoints:  
:N0 present : yes  
:N1 plus :I1/:J1 present : yes  
:N50000 plus :I50000/:J50000 present : yes  
:N99999 and :N100000 present : yes  
:A2 and success flag present : yes  

## Reason why  
The N3 source is a very deep rule chain. Each taxonomy-step rule consumes the same individual in class Ni and derives the next class N(i+1), plus two side labels I(i+1) and J(i+1). Once N100000 is present, the terminal rule derives A2; once A2 is present, the success rule derives :test :is true.  
source N3 starting fact assertions : 1  
source N3 taxonomy-step rules : 100000  
source N3 terminal/success rules : 2  
source N3 ARC check/report rules : 7  
source N3 total rules counted : 100009  
translated representation : compressed rule schema + 1563-word bit sets  
agenda pops : 100001  
taxonomy-step rule applications : 100000  
terminal rule applications : 1  
success rule applications : 1  
classification facts derived : 100001 N classes + 200000 side labels + A2 = 300002 type facts  
The side labels are not needed for the final A2 proof, but carrying both I and J at every level checks that the whole wide/deep expansion was performed, not just the main N-chain.  

## Check  
C1 OK - the requested taxonomy depth is 100000  
C2 OK - the terminal class N100000 is reported as reached  
C3 OK - the terminal A2 class and success flag are reported  
C4 OK - taxonomy-step rule count matches the JSON depth  
C5 OK - total source-rule count recomputes from start, step, terminal, success, and report rules  
C6 OK - agenda pops match one pop per N-class fact in the main chain  
C7 OK - taxonomy-step applications match the depth  
C8 OK - terminal and success rules fire exactly once  
C9 OK - classification fact total accounts for N, I/J side labels, and A2  
C10 OK - midpoint and endpoint checkpoints are present  
