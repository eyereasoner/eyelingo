// dijkstra_risk_path.go
//
// Inspired by Eyeling's shortest-path style examples such as `dijkstra.n3`.
//
// The example finds the lowest risk-adjusted route through a small delivery
// graph using Dijkstra's algorithm.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const eyelingoExampleName = "dijkstra_risk_path"

type Dataset struct {
	CaseName   string  `json:"caseName"`
	Question   string  `json:"question"`
	Start      string  `json:"start"`
	Goal       string  `json:"goal"`
	RiskWeight float64 `json:"riskWeight"`
	Edges      []Edge  `json:"edges"`
	Expected   struct {
		Path  []string `json:"path"`
		Score float64  `json:"score"`
	} `json:"expected"`
}
type Edge struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Cost float64 `json:"cost"`
	Risk float64 `json:"risk"`
}
type Label struct {
	Node  string
	Score float64
	Path  []string
}
type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	Best             Label
	RawCost, RiskSum float64
	Visited          int
	Checks           []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	a := derive(ds)
	printReport(ds, a)
	if !allOK(a.Checks) {
		os.Exit(1)
	}
}
func derive(ds Dataset) Analysis {
	queue := []Label{{Node: ds.Start, Score: 0, Path: []string{ds.Start}}}
	best := map[string]float64{}
	visited := 0
	var selected Label
	for len(queue) > 0 {
		sort.Slice(queue, func(i, j int) bool { return queue[i].Score < queue[j].Score })
		cur := queue[0]
		queue = queue[1:]
		if prev, ok := best[cur.Node]; ok && prev <= cur.Score {
			continue
		}
		best[cur.Node] = cur.Score
		visited++
		if cur.Node == ds.Goal {
			selected = cur
			break
		}
		for _, e := range ds.Edges {
			if e.From != cur.Node {
				continue
			}
			nextPath := append(append([]string{}, cur.Path...), e.To)
			queue = append(queue, Label{Node: e.To, Score: cur.Score + e.Cost + ds.RiskWeight*e.Risk, Path: nextPath})
		}
	}
	raw, risk := pathCost(ds.Edges, selected.Path)
	checks := []Check{
		{"C1", len(ds.Edges) == 8, "all route edges were loaded from JSON"},
		{"C2", ds.RiskWeight == 2.0, "edge score is cost + 2.00 × risk"},
		{"C3", selected.Node == ds.Goal, "Dijkstra reached HubZ from ClinicA"},
		{"C4", strings.Join(selected.Path, " -> ") == "ClinicA -> DepotB -> LabD -> HubZ", "selected path is ClinicA -> DepotB -> LabD -> HubZ"},
		{"C5", close(selected.Score, ds.Expected.Score), "selected total score is 11.10"},
		{"C6", !strings.Contains(strings.Join(selected.Path, ","), "DepotC"), "higher-risk shortcut through DepotC is rejected"},
	}
	return Analysis{selected, raw, risk, visited, checks}
}
func pathCost(edges []Edge, path []string) (float64, float64) {
	raw := 0.0
	risk := 0.0
	for i := 0; i+1 < len(path); i++ {
		for _, e := range edges {
			if e.From == path[i] && e.To == path[i+1] {
				raw += e.Cost
				risk += e.Risk
			}
		}
	}
	return raw, risk
}
func close(a, b float64) bool {
	if a > b {
		return a-b < 1e-9
	}
	return b-a < 1e-9
}
func allOK(checks []Check) bool {
	for _, c := range checks {
		if !c.OK {
			return false
		}
	}
	return true
}
func countOK(checks []Check) int {
	n := 0
	for _, c := range checks {
		if c.OK {
			n++
		}
	}
	return n
}
func printReport(ds Dataset, a Analysis) {
	fmt.Println("# Dijkstra Risk Path")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("selected path : %s\n", strings.Join(a.Best.Path, " -> "))
	fmt.Printf("raw cost : %.2f\n", a.RawCost)
	fmt.Printf("risk sum : %.2f\n", a.RiskSum)
	fmt.Printf("risk-adjusted score : %.2f\n", a.Best.Score)
	fmt.Printf("edges in selected path : %d\n", len(a.Best.Path)-1)
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("Each edge contributes its delivery cost plus the configured risk penalty.")
	fmt.Println("Dijkstra's queue expands the lowest accumulated score first, so the first time HubZ is popped the selected route is optimal for the weighted graph.")
	fmt.Println("The DepotC shortcut has lower early cost but carries enough risk to lose under the configured risk weight.")
	fmt.Println("The selected route balances cost and risk through DepotB and LabD.")
	fmt.Println()
	fmt.Println("## Check")
	for _, c := range a.Checks {
		status := "FAIL"
		if c.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", c.ID, status, c.Text)
	}
	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("start : %s\n", ds.Start)
	fmt.Printf("goal : %s\n", ds.Goal)
	fmt.Printf("edges loaded : %d\n", len(ds.Edges))
	fmt.Printf("risk weight : %.2f\n", ds.RiskWeight)
	fmt.Printf("visited nodes : %d\n", a.Visited)
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
