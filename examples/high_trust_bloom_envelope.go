// high_trust_bloom_envelope.go
//
// A self-contained Go translation inspired by Eyeling's
// `examples/high-trust-rdf-bloom-envelope.n3`.
//
// The scenario checks whether an RDF graph artifact can use a Bloom prefilter
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
	"os"
	"runtime"
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
	printChecks(analysis)
	printAudit(ds, analysis)
	if !allChecksOK(analysis.Checks) {
		os.Exit(1)
	}
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
	fmt.Println("## Reason why")
	fmt.Println("The canonical graph and the SPO index have the same triple count, so exact membership remains grounded in the graph snapshot.")
	fmt.Printf("The Bloom prefilter has n=%d triples, m=%d bits, and k=%d hash functions, giving lambda %.10f.\n", ds.Artifact.CanonicalTripleCount, ds.Artifact.BloomBits, ds.Artifact.HashFunctions, analysis.Lambda)
	fmt.Printf("Instead of asking the engine to know %s exactly, the input carries a decimal interval certificate for exp(-lambda).\n", ds.Artifact.ExactTranscendentalSymbol)
	fmt.Printf("That certificate bounds (1-exp(-lambda))^k below the %.3f false-positive budget and keeps extra exact confirmations below %.1f.\n", ds.Artifact.FPRateBudget, ds.Artifact.ExtraExactLookupsBudget)
	fmt.Println("Correctness never depends on the Bloom filter alone: maybe-positive answers must be confirmed against the canonical graph.")
	fmt.Println()
}

func printChecks(analysis Analysis) {
	fmt.Println("## Check")
	for _, check := range analysis.Checks {
		status := "FAIL"
		if check.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", check.ID, status, check.Text)
	}
	fmt.Println()
}

func printAudit(ds Dataset, analysis Analysis) {
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("artifact : %s\n", ds.Artifact.ID)
	fmt.Printf("triple counts : canonical=%d spoIndex=%d agreement=%v\n", ds.Artifact.CanonicalTripleCount, ds.Artifact.SPOIndexTripleCount, analysis.IndexAgreement)
	fmt.Printf("bloom parameters : bits=%d hashFunctions=%d lambda=%.10f\n", ds.Artifact.BloomBits, ds.Artifact.HashFunctions, analysis.Lambda)
	fmt.Printf("certificate : symbol=%s certifiedLambda=%.10f expLower=%.10f expUpper=%.10f certified=%v\n", ds.Artifact.ExactTranscendentalSymbol, ds.Artifact.CertifiedLambda, ds.Artifact.ExpMinusLambdaLower, ds.Artifact.ExpMinusLambdaUpper, analysis.IntervalCertified)
	fmt.Printf("fp envelope : lower=%.9f upper=%.9f budget=%.3f within=%v\n", analysis.FPRateLower, analysis.FPRateUpper, ds.Artifact.FPRateBudget, analysis.WithinFPRateBudget)
	fmt.Printf("lookup budget : negativeLookups=%d extraUpper=%.3f budget=%.1f within=%v\n", ds.Artifact.NegativeLookupsPerBatch, analysis.ExpectedExtraLookupUpper, ds.Artifact.ExtraExactLookupsBudget, analysis.WithinExactLookupBudget)
	fmt.Printf("policies : maybePositive=%s definiteNegative=%s\n", ds.Policies.MaybePositivePolicy, ds.Policies.DefiniteNegativePolicy)
	fmt.Printf("checks passed : %d/%d\n", countChecksOK(analysis.Checks), len(analysis.Checks))
	fmt.Printf("decision : %s\n", analysis.Decision)
}
