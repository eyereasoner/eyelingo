// gray_code_counter.go
//
// Inspired by Eyeling's `examples/gray-code-counter.n3`.
//
// The example generates a cyclic reflected binary Gray counter and audits its
// one-bit transition property.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "gray_code_counter"

type Dataset struct {
	CaseName string `json:"caseName"`
	Question string `json:"question"`
	Bits     int    `json:"bits"`
	Steps    int    `json:"steps"`
	Expected struct {
		UniqueStates        int `json:"uniqueStates"`
		WrapHammingDistance int `json:"wrapHammingDistance"`
	} `json:"expected"`
}
type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	Sequence  []int
	Distances []int
	Checks    []Check
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
	seq := []int{}
	for i := 0; i < ds.Steps; i++ {
		seq = append(seq, i^(i>>1))
	}
	dists := []int{}
	for i := 0; i < len(seq); i++ {
		dists = append(dists, hamming(seq[i], seq[(i+1)%len(seq)]))
	}
	checks := []Check{
		{"C1", len(seq) == 16 && ds.Bits == 4, "16 states were generated for a 4-bit counter"},
		{"C2", unique(seq) == ds.Expected.UniqueStates, "all generated states are unique"},
		{"C3", maxInt(dists) == 1, "each adjacent transition flips exactly one bit"},
		{"C4", dists[len(dists)-1] == ds.Expected.WrapHammingDistance, "the final state wraps to the first with one bit flip"},
		{"C5", fmtBits(seq[0], ds.Bits) == "0000" && fmtBits(seq[1], ds.Bits) == "0001" && fmtBits(seq[2], ds.Bits) == "0011" && fmtBits(seq[3], ds.Bits) == "0010", "first four states match the reflected binary Gray-code prefix"},
		{"C6", seq[7] == (7 ^ (7 >> 1)), "the numeric generator is n xor (n >> 1)"},
	}
	return Analysis{seq, dists, checks}
}
func hamming(a, b int) int {
	x := a ^ b
	n := 0
	for x > 0 {
		n += x & 1
		x >>= 1
	}
	return n
}
func unique(xs []int) int {
	seen := map[int]bool{}
	for _, x := range xs {
		seen[x] = true
	}
	return len(seen)
}
func maxInt(xs []int) int {
	m := 0
	for _, x := range xs {
		if x > m {
			m = x
		}
	}
	return m
}
func fmtBits(x, bits int) string { return fmt.Sprintf("%0*b", bits, x) }
func prefix(seq []int, bits int, n int) string {
	out := []string{}
	for i := 0; i < n && i < len(seq); i++ {
		out = append(out, fmtBits(seq[i], bits))
	}
	return strings.Join(out, ", ")
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
	fmt.Println("# Gray Code Counter")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("bits : %d\n", ds.Bits)
	fmt.Printf("states visited : %d\n", len(a.Sequence))
	fmt.Printf("unique states : %d\n", unique(a.Sequence))
	fmt.Printf("sequence prefix : %s\n", prefix(a.Sequence, ds.Bits, 8))
	fmt.Printf("wrap transition : %s -> %s\n", fmtBits(a.Sequence[len(a.Sequence)-1], ds.Bits), fmtBits(a.Sequence[0], ds.Bits))
	fmt.Printf("maximum adjacent Hamming distance : %d\n", maxInt(a.Distances))
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The counter maps each integer n to n xor (n >> 1), which is the reflected binary Gray-code construction.")
	fmt.Println("For 4 bits, the first 16 integers cover the full state space without duplicates.")
	fmt.Println("The Hamming-distance check compares each state with the next state, including the final wraparound transition.")
	fmt.Println("A valid cyclic Gray counter therefore changes exactly one bit at every step.")
	fmt.Println()
	return
	for _, c := range a.Checks {
		status := "FAIL"
		if c.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", c.ID, status, c.Text)
	}
	fmt.Println()
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("bits : %d\n", ds.Bits)
	fmt.Printf("requested steps : %d\n", ds.Steps)
	fmt.Printf("unique states : %d\n", unique(a.Sequence))
	fmt.Printf("adjacent transitions checked : %d\n", len(a.Distances))
	fmt.Printf("wrap hamming distance : %d\n", a.Distances[len(a.Distances)-1])
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
