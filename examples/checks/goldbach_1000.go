package main

func checkGoldbach(ctx *Context) []Check {
	d := ctx.M()
	maxEven := integer(d["maxEven"])
	primes := map[int]bool{}
	for n := 2; n <= maxEven; n++ {
		if primeInt(n) {
			primes[n] = true
		}
	}
	evens := []int{}
	witnesses := map[int][2]int{}
	failures := []int{}
	for e := 4; e <= maxEven; e += 2 {
		evens = append(evens, e)
		found := false
		for p := 2; p <= e/2; p++ {
			if primes[p] && primes[e-p] {
				witnesses[e] = [2]int{p, e - p}
				found = true
				break
			}
		}
		if !found {
			failures = append(failures, e)
		}
	}
	parsed := map[int][2]int{}
	if m := reFind(ctx.Answer, `sample witnesses\s*:\s*(.+)`); m != nil {
		for _, w := range reAll(m[0], `(\d+)=(\d+)\+(\d+)`) {
			parsed[parseInt(w[1])] = [2]int{parseInt(w[2]), parseInt(w[3])}
		}
	}
	samplesOK := len(parsed) == len(asSlice(d["sampleEvens"]))
	for _, v := range asSlice(d["sampleEvens"]) {
		e := integer(v)
		samplesOK = samplesOK && parsed[e] == witnesses[e]
	}
	reportedOK := true
	for e, pq := range parsed {
		reportedOK = reportedOK && primes[pq[0]] && primes[pq[1]] && pq[0]+pq[1] == e
	}
	return []Check{{"the configured upper bound is parsed from JSON as 1000", maxEven == 1000 && grabInt(ctx.Answer, `All \d+ even integers from 4 through (\d+) have`) == maxEven}, {"there are exactly 499 even integers from 4 through 1000", len(evens) == 499 && grabInt(ctx.Answer, `All (\d+) even integers`) == len(evens)}, {"trial division independently finds 168 primes at or below 1000", len(primes) == 168}, {"every checked even integer has a prime-pair witness", len(failures) == 0}, {"each requested sample even has the first witness found by the independent search", samplesOK}, {"every reported witness uses two primes whose sum is the reported even integer", reportedOK}, {"the bounded Goldbach result is derived from recomputed witnesses, not from static output text", contains(ctx.Reason, "No counterexample") && len(witnesses) == len(evens)}}
}
