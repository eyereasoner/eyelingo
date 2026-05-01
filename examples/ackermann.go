// ackermann.go
//
// A self-contained Go translation of ackermann.n3 from the Eyeling examples.
//
// The original N3 file defines Ackermann answers by first calling a more
// general three-argument helper:
//
//	ackermann(x,y) = ackermann3(x, y+3, 2) - 3
//
// That helper covers a ladder of operations: successor, addition,
// multiplication, exponentiation, and then larger repeated-power operations
// such as tetration and pentation. The test query asks for twelve values,
// including exact integers with hundreds and tens of thousands of digits.
//
// This Go version keeps the reduction and recursive rules explicit. It uses
// math/big for exact arithmetic and reports large answers by decimal digit
// count plus SHA-256 fingerprint so the output stays readable and traceable.
//
// Run:
//
//	go run ackermann.go
//
// The program has no third-party dependencies.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"eyelingo/internal/exampleinput"
	"fmt"
	"math/big"
)

const eyelingoExampleName = "ackermann"

const (
	baseZ          = 2
	largeValueCut  = 80
	fingerprintLen = 24
)

type Query struct {
	ID string
	X  int
	Y  int64
}

type DerivedFact struct {
	Query         Query
	TernaryY      *big.Int
	TernaryResult *big.Int
	Answer        *big.Int
	Summary       ValueSummary
	RulePath      string
}

type ValueSummary struct {
	Digits int
	Exact  string
	Lead   string
	Tail   string
	SHA256 string
}

type EngineStats struct {
	Calls               int
	MemoHits            int
	ComputedRules       int
	SuccessorRules      int
	AdditionRules       int
	MultiplicationRules int
	PowerRules          int
	OneRules            int
	RecursiveRules      int
	MaxX                int
	MaxYDigits          int
	MaxResultDigits     int
}

type Analysis struct {
	Question          string
	Facts             []DerivedFact
	Stats             EngineStats
	Checks            []Check
	PrimitiveQueries  int
	BinaryReductions  int
	DistinctTernaries int
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type engine struct {
	memo  map[string]*big.Int
	stats EngineStats
}

func main() {
	analysis := derive()
	printAnswer(analysis)
	printReason(analysis)
}

func derive() Analysis {
	queries := exampleinput.Load(eyelingoExampleName, []Query{
		{ID: "A0", X: 0, Y: 0},
		{ID: "A1", X: 0, Y: 6},
		{ID: "A2", X: 1, Y: 2},
		{ID: "A3", X: 1, Y: 7},
		{ID: "A4", X: 2, Y: 2},
		{ID: "A5", X: 2, Y: 9},
		{ID: "A6", X: 3, Y: 4},
		{ID: "A7", X: 3, Y: 1000},
		{ID: "A8", X: 4, Y: 0},
		{ID: "A9", X: 4, Y: 1},
		{ID: "A10", X: 4, Y: 2},
		{ID: "A11", X: 5, Y: 0},
	})

	e := &engine{memo: make(map[string]*big.Int)}
	facts := make([]DerivedFact, 0, len(queries))
	for _, q := range queries {
		ternaryY := new(big.Int).Add(big.NewInt(q.Y), big.NewInt(3))
		ternary := e.ack3(q.X, ternaryY, baseZ)
		answer := new(big.Int).Sub(ternary, big.NewInt(3))
		facts = append(facts, DerivedFact{
			Query:         q,
			TernaryY:      ternaryY,
			TernaryResult: ternary,
			Answer:        answer,
			Summary:       summarize(answer),
			RulePath:      rulePath(q.X, ternaryY),
		})
	}

	analysis := Analysis{
		Question:          "Evaluate the Ackermann facts queried by ackermann.n3.",
		Facts:             facts,
		Stats:             e.stats,
		PrimitiveQueries:  len(queries),
		BinaryReductions:  len(queries),
		DistinctTernaries: len(e.memo),
	}
	analysis.Checks = buildChecks(analysis)
	return analysis
}

func (e *engine) ack3(x int, y *big.Int, z int) *big.Int {
	e.stats.Calls++
	if x > e.stats.MaxX {
		e.stats.MaxX = x
	}
	if digits := decimalDigits(y); digits > e.stats.MaxYDigits {
		e.stats.MaxYDigits = digits
	}

	key := fmt.Sprintf("%d|%s|%d", x, y.String(), z)
	if cached, ok := e.memo[key]; ok {
		e.stats.MemoHits++
		return new(big.Int).Set(cached)
	}

	var result *big.Int
	switch {
	case x == 0:
		e.stats.SuccessorRules++
		result = new(big.Int).Add(y, big.NewInt(1))
	case x == 1:
		e.stats.AdditionRules++
		result = new(big.Int).Add(y, big.NewInt(int64(z)))
	case x == 2:
		e.stats.MultiplicationRules++
		result = new(big.Int).Mul(y, big.NewInt(int64(z)))
	case x == 3:
		e.stats.PowerRules++
		result = new(big.Int).Exp(big.NewInt(int64(z)), y, nil)
	case y.Sign() == 0:
		e.stats.OneRules++
		result = big.NewInt(1)
	default:
		e.stats.RecursiveRules++
		previousY := new(big.Int).Sub(y, big.NewInt(1))
		inner := e.ack3(x, previousY, z)
		result = e.ack3(x-1, inner, z)
	}

	e.stats.ComputedRules++
	if digits := decimalDigits(result); digits > e.stats.MaxResultDigits {
		e.stats.MaxResultDigits = digits
	}
	e.memo[key] = new(big.Int).Set(result)
	return new(big.Int).Set(result)
}

func rulePath(x int, ternaryY *big.Int) string {
	switch x {
	case 0:
		return fmt.Sprintf("binary offset -> T(0,%s,2) -> successor", ternaryY)
	case 1:
		return fmt.Sprintf("binary offset -> T(1,%s,2) -> addition", ternaryY)
	case 2:
		return fmt.Sprintf("binary offset -> T(2,%s,2) -> multiplication", ternaryY)
	case 3:
		return fmt.Sprintf("binary offset -> T(3,%s,2) -> exponentiation", ternaryY)
	case 4:
		return fmt.Sprintf("binary offset -> T(4,%s,2) -> tetration recursion", ternaryY)
	default:
		return fmt.Sprintf("binary offset -> T(%d,%s,2) -> higher hyperoperation recursion", x, ternaryY)
	}
}

func summarize(n *big.Int) ValueSummary {
	text := n.String()
	sum := sha256.Sum256([]byte(text))
	out := ValueSummary{Digits: len(text), SHA256: hex.EncodeToString(sum[:])}
	if len(text) <= largeValueCut {
		out.Exact = text
		return out
	}
	out.Lead = text[:fingerprintLen]
	out.Tail = text[len(text)-fingerprintLen:]
	return out
}

func decimalDigits(n *big.Int) int {
	return len(n.String())
}

func buildChecks(a Analysis) []Check {
	lookup := make(map[string]*big.Int, len(a.Facts))
	for _, fact := range a.Facts {
		lookup[fact.Query.ID] = fact.Answer
	}

	exp2 := func(n int64) *big.Int {
		return new(big.Int).Exp(big.NewInt(2), big.NewInt(n), nil)
	}
	minus3 := func(n *big.Int) *big.Int {
		return new(big.Int).Sub(n, big.NewInt(3))
	}

	return []Check{
		{
			Label: "C1",
			OK:    equals(lookup["A0"], big.NewInt(1)) && equals(lookup["A1"], big.NewInt(7)),
			Text:  "x=0 reduces to successor after the y+3 binary offset.",
		},
		{
			Label: "C2",
			OK:    equals(lookup["A2"], big.NewInt(4)) && equals(lookup["A3"], big.NewInt(9)),
			Text:  "x=1 reduces to addition after the y+3 binary offset.",
		},
		{
			Label: "C3",
			OK:    equals(lookup["A4"], big.NewInt(7)) && equals(lookup["A5"], big.NewInt(21)),
			Text:  "x=2 reduces to multiplication after the y+3 binary offset.",
		},
		{
			Label: "C4",
			OK:    equals(lookup["A6"], big.NewInt(125)) && equals(lookup["A7"], minus3(exp2(1003))),
			Text:  "x=3 reduces to exact BigInt exponentiation, including 2^1003-3.",
		},
		{
			Label: "C5",
			OK:    equals(lookup["A8"], big.NewInt(13)) && equals(lookup["A9"], big.NewInt(65533)),
			Text:  "x=4 derives the first tetration cases T(4,3,2)-3 and T(4,4,2)-3.",
		},
		{
			Label: "C6",
			OK:    equals(lookup["A10"], minus3(exp2(65536))),
			Text:  "A(4,2) is held exactly as 2^65536-3, not as a floating-point approximation.",
		},
		{
			Label: "C7",
			OK:    equals(lookup["A11"], lookup["A9"]),
			Text:  "the pentation query A(5,0) lands on the same value as A(4,1).",
		},
		{
			Label: "C8",
			OK:    a.Stats.MaxResultDigits == 19729 && a.DistinctTernaries == a.Stats.ComputedRules,
			Text:  "the evaluator reached the expected largest exact integer and memoized each distinct ternary fact once.",
		},
	}
}

func equals(a, b *big.Int) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Cmp(b) == 0
}

func printAnswer(a Analysis) {
	fmt.Println("# Ackermann")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("The ackermann.n3 test query derives %d Ackermann facts.\n", len(a.Facts))
	fmt.Println("Computed values:")
	for _, fact := range a.Facts {
		fmt.Printf("%s ackermann(%d,%d) = %s\n", fact.Query.ID, fact.Query.X, fact.Query.Y, formatSummary(fact.Summary))
	}
	fmt.Println("Large exact-value fingerprints:")
	for _, fact := range a.Facts {
		if fact.Summary.Exact == "" {
			fmt.Printf("%s digits=%d leading=%s trailing=%s sha256=%s\n", fact.Query.ID, fact.Summary.Digits, fact.Summary.Lead, fact.Summary.Tail, fact.Summary.SHA256)
		}
	}
	fmt.Println()
}

func printReason(a Analysis) {
	fmt.Println("## Reason")
	fmt.Println("The N3 source defines binary ackermann(x,y) by computing T(x,y+3,2) and subtracting 3. The ternary predicate T uses direct rules for successor, addition, multiplication, and exponentiation, then uses the recursive hyperoperation rule T(x,y,z)=T(x-1,T(x,y-1,z),z) when x>3 and y is non-zero.")
	fmt.Printf("primitive test queries : %d\n", a.PrimitiveQueries)
	fmt.Printf("binary reductions : %d\n", a.BinaryReductions)
	fmt.Printf("distinct ternary facts : %d\n", a.DistinctTernaries)
	fmt.Printf("memo hits : %d\n", a.Stats.MemoHits)
	fmt.Println("rule paths:")
	for _, fact := range a.Facts {
		fmt.Printf("%s %s gives T=%s, answer=T-3=%s\n", fact.Query.ID, fact.RulePath, formatSummary(summarize(fact.TernaryResult)), formatSummary(fact.Summary))
	}
	fmt.Println("hyperoperation highlights:")
	fmt.Println("A7 is 2^1003 - 3, an exact 302-digit integer.")
	fmt.Println("A10 is 2^65536 - 3, an exact 19,729-digit integer summarized by fingerprint.")
	fmt.Println("A11 reuses the pentation step T(5,3,2)=T(4,4,2)=65536, so A11 equals A9.")
	fmt.Println()
}

func formatSummary(s ValueSummary) string {
	if s.Exact != "" {
		return s.Exact
	}
	return fmt.Sprintf("%d-digit integer [%s...%s; sha256=%s]", s.Digits, s.Lead, s.Tail, s.SHA256)
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

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}
