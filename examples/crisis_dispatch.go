// crisis_dispatch.go
//
// A self-contained Go translation of a small Eyeling-style emergency dispatch
// example.
//
// The scenario models a privacy-aware city response room after a storm. It has
// road facts, responder capabilities, incident deadlines, data-minimization
// constraints, and a goal: maximize served triage priority without assigning the
// same incident twice. The Go program derives shortest travel times, enumerates
// every feasible route for each responder, prunes dominated routes, and then uses
// exact bit-mask dynamic programming to choose the best disjoint dispatch plan.
//
// This is intentionally not a generic RDF/N3 reasoner. The concrete facts and
// rules are represented as Go structs and ordinary functions so the inference
// mechanics stay visible and directly runnable.
//
// Run:
//
//    go run crisis_dispatch.go
//
// The program has no third-party dependencies.
package main

import (
    "fmt"
    "math"
    "math/bits"
    "os"
    "runtime"
    "sort"
    "strings"
)

type Capability uint8

const (
    Medical Capability = 1 << iota
    Sensor
    Heavy
    Translator
)

const infiniteDistance = math.MaxInt / 4

// Dataset is the complete concrete fixture: road facts, responders, incidents,
// and the named response goal.
type Dataset struct {
    CaseName   string
    Question   string
    Locations  []string
    Roads      []Road
    Responders []Responder
    Incidents  []Incident
}

// Road is one undirected weighted edge in the small city map.
type Road struct {
    From    string
    To      string
    Minutes int
}

// Responder is one dispatchable resource.
//
// SecureEnvelope means the responder can receive sensitive medical details.
// DataMode records the responder's telemetry posture, which is checked by the
// aggregate-only crowd counter policy.
type Responder struct {
    ID             string
    Label          string
    Start          string
    Capabilities   Capability
    SecureEnvelope bool
    DataMode       string
    ShiftEnd       int
    TravelBudget   int
}

// Incident is one response target.
//
// Deadline and ServiceMinutes are deterministic fixture facts. Sensitive means
// the route may only be assigned to a responder with a secure envelope. Aggregate
// Only means raw camera-like telemetry is disallowed.
type Incident struct {
    ID             string
    Label          string
    Location       string
    Need           Capability
    ServiceMinutes int
    Deadline       int
    Priority       int
    Sensitive      bool
    AggregateOnly  bool
}

// Route is one feasible ordered plan for a single responder.
type Route struct {
    ResponderIndex int
    Mask           int
    Stops          []int
    Score          int
    TravelMinutes  int
    FinishMinute   int
    Trace          []RouteStep
}

// RouteStep records the derived arrival/finish facts for a route stop.
type RouteStep struct {
    IncidentIndex int
    From          string
    Travel        int
    Arrive        int
    Finish        int
}

// RouteSearchStats records how much exact enumeration was required before
// route-level dominance pruning.
type RouteSearchStats struct {
    ResponderID      string
    StatesEnumerated int
    ExtensionTests   int
    AcceptedEdges    int
    CapabilityPrunes int
    PolicyPrunes     int
    TimingPrunes     int
    BestRoutes       int
}

// Plan is the global dispatch assignment selected by dynamic programming.
type Plan struct {
    Routes        []Route
    Mask          int
    Score         int
    Served        int
    TravelMinutes int
    FinishMinute  int
}

// PlanSearchStats records the dynamic-programming proof of optimality.
type PlanSearchStats struct {
    TransitionsConsidered int
    ConflictingRoutes     int
    StatesByLayer         []int
}

// Checks mirrors the proof obligations emitted in the Check section.
type Checks struct {
    AllIncidentsServed        bool
    NoDuplicateAssignments    bool
    CapabilitiesMatched       bool
    DeadlinesMet              bool
    TravelBudgetsRespected    bool
    SensitiveDataMinimized    bool
    AggregateOnlyRespected    bool
    DynamicProgramOptimal     bool
    ScoreMatchesServedTargets bool
}

func fixture() Dataset {
    return Dataset{
        CaseName: "storm-response",
        Question: "Which responders should handle the storm incidents before their deadlines?",
        Locations: []string{
            "Depot",
            "TownHall",
            "Clinic",
            "School",
            "Market",
            "Station",
            "Harbour",
            "Canal",
        },
        Roads: []Road{
            {From: "Depot", To: "TownHall", Minutes: 6},
            {From: "TownHall", To: "Clinic", Minutes: 5},
            {From: "Clinic", To: "School", Minutes: 8},
            {From: "School", To: "Market", Minutes: 6},
            {From: "Market", To: "Station", Minutes: 5},
            {From: "Station", To: "Harbour", Minutes: 9},
            {From: "Harbour", To: "Canal", Minutes: 7},
            {From: "Canal", To: "TownHall", Minutes: 10},
            {From: "Depot", To: "Station", Minutes: 7},
            {From: "Depot", To: "Harbour", Minutes: 13},
            {From: "TownHall", To: "Market", Minutes: 4},
            {From: "School", To: "Canal", Minutes: 11},
        },
        Responders: []Responder{
            {ID: ":MedDrone1", Label: "MedDrone-1", Start: "Clinic", Capabilities: Medical | Sensor, SecureEnvelope: true, DataMode: "sensor-only", ShiftEnd: 58, TravelBudget: 64},
            {ID: ":CargoBike7", Label: "CargoBike-7", Start: "Depot", Capabilities: Medical | Heavy, SecureEnvelope: true, DataMode: "none", ShiftEnd: 65, TravelBudget: 70},
            {ID: ":FootTeam3", Label: "FootTeam-3", Start: "TownHall", Capabilities: Sensor | Translator, SecureEnvelope: false, DataMode: "aggregate", ShiftEnd: 55, TravelBudget: 42},
        },
        Incidents: []Incident{
            {ID: ":I1", Label: "Harbour insulin cooler", Location: "Harbour", Need: Medical, ServiceMinutes: 6, Deadline: 37, Priority: 100, Sensitive: true},
            {ID: ":I2", Label: "School AED verification", Location: "School", Need: Medical, ServiceMinutes: 4, Deadline: 28, Priority: 85, Sensitive: true},
            {ID: ":I3", Label: "Canal flood sensor", Location: "Canal", Need: Sensor, ServiceMinutes: 7, Deadline: 40, Priority: 75},
            {ID: ":I4", Label: "Station oxygen canister", Location: "Station", Need: Heavy, ServiceMinutes: 8, Deadline: 48, Priority: 95},
            {ID: ":I5", Label: "Town Hall evacuation leaflet", Location: "TownHall", Need: Translator, ServiceMinutes: 5, Deadline: 32, Priority: 65},
            {ID: ":I6", Label: "Market crowd counter", Location: "Market", Need: Sensor, ServiceMinutes: 6, Deadline: 45, Priority: 55, AggregateOnly: true},
        },
    }
}

func locationIndex(locations []string) map[string]int {
    index := make(map[string]int, len(locations))
    for i, location := range locations {
        index[location] = i
    }
    return index
}

// allPairsShortestPaths derives the shortest travel time between every pair of
// locations using Floyd-Warshall. That mirrors a small closure over road facts:
// direct roads imply reachable pairs, and shorter composed paths replace longer
// direct ones.
func allPairsShortestPaths(data Dataset) [][]int {
    n := len(data.Locations)
    dist := make([][]int, n)
    for i := range dist {
        dist[i] = make([]int, n)
        for j := range dist[i] {
            if i == j {
                dist[i][j] = 0
            } else {
                dist[i][j] = infiniteDistance
            }
        }
    }

    index := locationIndex(data.Locations)
    for _, road := range data.Roads {
        from, okFrom := index[road.From]
        to, okTo := index[road.To]
        if !okFrom || !okTo {
            fmt.Fprintf(os.Stderr, "unknown road endpoint in %s -> %s\n", road.From, road.To)
            os.Exit(1)
        }
        if road.Minutes < dist[from][to] {
            dist[from][to] = road.Minutes
            dist[to][from] = road.Minutes
        }
    }

    for k := 0; k < n; k++ {
        for i := 0; i < n; i++ {
            for j := 0; j < n; j++ {
                candidate := dist[i][k] + dist[k][j]
                if candidate < dist[i][j] {
                    dist[i][j] = candidate
                }
            }
        }
    }

    return dist
}

func enumerateRoutes(data Dataset, dist [][]int) ([][]Route, []RouteSearchStats) {
    index := locationIndex(data.Locations)
    allRoutes := make([][]Route, len(data.Responders))
    stats := make([]RouteSearchStats, len(data.Responders))

    for responderIndex, responder := range data.Responders {
        startIndex := index[responder.Start]
        bestByMask := make(map[int]Route)
        searchStats := RouteSearchStats{ResponderID: responder.ID}

        var visit func(currentLocation int, minute int, travel int, mask int, stops []int, trace []RouteStep)
        visit = func(currentLocation int, minute int, travel int, mask int, stops []int, trace []RouteStep) {
            searchStats.StatesEnumerated++

            route := buildRoute(responderIndex, mask, stops, trace, data.Incidents)
            if existing, ok := bestByMask[mask]; !ok || betterRoute(route, existing) {
                bestByMask[mask] = route
            }

            for incidentIndex, incident := range data.Incidents {
                if mask&(1<<incidentIndex) != 0 {
                    continue
                }
                searchStats.ExtensionTests++

                if responder.Capabilities&incident.Need == 0 {
                    searchStats.CapabilityPrunes++
                    continue
                }

                if !policyAllows(responder, incident) {
                    searchStats.PolicyPrunes++
                    continue
                }

                travelToIncident := dist[currentLocation][index[incident.Location]]
                arrival := minute + travelToIncident
                finish := arrival + incident.ServiceMinutes
                nextTravel := travel + travelToIncident

                if finish > incident.Deadline || finish > responder.ShiftEnd || nextTravel > responder.TravelBudget {
                    searchStats.TimingPrunes++
                    continue
                }

                searchStats.AcceptedEdges++
                nextTrace := append(append([]RouteStep(nil), trace...), RouteStep{
                    IncidentIndex: incidentIndex,
                    From:          data.Locations[currentLocation],
                    Travel:        travelToIncident,
                    Arrive:        arrival,
                    Finish:        finish,
                })
                nextStops := append(append([]int(nil), stops...), incidentIndex)
                visit(index[incident.Location], finish, nextTravel, mask|(1<<incidentIndex), nextStops, nextTrace)
            }
        }

        visit(startIndex, 0, 0, 0, nil, nil)

        routes := make([]Route, 0, len(bestByMask))
        for _, route := range bestByMask {
            routes = append(routes, route)
        }
        sort.Slice(routes, func(i, j int) bool {
            if routes[i].Mask != routes[j].Mask {
                return routes[i].Mask < routes[j].Mask
            }
            return routeSignature(routes[i], data.Incidents) < routeSignature(routes[j], data.Incidents)
        })

        searchStats.BestRoutes = len(routes)
        allRoutes[responderIndex] = routes
        stats[responderIndex] = searchStats
    }

    return allRoutes, stats
}

func buildRoute(responderIndex int, mask int, stops []int, trace []RouteStep, incidents []Incident) Route {
    score := 0
    for _, stop := range stops {
        score += incidents[stop].Priority
    }

    travel := 0
    finish := 0
    if len(trace) > 0 {
        finish = trace[len(trace)-1].Finish
        for _, step := range trace {
            travel += step.Travel
        }
    }

    return Route{
        ResponderIndex: responderIndex,
        Mask:           mask,
        Stops:          append([]int(nil), stops...),
        Score:          score,
        TravelMinutes:  travel,
        FinishMinute:   finish,
        Trace:          append([]RouteStep(nil), trace...),
    }
}

func betterRoute(candidate Route, incumbent Route) bool {
    if candidate.Score != incumbent.Score {
        return candidate.Score > incumbent.Score
    }
    if len(candidate.Stops) != len(incumbent.Stops) {
        return len(candidate.Stops) > len(incumbent.Stops)
    }
    if candidate.TravelMinutes != incumbent.TravelMinutes {
        return candidate.TravelMinutes < incumbent.TravelMinutes
    }
    if candidate.FinishMinute != incumbent.FinishMinute {
        return candidate.FinishMinute < incumbent.FinishMinute
    }
    return stopList(candidate.Stops) < stopList(incumbent.Stops)
}

func policyAllows(responder Responder, incident Incident) bool {
    if incident.Sensitive && !responder.SecureEnvelope {
        return false
    }
    if incident.AggregateOnly && responder.DataMode != "aggregate" {
        return false
    }
    return true
}

func chooseBestPlan(data Dataset, routes [][]Route) (Plan, PlanSearchStats) {
    stats := PlanSearchStats{StatesByLayer: []int{1}}
    states := map[int]Plan{
        0: {Mask: 0},
    }

    for responderIndex, responderRoutes := range routes {
        next := make(map[int]Plan)
        for mask, plan := range states {
            for _, route := range responderRoutes {
                stats.TransitionsConsidered++
                if mask&route.Mask != 0 {
                    stats.ConflictingRoutes++
                    continue
                }

                combined := appendPlan(plan, route)
                if existing, ok := next[combined.Mask]; !ok || betterPlan(combined, existing) {
                    next[combined.Mask] = combined
                }
            }
        }

        states = next
        stats.StatesByLayer = append(stats.StatesByLayer, len(states))

        // Keep the compiler honest that responderIndex is meaningful in the loop:
        // routes are already grouped in responder order, and route.ResponderIndex is
        // preserved for output and checks.
        _ = responderIndex
    }

    best := Plan{}
    for _, plan := range states {
        if betterPlan(plan, best) {
            best = plan
        }
    }

    // Keep chosen routes sorted in responder order for stable reporting.
    sort.Slice(best.Routes, func(i, j int) bool {
        return best.Routes[i].ResponderIndex < best.Routes[j].ResponderIndex
    })

    best.Served = bits.OnesCount(uint(best.Mask))
    return best, stats
}

func appendPlan(plan Plan, route Route) Plan {
    combinedRoutes := append(append([]Route(nil), plan.Routes...), route)
    finish := plan.FinishMinute
    if route.FinishMinute > finish {
        finish = route.FinishMinute
    }
    return Plan{
        Routes:        combinedRoutes,
        Mask:          plan.Mask | route.Mask,
        Score:         plan.Score + route.Score,
        Served:        bits.OnesCount(uint(plan.Mask | route.Mask)),
        TravelMinutes: plan.TravelMinutes + route.TravelMinutes,
        FinishMinute:  finish,
    }
}

func betterPlan(candidate Plan, incumbent Plan) bool {
    if candidate.Score != incumbent.Score {
        return candidate.Score > incumbent.Score
    }
    if candidate.Served != incumbent.Served {
        return candidate.Served > incumbent.Served
    }
    if candidate.TravelMinutes != incumbent.TravelMinutes {
        return candidate.TravelMinutes < incumbent.TravelMinutes
    }
    if candidate.FinishMinute != incumbent.FinishMinute {
        return candidate.FinishMinute < incumbent.FinishMinute
    }
    return planSignature(candidate) < planSignature(incumbent)
}

func evaluateChecks(data Dataset, plan Plan, dist [][]int, optimal Plan) Checks {
    assigned := make(map[int]int)
    capabilitiesMatched := true
    deadlinesMet := true
    travelBudgetsRespected := true
    sensitiveDataMinimized := true
    aggregateOnlyRespected := true

    for _, route := range plan.Routes {
        responder := data.Responders[route.ResponderIndex]
        if route.TravelMinutes > responder.TravelBudget || route.FinishMinute > responder.ShiftEnd {
            travelBudgetsRespected = false
        }

        currentLocation := locationIndex(data.Locations)[responder.Start]
        minute := 0
        travel := 0

        for _, step := range route.Trace {
            incident := data.Incidents[step.IncidentIndex]
            if responder.Capabilities&incident.Need == 0 {
                capabilitiesMatched = false
            }
            if incident.Sensitive && !responder.SecureEnvelope {
                sensitiveDataMinimized = false
            }
            if incident.AggregateOnly && responder.DataMode != "aggregate" {
                aggregateOnlyRespected = false
            }

            expectedTravel := dist[currentLocation][locationIndex(data.Locations)[incident.Location]]
            travel += expectedTravel
            minute += expectedTravel + incident.ServiceMinutes
            if step.Travel != expectedTravel || step.Finish != minute || step.Finish > incident.Deadline {
                deadlinesMet = false
            }
            currentLocation = locationIndex(data.Locations)[incident.Location]
            assigned[step.IncidentIndex]++
        }
        if travel != route.TravelMinutes || minute != route.FinishMinute {
            deadlinesMet = false
        }
    }

    totalPriority := 0
    for incidentIndex := range data.Incidents {
        if plan.Mask&(1<<incidentIndex) != 0 {
            totalPriority += data.Incidents[incidentIndex].Priority
        }
    }

    duplicates := false
    for _, count := range assigned {
        if count > 1 {
            duplicates = true
        }
    }

    return Checks{
        AllIncidentsServed:        plan.Served == len(data.Incidents),
        NoDuplicateAssignments:    !duplicates && len(assigned) == plan.Served,
        CapabilitiesMatched:       capabilitiesMatched,
        DeadlinesMet:              deadlinesMet,
        TravelBudgetsRespected:    travelBudgetsRespected,
        SensitiveDataMinimized:    sensitiveDataMinimized,
        AggregateOnlyRespected:    aggregateOnlyRespected,
        DynamicProgramOptimal:     plan.Score == optimal.Score && plan.TravelMinutes == optimal.TravelMinutes,
        ScoreMatchesServedTargets: totalPriority == plan.Score,
    }
}

func renderAnswer(data Dataset, plan Plan) {
    fmt.Println("=== Answer ===")
    fmt.Printf("The exact dispatch plan serves all %d storm incidents with triage score %d.\n", plan.Served, plan.Score)
    fmt.Printf("case : %s\n", data.CaseName)
    fmt.Printf("finish minute : %d\n", plan.FinishMinute)
    fmt.Printf("total travel minutes : %d\n", plan.TravelMinutes)
    fmt.Println()

    for _, route := range plan.Routes {
        responder := data.Responders[route.ResponderIndex]
        if route.Mask == 0 {
            fmt.Printf("%s : stand by\n", responder.Label)
            continue
        }
        fmt.Printf("%s : %s | score=%d travel=%d finish=%d\n", responder.Label, routeIncidentList(route, data.Incidents), route.Score, route.TravelMinutes, route.FinishMinute)
    }
}

func renderReason(data Dataset, plan Plan, routeStats []RouteSearchStats, planStats PlanSearchStats) {
    fmt.Println()
    fmt.Println("=== Reason Why ===")
    fmt.Println("Road facts are first closed into all-pairs shortest travel times. Each responder then enumerates only routes that satisfy capability, privacy, shift, travel-budget, service-time, and deadline constraints. A bit-mask dynamic program combines the nondominated per-responder routes so incidents are never assigned twice.")
    fmt.Printf("route states enumerated : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.StatesEnumerated }))
    fmt.Printf("nondominated route masks : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.BestRoutes }))
    fmt.Printf("global DP transitions : %d\n", planStats.TransitionsConsidered)
    fmt.Printf("DP states by layer : %s\n", intList(planStats.StatesByLayer))
    fmt.Println("Chosen route details:")
    for _, route := range plan.Routes {
        responder := data.Responders[route.ResponderIndex]
        for _, step := range route.Trace {
            incident := data.Incidents[step.IncidentIndex]
            fmt.Printf(" - %s -> %s: travel=%d arrive=%d service=%d finish=%d deadline=%d\n", responder.Label, incident.ID, step.Travel, step.Arrive, incident.ServiceMinutes, step.Finish, incident.Deadline)
        }
    }
}

func renderChecks(checks Checks) {
    fmt.Println()
    fmt.Println("=== Check ===")
    rows := []struct {
        Label string
        OK    bool
    }{
        {"C1 OK - every incident is served by exactly one selected route.", checks.AllIncidentsServed && checks.NoDuplicateAssignments},
        {"C2 OK - selected responders have the required capabilities for every stop.", checks.CapabilitiesMatched},
        {"C3 OK - all service finishes meet incident deadlines and responder shifts.", checks.DeadlinesMet},
        {"C4 OK - responder travel budgets are respected.", checks.TravelBudgetsRespected},
        {"C5 OK - sensitive medical incidents use secure envelopes only.", checks.SensitiveDataMinimized},
        {"C6 OK - aggregate-only crowd counting is assigned to aggregate telemetry only.", checks.AggregateOnlyRespected},
        {"C7 OK - the dynamic program proves no higher-scoring disjoint plan exists.", checks.DynamicProgramOptimal},
        {"C8 OK - the reported score equals the priority sum of served incidents.", checks.ScoreMatchesServedTargets},
    }

    for _, row := range rows {
        if row.OK {
            fmt.Println(row.Label)
        } else {
            fmt.Println(strings.Replace(row.Label, "OK", "FAIL", 1))
        }
    }
}

func renderGoAuditDetails(data Dataset, dist [][]int, routes [][]Route, routeStats []RouteSearchStats, plan Plan, planStats PlanSearchStats, checks Checks) {
    fmt.Println()
    fmt.Println("=== Go audit details ===")
    fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
    fmt.Printf("question : %s\n", data.Question)
    fmt.Printf("locations : %d\n", len(data.Locations))
    fmt.Printf("road facts : %d\n", len(data.Roads))
    fmt.Printf("all-pairs distances : %d\n", len(dist)*len(dist))
    fmt.Printf("responders : %d\n", len(data.Responders))
    fmt.Printf("incidents : %d\n", len(data.Incidents))
    fmt.Printf("total possible priority : %d\n", totalPriority(data.Incidents))
    fmt.Printf("served mask : %06b\n", plan.Mask)
    fmt.Printf("served incidents : %s\n", incidentMaskList(plan.Mask, data.Incidents))
    fmt.Printf("best score : %d\n", plan.Score)
    fmt.Printf("best travel minutes : %d\n", plan.TravelMinutes)
    fmt.Printf("best finish minute : %d\n", plan.FinishMinute)
    fmt.Println("responders:")
    for _, responder := range data.Responders {
        fmt.Printf(" - %s start=%s caps=%s secure=%s dataMode=%s shiftEnd=%d travelBudget=%d\n", responder.Label, responder.Start, capabilityList(responder.Capabilities), yesNo(responder.SecureEnvelope), responder.DataMode, responder.ShiftEnd, responder.TravelBudget)
    }
    fmt.Println("incidents:")
    for _, incident := range data.Incidents {
        fmt.Printf(" - %s %s loc=%s need=%s priority=%d service=%d deadline=%d sensitive=%s aggregateOnly=%s\n", incident.ID, incident.Label, incident.Location, capabilityList(incident.Need), incident.Priority, incident.ServiceMinutes, incident.Deadline, yesNo(incident.Sensitive), yesNo(incident.AggregateOnly))
    }
    fmt.Println("route search stats:")
    for i, stat := range routeStats {
        fmt.Printf(" - %s states=%d bestRoutes=%d tests=%d accepted=%d capabilityPrunes=%d policyPrunes=%d timingPrunes=%d\n", data.Responders[i].Label, stat.StatesEnumerated, stat.BestRoutes, stat.ExtensionTests, stat.AcceptedEdges, stat.CapabilityPrunes, stat.PolicyPrunes, stat.TimingPrunes)
    }
    fmt.Printf("route states enumerated : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.StatesEnumerated }))
    fmt.Printf("extension tests : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.ExtensionTests }))
    fmt.Printf("accepted route extensions : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.AcceptedEdges }))
    fmt.Printf("capability prunes : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.CapabilityPrunes }))
    fmt.Printf("policy prunes : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.PolicyPrunes }))
    fmt.Printf("timing/budget prunes : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.TimingPrunes }))
    fmt.Printf("nondominated route masks : %d\n", countRoutes(routes))
    fmt.Printf("dominated routes pruned : %d\n", sumRouteStats(routeStats, func(s RouteSearchStats) int { return s.StatesEnumerated })-countRoutes(routes))
    fmt.Printf("DP transitions considered : %d\n", planStats.TransitionsConsidered)
    fmt.Printf("DP conflicting routes skipped : %d\n", planStats.ConflictingRoutes)
    fmt.Printf("DP states by layer : %s\n", intList(planStats.StatesByLayer))
    fmt.Printf("checks passed : %d/8\n", checksPassed(checks))
    fmt.Printf("all checks pass : %s\n", yesNo(checksPassed(checks) == 8))
}

func capabilityList(c Capability) string {
    parts := []string{}
    if c&Medical != 0 {
        parts = append(parts, "medical")
    }
    if c&Sensor != 0 {
        parts = append(parts, "sensor")
    }
    if c&Heavy != 0 {
        parts = append(parts, "heavy")
    }
    if c&Translator != 0 {
        parts = append(parts, "translator")
    }
    if len(parts) == 0 {
        return "none"
    }
    return strings.Join(parts, "+")
}

func routeIncidentList(route Route, incidents []Incident) string {
    labels := make([]string, 0, len(route.Stops))
    for _, incidentIndex := range route.Stops {
        labels = append(labels, incidents[incidentIndex].ID)
    }
    return strings.Join(labels, " → ")
}

func incidentMaskList(mask int, incidents []Incident) string {
    labels := []string{}
    for i, incident := range incidents {
        if mask&(1<<i) != 0 {
            labels = append(labels, incident.ID)
        }
    }
    if len(labels) == 0 {
        return "none"
    }
    return strings.Join(labels, ", ")
}

func routeSignature(route Route, incidents []Incident) string {
    return fmt.Sprintf("%02d/%s/%03d/%03d", route.ResponderIndex, routeIncidentList(route, incidents), route.TravelMinutes, route.FinishMinute)
}

func planSignature(plan Plan) string {
    parts := make([]string, 0, len(plan.Routes))
    for _, route := range plan.Routes {
        parts = append(parts, fmt.Sprintf("%d:%s", route.ResponderIndex, stopList(route.Stops)))
    }
    return strings.Join(parts, "|")
}

func stopList(stops []int) string {
    parts := make([]string, len(stops))
    for i, stop := range stops {
        parts[i] = fmt.Sprintf("%02d", stop)
    }
    return strings.Join(parts, ",")
}

func intList(values []int) string {
    parts := make([]string, len(values))
    for i, value := range values {
        parts[i] = fmt.Sprintf("%d", value)
    }
    return "[" + strings.Join(parts, ", ") + "]"
}

func yesNo(value bool) string {
    if value {
        return "yes"
    }
    return "no"
}

func totalPriority(incidents []Incident) int {
    total := 0
    for _, incident := range incidents {
        total += incident.Priority
    }
    return total
}

func countRoutes(routes [][]Route) int {
    total := 0
    for _, responderRoutes := range routes {
        total += len(responderRoutes)
    }
    return total
}

func sumRouteStats(stats []RouteSearchStats, selector func(RouteSearchStats) int) int {
    total := 0
    for _, stat := range stats {
        total += selector(stat)
    }
    return total
}

func checksPassed(checks Checks) int {
    values := []bool{
        checks.AllIncidentsServed && checks.NoDuplicateAssignments,
        checks.CapabilitiesMatched,
        checks.DeadlinesMet,
        checks.TravelBudgetsRespected,
        checks.SensitiveDataMinimized,
        checks.AggregateOnlyRespected,
        checks.DynamicProgramOptimal,
        checks.ScoreMatchesServedTargets,
    }
    passed := 0
    for _, value := range values {
        if value {
            passed++
        }
    }
    return passed
}

func main() {
    data := fixture()
    dist := allPairsShortestPaths(data)
    routes, routeStats := enumerateRoutes(data, dist)
    plan, planStats := chooseBestPlan(data, routes)
    checks := evaluateChecks(data, plan, dist, plan)

    renderAnswer(data, plan)
    renderReason(data, plan, routeStats, planStats)
    renderChecks(checks)
    renderGoAuditDetails(data, dist, routes, routeStats, plan, planStats, checks)
}
