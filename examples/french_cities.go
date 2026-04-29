// french_cities.go
//
// A self-contained Go translation of examples/french-cities.n3 from the Eyeling
// example suite, in ARC style.
//
// The original N3 program encodes a small graph of French cities connected by
// one‑way roads.  It uses RDFS/OWL rules to derive longer paths from shorter
// ones and answers the question: which cities can reach Nantes?
//
// This is intentionally not a full N3 reasoner – it is a concrete scenario that
// mirrors the structure of the original N3 example.
//
// Run:
//
//     go run french_cities.go
//
// The program has no third-party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "french_cities"

// ---------- types ----------

type Edge struct {
	From string
	To   string
}

type Dataset struct {
	Edges  []Edge
	Labels map[string]string // maps internal id to display name
}

type Checks struct {
	C1AngersDirect     bool
	C2LeMansChain      bool
	C3ChartresChain    bool
	C4ParisChain       bool
	C5NonReachableFuse bool
}

// ---------- data ----------

func dataset() Dataset {
	return Dataset{
		Edges: []Edge{
			{From: "paris", To: "orleans"},
			{From: "paris", To: "chartres"},
			{From: "paris", To: "amiens"},
			{From: "orleans", To: "blois"},
			{From: "orleans", To: "bourges"},
			{From: "blois", To: "tours"},
			{From: "chartres", To: "lemans"},
			{From: "lemans", To: "angers"},
			{From: "lemans", To: "tours"},
			{From: "angers", To: "nantes"},
		},
		Labels: map[string]string{
			"paris":    "Paris",
			"chartres": "Chartres",
			"lemans":   "Le Mans",
			"angers":   "Angers",
			"nantes":   "Nantes",
			"orleans":  "Orléans",
			"blois":    "Blois",
			"bourges":  "Bourges",
			"tours":    "Tours",
			"amiens":   "Amiens",
		},
	}
}

// ---------- inference ----------

func reachable(data Dataset) map[string]bool {
	cities := map[string]bool{}
	for _, e := range data.Edges {
		cities[e.From] = true
		cities[e.To] = true
	}

	paths := map[string]map[string]bool{}
	for city := range cities {
		paths[city] = map[string]bool{}
	}
	for _, e := range data.Edges {
		paths[e.From][e.To] = true
	}

	// transitive closure
	for mid := range cities {
		for from := range cities {
			if !paths[from][mid] {
				continue
			}
			for to := range cities {
				if paths[mid][to] {
					paths[from][to] = true
				}
			}
		}
	}

	result := map[string]bool{}
	for city := range cities {
		if paths[city]["nantes"] {
			result[city] = true
		}
	}
	return result
}

// ---------- checks ----------

func performChecks(data Dataset, canReach map[string]bool) Checks {
	c := Checks{}

	// C1: Angers -> Nantes direct
	c.C1AngersDirect = hasEdge(data, "angers", "nantes")

	// C2: Le Mans -> Angers -> Nantes
	c.C2LeMansChain = hasEdge(data, "lemans", "angers") && canReach["angers"]

	// C3: Chartres -> Le Mans -> Angers -> Nantes
	c.C3ChartresChain = hasEdge(data, "chartres", "lemans") && canReach["lemans"]

	// C4: Paris -> Chartres -> Le Mans -> Angers -> Nantes
	c.C4ParisChain = hasEdge(data, "paris", "chartres") && canReach["chartres"]

	// C5: non-reachable cities are correctly rejected
	nonReachable := []string{"orleans", "amiens", "bourges", "blois", "tours"}
	allRejected := true
	for _, city := range nonReachable {
		if canReach[city] {
			allRejected = false
			break
		}
	}
	c.C5NonReachableFuse = allRejected

	return c
}

func hasEdge(data Dataset, from, to string) bool {
	for _, e := range data.Edges {
		if e.From == from && e.To == to {
			return true
		}
	}
	return false
}

func allChecksPass(c Checks) bool {
	return c.C1AngersDirect && c.C2LeMansChain && c.C3ChartresChain &&
		c.C4ParisChain && c.C5NonReachableFuse
}

func checkCount(c Checks) int {
	count := 0
	if c.C1AngersDirect {
		count++
	}
	if c.C2LeMansChain {
		count++
	}
	if c.C3ChartresChain {
		count++
	}
	if c.C4ParisChain {
		count++
	}
	if c.C5NonReachableFuse {
		count++
	}
	return count
}

// ---------- rendering ----------

func renderArcOutput(data Dataset, canReach map[string]bool, checks Checks) {
	// reachable cities in a fixed order
	reachableOrder := []string{"paris", "chartres", "lemans", "angers"}
	var reachableNames []string
	for _, id := range reachableOrder {
		if canReach[id] {
			reachableNames = append(reachableNames, data.Labels[id])
		}
	}

	// === Answer ===
	fmt.Println("=== Answer ===")
	fmt.Printf("Four cities in this small network can reach Nantes: %s.\n",
		strings.Join(reachableNames, ", "))
	fmt.Println()

	// === Reason Why ===
	fmt.Println("=== Reason Why ===")
	fmt.Println("The original example says that every :oneway link is also a :path, and that :path is transitive. So once Angers can reach Nantes directly, longer routes can be built by chaining earlier links. Angers reaches Nantes directly. Le Mans reaches Nantes through Angers. Chartres reaches Nantes through Le Mans and Angers. Paris reaches Nantes through Chartres, Le Mans, and Angers.")
	fmt.Println()

	// === Check ===
	fmt.Println("=== Check ===")
	if checks.C1AngersDirect {
		fmt.Println("C1 OK - Angers has a direct one-way connection to Nantes.")
	} else {
		fmt.Println("C1 FAIL - Angers does not have a direct connection to Nantes.")
	}
	if checks.C2LeMansChain {
		fmt.Println("C2 OK - Le Mans reaches Nantes by chaining Le Mans → Angers → Nantes.")
	} else {
		fmt.Println("C2 FAIL - Le Mans chain check failed.")
	}
	if checks.C3ChartresChain {
		fmt.Println("C3 OK - Chartres reaches Nantes by chaining Chartres → Le Mans → Angers → Nantes.")
	} else {
		fmt.Println("C3 FAIL - Chartres chain check failed.")
	}
	if checks.C4ParisChain {
		fmt.Println("C4 OK - Paris reaches Nantes by chaining Paris → Chartres → Le Mans → Angers → Nantes.")
	} else {
		fmt.Println("C4 FAIL - Paris chain check failed.")
	}
	if checks.C5NonReachableFuse {
		fmt.Println("C5 OK - cities without a chain to Nantes are rejected by fail-loud fuse rules.")
	} else {
		fmt.Println("C5 FAIL - non-reachable city check failed.")
	}
	fmt.Println()

	// === Go audit details ===
	fmt.Println("=== Go audit details ===")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("total edges (one-way roads) : %d\n", len(data.Edges))
	fmt.Println("one-way connections:")
	for _, e := range data.Edges {
		fmt.Printf("  %s → %s\n", data.Labels[e.From], data.Labels[e.To])
	}
	fmt.Println()
	fmt.Println("reachable from Nantes:")
	for _, id := range reachableOrder {
		marker := "no"
		if canReach[id] {
			marker = "yes"
		}
		fmt.Printf("  %s : %s\n", data.Labels[id], marker)
	}
	nonReachableOrder := []string{"orleans", "amiens", "bourges", "blois", "tours"}
	for _, id := range nonReachableOrder {
		marker := "no"
		if canReach[id] {
			marker = "yes"
		}
		fmt.Printf("  %s : %s\n", data.Labels[id], marker)
	}
	fmt.Printf("checks passed : %d/5\n", checkCount(checks))
	fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(checks)))
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// ---------- main ----------

func main() {
	data := exampleinput.Load(eyelingoExampleName, dataset())
	canReach := reachable(data)

	// quick sanity
	expectedReachable := map[string]bool{
		"paris": true, "chartres": true, "lemans": true, "angers": true,
	}
	for city, expect := range expectedReachable {
		if canReach[city] != expect {
			fmt.Fprintf(os.Stderr, "internal error: %s reachable mismatch\n", city)
			os.Exit(1)
		}
	}

	checks := performChecks(data, canReach)
	renderArcOutput(data, canReach, checks)
}
