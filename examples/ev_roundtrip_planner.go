// ev_roundtrip_planner.go
//
// A self-contained Go translation inspired by Eyeling's
// `examples/ev-roundtrip-planner.n3`.
//
// The scenario models a bounded, GPS-style EV journey planner from Brussels to
// Cologne. Actions are represented as state transitions with duration, cost,
// belief, and comfort scores. The planner composes transitions while consuming a
// finite search budget, then keeps only plans that satisfy the query thresholds.
//
// Run:
//
//	go run examples/ev_roundtrip_planner.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "ev_roundtrip_planner"

type Dataset struct {
	CaseName   string
	Question   string
	Vehicle    Vehicle
	Goal       Goal
	FuelSteps  int
	Thresholds Thresholds
	Actions    []Action
}

type Vehicle struct {
	ID      string
	At      string
	Battery string
	Pass    string
}

type Goal struct {
	At      string
	Battery string
	Pass    string
}

type Thresholds struct {
	MinBelief   float64
	MaxCost     float64
	MaxDuration float64
}

type Action struct {
	Name     string
	From     State
	To       State
	Duration float64
	Cost     float64
	Belief   float64
	Comfort  float64
}

type State struct {
	At      string
	Battery string
	Pass    string
}

type Plan struct {
	Actions       []string
	FinalState    State
	Duration      float64
	Cost          float64
	Belief        float64
	Comfort       float64
	FuelRemaining int
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type SearchStats struct {
	RecursiveCalls int
	ActionTests    int
	ActionsTaken   int
	StatePrunes    int
	MaxDepth       int
}

type Analysis struct {
	Start           State
	AcceptablePlans []Plan
	BestPlan        Plan
	Checks          []Check
	Stats           SearchStats
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
	start := State{At: ds.Vehicle.At, Battery: ds.Vehicle.Battery, Pass: ds.Vehicle.Pass}
	var stats SearchStats
	plans := search(ds.Actions, start, ds.Goal, ds.Thresholds, ds.FuelSteps, &stats)
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].Duration != plans[j].Duration {
			return plans[i].Duration < plans[j].Duration
		}
		if plans[i].Cost != plans[j].Cost {
			return plans[i].Cost < plans[j].Cost
		}
		return strings.Join(plans[i].Actions, "/") < strings.Join(plans[j].Actions, "/")
	})

	best := Plan{}
	if len(plans) > 0 {
		best = plans[0]
	}

	checks := []Check{
		{ID: "C1", OK: len(plans) > 0, Text: fmt.Sprintf("%d acceptable Brussels-to-Cologne plans were derived", len(plans))},
		{ID: "C2", OK: best.Duration < ds.Thresholds.MaxDuration, Text: fmt.Sprintf("selected plan duration %.1f is below %.1f", best.Duration, ds.Thresholds.MaxDuration)},
		{ID: "C3", OK: best.Cost < ds.Thresholds.MaxCost, Text: fmt.Sprintf("selected plan cost %.3f is below %.3f", best.Cost, ds.Thresholds.MaxCost)},
		{ID: "C4", OK: best.Belief > ds.Thresholds.MinBelief, Text: fmt.Sprintf("selected plan belief %.6f is above %.2f", best.Belief, ds.Thresholds.MinBelief)},
		{ID: "C5", OK: best.FinalState.At == ds.Goal.At, Text: fmt.Sprintf("selected plan reaches %s", ds.Goal.At)},
		{ID: "C6", OK: containsAction(best.Actions, "shuttle_aachen_cologne"), Text: "selected plan uses the high-belief Aachen-Cologne shuttle for the last mile"},
		{ID: "C7", OK: stats.MaxDepth <= ds.FuelSteps, Text: fmt.Sprintf("bounded search consumed at most %d of %d fuel tokens", stats.MaxDepth, ds.FuelSteps)},
	}

	return Analysis{Start: start, AcceptablePlans: plans, BestPlan: best, Checks: checks, Stats: stats}
}

func search(actions []Action, start State, goal Goal, thresholds Thresholds, fuel int, stats *SearchStats) []Plan {
	var plans []Plan
	var walk func(State, []string, float64, float64, float64, float64, int, map[State]bool)
	walk = func(state State, path []string, duration, cost, belief, comfort float64, fuelLeft int, seen map[State]bool) {
		stats.RecursiveCalls++
		if len(path) > stats.MaxDepth {
			stats.MaxDepth = len(path)
		}
		if matchesGoal(state, goal) {
			if belief > thresholds.MinBelief && cost < thresholds.MaxCost && duration < thresholds.MaxDuration {
				plans = append(plans, Plan{Actions: append([]string{}, path...), FinalState: state, Duration: duration, Cost: cost, Belief: belief, Comfort: comfort, FuelRemaining: fuelLeft})
			}
			return
		}
		if fuelLeft == 0 {
			return
		}
		for _, action := range actions {
			stats.ActionTests++
			if !matchesState(action.From, state) {
				continue
			}
			next := applyState(action.To, state)
			if seen[next] && next != state {
				stats.StatePrunes++
				continue
			}
			stats.ActionsTaken++
			nextSeen := cloneSeen(seen)
			nextSeen[next] = true
			walk(next,
				append(path, action.Name),
				duration+action.Duration,
				cost+action.Cost,
				belief*action.Belief,
				comfort*action.Comfort,
				fuelLeft-1,
				nextSeen,
			)
		}
	}
	walk(start, nil, 0, 0, 1, 1, fuel, map[State]bool{start: true})
	return plans
}

func matchesState(pattern State, state State) bool {
	return wildcardMatch(pattern.At, state.At) && wildcardMatch(pattern.Battery, state.Battery) && wildcardMatch(pattern.Pass, state.Pass)
}

func matchesGoal(state State, goal Goal) bool {
	return wildcardMatch(goal.At, state.At) && wildcardMatch(goal.Battery, state.Battery) && wildcardMatch(goal.Pass, state.Pass)
}

func wildcardMatch(pattern, value string) bool {
	return pattern == "*" || pattern == value
}

func applyState(pattern State, current State) State {
	return State{At: choose(pattern.At, current.At), Battery: choose(pattern.Battery, current.Battery), Pass: choose(pattern.Pass, current.Pass)}
}

func choose(pattern, current string) string {
	if pattern == "*" {
		return current
	}
	return pattern
}

func cloneSeen(seen map[State]bool) map[State]bool {
	copy := make(map[State]bool, len(seen)+1)
	for key, value := range seen {
		copy[key] = value
	}
	return copy
}

func containsAction(actions []string, want string) bool {
	for _, action := range actions {
		if action == want {
			return true
		}
	}
	return false
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
	best := analysis.BestPlan
	fmt.Println("# EV Roadtrip Planner")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("Select plan : %s.\n", strings.Join(best.Actions, " -> "))
	fmt.Printf("route result : %s battery=%s pass=%s\n", best.FinalState.At, best.FinalState.Battery, best.FinalState.Pass)
	fmt.Printf("duration : %.1f minutes\n", best.Duration)
	fmt.Printf("cost : %.3f\n", best.Cost)
	fmt.Printf("belief : %.6f\n", best.Belief)
	fmt.Printf("comfort : %.6f\n", best.Comfort)
	fmt.Printf("acceptable plans : %d\n", len(analysis.AcceptablePlans))
	fmt.Printf("fuel remaining : %d of %d\n", best.FuelRemaining, ds.FuelSteps)
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	best := analysis.BestPlan
	fmt.Println("## Reason why")
	fmt.Printf("The planner starts with %s at %s, battery=%s, pass=%s, then composes action descriptions until the goal city %s is reached.\n", ds.Vehicle.ID, analysis.Start.At, analysis.Start.Battery, analysis.Start.Pass, ds.Goal.At)
	fmt.Printf("Duration and cost are summed across each candidate; belief and comfort are multiplied, matching the N3 planner pattern.\n")
	fmt.Printf("The selected plan is the fastest acceptable candidate under belief > %.2f, cost < %.3f, and duration < %.1f.\n", ds.Thresholds.MinBelief, ds.Thresholds.MaxCost, ds.Thresholds.MaxDuration)
	fmt.Printf("It uses the shuttle from Aachen to Cologne, avoiding an extra charge stop while keeping belief at %.6f.\n", best.Belief)
	fmt.Println()
	fmt.Println("Top acceptable plans:")
	for i, plan := range analysis.AcceptablePlans {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s | duration=%.1f cost=%.3f belief=%.6f comfort=%.6f final=%s/%s/%s\n",
			i+1,
			strings.Join(plan.Actions, " -> "),
			plan.Duration,
			plan.Cost,
			plan.Belief,
			plan.Comfort,
			plan.FinalState.At,
			plan.FinalState.Battery,
			plan.FinalState.Pass,
		)
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
	best := analysis.BestPlan
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("vehicle : %s start=%s battery=%s pass=%s\n", ds.Vehicle.ID, ds.Vehicle.At, ds.Vehicle.Battery, ds.Vehicle.Pass)
	fmt.Printf("goal : at=%s battery=%s pass=%s\n", ds.Goal.At, ds.Goal.Battery, ds.Goal.Pass)
	fmt.Printf("thresholds : minBelief=%.2f maxCost=%.3f maxDuration=%.1f fuel=%d\n", ds.Thresholds.MinBelief, ds.Thresholds.MaxCost, ds.Thresholds.MaxDuration, ds.FuelSteps)
	fmt.Printf("actions : %d\n", len(ds.Actions))
	fmt.Printf("acceptable plans : %d\n", len(analysis.AcceptablePlans))
	fmt.Printf("selected final state : %s battery=%s pass=%s\n", best.FinalState.At, best.FinalState.Battery, best.FinalState.Pass)
	fmt.Printf("selected actions : %s\n", strings.Join(best.Actions, " -> "))
	fmt.Printf("selected metrics : duration=%.1f cost=%.3f belief=%.6f comfort=%.6f fuelRemaining=%d\n", best.Duration, best.Cost, best.Belief, best.Comfort, best.FuelRemaining)
	fmt.Printf("search recursive calls : %d\n", analysis.Stats.RecursiveCalls)
	fmt.Printf("search action tests : %d\n", analysis.Stats.ActionTests)
	fmt.Printf("search actions taken : %d\n", analysis.Stats.ActionsTaken)
	fmt.Printf("search state prunes : %d\n", analysis.Stats.StatePrunes)
	fmt.Printf("search max depth : %d\n", analysis.Stats.MaxDepth)
	fmt.Printf("checks passed : %d/%d\n", countChecksOK(analysis.Checks), len(analysis.Checks))
}
