// deep_taxonomy_100000.go
//
// A self-contained Go translation of deep-taxonomy-100000.n3 from the Eyeling
// examples.
//
// The original N3 file is an ARC-style stress test with 100,000 nearly
// identical taxonomy-step rules:
//
//	{?X a :N42} => {?X a :N43, :I43, :J43}.
//
// Starting from :ind a :N0, the chain must reach :N100000. The terminal rule
// then derives :A2, and the success rule derives :test :is true.
//
// This Go version translates the concrete rule family into a compact bit-set
// forward chainer. It is intentionally not a generic RDF/N3 reasoner. The goal
// is to make the deep classification derivation visible, deterministic, and
// fast without embedding five megabytes of repetitive Go literals.
//
// Run:
//
//	go run deep_taxonomy_100000.go
//
// The program has no third-party dependencies.
package main

import (
	"fmt"
	"math/bits"
	"os"
	"runtime"
	"strings"
)

const (
	taxonomyDepth = 100000
	midpointDepth = taxonomyDepth / 2

	sourceStepRules      = taxonomyDepth
	sourceTerminalRules  = 1
	sourceSuccessRules   = 1
	sourceARCCheckRules  = 6
	sourceARCReportRules = 1
	sideLabelsPerStep    = 2
	expectedNClasses     = taxonomyDepth + 1
	expectedSideLabels   = taxonomyDepth * sideLabelsPerStep
	expectedTypeFacts    = expectedNClasses + expectedSideLabels + 1 // + A2
	sourceFactAssertions = 1
	sourceTotalRules     = sourceStepRules + sourceTerminalRules + sourceSuccessRules + sourceARCCheckRules + sourceARCReportRules
)

type bitset []uint64

func newBitset(size int) bitset {
	return make(bitset, (size+63)/64)
}

func (b bitset) set(i int) bool {
	word := i / 64
	mask := uint64(1) << uint(i%64)
	if b[word]&mask != 0 {
		return false
	}
	b[word] |= mask
	return true
}

func (b bitset) has(i int) bool {
	word := i / 64
	mask := uint64(1) << uint(i%64)
	return b[word]&mask != 0
}

func (b bitset) count() int {
	n := 0
	for _, word := range b {
		n += bits.OnesCount64(word)
	}
	return n
}

type TaxonomyFacts struct {
	N    bitset
	I    bitset
	J    bitset
	A2   bool
	Test bool
}

type RunStats struct {
	AgendaPops               int
	StepRuleTests            int
	StepRuleApplications     int
	TerminalRuleTests        int
	TerminalRuleApplications int
	SuccessRuleTests         int
	SuccessRuleApplications  int
	NewNClasses              int
	NewSideI                 int
	NewSideJ                 int
	DuplicateAssertions      int
	MaxNReached              int
	LastAgenda               int
}

type Checks struct {
	StartPresent           bool
	FirstExpansionComplete bool
	MidpointComplete       bool
	FinalStepComplete      bool
	TerminalClassDerived   bool
	SuccessFlagRaised      bool
	ClassCountsCorrect     bool
	SideLabelCountsCorrect bool
	NoSkippedLevel         bool
	RuleApplicationsExact  bool
}

type Result struct {
	Facts  TaxonomyFacts
	Stats  RunStats
	Checks Checks
}

func runForwardChain(depth int) Result {
	facts := TaxonomyFacts{
		N: newBitset(depth + 1),
		I: newBitset(depth + 1),
		J: newBitset(depth + 1),
	}
	stats := RunStats{MaxNReached: 0, LastAgenda: -1}

	agenda := make([]int, 0, depth+1)
	if facts.N.set(0) {
		stats.NewNClasses++
		agenda = append(agenda, 0)
	}

	for len(agenda) > 0 {
		current := agenda[0]
		copy(agenda, agenda[1:])
		agenda = agenda[:len(agenda)-1]

		stats.AgendaPops++
		stats.LastAgenda = current
		if current > stats.MaxNReached {
			stats.MaxNReached = current
		}

		if current < depth {
			stats.StepRuleTests++
			next := current + 1
			applied := false

			if facts.N.set(next) {
				stats.NewNClasses++
				agenda = append(agenda, next)
				applied = true
			} else {
				stats.DuplicateAssertions++
			}
			if facts.I.set(next) {
				stats.NewSideI++
				applied = true
			} else {
				stats.DuplicateAssertions++
			}
			if facts.J.set(next) {
				stats.NewSideJ++
				applied = true
			} else {
				stats.DuplicateAssertions++
			}
			if applied {
				stats.StepRuleApplications++
			}
			continue
		}

		stats.TerminalRuleTests++
		if facts.N.has(depth) && !facts.A2 {
			facts.A2 = true
			stats.TerminalRuleApplications++
		}
	}

	stats.SuccessRuleTests++
	if facts.A2 && !facts.Test {
		facts.Test = true
		stats.SuccessRuleApplications++
	}

	checks := buildChecks(facts, stats, depth)
	return Result{Facts: facts, Stats: stats, Checks: checks}
}

func buildChecks(facts TaxonomyFacts, stats RunStats, depth int) Checks {
	return Checks{
		StartPresent:           facts.N.has(0),
		FirstExpansionComplete: facts.N.has(1) && facts.I.has(1) && facts.J.has(1),
		MidpointComplete:       facts.N.has(depth/2) && facts.I.has(depth/2) && facts.J.has(depth/2),
		FinalStepComplete:      facts.N.has(depth-1) && facts.N.has(depth),
		TerminalClassDerived:   facts.N.has(depth) && facts.A2,
		SuccessFlagRaised:      facts.A2 && facts.Test,
		ClassCountsCorrect:     facts.N.count() == depth+1,
		SideLabelCountsCorrect: facts.I.count() == depth && facts.J.count() == depth,
		NoSkippedLevel:         firstMissingN(facts.N, depth) == -1,
		RuleApplicationsExact:  stats.StepRuleApplications == depth && stats.TerminalRuleApplications == 1 && stats.SuccessRuleApplications == 1,
	}
}

func firstMissingN(n bitset, depth int) int {
	for i := 0; i <= depth; i++ {
		if !n.has(i) {
			return i
		}
	}
	return -1
}

func renderAnswer(result Result) {
	fmt.Println("=== Answer ===")
	if result.Facts.Test {
		fmt.Println("The deep taxonomy test succeeds.")
		fmt.Printf("Starting fact : :ind a :N0\n")
		fmt.Printf("Reached class : :ind a :N%d\n", taxonomyDepth)
		fmt.Println("Terminal class : :ind a :A2")
		fmt.Println("Success flag : :test :is true")
	} else {
		fmt.Println("The deep taxonomy test does not succeed.")
	}
	fmt.Println()
	fmt.Println("Proof checkpoints:")
	fmt.Printf(" - :N0 present : %s\n", yesNo(result.Facts.N.has(0)))
	fmt.Printf(" - :N1 plus :I1/:J1 present : %s\n", yesNo(result.Facts.N.has(1) && result.Facts.I.has(1) && result.Facts.J.has(1)))
	fmt.Printf(" - :N%d plus :I%d/:J%d present : %s\n", midpointDepth, midpointDepth, midpointDepth, yesNo(result.Facts.N.has(midpointDepth) && result.Facts.I.has(midpointDepth) && result.Facts.J.has(midpointDepth)))
	fmt.Printf(" - :N99999 and :N100000 present : %s\n", yesNo(result.Facts.N.has(99999) && result.Facts.N.has(100000)))
	fmt.Printf(" - :A2 and success flag present : %s\n", yesNo(result.Facts.A2 && result.Facts.Test))
	fmt.Println()
}

func renderReason(result Result) {
	fmt.Println("=== Reason Why ===")
	fmt.Println("The N3 source is a very deep rule chain. Each taxonomy-step rule consumes the same individual in class Ni and derives the next class N(i+1), plus two side labels I(i+1) and J(i+1). Once N100000 is present, the terminal rule derives A2; once A2 is present, the success rule derives :test :is true.")
	fmt.Printf("source N3 starting fact assertions : %d\n", sourceFactAssertions)
	fmt.Printf("source N3 taxonomy-step rules : %d\n", sourceStepRules)
	fmt.Printf("source N3 terminal/success rules : %d\n", sourceTerminalRules+sourceSuccessRules)
	fmt.Printf("source N3 ARC check/report rules : %d\n", sourceARCCheckRules+sourceARCReportRules)
	fmt.Printf("source N3 total rules counted : %d\n", sourceTotalRules)
	fmt.Printf("translated representation : compressed rule schema + %d-word bit sets\n", len(result.Facts.N))
	fmt.Printf("agenda pops : %d\n", result.Stats.AgendaPops)
	fmt.Printf("taxonomy-step rule applications : %d\n", result.Stats.StepRuleApplications)
	fmt.Printf("terminal rule applications : %d\n", result.Stats.TerminalRuleApplications)
	fmt.Printf("success rule applications : %d\n", result.Stats.SuccessRuleApplications)
	fmt.Printf("classification facts derived : %d N classes + %d side labels + A2 = %d type facts\n", result.Facts.N.count(), result.Facts.I.count()+result.Facts.J.count(), expectedTypeFacts)
	fmt.Println("The side labels are not needed for the final A2 proof, but carrying both I and J at every level checks that the whole wide/deep expansion was performed, not just the main N-chain.")
	fmt.Println()
}

func renderChecks(checks Checks) {
	fmt.Println("=== Check ===")
	lines := []struct {
		ok   bool
		text string
	}{
		{checks.StartPresent, "C1 OK - the starting classification N0 is present."},
		{checks.FirstExpansionComplete, "C2 OK - the first expansion produced N1 together with side labels I1 and J1."},
		{checks.MidpointComplete, "C3 OK - the chain reaches the midpoint N50000 and still carries both side-label branches."},
		{checks.FinalStepComplete, "C4 OK - the final taxonomy step from N99999 to N100000 was completed."},
		{checks.TerminalClassDerived, "C5 OK - once N100000 is reached, the terminal class A2 is derived."},
		{checks.SuccessFlagRaised, "C6 OK - the success flag is raised only after the terminal class A2 is present."},
		{checks.ClassCountsCorrect, "C7 OK - all 100001 N-level classifications are present with no missing level."},
		{checks.SideLabelCountsCorrect, "C8 OK - all 200000 side-label classifications I1..I100000 and J1..J100000 are present."},
		{checks.RuleApplicationsExact, "C9 OK - exactly 100000 taxonomy-step rules, one terminal rule, and one success rule fired."},
		{checks.NoSkippedLevel, "C10 OK - a linear scan finds no skipped taxonomy level."},
	}
	for _, line := range lines {
		if line.ok {
			fmt.Println(line.text)
		} else {
			fmt.Println(strings.Replace(line.text, "OK", "FAIL", 1))
		}
	}
	fmt.Println()
}

func renderAudit(result Result) {
	fmt.Println("=== Go audit details ===")
	fmt.Printf("go runtime : %s\n", runtime.Version())
	fmt.Printf("go os/arch : %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("source file : deep-taxonomy-100000.n3")
	fmt.Println("question : Can one starting individual be classified through 100000 taxonomy levels and trigger A2?")
	fmt.Printf("taxonomy depth : %d\n", taxonomyDepth)
	fmt.Printf("midpoint depth : %d\n", midpointDepth)
	fmt.Printf("source fact assertions : %d\n", sourceFactAssertions)
	fmt.Printf("source taxonomy-step rules : %d\n", sourceStepRules)
	fmt.Printf("source terminal rules : %d\n", sourceTerminalRules)
	fmt.Printf("source success rules : %d\n", sourceSuccessRules)
	fmt.Printf("source ARC check rules : %d\n", sourceARCCheckRules)
	fmt.Printf("source ARC report rules : %d\n", sourceARCReportRules)
	fmt.Printf("source total rules counted : %d\n", sourceTotalRules)
	fmt.Printf("bitset words per N/I/J column : %d\n", len(result.Facts.N))
	fmt.Printf("N classifications present : %d\n", result.Facts.N.count())
	fmt.Printf("I side labels present : %d\n", result.Facts.I.count())
	fmt.Printf("J side labels present : %d\n", result.Facts.J.count())
	fmt.Printf("A2 present : %s\n", yesNo(result.Facts.A2))
	fmt.Printf("test true : %s\n", yesNo(result.Facts.Test))
	fmt.Printf("expected type facts : %d\n", expectedTypeFacts)
	fmt.Printf("actual type facts : %d\n", result.Facts.N.count()+result.Facts.I.count()+result.Facts.J.count()+boolInt(result.Facts.A2))
	fmt.Printf("agenda pops : %d\n", result.Stats.AgendaPops)
	fmt.Printf("step rule tests : %d\n", result.Stats.StepRuleTests)
	fmt.Printf("step rule applications : %d\n", result.Stats.StepRuleApplications)
	fmt.Printf("terminal rule tests : %d\n", result.Stats.TerminalRuleTests)
	fmt.Printf("terminal rule applications : %d\n", result.Stats.TerminalRuleApplications)
	fmt.Printf("success rule tests : %d\n", result.Stats.SuccessRuleTests)
	fmt.Printf("success rule applications : %d\n", result.Stats.SuccessRuleApplications)
	fmt.Printf("new N classes : %d\n", result.Stats.NewNClasses)
	fmt.Printf("new I side labels : %d\n", result.Stats.NewSideI)
	fmt.Printf("new J side labels : %d\n", result.Stats.NewSideJ)
	fmt.Printf("duplicate assertions ignored : %d\n", result.Stats.DuplicateAssertions)
	fmt.Printf("max N reached : N%d\n", result.Stats.MaxNReached)
	fmt.Printf("last agenda item : N%d\n", result.Stats.LastAgenda)
	fmt.Printf("checks passed : %d/10\n", countChecks(result.Checks))
	fmt.Printf("all checks pass : %s\n", yesNo(countChecks(result.Checks) == 10))
}

func countChecks(checks Checks) int {
	values := []bool{
		checks.StartPresent,
		checks.FirstExpansionComplete,
		checks.MidpointComplete,
		checks.FinalStepComplete,
		checks.TerminalClassDerived,
		checks.SuccessFlagRaised,
		checks.ClassCountsCorrect,
		checks.SideLabelCountsCorrect,
		checks.RuleApplicationsExact,
		checks.NoSkippedLevel,
	}
	n := 0
	for _, value := range values {
		if value {
			n++
		}
	}
	return n
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func main() {
	result := runForwardChain(taxonomyDepth)
	if !result.Facts.Test {
		fmt.Fprintln(os.Stderr, "deep taxonomy proof did not reach :test :is true")
		os.Exit(1)
	}

	renderAnswer(result)
	renderReason(result)
	renderChecks(result.Checks)
	renderAudit(result)
}
