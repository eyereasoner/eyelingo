// gradient_descent_step.go
//
// Inspired by Eyeling's `examples/gd-step-certified.n3`.
//
// The example certifies one bounded gradient-descent step on a convex
// quadratic instead of running an open-ended optimizer.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"os"
	"runtime"
)

const eyelingoExampleName = "gradient_descent_step"

type Dataset struct {
	CaseName string `json:"caseName"`
	Question string `json:"question"`
	Function struct {
		CenterX float64 `json:"centerX"`
		CenterY float64 `json:"centerY"`
		WeightY float64 `json:"weightY"`
	} `json:"function"`
	Start struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"start"`
	StepSize    float64 `json:"stepSize"`
	MaxStepNorm float64 `json:"maxStepNorm"`
	Expected    struct {
		Decreases bool `json:"decreases"`
	} `json:"expected"`
}
type Check struct {
	ID   string
	OK   bool
	Text string
}
type Analysis struct {
	GX, GY, NX, NY, Before, After, StepNorm float64
	Checks                                  []Check
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
	x, y := ds.Start.X, ds.Start.Y
	gx := 2 * (x - ds.Function.CenterX)
	gy := 2 * ds.Function.WeightY * (y - ds.Function.CenterY)
	nx := x - ds.StepSize*gx
	ny := y - ds.StepSize*gy
	before := objective(ds, x, y)
	after := objective(ds, nx, ny)
	stepNorm := math.Hypot(nx-x, ny-y)
	checks := []Check{
		{"C1", ds.Function.CenterX == 3 && ds.Function.CenterY == -1 && ds.Function.WeightY == 2, "gradient was derived from f(x,y) = (x-3)^2 + 2(y+1)^2"},
		{"C2", ds.StepSize > 0 && ds.StepSize < 0.25, "step size is positive and below the conservative bound"},
		{"C3", close(nx, 7.0) && close(ny, 2.6), "new point is (7.000, 2.600)"},
		{"C4", after < before && close(before, 97.0) && close(after, 41.92), "objective value decreases from 97.000 to 41.920"},
		{"C5", stepNorm < ds.MaxStepNorm, "step norm stays below 3.000"},
		{"C6", ds.Expected.Decreases && after < before, "the JSON expected decrease flag is satisfied"},
	}
	return Analysis{gx, gy, nx, ny, before, after, stepNorm, checks}
}
func objective(ds Dataset, x, y float64) float64 {
	dx := x - ds.Function.CenterX
	dy := y - ds.Function.CenterY
	return dx*dx + ds.Function.WeightY*dy*dy
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
	fmt.Println("# Gradient Descent Step")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("start point : (%.3f, %.3f)\n", ds.Start.X, ds.Start.Y)
	fmt.Printf("gradient : (%.3f, %.3f)\n", a.GX, a.GY)
	fmt.Printf("step size : %.3f\n", ds.StepSize)
	fmt.Printf("next point : (%.3f, %.3f)\n", a.NX, a.NY)
	fmt.Printf("objective before : %.3f\n", a.Before)
	fmt.Printf("objective after : %.3f\n", a.After)
	fmt.Printf("decrease : %.3f\n", a.Before-a.After)
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("The quadratic is convex and its gradient is computed symbolically from the JSON coefficients.")
	fmt.Println("The update uses x_next = x - alpha × gradient, with alpha fixed by the input.")
	fmt.Println("The new point stays within the declared step-norm bound and produces a strictly smaller objective value.")
	fmt.Println("That gives a small finite certificate for one safe descent step rather than an open-ended optimizer.")
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
	fmt.Printf("center : (%.3f, %.3f)\n", ds.Function.CenterX, ds.Function.CenterY)
	fmt.Printf("weightY : %.3f\n", ds.Function.WeightY)
	fmt.Printf("step norm : %.3f\n", a.StepNorm)
	fmt.Printf("max step norm : %.3f\n", ds.MaxStepNorm)
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
