// path_discovery.go
//
// A self-contained Go translation of path-discovery.n3 from the Eyeling
// examples.
//
// The original N3 file contains a large Neptune air-routes graph and a bounded
// recursive rule:
//
//	(?source ?destination () 0 2) :route ?airports
//
// This Go version translates the query, the recursive no-revisit route rule,
// and the complete airport/flight data loaded from examples/input/path_discovery.json:
// 7,698 airport labels and 37,505 nepo:hasOutboundRouteTo facts. The bounded
// query still touches only the 338 outbound candidates reachable from
// Ostend-Bruges within the two-stopover bound, but the full graph is loaded
// and checked.
//
// This is intentionally not a generic RDF/N3 reasoner. The concrete route rules
// are represented as ordinary Go functions, and the concrete airport labels and
// flight facts are loaded from JSON input so the derivation remains visible and
// directly runnable.
//
// Run:
//
//	go run path_discovery.go
//
// The program has no third-party dependencies.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "path_discovery"

const (
	sourceGraphAirportLabels = 7698
	sourceGraphOutboundFacts = 37505
	maxStopovers             = 2
	maxHops                  = maxStopovers + 1

	sourceAirport      = "res:AIRPORT_310"
	destinationAirport = "res:AIRPORT_1587"
	mandatoryFirstHop  = "res:AIRPORT_309"
)

type Dataset struct {
	Question      string
	SourceID      string
	DestinationID string
	Labels        map[string]string
	Edges         []Edge
}

type Edge struct {
	From string
	To   string
}

type Route struct {
	Airports []string
}

type SearchStats struct {
	RecursiveCalls   int
	EdgeTests        int
	EdgesExtended    int
	RevisitPrunes    int
	DepthLimitLeaves int
	DeadEnds         int
	RoutesEmitted    int
	MaxDepth         int
}

type Checks struct {
	SourceAndDestinationKnown bool
	FirstHopMatchesN3Facts    bool
	RouteSetMatchesN3Query    bool
	NoShorterRouteExists      bool
	RoutesWithinStopoverLimit bool
	EveryHopHasFact           bool
	NoAirportRevisited        bool
	FullSourceGraphLoaded     bool
	RoutesSortedDeterministic bool
}

type InferenceResult struct {
	Routes               []Route
	Stats                SearchStats
	Checks               Checks
	LoadedAirportLabels  int
	LoadedOutboundFacts  int
	EdgeEndpointAirports int
	ExpandedAirports     []string
	SourceOut            []string
	FirstHopOut          []string
	DirectRoutes         int
	OneStopRoutes        int
	TwoStopRoutes        int
}

func infer(data Dataset) InferenceResult {
	adj := buildAdjacency(data)
	edgeSet := buildEdgeSet(data.Edges)
	stats := SearchStats{}
	routes := []Route{}
	dfs(data.SourceID, data.DestinationID, adj, []string{data.SourceID}, &routes, &stats)
	sortRoutes(data, routes)

	direct, oneStop, twoStop := routeDistribution(routes)
	expanded := expandedAirports(data.SourceID, adj, 2)
	firstHopOut := append([]string{}, adj[mandatoryFirstHop]...)
	sortTermsByLabel(data, firstHopOut)

	checks := Checks{
		SourceAndDestinationKnown: data.label(data.SourceID) != data.SourceID && data.label(data.DestinationID) != data.DestinationID,
		FirstHopMatchesN3Facts:    len(adj[data.SourceID]) == 1 && adj[data.SourceID][0] == mandatoryFirstHop,
		RouteSetMatchesN3Query:    routeSetMatches(data, routes, expectedRoutes()),
		NoShorterRouteExists:      direct == 0 && oneStop == 0,
		RoutesWithinStopoverLimit: routesWithinStopovers(routes, maxStopovers),
		EveryHopHasFact:           everyHopHasFact(routes, edgeSet),
		NoAirportRevisited:        noAirportRevisited(routes),
		FullSourceGraphLoaded:     len(data.Labels) == sourceGraphAirportLabels && len(data.Edges) == sourceGraphOutboundFacts,
		RoutesSortedDeterministic: routesSorted(data, routes),
	}

	return InferenceResult{
		Routes:               routes,
		Stats:                stats,
		Checks:               checks,
		LoadedAirportLabels:  len(data.Labels),
		LoadedOutboundFacts:  len(data.Edges),
		EdgeEndpointAirports: countAirportsInEdges(data.Edges),
		ExpandedAirports:     expanded,
		SourceOut:            append([]string{}, adj[data.SourceID]...),
		FirstHopOut:          firstHopOut,
		DirectRoutes:         direct,
		OneStopRoutes:        oneStop,
		TwoStopRoutes:        twoStop,
	}
}

func buildAdjacency(data Dataset) map[string][]string {
	adj := map[string][]string{}
	for _, edge := range data.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
	}
	for from := range adj {
		sortTermsByLabel(data, adj[from])
	}
	return adj
}

func buildEdgeSet(edges []Edge) map[string]bool {
	set := map[string]bool{}
	for _, edge := range edges {
		set[edge.From+" -> "+edge.To] = true
	}
	return set
}

func dfs(current, destination string, adj map[string][]string, path []string, routes *[]Route, stats *SearchStats) {
	depth := len(path) - 1
	stats.RecursiveCalls++
	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	if len(path) > 1 && current == destination {
		stats.RoutesEmitted++
		copyPath := append([]string{}, path...)
		*routes = append(*routes, Route{Airports: copyPath})
		return
	}

	if depth >= maxHops {
		stats.DepthLimitLeaves++
		return
	}

	outbound := adj[current]
	stats.EdgeTests += len(outbound)
	if len(outbound) == 0 {
		stats.DeadEnds++
	}

	for _, next := range outbound {
		if contains(path, next) {
			stats.RevisitPrunes++
			continue
		}
		stats.EdgesExtended++
		nextPath := append(append([]string{}, path...), next)
		dfs(next, destination, adj, nextPath, routes, stats)
	}
}

func expectedRoutes() [][]string {
	return [][]string{
		{sourceAirport, mandatoryFirstHop, "res:AIRPORT_1472", destinationAirport},
		{sourceAirport, mandatoryFirstHop, "res:AIRPORT_1452", destinationAirport},
		{sourceAirport, mandatoryFirstHop, "res:AIRPORT_3998", destinationAirport},
	}
}

func routeSetMatches(data Dataset, actual []Route, expected [][]string) bool {
	actualKeys := make([]string, 0, len(actual))
	for _, route := range actual {
		actualKeys = append(actualKeys, strings.Join(route.Airports, "|"))
	}
	expectedKeys := make([]string, 0, len(expected))
	for _, route := range expected {
		expectedKeys = append(expectedKeys, strings.Join(route, "|"))
	}
	sort.Strings(actualKeys)
	sort.Strings(expectedKeys)
	if len(actualKeys) != len(expectedKeys) {
		return false
	}
	for i := range actualKeys {
		if actualKeys[i] != expectedKeys[i] {
			return false
		}
	}
	return routesSorted(data, actual)
}

func routesWithinStopovers(routes []Route, limit int) bool {
	for _, route := range routes {
		if route.Stopovers() > limit {
			return false
		}
	}
	return true
}

func everyHopHasFact(routes []Route, edgeSet map[string]bool) bool {
	for _, route := range routes {
		for i := 0; i < len(route.Airports)-1; i++ {
			key := route.Airports[i] + " -> " + route.Airports[i+1]
			if !edgeSet[key] {
				return false
			}
		}
	}
	return true
}

func noAirportRevisited(routes []Route) bool {
	for _, route := range routes {
		seen := map[string]bool{}
		for _, airport := range route.Airports {
			if seen[airport] {
				return false
			}
			seen[airport] = true
		}
	}
	return true
}

func routesSorted(data Dataset, routes []Route) bool {
	for i := 1; i < len(routes); i++ {
		if routeLabel(data, routes[i-1]) > routeLabel(data, routes[i]) {
			return false
		}
	}
	return true
}

func sortRoutes(data Dataset, routes []Route) {
	sort.Slice(routes, func(i, j int) bool {
		return routeLabel(data, routes[i]) < routeLabel(data, routes[j])
	})
}

func sortTermsByLabel(data Dataset, terms []string) {
	sort.Slice(terms, func(i, j int) bool {
		left := data.label(terms[i])
		right := data.label(terms[j])
		if left == right {
			return terms[i] < terms[j]
		}
		return left < right
	})
}

func routeDistribution(routes []Route) (direct, oneStop, twoStop int) {
	for _, route := range routes {
		switch route.Stopovers() {
		case 0:
			direct++
		case 1:
			oneStop++
		case 2:
			twoStop++
		}
	}
	return direct, oneStop, twoStop
}

func expandedAirports(source string, adj map[string][]string, maxDepth int) []string {
	seen := map[string]bool{}
	ordered := []string{}
	var walk func(string, []string)
	walk = func(current string, path []string) {
		depth := len(path) - 1
		if depth > maxDepth {
			return
		}
		if !seen[current] {
			seen[current] = true
			ordered = append(ordered, current)
		}
		if depth == maxDepth {
			return
		}
		for _, next := range adj[current] {
			if contains(path, next) {
				continue
			}
			walk(next, append(append([]string{}, path...), next))
		}
	}
	walk(source, []string{source})
	return ordered
}

func countAirportsInEdges(edges []Edge) int {
	seen := map[string]bool{}
	for _, edge := range edges {
		seen[edge.From] = true
		seen[edge.To] = true
	}
	return len(seen)
}

func countChecks(checks Checks) (passed, total int) {
	values := []bool{
		checks.SourceAndDestinationKnown,
		checks.FirstHopMatchesN3Facts,
		checks.RouteSetMatchesN3Query,
		checks.NoShorterRouteExists,
		checks.RoutesWithinStopoverLimit,
		checks.EveryHopHasFact,
		checks.NoAirportRevisited,
		checks.FullSourceGraphLoaded,
		checks.RoutesSortedDeterministic,
	}
	for _, ok := range values {
		if ok {
			passed++
		}
	}
	return passed, len(values)
}

func contains(path []string, term string) bool {
	for _, item := range path {
		if item == term {
			return true
		}
	}
	return false
}

func (data Dataset) label(term string) string {
	if label, ok := data.Labels[term]; ok {
		return label
	}
	return term
}

func (route Route) Stopovers() int {
	if len(route.Airports) < 2 {
		return 0
	}
	return len(route.Airports) - 2
}

func (route Route) Hops() int {
	if len(route.Airports) == 0 {
		return 0
	}
	return len(route.Airports) - 1
}

func routeLabel(data Dataset, route Route) string {
	labels := make([]string, 0, len(route.Airports))
	for _, airport := range route.Airports {
		labels = append(labels, data.label(airport))
	}
	return strings.Join(labels, " -> ")
}

func routeTerms(route Route) string {
	return strings.Join(route.Airports, " -> ")
}

func status(ok bool) string {
	if ok {
		return "OK"
	}
	return "FAIL"
}

func main() {
	data := exampleinput.Load(eyelingoExampleName, Dataset{})
	result := infer(data)
	checksPassed, checksTotal := countChecks(result.Checks)

	fmt.Println("# Path Discovery")

	fmt.Println()

	fmt.Println("## Answer")
	fmt.Printf("The path discovery query finds %d air routes with at most %d stopovers.\n", len(result.Routes), maxStopovers)
	fmt.Printf("from : %s\n", data.label(data.SourceID))
	fmt.Printf("to : %s\n", data.label(data.DestinationID))
	fmt.Printf("max stopovers : %d\n", maxStopovers)
	fmt.Println()
	fmt.Println("Discovered routes:")
	for i, route := range result.Routes {
		fmt.Printf("route %d (%d stopovers): %s\n", i+1, route.Stopovers(), routeLabel(data, route))
	}

	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The N3 source defines a recursive :route relation over nepo:hasOutboundRouteTo facts. A route can use a direct edge when the current length is within the maximum, or extend through a non-visited intermediate airport and recurse with length+1. The final log:collectAllIn query collects the labels of each airport in every route from the source to the destination.")
	fmt.Printf("source N3 airport labels : %d\n", sourceGraphAirportLabels)
	fmt.Printf("source N3 outbound-route facts : %d\n", sourceGraphOutboundFacts)
	fmt.Printf("translated full airport labels : %d\n", result.LoadedAirportLabels)
	fmt.Printf("translated full outbound-route facts : %d\n", result.LoadedOutboundFacts)
	fmt.Printf("airport terms appearing in outbound facts : %d\n", result.EdgeEndpointAirports)
	fmt.Printf("frontier airports expanded : %d\n", len(result.ExpandedAirports))
	fmt.Printf("bounded search outbound facts touched : %d\n", result.Stats.EdgeTests)
	fmt.Printf("source outbound candidates : %d\n", len(result.SourceOut))
	fmt.Printf("Liège outbound candidates : %d\n", len(result.FirstHopOut))
	fmt.Printf("direct routes : %d\n", result.DirectRoutes)
	fmt.Printf("one-stop routes : %d\n", result.OneStopRoutes)
	fmt.Printf("two-stopover routes : %d\n", result.TwoStopRoutes)
	fmt.Printf("search recursive calls : %d\n", result.Stats.RecursiveCalls)
	fmt.Printf("search edge tests : %d\n", result.Stats.EdgeTests)
	fmt.Printf("search depth-limit leaves : %d\n", result.Stats.DepthLimitLeaves)
	fmt.Println("Second-hop candidates from Liège:")
	for _, airport := range result.FirstHopOut {
		fmt.Printf("%s (%s)\n", data.label(airport), airport)
	}

	fmt.Println()
	fmt.Println("## Check")
	fmt.Printf("C1 %s - source and destination airport labels are known.\n", status(result.Checks.SourceAndDestinationKnown))
	fmt.Printf("C2 %s - Ostend-Bruges has one outbound route in the full N3 graph, to Liège Airport.\n", status(result.Checks.FirstHopMatchesN3Facts))
	fmt.Printf("C3 %s - the discovered route set matches the N3 query answer.\n", status(result.Checks.RouteSetMatchesN3Query))
	fmt.Printf("C4 %s - no direct or one-stop route exists under the same bound.\n", status(result.Checks.NoShorterRouteExists))
	fmt.Printf("C5 %s - every discovered route has at most two stopovers.\n", status(result.Checks.RoutesWithinStopoverLimit))
	fmt.Printf("C6 %s - every hop is backed by a nepo:hasOutboundRouteTo fact.\n", status(result.Checks.EveryHopHasFact))
	fmt.Printf("C7 %s - no route revisits an airport.\n", status(result.Checks.NoAirportRevisited))
	fmt.Printf("C8 %s - the Go translation loaded every airport label and outbound-route fact from the N3 source.\n", status(result.Checks.FullSourceGraphLoaded))
	fmt.Printf("C9 %s - route output is sorted deterministically by airport labels.\n", status(result.Checks.RoutesSortedDeterministic))

	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("question : %s\n", data.Question)
	fmt.Printf("source airport : %s (%s)\n", data.label(data.SourceID), data.SourceID)
	fmt.Printf("destination airport : %s (%s)\n", data.label(data.DestinationID), data.DestinationID)
	fmt.Printf("source graph airport labels : %d\n", sourceGraphAirportLabels)
	fmt.Printf("source graph outbound facts : %d\n", sourceGraphOutboundFacts)
	fmt.Printf("translated full airport labels : %d\n", result.LoadedAirportLabels)
	fmt.Printf("translated full outbound-route facts : %d\n", result.LoadedOutboundFacts)
	fmt.Printf("airport terms appearing in outbound facts : %d\n", result.EdgeEndpointAirports)
	fmt.Printf("bounded search outbound facts touched : %d\n", result.Stats.EdgeTests)
	fmt.Printf("max stopovers : %d\n", maxStopovers)
	fmt.Printf("max hops : %d\n", maxHops)
	fmt.Printf("routes discovered : %d\n", len(result.Routes))
	fmt.Printf("mandatory first hop : %s (%s)\n", data.label(mandatoryFirstHop), mandatoryFirstHop)
	fmt.Println("expanded airports:")
	for _, airport := range result.ExpandedAirports {
		fmt.Printf("%s (%s)\n", data.label(airport), airport)
	}
	for i, route := range result.Routes {
		fmt.Printf("route %d terms : %s\n", i+1, routeTerms(route))
		fmt.Printf("route %d labels : %s\n", i+1, routeLabel(data, route))
		fmt.Printf("route %d hops : %d\n", i+1, route.Hops())
		fmt.Printf("route %d stopovers : %d\n", i+1, route.Stopovers())
	}
	fmt.Printf("search recursive calls : %d\n", result.Stats.RecursiveCalls)
	fmt.Printf("search edge tests : %d\n", result.Stats.EdgeTests)
	fmt.Printf("search edges extended : %d\n", result.Stats.EdgesExtended)
	fmt.Printf("search revisit prunes : %d\n", result.Stats.RevisitPrunes)
	fmt.Printf("search depth-limit leaves : %d\n", result.Stats.DepthLimitLeaves)
	fmt.Printf("search dead ends : %d\n", result.Stats.DeadEnds)
	fmt.Printf("search routes emitted : %d\n", result.Stats.RoutesEmitted)
	fmt.Printf("search max depth : %d\n", result.Stats.MaxDepth)
	fmt.Printf("checks passed : %d/%d\n", checksPassed, checksTotal)
	fmt.Printf("all checks pass : %s\n", map[bool]string{true: "yes", false: "no"}[checksPassed == checksTotal])

	if checksPassed != checksTotal {
		os.Exit(1)
	}
}
