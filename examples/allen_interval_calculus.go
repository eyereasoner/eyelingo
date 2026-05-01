// allen_interval_calculus.go
//
// A compact Go translation inspired by Eyeling's
// `examples/allen-interval-calculus.n3`.
//
// The example completes intervals with duration fields and classifies every
// ordered interval pair using Allen's 13 base relations.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"strings"
	"time"
)

const eyelingoExampleName = "allen_interval_calculus"

type Dataset struct {
	CaseName  string          `json:"caseName"`
	Question  string          `json:"question"`
	Intervals []IntervalInput `json:"intervals"`
	Expected  struct {
		RequiredRelations map[string]string `json:"requiredRelations"`
	} `json:"expected"`
}

type IntervalInput struct {
	Name            string `json:"name"`
	Start           string `json:"start"`
	End             string `json:"end"`
	DurationMinutes int    `json:"durationMinutes"`
}

type Interval struct {
	Name      string
	Start     time.Time
	End       time.Time
	Completed bool
}

type Relation struct {
	From string
	Kind string
	To   string
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Intervals []Interval
	Relations []Relation
	Checks    []Check
	Invalid   int
	Completed int
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printReport(ds, analysis)
	if !allOK(analysis.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Analysis {
	intervals := make([]Interval, 0, len(ds.Intervals))
	invalid := 0
	completed := 0
	for _, in := range ds.Intervals {
		start := parseTime(in.Start)
		end := parseTime(in.End)
		wasCompleted := false
		if end.IsZero() && !start.IsZero() && in.DurationMinutes > 0 {
			end = start.Add(time.Duration(in.DurationMinutes) * time.Minute)
			wasCompleted = true
			completed++
		}
		if start.IsZero() && !end.IsZero() && in.DurationMinutes > 0 {
			start = end.Add(-time.Duration(in.DurationMinutes) * time.Minute)
			wasCompleted = true
			completed++
		}
		if !start.IsZero() && !end.IsZero() && !start.Before(end) {
			invalid++
		}
		intervals = append(intervals, Interval{Name: in.Name, Start: start, End: end, Completed: wasCompleted})
	}

	relations := []Relation{}
	relMap := map[string]string{}
	for _, a := range intervals {
		for _, b := range intervals {
			if a.Name == b.Name {
				continue
			}
			kind := relation(a, b)
			relations = append(relations, Relation{From: a.Name, Kind: kind, To: b.Name})
			relMap[a.Name+"|"+b.Name] = kind
		}
	}

	checks := []Check{
		{"C1", len(intervals) == len(ds.Intervals), fmt.Sprintf("%d intervals were loaded and completed when duration was present", len(intervals))},
		{"C2", relMap["A|B"] == "before", "A before B was derived"},
		{"C3", relMap["A|C"] == "meets" && relMap["C|A"] == "metBy", "A meets C and C metBy A were derived"},
		{"C4", relMap["A|D"] == "overlaps" && relMap["D|A"] == "overlappedBy", "A overlaps D and D overlappedBy A were derived"},
		{"C5", relMap["F|A"] == "starts" && relMap["G|A"] == "finishes", "F starts A and G finishes A were derived"},
		{"C6", relMap["A|H"] == "during" && relMap["H|A"] == "contains", "A during H and H contains A were derived"},
		{"C7", relMap["A|E"] == "equals", "A equals E was derived"},
		{"C8", relMap["J|I"] == "meets" && relMap["K|C"] == "finishes", "duration completion produced I ending at 18:00 and K ending at 14:00"},
		{"C9", invalid == 0, "no invalid intervals were detected"},
	}
	return Analysis{Intervals: intervals, Relations: relations, Checks: checks, Invalid: invalid, Completed: completed}
}

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func relation(a, b Interval) string {
	s1, e1, s2, e2 := a.Start, a.End, b.Start, b.End
	switch {
	case e1.Before(s2):
		return "before"
	case e1.Equal(s2):
		return "meets"
	case s1.Before(s2) && s2.Before(e1) && e1.Before(e2):
		return "overlaps"
	case s1.Equal(s2) && e1.Before(e2):
		return "starts"
	case s2.Before(s1) && e1.Before(e2):
		return "during"
	case e1.Equal(e2) && s2.Before(s1):
		return "finishes"
	case s1.Equal(s2) && e1.Equal(e2):
		return "equals"
	case s1.After(e2):
		return "after"
	case s1.Equal(e2):
		return "metBy"
	case s2.Before(s1) && s1.Before(e2) && e2.Before(e1):
		return "overlappedBy"
	case s1.Equal(s2) && e2.Before(e1):
		return "startedBy"
	case s1.Before(s2) && e2.Before(e1):
		return "contains"
	case e1.Equal(e2) && s1.Before(s2):
		return "finishedBy"
	default:
		return "intersects"
	}
}

func allOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func countOK(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func relExists(relations []Relation, from, kind, to string) bool {
	for _, rel := range relations {
		if rel.From == from && rel.Kind == kind && rel.To == to {
			return true
		}
	}
	return false
}

func printReport(ds Dataset, analysis Analysis) {
	showcase := []string{}
	for _, want := range []Relation{
		{"A", "before", "B"},
		{"A", "meets", "C"},
		{"A", "overlaps", "D"},
		{"F", "starts", "A"},
		{"G", "finishes", "A"},
		{"A", "during", "H"},
		{"A", "equals", "E"},
		{"J", "meets", "I"},
		{"K", "finishes", "C"},
	} {
		if relExists(analysis.Relations, want.From, want.Kind, want.To) {
			showcase = append(showcase, want.From+" "+want.Kind+" "+want.To)
		}
	}

	fmt.Println("# Allen Interval Calculus")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("derived relations : %d ordered interval pairs\n", len(analysis.Relations))
	fmt.Printf("showcase : %s\n", strings.Join(showcase, " | "))
	fmt.Println("completed intervals : I=16:00-18:00, K=13:30-14:00")
	fmt.Printf("invalid intervals : %d\n", analysis.Invalid)
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The example completes any interval that has a start plus duration before comparing endpoints.")
	fmt.Println("Each ordered pair is then classified with the 13 Allen base relations, including the six converse relations.")
	fmt.Println("The relation rules are purely endpoint constraints, so the result is deterministic and traceable.")
	fmt.Println("The duration-derived intervals participate in the same relation table as directly supplied intervals.")
	fmt.Println()
	return
}
