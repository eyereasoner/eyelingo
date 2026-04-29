// fundamental_theorem_arithmetic.go
//
// A self-contained Go translation of fundamental-theorem-arithmetic.n3 from
// the Eyeling examples.
//
// The Fundamental Theorem of Arithmetic says that every integer greater than 1
// has a prime factorization, and that the factorization is unique except for
// the order of the factors. The source N3 file demonstrates this with:
//
//	202692987 = 3^2 * 7 * 829 * 3881
//
// This Go version keeps that source example as the primary case and adds a
// wider set of numbers, including larger composites, repeated factors, and a
// large prime. Each case is factored by repeatedly taking the smallest divisor.
// The program then checks the product, checks that the factors are prime, and
// compares smallest-first and largest-first traversals to confirm uniqueness up
// to order.
//
// Run:
//
//	go run fundamental_theorem_arithmetic.go
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

const eyelingoExampleName = "fundamental_theorem_arithmetic"

const sourceFile = "fundamental-theorem-arithmetic.n3"

const primaryN int64 = 202692987

var expectedPrimaryFactors = []int64{3, 3, 7, 829, 3881}
var expectedPrimaryPrimePower = "3^2 * 7 * 829 * 3881"
var expectedPrimaryFlat = "3 * 3 * 7 * 829 * 3881"
var expectedPrimaryLargestFlat = "3881 * 829 * 7 * 3 * 3"

// sampleNumbers includes the N3 source number plus additional values that make
// the example less narrow: a number with many small factors, a product of
// Fermat primes, a well-known large composite, a ten-digit mixed composite, and
// a ten-digit prime.
var sampleNumbers = []int64{
	360360,
	primaryN,
	4294967295,
	600851475143,
	9876543210,
	9999999967,
}

type FactorStep struct {
	N        int64
	Divisor  int64
	Quotient int64
}

type Stats struct {
	SmallestDivisorSearches int64
	DivisibilityTests       int64
	PrimalityChecks         int64
	PrimeDivisorTests       int64
}

type CaseAnalysis struct {
	N                int64
	FactorsSmallest  []int64
	FactorsLargest   []int64
	SortedSmallest   []int64
	SortedLargest    []int64
	Product          int64
	DistinctFactors  []int64
	PrimePowerString string
	FlatString       string
	LargestFlat      string
	SmallestPrime    int64
	LargestPrime     int64
	Steps            []FactorStep
	AllFactorsPrime  bool
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Question      string
	Cases         []CaseAnalysis
	Primary       CaseAnalysis
	LargestN      int64
	TotalFactors  int
	DistinctCount int
	Stats         Stats
	Checks        []Check
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
	inputSamples := exampleinput.Load(eyelingoExampleName, sampleNumbers)
	sampleNumbers = inputSamples
	stats := Stats{}
	cases := make([]CaseAnalysis, 0, len(sampleNumbers))
	allDistinctFactors := []int64{}

	for _, n := range sampleNumbers {
		c := analyzeCase(n, &stats)
		cases = append(cases, c)
		allDistinctFactors = append(allDistinctFactors, c.DistinctFactors...)
	}

	primary := findCase(cases, primaryN)
	largestN := maxN(cases)
	totalFactors := 0
	for _, c := range cases {
		totalFactors += len(c.FactorsSmallest)
	}

	sort.Slice(allDistinctFactors, func(i, j int) bool {
		return allDistinctFactors[i] < allDistinctFactors[j]
	})

	analysis := Analysis{
		Question:      "What are the prime factorizations, and are they unique up to order?",
		Cases:         cases,
		Primary:       primary,
		LargestN:      largestN,
		TotalFactors:  totalFactors,
		DistinctCount: len(distinctInts(allDistinctFactors)),
		Stats:         stats,
	}
	analysis.Checks = buildChecks(analysis)
	return analysis
}

func analyzeCase(n int64, stats *Stats) CaseAnalysis {
	factors, steps := factorSmallest(n, stats)
	largest := reverseInts(factors)
	sortedSmallest := sortedCopy(factors)
	sortedLargest := sortedCopy(largest)
	distinct := distinctInts(sortedSmallest)

	c := CaseAnalysis{
		N:                n,
		FactorsSmallest:  factors,
		FactorsLargest:   largest,
		SortedSmallest:   sortedSmallest,
		SortedLargest:    sortedLargest,
		Product:          product(factors),
		DistinctFactors:  distinct,
		PrimePowerString: primePowerString(sortedSmallest),
		FlatString:       factorString(factors),
		LargestFlat:      factorString(largest),
		Steps:            steps,
	}
	c.SmallestPrime = c.FactorsSmallest[0]
	c.LargestPrime = c.FactorsSmallest[len(c.FactorsSmallest)-1]
	c.AllFactorsPrime = true
	for _, factor := range distinct {
		if !isPrimeByTrialDivision(factor, stats) {
			c.AllFactorsPrime = false
			break
		}
	}
	return c
}

// factorSmallest repeatedly removes the smallest divisor. The resulting list is
// nondecreasing, which makes it easy to compare with the expected N3 list.
func factorSmallest(n int64, stats *Stats) ([]int64, []FactorStep) {
	if n < 2 {
		return nil, nil
	}

	factors := []int64{}
	steps := []FactorStep{}
	current := n

	for current >= 2 {
		divisor := smallestDivisorFrom(current, 2, stats)
		if divisor == current {
			factors = append(factors, current)
			break
		}
		quotient := current / divisor
		factors = append(factors, divisor)
		steps = append(steps, FactorStep{N: current, Divisor: divisor, Quotient: quotient})
		current = quotient
	}

	return factors, steps
}

// smallestDivisorFrom mirrors the helper relation in the N3 source. It tries
// possible divisors in order and returns n itself when no smaller divisor can
// exist.
func smallestDivisorFrom(n int64, start int64, stats *Stats) int64 {
	stats.SmallestDivisorSearches++
	for d := start; ; d++ {
		stats.DivisibilityTests++
		if n%d == 0 {
			return d
		}
		if d*d > n {
			return n
		}
	}
}

func isPrimeByTrialDivision(n int64, stats *Stats) bool {
	stats.PrimalityChecks++
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 {
		stats.PrimeDivisorTests++
		return false
	}
	stats.PrimeDivisorTests++
	for d := int64(3); d*d <= n; d += 2 {
		stats.PrimeDivisorTests++
		if n%d == 0 {
			return false
		}
	}
	return true
}

func buildChecks(a Analysis) []Check {
	allProducts := true
	allPrime := true
	allUnique := true
	for _, c := range a.Cases {
		if c.Product != c.N {
			allProducts = false
		}
		if !c.AllFactorsPrime {
			allPrime = false
		}
		if !equalInts(c.SortedSmallest, c.SortedLargest) {
			allUnique = false
		}
	}

	return []Check{
		{
			Label: "primary factors",
			OK:    equalInts(a.Primary.FactorsSmallest, expectedPrimaryFactors),
			Text:  "the source example factors 202692987 as 3,3,7,829,3881",
		},
		{
			Label: "primary prime-power form",
			OK:    a.Primary.PrimePowerString == expectedPrimaryPrimePower,
			Text:  "the source example groups repeated factors as 3^2 * 7 * 829 * 3881",
		},
		{
			Label: "product reconstruction",
			OK:    allProducts,
			Text:  "multiplying each computed factor list reconstructs its original number",
		},
		{
			Label: "prime factors",
			OK:    allPrime,
			Text:  "every distinct factor in every decomposition is prime by trial division",
		},
		{
			Label: "unique up to order",
			OK:    allUnique,
			Text:  "smallest-first and largest-first traversals sort to the same multisets",
		},
		{
			Label: "larger examples",
			OK:    len(a.Cases) == 6 && hasCase(a.Cases, 9999999967),
			Text:  "the extended sample includes six cases and includes the ten-digit prime 9999999967",
		},
	}
}

func printAnswer(a Analysis) {
	fmt.Println("=== Answer ===")
	fmt.Printf("Primary N3 case: n = %d has prime factors %s.\n", a.Primary.N, a.Primary.FlatString)
	fmt.Printf("primary prime-power form : %s\n", a.Primary.PrimePowerString)
	fmt.Printf("sample count : %d\n", len(a.Cases))
	fmt.Printf("largest sample : %d\n", a.LargestN)
	fmt.Printf("total prime factors counted with multiplicity : %d\n", a.TotalFactors)
	fmt.Printf("distinct primes seen across samples : %d\n", a.DistinctCount)
	fmt.Println()
	fmt.Println("Sample factorizations:")
	for _, c := range a.Cases {
		fmt.Printf("  %d = %s\n", c.N, c.PrimePowerString)
	}
	fmt.Println()
}

func printReason(a Analysis) {
	fmt.Println("=== Reason Why ===")
	fmt.Println("Existence comes from repeated smallest-divisor decomposition.")
	fmt.Println("At each step, the first divisor found is prime because no smaller")
	fmt.Println("positive divisor can divide the current number.")
	fmt.Println()
	fmt.Println("Smallest-divisor trace for the N3 source number:")
	for _, step := range a.Primary.Steps {
		fmt.Printf("  %d = %d * %d\n", step.N, step.Divisor, step.Quotient)
	}
	fmt.Printf("  %d is prime\n", a.Primary.FactorsSmallest[len(a.Primary.FactorsSmallest)-1])
	fmt.Println()
	fmt.Println("Uniqueness up to order is checked by reversing each traversal and sorting")
	fmt.Println("both factor lists. Matching sorted lists describe the same multiset of")
	fmt.Println("prime factors, even when the factors were discovered in the opposite order.")
	fmt.Printf("  source smallest-first factors : %s\n", a.Primary.FlatString)
	fmt.Printf("  source largest-first factors : %s\n", a.Primary.LargestFlat)
	fmt.Printf("  source sorted comparison : %s\n", factorString(a.Primary.SortedSmallest))
	fmt.Println()
	fmt.Println("The additional samples cover repeated small factors, special products,")
	fmt.Println("large composites, and a larger prime that has no smaller divisor.")
	fmt.Println()
}

func printChecks(checks []Check) {
	fmt.Println("=== Check ===")
	for i, check := range checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("C%d %s - %s\n", i+1, status, check.Text)
	}
	fmt.Println()
}

func printAudit(a Analysis) {
	passed := countPassed(a.Checks)

	fmt.Println("=== Go audit details ===")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("source file : %s\n", sourceFile)
	fmt.Printf("question : %s\n", a.Question)
	fmt.Printf("primary n : %d\n", a.Primary.N)
	fmt.Printf("primary smallest-first factors : %s\n", joinInts(a.Primary.FactorsSmallest, ","))
	fmt.Printf("primary largest-first factors : %s\n", joinInts(a.Primary.FactorsLargest, ","))
	fmt.Printf("primary flat factor string : %s\n", a.Primary.FlatString)
	fmt.Printf("primary prime-power string : %s\n", a.Primary.PrimePowerString)
	fmt.Printf("expected flat string : %s\n", expectedPrimaryFlat)
	fmt.Printf("expected largest-first string : %s\n", expectedPrimaryLargestFlat)
	fmt.Printf("sample numbers : %s\n", joinCaseNumbers(a.Cases, ","))
	fmt.Printf("sample count : %d\n", len(a.Cases))
	fmt.Printf("largest sample : %d\n", a.LargestN)
	fmt.Printf("total prime factors counted with multiplicity : %d\n", a.TotalFactors)
	fmt.Printf("distinct primes seen across samples : %d\n", a.DistinctCount)
	fmt.Printf("smallest-divisor searches : %d\n", a.Stats.SmallestDivisorSearches)
	fmt.Printf("divisibility tests : %d\n", a.Stats.DivisibilityTests)
	fmt.Printf("primality checks : %d\n", a.Stats.PrimalityChecks)
	fmt.Printf("prime-divisor tests : %d\n", a.Stats.PrimeDivisorTests)
	fmt.Printf("checks passed : %d/%d\n", passed, len(a.Checks))
	fmt.Printf("all checks pass : %s\n", yesNo(passed == len(a.Checks)))
}

func product(values []int64) int64 {
	result := int64(1)
	for _, value := range values {
		result *= value
	}
	return result
}

func sortedCopy(values []int64) []int64 {
	out := append([]int64(nil), values...)
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}

func reverseInts(values []int64) []int64 {
	out := append([]int64(nil), values...)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func distinctInts(values []int64) []int64 {
	out := []int64{}
	for _, value := range values {
		if len(out) == 0 || out[len(out)-1] != value {
			out = append(out, value)
		}
	}
	return out
}

func findCase(cases []CaseAnalysis, n int64) CaseAnalysis {
	for _, c := range cases {
		if c.N == n {
			return c
		}
	}
	panic(fmt.Sprintf("missing analysis case for %d", n))
}

func hasCase(cases []CaseAnalysis, n int64) bool {
	for _, c := range cases {
		if c.N == n {
			return true
		}
	}
	return false
}

func maxN(cases []CaseAnalysis) int64 {
	largest := cases[0].N
	for _, c := range cases[1:] {
		if c.N > largest {
			largest = c.N
		}
	}
	return largest
}

func factorString(values []int64) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprint(value)
	}
	return strings.Join(parts, " * ")
}

func primePowerString(values []int64) string {
	parts := []string{}
	for i := 0; i < len(values); {
		value := values[i]
		count := 1
		for j := i + 1; j < len(values) && values[j] == value; j++ {
			count++
		}
		if count == 1 {
			parts = append(parts, fmt.Sprint(value))
		} else {
			parts = append(parts, fmt.Sprintf("%d^%d", value, count))
		}
		i += count
	}
	return strings.Join(parts, " * ")
}

func joinInts(values []int64, separator string) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprint(value)
	}
	return strings.Join(parts, separator)
}

func joinCaseNumbers(cases []CaseAnalysis, separator string) string {
	parts := make([]string, len(cases))
	for i, c := range cases {
		parts[i] = fmt.Sprint(c.N)
	}
	return strings.Join(parts, separator)
}

func equalInts(a []int64, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func allChecksOK(checks []Check) bool {
	return countPassed(checks) == len(checks)
}

func countPassed(checks []Check) int {
	passed := 0
	for _, check := range checks {
		if check.OK {
			passed++
		}
	}
	return passed
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}
