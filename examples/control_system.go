// control_system.go
//
// A self-contained Go translation of examples/control-system.n3 from the
// Eyeling example suite, with the rule structure kept close to
// reasoning/control-system/rules-001.n3 from EYE.
//
// The original N3 program has input/disturbance/state/output facts, two
// forward control rules, and two backward rules that derive :measurement10
// from :measurement1 pairs. This translation keeps the concrete data in
// examples/input/control_system.json and represents the rules as ordinary Go
// functions.
//
// Run:
//
//     go run examples/control_system.go
//
// The program has no third-party dependencies.

package main

import (
	"see/internal/exampleinput"
	"fmt"
	"math"
	"os"
)

const seeExampleName = "control_system"

// ---------- types ----------

type Measurement1Fact struct {
	Measurement1 [2]float64 `json:"measurement1"`
}

type Measurement2Fact struct {
	Measurement2 bool `json:"measurement2"`
}

type Measurement3Fact struct {
	Measurement3 float64 `json:"measurement3"`
}

type Observation1Fact struct {
	Observation1 float64 `json:"observation1"`
}

type Observation2Fact struct {
	Observation2 bool `json:"observation2"`
}

type Observation3Fact struct {
	Observation3 float64 `json:"observation3"`
}

type Output2Fact struct {
	Measurement4 float64 `json:"measurement4"`
	Target2      float64 `json:"target2"`
}

// Input stores the editable fact block from control-system.n3.
type Input struct {
	Input1       Measurement1Fact `json:"input1"`
	Input2       Measurement2Fact `json:"input2"`
	Input3       Measurement3Fact `json:"input3"`
	Disturbance1 Measurement3Fact `json:"disturbance1"`
	Disturbance2 Measurement1Fact `json:"disturbance2"`
	State1       Observation1Fact `json:"state1"`
	State2       Observation2Fact `json:"state2"`
	State3       Observation3Fact `json:"state3"`
	Output2      Output2Fact      `json:"output2"`
}

type Measurement10Derivation struct {
	Subject string
	Input   [2]float64
	Value   float64
	Branch  string
	Reason  string
}

type FeedForwardControl struct {
	InputMeasurement10 float64
	Disturbance3       float64
	ProductPart        float64
	CompensationPart   float64
	Control            float64
}

type PNDFeedbackControl struct {
	Input3Measurement3 float64
	StateObservation3  float64
	OutputMeasurement4 float64
	OutputTarget2      float64
	Error              float64
	DifferentialError  float64
	ProportionalPart   float64
	NonlinearFactor    float64
	NonlinearPart      float64
	Control            float64
}

type Checks struct {
	C1BackwardLessThan      string
	C2BackwardNotLessThan   string
	C3FeedForwardGuard      string
	C4FeedForwardArithmetic string
	C5FeedbackError         string
	C6DifferentialError     string
	C7NonlinearPart         string
	C8FeedbackControlSum    string
}

type InferenceResult struct {
	Input1M10   Measurement10Derivation
	Disturb2M10 Measurement10Derivation
	FeedForward FeedForwardControl
	PNDFeedback PNDFeedbackControl
	Checks      Checks
}

// ---------- dataset ----------

func inputData() Input {
	return Input{
		Input1:       Measurement1Fact{Measurement1: [2]float64{6, 11}},
		Input2:       Measurement2Fact{Measurement2: true},
		Input3:       Measurement3Fact{Measurement3: 56967},
		Disturbance1: Measurement3Fact{Measurement3: 35766},
		Disturbance2: Measurement1Fact{Measurement1: [2]float64{45, 39}},
		State1:       Observation1Fact{Observation1: 80},
		State2:       Observation2Fact{Observation2: false},
		State3:       Observation3Fact{Observation3: 22},
		Output2:      Output2Fact{Measurement4: 24, Target2: 29},
	}
}

// ---------- inference ----------

func infer(in Input) (InferenceResult, error) {
	input1M10, err := deriveMeasurement10("input1", in.Input1.Measurement1)
	if err != nil {
		return InferenceResult{}, err
	}
	disturb2M10, err := deriveMeasurement10("disturbance2", in.Disturbance2.Measurement1)
	if err != nil {
		return InferenceResult{}, err
	}

	feedForward, err := applyFeedForwardRule(in, input1M10.Value)
	if err != nil {
		return InferenceResult{}, err
	}
	pndFeedback, err := applyPNDFeedbackRule(in)
	if err != nil {
		return InferenceResult{}, err
	}

	result := InferenceResult{
		Input1M10:   input1M10,
		Disturb2M10: disturb2M10,
		FeedForward: feedForward,
		PNDFeedback: pndFeedback,
	}
	result.Checks = performChecks(in, result)
	return result, nil
}

// deriveMeasurement10 implements the two backward rules from rules-001.n3:
//
//	{ ?I :measurement10 ?M. } <= { ?I :measurement1 (?M1 ?M2). ?M1 math:lessThan ?M2.
//	                               (?M2 ?M1) math:difference ?M3.
//	                               (?M3 0.5) math:exponentiation ?M. }.
//	{ ?I :measurement10 ?M1. } <= { ?I :measurement1 (?M1 ?M2). ?M1 math:notLessThan ?M2. }.
func deriveMeasurement10(subject string, pair [2]float64) (Measurement10Derivation, error) {
	m1 := pair[0]
	m2 := pair[1]
	if m1 < m2 {
		difference := m2 - m1
		if difference < 0 {
			return Measurement10Derivation{}, fmt.Errorf("%s has negative difference", subject)
		}
		value := math.Pow(difference, 0.5)
		return Measurement10Derivation{
			Subject: subject,
			Input:   pair,
			Value:   value,
			Branch:  "lessThan",
			Reason:  "measurement10 = (second measurement1 value - first measurement1 value)^0.5",
		}, nil
	}

	return Measurement10Derivation{
		Subject: subject,
		Input:   pair,
		Value:   m1,
		Branch:  "notLessThan",
		Reason:  "measurement10 = first measurement1 value",
	}, nil
}

// applyFeedForwardRule implements the first forward rule from rules-001.n3:
// input1 measurement10, input2 measurement2 true, disturbance1 measurement3,
// product, exponentiation, difference -> actuator1 control1.
func applyFeedForwardRule(in Input, measurement10 float64) (FeedForwardControl, error) {
	if !in.Input2.Measurement2 {
		return FeedForwardControl{}, fmt.Errorf("feedforward guard failed: input2 measurement2 is false")
	}
	if in.Disturbance1.Measurement3 <= 0 {
		return FeedForwardControl{}, fmt.Errorf("disturbance1 measurement3 must be positive")
	}

	productPart := measurement10 * 19.6
	compensationPart := math.Log10(in.Disturbance1.Measurement3)
	control := productPart - compensationPart

	return FeedForwardControl{
		InputMeasurement10: measurement10,
		Disturbance3:       in.Disturbance1.Measurement3,
		ProductPart:        productPart,
		CompensationPart:   compensationPart,
		Control:            control,
	}, nil
}

// applyPNDFeedbackRule implements the second forward rule from rules-001.n3:
// input3/state3/output2 facts, target error, differential error,
// proportional term, nonlinear factor, nonlinear differential term, sum ->
// actuator2 control1.
func applyPNDFeedbackRule(in Input) (PNDFeedbackControl, error) {
	measurement4 := in.Output2.Measurement4
	target2 := in.Output2.Target2
	errorTerm := target2 - measurement4
	if errorTerm == 0 {
		return PNDFeedbackControl{}, fmt.Errorf("target and measurement are equal; quotient term would divide by zero")
	}

	differentialError := in.State3.Observation3 - measurement4
	proportionalPart := 5.8 * errorTerm
	nonlinearFactor := 7.3 / errorTerm
	nonlinearPart := nonlinearFactor * differentialError
	control := proportionalPart + nonlinearPart

	return PNDFeedbackControl{
		Input3Measurement3: in.Input3.Measurement3,
		StateObservation3:  in.State3.Observation3,
		OutputMeasurement4: measurement4,
		OutputTarget2:      target2,
		Error:              errorTerm,
		DifferentialError:  differentialError,
		ProportionalPart:   proportionalPart,
		NonlinearFactor:    nonlinearFactor,
		NonlinearPart:      nonlinearPart,
		Control:            control,
	}, nil
}

func performChecks(in Input, result InferenceResult) Checks {
	const tolerance = 1e-9
	var checks Checks

	if result.Input1M10.Branch == "lessThan" && nearlyEqual(result.Input1M10.Value, math.Sqrt(5), tolerance) {
		checks.C1BackwardLessThan = "OK - input1 measurement10 follows the lessThan backward rule."
	} else {
		checks.C1BackwardLessThan = "FAIL - input1 measurement10 lessThan rule mismatch."
	}

	if result.Disturb2M10.Branch == "notLessThan" && nearlyEqual(result.Disturb2M10.Value, 45, tolerance) {
		checks.C2BackwardNotLessThan = "OK - disturbance2 measurement10 follows the notLessThan backward rule."
	} else {
		checks.C2BackwardNotLessThan = "FAIL - disturbance2 measurement10 notLessThan rule mismatch."
	}

	if in.Input2.Measurement2 {
		checks.C3FeedForwardGuard = "OK - input2 boolean guard is true for the feedforward rule."
	} else {
		checks.C3FeedForwardGuard = "FAIL - input2 boolean guard is false."
	}

	expectedFeedForward := result.FeedForward.InputMeasurement10*19.6 - math.Log10(result.FeedForward.Disturbance3)
	if nearlyEqual(result.FeedForward.Control, expectedFeedForward, tolerance) {
		checks.C4FeedForwardArithmetic = "OK - feedforward control equals product minus log10 compensation."
	} else {
		checks.C4FeedForwardArithmetic = "FAIL - feedforward arithmetic mismatch."
	}

	expectedError := in.Output2.Target2 - in.Output2.Measurement4
	if nearlyEqual(result.PNDFeedback.Error, expectedError, tolerance) {
		checks.C5FeedbackError = "OK - feedback error is target minus measurement."
	} else {
		checks.C5FeedbackError = "FAIL - feedback error mismatch."
	}

	expectedDifferential := in.State3.Observation3 - in.Output2.Measurement4
	if nearlyEqual(result.PNDFeedback.DifferentialError, expectedDifferential, tolerance) {
		checks.C6DifferentialError = "OK - differential error is state observation minus output measurement."
	} else {
		checks.C6DifferentialError = "FAIL - differential error mismatch."
	}

	expectedNonlinear := (7.3 / result.PNDFeedback.Error) * result.PNDFeedback.DifferentialError
	if nearlyEqual(result.PNDFeedback.NonlinearPart, expectedNonlinear, tolerance) {
		checks.C7NonlinearPart = "OK - nonlinear differential part equals (7.3 / error) times differential error."
	} else {
		checks.C7NonlinearPart = "FAIL - nonlinear differential part mismatch."
	}

	expectedFeedback := result.PNDFeedback.ProportionalPart + result.PNDFeedback.NonlinearPart
	if nearlyEqual(result.PNDFeedback.Control, expectedFeedback, tolerance) {
		checks.C8FeedbackControlSum = "OK - actuator2 control is proportional part plus nonlinear differential part."
	} else {
		checks.C8FeedbackControlSum = "FAIL - actuator2 control sum mismatch."
	}

	return checks
}

func allChecksPass(checks Checks) bool {
	return checks.C1BackwardLessThan[:2] == "OK" &&
		checks.C2BackwardNotLessThan[:2] == "OK" &&
		checks.C3FeedForwardGuard[:2] == "OK" &&
		checks.C4FeedForwardArithmetic[:2] == "OK" &&
		checks.C5FeedbackError[:2] == "OK" &&
		checks.C6DifferentialError[:2] == "OK" &&
		checks.C7NonlinearPart[:2] == "OK" &&
		checks.C8FeedbackControlSum[:2] == "OK"
}

func checkCount(checks Checks) int {
	count := 0
	if checks.C1BackwardLessThan[:2] == "OK" {
		count++
	}
	if checks.C2BackwardNotLessThan[:2] == "OK" {
		count++
	}
	if checks.C3FeedForwardGuard[:2] == "OK" {
		count++
	}
	if checks.C4FeedForwardArithmetic[:2] == "OK" {
		count++
	}
	if checks.C5FeedbackError[:2] == "OK" {
		count++
	}
	if checks.C6DifferentialError[:2] == "OK" {
		count++
	}
	if checks.C7NonlinearPart[:2] == "OK" {
		count++
	}
	if checks.C8FeedbackControlSum[:2] == "OK" {
		count++
	}
	return count
}

func nearlyEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// ---------- rendering ----------

func renderArcOutput(data Input, result InferenceResult) {
	ff := result.FeedForward
	fb := result.PNDFeedback

	fmt.Println("# Control System — ARC-style control-system example")
	fmt.Println()

	fmt.Println("## Answer")
	fmt.Printf("actuator1 control1 = %.6f\n", ff.Control)
	fmt.Printf("actuator2 control1 = %.6f\n", fb.Control)
	fmt.Printf("input1 measurement10 = %.6f (%s branch)\n", result.Input1M10.Value, result.Input1M10.Branch)
	fmt.Printf("disturbance2 measurement10 = %.6f (%s branch)\n", result.Disturb2M10.Value, result.Disturb2M10.Branch)
	fmt.Println()

	fmt.Println("## Reason")
	fmt.Println("The backward measurement10 rules first normalize measurement1 pairs into scalar measurement10 values.")
	fmt.Printf("For input1, %.0f < %.0f, so sqrt(%.0f - %.0f) = %.6f.\n",
		data.Input1.Measurement1[0], data.Input1.Measurement1[1],
		data.Input1.Measurement1[1], data.Input1.Measurement1[0], result.Input1M10.Value)
	fmt.Println("The first forward rule computes feedforward control from input1 measurement10, the input2 boolean guard, and disturbance1 compensation.")
	fmt.Println("The second forward rule computes PND feedback control from target/measurement error and the state/output differential error.")
	fmt.Println()

	return
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// ---------- main ----------

func main() {
	data := exampleinput.Load(seeExampleName, inputData())
	result, err := infer(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "control-system inference failed: %v\n", err)
		os.Exit(1)
	}
	renderArcOutput(data, result)
}
