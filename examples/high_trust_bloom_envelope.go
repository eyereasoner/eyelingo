// high_trust_bloom_envelope.go
//
// A self-contained Go translation inspired by Eyeling's
// `examples/high-trust-rdf-bloom-envelope.n3`.
//
// without making correctness depend on it. Exact maybe-positive results must be
// confirmed against the canonical graph, and a decimal certificate bounds the
// expected false-positive workload.
//
// Run:
//
//	go run examples/high_trust_bloom_envelope.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
)

const eyelingoExampleName = "high_trust_bloom_envelope"

type Dataset struct {
	CaseName string
	Question string
	Artifact Artifact
	Policies Policies
	Expected Expected
}

type Artifact struct {
	ID                        string
	CanonicalTripleCount      int
	SPOIndexTripleCount       int
	BloomBits                 int
	HashFunctions             int
	NegativeLookupsPerBatch   int
	FPRateBudget              float64
	ExtraExactLookupsBudget   float64
	ExactTranscendentalSymbol string
	CertifiedLambda           float64
	ExpMinusLambdaLower       float64
	ExpMinusLambdaUpper       float64
}

type Policies struct {
	MaybePositivePolicy    string
	DefiniteNegativePolicy string
}

type Expected struct {
	Decision string
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	ParameterSanity          bool
	IndexAgreement           bool
	Lambda                   float64
	IntervalCertified        bool
	FPRateLower              float64
	FPRateUpper              float64
	ExpectedExtraLookupUpper float64
	WithinFPRateBudget       bool
	WithinExactLookupBudget  bool
	Decision                 string
	Checks                   []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func derive(ds Dataset) Analysis {
	a := ds.Artifact
	parameterSanity := a.CanonicalTripleCount > 0 && a.SPOIndexTripleCount > 0 && a.BloomBits > 0 && a.HashFunctions > 0 && a.NegativeLookupsPerBatch > 0 && a.FPRateBudget > 0 && a.ExtraExactLookupsBudget > 0
	indexAgreement := a.CanonicalTripleCount == a.SPOIndexTripleCount
	lambda := float64(a.HashFunctions*a.CanonicalTripleCount) / float64(a.BloomBits)
	intervalCertified := lambda > 0 && almostEqual(lambda, a.CertifiedLambda) && a.ExpMinusLambdaLower < a.ExpMinusLambdaUpper && a.ExpMinusLambdaLower > 0 && a.ExpMinusLambdaUpper < 1
	fpLower := math.Pow(1-a.ExpMinusLambdaUpper, float64(a.HashFunctions))
	fpUpper := math.Pow(1-a.ExpMinusLambdaLower, float64(a.HashFunctions))
	extraUpper := float64(a.NegativeLookupsPerBatch) * fpUpper
	withinFPRateBudget := fpUpper > 0 && fpUpper < a.FPRateBudget
	withinExactLookupBudget := extraUpper > 0 && extraUpper < a.ExtraExactLookupsBudget

	decision := "RejectForHighTrustUse"
	if parameterSanity && indexAgreement && intervalCertified && withinFPRateBudget && withinExactLookupBudget && ds.Policies.MaybePositivePolicy == "ConfirmAgainstCanonicalGraph" && ds.Policies.DefiniteNegativePolicy == "ReturnAbsent" {
		decision = "AcceptForHighTrustUse"
	}

	analysis := Analysis{
		ParameterSanity:          parameterSanity,
		IndexAgreement:           indexAgreement,
		Lambda:                   lambda,
		IntervalCertified:        intervalCertified,
		FPRateLower:              fpLower,
		FPRateUpper:              fpUpper,
		ExpectedExtraLookupUpper: extraUpper,
		WithinFPRateBudget:       withinFPRateBudget,
		WithinExactLookupBudget:  withinExactLookupBudget,
		Decision:                 decision,
	}
	analysis.Checks = []Check{
		{ID: "C1", OK: parameterSanity, Text: "numeric Bloom and workload parameters are positive"},
		{ID: "C2", OK: indexAgreement, Text: fmt.Sprintf("canonical graph and SPO index agree on %d triples", a.CanonicalTripleCount)},
		{ID: "C3", OK: almostEqual(lambda, a.CertifiedLambda), Text: fmt.Sprintf("derived lambda %.10f matches the certified lambda", lambda)},
		{ID: "C4", OK: intervalCertified, Text: fmt.Sprintf("decimal interval %.10f..%.10f is a valid exp(-lambda) certificate", a.ExpMinusLambdaLower, a.ExpMinusLambdaUpper)},
		{ID: "C5", OK: withinFPRateBudget, Text: fmt.Sprintf("false-positive upper bound %.9f is below %.3f", fpUpper, a.FPRateBudget)},
		{ID: "C6", OK: withinExactLookupBudget, Text: fmt.Sprintf("expected extra exact lookups %.3f stay below %.1f", extraUpper, a.ExtraExactLookupsBudget)},
		{ID: "C7", OK: ds.Policies.MaybePositivePolicy == "ConfirmAgainstCanonicalGraph", Text: "maybe-positive Bloom hits are confirmed against the canonical graph"},
		{ID: "C8", OK: ds.Policies.DefiniteNegativePolicy == "ReturnAbsent", Text: "definite Bloom negatives may be returned absent without exact lookup"},
		{ID: "C9", OK: decision == ds.Expected.Decision, Text: fmt.Sprintf("deployment decision is %s", decision)},
	}
	return analysis
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-12
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
	fmt.Println("# High Trust RDF Bloom Envelope")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("Deployment decision : %s for %s.\n", analysis.Decision, ds.Artifact.ID)
	fmt.Printf("lambda : %.10f\n", analysis.Lambda)
	fmt.Printf("false-positive envelope : %.9f .. %.9f\n", analysis.FPRateLower, analysis.FPRateUpper)
	fmt.Printf("expected extra exact lookups upper : %.3f per %d negative lookups\n", analysis.ExpectedExtraLookupUpper, ds.Artifact.NegativeLookupsPerBatch)
	fmt.Printf("maybe-positive policy : %s\n", ds.Policies.MaybePositivePolicy)
	fmt.Printf("definite-negative policy : %s\n", ds.Policies.DefiniteNegativePolicy)
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	fmt.Println("The canonical graph and the SPO index have the same triple count, so exact membership remains grounded in the graph snapshot.")
	fmt.Printf("The Bloom prefilter has n=%d triples, m=%d bits, and k=%d hash functions, giving lambda %.10f.\n", ds.Artifact.CanonicalTripleCount, ds.Artifact.BloomBits, ds.Artifact.HashFunctions, analysis.Lambda)
	fmt.Printf("Instead of asking the engine to know %s exactly, the input carries a decimal interval certificate for exp(-lambda).\n", ds.Artifact.ExactTranscendentalSymbol)
	fmt.Printf("That certificate bounds (1-exp(-lambda))^k below the %.3f false-positive budget and keeps extra exact confirmations below %.1f.\n", ds.Artifact.FPRateBudget, ds.Artifact.ExtraExactLookupsBudget)
	fmt.Println("Correctness never depends on the Bloom filter alone: maybe-positive answers must be confirmed against the canonical graph.")
	fmt.Println()
}
