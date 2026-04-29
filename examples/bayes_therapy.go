// bayes_therapy.go
//
// A self-contained Go translation of examples/bayes-therapy.n3 from the Eyeling
// example suite, in ARC style.
//
// The original N3 program extends a Naive‑Bayes diagnostic model with a
// decision‑theoretic layer that scores five therapies by expected utility.
// It then recommends the therapy with the highest utility.
//
// This is intentionally not a generic N3 reasoner. The concrete N3 facts and
// rules are represented as ordinary Go data and functions so the probabilistic
// inference and decision logic are easy to read and directly runnable.
//
// Run:
//
//      go run bayes_therapy.go
//
// The program has no third-party dependencies.

package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "bayes_therapy"

// ---------- types ----------

// Disease holds the name and prior probability of one disease.
type Disease struct {
	Name  string
	Prior float64
}

// Therapy holds the name, success rate per disease, and adverse event rate.
type Therapy struct {
	Name             string
	SuccessByDisease []float64 // order must match the Case disease list
	Adverse          float64
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

// TherapyResult stores the computed values for a single therapy.
type TherapyResult struct {
	Therapy         string
	ExpectedSuccess float64
	ExpectedAdverse float64
	Utility         float64
}

// Dataset contains the fixed facts translated from bayes-therapy.n3.
type Dataset struct {
	Diseases      []Disease
	Therapies     []Therapy
	ProbGiven     map[string]map[string]float64
	Evidence      []EvidenceItem
	BenefitWeight float64
	HarmWeight    float64
}

// InferenceResult holds all derived facts needed by the renderer.
type InferenceResult struct {
	Scores          []float64
	EvidenceTotal   float64
	Posteriors      []float64
	PosteriorDetail []PosteriorResult
	TherapyResults  []TherapyResult
	BestTherapy     string
}

// Checks mirrors the proof obligations at the bottom of bayes-therapy.n3.
// In the N3 source, guards fire if any probability is outside [0,1].
type Checks struct {
	PriorsInRange        bool
	CondProbsInRange     bool
	AdverseInRange       bool
	SuccessInRange       bool
	EvidenceTotalNonZero bool
	DiseaseCountMatch    bool
	TherapyCountMatch    bool
}

// ---------- helpers ----------

// probInRange returns true if 0 <= p <= 1.
func probInRange(p float64) bool {
	return p >= 0 && p <= 1
}

// factor returns the likelihood factor for one disease and one evidence item.
func factor(d Disease, e EvidenceItem, probGiven map[string]map[string]float64) float64 {
	p := probGiven[d.Name][e.Symptom]
	if e.Present {
		return p
	}
	return 1 - p
}

// product returns the product of a slice of floats.
func product(vals []float64) float64 {
	prod := 1.0
	for _, v := range vals {
		prod *= v
	}
	return prod
}

// sum returns the sum of a slice of floats.
func sum(vals []float64) float64 {
	total := 0.0
	for _, v := range vals {
		total += v
	}
	return total
}

// zip returns elementwise pairs as a slice of [2]float64.
func zip(a, b []float64) [][2]float64 {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	pairs := make([][2]float64, n)
	for i := 0; i < n; i++ {
		pairs[i] = [2]float64{a[i], b[i]}
	}
	return pairs
}

// formatEvidence returns a human‑readable string of the evidence list.
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

// round6 rounds a float to 6 decimal places for display.
func round6(v float64) float64 {
	return float64(int(v*1e6+0.5)) / 1e6
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
		Therapies: []Therapy{
			{"Paxlovid", []float64{0.75, 0.05, 0.02, 0.05}, 0.10},
			{"Oseltamivir", []float64{0.05, 0.60, 0.02, 0.05}, 0.08},
			{"Antihistamine", []float64{0.10, 0.10, 0.75, 0.05}, 0.03},
			{"Antibiotic", []float64{0.05, 0.05, 0.02, 0.80}, 0.07},
			{"SupportiveCare", []float64{0.30, 0.30, 0.25, 0.20}, 0.01},
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
			{"LossOfSmell", false},
			{"Sneezing", false},
			{"ShortBreath", false},
		},
		BenefitWeight: 10,
		HarmWeight:    3,
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
				fmt.Fprintf(os.Stderr, "guard failed: pGiven(%s,%s) = %v\n", dname, sname, p)
				os.Exit(1)
			}
		}
	}
	for _, t := range data.Therapies {
		if !probInRange(t.Adverse) {
			fmt.Fprintf(os.Stderr, "guard failed: adverse(%s) = %v\n", t.Name, t.Adverse)
			os.Exit(1)
		}
		for _, p := range t.SuccessByDisease {
			if !probInRange(p) {
				fmt.Fprintf(os.Stderr, "guard failed: success(%s) = %v\n", t.Name, p)
				os.Exit(1)
			}
		}
	}
}

// ---------- inference ----------

func infer(data Dataset) InferenceResult {
	runGuards(data)

	// 1) Compute unnormalized scores for each disease.
	scores := make([]float64, len(data.Diseases))
	for i, d := range data.Diseases {
		// Build the list of factors for this disease vs all evidence items.
		factors := make([]float64, len(data.Evidence))
		for j, e := range data.Evidence {
			factors[j] = factor(d, e, data.ProbGiven)
		}
		// score = prior * product(factors)
		scores[i] = d.Prior * product(factors)
	}

	// 2) Evidence total = sum of scores.
	total := sum(scores)

	// 3) Posteriors = score / total.
	posteriors := make([]float64, len(scores))
	for i, s := range scores {
		if total > 0 {
			posteriors[i] = s / total
		}
	}

	// Build a human‑readable posterior detail list.
	posteriorDetail := make([]PosteriorResult, len(data.Diseases))
	for i, d := range data.Diseases {
		posteriorDetail[i] = PosteriorResult{
			Disease:      d.Name,
			Unnormalized: scores[i],
			Posterior:    posteriors[i],
		}
	}

	// 4) Compute expected success and utility for each therapy.
	therapyResults := make([]TherapyResult, len(data.Therapies))
	bestUtility := -1e12
	bestTherapy := ""
	for i, t := range data.Therapies {
		// expectedSuccess = sum(posterior[j] * successByDisease[j])
		expectedSuccess := 0.0
		for j, post := range posteriors {
			expectedSuccess += post * t.SuccessByDisease[j]
		}
		// utility = benefitWeight * expectedSuccess - harmWeight * adverse
		utility := data.BenefitWeight*expectedSuccess - data.HarmWeight*t.Adverse
		therapyResults[i] = TherapyResult{
			Therapy:         t.Name,
			ExpectedSuccess: expectedSuccess,
			ExpectedAdverse: t.Adverse,
			Utility:         utility,
		}
		if utility > bestUtility {
			bestUtility = utility
			bestTherapy = t.Name
		}
	}

	return InferenceResult{
		Scores:          scores,
		EvidenceTotal:   total,
		Posteriors:      posteriors,
		PosteriorDetail: posteriorDetail,
		TherapyResults:  therapyResults,
		BestTherapy:     bestTherapy,
	}
}

// ---------- checks ----------

func performChecks(data Dataset, result InferenceResult) Checks {
	c := Checks{}

	// Check priors are in [0,1].
	c.PriorsInRange = true
	for _, d := range data.Diseases {
		if !probInRange(d.Prior) {
			c.PriorsInRange = false
			break
		}
	}

	// Check conditional probabilities are in [0,1].
	c.CondProbsInRange = true
	for _, syms := range data.ProbGiven {
		for _, p := range syms {
			if !probInRange(p) {
				c.CondProbsInRange = false
				break
			}
		}
	}

	// Check adverse rates are in [0,1].
	c.AdverseInRange = true
	for _, t := range data.Therapies {
		if !probInRange(t.Adverse) {
			c.AdverseInRange = false
			break
		}
	}

	// Check success rates are in [0,1].
	c.SuccessInRange = true
	for _, t := range data.Therapies {
		for _, p := range t.SuccessByDisease {
			if !probInRange(p) {
				c.SuccessInRange = false
				break
			}
		}
	}

	// Check that evidence total is non‑zero.
	c.EvidenceTotalNonZero = result.EvidenceTotal > 0

	// Check that the number of diseases matches the success list length.
	c.DiseaseCountMatch = len(data.Diseases) == len(data.Therapies[0].SuccessByDisease)

	// Check that the number of therapies matches the expected count.
	c.TherapyCountMatch = len(data.Therapies) == 5

	return c
}

func allChecksPass(c Checks) bool {
	return c.PriorsInRange && c.CondProbsInRange && c.AdverseInRange &&
		c.SuccessInRange && c.EvidenceTotalNonZero && c.DiseaseCountMatch &&
		c.TherapyCountMatch
}

func checkCount(c Checks) int {
	count := 0
	if c.PriorsInRange {
		count++
	}
	if c.CondProbsInRange {
		count++
	}
	if c.AdverseInRange {
		count++
	}
	if c.SuccessInRange {
		count++
	}
	if c.EvidenceTotalNonZero {
		count++
	}
	if c.DiseaseCountMatch {
		count++
	}
	if c.TherapyCountMatch {
		count++
	}
	return count
}

// ---------- rendering ----------

func renderArcOutput(data Dataset, result InferenceResult, checks Checks) {
	fmt.Println("# Bayes Therapy Decision Support")
	fmt.Println()

	// --- Answer ---
	fmt.Println("## Answer")
	fmt.Printf("The recommended therapy is %s (utility = %.6f).\n",
		result.BestTherapy,
		therapyUtility(result.TherapyResults, result.BestTherapy))
	fmt.Println()
	fmt.Println("Full posterior distribution:")
	for _, r := range result.PosteriorDetail {
		fmt.Printf("  %-20s  posterior = %.6f  (unnormalized = %.8f)\n",
			r.Disease, r.Posterior, r.Unnormalized)
	}
	fmt.Println()
	fmt.Println("Therapy utility scores:")
	for _, tr := range result.TherapyResults {
		fmt.Printf("  %-20s  expectedSuccess = %.6f  adverse = %.2f  utility = %.6f\n",
			tr.Therapy, tr.ExpectedSuccess, tr.ExpectedAdverse, tr.Utility)
	}
	fmt.Println()

	// --- Reason Why ---
	fmt.Println("## Reason why")
	fmt.Printf("Evidence: %s.\n", formatEvidence(data.Evidence))
	fmt.Printf("Evidence total (normalizing constant) = %.8f.\n", result.EvidenceTotal)
	fmt.Println()
	fmt.Println("The posterior for each disease is computed as:")
	fmt.Println("  posterior(d) = prior(d) × ∏ P(symptom|d) / evidenceTotal")
	fmt.Println("where for an absent symptom the factor is 1 − P(symptom|d).")
	fmt.Println()
	fmt.Println("For each therapy, expected success is:")
	fmt.Println("  expectedSuccess(t) = Σ_i posterior(i) × successByDisease(i)")
	fmt.Println("and utility = benefitWeight × expectedSuccess − harmWeight × adverse.")
	fmt.Println("The recommended therapy is the one with the highest utility.")
	fmt.Println()

	// --- Check ---
	fmt.Println("## Check")
	if checks.PriorsInRange {
		fmt.Println("C1 OK - all prior probabilities are in [0,1].")
	} else {
		fmt.Println("C1 FAIL - one or more prior probabilities are outside [0,1].")
	}
	if checks.CondProbsInRange {
		fmt.Println("C2 OK - all conditional probabilities are in [0,1].")
	} else {
		fmt.Println("C2 FAIL - one or more conditional probabilities are outside [0,1].")
	}
	if checks.AdverseInRange {
		fmt.Println("C3 OK - all adverse probabilities are in [0,1].")
	} else {
		fmt.Println("C3 FAIL - one or more adverse probabilities are outside [0,1].")
	}
	if checks.SuccessInRange {
		fmt.Println("C4 OK - all success probabilities are in [0,1].")
	} else {
		fmt.Println("C4 FAIL - one or more success probabilities are outside [0,1].")
	}
	if checks.EvidenceTotalNonZero {
		fmt.Println("C5 OK - evidence total is non‑zero.")
	} else {
		fmt.Println("C5 FAIL - evidence total is zero.")
	}
	if checks.DiseaseCountMatch {
		fmt.Println("C6 OK - number of diseases matches success list length.")
	} else {
		fmt.Println("C6 FAIL - disease count mismatch.")
	}
	if checks.TherapyCountMatch {
		fmt.Println("C7 OK - number of therapies is correct.")
	} else {
		fmt.Println("C7 FAIL - therapy count mismatch.")
	}
	fmt.Println()

	// --- Go audit details ---
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("diseases : %d\n", len(data.Diseases))
	fmt.Printf("symptoms : %d\n", len(data.ProbGiven[data.Diseases[0].Name]))
	fmt.Printf("therapies : %d\n", len(data.Therapies))
	fmt.Printf("evidence items : %d\n", len(data.Evidence))
	fmt.Printf("evidence total : %.8f\n", result.EvidenceTotal)
	fmt.Println("posteriors :")
	for _, r := range result.PosteriorDetail {
		fmt.Printf("  %-20s  unnormalized=%.8f  posterior=%.6f\n",
			r.Disease, r.Unnormalized, r.Posterior)
	}
	fmt.Println("therapy scores :")
	for _, tr := range result.TherapyResults {
		fmt.Printf("  %-20s  expSucc=%.6f  adverse=%.2f  utility=%.6f\n",
			tr.Therapy, tr.ExpectedSuccess, tr.ExpectedAdverse, tr.Utility)
	}
	fmt.Printf("best therapy : %s\n", result.BestTherapy)
	fmt.Printf("checks passed : %d/7\n", checkCount(checks))
	fmt.Printf("recommendation consistent : %s\n", yesNo(allChecksPass(checks)))
}

func therapyUtility(results []TherapyResult, name string) float64 {
	for _, r := range results {
		if r.Therapy == name {
			return r.Utility
		}
	}
	return 0
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
	checks := performChecks(data, result)

	// Render ARC‑style output.
	renderArcOutput(data, result, checks)
}
