// gps.go
//
// A self-contained Go translation of examples/gps.n3 from the Eyeling example
// suite.
//
// The original N3 example is an ARC-style, goal-driven route planner for a tiny
// western-Belgium map. It starts from Gent, derives possible paths to Oostende,
// compares their computed metrics, recommends the better route, renders a small
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
	"fmt"
	"os"
	"see/internal/exampleinput"
	"strings"
)

const seeExampleName = "gps"

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
	Stops    []string
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

// routeStops parses the route label fact into the sequence of locations it names.
// This keeps the Go derivation from hard-coding action lists for the known
// answers: candidate routes are matched against paths inferred from map edges.
func routeStops(label string) []string {
	parts := strings.Split(label, "→")
	if len(parts) == 1 {
		parts = strings.Split(label, "->")
	}
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

// infer runs the translated N3 rules in a deterministic order.
func infer(data Dataset) (InferenceResult, error) {
	// Query the path relation for the traveller starting in Gent and aiming for
	// Oostende. This corresponds to the N3 rule that derives :i1 gps:path ... .
	paths, stats := findPaths(data.Edges, data.Traveller.Location, locationOostende)

	// Match the named route facts to paths derived from the map. The labels live
	// in the input data; the action sequences are not baked into the answer path.
	directRoute, directOK := findPathByStops(paths, routeStops(data.Routes[routeDirectID].Label))
	alternativeRoute, alternativeOK := findPathByStops(paths, routeStops(data.Routes[routeViaKortrijkID].Label))

	checks := Checks{
		DirectRouteDerived:      directOK,
		AlternativeRouteDerived: alternativeOK,
	}
	if !directOK || !alternativeOK {
		return InferenceResult{AllPaths: paths, Checks: checks, Stats: stats}, errors.New("expected routes were not derived")
	}

	// Decision rule: recommend a route only from computed metrics. The N3 fixture
	// contains a rule for :routeDirect that fires when it is faster, cheaper,
	// more reliable, and more comfortable than :routeViaKortrijk; the Go code
	// must therefore derive the recommendation from those inequalities rather
	// than unconditionally naming routeDirect.
	checks.RecommendedIsFaster = directRoute.Duration < alternativeRoute.Duration
	checks.RecommendedIsCheaper = directRoute.Cost < alternativeRoute.Cost
	checks.RecommendedScoresHigher = directRoute.Belief > alternativeRoute.Belief &&
		directRoute.Comfort > alternativeRoute.Comfort

	decision := chooseRoute(data, directRoute, alternativeRoute, checks)

	return InferenceResult{
		AllPaths:         paths,
		DirectRoute:      directRoute,
		AlternativeRoute: alternativeRoute,
		Decision:         decision,
		Checks:           checks,
		Stats:            stats,
	}, nil
}

// chooseRoute mirrors the N3 decision rule without baking the conclusion into
// the Answer path. If the translated facts change so that the inequalities no
// longer hold, no direct-route recommendation is emitted.
func chooseRoute(data Dataset, directRoute, alternativeRoute Path, checks Checks) Decision {
	if checks.RecommendedIsFaster && checks.RecommendedIsCheaper && checks.RecommendedScoresHigher {
		return Decision{
			RecommendedRoute: routeDirectID,
			Outcome:          fmt.Sprintf("Take the %s.", routeSummary(data.Routes[routeDirectID].Label)),
		}
	}

	if alternativeRoute.Duration < directRoute.Duration &&
		alternativeRoute.Cost < directRoute.Cost &&
		alternativeRoute.Belief > directRoute.Belief &&
		alternativeRoute.Comfort > directRoute.Comfort {
		return Decision{
			RecommendedRoute: routeViaKortrijkID,
			Outcome:          fmt.Sprintf("Take the %s.", routeSummary(data.Routes[routeViaKortrijkID].Label)),
		}
	}

	return Decision{
		RecommendedRoute: "",
		Outcome:          "No route is recommended because no candidate dominates on all metrics.",
	}
}

func routeSummary(label string) string {
	if strings.Contains(label, "Brugge") && !strings.Contains(label, "Kortrijk") {
		return "direct route via Brugge"
	}
	return label
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
		Stops:   []string{from},
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
			Stops:    appendStop(partial.Stops, edge.To),
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

func appendStop(stops []string, stop string) []string {
	out := make([]string, 0, len(stops)+1)
	out = append(out, stops...)
	out = append(out, stop)
	return out
}

// findPathByStops identifies a derived path by the location sequence supplied
// as input facts, rather than by hard-coded action names in the Go program.
func findPathByStops(paths []Path, stops []string) (Path, bool) {
	for _, path := range paths {
		if sameStrings(path.Stops, stops) {
			return path, true
		}
	}
	return Path{}, false
}

func sameStrings(left, right []string) bool {
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

// loudly when a recommended route contradicts the computed metrics.
func allChecksPass(checks Checks) bool {
	return checks.DirectRouteDerived &&
		checks.AlternativeRouteDerived &&
		checks.RecommendedIsFaster &&
		checks.RecommendedIsCheaper &&
		checks.RecommendedScoresHigher
}

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

// same candidates that the N3 recursive path rule would derive.

// string:concatenation + log:outputString section.
func renderArcOutput(data Dataset, result InferenceResult) {
	directLabel := data.Routes[routeDirectID].Label
	alternativeLabel := data.Routes[routeViaKortrijkID].Label
	recommendedLabel := "none"
	if result.Decision.RecommendedRoute != "" {
		recommendedLabel = data.Routes[result.Decision.RecommendedRoute].Label
	}

	fmt.Println("# GPS — Goal driven route planning")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println(result.Decision.Outcome)
	fmt.Printf("Recommended route: %s\n", recommendedLabel)

	fmt.Println()
	fmt.Println("## Reason")
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
	if result.Decision.RecommendedRoute == routeDirectID {
		fmt.Println("So the direct route is faster, cheaper, more reliable, and slightly more comfortable.")
	} else {
		fmt.Println("No candidate dominates on all route metrics, so the planner withholds a recommendation.")
	}

	fmt.Println()
	return
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
	data := exampleinput.Load(seeExampleName, Dataset{})
	result, err := infer(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GPS inference failed: %v\n", err)
		os.Exit(1)
	}
	renderArcOutput(data, result)
}
