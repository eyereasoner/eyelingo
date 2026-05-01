# Goldbach 1000  

## Answer  
All 499 even integers from 4 through 1000 have a Goldbach witness.  
sample witnesses : 4=2+2; 28=5+23; 100=3+97; 998=7+991; 1000=3+997  

## Reason why  
The checker caches primes up to the configured bound and then searches each even number E for a prime P not greater than E/2 where E-P is also prime.  
No counterexample is found in the bounded range, so the bounded Goldbach condition succeeds for this dataset.  

## Check  
C1 OK - the configured upper bound is 1000  
C2 OK - there are 499 even integers from 4 through 1000  
C3 OK - every checked even integer has a prime-pair witness  
C4 OK - each requested sample even has a witness  
C5 OK - there are 168 primes at or below 1000  
