package main

import (
	"strings"
)

func checkFTA(ctx *Context) []Check {
	nums := []int{}
	for _, v := range ctx.A() {
		nums = append(nums, integer(v))
	}
	facts := map[int][]int{}
	for _, n := range nums {
		facts[n] = ftaFactor(n)
	}
	rows := map[int]string{}
	for _, m := range reAll(ctx.Answer, `(?m)^\s+([0-9]+)\s+=\s+(.+)$`) {
		rows[parseInt(m[1])] = strings.TrimSpace(m[2])
	}
	distinct := map[int]bool{}
	total := 0
	prodOK := true
	for n, fs := range facts {
		prod := 1
		for _, p := range fs {
			prod *= p
			distinct[p] = true
		}
		prodOK = prodOK && prod == n
		total += len(fs)
	}
	allPrime := true
	for p := range distinct {
		allPrime = allPrime && ftaPrime(p)
	}
	rowsOK := len(rows) == len(nums)
	for n, row := range rows {
		rowsOK = rowsOK && row == ftaFormat(facts[n])
	}
	return []Check{{"the primary source number is factored independently", intsEq(facts[202692987], []int{3, 3, 7, 829, 3881})}, {"the reported primary prime-power form matches grouped exponents", contains(ctx.Answer, "primary prime-power form : "+ftaFormat(facts[202692987]))}, {"multiplying every factor list reconstructs its source integer", prodOK}, {"every distinct factor found by trial division is prime", allPrime}, {"the report includes one sample factorization row for every JSON number", rowsOK}, {"every reported sample row matches the independently formatted factorization", rowsOK}, {"smallest-first and largest-first traversals sort to the same primary multiset", contains(ctx.Reason, "source sorted comparison : 3 * 3 * 7 * 829 * 3881")}, {"reported sample count and largest sample match JSON", fieldInt(ctx.Answer, "sample count") == len(nums) && fieldInt(ctx.Answer, "largest sample") == maxInt(nums)}, {"reported total factor multiplicity and distinct-prime count match recomputation", fieldInt(ctx.Answer, "total prime factors counted with multiplicity") == total && fieldInt(ctx.Answer, "distinct primes seen across samples") == len(distinct)}, {"the ten-digit prime remains unfactored after trial division", intsEq(facts[9999999967], []int{9999999967}) && ftaPrime(9999999967)}}
}
