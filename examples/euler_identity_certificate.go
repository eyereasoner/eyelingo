// euler_identity_certificate.go
//
// Inspired by Eyeling's mathematical certificate examples such as
// `examples/euler-identity.n3`.
//
// The example verifies a finite residual certificate for exp(i*pi)+1.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"math/cmplx"
	"os"
)

const seeExampleName = "euler_identity_certificate"

type Dataset struct {
	CaseName  string  `json:"caseName"`
	Question  string  `json:"question"`
	Angle     float64 `json:"angle"`
	Tolerance float64 `json:"tolerance"`
	Terms     int     `json:"terms"`
	Expected  struct {
		ResidualBelowTolerance bool `json:"residualBelowTolerance"`
	} `json:"expected"`
}
type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	Value    complex128
	Residual float64
	Checks   []Check
}

func main() {
	ds := exampleinput.Load(seeExampleName, Dataset{})
	a := derive(ds)
	printReport(ds, a)
	if !allOK(a.Checks) {
		os.Exit(1)
	}
}
func derive(ds Dataset) Analysis {
	z := complex(0, ds.Angle)
	sum := complex(0, 0)
	power := complex(1, 0)
	factorial := 1.0
	for n := 0; n < ds.Terms; n++ {
		if n > 0 {
			power *= z
			factorial *= float64(n)
		}
		sum += power / complex(factorial, 0)
	}
	residual := cmplx.Abs(sum + 1)
	checks := []Check{
		{"C1", ds.Terms == 28, "Taylor expansion used 28 terms from the JSON input"},
		{"C2", real(sum) < -0.999999999999 && real(sum) > -1.000000000001, "computed real part is close to -1"},
		{"C3", imag(sum) < 1e-12 && imag(sum) > -1e-12, "computed imaginary part is close to 0"},
		{"C4", residual < ds.Tolerance && ds.Expected.ResidualBelowTolerance, "|exp(iπ)+1| is below the configured tolerance"},
		{"C5", residual > 0, "the audit records the finite residual rather than asserting exact real arithmetic"},
	}
	return Analysis{sum, residual, checks}
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
	fmt.Println("# Euler Identity Certificate")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Println("expression : exp(iπ) + 1")
	fmt.Printf("terms used : %d\n", ds.Terms)
	fmt.Printf("computed real part of exp(iπ) : %.15f\n", real(a.Value))
	fmt.Printf("computed imaginary part of exp(iπ) : %.15f\n", imag(a.Value))
	fmt.Printf("residual magnitude : %.3e\n", a.Residual)
	fmt.Printf("within tolerance : %t\n", a.Residual < ds.Tolerance)
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The example approximates exp(iπ) by a finite Taylor series over complex numbers.")
	fmt.Println("The resulting residual is not claimed to be mathematically exact zero; it is checked against the explicit tolerance from JSON.")
	fmt.Println("The computed real part is effectively -1 and the imaginary part is near 0 at the chosen precision.")
	fmt.Println("That gives a reproducible finite certificate for the familiar Euler-identity witness.")
	fmt.Println()
	return
}
