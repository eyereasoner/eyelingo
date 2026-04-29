// lldm.go
//
// A self-contained Go translation of examples/lldm.n3 from the Eyeling
// example suite, in ARC style.
//
// The original N3 program encodes a geometric model for Leg Length
// Discrepancy Measurement.  Four measurement points on a medical image are
// processed through a pipeline of coordinate differences, line slopes,
// intersection points, and Euclidean distances to decide whether an alarm
// should be raised.
//
// This is intentionally not a generic N3 reasoner.  The concrete N3 facts
// and derivation rules are represented as ordinary Go data and functions so
// the decision logic is easy to read and directly runnable.
//
// Run:
//
//     go run lldm.go
//
// The program has no third-party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"runtime"
)

const eyelingoExampleName = "lldm"

// ---------- types ----------

// Point holds the (x, y) coordinates of one measurement point.
type Point struct {
	Name string
	X    float64
	Y    float64
}

// Measurement collects the eight input coordinates that appear in
// the original N3 data block.
type Measurement struct {
	P1x, P1y float64
	P2x, P2y float64
	P3x, P3y float64
	P4x, P4y float64
}

// Derived holds every intermediate and final quantity that the N3 rules
// compute.  The field names mirror the N3 predicate names.
type Derived struct {
	Dx12, Dy12             float64 // coordinate differences
	Dx51, Dx53, Dx62, Dx64 float64
	Dy13, Dy24, Dy53, Dy64 float64
	CL1                    float64 // slope of L1 = dy12 / dx12
	DL3m                   float64 // 1 / cL1
	CL3                    float64 // -dL3m  (slope of L3 and L4, perpendicular to L1)
	PL1x1, PL1x2           float64 // cL1 * p1x , cL1 * p2x
	PL3x3, PL3x4           float64 // cL3 * p3x , cL3 * p4x
	Dd13, Ddy13            float64
	Dd24, Ddy24            float64
	DdL13                  float64 // cL1 - cL3
	P5x, P5y               float64 // intersection of L1 and L3
	P6x, P6y               float64 // intersection of L1 and L4
	PL1dx51, PL1dx62       float64
	Sdx53, Sdx64           float64 // squared
	Sdy53, Sdy64           float64
	Ssd53, Ssd64           float64 // sums of squares
	D53, D64               float64 // Euclidean distances
	D                      float64 // leg-length discrepancy
}

// InferenceResult holds all derived facts needed by the renderer.
type InferenceResult struct {
	Input Measurement
	D     Derived
	Alarm bool
}

// Checks mirrors the proof obligations from the N3 axioms.
type Checks struct {
	PointsOnLines        bool
	Perpendicular        bool
	P5Intersection       bool
	P6Intersection       bool
	DistancesNonNegative bool
	DcmInRange           bool
}

// ---------- dataset ----------

func inputMeasurement() Measurement {
	return Measurement{
		P1x: 10.1, P1y: 7.8,
		P2x: 45.1, P2y: 5.6,
		P3x: 3.6, P3y: 29.8,
		P4x: 54.7, P4y: 28.5,
	}
}

// ---------- inference ----------

func infer(m Measurement) InferenceResult {
	var d Derived

	// === Block 1: differences that depend only on input coordinates ===
	// N3: (?p1x ?p2x) math:difference → dx12  (dx12 = p1x - p2x)
	d.Dx12 = m.P1x - m.P2x
	// N3: (?p1y ?p2y) math:difference → dy12
	d.Dy12 = m.P1y - m.P2y
	// N3: (?p1y ?p3y) math:difference → dy13
	d.Dy13 = m.P1y - m.P3y
	// N3: (?p2y ?p4y) math:difference → dy24
	d.Dy24 = m.P2y - m.P4y

	// === Block 2: slopes of the lines ===
	// cL1 = dy12 / dx12
	d.CL1 = d.Dy12 / d.Dx12
	// dL3m = 1 / cL1
	d.DL3m = 1.0 / d.CL1
	// cL3 = 0 - dL3m  (i.e. cL3 = -1/cL1)
	d.CL3 = 0.0 - d.DL3m

	// === Block 3: projected x-values ===
	d.PL1x1 = d.CL1 * m.P1x
	d.PL1x2 = d.CL1 * m.P2x
	d.PL3x3 = d.CL3 * m.P3x
	d.PL3x4 = d.CL3 * m.P4x

	// === Block 4: intermediate differences ===
	d.Dd13 = d.PL1x1 - d.PL3x3
	d.Ddy13 = d.Dd13 - d.Dy13
	d.Dd24 = d.PL1x2 - d.PL3x4
	d.Ddy24 = d.Dd24 - d.Dy24
	d.DdL13 = d.CL1 - d.CL3

	// === Block 5: intersection points p5 and p6 ===
	d.P5x = d.Ddy13 / d.DdL13
	d.P6x = d.Ddy24 / d.DdL13

	// differences that depend on p5x and p6x
	d.Dx51 = d.P5x - m.P1x
	d.Dx53 = d.P5x - m.P3x
	d.Dx62 = d.P6x - m.P2x
	d.Dx64 = d.P6x - m.P4x

	d.PL1dx51 = d.CL1 * d.Dx51
	d.PL1dx62 = d.CL1 * d.Dx62

	d.P5y = d.PL1dx51 + m.P1y
	d.P6y = d.PL1dx62 + m.P2y

	// === Block 6: remaining coordinate differences ===
	d.Dy53 = d.P5y - m.P3y
	d.Dy64 = d.P6y - m.P4y

	// === Block 7: squares of differences ===
	d.Sdx53 = d.Dx53 * d.Dx53
	d.Sdx64 = d.Dx64 * d.Dx64
	d.Sdy53 = d.Dy53 * d.Dy53
	d.Sdy64 = d.Dy64 * d.Dy64

	// === Block 8: sums of squares → squared distances ===
	d.Ssd53 = d.Sdx53 + d.Sdy53
	d.Ssd64 = d.Sdx64 + d.Sdy64

	// === Block 9: Euclidean distances ===
	d.D53 = math.Sqrt(d.Ssd53)
	d.D64 = math.Sqrt(d.Ssd64)

	// === Block 10: leg-length discrepancy and alarm decision ===
	d.D = d.D53 - d.D64
	alarm := d.D < -1.25 || d.D > 1.25

	return InferenceResult{
		Input: m,
		D:     d,
		Alarm: alarm,
	}
}

// ---------- checks ----------

func performChecks(m Measurement, d Derived) Checks {
	const eps = 1e-9

	// Check that p5 is incident on L1 and L3.
	yL1atP5 := d.CL1*d.P5x + (m.P1y - d.CL1*m.P1x)
	yL3atP5 := d.CL3*d.P5x + (m.P3y - d.CL3*m.P3x)
	p5onLines := math.Abs(d.P5y-yL1atP5) < eps && math.Abs(d.P5y-yL3atP5) < eps

	// Check that p6 is incident on L1 and L4.
	yL1atP6 := d.CL1*d.P6x + (m.P1y - d.CL1*m.P1x)
	yL4atP6 := d.CL3*d.P6x + (m.P4y - d.CL3*m.P4x)
	p6onLines := math.Abs(d.P6y-yL1atP6) < eps && math.Abs(d.P6y-yL4atP6) < eps

	// Check perpendicular: product of slopes should be -1.
	perp := math.Abs(d.CL1*d.CL3+1.0) < eps

	// Check that squared sums are non‑negative (always true for squares).
	distsNonNeg := d.Ssd53 >= 0 && d.Ssd64 >= 0

	// Check that dCm is within reasonable range (not NaN or Inf).
	dcmOk := !math.IsNaN(d.D) && !math.IsInf(d.D, 0)

	return Checks{
		PointsOnLines:        p5onLines && p6onLines,
		Perpendicular:        perp,
		P5Intersection:       p5onLines,
		P6Intersection:       p6onLines,
		DistancesNonNegative: distsNonNeg,
		DcmInRange:           dcmOk,
	}
}

func allChecksPass(c Checks) bool {
	return c.PointsOnLines && c.Perpendicular &&
		c.P5Intersection && c.P6Intersection &&
		c.DistancesNonNegative && c.DcmInRange
}

func checkCount(c Checks) int {
	n := 0
	if c.Perpendicular {
		n++
	}
	if c.P5Intersection {
		n++
	}
	if c.P6Intersection {
		n++
	}
	if c.DistancesNonNegative {
		n++
	}
	if c.DcmInRange {
		n++
	}
	return n
}

// ---------- rendering ----------

func renderArcOutput(res InferenceResult, checks Checks) {
	fmt.Println("# LLD — Leg Length Discrepancy Measurement")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	if res.Alarm {
		fmt.Printf("LLD Alarm = TRUE  (discrepancy dCm = %.6f, threshold ±1.25)\n", res.D.D)
	} else {
		fmt.Printf("LLD Alarm = FALSE (discrepancy dCm = %.6f, threshold ±1.25)\n", res.D.D)
	}
	fmt.Println()
	fmt.Println("Key computed values:")
	fmt.Printf("  SL1 = %.6f  SL3 = SL4 = %.6f\n", res.D.CL1, res.D.CL3)
	fmt.Printf("  p5  = (%.4f, %.4f)\n", res.D.P5x, res.D.P5y)
	fmt.Printf("  p6  = (%.4f, %.4f)\n", res.D.P6x, res.D.P6y)
	fmt.Printf("  d53 = %.6f cm\n", res.D.D53)
	fmt.Printf("  d64 = %.6f cm\n", res.D.D64)
	fmt.Printf("  dCm = %.6f cm\n", res.D.D)
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Println("The measurement points p1-p4 define two parallel lines L1 and L3/L4.")
	fmt.Println("L1 is defined by p1 and p2; L3 passes through p3; L4 passes through p4.")
	fmt.Println("L3 and L4 are perpendicular to L1.  Intersection points p5 (L1∩L3)")
	fmt.Println("and p6 (L1∩L4) are computed analytically.  The Euclidean distances")
	fmt.Println("d53 (p5–p3) and d64 (p6–p4) are then computed, and their difference")
	fmt.Println("dCm = d53 − d64 is the leg-length discrepancy.  An alarm fires when")
	fmt.Println("|dCm| > 1.25 cm.")
	fmt.Println()

	// --- Check ---
	fmt.Println("## Check")
	if checks.Perpendicular {
		fmt.Println("C1 OK - L1 is perpendicular to L3 and L4 (slopes product ≈ -1).")
	} else {
		fmt.Println("C1 FAIL - L1 not perpendicular to L3/L4.")
	}
	if checks.P5Intersection {
		fmt.Println("C2 OK - p5 lies on both L1 and L3.")
	} else {
		fmt.Println("C2 FAIL - p5 not incident on L1∩L3.")
	}
	if checks.P6Intersection {
		fmt.Println("C3 OK - p6 lies on both L1 and L4.")
	} else {
		fmt.Println("C3 FAIL - p6 not incident on L1∩L4.")
	}
	if checks.DistancesNonNegative {
		fmt.Println("C4 OK - squared distances are non‑negative.")
	} else {
		fmt.Println("C4 FAIL - negative squared distance.")
	}
	if checks.DcmInRange {
		fmt.Println("C5 OK - dCm is a finite number.")
	} else {
		fmt.Println("C5 FAIL - dCm is NaN or Inf.")
	}
	fmt.Println()

	// --- Go audit details ---
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("input points : p1(%.1f, %.1f) p2(%.1f, %.1f) p3(%.1f, %.1f) p4(%.1f, %.1f)\n",
		res.Input.P1x, res.Input.P1y,
		res.Input.P2x, res.Input.P2y,
		res.Input.P3x, res.Input.P3y,
		res.Input.P4x, res.Input.P4y,
	)
	fmt.Printf("cL1 (slope L1)          : %.8f\n", res.D.CL1)
	fmt.Printf("cL3 (slope L3, L4)      : %.8f\n", res.D.CL3)
	fmt.Printf("p5 (x, y)               : %.8f, %.8f\n", res.D.P5x, res.D.P5y)
	fmt.Printf("p6 (x, y)               : %.8f, %.8f\n", res.D.P6x, res.D.P6y)
	fmt.Printf("d53                     : %.8f cm\n", res.D.D53)
	fmt.Printf("d64                     : %.8f cm\n", res.D.D64)
	fmt.Printf("dCm (LL discrepancy)    : %.8f cm\n", res.D.D)
	fmt.Printf("alarm threshold         : ±1.25 cm\n")
	fmt.Printf("LLD alarm               : %v\n", res.Alarm)
	fmt.Printf("checks passed           : %d/5\n", checkCount(checks))
	fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(checks)))
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// ---------- main ----------

func main() {
	input := exampleinput.Load(eyelingoExampleName, inputMeasurement())
	result := infer(input)
	checks := performChecks(result.Input, result.D)
	renderArcOutput(result, checks)
}
