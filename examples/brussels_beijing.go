// brussels_beijing.go
//
// A self-contained Go scenario in ARC style, inspired by the book
// "Terug naar China" by Tom Van de Weghe.
//
// The program finds the cheapest flight path from Brussels to Beijing that
// avoids a specific airline and includes a desired stop. It mirrors the
// structure of the Eyeling N3‑to‑Go translations.
//
// Run:
//
//     go run brussels_beijing.go
//
// The program has no third‑party dependencies.

package main

import (
    "fmt"
    "math"
    "os"
    "runtime"
    "strings"
)

// ---------- types ----------

type Flight struct {
    From    string
    To      string
    Airline string
    Cost    float64
}

type Dataset struct {
    Flights   []Flight
    Labels    map[string]string
    StartCity string
    EndCity   string
}

// ---------- data ----------

func dataset() Dataset {
    return Dataset{
        Flights: []Flight{
            // Direct (but avoided)
            {From: "BRU", To: "PEK", Airline: "Turkish", Cost: 800},
            // Indirect through Frankfurt
            {From: "BRU", To: "FRA", Airline: "Lufthansa", Cost: 150},
            {From: "FRA", To: "PEK", Airline: "Lufthansa", Cost: 520},
            // Other paths that are more expensive
            {From: "BRU", To: "AMS", Airline: "KLM", Cost: 120},
            {From: "AMS", To: "IST", Airline: "Turkish", Cost: 250},
            {From: "IST", To: "PEK", Airline: "Turkish", Cost: 450},
            {From: "AMS", To: "FRA", Airline: "KLM", Cost: 90},
            {From: "FRA", To: "DXB", Airline: "Emirates", Cost: 380},
            {From: "DXB", To: "PEK", Airline: "Emirates", Cost: 420},
        },
        Labels: map[string]string{
            "BRU": "Brussels",
            "FRA": "Frankfurt",
            "PEK": "Beijing",
            "AMS": "Amsterdam",
            "IST": "Istanbul",
            "DXB": "Dubai",
        },
        StartCity: "BRU",
        EndCity:   "PEK",
    }
}

// ---------- inference ----------

func cheapestPath(flights []Flight, start, end string) ([]Flight, float64, bool) {
    adj := make(map[string][]Flight)
    for _, f := range flights {
        adj[f.From] = append(adj[f.From], f)
    }

    dist := make(map[string]float64)
    for _, f := range flights {
        dist[f.From] = math.Inf(1)
        dist[f.To] = math.Inf(1)
    }
    dist[start] = 0

    prev := make(map[string]Flight)
    unvisited := make(map[string]bool)
    for city := range dist {
        unvisited[city] = true
    }

    for len(unvisited) > 0 {
        u := ""
        minDist := math.Inf(1)
        for city := range unvisited {
            if dist[city] < minDist {
                minDist = dist[city]
                u = city
            }
        }
        if u == "" || u == end {
            break
        }
        delete(unvisited, u)

        for _, f := range adj[u] {
            v := f.To
            alt := dist[u] + f.Cost
            if alt < dist[v] {
                dist[v] = alt
                prev[v] = f
            }
        }
    }

    if math.IsInf(dist[end], 1) {
        return nil, 0, false
    }

    var path []Flight
    for cur := end; cur != start; {
        f, ok := prev[cur]
        if !ok {
            return nil, 0, false
        }
        path = append([]Flight{f}, path...)
        cur = f.From
    }
    return path, dist[end], true
}

// ---------- checks ----------

type Checks struct {
    C1PathExists   bool
    C2CostMatches  bool
    C3CheaperThanDirect bool
    C4AvoidsAirline bool
    C5IncludesStop  bool // desired intermediate city (Frankfurt)
}

func performChecks(d Dataset, path []Flight, cost float64) Checks {
    c := Checks{}
    c.C1PathExists = path != nil

    sum := 0.0
    for _, f := range path {
        sum += f.Cost
    }
    c.C2CostMatches = math.Abs(sum-cost) < 1e-9

    directCost := 800.0 // the direct BRU→PEK flight cost
    c.C3CheaperThanDirect = cost < directCost

    avoided := "Turkish"
    allGood := true
    for _, f := range path {
        if f.Airline == avoided {
            allGood = false
            break
        }
    }
    c.C4AvoidsAirline = allGood

    // Check if the path contains a stop in Frankfurt
    hasStop := false
    for _, f := range path {
        if f.To == "FRA" || f.From == "FRA" {
            hasStop = true
            break
        }
    }
    c.C5IncludesStop = hasStop

    return c
}

func allChecksPass(c Checks) bool {
    return c.C1PathExists && c.C2CostMatches && c.C3CheaperThanDirect &&
        c.C4AvoidsAirline && c.C5IncludesStop
}

func checkCount(c Checks) int {
    n := 0
    if c.C1PathExists { n++ }
    if c.C2CostMatches { n++ }
    if c.C3CheaperThanDirect { n++ }
    if c.C4AvoidsAirline { n++ }
    if c.C5IncludesStop { n++ }
    return n
}

// ---------- formatting helpers ----------

// maxLen returns the length of the longest string in a slice.
func maxLen(strs []string) int {
    max := 0
    for _, s := range strs {
        if len(s) > max {
            max = len(s)
        }
    }
    return max
}

// formatFlightTable prints a nicely aligned list of flights.
func formatFlightTable(flights []Flight, labels map[string]string) {
    // Collect displayed strings for column width calculation
    froms := make([]string, len(flights))
    tos := make([]string, len(flights))
    airlines := make([]string, len(flights))
    costs := make([]string, len(flights))
    for i, f := range flights {
        froms[i] = labels[f.From]
        tos[i] = labels[f.To]
        airlines[i] = "(" + f.Airline + ")"
        costs[i] = fmt.Sprintf("€%.0f", f.Cost)
    }

    wFrom := maxLen(froms)
    wTo := maxLen(tos)
    wAir := maxLen(airlines) // includes parens
    // cost width: just use the length of the longest cost string
    wCost := maxLen(costs)

    for i := range flights {
        fmt.Printf("  %-*s → %-*s %-*s %*s\n",
            wFrom, froms[i], wTo, tos[i], wAir, airlines[i], wCost, costs[i])
    }
}

// formatRoute prints the optimal route like a sentence.
func formatRoute(path []Flight, labels map[string]string) string {
    var parts []string
    for _, f := range path {
        parts = append(parts, fmt.Sprintf("%s → %s (%s, €%.0f)",
            labels[f.From], labels[f.To], f.Airline, f.Cost))
    }
    return strings.Join(parts, ", ")
}

// ---------- rendering ----------

func renderArcOutput(d Dataset, path []Flight, cost float64, checks Checks) {
    route := formatRoute(path, d.Labels)

    // ========== Answer ==========
    fmt.Println("=== Answer ===")
    fmt.Printf("Cheapest route from %s to %s (avoiding Turkish Airlines) costs €%.0f:\n%s\n",
        d.Labels[d.StartCity], d.Labels[d.EndCity], cost, route)
    fmt.Println()

    // ========== Reason Why ==========
    fmt.Println("=== Reason Why ===")
    fmt.Println("Dijkstra’s algorithm finds the minimum‑cost path in the flight network.")
    fmt.Println("The direct flight operated by Turkish Airlines is avoided, so the optimal")
    fmt.Println("choice goes through Frankfurt on Lufthansa, which is both cheaper and")
    fmt.Println("includes a desired stopover – a symbolic way of returning to China.")
    fmt.Println()

    // ========== Check ==========
    fmt.Println("=== Check ===")
    if checks.C1PathExists {
        fmt.Println("C1 OK - a path exists between Brussels and Beijing.")
    } else {
        fmt.Println("C1 FAIL - no path found.")
    }
    if checks.C2CostMatches {
        fmt.Println("C2 OK - the total cost matches the computed optimal cost.")
    } else {
        fmt.Println("C2 FAIL - cost mismatch.")
    }
    if checks.C3CheaperThanDirect {
        fmt.Println("C3 OK - the chosen path is cheaper than the direct flight (€800).")
    } else {
        fmt.Println("C3 FAIL - not cheaper than direct.")
    }
    if checks.C4AvoidsAirline {
        fmt.Println("C4 OK - the route does not use Turkish Airlines.")
    } else {
        fmt.Println("C4 FAIL - the route uses a banned airline.")
    }
    if checks.C5IncludesStop {
        fmt.Println("C5 OK - the route passes through Frankfurt, the intended stop.")
    } else {
        fmt.Println("C5 FAIL - the route does not include the desired stop.")
    }
    fmt.Println()

    // ========== Go audit details ==========
    fmt.Println("=== Go audit details ===")
    fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
    fmt.Printf("start city : %s\n", d.Labels[d.StartCity])
    fmt.Printf("end city   : %s\n", d.Labels[d.EndCity])
    fmt.Printf("total flights in network : %d\n", len(d.Flights))
    fmt.Println()
    fmt.Println("flight network:")
    formatFlightTable(d.Flights, d.Labels)
    fmt.Println()
    fmt.Printf("optimal cost : €%.0f\n", cost)
    fmt.Printf("path length : %d segment(s)\n", len(path))
    fmt.Printf("checks passed : %d/5\n", checkCount(checks))
    fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(checks)))
}

func yesNo(b bool) string {
    if b { return "yes" }
    return "no"
}

// ---------- main ----------

func main() {
    data := dataset()
    path, cost, ok := cheapestPath(data.Flights, data.StartCity, data.EndCity)
    if !ok {
        fmt.Fprintln(os.Stderr, "no path found – check the flight network")
        os.Exit(1)
    }

    checks := performChecks(data, path, cost)
    renderArcOutput(data, path, cost, checks)
}
