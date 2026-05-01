// fft32_numeric.go
//
// Inspired by Eyeling's `examples/fft32-numeric.n3`.
//
// The example computes full 32-point Fourier spectra for several waveform
// fixtures and audits the dominant bins, flat-spectrum impulse behavior,
// conjugate symmetry, and Parseval energy preservation.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "fft32_numeric"

type Dataset struct {
	CaseName      string         `json:"caseName"`
	SourceExample string         `json:"sourceExample"`
	Question      string         `json:"question"`
	Length        int            `json:"length"`
	Tolerance     float64        `json:"tolerance"`
	Waveforms     []WaveformSpec `json:"waveforms"`
	ExpectedOps   int            `json:"expectedOperations"`
}

type WaveformSpec struct {
	Name                      string  `json:"name"`
	Kind                      string  `json:"kind"`
	Bin                       int     `json:"bin,omitempty"`
	ExpectedDominantBins      []int   `json:"expectedDominantBins"`
	ExpectedDominantMagnitude float64 `json:"expectedDominantMagnitude"`
	ExpectedFlatSpectrum      bool    `json:"expectedFlatSpectrum,omitempty"`
}

type Bin struct {
	K         int
	Value     complex128
	Magnitude float64
	Phase     float64
}

type Spectrum struct {
	Spec       WaveformSpec
	Samples    []float64
	Bins       []Bin
	Dominant   []int
	MaxMag     float64
	EnergyTime float64
	EnergyFreq float64
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Spectra    []Spectrum
	Operations int
	Checks     []Check
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
	spectra := make([]Spectrum, 0, len(ds.Waveforms))
	operations := 0
	for _, wf := range ds.Waveforms {
		samples := makeSamples(ds.Length, wf)
		bins := dft(samples)
		operations += len(samples) * len(samples)
		maxMag := maxMagnitude(bins)
		dominant := dominantBins(bins, maxMag, ds.Tolerance)
		energyTime := sampleEnergy(samples)
		energyFreq := spectrumEnergy(bins)
		spectra = append(spectra, Spectrum{Spec: wf, Samples: samples, Bins: bins, Dominant: dominant, MaxMag: maxMag, EnergyTime: energyTime, EnergyFreq: energyFreq})
	}

	checks := []Check{
		{"C1", allLength(spectra, ds.Length), "each generated waveform has exactly 32 samples"},
		{"C2", allDominantBinsMatch(spectra), "dominant bins match the FFT32 fixture expectations"},
		{"C3", allDominantMagnitudesMatch(spectra, ds.Tolerance), "dominant magnitudes match the configured certificates"},
		{"C4", allParseval(spectra, ds.Tolerance), "Parseval energy is preserved for every waveform under the unnormalized DFT convention"},
		{"C5", allConjugateSymmetric(spectra, ds.Tolerance), "real-valued waveforms produce conjugate-symmetric spectra"},
		{"C6", allFlatSpectraOK(spectra, ds.Tolerance), "the impulse waveform has a flat unit-magnitude spectrum"},
		{"C7", operations == ds.ExpectedOps, "the full spectrum is computed once per waveform"},
	}
	return Analysis{Spectra: spectra, Operations: operations, Checks: checks}
}

func makeSamples(n int, wf WaveformSpec) []float64 {
	samples := make([]float64, n)
	for j := 0; j < n; j++ {
		switch wf.Kind {
		case "alternating":
			if j%2 == 0 {
				samples[j] = 1
			} else {
				samples[j] = -1
			}
		case "constant":
			samples[j] = 1
		case "cosine":
			samples[j] = math.Cos(2 * math.Pi * float64(wf.Bin*j) / float64(n))
		case "impulse":
			if j == 0 {
				samples[j] = 1
			}
		case "ramp":
			samples[j] = float64(j)
		case "sine":
			samples[j] = math.Sin(2 * math.Pi * float64(wf.Bin*j) / float64(n))
		default:
			panic(fmt.Sprintf("unknown waveform kind %q", wf.Kind))
		}
	}
	return samples
}

func dft(samples []float64) []Bin {
	n := len(samples)
	bins := make([]Bin, 0, n)
	for k := 0; k < n; k++ {
		var sum complex128
		for j, x := range samples {
			angle := -2 * math.Pi * float64(k*j) / float64(n)
			sum += complex(x, 0) * cmplx.Exp(complex(0, angle))
		}
		bins = append(bins, Bin{K: k, Value: sum, Magnitude: cmplx.Abs(sum), Phase: cleanZero(cmplx.Phase(sum))})
	}
	return bins
}

func cleanZero(x float64) float64 {
	if math.Abs(x) < 5e-12 {
		return 0
	}
	return x
}

func maxMagnitude(bins []Bin) float64 {
	maxMag := 0.0
	for _, b := range bins {
		if b.Magnitude > maxMag {
			maxMag = b.Magnitude
		}
	}
	return maxMag
}

func dominantBins(bins []Bin, maxMag, tol float64) []int {
	dominant := []int{}
	for _, b := range bins {
		if nearly(b.Magnitude, maxMag, tol) {
			dominant = append(dominant, b.K)
		}
	}
	return dominant
}

func sampleEnergy(samples []float64) float64 {
	energy := 0.0
	for _, x := range samples {
		energy += x * x
	}
	return energy
}

func spectrumEnergy(bins []Bin) float64 {
	energy := 0.0
	for _, b := range bins {
		energy += b.Magnitude * b.Magnitude
	}
	return energy / float64(len(bins))
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

func allLength(spectra []Spectrum, n int) bool {
	if n != 32 || len(spectra) == 0 {
		return false
	}
	for _, s := range spectra {
		if len(s.Samples) != n || len(s.Bins) != n {
			return false
		}
	}
	return true
}

func allDominantBinsMatch(spectra []Spectrum) bool {
	for _, s := range spectra {
		if !sameInts(s.Dominant, s.Spec.ExpectedDominantBins) {
			return false
		}
	}
	return true
}

func allDominantMagnitudesMatch(spectra []Spectrum, tol float64) bool {
	for _, s := range spectra {
		if !nearly(s.MaxMag, s.Spec.ExpectedDominantMagnitude, tol) {
			return false
		}
	}
	return true
}

func allParseval(spectra []Spectrum, tol float64) bool {
	for _, s := range spectra {
		if !nearly(s.EnergyTime, s.EnergyFreq, tol) {
			return false
		}
	}
	return true
}

func allConjugateSymmetric(spectra []Spectrum, tol float64) bool {
	for _, s := range spectra {
		n := len(s.Bins)
		for k := 1; k < n; k++ {
			a := s.Bins[k].Value
			b := s.Bins[(n-k)%n].Value
			if math.Abs(real(a)-real(b)) > tol || math.Abs(imag(a)+imag(b)) > tol {
				return false
			}
		}
	}
	return true
}

func allFlatSpectraOK(spectra []Spectrum, tol float64) bool {
	for _, s := range spectra {
		if !s.Spec.ExpectedFlatSpectrum {
			continue
		}
		for _, b := range s.Bins {
			if !nearly(b.Magnitude, s.Spec.ExpectedDominantMagnitude, tol) {
				return false
			}
		}
	}
	return true
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

func binSummary(s Spectrum) string {
	if len(s.Dominant) == len(s.Bins) {
		return fmt.Sprintf("all %d bins magnitude=%.6f", len(s.Bins), s.MaxMag)
	}
	parts := make([]string, 0, len(s.Dominant))
	for _, k := range s.Dominant {
		b := s.Bins[k]
		parts = append(parts, fmt.Sprintf("k=%d magnitude=%.6f phase=%.6f", b.K, b.Magnitude, b.Phase))
	}
	return strings.Join(parts, "; ")
}

func printReport(ds Dataset, a Analysis) {
	fmt.Println("# FFT32 Numeric")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("waveforms evaluated : %d\n", len(a.Spectra))
	for _, s := range a.Spectra {
		fmt.Printf("%s : %s; energy=%.6f\n", s.Spec.Name, binSummary(s), s.EnergyTime)
	}
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The upstream FFT32 fixture defines several 32-sample waveforms and asks for the whole spectrum of each waveform.")
	fmt.Println("The Go translation evaluates every frequency bin by summing samples against the corresponding complex root of unity.")
	fmt.Println("Constant, alternating, cosine, sine, impulse, and ramp fixtures exercise different spectral shapes.")
	fmt.Println("The checks verify dominant bins, magnitudes, flat impulse spectrum, conjugate symmetry, and energy preservation.")
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
	fmt.Printf("source example : %s\n", ds.SourceExample)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("sample count per waveform : %d\n", ds.Length)
	fmt.Printf("waveform count : %d\n", len(ds.Waveforms))
	fmt.Printf("complex bin sums : %d\n", a.Operations)
	fmt.Printf("tolerance : %.1e\n", ds.Tolerance)
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
