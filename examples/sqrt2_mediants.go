// sqrt2_mediants.go
//
// Inspired by Eyeling's `examples/integer-first-sqrt2-mediants.n3`.
//
// The example builds an integer-certified rational bracket for sqrt(2) using
// convergents of the continued fraction [1; 2, 2, 2, ...].
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "sqrt2_mediants"

type Dataset struct {
	CaseName       string `json:"caseName"`
	Question       string `json:"question"`
	MaxDenominator int    `json:"maxDenominator"`
	Expected       struct {
		Lower string `json:"lower"`
		Upper string `json:"upper"`
	} `json:"expected"`
}

type Fraction struct {
	P int
	Q int
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Convergents []Fraction
	Lower       Fraction
	Upper       Fraction
	Checks      []Check
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
	convs := []Fraction{}
	p0, p1 := 0, 1
	q0, q1 := 1, 0
	for i := 0; ; i++ {
		a := 2
		if i == 0 {
			a = 1
		}
		p := a*p1 + p0
		q := a*q1 + q0
		if q > ds.MaxDenominator {
			break
		}
		convs = append(convs, Fraction{P: p, Q: q})
		p0, p1 = p1, p
		q0, q1 = q1, q
	}
	lower := Fraction{}
	upper := Fraction{}
	for _, f := range convs {
		cmp := compareToSqrt2(f)
		if cmp < 0 {
			lower = f
		} else if cmp > 0 {
			upper = f
		}
	}
	checks := []Check{
		{"C1", len(convs) == 9, "nine convergents stay within the denominator limit"},
		{"C2", lower.String() == ds.Expected.Lower, "the best lower bound is 1393/985"},
		{"C3", upper.String() == ds.Expected.Upper, "the best upper bound is 577/408"},
		{"C4", compareToSqrt2(lower) < 0 && compareToSqrt2(upper) > 0, "the bracket is certified by integer square comparisons"},
		{"C5", lower.P*upper.Q < upper.P*lower.Q, "the lower rational is strictly below the upper rational"},
		{"C6", upper.P*lower.Q-lower.P*upper.Q == 1, "the chosen bounds are adjacent Stern-Brocot neighbors"},
	}
	return Analysis{Convergents: convs, Lower: lower, Upper: upper, Checks: checks}
}

func compareToSqrt2(f Fraction) int {
	left := f.P * f.P
	right := 2 * f.Q * f.Q
	if left < right {
		return -1
	}
	if left > right {
		return 1
	}
	return 0
}

func (f Fraction) String() string { return fmt.Sprintf("%d/%d", f.P, f.Q) }

func decimal(f Fraction) float64 { return float64(f.P) / float64(f.Q) }

func width(lower, upper Fraction) float64 { return decimal(upper) - decimal(lower) }

func convergentLine(xs []Fraction) string {
	parts := make([]string, 0, len(xs))
	for _, f := range xs {
		parts = append(parts, f.String())
	}
	return strings.Join(parts, ", ")
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
	fmt.Println("# Integer-First Sqrt2 Mediants")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("lower bound : %s = %.12f\n", a.Lower.String(), decimal(a.Lower))
	fmt.Printf("upper bound : %s = %.12f\n", a.Upper.String(), decimal(a.Upper))
	fmt.Printf("certified interval width : %.12f\n", width(a.Lower, a.Upper))
	fmt.Printf("convergents used : %s\n", convergentLine(a.Convergents))
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The continued fraction for sqrt(2) is [1; 2, 2, 2, ...], so each convergent is derived by an integer recurrence.")
	fmt.Println("Each rational p/q is classified without floating-point square roots by comparing p*p with 2*q*q.")
	fmt.Println("The latest lower and upper convergents within the denominator limit form a tight certified bracket.")
	fmt.Println("Their cross-difference is one, so no simpler Stern-Brocot rational lies strictly between them.")
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
	fmt.Printf("max denominator : %d\n", ds.MaxDenominator)
	fmt.Printf("convergents generated : %d\n", len(a.Convergents))
	fmt.Printf("lower square comparison : %d^2 < 2*%d^2\n", a.Lower.P, a.Lower.Q)
	fmt.Printf("upper square comparison : %d^2 > 2*%d^2\n", a.Upper.P, a.Upper.Q)
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
