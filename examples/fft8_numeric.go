// fft8_numeric.go
//
// Inspired by Eyeling's `examples/fft8-numeric.n3`.
//
// dominant frequency bins of a single-cycle sine wave.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

const seeExampleName = "fft8_numeric"

type Dataset struct {
	CaseName string    `json:"caseName"`
	Question string    `json:"question"`
	Samples  []float64 `json:"samples"`
	Expected struct {
		DominantBins []int   `json:"dominantBins"`
		Tolerance    float64 `json:"tolerance"`
	} `json:"expected"`
}

type Bin struct {
	K         int
	Value     complex128
	Magnitude float64
	Phase     float64
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Bins         []Bin
	DominantBins []int
	EnergyTime   float64
	EnergyFreq   float64
	Checks       []Check
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
	n := len(ds.Samples)
	bins := make([]Bin, 0, n)
	for k := 0; k < n; k++ {
		var sum complex128
		for j, x := range ds.Samples {
			angle := -2 * math.Pi * float64(k*j) / float64(n)
			sum += complex(x, 0) * cmplx.Exp(complex(0, angle))
		}
		bins = append(bins, Bin{K: k, Value: sum, Magnitude: cmplx.Abs(sum), Phase: cmplx.Phase(sum)})
	}
	maxMag := 0.0
	for _, b := range bins {
		if b.Magnitude > maxMag {
			maxMag = b.Magnitude
		}
	}
	dominant := []int{}
	for _, b := range bins {
		if nearly(b.Magnitude, maxMag, ds.Expected.Tolerance) {
			dominant = append(dominant, b.K)
		}
	}
	energyTime := 0.0
	for _, x := range ds.Samples {
		energyTime += x * x
	}
	energyFreq := 0.0
	for _, b := range bins {
		energyFreq += b.Magnitude * b.Magnitude
	}
	energyFreq /= float64(n)
	checks := []Check{
		{"C1", n == 8, "the input contains exactly 8 time-domain samples"},
		{"C2", sameInts(dominant, ds.Expected.DominantBins), "the dominant bins are k=1 and k=7"},
		{"C3", nearly(bins[0].Magnitude, 0, ds.Expected.Tolerance), "the DC component is zero"},
		{"C4", nearly(bins[1].Magnitude, 4, 1e-9) && nearly(bins[7].Magnitude, 4, 1e-9), "the two conjugate sine bins have magnitude 4"},
		{"C5", nearly(energyTime, energyFreq, 1e-9), "Parseval energy is preserved by the unnormalized DFT convention"},
		{"C6", conjugatePair(bins[1].Value, bins[7].Value, 1e-9), "real-valued samples produce conjugate-symmetric bins"},
	}
	return Analysis{Bins: bins, DominantBins: dominant, EnergyTime: energyTime, EnergyFreq: energyFreq, Checks: checks}
}

func nearly(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

func sameInts(a, b []int) bool {
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

func conjugatePair(a, b complex128, tol float64) bool {
	return math.Abs(real(a)-real(b)) <= tol && math.Abs(imag(a)+imag(b)) <= tol
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

func sampleLine(samples []float64) string {
	parts := make([]string, 0, len(samples))
	for _, x := range samples {
		parts = append(parts, fmt.Sprintf("%.6f", x))
	}
	return strings.Join(parts, ", ")
}

func dominantLine(a Analysis) string {
	parts := []string{}
	for _, k := range a.DominantBins {
		b := a.Bins[k]
		parts = append(parts, fmt.Sprintf("k=%d magnitude=%.6f phase=%.6f", b.K, b.Magnitude, b.Phase))
	}
	return strings.Join(parts, "; ")
}

func printReport(ds Dataset, a Analysis) {
	fmt.Println("# FFT8 Numeric")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("sample vector : %s\n", sampleLine(ds.Samples))
	fmt.Printf("dominant bins : %s\n", dominantLine(a))
	fmt.Printf("time-domain energy : %.6f\n", a.EnergyTime)
	fmt.Printf("frequency-domain energy / 8 : %.6f\n", a.EnergyFreq)
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The input samples describe one sine cycle over eight equally spaced samples.")
	fmt.Println("The DFT projects the signal onto eight complex roots of unity.")
	fmt.Println("A real sine wave has equal magnitude at the positive and negative frequency bins.")
	fmt.Println("All non-dominant bins cancel to zero within the configured numerical tolerance.")
	fmt.Println()
	return
}
