# Fundamental Theorem Arithmetic

## Answer
- Primary N3 case: n = 202692987 has prime factors 3 * 3 * 7 * 829 * 3881.
- primary prime-power form : 3^2 * 7 * 829 * 3881
- sample count : 6
- largest sample : 600851475143
- total prime factors counted with multiplicity : 31
- distinct primes seen across samples : 17

- Sample factorizations:
  360360 = 2^3 * 3^2 * 5 * 7 * 11 * 13
  202692987 = 3^2 * 7 * 829 * 3881
  4294967295 = 3 * 5 * 17 * 257 * 65537
  600851475143 = 71 * 839 * 1471 * 6857
  9876543210 = 2 * 3^2 * 5 * 17^2 * 379721
  9999999967 = 9999999967

## Reason why
Existence comes from repeated smallest-divisor decomposition.
At each step, the first divisor found is prime because no smaller
positive divisor can divide the current number.

- Smallest-divisor trace for the N3 source number:
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

## Check
C1 OK - the source example factors 202692987 as 3,3,7,829,3881
C2 OK - the source example groups repeated factors as 3^2 * 7 * 829 * 3881
C3 OK - multiplying each computed factor list reconstructs its original number
C4 OK - every distinct factor in every decomposition is prime by trial division
C5 OK - smallest-first and largest-first traversals sort to the same multisets
C6 OK - the extended sample includes six cases and includes the ten-digit prime 9999999967

## Go audit details
- platform : go1.26.2 linux/amd64
- source file : fundamental-theorem-arithmetic.n3
- question : What are the prime factorizations, and are they unique up to order?
- primary n : 202692987
- primary smallest-first factors : 3,3,7,829,3881
- primary largest-first factors : 3881,829,7,3,3
- primary flat factor string : 3 * 3 * 7 * 829 * 3881
- primary prime-power string : 3^2 * 7 * 829 * 3881
- expected flat string : 3 * 3 * 7 * 829 * 3881
- expected largest-first string : 3881 * 829 * 7 * 3 * 3
- sample numbers : 360360,202692987,4294967295,600851475143,9876543210,9999999967
- sample count : 6
- largest sample : 600851475143
- total prime factors counted with multiplicity : 31
- distinct primes seen across samples : 17
- smallest-divisor searches : 31
- divisibility tests : 104580
- primality checks : 25
- prime-divisor tests : 50580
- checks passed : 6/6
- all checks pass : yes
