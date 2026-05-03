package main

func checkDrone(ctx *Context) []Check {
	d := ctx.M()
	plans := droneSearch(d)
	best := plans[0]
	reported := parseDroneAnswer(ctx.Answer)
	next := plans[1]
	exp := asMap(d["Expected"])
	unique := setOf([]string{droneStateKey(asMap(d["Start"])), droneStateKey(best.end), "Brugge|mid|none"})
	return []Check{{"14 corridor actions are loaded from JSON", len(maps(d["Actions"])) == 14}, {"bounded search independently finds 17 surviving plans", len(plans) == integer(exp["SurvivingPlans"]) && reported.fuel == len(plans)}, {"the selected path is the lowest-cost survivor", sliceEq(reported.actions, best.actions) && sliceEq(best.actions, []string{"fly_gent_brugge", "public_coastline_brugge_oostende"})}, {"the selected path starts with the expected first action", best.actions[0] == str(exp["SelectedFirstAction"])}, {"duration, cost, belief, and comfort are recomputed along the selected path", reported.duration == best.duration && close(reported.cost, best.cost, 5e-4) && close(reported.belief, best.belief, 5e-7) && close(reported.comfort, best.comfort, 5e-7)}, {"the selected end state reaches Oostende with low battery and no permit", droneStateKey(best.end) == "Oostende|low|none"}, {"the selected belief and cost satisfy the thresholds", best.belief > num(asMap(d["Thresholds"])["MinBelief"]) && best.cost < num(asMap(d["Thresholds"])["MaxCost"])}, {"state-cycle pruning keeps every selected-path state unique", len(best.actions)+1 == len(unique)}, {"the next cheapest survivor costs 0.014 as stated", close(next.cost, 0.014, 1e-12) && contains(ctx.Reason, "next cheapest")}}
}
