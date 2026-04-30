// complex_matrix_stability.go
//
// Inspired by Eyeling's `examples/complex-matrix-stability.n3`.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"os"
	"runtime"
)

const eyelingoExampleName = "complex_matrix_stability"

type Complex struct {
	Re float64 `json:"re"`
	Im float64 `json:"im"`
}

type Matrix struct {
	Name     string    `json:"name"`
	Diagonal []Complex `json:"diagonal"`
}

type SampleProduct struct {
	Z Complex `json:"z"`
	W Complex `json:"w"`
}

type Dataset struct {
	CaseName      string        `json:"caseName"`
	Source        string        `json:"source"`
	Question      string        `json:"question"`
	Scale         float64       `json:"scale"`
	Matrices      []Matrix      `json:"matrices"`
	SampleProduct SampleProduct `json:"sampleProduct"`
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Result struct {
	Radii          map[string]float64
	Classes        map[string]string
	Checks         []Check
	ScaledRadiusSq float64
	ProductOK      bool
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	r := derive(ds)
	printReport(ds, r)
	if !allOK(r.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Result {
	r := Result{Radii: map[string]float64{}, Classes: map[string]string{}}
	unstableSq := 0.0

	for _, m := range ds.Matrices {
		maxAbs2 := 0.0
		for _, z := range m.Diagonal {
			if a := abs2(z); a > maxAbs2 {
				maxAbs2 = a
			}
		}
		radius := math.Sqrt(maxAbs2)
		r.Radii[m.Name] = radius
		r.Classes[m.Name] = classify(radius)
		if m.Name == "A_unstable" {
			unstableSq = maxAbs2
		}
	}

	z := ds.SampleProduct.Z
	w := ds.SampleProduct.W
	zw := Complex{Re: z.Re*w.Re - z.Im*w.Im, Im: z.Re*w.Im + z.Im*w.Re}
	r.ProductOK = nearly(abs2(zw), abs2(z)*abs2(w))
	r.ScaledRadiusSq = ds.Scale * ds.Scale * unstableSq

	r.Checks = []Check{
		{"C1", r.Classes["A_unstable"] == "unstable" && nearly(r.Radii["A_unstable"], 2), "A_unstable has spectral radius 2, so it is unstable"},
		{"C2", r.Classes["A_stable"] == "marginally stable" && nearly(r.Radii["A_stable"], 1), "A_stable has spectral radius 1, so it is marginally stable"},
		{"C3", r.Classes["A_damped"] == "damped" && nearly(r.Radii["A_damped"], 0), "A_damped has spectral radius 0, so every mode decays"},
		{"C4", r.ProductOK, "squared modulus of z*w equals the product of squared moduli"},
		{"C5", nearly(r.ScaledRadiusSq, 4*unstableSq), "scaling A_unstable by 2 multiplies spectral-radius-squared by 4"},
	}
	return r
}

func printReport(ds Dataset, r Result) {
	fmt.Println("# Complex Matrix Stability")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("A_unstable : spectral radius %.0f -> %s\n", r.Radii["A_unstable"], r.Classes["A_unstable"])
	fmt.Printf("A_stable : spectral radius %.0f -> %s\n", r.Radii["A_stable"], r.Classes["A_stable"])
	fmt.Printf("A_damped : spectral radius %.0f -> %s\n", r.Radii["A_damped"], r.Classes["A_damped"])
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("For a discrete-time linear system x_{k+1} = A x_k, the eigenvalues of A govern the modal behaviour.")
	fmt.Println("Because the matrices are diagonal, the eigenvalues are the diagonal entries; the largest modulus gives the spectral radius and therefore the stability class.")
	fmt.Println()
	fmt.Println("## Check")
	for _, c := range r.Checks {
		status := "FAIL"
		if c.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", c.ID, status, c.Text)
	}
	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("translated source : %s\n", ds.Source)
	fmt.Printf("matrices checked : %d\n", len(ds.Matrices))
	fmt.Printf("scale factor : %.0f\n", ds.Scale)
	fmt.Printf("scaled unstable radius squared : %.0f\n", r.ScaledRadiusSq)
	fmt.Printf("checks passed : %d/%d\n", countOK(r.Checks), len(r.Checks))
}

func abs2(z Complex) float64 { return z.Re*z.Re + z.Im*z.Im }

func classify(radius float64) string {
	if radius > 1 {
		return "unstable"
	}
	if nearly(radius, 1) {
		return "marginally stable"
	}
	return "damped"
}

func nearly(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func allOK(cs []Check) bool {
	for _, c := range cs {
		if !c.OK {
			return false
		}
	}
	return true
}

func countOK(cs []Check) int {
	n := 0
	for _, c := range cs {
		if c.OK {
			n++
		}
	}
	return n
}
