// gps.go
//
// A self-contained Go translation of examples/gps.n3 from the Eyeling example
// suite.
//
// The original N3 example is an ARC-style, goal-driven route planner for a tiny
// western-Belgium map. It starts from Gent, derives possible paths to Oostende,
// compares their computed metrics, recommends the better route, renders a small
// explanation, and checks that the recommendation is consistent with the route
// metrics.
//
// This is intentionally not a generic RDF/N3 reasoner. The concrete N3 facts
// and rules are represented as Go structs and ordinary functions so the path
// derivation and decision logic are easy to read and directly runnable.
//
// Run:
//
//	go run gps.go
//
// The program has no third-party dependencies.
package main

import (
	"errors"
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "gps"

const (
	locationGent       = "Gent"
	locationBrugge     = "Brugge"
	locationKortrijk   = "Kortrijk"
	locationOostende   = "Oostende"
	routeDirectID      = "routeDirect"
	routeViaKortrijkID = "routeViaKortrijk"
)

// Dataset contains the fixed facts from gps.n3.
//
// The N3 source stores terms such as :Gent and :drive_gent_brugge. In this
// concrete translation, those symbols become ordinary strings.
type Dataset struct {
	Traveller Traveller
	Question  string
	Routes    map[string]RouteDefinition
	Edges     []MapDescription
}

// Traveller represents the subject whose current location is used as the start
// state for the path query.
type Traveller struct {
	ID       string
	Location string
}

// RouteDefinition stores the human-readable labels assigned to the two known
// concrete routes in the N3 fixture.
type RouteDefinition struct {
	ID    string
	Label string
}

// MapDescription is the Go form of one gps:description fact.
//
// Each edge says that, from one location, a single action can drive to another
// location with a duration, cost, belief score, and comfort score.
type MapDescription struct {
	From     string
	To       string
	Action   string
	Duration float64
	Cost     float64
	Belief   float64
	Comfort  float64
}

// Path is a derived route from one location to another.
//
// Duration and Cost are summed across edges, while Belief and Comfort are
// multiplied, matching the math:sum and math:product rules in the N3 source.
type Path struct {
	From     string
	To       string
	Actions  []string
	Duration float64
	Cost     float64
	Belief   float64
	Comfort  float64
}

// Decision is the final recommendation produced after comparing the named
// routes.
type Decision struct {
	RecommendedRoute string
	Outcome          string
}

// Checks mirrors the proof obligations at the bottom of gps.n3.
//
// The first two checks confirm that both concrete routes were derived. The last
// three checks confirm that the direct route really dominates the alternative.
type Checks struct {
	DirectRouteDerived      bool
	AlternativeRouteDerived bool
	RecommendedIsFaster     bool
	RecommendedIsCheaper    bool
	RecommendedScoresHigher bool
}

// SearchStats records operational details from the Go depth-first path search.
//
// These counters are not part of the original N3 proof. They make the Go
// translation easier to audit by showing how much graph exploration was needed
// to derive the candidate routes.
type SearchStats struct {
	RecursiveCalls int
	EdgeTests      int
	EdgesExtended  int
	RevisitPrunes  int
	MaxDepth       int
}

// InferenceResult contains all derived facts needed by the renderer.
type InferenceResult struct {
	AllPaths         []Path
	DirectRoute      Path
	AlternativeRoute Path
	Decision         Decision
	Checks           Checks
	Stats            SearchStats
}

// dataset returns the hard-coded facts translated from gps.n3.
func dataset() Dataset {
	return Dataset{
		Traveller: Traveller{
			ID:       "i1",
			Location: locationGent,
		},
		Question: "Which route should we take from Gent to Oostende?",
		Routes: map[string]RouteDefinition{
			routeDirectID: {
				ID:    routeDirectID,
				Label: "Gent → Brugge → Oostende",
			},
			routeViaKortrijkID: {
				ID:    routeViaKortrijkID,
				Label: "Gent → Kortrijk → Brugge → Oostende",
			},
		},
		Edges: []MapDescription{
			{
				From:     locationGent,
				To:       locationBrugge,
				Action:   "drive_gent_brugge",
				Duration: 1500.0,
				Cost:     0.006,
				Belief:   0.96,
				Comfort:  0.99,
			},
			{
				From:     locationGent,
				To:       locationKortrijk,
				Action:   "drive_gent_kortrijk",
				Duration: 1600.0,
				Cost:     0.007,
				Belief:   0.96,
				Comfort:  0.99,
			},
			{
				From:     locationKortrijk,
				To:       locationBrugge,
				Action:   "drive_kortrijk_brugge",
				Duration: 1600.0,
				Cost:     0.007,
				Belief:   0.96,
				Comfort:  0.99,
			},
			{
				From:     locationBrugge,
				To:       locationOostende,
				Action:   "drive_brugge_oostende",
				Duration: 900.0,
				Cost:     0.004,
				Belief:   0.98,
				Comfort:  1.0,
			},
		},
	}
}

// infer runs the translated N3 rules in a deterministic order.
func infer(data Dataset) (InferenceResult, error) {
	// Query the path relation for the traveller starting in Gent and aiming for
	// Oostende. This corresponds to the N3 rule that derives :i1 gps:path ... .
	paths, stats := findPaths(data.Edges, data.Traveller.Location, locationOostende)

	// Name the two concrete routes that exist in this tiny map.
	directRoute, directOK := findPathByActions(paths, []string{
		"drive_gent_brugge",
		"drive_brugge_oostende",
	})
	alternativeRoute, alternativeOK := findPathByActions(paths, []string{
		"drive_gent_kortrijk",
		"drive_kortrijk_brugge",
		"drive_brugge_oostende",
	})

	checks := Checks{
		DirectRouteDerived:      directOK,
		AlternativeRouteDerived: alternativeOK,
	}
	if !directOK || !alternativeOK {
		return InferenceResult{AllPaths: paths, Checks: checks, Stats: stats}, errors.New("expected routes were not derived")
	}

	// Decision rule: the direct route is recommended only if it is faster,
	// cheaper, more reliable, and more comfortable than the alternative.
	checks.RecommendedIsFaster = directRoute.Duration < alternativeRoute.Duration
	checks.RecommendedIsCheaper = directRoute.Cost < alternativeRoute.Cost
	checks.RecommendedScoresHigher = directRoute.Belief > alternativeRoute.Belief &&
		directRoute.Comfort > alternativeRoute.Comfort

	if !allChecksPass(checks) {
		return InferenceResult{
			AllPaths:         paths,
			DirectRoute:      directRoute,
			AlternativeRoute: alternativeRoute,
			Checks:           checks,
			Stats:            stats,
		}, fmt.Errorf("recommendation checks failed: %+v", checks)
	}

	decision := Decision{
		RecommendedRoute: routeDirectID,
		Outcome:          "Take the direct route via Brugge.",
	}

	return InferenceResult{
		AllPaths:         paths,
		DirectRoute:      directRoute,
		AlternativeRoute: alternativeRoute,
		Decision:         decision,
		Checks:           checks,
		Stats:            stats,
	}, nil
}

// findPaths computes all simple paths between two locations.
//
// The base case mirrors the N3 rule saying that a single gps:description is
// already a path. The recursive case extends an edge with the rest of a path,
// appending actions and combining metrics along the way.
func findPaths(edges []MapDescription, from, to string) ([]Path, SearchStats) {
	visited := map[string]bool{from: true}
	stats := SearchStats{}
	paths := findPathsRecursive(edges, from, to, visited, Path{
		From:    from,
		To:      from,
		Belief:  1.0,
		Comfort: 1.0,
	}, 0, &stats)
	return paths, stats
}

func findPathsRecursive(edges []MapDescription, current, target string, visited map[string]bool, partial Path, depth int, stats *SearchStats) []Path {
	stats.RecursiveCalls++
	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	var paths []Path

	for _, edge := range edges {
		stats.EdgeTests++
		if edge.From != current {
			continue
		}
		if visited[edge.To] {
			stats.RevisitPrunes++
			continue
		}

		// Extend the partial path with this edge. Sums/products mirror the N3
		// math:sum and math:product built-ins.
		stats.EdgesExtended++
		next := Path{
			From:     partial.From,
			To:       edge.To,
			Actions:  appendAction(partial.Actions, edge.Action),
			Duration: partial.Duration + edge.Duration,
			Cost:     partial.Cost + edge.Cost,
			Belief:   partial.Belief * edge.Belief,
			Comfort:  partial.Comfort * edge.Comfort,
		}

		if edge.To == target {
			paths = append(paths, next)
			continue
		}

		visited[edge.To] = true
		paths = append(paths, findPathsRecursive(edges, edge.To, target, visited, next, depth+1, stats)...)
		delete(visited, edge.To)
	}

	return paths
}

// appendAction copies the existing action list before appending so recursive
// branches cannot accidentally share and mutate the same backing array.
func appendAction(actions []string, action string) []string {
	out := make([]string, 0, len(actions)+1)
	out = append(out, actions...)
	out = append(out, action)
	return out
}

// findPathByActions identifies the derived path corresponding to one named
// concrete route in the N3 source.
func findPathByActions(paths []Path, actions []string) (Path, bool) {
	for _, path := range paths {
		if sameActions(path.Actions, actions) {
			return path, true
		}
	}
	return Path{}, false
}

func sameActions(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// allChecksPass is the Go counterpart of the N3 verification rules that fail
// loudly when a recommended route contradicts the computed metrics.
func allChecksPass(checks Checks) bool {
	return checks.DirectRouteDerived &&
		checks.AlternativeRouteDerived &&
		checks.RecommendedIsFaster &&
		checks.RecommendedIsCheaper &&
		checks.RecommendedScoresHigher
}

// checkCount returns the number of passed verification checks for compact audit
// output. The total is fixed because gps.n3 has five explicit consistency checks.
func checkCount(checks Checks) int {
	count := 0
	if checks.DirectRouteDerived {
		count++
	}
	if checks.AlternativeRouteDerived {
		count++
	}
	if checks.RecommendedIsFaster {
		count++
	}
	if checks.RecommendedIsCheaper {
		count++
	}
	if checks.RecommendedScoresHigher {
		count++
	}
	return count
}

// actionPath renders a derived path as a readable action sequence.
func actionPath(path Path) string {
	if len(path.Actions) == 0 {
		return "(none)"
	}
	return strings.Join(path.Actions, " -> ")
}

// metricDelta reports how much better the direct route is for metrics where
// lower values are better.
func metricDelta(alternative, direct float64) float64 {
	return alternative - direct
}

// scoreDelta reports how much better the direct route is for metrics where
// higher values are better.
func scoreDelta(direct, alternative float64) float64 {
	return direct - alternative
}

// renderPathAudit prints every derived path with its action sequence and
// accumulated metrics. This is useful for checking that the Go DFS found the
// same candidates that the N3 recursive path rule would derive.
func renderPathAudit(paths []Path) {
	for i, path := range paths {
		fmt.Printf(
			"derived path %d : %s | duration=%s cost=%s belief=%s comfort=%s\n",
			i+1,
			actionPath(path),
			formatDuration(path.Duration),
			formatDecimal(path.Cost, 3),
			formatDecimal(path.Belief, 6),
			formatDecimal(path.Comfort, 4),
		)
	}
}

// renderArcOutput prints the same answer / reason / check style as the N3
// string:concatenation + log:outputString section.
func renderArcOutput(data Dataset, result InferenceResult) {
	directLabel := data.Routes[routeDirectID].Label
	alternativeLabel := data.Routes[routeViaKortrijkID].Label

	fmt.Println("# GPS — Goal driven route planning")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println(result.Decision.Outcome)
	fmt.Printf("Recommended route: %s\n", directLabel)

	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("From Gent to Oostende, the planner found two routes in this small map.")
	fmt.Printf(
		"The direct route (%s) takes %s seconds at cost %s, with belief %s and comfort %s. ",
		directLabel,
		formatDuration(result.DirectRoute.Duration),
		formatDecimal(result.DirectRoute.Cost, 3),
		formatDecimal(result.DirectRoute.Belief, 6),
		formatDecimal(result.DirectRoute.Comfort, 4),
	)
	fmt.Printf(
		"The alternative (%s) takes %s seconds at cost %s, with belief %s and comfort %s.\n",
		alternativeLabel,
		formatDuration(result.AlternativeRoute.Duration),
		formatDecimal(result.AlternativeRoute.Cost, 3),
		formatDecimal(result.AlternativeRoute.Belief, 6),
		formatDecimal(result.AlternativeRoute.Comfort, 4),
	)
	fmt.Println("So the direct route is faster, cheaper, more reliable, and slightly more comfortable.")

	fmt.Println()
	fmt.Println("## Check")
	fmt.Println("C1 OK - the direct Gent → Brugge → Oostende route was derived.")
	fmt.Println("C2 OK - the alternative Gent → Kortrijk → Brugge → Oostende route was derived.")
	fmt.Println("C3 OK - the recommended route is faster than the alternative.")
	fmt.Println("C4 OK - the recommended route is cheaper than the alternative.")
	fmt.Println("C5 OK - the recommended route has higher belief and comfort scores.")

	// Extra audit details are useful in this Go translation because they show the
	// concrete graph search, route comparisons, and runtime context that led to
	// the recommendation.
	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("question : %s\n", data.Question)
	fmt.Printf("traveller : %s\n", data.Traveller.ID)
	fmt.Printf("start : %s\n", data.Traveller.Location)
	fmt.Printf("goal : %s\n", locationOostende)
	fmt.Printf("map edges : %d\n", len(data.Edges))
	fmt.Printf("named routes : %d\n", len(data.Routes))
	fmt.Printf("paths derived : %d\n", len(result.AllPaths))
	renderPathAudit(result.AllPaths)
	fmt.Printf("direct actions : %s\n", actionPath(result.DirectRoute))
	fmt.Printf("alternative actions : %s\n", actionPath(result.AlternativeRoute))
	fmt.Printf("duration advantage seconds : %s\n", formatDuration(metricDelta(result.AlternativeRoute.Duration, result.DirectRoute.Duration)))
	fmt.Printf("cost advantage : %s\n", formatDecimal(metricDelta(result.AlternativeRoute.Cost, result.DirectRoute.Cost), 3))
	fmt.Printf("belief advantage : %s\n", formatDecimal(scoreDelta(result.DirectRoute.Belief, result.AlternativeRoute.Belief), 6))
	fmt.Printf("comfort advantage : %s\n", formatDecimal(scoreDelta(result.DirectRoute.Comfort, result.AlternativeRoute.Comfort), 4))
	fmt.Printf("search recursive calls : %d\n", result.Stats.RecursiveCalls)
	fmt.Printf("search edge tests : %d\n", result.Stats.EdgeTests)
	fmt.Printf("search edges extended : %d\n", result.Stats.EdgesExtended)
	fmt.Printf("search revisit prunes : %d\n", result.Stats.RevisitPrunes)
	fmt.Printf("search max depth : %d\n", result.Stats.MaxDepth)
	fmt.Printf("checks passed : %d/5\n", checkCount(result.Checks))
	fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(result.Checks)))
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func formatDuration(value float64) string {
	return fmt.Sprintf("%.1f", value)
}

// formatDecimal trims insignificant trailing zeroes while keeping enough
// precision to show the route comparisons clearly.
func formatDecimal(value float64, decimals int) string {
	text := fmt.Sprintf("%.*f", decimals, value)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	if text == "" || text == "-0" {
		return "0"
	}
	return text
}

func main() {
	data := exampleinput.Load(eyelingoExampleName, dataset())
	result, err := infer(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GPS inference failed: %v\n", err)
		os.Exit(1)
	}
	renderArcOutput(data, result)
}
