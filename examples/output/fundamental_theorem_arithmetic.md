# Fundamental Theorem Arithmetic  

## Insight  
Primary N3 case: n = 202692987 has prime factors 3 * 3 * 7 * 829 * 3881.  
primary prime-power form : 3^2 * 7 * 829 * 3881  
sample count : 6  
largest sample : 600851475143  
total prime factors counted with multiplicity : 31  
distinct primes seen across samples : 17  

Sample factorizations:  
  360360 = 2^3 * 3^2 * 5 * 7 * 11 * 13  
  202692987 = 3^2 * 7 * 829 * 3881  
  4294967295 = 3 * 5 * 17 * 257 * 65537  
  600851475143 = 71 * 839 * 1471 * 6857  
  9876543210 = 2 * 3^2 * 5 * 17^2 * 379721  
  9999999967 = 9999999967  

## Explanation  
Existence comes from repeated smallest-divisor decomposition.  
At each step, the first divisor found is prime because no smaller  
positive divisor can divide the current number.  

Smallest-divisor trace for the N3 source number:  
  202692987 = 3 * 67564329  
  67564329 = 3 * 22521443  
  22521443 = 7 * 3217349  
  3217349 = 829 * 3881  
  3881 is prime  

Uniqueness up to order is checked by reversing each traversal and sorting  
both factor lists. Matching sorted lists describe the same multiset of  
prime factors, even when the factors were discovered in the opposite order.  
  source smallest-first factors : 3 * 3 * 7 * 829 * 3881  
  source largest-first factors : 3881 * 829 * 7 * 3 * 3  
  source sorted comparison : 3 * 3 * 7 * 829 * 3881  

The additional samples cover repeated small factors, special products,  
large composites, and a larger prime that has no smaller divisor.  
