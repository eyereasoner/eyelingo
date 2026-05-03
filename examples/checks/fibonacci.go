package main

import (
	"math/big"
	"strings"
)

func checkFibonacci(ctx *Context) []Check {
	raw := asMap(ctx.Load())
	expected := map[int]*big.Int{}
	largest := 0
	for k, v := range raw {
		n := parseInt(k)
		z := new(big.Int)
		z.SetString(str(v), 10)
		expected[n] = z
		if n > largest {
			largest = n
		}
	}
	computed := map[int]*big.Int{}
	for n := range expected {
		computed[n] = fibBig(n)
	}
	idx, val := 0, ""
	if m := reFind(ctx.Answer, `index\s+(\d+)\s+is:\s*\n\s*([0-9]+)`); m != nil {
		idx = parseInt(m[0])
		val = m[1]
	}
	recur := true
	for n := 2; n < 60; n++ {
		s := new(big.Int).Add(fibBig(n-1), fibBig(n-2))
		recur = recur && s.Cmp(fibBig(n)) == 0
	}
	mono := true
	for n := 2; n < 200; n++ {
		mono = mono && fibBig(n+1).Cmp(fibBig(n)) >= 0
	}
	cass := new(big.Int).Sub(new(big.Int).Mul(fibBig(1001), fibBig(999)), new(big.Int).Mul(fibBig(1000), fibBig(1000)))
	allMatch := true
	nonneg := true
	for n, v := range computed {
		allMatch = allMatch && v.Cmp(expected[n]) == 0
		nonneg = nonneg && v.Sign() >= 0
	}
	f1000 := fibBig(1000).String()
	f10000 := fibBig(10000).String()
	return []Check{{"base cases are recomputed as F(0)=0 and F(1)=1", computed[0].Sign() == 0 && computed[1].Cmp(big.NewInt(1)) == 0}, {"the recurrence F(n)=F(n-1)+F(n-2) holds over an independent prefix", recur}, {"every JSON expected Fibonacci value matches fast-doubling recomputation", allMatch}, {"the report answers the largest requested index", idx == largest}, {"the reported F(10000) exactly matches the independent big-integer value", val == computed[largest].String()}, {"F(1000) has the expected exact decimal length and endpoints", len(f1000) == 209 && strings.HasPrefix(f1000, "434665576869") && strings.HasSuffix(f1000, "76137795166849228875")}, {"F(10000) has the expected exact decimal length and endpoints", len(f10000) == 2090 && strings.HasPrefix(f10000, "336447648764") && strings.HasSuffix(f10000, "6073310059947366875")}, {"all requested Fibonacci numbers are nonnegative", nonneg}, {"the Fibonacci sequence is nondecreasing from F(2) onward", mono}, {"Cassini's identity holds at n=1000 for the same independent generator", cass.Cmp(big.NewInt(1)) == 0}, {"the explanation names arbitrary-precision arithmetic for the large integer", contains(ctx.Reason, "Arbitrary") && contains(ctx.Reason, "precision") && contains(ctx.Reason, "without overflow")}}
}
