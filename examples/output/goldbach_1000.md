# Goldbach 1000  

## Answer  
All 499 even integers from 4 through 1000 have a Goldbach witness.  
sample witnesses : 4=2+2; 28=5+23; 100=3+97; 998=7+991; 1000=3+997  

## Reason  
The checker caches primes up to the configured bound and then searches each even number E for a prime P not greater than E/2 where E-P is also prime.  
No counterexample is found in the bounded range, so the bounded Goldbach condition succeeds for this dataset.  

## Check  
C1 OK - the configured upper bound is parsed from JSON as 1000  
C2 OK - there are exactly 499 even integers from 4 through 1000  
C3 OK - trial division independently finds 168 primes at or below 1000  
C4 OK - every checked even integer has a prime-pair witness  
C5 OK - each requested sample even has the first witness found by the independent search  
C6 OK - every reported witness uses two primes whose sum is the reported even integer  
C7 OK - the bounded Goldbach result is derived from recomputed witnesses, not from static output text  
