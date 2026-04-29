# Ackermann

## Answer
The ackermann.n3 test query derives 12 Ackermann facts.
- Computed values:
 - A0 ackermann(0,0) = 1
 - A1 ackermann(0,6) = 7
 - A2 ackermann(1,2) = 4
 - A3 ackermann(1,7) = 9
 - A4 ackermann(2,2) = 7
 - A5 ackermann(2,9) = 21
 - A6 ackermann(3,4) = 125
 - A7 ackermann(3,1000) = 302-digit integer [857206885749013856758740...220995094697645344555005; sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e]
 - A8 ackermann(4,0) = 13
 - A9 ackermann(4,1) = 65533
 - A10 ackermann(4,2) = 19729-digit integer [200352993040684646497907...339445587895905719156733; sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b]
 - A11 ackermann(5,0) = 65533
- Large exact-value fingerprints:
 - A7 digits=302 leading=857206885749013856758740 trailing=220995094697645344555005 sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e
 - A10 digits=19729 leading=200352993040684646497907 trailing=339445587895905719156733 sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b

## Reason why
The N3 source defines binary ackermann(x,y) by computing T(x,y+3,2) and subtracting 3. The ternary predicate T uses direct rules for successor, addition, multiplication, and exponentiation, then uses the recursive hyperoperation rule T(x,y,z)=T(x-1,T(x,y-1,z),z) when x>3 and y is non-zero.
- primitive test queries : 12
- binary reductions : 12
- distinct ternary facts : 23
- memo hits : 5
- rule paths:
 - A0 binary offset -> T(0,3,2) -> successor gives T=4, answer=T-3=1
 - A1 binary offset -> T(0,9,2) -> successor gives T=10, answer=T-3=7
 - A2 binary offset -> T(1,5,2) -> addition gives T=7, answer=T-3=4
 - A3 binary offset -> T(1,10,2) -> addition gives T=12, answer=T-3=9
 - A4 binary offset -> T(2,5,2) -> multiplication gives T=10, answer=T-3=7
 - A5 binary offset -> T(2,12,2) -> multiplication gives T=24, answer=T-3=21
 - A6 binary offset -> T(3,7,2) -> exponentiation gives T=128, answer=T-3=125
 - A7 binary offset -> T(3,1003,2) -> exponentiation gives T=302-digit integer [857206885749013856758740...220995094697645344555008; sha256=00800470f0cc96fd412c9f2f9cb3e7d9c6d3a244820bd5c4b65928c8d9f899cf], answer=T-3=302-digit integer [857206885749013856758740...220995094697645344555005; sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e]
 - A8 binary offset -> T(4,3,2) -> tetration recursion gives T=16, answer=T-3=13
 - A9 binary offset -> T(4,4,2) -> tetration recursion gives T=65536, answer=T-3=65533
 - A10 binary offset -> T(4,5,2) -> tetration recursion gives T=19729-digit integer [200352993040684646497907...339445587895905719156736; sha256=64829919027d6b545f931768c171c25c3022620a646a692116639aecb45f0ba8], answer=T-3=19729-digit integer [200352993040684646497907...339445587895905719156733; sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b]
 - A11 binary offset -> T(5,3,2) -> higher hyperoperation recursion gives T=65536, answer=T-3=65533
- hyperoperation highlights:
 - A7 is 2^1003 - 3, an exact 302-digit integer.
 - A10 is 2^65536 - 3, an exact 19,729-digit integer summarized by fingerprint.
 - A11 reuses the pentation step T(5,3,2)=T(4,4,2)=65536, so A11 equals A9.

## Check
- C1 OK - x=0 reduces to successor after the y+3 binary offset.
- C2 OK - x=1 reduces to addition after the y+3 binary offset.
- C3 OK - x=2 reduces to multiplication after the y+3 binary offset.
- C4 OK - x=3 reduces to exact BigInt exponentiation, including 2^1003-3.
- C5 OK - x=4 derives the first tetration cases T(4,3,2)-3 and T(4,4,2)-3.
- C6 OK - A(4,2) is held exactly as 2^65536-3, not as a floating-point approximation.
- C7 OK - the pentation query A(5,0) lands on the same value as A(4,1).
- C8 OK - the evaluator reached the expected largest exact integer and memoized each distinct ternary fact once.

## Go audit details
- platform : go1.26.2 linux/amd64
- question : Evaluate the Ackermann facts queried by ackermann.n3.
- translated source : ackermann.n3
- binary definition : ackermann(x,y) = T(x,y+3,2)-3
- primitive test queries : 12
- computed ternary facts : 23
- calls including memo hits : 28
- memo hits : 5
- successor rules : 2
- addition rules : 2
- multiplication rules : 2
- power rules : 7
- one/base rules : 2
- recursive hyperoperation rules : 8
- max x reached : 5
- max y decimal digits : 5
- max result decimal digits : 19729
- N3 test expressions:
 - A0 (0 0) :ackermann ?A0
 - A1 (0 6) :ackermann ?A1
 - A2 (1 2) :ackermann ?A2
 - A3 (1 7) :ackermann ?A3
 - A4 (2 2) :ackermann ?A4
 - A5 (2 9) :ackermann ?A5
 - A6 (3 4) :ackermann ?A6
 - A7 (3 1000) :ackermann ?A7
 - A8 (4 0) :ackermann ?A8
 - A9 (4 1) :ackermann ?A9
 - A10 (4 2) :ackermann ?A10
 - A11 (5 0) :ackermann ?A11
- derived fact fingerprints:
 - A0 ackermann(0,0) digits=1 sha256=6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b
 - A1 ackermann(0,6) digits=1 sha256=7902699be42c8a8e46fbbb4501726517e86b22c56a189f7625a6da49081b2451
 - A2 ackermann(1,2) digits=1 sha256=4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a
 - A3 ackermann(1,7) digits=1 sha256=19581e27de7ced00ff1ce50b2047e7a567c76b1cbaebabe5ef03f7c3017bb5b7
 - A4 ackermann(2,2) digits=1 sha256=7902699be42c8a8e46fbbb4501726517e86b22c56a189f7625a6da49081b2451
 - A5 ackermann(2,9) digits=2 sha256=6f4b6612125fb3a0daecd2799dfd6c9c299424fd920f9b308110a2c1fbd8f443
 - A6 ackermann(3,4) digits=3 sha256=0f8ef3377b30fc47f96b48247f463a726a802f62f3faa03d56403751d2f66c67
 - A7 ackermann(3,1000) digits=302 sha256=365b71740f457d3277649054c99e1e3632b6590da448939739a259cf1152377e
 - A8 ackermann(4,0) digits=2 sha256=3fdba35f04dc8c462986c992bcf875546257113072a909c162f7e470e581e278
 - A9 ackermann(4,1) digits=5 sha256=786c4b42f717ce4387fa6a85e7d908d84b64fc623c98bd6d14ec0b5717d6837e
 - A10 ackermann(4,2) digits=19729 sha256=daeafdfae592d817542850bc06b0ae8c808ab7727edd8a71f0c5d588c7b7083b
 - A11 ackermann(5,0) digits=5 sha256=786c4b42f717ce4387fa6a85e7d908d84b64fc623c98bd6d14ec0b5717d6837e
- checks passed : 8/8
- all checks pass : yes
