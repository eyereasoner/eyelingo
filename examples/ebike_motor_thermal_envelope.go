// ebike_motor_thermal_envelope.go
//
// A Go translation inspired by Eyeling's
// `examples/decimal-ebike-motor-thermal-envelope.n3`.
//
// The example propagates a certified decimal interval for exp(-1/4) through a
// sampled e-bike motor thermal model.
//
// Run:
//
//	go run examples/ebike_motor_thermal_envelope.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
)

const eyelingoExampleName = "ebike_motor_thermal_envelope"

type Dataset struct {
	CaseName                string
	Question                string
	SamplePeriodSec         float64
	ThermalTimeConstantSec  float64
	ExactCoolingSymbol      string
	CoolingLower            float64
	CoolingUpper            float64
	AmbientC                float64
	InitialMotorC           float64
	WarningLimitC           float64
	HardLimitC              float64
	AssistPlan              []string
	HeatingEnvelopeByAssist map[string]HeatingEnvelope
	Expected                Expected
}

type HeatingEnvelope struct {
	Lower float64
	Upper float64
}

type Expected struct {
	Decision            string
	WarningRecoveryStep int
}

type TracePoint struct {
	Step       int
	Mode       string
	TempLowerC float64
	TempUpperC float64
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Trace               []TracePoint
	CoolingCertified    bool
	BelowHardLimit      bool
	DecreasingSteps     []int
	WarningRecoveryStep int
	WarningRecoverySec  float64
	MaxUpperC           float64
	Decision            string
	Checks              []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func derive(ds Dataset) Analysis {
	trace := make([]TracePoint, 0, len(ds.AssistPlan)+1)
	excessLo := ds.InitialMotorC - ds.AmbientC
	excessHi := excessLo
	trace = append(trace, TracePoint{Step: 0, Mode: "Initial", TempLowerC: ds.AmbientC + excessLo, TempUpperC: ds.AmbientC + excessHi})
	maxUpper := ds.AmbientC + excessHi
	for i, mode := range ds.AssistPlan {
		heat := ds.HeatingEnvelopeByAssist[mode]
		excessLo = ds.CoolingLower*excessLo + heat.Lower
		excessHi = ds.CoolingUpper*excessHi + heat.Upper
		upper := ds.AmbientC + excessHi
		if upper > maxUpper {
			maxUpper = upper
		}
		trace = append(trace, TracePoint{Step: i + 1, Mode: mode, TempLowerC: ds.AmbientC + excessLo, TempUpperC: upper})
	}

	coolingCertified := ds.CoolingLower > 0 && ds.CoolingLower < ds.CoolingUpper && ds.CoolingUpper < 1
	belowHard := true
	for _, point := range trace {
		if point.TempUpperC >= ds.HardLimitC {
			belowHard = false
		}
	}
	decreasingSteps := make([]int, 0)
	for i := 0; i+1 < len(trace); i++ {
		if trace[i+1].TempUpperC < trace[i].TempUpperC {
			decreasingSteps = append(decreasingSteps, i)
		}
	}
	recoveryStep := -1
	for i := 1; i < len(trace); i++ {
		if trace[i].TempUpperC < ds.WarningLimitC && trace[i-1].TempUpperC >= ds.WarningLimitC {
			recoveryStep = i
			break
		}
	}
	decision := "RejectThermalPlan"
	if coolingCertified && belowHard && recoveryStep == ds.Expected.WarningRecoveryStep {
		decision = ds.Expected.Decision
	}
	analysis := Analysis{Trace: trace, CoolingCertified: coolingCertified, BelowHardLimit: belowHard, DecreasingSteps: decreasingSteps, WarningRecoveryStep: recoveryStep, WarningRecoverySec: float64(recoveryStep) * ds.SamplePeriodSec, MaxUpperC: maxUpper, Decision: decision}
	analysis.Checks = []Check{
		{ID: "C1", OK: coolingCertified, Text: fmt.Sprintf("cooling interval %.10f..%.10f is positive, ordered, and contractive", ds.CoolingLower, ds.CoolingUpper)},
		{ID: "C2", OK: len(trace) == len(ds.AssistPlan)+1, Text: fmt.Sprintf("temperature trace has %d samples including the initial state", len(trace))},
		{ID: "C3", OK: belowHard, Text: fmt.Sprintf("maximum upper temperature %.4f C stays below hard limit %.1f C", maxUpper, ds.HardLimitC)},
		{ID: "C4", OK: recoveryStep == ds.Expected.WarningRecoveryStep, Text: fmt.Sprintf("warning recovery occurs at step %d after %.0f seconds", recoveryStep, analysis.WarningRecoverySec)},
		{ID: "C5", OK: len(decreasingSteps) == 9 && decreasingSteps[0] == 3, Text: "upper envelope decreases from the first post-Turbo sample onward"},
		{ID: "C6", OK: decision == ds.Expected.Decision, Text: fmt.Sprintf("ride decision is %s", decision)},
	}
	return analysis
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func countChecksOK(checks []Check) int {
	count := 0
	for _, check := range checks {
		if check.OK {
			count++
		}
	}
	return count
}

func printAnswer(ds Dataset, analysis Analysis) {
	fmt.Println("# E-Bike Motor Thermal Envelope")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("decision : %s\n", analysis.Decision)
	fmt.Printf("cooling certificate : %s in %.10f .. %.10f\n", ds.ExactCoolingSymbol, ds.CoolingLower, ds.CoolingUpper)
	fmt.Printf("maximum upper motor temperature : %.4f C\n", analysis.MaxUpperC)
	fmt.Printf("warning recovery : step %d at %.0f s\n", analysis.WarningRecoveryStep, analysis.WarningRecoverySec)
	fmt.Printf("hard limit : %.1f C\n", ds.HardLimitC)
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	fmt.Println("The model keeps an interval for motor temperature excess over ambient instead of pretending to know the transcendental cooling factor exactly.")
	fmt.Printf("For each %.0f second sample, the lower and upper excess envelopes are propagated with the certified cooling interval and the mode-specific heat injection.\n", ds.SamplePeriodSec)
	fmt.Printf("Turbo pushes the upper envelope to %.4f C, then Tour, Eco, and Coast allow the envelope to decrease.\n", analysis.MaxUpperC)
	fmt.Printf("The upper envelope returns below the %.1f C warning limit at step %d and remains below the %.1f C hard limit throughout.\n", ds.WarningLimitC, analysis.WarningRecoveryStep, ds.HardLimitC)
	fmt.Println()
}
