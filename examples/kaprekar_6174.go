// kaprekar_6174.go
//
// A self-contained Go translation of kaprekar-6174.n3 from the Eyeling
// examples.
//
// Kaprekar's routine starts with any four digits, including leading zeroes.
// At each step it sorts those digits from high to low to make one number,
// sorts them from low to high to make another number, and subtracts the smaller
// from the larger. Most starts eventually reach 6174, Kaprekar's constant.
// Starts with four equal digits, such as 0000 or 2222, fall into 0000 instead
// and are intentionally omitted from the final :kaprekar facts.
//
// The N3 example unrolls the proof to a fixed seven-step bound, because no
// four-digit start that reaches 6174 needs more than seven Kaprekar steps. This
// Go version keeps the same bounded, explicit approach while using ordinary Go
// data structures instead of a general RDF/N3 reasoner.
//
// Run:
//
//	go run kaprekar_6174.go
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

const eyelingoExampleName = "kaprekar_6174"

var (
	startCount       = 10000
	targetConstant   = 6174
	zeroBasin        = 0
	maxKaprekarSteps = 7
)

type Chain struct {
	Start int
	Path  []int
}

type Stats struct {
	StartsEnumerated     int
	StepFactsComputed    int
	ChainsEmitted        int
	ZeroBasinStarts      int
	MaxStepsToTarget     int
	IdentityChecks       int
	DirectStepMismatches int
	HistogramByStepCount map[int]int
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Question string
	Chains   []Chain
	ByStart  map[int]Chain
	Omitted  map[int][]int
	Stats    Stats
	Checks   []Check
}

func main() {
	analysis := derive()
	printAnswer(analysis)
	printReason(analysis)
	printChecks(analysis.Checks)
	printAudit(analysis)
	if !allChecksOK(analysis.Checks) {
		os.Exit(1)
	}
}

func derive() Analysis {
	cfg := exampleinput.Load(eyelingoExampleName, struct {
		StartCount       int
		TargetConstant   int
		ZeroBasin        int
		MaxKaprekarSteps int
	}{StartCount: startCount, TargetConstant: targetConstant, ZeroBasin: zeroBasin, MaxKaprekarSteps: maxKaprekarSteps})
	startCount = cfg.StartCount
	targetConstant = cfg.TargetConstant
	zeroBasin = cfg.ZeroBasin
	maxKaprekarSteps = cfg.MaxKaprekarSteps
	chains := make([]Chain, 0, startCount)
	byStart := make(map[int]Chain)
	omitted := make(map[int][]int)
	stats := Stats{HistogramByStepCount: make(map[int]int)}

	for start := 0; start < startCount; start++ {
		stats.StartsEnumerated++
		if kaprekarStep(start) != directKaprekarStep(start) {
			stats.DirectStepMismatches++
		}
		stats.IdentityChecks++

		if start == targetConstant {
			chain := Chain{Start: start, Path: []int{targetConstant}}
			chains = append(chains, chain)
			byStart[start] = chain
			stats.ChainsEmitted++
			stats.HistogramByStepCount[0]++
			continue
		}

		current := start
		path := make([]int, 0, maxKaprekarSteps)
		reachedTarget := false
		reachedZero := false

		for stepNumber := 1; stepNumber <= maxKaprekarSteps; stepNumber++ {
			next := kaprekarStep(current)
			stats.StepFactsComputed++
			path = append(path, next)
			current = next

			if next == targetConstant {
				reachedTarget = true
				if stepNumber > stats.MaxStepsToTarget {
					stats.MaxStepsToTarget = stepNumber
				}
				stats.HistogramByStepCount[stepNumber]++
				break
			}
			if next == zeroBasin {
				reachedZero = true
				break
			}
		}

		if reachedTarget {
			chain := Chain{Start: start, Path: path}
			chains = append(chains, chain)
			byStart[start] = chain
		} else if reachedZero {
			stats.ZeroBasinStarts++
			omitted[start] = path
		} else {
			omitted[start] = path
		}
	}

	stats.ChainsEmitted = len(chains)
	analysis := Analysis{
		Question: "Which four-digit starts have a Kaprekar chain ending at 6174?",
		Chains:   chains,
		ByStart:  byStart,
		Omitted:  omitted,
		Stats:    stats,
	}
	analysis.Checks = buildChecks(analysis)
	return analysis
}

// kaprekarStep implements the optimized identity used in the N3 comments.
// If the sorted digits are a0 <= a1 <= a2 <= a3, then descending minus
// ascending is 999*(a3-a0) + 90*(a2-a1).
func kaprekarStep(n int) int {
	digits := digitsAscending(n)
	return 999*(digits[3]-digits[0]) + 90*(digits[2]-digits[1])
}

// directKaprekarStep is a slower, more literal version of the rule. It is used
// only as a check that the optimized identity computes the same value.
func directKaprekarStep(n int) int {
	digits := digitsAscending(n)
	ascending := 1000*digits[0] + 100*digits[1] + 10*digits[2] + digits[3]
	descending := 1000*digits[3] + 100*digits[2] + 10*digits[1] + digits[0]
	return descending - ascending
}

func digitsAscending(n int) [4]int {
	digits := [4]int{
		(n / 1000) % 10,
		(n / 100) % 10,
		(n / 10) % 10,
		n % 10,
	}
	sort.Ints(digits[:])
	return digits
}

func buildChecks(a Analysis) []Check {
	chain3524 := a.ByStart[3524]
	chain1 := a.ByStart[1]
	_, has1111 := a.ByStart[1111]
	_, has0000 := a.ByStart[0]

	checks := []Check{
		{
			Label: "complete enumeration",
			OK:    a.Stats.StartsEnumerated == startCount,
			Text:  "all digit patterns from 0000 through 9999 were considered",
		},
		{
			Label: "optimized step formula",
			OK:    a.Stats.DirectStepMismatches == 0 && a.Stats.IdentityChecks == startCount,
			Text:  "the identity-based step matches direct digit sorting for every start",
		},
		{
			Label: "known sample 3524",
			OK:    equalInts(chain3524.Path, []int{3087, 8352, 6174}),
			Text:  "3524 follows the classic 3087 -> 8352 -> 6174 chain",
		},
		{
			Label: "leading zero sample",
			OK:    equalInts(chain1.Path, []int{999, 8991, 8082, 8532, 6174}),
			Text:  "0001 is accepted as a four-digit start by treating it as 0,0,0,1",
		},
		{
			Label: "zero basin omitted",
			OK:    !has0000 && !has1111 && a.Stats.ZeroBasinStarts == 10,
			Text:  "0000 and the nine non-zero repdigits fall to 0000 and are not emitted",
		},
		{
			Label: "all emitted chains end at 6174",
			OK:    allChainsEndAt(a.Chains, targetConstant),
			Text:  "every :kaprekar fact kept by the translation reaches Kaprekar's constant",
		},
		{
			Label: "seven-step bound",
			OK:    a.Stats.MaxStepsToTarget == maxKaprekarSteps,
			Text:  "no emitted chain needs more than the seven steps unrolled in the N3 source",
		},
		{
			Label: "non-repdigit coverage",
			OK:    a.Stats.ChainsEmitted == startCount-10,
			Text:  "all 9990 non-repdigit starts are emitted, including 6174 itself",
		},
	}
	return checks
}

func printAnswer(a Analysis) {
	fmt.Println("=== Answer ===")
	fmt.Println("Kaprekar chains that end at 6174 are emitted as :kaprekar facts.")
	fmt.Printf("total emitted : %d\n", a.Stats.ChainsEmitted)
	fmt.Printf("omitted 0000 basin : %d\n", a.Stats.ZeroBasinStarts)
	fmt.Printf("maximum steps to 6174 : %d\n", a.Stats.MaxStepsToTarget)
	fmt.Println()
	fmt.Println("Selected facts, shown with four-digit padding for readability:")
	for _, start := range []int{1, 3524, 6174, 9831, 9998} {
		if chain, ok := a.ByStart[start]; ok {
			fmt.Printf("  %s :kaprekar (%s)\n", fourDigits(chain.Start), formatPath(chain.Path))
		}
	}
	fmt.Println()
}

func printReason(a Analysis) {
	fmt.Println("=== Reason Why ===")
	fmt.Println("Each start is read as four digits, so 1 is treated as 0001.")
	fmt.Println("The digits are sorted once, then the optimized identity computes the")
	fmt.Println("same result as descending-number minus ascending-number.")
	fmt.Println("The search is bounded to seven steps, matching the N3 source: any")
	fmt.Println("four-digit start that reaches 6174 does so within that bound.")
	fmt.Println()
	fmt.Println("Step-count distribution for emitted starts:")
	for steps := 0; steps <= maxKaprekarSteps; steps++ {
		fmt.Printf("  %d step(s) : %d start(s)\n", steps, a.Stats.HistogramByStepCount[steps])
	}
	fmt.Println()
	fmt.Println("Examples omitted because they fall to 0000:")
	for _, start := range []int{0, 1111, 2222, 9999} {
		fmt.Printf("  %s -> (%s)\n", fourDigits(start), formatPath(a.Omitted[start]))
	}
	fmt.Println()
}

func printChecks(checks []Check) {
	fmt.Println("=== Check ===")
	for index, check := range checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("C%d %s - %s\n", index+1, status, check.Text)
	}
	fmt.Println()
}

func printAudit(a Analysis) {
	fmt.Println("=== Go audit details ===")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Println("source file : kaprekar-6174.n3")
	fmt.Printf("question : %s\n", a.Question)
	fmt.Printf("starts enumerated : %d\n", a.Stats.StartsEnumerated)
	fmt.Printf("optimized step checks : %d\n", a.Stats.IdentityChecks)
	fmt.Printf("direct step mismatches : %d\n", a.Stats.DirectStepMismatches)
	fmt.Printf("bounded step applications : %d\n", a.Stats.StepFactsComputed)
	fmt.Printf("emitted kaprekar facts : %d\n", a.Stats.ChainsEmitted)
	fmt.Printf("omitted zero-basin starts : %d\n", a.Stats.ZeroBasinStarts)
	fmt.Printf("max steps to target : %d\n", a.Stats.MaxStepsToTarget)
	fmt.Printf("histogram : %s\n", formatHistogram(a.Stats.HistogramByStepCount))
	fmt.Printf("checks passed : %d/%d\n", countChecks(a.Checks), len(a.Checks))
	fmt.Printf("all checks pass : %s\n", yesNo(allChecksOK(a.Checks)))
}

func formatPath(path []int) string {
	values := make([]string, 0, len(path))
	for _, n := range path {
		values = append(values, fourDigits(n))
	}
	return strings.Join(values, " ")
}

func formatHistogram(histogram map[int]int) string {
	parts := make([]string, 0, maxKaprekarSteps+1)
	for steps := 0; steps <= maxKaprekarSteps; steps++ {
		parts = append(parts, fmt.Sprintf("%d:%d", steps, histogram[steps]))
	}
	return strings.Join(parts, ", ")
}

func fourDigits(n int) string {
	return fmt.Sprintf("%04d", n)
}

func equalInts(left []int, right []int) bool {
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

func allChainsEndAt(chains []Chain, target int) bool {
	for _, chain := range chains {
		if len(chain.Path) == 0 || chain.Path[len(chain.Path)-1] != target {
			return false
		}
	}
	return true
}

func countChecks(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
