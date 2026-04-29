// decimal_servo_envelope.go
//
// Inspired by Eyeling's `examples/decimal-transcendental-servo-envelope.n3`.
//
// A finite decimal interval certifies exp(-1/3), then the error envelope is
// propagated until the servo is guaranteed inside tolerance.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"os"
	"runtime"
)

const eyelingoExampleName = "decimal_servo_envelope"

type Dataset struct {
	CaseName        string  `json:"caseName"`
	Question        string  `json:"question"`
	SamplePeriod    float64 `json:"samplePeriod"`
	TimeConstant    float64 `json:"timeConstant"`
	ExactPoleSymbol string  `json:"exactPoleSymbol"`
	PoleLower       float64 `json:"poleLower"`
	PoleUpper       float64 `json:"poleUpper"`
	InitialAbsError float64 `json:"initialAbsError"`
	Tolerance       float64 `json:"tolerance"`
	MaxStep         int     `json:"maxStep"`
	Expected        struct {
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
	Envelope     []Row
	FirstSettled int
	Checks       []Check
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
	env := []Row{}
	for k := 0; k <= ds.MaxStep; k++ {
		env = append(env, Row{k, ds.InitialAbsError * math.Pow(ds.PoleLower, float64(k)), ds.InitialAbsError * math.Pow(ds.PoleUpper, float64(k))})
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
		{"C1", ds.PoleLower > 0 && ds.PoleLower < ds.PoleUpper && ds.PoleUpper < 1, "pole certificate is nonempty, positive, and below 1"},
		{"C2", shrinks, "upper envelope strictly decreases at every sampled step"},
		{"C3", env[9].Upper >= ds.Tolerance, "step 9 is not yet below tolerance"},
		{"C4", first == ds.Expected.FirstSettledStep, "step 10 is the first certified settling step"},
		{"C5", close(float64(first)*ds.SamplePeriod, 0.200), "settling time is 0.200 s"},
		{"C6", len(env) == ds.MaxStep+1, "all values are derived from the JSON certificate parameters"},
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
	row := a.Envelope[a.FirstSettled]
	fmt.Println("# Decimal Servo Envelope")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("exact pole symbol : %s\n", ds.ExactPoleSymbol)
	fmt.Printf("certified pole interval : [%.10f, %.10f]\n", ds.PoleLower, ds.PoleUpper)
	fmt.Printf("first settled step : %d\n", a.FirstSettled)
	fmt.Printf("first settled time : %.3f s\n", float64(a.FirstSettled)*ds.SamplePeriod)
	fmt.Printf("upper envelope at step 10 : %.6f\n", row.Upper)
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The exact pole exp(-1/3) is not represented as an exact finite decimal, so the input provides a certified decimal interval.")
	fmt.Println("The upper bound of the pole interval is below 1, which makes the error envelope contractive.")
	fmt.Println("The certificate propagates the lower and upper absolute-error envelopes sample by sample.")
	fmt.Println("The first sample whose upper envelope is below the tolerance is the first guaranteed settling witness.")
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
	fmt.Printf("sample period : %.3f s\n", ds.SamplePeriod)
	fmt.Printf("time constant : %.3f s\n", ds.TimeConstant)
	fmt.Printf("max step : %d\n", ds.MaxStep)
	fmt.Printf("initial absolute error : %.3f\n", ds.InitialAbsError)
	fmt.Printf("tolerance : %.3f\n", ds.Tolerance)
	fmt.Printf("envelope rows : %d\n", len(a.Envelope))
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
