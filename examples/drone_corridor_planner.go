// drone_corridor_planner.go
//
// A compact Go translation inspired by Eyeling's
// `examples/drone-corridor-planner.n3`.
//
// The example composes corridor actions into bounded plans from Gent to
// Oostende. Duration and cost are summed; belief and comfort are multiplied.
// Plans are kept only when they satisfy the same pruning style as the N3 file.
//
// Run:
//
//	go run examples/drone_corridor_planner.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "drone_corridor_planner"

type Dataset struct {
	CaseName     string
	Question     string
	Start        State
	GoalLocation string
	Fuel         int
	Thresholds   Thresholds
	Actions      []Action
	Expected     Expected
}

type State struct {
	Location string
	Battery  string
	Permit   string
}

type Action struct {
	Name        string
	From        State
	To          State
	DurationSec int
	Cost        float64
	Belief      float64
	Comfort     float64
}

type Thresholds struct {
	MinBelief float64
	MaxCost   float64
}

type Expected struct {
	SelectedFirstAction string
	SurvivingPlans      int
}

type Plan struct {
	Actions     []string
	End         State
	DurationSec int
	Cost        float64
	Belief      float64
	Comfort     float64
	FuelLeft    int
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Plans  []Plan
	Best   Plan
	Checks []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
	printChecks(analysis)
	printAudit(ds, analysis)
	if !allChecksOK(analysis.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Analysis {
	plans := make([]Plan, 0)
	var walk func(State, int, []string, int, float64, float64, float64, map[string]bool)
	walk = func(st State, fuel int, acts []string, dur int, cost float64, belief float64, comfort float64, seen map[string]bool) {
		if st.Location == ds.GoalLocation {
			if belief > ds.Thresholds.MinBelief && cost < ds.Thresholds.MaxCost {
				plans = append(plans, Plan{Actions: append([]string(nil), acts...), End: st, DurationSec: dur, Cost: cost, Belief: belief, Comfort: comfort, FuelLeft: fuel})
			}
			return
		}
		if fuel == 0 {
			return
		}
		for _, action := range ds.Actions {
			if !matches(action.From, st) {
				continue
			}
			next := apply(action.To, st)
			key := stateKey(next)
			if seen[key] {
				continue
			}
			nextSeen := copySeen(seen)
			nextSeen[key] = true
			walk(next, fuel-1, append(acts, action.Name), dur+action.DurationSec, cost+action.Cost, belief*action.Belief, comfort*action.Comfort, nextSeen)
		}
	}
	walk(ds.Start, ds.Fuel, nil, 0, 0, 1, 1, map[string]bool{stateKey(ds.Start): true})

	sort.Slice(plans, func(i, j int) bool {
		if !almostEqual(plans[i].Cost, plans[j].Cost) {
			return plans[i].Cost < plans[j].Cost
		}
		if plans[i].DurationSec != plans[j].DurationSec {
			return plans[i].DurationSec < plans[j].DurationSec
		}
		if !almostEqual(plans[i].Belief, plans[j].Belief) {
			return plans[i].Belief > plans[j].Belief
		}
		return strings.Join(plans[i].Actions, ",") < strings.Join(plans[j].Actions, ",")
	})

	best := Plan{}
	if len(plans) > 0 {
		best = plans[0]
	}
	checks := []Check{
		{ID: "C1", OK: len(ds.Actions) >= 10, Text: fmt.Sprintf("%d corridor actions were loaded from JSON", len(ds.Actions))},
		{ID: "C2", OK: len(plans) == ds.Expected.SurvivingPlans, Text: fmt.Sprintf("bounded search found %d plans meeting belief and cost thresholds", len(plans))},
		{ID: "C3", OK: len(plans) > 0 && best.Actions[0] == ds.Expected.SelectedFirstAction, Text: fmt.Sprintf("lowest-cost selected plan starts with %s", firstAction(best))},
		{ID: "C4", OK: len(plans) > 0 && best.End.Location == ds.GoalLocation, Text: fmt.Sprintf("selected plan reaches %s", ds.GoalLocation)},
		{ID: "C5", OK: len(plans) > 0 && best.Belief > ds.Thresholds.MinBelief, Text: fmt.Sprintf("selected belief %.6f is above %.2f", best.Belief, ds.Thresholds.MinBelief)},
		{ID: "C6", OK: len(plans) > 0 && best.Cost < ds.Thresholds.MaxCost, Text: fmt.Sprintf("selected cost %.3f is below %.2f", best.Cost, ds.Thresholds.MaxCost)},
	}
	return Analysis{Plans: plans, Best: best, Checks: checks}
}

func matches(pattern State, actual State) bool {
	return matchesField(pattern.Location, actual.Location) && matchesField(pattern.Battery, actual.Battery) && matchesField(pattern.Permit, actual.Permit)
}

func matchesField(pattern string, value string) bool {
	return pattern == "*" || pattern == value
}

func apply(pattern State, actual State) State {
	return State{Location: choose(pattern.Location, actual.Location), Battery: choose(pattern.Battery, actual.Battery), Permit: choose(pattern.Permit, actual.Permit)}
}

func choose(pattern string, value string) string {
	if pattern == "*" {
		return value
	}
	return pattern
}

func copySeen(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in)+1)
	for k, v := range in {
		out[k] = v
	}
	return out
}

func stateKey(st State) string {
	return st.Location + "|" + st.Battery + "|" + st.Permit
}

func almostEqual(a, b float64) bool {
	if a > b {
		return a-b < 1e-12
	}
	return b-a < 1e-12
}

func firstAction(plan Plan) string {
	if len(plan.Actions) == 0 {
		return ""
	}
	return plan.Actions[0]
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func countChecksOK(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# Drone Corridor Planner")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("selected plan : %s\n", strings.Join(analysis.Best.Actions, " -> "))
	fmt.Printf("duration : %d s\n", analysis.Best.DurationSec)
	fmt.Printf("cost : %.3f\n", analysis.Best.Cost)
	fmt.Printf("belief : %.6f\n", analysis.Best.Belief)
	fmt.Printf("comfort : %.6f\n", analysis.Best.Comfort)
	fmt.Printf("end state : location=%s battery=%s permit=%s fuelLeft=%d\n", analysis.Best.End.Location, analysis.Best.End.Battery, analysis.Best.End.Permit, analysis.Best.FuelLeft)
	fmt.Printf("surviving plans : %d\n", len(analysis.Plans))
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason why")
	fmt.Println("The planner treats each corridor description as a state transition over location, battery, and permit state.")
	fmt.Printf("Search is fuel-bounded to %d steps, which keeps cycles finite while still allowing charging and permit actions.\n", ds.Fuel)
	fmt.Println("For composed plans, duration and cost are summed while belief and comfort are multiplied along the path.")
	fmt.Printf("Plans are retained only when belief is greater than %.2f and cost is less than %.2f.\n", ds.Thresholds.MinBelief, ds.Thresholds.MaxCost)
	if len(analysis.Plans) > 1 {
		fmt.Printf("The selected plan is the lowest-cost surviving plan; the next cheapest starts with %s and costs %.3f.\n", firstAction(analysis.Plans[1]), analysis.Plans[1].Cost)
	}
	fmt.Println()
}

func printChecks(analysis Analysis) {
	fmt.Println("## Check")
	for _, check := range analysis.Checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", check.ID, status, check.Text)
	}
	fmt.Println()
}

func printAudit(ds Dataset, analysis Analysis) {
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("start : location=%s battery=%s permit=%s\n", ds.Start.Location, ds.Start.Battery, ds.Start.Permit)
	fmt.Printf("goal location : %s\n", ds.GoalLocation)
	fmt.Printf("fuel tokens : %d\n", ds.Fuel)
	fmt.Printf("actions loaded : %d\n", len(ds.Actions))
	fmt.Printf("thresholds : minBelief=%.2f maxCost=%.2f\n", ds.Thresholds.MinBelief, ds.Thresholds.MaxCost)
	fmt.Printf("plans passing thresholds : %d\n", len(analysis.Plans))
	limit := len(analysis.Plans)
	if limit > 5 {
		limit = 5
	}
	for i := 0; i < limit; i++ {
		plan := analysis.Plans[i]
		fmt.Printf("rank %d : %s duration=%d cost=%.3f belief=%.6f comfort=%.6f endBattery=%s permit=%s fuelLeft=%d\n", i+1, strings.Join(plan.Actions, " -> "), plan.DurationSec, plan.Cost, plan.Belief, plan.Comfort, plan.End.Battery, plan.End.Permit, plan.FuelLeft)
	}
	fmt.Printf("checks passed : %d/%d\n", countChecksOK(analysis.Checks), len(analysis.Checks))
}
