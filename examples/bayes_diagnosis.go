// bayes_diagnosis.go
//
// A self-contained Go translation of examples/bayes-diagnosis.n3 from the Eyeling
// example suite, in ARC style.
//
// The original N3 program encodes a small Bayesian diagnostic model (four
// diseases, five symptoms), plugs in a patient-case evidence list, and computes
// posterior probabilities by multiplying prior × likelihood and normalising.
//
// This is intentionally not a generic N3 reasoner. The concrete N3 facts and
// rules are represented as ordinary Go data and functions so the probabilistic
// inference is easy to read and directly runnable.
//
// Run:
//
//	go run bayes_diagnosis.go
//
// The program has no third-party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "bayes_diagnosis"

// ---------- types ----------

// Disease holds the name and prior probability of one disease.
type Disease struct {
	Name  string
	Prior float64
}

// EvidenceItem is one observed symptom.
type EvidenceItem struct {
	Symptom string
	Present bool
}

// PosteriorResult stores the computed values for a single disease.
type PosteriorResult struct {
	Disease      string
	Unnormalized float64
	Posterior    float64
}

// Dataset contains the fixed facts translated from bayes-diagnosis.n3.
type Dataset struct {
	Diseases []Disease
	// probGiven[disease][symptom] = P(symptom | disease)
	ProbGiven map[string]map[string]float64
	Evidence  []EvidenceItem
}

// InferenceResult holds all derived facts needed by the renderer.
type InferenceResult struct {
	Scores  []float64
	Total   float64
	Results []PosteriorResult
}

// checks mirrors the proof obligations at the bottom of bayes-diagnosis.n3.
// In the N3 source, guards fire if any probability is outside [0,1].
type Checks struct {
	PriorsInRange    bool
	CondProbsInRange bool
}

// ---------- helpers ----------

// round3 rounds a float to 3 decimal places for display.
func round3(v float64) float64 {
	return math.Round(v*1000) / 1000
}

// round6 rounds a float to 6 decimal places for display.
func round6(v float64) float64 {
	return math.Round(v*1e6) / 1e6
}

// probInRange returns true if 0 <= p <= 1.
func probInRange(p float64) bool {
	return p >= 0 && p <= 1
}

// factor returns the likelihood factor for one disease and one evidence item.
//   - if evidence present → P(symptom | disease)
//   - if evidence absent  → 1 − P(symptom | disease)
func factor(d Disease, e EvidenceItem, probGiven map[string]map[string]float64) float64 {
	p := probGiven[d.Name][e.Symptom]
	if e.Present {
		return p
	}
	return 1 - p
}

// ---------- dataset ----------

func dataset() Dataset {
	return Dataset{
		Diseases: []Disease{
			{"COVID19", 0.05},
			{"Influenza", 0.03},
			{"AllergicRhinitis", 0.10},
			{"BacterialPneumonia", 0.01},
		},
		ProbGiven: map[string]map[string]float64{
			"COVID19": {
				"Fever": 0.70, "DryCough": 0.65, "LossOfSmell": 0.40,
				"Sneezing": 0.15, "ShortBreath": 0.20,
			},
			"Influenza": {
				"Fever": 0.80, "DryCough": 0.50, "LossOfSmell": 0.05,
				"Sneezing": 0.20, "ShortBreath": 0.10,
			},
			"AllergicRhinitis": {
				"Fever": 0.05, "DryCough": 0.15, "LossOfSmell": 0.10,
				"Sneezing": 0.80, "ShortBreath": 0.05,
			},
			"BacterialPneumonia": {
				"Fever": 0.70, "DryCough": 0.60, "LossOfSmell": 0.02,
				"Sneezing": 0.05, "ShortBreath": 0.60,
			},
		},
		Evidence: []EvidenceItem{
			{"Fever", true},
			{"DryCough", true},
			{"LossOfSmell", true},
			{"Sneezing", false},
			{"ShortBreath", true},
		},
	}
}

// ---------- guards ----------

func runGuards(data Dataset) {
	for _, d := range data.Diseases {
		if !probInRange(d.Prior) {
			fmt.Fprintf(os.Stderr, "guard failed: prior(%s) = %v\n", d.Name, d.Prior)
			os.Exit(1)
		}
	}
	for dname, syms := range data.ProbGiven {
		for sname, p := range syms {
			if !probInRange(p) {
				fmt.Fprintf(os.Stderr, "guard failed: pGiven(%s, %s) = %v\n", dname, sname, p)
				os.Exit(1)
			}
		}
	}
}

// ---------- inference ----------

func infer(data Dataset) InferenceResult {
	runGuards(data)

	// unnormalized score(d) = prior(d) * ∏ factor(d, e)
	scores := make([]float64, len(data.Diseases))
	for i, d := range data.Diseases {
		s := d.Prior
		for _, e := range data.Evidence {
			s *= factor(d, e, data.ProbGiven)
		}
		scores[i] = s
	}

	// evidence total = Σ scores
	total := 0.0
	for _, s := range scores {
		total += s
	}

	// posterior = score / total
	results := make([]PosteriorResult, len(data.Diseases))
	for i, d := range data.Diseases {
		post := 0.0
		if total > 0 {
			post = scores[i] / total
		}
		results[i] = PosteriorResult{
			Disease:      d.Name,
			Unnormalized: scores[i],
			Posterior:    post,
		}
	}
	return InferenceResult{Scores: scores, Total: total, Results: results}
}

// ---------- checks ----------

func performChecks(data Dataset) Checks {
	priorsOK := true
	for _, d := range data.Diseases {
		if !probInRange(d.Prior) {
			priorsOK = false
			break
		}
	}
	condOK := true
	for _, syms := range data.ProbGiven {
		for _, p := range syms {
			if !probInRange(p) {
				condOK = false
				break
			}
		}
	}
	return Checks{PriorsInRange: priorsOK, CondProbsInRange: condOK}
}

func allChecksPass(c Checks) bool {
	return c.PriorsInRange && c.CondProbsInRange
}

func checkCount(c Checks) int {
	count := 0
	if c.PriorsInRange {
		count++
	}
	if c.CondProbsInRange {
		count++
	}
	return count
}

// ---------- rendering ----------

// formatEvidence returns a human-readable string of the evidence list.
func formatEvidence(ev []EvidenceItem) string {
	var parts []string
	for _, e := range ev {
		val := "present"
		if !e.Present {
			val = "absent"
		}
		parts = append(parts, fmt.Sprintf("%s=%s", e.Symptom, val))
	}
	return strings.Join(parts, ", ")
}

// renderArcOutput prints the answer / reason / check / audit style.
func renderArcOutput(data Dataset, result InferenceResult, checks Checks) {
	fmt.Println("# Bayes Diagnosis")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	var best PosteriorResult
	for _, r := range result.Results {
		if r.Posterior > best.Posterior {
			best = r
		}
	}
	fmt.Printf("The most likely disease is %s (posterior = %.6f).\n", best.Disease, best.Posterior)
	fmt.Println()
	fmt.Println("- Full posterior distribution:")
	for _, r := range result.Results {
		fmt.Printf("  %-20s  posterior = %.6f  (unnormalized = %.8f)\n",
			r.Disease, r.Posterior, r.Unnormalized)
	}
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Printf("- Evidence: %s.\n", formatEvidence(data.Evidence))
	fmt.Printf("Evidence total (normalizing constant) = %.8f.\n", result.Total)
	fmt.Println("- The posterior for each disease is computed as:")
	fmt.Println("  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal")
	fmt.Println("where for an absent symptom the factor is 1 − P(symptom|d).")
	fmt.Println()

	// --- Check ---
	fmt.Println("## Check")
	if checks.PriorsInRange {
		fmt.Println("- C1 OK - all prior probabilities are in [0,1].")
	} else {
		fmt.Println("- C1 FAIL - one or more prior probabilities are outside [0,1].")
	}
	if checks.CondProbsInRange {
		fmt.Println("- C2 OK - all conditional probabilities are in [0,1].")
	} else {
		fmt.Println("- C2 FAIL - one or more conditional probabilities are outside [0,1].")
	}
	fmt.Println()

	// --- Go audit details ---
	fmt.Println("## Go audit details")
	fmt.Printf("- platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("- diseases : %d\n", len(data.Diseases))
	fmt.Printf("- symptoms : %d\n", len(data.ProbGiven[data.Diseases[0].Name]))
	fmt.Printf("- evidence items : %d\n", len(data.Evidence))
	fmt.Printf("- evidence total : %.8f\n", result.Total)
	fmt.Println("- posteriors :")
	for _, r := range result.Results {
		fmt.Printf("  %-20s  unnormalized=%.8f  posterior=%.6f\n",
			r.Disease, r.Unnormalized, r.Posterior)
	}
	fmt.Printf("- checks passed : %d/2\n", checkCount(checks))
	fmt.Printf("- recommendation consistent : %s\n", yesNo(allChecksPass(checks)))
}

func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

// ---------- main ----------

func main() {
	data := exampleinput.Load(eyelingoExampleName, dataset())

	// Run inference (guards are checked inside infer).
	result := infer(data)

	// Perform explicit consistency checks.
	checks := performChecks(data)

	// Render ARC-style output.
	renderArcOutput(data, result, checks)
}
