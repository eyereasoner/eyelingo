package main

func checkDining(ctx *Context) []Check {
	sched := ctx.A()
	traces, counts, final, reqTotal, transTotal := diningDerive(sched)
	reported := parseDiningRows(ctx.Answer)
	mealTotal := 0
	noShared := true
	patternOK := len(reported) == len(traces)
	for i, t := range traces {
		mealTotal += len(t.meals)
		noShared = noShared && len(t.forksUsed) == len(setOf(t.forksUsed))
		if i < len(reported) {
			patternOK = patternOK && sliceEq(reported[i], t.meals)
		}
	}
	all3 := true
	for _, p := range philos {
		all3 = all3 && counts[p] == 3
	}
	numsOK := len(sched) == 9
	for i, r := range sched {
		numsOK = numsOK && integer(asMap(r)["Number"]) == i+1
	}
	finalH := map[string]string{}
	finalDirty := true
	for f, st := range final {
		finalH[f] = st.holder
		finalDirty = finalDirty && st.clean == "Dirty"
	}
	wantFinal := map[string]string{"F12": "P2", "F23": "P2", "F34": "P4", "F45": "P5", "F51": "P5"}
	return []Check{{"the JSON schedule contains nine deterministic rounds", numsOK}, {"the Chandy-Misra trace yields exactly 15 meals", mealTotal == 15 && contains(ctx.Answer, "meals : 15")}, {"each philosopher eats exactly three times", all3 && contains(ctx.Answer, "everyone ate 3 times : yes")}, {"no two meals in the same round use the same fork", noShared}, {"dirty-fork transfer simulation reproduces the reported meal pattern", patternOK}, {"the first round transfers only F23 to P3 and lets P1/P3 eat", len(traces[0].transfers) == 1 && traces[0].transfers[0] == [3]string{"P2", "P3", "F23"} && sliceEq(traces[0].meals, []string{"P1", "P3"})}, {"the P5-only rounds derive one meal each", sliceEq(traces[2].meals, []string{"P5"}) && sliceEq(traces[5].meals, []string{"P5"}) && sliceEq(traces[8].meals, []string{"P5"})}, {"all forks end dirty after the final phase", finalDirty}, {"the final fork holders match the independent simulation", mapStrEq(finalH, wantFinal)}, {"request and transfer counts are nontrivial and internally consistent", reqTotal == 26 && transTotal == 26}}
}
