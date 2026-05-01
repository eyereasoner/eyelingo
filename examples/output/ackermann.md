# Ackermann  

## Answer  
The ackermann.n3 test query derives 12 Ackermann facts.  
Computed values:  
A0 ackermann(0,0) = 1  
A1 ackermann(0,6) = 7  
A2 ackermann(1,2) = 4  
A3 ackermann(1,7) = 9  
A4 ackermann(2,2) = 7  
A5 ackermann(2,9) = 21  
A6 ackermann(3,4) = 125  
A7 ackermann(3,1000) = 302-digit integer [857206885749013856758740...220995094697645344555005; sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e]  
A8 ackermann(4,0) = 13  
A9 ackermann(4,1) = 65533  
A10 ackermann(4,2) = 19729-digit integer [200352993040684646497907...339445587895905719156733; sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b]  
A11 ackermann(5,0) = 65533  
Large exact-value fingerprints:  
A7 digits=302 leading=857206885749013856758740 trailing=220995094697645344555005 sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e  
A10 digits=19729 leading=200352993040684646497907 trailing=339445587895905719156733 sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b  

## Reason why  
The N3 source defines binary ackermann(x,y) by computing T(x,y+3,2) and subtracting 3. The ternary predicate T uses direct rules for successor, addition, multiplication, and exponentiation, then uses the recursive hyperoperation rule T(x,y,z)=T(x-1,T(x,y-1,z),z) when x>3 and y is non-zero.  
primitive test queries : 12  
binary reductions : 12  
distinct ternary facts : 23  
memo hits : 5  
rule paths:  
A0 binary offset -> T(0,3,2) -> successor gives T=4, answer=T-3=1  
A1 binary offset -> T(0,9,2) -> successor gives T=10, answer=T-3=7  
A2 binary offset -> T(1,5,2) -> addition gives T=7, answer=T-3=4  
A3 binary offset -> T(1,10,2) -> addition gives T=12, answer=T-3=9  
A4 binary offset -> T(2,5,2) -> multiplication gives T=10, answer=T-3=7  
A5 binary offset -> T(2,12,2) -> multiplication gives T=24, answer=T-3=21  
A6 binary offset -> T(3,7,2) -> exponentiation gives T=128, answer=T-3=125  
A7 binary offset -> T(3,1003,2) -> exponentiation gives T=302-digit integer [857206885749013856758740...220995094697645344555008; sha256=00800470f0cc96fd412c9f2f9cb3e7d9c6d3a244820bd5c4b65928c8d9f899cf], answer=T-3=302-digit integer [857206885749013856758740...220995094697645344555005; sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e]  
A8 binary offset -> T(4,3,2) -> tetration recursion gives T=16, answer=T-3=13  
A9 binary offset -> T(4,4,2) -> tetration recursion gives T=65536, answer=T-3=65533  
A10 binary offset -> T(4,5,2) -> tetration recursion gives T=19729-digit integer [200352993040684646497907...339445587895905719156736; sha256=64829919027d6b545f931768c171c25c3022620a646a692116639aecb45f0ba8], answer=T-3=19729-digit integer [200352993040684646497907...339445587895905719156733; sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b]  
A11 binary offset -> T(5,3,2) -> higher hyperoperation recursion gives T=65536, answer=T-3=65533  
hyperoperation highlights:  
A7 is 2^1003 - 3, an exact 302-digit integer.  
A10 is 2^65536 - 3, an exact 19,729-digit integer summarized by fingerprint.  
A11 reuses the pentation step T(5,3,2)=T(4,4,2)=65536, so A11 equals A9.  

## Check  
C1 OK - all twelve JSON Ackermann queries are recomputed  
C2 OK - x=0 queries reduce to successor after the +3 binary offset  
C3 OK - x=1 queries reduce to addition after the +3 binary offset  
C4 OK - x=2 queries reduce to multiplication after the +3 binary offset  
C5 OK - x=3 queries reduce to exact base-2 exponentiation  
C6 OK - A(4,0) and A(4,1) match the first tetration cases  
C7 OK - A(4,2) is held exactly as 2^65536-3 with the reported fingerprint  
C8 OK - A(5,0) lands on the same value as A(4,1)  
C9 OK - all non-huge reported exact values match recomputation  
C10 OK - the reported proof statistics match the query and memo structure  
