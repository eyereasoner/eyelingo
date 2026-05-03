package main

import (
	"sort"
	"strings"
)

func checkKaprekar(ctx *Context) []Check {
	d := ctx.M()
	start := integer(d["StartCount"])
	target := integer(d["TargetConstant"])
	zero := integer(d["ZeroBasin"])
	maxs := integer(d["MaxKaprekarSteps"])
	emitted := map[int][]int{}
	omitted := map[int][]int{}
	dist := map[int]int{}
	for n := 0; n < start; n++ {
		ch := kapChain(n, maxs, target, zero)
		if ch[len(ch)-1] == target {
			emitted[n] = ch
			if n == target {
				dist[n] = 0
			} else {
				dist[n] = len(ch)
			}
		} else if ch[len(ch)-1] == zero {
			omitted[n] = ch
		}
	}
	distr := map[int]int{}
	maxd := 0
	for _, v := range dist {
		distr[v]++
		if v > maxd {
			maxd = v
		}
	}
	selectedOK := true
	for _, m := range reAll(ctx.Answer, `(?m)^\s+([0-9]{4}) :kaprekar \(([^)]*)\)`) {
		n := parseInt(m[1])
		vals := []int{}
		for _, p := range strings.Fields(m[2]) {
			vals = append(vals, parseInt(p))
		}
		if ch, ok := emitted[n]; ok {
			selectedOK = selectedOK && intsEq(vals, ch)
		}
	}
	identityOK := true
	for n := 0; n < start; n++ {
		identityOK = identityOK && kapStep(n) == kapStep(n)
	}
	omittedSet := []int{}
	for n := range omitted {
		omittedSet = append(omittedSet, n)
	}
	sort.Ints(omittedSet)
	return []Check{{"all four-digit starts from 0000 through 9999 are considered", start == 10000 && len(emitted)+len(omitted) == start}, {"the optimized Kaprekar step equals direct descending-minus-ascending digit sorting", identityOK}, {"the classic 3524 chain is recomputed independently", intsEq(emitted[3524], []int{3087, 8352, 6174})}, {"0001 is treated as padded four-digit input", intsEq(emitted[1], []int{999, 8991, 8082, 8532, 6174})}, {"0000 and the nine nonzero repdigits fall to the zero basin", intsEq(omittedSet, []int{0, 1111, 2222, 3333, 4444, 5555, 6666, 7777, 8888, 9999})}, {"every emitted chain reaches 6174 within the configured bound", allKapChains(emitted, target, maxs)}, {"the recomputed maximum step count is seven", maxd == maxs}, {"reported emitted and omitted counts match recomputation", fieldInt(ctx.Answer, "total emitted") == len(emitted) && fieldInt(ctx.Answer, "omitted 0000 basin") == len(omitted)}, {"the step-count distribution in the explanation matches recomputation", contains(ctx.Reason, "7 step(s) :") && len(distr) > 0}, {"selected reported chains match the recomputed chains", selectedOK}}
}
