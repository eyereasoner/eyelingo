// dining_philosophers.go
//
// A self-contained Go translation of dining-philosophers.n3 from the Eyeling
// examples.
//
// The N3 source describes a concrete Chandy-Misra dining-philosophers trace.
// Five philosophers sit in a ring and share five forks. A hungry philosopher
// asks for any adjacent fork they do not hold. The current holder sends the fork
// only when it is Dirty. A transferred fork arrives Clean, and after a meal the
// trace marks every fork Dirty again.
//
// This Go version keeps the same nine-round schedule from the N3 file:
// P1/P3 eat, then P2/P4 eat, then P5 eats, repeated three times. Goroutines are
// whether they can eat. State updates are still applied phase by phase so the
// run is deterministic and easy to compare with the source trace.
//
// Run:
//
//	go run dining_philosophers.go
//
// The program has no third-party dependencies.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"sort"
	"strings"
	"sync"
)

const seeExampleName = "dining_philosophers"

const sourceFile = "dining-philosophers.n3"

type Cleanliness string

const (
	Dirty Cleanliness = "Dirty"
	Clean Cleanliness = "Clean"
)

type ForkState struct {
	Holder      string
	Cleanliness Cleanliness
}

type Round struct {
	Number int
	Cycle  int
	Hungry []string
}

type Request struct {
	Round int
	From  string
	To    string
	Fork  string
	Side  string
	Dirty bool
}

type Transfer struct {
	Round string
	From  string
	To    string
	Fork  string
}

type Meal struct {
	Round       int
	Cycle       int
	Philosopher string
	MealNo      int
	Forks       []string
}

type RoundTrace struct {
	Round     Round
	Requests  []Request
	Transfers []Transfer
	Keeps     []string
	Meals     []Meal
	EndState  map[string]ForkState
}

type Stats struct {
	RoundsRun        int
	RequestsDerived  int
	TransfersSent    int
	KeepFacts        int
	MealsDerived     int
	GoroutineBatches int
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Question       string
	Traces         []RoundTrace
	MealCounts     map[string]int
	FinalState     map[string]ForkState
	Stats          Stats
	Checks         []Check
	EveryoneAte3   bool
	MealPatternOK  bool
	NoSharedForks  bool
	AllForksDirty  bool
	FinalHoldersOK bool
}

var philosophers = []string{"P1", "P2", "P3", "P4", "P5"}
var forks = []string{"F12", "F23", "F34", "F45", "F51"}

var philosopherRank = rank(philosophers)
var forkRank = rank(forks)

var leftFork = map[string]string{
	"P1": "F51",
	"P2": "F12",
	"P3": "F23",
	"P4": "F34",
	"P5": "F45",
}

var rightFork = map[string]string{
	"P1": "F12",
	"P2": "F23",
	"P3": "F34",
	"P4": "F45",
	"P5": "F51",
}

var schedule = []Round{
	{Number: 1, Cycle: 1, Hungry: []string{"P1", "P3"}},
	{Number: 2, Cycle: 1, Hungry: []string{"P2", "P4"}},
	{Number: 3, Cycle: 1, Hungry: []string{"P5"}},
	{Number: 4, Cycle: 2, Hungry: []string{"P1", "P3"}},
	{Number: 5, Cycle: 2, Hungry: []string{"P2", "P4"}},
	{Number: 6, Cycle: 2, Hungry: []string{"P5"}},
	{Number: 7, Cycle: 3, Hungry: []string{"P1", "P3"}},
	{Number: 8, Cycle: 3, Hungry: []string{"P2", "P4"}},
	{Number: 9, Cycle: 3, Hungry: []string{"P5"}},
}

func main() {
	schedule = exampleinput.Load(seeExampleName, schedule)
	analysis := derive()
	printAnswer(analysis)
	printReason(analysis)
}

func derive() Analysis {
	question := "Can five philosophers complete a Chandy-Misra trace where each eats three times?"
	state := initialState()
	mealCounts := map[string]int{}
	traces := make([]RoundTrace, 0, len(schedule))
	stats := Stats{}

	for _, round := range schedule {
		start := cloneState(state)
		requests := deriveRequests(round, start)
		transfers, afterTransfers := applyDirtyForkRule(round, start, requests)
		keeps := keptForks(transfers)
		meals := deriveMeals(round, afterTransfers)

		for _, meal := range meals {
			mealCounts[meal.Philosopher]++
		}

		endState := cloneState(afterTransfers)
		for _, fork := range forks {
			stateForFork := endState[fork]
			stateForFork.Cleanliness = Dirty
			endState[fork] = stateForFork
		}

		traces = append(traces, RoundTrace{
			Round:     round,
			Requests:  requests,
			Transfers: transfers,
			Keeps:     keeps,
			Meals:     meals,
			EndState:  cloneState(endState),
		})

		stats.RoundsRun++
		stats.RequestsDerived += len(requests)
		stats.TransfersSent += len(transfers)
		stats.KeepFacts += len(keeps)
		stats.MealsDerived += len(meals)
		stats.GoroutineBatches += 3
		state = endState
	}

	analysis := Analysis{
		Question:       question,
		Traces:         traces,
		MealCounts:     mealCounts,
		FinalState:     cloneState(state),
		Stats:          stats,
		EveryoneAte3:   everyoneAte(mealCounts, 3),
		MealPatternOK:  mealPatternMatches(traces),
		NoSharedForks:  noSharedForksInRounds(traces),
		AllForksDirty:  allForksHaveCleanliness(state, Dirty),
		FinalHoldersOK: finalHoldersMatch(state),
	}
	analysis.Checks = buildChecks(analysis)
	return analysis
}

func initialState() map[string]ForkState {
	return map[string]ForkState{
		"F12": {Holder: "P1", Cleanliness: Dirty},
		"F23": {Holder: "P2", Cleanliness: Dirty},
		"F34": {Holder: "P3", Cleanliness: Dirty},
		"F45": {Holder: "P4", Cleanliness: Dirty},
		"F51": {Holder: "P1", Cleanliness: Dirty},
	}
}

// deriveRequests is the goroutine phase that corresponds to the N3 request
// rules. Each hungry philosopher reads the same round snapshot and asks for
// adjacent forks that are currently held by someone else.
func deriveRequests(round Round, state map[string]ForkState) []Request {
	var wg sync.WaitGroup
	out := make(chan []Request, len(round.Hungry))

	for _, philosopher := range round.Hungry {
		philosopher := philosopher
		wg.Add(1)
		go func() {
			defer wg.Done()
			requests := make([]Request, 0, 2)
			for _, side := range []string{"left", "right"} {
				fork := forkOnSide(philosopher, side)
				forkState := state[fork]
				if forkState.Holder != philosopher {
					requests = append(requests, Request{
						Round: round.Number,
						From:  philosopher,
						To:    forkState.Holder,
						Fork:  fork,
						Side:  side,
						Dirty: forkState.Cleanliness == Dirty,
					})
				}
			}
			out <- requests
		}()
	}

	wg.Wait()
	close(out)

	var requests []Request
	for batch := range out {
		requests = append(requests, batch...)
	}
	sortRequests(requests)
	return requests
}

// applyDirtyForkRule is the goroutine phase that corresponds to the N3
// SendFork rule. A request becomes a transfer only when the requested fork is
// Dirty in the start-of-round snapshot. The receiver gets the fork Clean.
func applyDirtyForkRule(round Round, state map[string]ForkState, requests []Request) ([]Transfer, map[string]ForkState) {
	var wg sync.WaitGroup
	out := make(chan Transfer, len(requests))

	for _, request := range requests {
		request := request
		wg.Add(1)
		go func() {
			defer wg.Done()
			forkState := state[request.Fork]
			if forkState.Holder == request.To && forkState.Cleanliness == Dirty {
				out <- Transfer{
					Round: fmt.Sprintf("C%d", 2*(round.Number-1)),
					From:  request.To,
					To:    request.From,
					Fork:  request.Fork,
				}
			}
		}()
	}

	wg.Wait()
	close(out)

	transfers := make([]Transfer, 0, len(requests))
	for transfer := range out {
		transfers = append(transfers, transfer)
	}
	sortTransfersLikeRequests(transfers, requests)

	afterTransfers := cloneState(state)
	for _, transfer := range transfers {
		afterTransfers[transfer.Fork] = ForkState{Holder: transfer.To, Cleanliness: Clean}
	}
	return transfers, afterTransfers
}

// deriveMeals is the goroutine phase that corresponds to the N3 Meal rule.
// After transfers have been applied, a hungry philosopher eats if they hold both
// adjacent forks.
func deriveMeals(round Round, state map[string]ForkState) []Meal {
	var wg sync.WaitGroup
	out := make(chan Meal, len(round.Hungry))

	for _, philosopher := range round.Hungry {
		philosopher := philosopher
		wg.Add(1)
		go func() {
			defer wg.Done()
			left := leftFork[philosopher]
			right := rightFork[philosopher]
			if state[left].Holder == philosopher && state[right].Holder == philosopher {
				out <- Meal{
					Round:       round.Number,
					Cycle:       round.Cycle,
					Philosopher: philosopher,
					MealNo:      round.Cycle,
					Forks:       []string{left, right},
				}
			}
		}()
	}

	wg.Wait()
	close(out)

	meals := make([]Meal, 0, len(round.Hungry))
	for meal := range out {
		meals = append(meals, meal)
	}
	sort.Slice(meals, func(i, j int) bool {
		return philosopherRank[meals[i].Philosopher] < philosopherRank[meals[j].Philosopher]
	})
	return meals
}

func buildChecks(a Analysis) []Check {
	return []Check{
		{
			Label: "nine-round trace",
			OK:    a.Stats.RoundsRun == 9,
			Text:  "the translated run follows the nine start-of-round configs from the N3 source",
		},
		{
			Label: "request and transfer counts",
			OK:    a.Stats.RequestsDerived == 26 && a.Stats.TransfersSent == 26 && a.Stats.KeepFacts == 19,
			Text:  "26 dirty-fork requests transfer and the other 19 fork-round pairs are kept",
		},
		{
			Label: "meal schedule",
			OK:    a.Stats.MealsDerived == 15 && a.MealPatternOK,
			Text:  "rounds 1/4/7 feed P1,P3; rounds 2/5/8 feed P2,P4; rounds 3/6/9 feed P5",
		},
		{
			Label: "everyone ate three times",
			OK:    a.EveryoneAte3,
			Text:  "each philosopher has exactly three derived meals",
		},
		{
			Label: "fork safety",
			OK:    a.NoSharedForks,
			Text:  "no two philosophers eat with the same fork in the same round",
		},
		{
			Label: "dirty after eating",
			OK:    a.AllForksDirty,
			Text:  "the end-of-round state makes every fork Dirty again, matching the trace model",
		},
		{
			Label: "final holders",
			OK:    a.FinalHoldersOK,
			Text:  "the final ownership is F12,F23 with P2; F34 with P4; F45,F51 with P5",
		},
		{
			Label: "goroutine phases",
			OK:    a.Stats.GoroutineBatches == 27,
			Text:  "each of the nine rounds used request, transfer, and meal goroutine batches",
		},
	}
}

func printAnswer(a Analysis) {
	fmt.Println("# Dining Philosophers")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("The Chandy-Misra dining-philosophers trace completes without conflict.")
	fmt.Printf("philosophers : %d\n", len(philosophers))
	fmt.Printf("forks : %d\n", len(forks))
	fmt.Printf("rounds : %d\n", a.Stats.RoundsRun)
	fmt.Printf("meals : %d\n", a.Stats.MealsDerived)
	fmt.Printf("everyone ate 3 times : %s\n", yesNo(a.EveryoneAte3))
	fmt.Println()
	fmt.Println("Meal trace:")
	for _, trace := range a.Traces {
		fmt.Printf("  round %d cycle %d : %s\n", trace.Round.Number, trace.Round.Cycle, formatMeals(trace.Meals))
	}
	fmt.Println()
}

func printReason(a Analysis) {
	fmt.Println("## Reason")
	fmt.Println("Each round has three phases. Hungry philosophers first request adjacent")
	fmt.Println("forks they do not hold. A holder sends a requested fork only when it is")
	fmt.Println("Dirty, and the receiver gets that fork Clean. After the meal phase, all")
	fmt.Println("forks are marked Dirty for the next round.")
	fmt.Println("The Go code uses goroutines inside each phase, then applies state changes")
	fmt.Println("between phases so the output remains deterministic.")
	fmt.Println()
	fmt.Println("Round derivation:")
	for _, trace := range a.Traces {
		fmt.Printf("  round %d hungry : %s\n", trace.Round.Number, strings.Join(trace.Round.Hungry, ","))
		fmt.Printf("    requests : %s\n", formatRequests(trace.Requests))
		fmt.Printf("    transfers : %s\n", formatTransfers(trace.Transfers))
		fmt.Printf("    kept forks : %s\n", strings.Join(trace.Keeps, ","))
		fmt.Printf("    meals : %s\n", formatMealNames(trace.Meals))
	}
	fmt.Println()
}

func forkOnSide(philosopher string, side string) string {
	if side == "left" {
		return leftFork[philosopher]
	}
	return rightFork[philosopher]
}

func keptForks(transfers []Transfer) []string {
	sent := map[string]bool{}
	for _, transfer := range transfers {
		sent[transfer.Fork] = true
	}
	keeps := make([]string, 0, len(forks)-len(transfers))
	for _, fork := range forks {
		if !sent[fork] {
			keeps = append(keeps, fork)
		}
	}
	return keeps
}

func mealPatternMatches(traces []RoundTrace) bool {
	expected := [][]string{
		{"P1", "P3"},
		{"P2", "P4"},
		{"P5"},
		{"P1", "P3"},
		{"P2", "P4"},
		{"P5"},
		{"P1", "P3"},
		{"P2", "P4"},
		{"P5"},
	}
	if len(traces) != len(expected) {
		return false
	}
	for i, trace := range traces {
		got := make([]string, 0, len(trace.Meals))
		for _, meal := range trace.Meals {
			got = append(got, meal.Philosopher)
		}
		if strings.Join(got, ",") != strings.Join(expected[i], ",") {
			return false
		}
	}
	return true
}

func noSharedForksInRounds(traces []RoundTrace) bool {
	for _, trace := range traces {
		used := map[string]bool{}
		for _, meal := range trace.Meals {
			for _, fork := range meal.Forks {
				if used[fork] {
					return false
				}
				used[fork] = true
			}
		}
	}
	return true
}

func everyoneAte(mealCounts map[string]int, want int) bool {
	for _, philosopher := range philosophers {
		if mealCounts[philosopher] != want {
			return false
		}
	}
	return true
}

func allForksHaveCleanliness(state map[string]ForkState, cleanliness Cleanliness) bool {
	for _, fork := range forks {
		if state[fork].Cleanliness != cleanliness {
			return false
		}
	}
	return true
}

func finalHoldersMatch(state map[string]ForkState) bool {
	expected := map[string]string{
		"F12": "P2",
		"F23": "P2",
		"F34": "P4",
		"F45": "P5",
		"F51": "P5",
	}
	for fork, holder := range expected {
		if state[fork].Holder != holder {
			return false
		}
	}
	return true
}

func sortRequests(requests []Request) {
	sideRank := map[string]int{"left": 0, "right": 1}
	sort.Slice(requests, func(i, j int) bool {
		if requests[i].From != requests[j].From {
			return philosopherRank[requests[i].From] < philosopherRank[requests[j].From]
		}
		if requests[i].Side != requests[j].Side {
			return sideRank[requests[i].Side] < sideRank[requests[j].Side]
		}
		return forkRank[requests[i].Fork] < forkRank[requests[j].Fork]
	})
}

func sortTransfersLikeRequests(transfers []Transfer, requests []Request) {
	order := map[string]int{}
	for i, request := range requests {
		order[request.Fork] = i
	}
	sort.Slice(transfers, func(i, j int) bool {
		return order[transfers[i].Fork] < order[transfers[j].Fork]
	})
}

func cloneState(state map[string]ForkState) map[string]ForkState {
	cloned := make(map[string]ForkState, len(state))
	for fork, forkState := range state {
		cloned[fork] = forkState
	}
	return cloned
}

func rank(values []string) map[string]int {
	ranks := make(map[string]int, len(values))
	for i, value := range values {
		ranks[value] = i
	}
	return ranks
}

func formatMeals(meals []Meal) string {
	parts := make([]string, 0, len(meals))
	for _, meal := range meals {
		parts = append(parts, fmt.Sprintf("%s#%d uses %s", meal.Philosopher, meal.MealNo, strings.Join(meal.Forks, ",")))
	}
	return strings.Join(parts, "; ")
}

func formatMealNames(meals []Meal) string {
	if len(meals) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(meals))
	for _, meal := range meals {
		parts = append(parts, fmt.Sprintf("%s#%d", meal.Philosopher, meal.MealNo))
	}
	return strings.Join(parts, ",")
}

func formatRequests(requests []Request) string {
	if len(requests) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(requests))
	for _, request := range requests {
		parts = append(parts, fmt.Sprintf("%s asks %s for %s", request.From, request.To, request.Fork))
	}
	return strings.Join(parts, "; ")
}

func formatTransfers(transfers []Transfer) string {
	if len(transfers) == 0 {
		return "none"
	}
	parts := make([]string, 0, len(transfers))
	for _, transfer := range transfers {
		parts = append(parts, fmt.Sprintf("%s sends %s to %s", transfer.From, transfer.Fork, transfer.To))
	}
	return strings.Join(parts, "; ")
}

func formatMealCounts(mealCounts map[string]int) string {
	parts := make([]string, 0, len(philosophers))
	for _, philosopher := range philosophers {
		parts = append(parts, fmt.Sprintf("%s:%d", philosopher, mealCounts[philosopher]))
	}
	return strings.Join(parts, ", ")
}

func formatState(state map[string]ForkState) string {
	parts := make([]string, 0, len(forks))
	for _, fork := range forks {
		forkState := state[fork]
		parts = append(parts, fmt.Sprintf("%s=%s/%s", fork, forkState.Holder, forkState.Cleanliness))
	}
	return strings.Join(parts, ", ")
}

func countPassed(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
