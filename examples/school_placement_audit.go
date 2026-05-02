// school_placement_audit.go
//
// A small algorithmic-governance example inspired by school-placement disputes
// where straight-line distance can hide real walking routes across rivers,
// highways, and other barriers.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"sort"
	"strings"
)

const eyelingoExampleName = "school_placement_audit"

type Dataset struct {
	CaseName  string     `json:"caseName"`
	Context   Context    `json:"context"`
	Policy    Policy     `json:"policy"`
	Students  []Student  `json:"students"`
	Schools   []School   `json:"schools"`
	Distances []Distance `json:"distances"`
}

type Context struct {
	City                 string `json:"city"`
	DecisionSupportLabel string `json:"decisionSupportLabel"`
	Issue                string `json:"issue"`
}

type Policy struct {
	MaxWalkingMeters        int `json:"maxWalkingMeters"`
	PreferencePenaltyMeters int `json:"preferencePenaltyMeters"`
}

type Student struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	HomeArea    string   `json:"homeArea"`
	Preferences []string `json:"preferences"`
}

type School struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Distance struct {
	Student        string `json:"student"`
	School         string `json:"school"`
	StraightMeters int    `json:"straightMeters"`
	WalkingMeters  int    `json:"walkingMeters"`
	Barrier        string `json:"barrier"`
}

type Assignment struct {
	Student  Student
	School   string
	Distance Distance
	Score    int
}

type Change struct {
	Student  Student
	Straight Assignment
	Audited  Assignment
	Flagged  bool
}

type Analysis struct {
	Straight            map[string]Assignment
	Audited             map[string]Assignment
	Changes             []Change
	LargestHiddenDetour Change
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printReport(ds, analysis)
}

func derive(ds Dataset) Analysis {
	straight := map[string]Assignment{}
	audited := map[string]Assignment{}
	for _, student := range ds.Students {
		straight[student.ID] = chooseSchool(ds, student, false)
		audited[student.ID] = chooseSchool(ds, student, true)
	}
	changes := []Change{}
	var largest Change
	largestSet := false
	for _, student := range ds.Students {
		flawed := straight[student.ID]
		better := audited[student.ID]
		flagged := flawed.School != better.School || flawed.Distance.WalkingMeters > ds.Policy.MaxWalkingMeters
		if flagged {
			change := Change{Student: student, Straight: flawed, Audited: better, Flagged: true}
			changes = append(changes, change)
			detour := hiddenDetour(flawed)
			if !largestSet || detour > hiddenDetour(largest.Straight) {
				largest = change
				largestSet = true
			}
		}
	}
	return Analysis{Straight: straight, Audited: audited, Changes: changes, LargestHiddenDetour: largest}
}

func chooseSchool(ds Dataset, student Student, audited bool) Assignment {
	best := Assignment{Score: math.MaxInt}
	for _, distance := range ds.Distances {
		if distance.Student != student.ID {
			continue
		}
		rank := preferenceRank(student, distance.School)
		score := distance.StraightMeters
		if audited {
			score = distance.WalkingMeters + rank*ds.Policy.PreferencePenaltyMeters
		}
		candidate := Assignment{Student: student, School: distance.School, Distance: distance, Score: score}
		if candidateLess(candidate, best, student) {
			best = candidate
		}
	}
	return best
}

func candidateLess(a, b Assignment, student Student) bool {
	if a.Score != b.Score {
		return a.Score < b.Score
	}
	ar := preferenceRank(student, a.School)
	br := preferenceRank(student, b.School)
	if ar != br {
		return ar < br
	}
	return a.School < b.School
}

func preferenceRank(student Student, school string) int {
	for index, preferred := range student.Preferences {
		if preferred == school {
			return index
		}
	}
	return 999
}

func hiddenDetour(a Assignment) int {
	return a.Distance.WalkingMeters - a.Distance.StraightMeters
}

func printReport(ds Dataset, analysis Analysis) {
	fmt.Println("# School Placement Route Audit")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("audit result : fail")
	fmt.Printf("children affected by straight-line rule : %s\n", names(analysis.Changes))
	fmt.Printf("largest hidden detour : %s, %d m\n", analysis.LargestHiddenDetour.Student.Name, hiddenDetour(analysis.LargestHiddenDetour.Straight))
	fmt.Printf("recommended assignments : %s\n", assignmentList(ds.Students, analysis.Audited))
	fmt.Println("explanation requested : yes")
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The support-tool rule chooses the school with the smallest straight-line distance, using preference rank only as a tie-breaker.")
	fmt.Printf("The independent audit recomputes each candidate with walking-route distance plus %d m per preference step.\n", ds.Policy.PreferencePenaltyMeters)
	fmt.Printf("Any provisional assignment that is not the audited best, or that requires more than %d m of walking, is flagged.\n", ds.Policy.MaxWalkingMeters)
	fmt.Println("Ada and Björn look close to Centrum on a map, but their walking routes cross barriers and exceed the walking limit; Davi is also better served by the first-preference Haga route.")
	fmt.Println("This illustrates why a decision-support label is not enough: route geometry, preferences, and audit records must be inspectable.")
	fmt.Println()
}

func names(changes []Change) string {
	out := []string{}
	for _, change := range changes {
		out = append(out, change.Student.Name)
	}
	sort.Strings(out)
	return strings.Join(out, ", ")
}

func assignmentList(students []Student, assignments map[string]Assignment) string {
	parts := []string{}
	for _, student := range students {
		parts = append(parts, fmt.Sprintf("%s -> %s", student.Name, assignments[student.ID].School))
	}
	return strings.Join(parts, "; ")
}
