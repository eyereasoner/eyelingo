// goldbach_1000.go
//
// Inspired by Eyeling's `examples/goldbach-1000.n3`.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"os"
)

const seeExampleName = "goldbach_1000"

type Dataset struct {
	CaseName    string `json:"caseName"`
	Source      string `json:"source"`
	Question    string `json:"question"`
	MaxEven     int    `json:"maxEven"`
	SampleEvens []int  `json:"sampleEvens"`
}

type Witness struct {
	E int
	P int
	Q int
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Result struct {
	PrimeCount int
	EvenCount  int
	Failures   []int
	Samples    []Witness
	Checks     []Check
}

func main() {
	ds := exampleinput.Load(seeExampleName, Dataset{})
	r := derive(ds)
	printReport(ds, r)
	if !allOK(r.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Result {
	primes := map[int]bool{}
	for n := 2; n <= ds.MaxEven; n++ {
		if isPrime(n) {
			primes[n] = true
		}
	}

	r := Result{PrimeCount: len(primes)}
	for e := 4; e <= ds.MaxEven; e += 2 {
		r.EvenCount++
		if _, ok := firstWitness(e, primes); !ok {
			r.Failures = append(r.Failures, e)
		}
	}
	for _, e := range ds.SampleEvens {
		if w, ok := firstWitness(e, primes); ok {
			r.Samples = append(r.Samples, w)
		}
	}

	r.Checks = []Check{
		{"C1", ds.MaxEven == 1000, "the configured upper bound is 1000"},
		{"C2", r.EvenCount == 499, "there are 499 even integers from 4 through 1000"},
		{"C3", len(r.Failures) == 0, "every checked even integer has a prime-pair witness"},
		{"C4", len(r.Samples) == len(ds.SampleEvens), "each requested sample even has a witness"},
		{"C5", r.PrimeCount == 168, "there are 168 primes at or below 1000"},
	}
	return r
}

func printReport(ds Dataset, r Result) {
	fmt.Println("# Goldbach 1000")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("All %d even integers from 4 through %d have a Goldbach witness.\n", r.EvenCount, ds.MaxEven)
	fmt.Print("sample witnesses : ")
	for i, w := range r.Samples {
		if i > 0 {
			fmt.Print("; ")
		}
		fmt.Printf("%d=%d+%d", w.E, w.P, w.Q)
	}
	fmt.Println()
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The checker caches primes up to the configured bound and then searches each even number E for a prime P not greater than E/2 where E-P is also prime.")
	fmt.Println("No counterexample is found in the bounded range, so the bounded Goldbach condition succeeds for this dataset.")
	fmt.Println()
	return
}

func firstWitness(e int, primes map[int]bool) (Witness, bool) {
	for p := 2; p <= e/2; p++ {
		if primes[p] && primes[e-p] {
			return Witness{E: e, P: p, Q: e - p}, true
		}
	}
	return Witness{}, false
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for d := 3; d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func allOK(cs []Check) bool {
	for _, c := range cs {
		if !c.OK {
			return false
		}
	}
	return true
}

func countOK(cs []Check) int {
	n := 0
	for _, c := range cs {
		if c.OK {
			n++
		}
	}
	return n
}
