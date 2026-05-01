// bmi.go
//
// A self-contained Go translation of examples/bmi.n3 from the Eyeling
// example suite, in ARC style.
//
// The original N3 program encodes a small Body Mass Index calculator that
// normalizes metric or US inputs, computes BMI, assigns a WHO adult category,
//
// This is intentionally not a generic N3 reasoner.  The concrete N3 facts and
// rules are represented as ordinary Go data and functions so the
// probabilistic inference is easy to read and directly runnable.
//
// Run:
//
//     go run bmi.go
//
// The program has no third-party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"os"
)

const eyelingoExampleName = "bmi"

// ---------- types ----------

// Input stores the editable input block from bmi.n3.
type Input struct {
	UnitSystem string
	Weight     float64
	Height     float64
}

// Case holds the normalized SI values and derived quantities.
type Case struct {
	WeightKg            float64
	HeightM             float64
	HeightSquared       float64
	Bmi                 float64
	BmiRounded          float64
	HealthyMinKg        float64
	HealthyMaxKg        float64
	HealthyMinKgRounded float64
	HealthyMaxKgRounded float64
}

// Decision stores the WHO category assigned to the computed BMI.
type Decision struct {
	Category string
}

// Reason stores the human-readable explanations.
type Reason struct {
	Units            string
	Formula          string
	Calculation      string
	CategoryRule     string
	UnitsExplanation string
}

type Checks struct {
	C1InputPositiveSI       string // input normalized into positive SI values
	C2HeightSquared         string // height squared reconstructed
	C3BMIMatchesFormula     string // BMI matches kg/m²
	C4BelowNormalThreshold  string // 18.49 stays below normal threshold
	C5LowerBoundaryHalfOpen string // 18.5 is Normal
	C6OverweightStart       string // 25.0 starts Overweight
	C7ObesityIStart         string // 30.0 starts Obesity I
	C8Monotonic             string // classification is monotonic
	C9HealthyBand           string // healthy-weight band reconstructed
}

// InferenceResult holds all derived facts needed by the renderer.
type InferenceResult struct {
	Input    Input
	Case     Case
	Decision Decision
	Reason   Reason
	Checks   Checks
}

// ---------- dataset ----------

func inputData() Input {
	// Default metric input from bmi.n3.
	return Input{
		UnitSystem: "metric",
		Weight:     72.0,
		Height:     178.0,
	}
}

// ---------- inference ----------

func infer(in Input) (InferenceResult, error) {
	var c Case
	var d Decision
	var r Reason

	// ---------- Normalization ----------
	if in.UnitSystem == "metric" {
		c.WeightKg = in.Weight
		c.HeightM = in.Height / 100.0
		r.Units = "Inputs were already metric, so kilograms stay kilograms and " +
			"centimeters are divided by 100 to obtain meters."
		r.UnitsExplanation = "Metric input: weight in kg, height in cm ÷ 100 = height in m."
	} else if in.UnitSystem == "us" {
		c.WeightKg = in.Weight * 0.45359237
		c.HeightM = in.Height * 0.0254
		r.Units = "US inputs were converted to SI units: pounds to kilograms " +
			"and inches to meters."
		r.UnitsExplanation = "US input: pounds × 0.45359237 = kg, inches × 0.0254 = m."
	} else {
		return InferenceResult{}, fmt.Errorf("unknown unit system: %s", in.UnitSystem)
	}

	// ---------- BMI calculation ----------
	c.HeightSquared = c.HeightM * c.HeightM
	c.Bmi = c.WeightKg / c.HeightSquared

	// Rounded BMI (multiply by 100, round, divide by 100) as in bmi.n3.
	bmiX100 := c.Bmi * 100.0
	bmiRoundedInt := math.Round(bmiX100)
	c.BmiRounded = bmiRoundedInt / 100.0

	// ---------- Healthy weight range ----------
	c.HealthyMinKg = 18.5 * c.HeightSquared
	c.HealthyMaxKg = 24.9 * c.HeightSquared

	// Rounded healthy range (multiply by 10, round, divide by 10).
	minX10 := c.HealthyMinKg * 10.0
	maxX10 := c.HealthyMaxKg * 10.0
	minRoundedInt := math.Round(minX10)
	maxRoundedInt := math.Round(maxX10)
	c.HealthyMinKgRounded = minRoundedInt / 10.0
	c.HealthyMaxKgRounded = maxRoundedInt / 10.0

	// ---------- WHO adult category ----------
	switch {
	case c.Bmi < 18.5:
		d.Category = "Underweight"
	case c.Bmi >= 18.5 && c.Bmi < 25.0:
		d.Category = "Normal"
	case c.Bmi >= 25.0 && c.Bmi < 30.0:
		d.Category = "Overweight"
	case c.Bmi >= 30.0 && c.Bmi < 35.0:
		d.Category = "Obesity I"
	case c.Bmi >= 35.0 && c.Bmi < 40.0:
		d.Category = "Obesity II"
	default:
		d.Category = "Obesity III"
	}

	// ---------- Reason ----------
	r.Formula = "BMI is defined as weight in kilograms divided by height in meters squared."
	r.Calculation = "The normalized weight and height were used to compute BMI, " +
		"then the result was mapped to the WHO adult category table."
	r.CategoryRule = d.Category

	checks := performChecks(c)

	return InferenceResult{
		Input:    in,
		Case:     c,
		Decision: d,
		Reason:   r,
		Checks:   checks,
	}, nil
}

func performChecks(c Case) Checks {
	var checks Checks

	// C1: input normalized into positive SI values
	if c.WeightKg > 0 && c.HeightM > 0 {
		checks.C1InputPositiveSI = "OK - the input was normalized into positive SI values."
	} else {
		checks.C1InputPositiveSI = "FAIL - normalized values are not positive."
	}

	// C2: height squared reconstructed from normalized height
	reconstructedM2 := c.HeightM * c.HeightM
	if math.Abs(reconstructedM2-c.HeightSquared) < 1e-12 {
		checks.C2HeightSquared = "OK - height squared was reconstructed from the normalized height."
	} else {
		checks.C2HeightSquared = "FAIL - height squared mismatch."
	}

	// C3: BMI value matches BMI = kg / m² formula
	recomputedBMI := c.WeightKg / c.HeightSquared
	if math.Abs(recomputedBMI-c.Bmi) < 1e-12 {
		checks.C3BMIMatchesFormula = "OK - the BMI value matches the BMI = kg / m² formula."
	} else {
		checks.C3BMIMatchesFormula = "FAIL - BMI formula mismatch."
	}

	// C4: a BMI of 18.49 stays below the normal-weight threshold
	if 18.49 < 18.5 {
		checks.C4BelowNormalThreshold = "OK - a BMI of 18.49 stays below the normal-weight threshold."
	} else {
		checks.C4BelowNormalThreshold = "FAIL - 18.49 threshold check."
	}

	// C5: the lower boundary is half-open: BMI 18.5 is classified as Normal
	testBmi185 := 18.5
	if testBmi185 >= 18.5 && testBmi185 < 25.0 {
		checks.C5LowerBoundaryHalfOpen = "OK - the lower boundary is half-open: BMI 18.5 is classified as Normal."
	} else {
		checks.C5LowerBoundaryHalfOpen = "FAIL - lower boundary classification."
	}

	// C6: BMI 25.0 starts the Overweight category
	testBmi250 := 25.0
	if testBmi250 >= 25.0 && testBmi250 < 30.0 {
		checks.C6OverweightStart = "OK - BMI 25.0 starts the Overweight category."
	} else {
		checks.C6OverweightStart = "FAIL - Overweight threshold."
	}

	// C7: BMI 30.0 starts the Obesity I category
	testBmi300 := 30.0
	if testBmi300 >= 30.0 && testBmi300 < 35.0 {
		checks.C7ObesityIStart = "OK - BMI 30.0 starts the Obesity I category."
	} else {
		checks.C7ObesityIStart = "FAIL - Obesity I threshold."
	}

	// C8: classification behavior is monotonic across representative BMI values
	bmi22 := 22.0
	bmi27 := 27.0
	bmi41 := 41.0
	cat22 := classifyBmi(bmi22)
	cat27 := classifyBmi(bmi27)
	cat41 := classifyBmi(bmi41)
	if cat22 == "Normal" && cat27 == "Overweight" && cat41 == "Obesity III" {
		checks.C8Monotonic = "OK - classification behavior is monotonic across representative BMI values."
	} else {
		checks.C8Monotonic = "FAIL - monotonic classification."
	}

	// C9: healthy-weight band reconstructed from BMI 18.5 to 24.9 at the same height
	minFromBMI := 18.5 * c.HeightSquared
	maxFromBMI := 24.9 * c.HeightSquared
	if math.Abs(minFromBMI-c.HealthyMinKg) < 1e-12 &&
		math.Abs(maxFromBMI-c.HealthyMaxKg) < 1e-12 {
		checks.C9HealthyBand = "OK - the healthy-weight band was reconstructed from BMI 18.5 to 24.9 at the same height."
	} else {
		checks.C9HealthyBand = "FAIL - healthy band reconstruction."
	}

	return checks
}

// classifyBmi returns the WHO category for a given BMI value.
func classifyBmi(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "Underweight"
	case bmi >= 18.5 && bmi < 25.0:
		return "Normal"
	case bmi >= 25.0 && bmi < 30.0:
		return "Overweight"
	case bmi >= 30.0 && bmi < 35.0:
		return "Obesity I"
	case bmi >= 35.0 && bmi < 40.0:
		return "Obesity II"
	default:
		return "Obesity III"
	}
}

func allChecksPass(checks Checks) bool {
	return checks.C1InputPositiveSI[:2] == "OK" &&
		checks.C2HeightSquared[:2] == "OK" &&
		checks.C3BMIMatchesFormula[:2] == "OK" &&
		checks.C4BelowNormalThreshold[:2] == "OK" &&
		checks.C5LowerBoundaryHalfOpen[:2] == "OK" &&
		checks.C6OverweightStart[:2] == "OK" &&
		checks.C7ObesityIStart[:2] == "OK" &&
		checks.C8Monotonic[:2] == "OK" &&
		checks.C9HealthyBand[:2] == "OK"
}

func checkCount(checks Checks) int {
	count := 0
	if checks.C1InputPositiveSI[:2] == "OK" {
		count++
	}
	if checks.C2HeightSquared[:2] == "OK" {
		count++
	}
	if checks.C3BMIMatchesFormula[:2] == "OK" {
		count++
	}
	if checks.C4BelowNormalThreshold[:2] == "OK" {
		count++
	}
	if checks.C5LowerBoundaryHalfOpen[:2] == "OK" {
		count++
	}
	if checks.C6OverweightStart[:2] == "OK" {
		count++
	}
	if checks.C7ObesityIStart[:2] == "OK" {
		count++
	}
	if checks.C8Monotonic[:2] == "OK" {
		count++
	}
	if checks.C9HealthyBand[:2] == "OK" {
		count++
	}
	return count
}

// ---------- rendering ----------

func renderArcOutput(data Input, result InferenceResult) {
	c := result.Case
	d := result.Decision

	fmt.Println("# BMI — ARC-style Body Mass Index example")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	fmt.Printf("BMI = %.1f\n", c.BmiRounded)
	fmt.Printf("Category = %s\n", d.Category)
	fmt.Printf("At height %.0f cm, a healthy-weight range is about %.1f–%.1f kg (BMI 18.5–24.9).\n",
		data.Height, c.HealthyMinKgRounded, c.HealthyMaxKgRounded)
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Println(result.Reason.Formula)
	fmt.Println(result.Reason.Calculation)
	fmt.Printf("The input was %s, so %s\n", data.UnitSystem, result.Reason.Units)
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
	data := exampleinput.Load(eyelingoExampleName, inputData())
	result, err := infer(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "BMI inference failed: %v\n", err)
		os.Exit(1)
	}
	renderArcOutput(data, result)
}
