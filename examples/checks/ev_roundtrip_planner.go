package main

import (
	"strings"
)

func checkEVRoundtrip(ctx *Context) []Check {
	d := ctx.M()
	plans, depth := evSearch(d)
	best := plans[0]
	ansPlan := []string{}
	if m := reFind(ctx.Answer, `Select plan : ([^.]+)\.`); m != nil {
		for _, p := range strings.Split(m[0], "->") {
			ansPlan = append(ansPlan, strings.TrimSpace(p))
		}
	}
	return []Check{{"bounded search finds eight acceptable Brussels-to-Cologne plans", len(plans) == 8 && grabInt(ctx.Answer, `acceptable plans : (\d+)`) == 8}, {"the selected plan is the fastest acceptable candidate", sliceEq(ansPlan, best.actions) && sliceEq(best.actions, []string{"drive_bru_liege", "drive_liege_aachen", "shuttle_aachen_cologne"})}, {"selected duration and fuel remaining are recomputed", close(fieldFloat(ctx.Answer, "duration"), best.duration, 5e-7) && grabInt(ctx.Answer, `fuel remaining : (\d+) of`) == best.fuel}, {"selected cost, belief, and comfort are recomputed by summing/multiplying actions", close(fieldFloat(ctx.Answer, "cost"), best.cost, 5e-4) && close(fieldFloat(ctx.Answer, "belief"), best.belief, 5e-7) && close(fieldFloat(ctx.Answer, "comfort"), best.comfort, 5e-7)}, {"the selected final state satisfies the wildcard goal", evMatches(asMap(d["Goal"]), best.state) && evKey(best.state) == "Cologne|low|none"}, {"the selected plan satisfies reliability, cost, and duration thresholds", best.belief > num(asMap(d["Thresholds"])["MinBelief"]) && best.cost < num(asMap(d["Thresholds"])["MaxCost"]) && best.duration < num(asMap(d["Thresholds"])["MaxDuration"])}, {"the last mile uses the high-belief Aachen-Cologne shuttle", best.actions[len(best.actions)-1] == "shuttle_aachen_cologne"}, {"search depth stays within the fuel-step bound", depth <= integer(d["FuelSteps"])}, {"the top two acceptable plans are ordered by duration then cost", plans[0].duration <= plans[1].duration && plans[0].cost < plans[1].cost}}
}
