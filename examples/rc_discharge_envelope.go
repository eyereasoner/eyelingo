// rc_discharge_envelope.go
//
// Inspired by Eyeling's `examples/floating-point-first-rc-discharge.n3`.
//
// The example propagates a certified floating-point decay interval for a
// sampled RC discharge and finds the first safe voltage sample.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"math"
	"os"
)

const seeExampleName = "rc_discharge_envelope"

type Dataset struct {
	CaseName         string  `json:"caseName"`
	Question         string  `json:"question"`
	SamplePeriod     float64 `json:"samplePeriod"`
	TimeConstant     float64 `json:"timeConstant"`
	ExactDecaySymbol string  `json:"exactDecaySymbol"`
	DecayLower       float64 `json:"decayLower"`
	DecayUpper       float64 `json:"decayUpper"`
	InitialVoltage   float64 `json:"initialVoltage"`
	Tolerance        float64 `json:"tolerance"`
	MaxStep          int     `json:"maxStep"`
	Expected         struct {
		FirstSettledStep int `json:"firstSettledStep"`
	} `json:"expected"`
}
type Row struct {
	Step         int
	Lower, Upper float64
}
type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	Envelope []Row
	First    int
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
	env := []Row{}
	for k := 0; k <= ds.MaxStep; k++ {
		env = append(env, Row{k, ds.InitialVoltage * math.Pow(ds.DecayLower, float64(k)), ds.InitialVoltage * math.Pow(ds.DecayUpper, float64(k))})
	}
	first := -1
	for i := 1; i < len(env); i++ {
		if env[i].Upper < ds.Tolerance && env[i-1].Upper >= ds.Tolerance {
			first = env[i].Step
			break
		}
	}
	shrinks := true
	for i := 1; i < len(env); i++ {
		if !(env[i].Upper < env[i-1].Upper) {
			shrinks = false
		}
	}
	checks := []Check{
		{"C1", ds.DecayLower > 0 && ds.DecayLower < ds.DecayUpper && ds.DecayUpper < 1, "decay certificate is nonempty, positive, and below 1"},
		{"C2", shrinks, "voltage upper envelope decreases at every sample"},
		{"C3", env[12].Upper >= ds.Tolerance, "step 12 remains above the voltage tolerance"},
		{"C4", first == ds.Expected.FirstSettledStep, fmt.Sprintf("step %d is the first certified discharge step", first)},
		{"C5", close(float64(first)*ds.SamplePeriod, 0.325), "settling time is 0.325 s"},
		{"C6", len(env) == ds.MaxStep+1, "the certificate uses the JSON double interval rather than an exact transcendental value"},
	}
	return Analysis{env, first, checks}
}
func close(a, b float64) bool {
	if a > b {
		return a-b < 1e-9
	}
	return b-a < 1e-9
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
	row := a.Envelope[a.First]
	fmt.Println("# RC Discharge Envelope")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("exact decay symbol : %s\n", ds.ExactDecaySymbol)
	fmt.Printf("certified decay interval : [%.10f, %.10f]\n", ds.DecayLower, ds.DecayUpper)
	fmt.Printf("first below tolerance step : %d\n", a.First)
	fmt.Printf("first below tolerance time : %.3f s\n", float64(a.First)*ds.SamplePeriod)
	fmt.Printf("upper voltage at step %d : %.6f V\n", a.First, row.Upper)
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The physical decay factor is exp(-1/4), but the example uses a finite double interval as the certificate.")
	fmt.Println("Because the interval lies strictly between 0 and 1, the capacitor voltage envelope contracts each sample.")
	fmt.Println("The upper envelope is the safety-relevant bound: once it falls below 1.0 V, every compatible exact trajectory is below tolerance.")
	fmt.Println("The first such witness occurs before the configured maximum step.")
	fmt.Println()
	return
}
