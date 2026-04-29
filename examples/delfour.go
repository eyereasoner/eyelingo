// delfour.go
//
// A self-contained Go translation of delfour.js, which itself is a concrete
// JavaScript translation of examples/delfour.n3.
//
// The scenario models a privacy-preserving retail interaction. A phone turns a
// sensitive household condition into a narrow, expiring shopping insight; a
// store scanner checks the envelope and policy; and, if authorized, suggests a
// lower-sugar product.
//
// This is intentionally not a generic RDF/N3 reasoner. The concrete N3 facts
// and rules are represented as Go structs and ordinary functions so the data
// flow is clear and directly runnable.
//
// Run:
//
//	go run delfour.go
//
// The program has no third-party dependencies. It only uses the standard
// library to compute the SHA-256 payload hash used by the fixture checks.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "delfour"

const (
	insightID              = "https://example.org/insight/delfour"
	expectedFilesWritten   = 6
	expectedPayloadSHA256  = "34ad35638dfd7c67d031eeca8abb235ec24280740f863f3f31cd9d7b6517f098"
	trustedPrecomputedHMAC = "trusted-precomputed-input"
)

// The N3 source stores JSON inside an N3 string literal. In the raw file, the
// quotes inside that string are escaped. The fixture hashes that escaped string
// form, so the program stores readable JSON and escapes it before hashing.
const canonicalJSON = `{"insight":{"createdAt":"2025-10-05T20:33:48.907163+00:00","expiresAt":"2025-10-05T22:33:48.907185+00:00","id":"https://example.org/insight/delfour","metric":"sugar_g_per_serving","retailer":"Delfour","scopeDevice":"self-scanner","scopeEvent":"pick_up_scanner","suggestionPolicy":"lower_metric_first_higher_price_ok","threshold":10.0,"type":"ins:Insight"},"policy":{"duty":{"action":"odrl:delete","constraint":{"leftOperand":"odrl:dateTime","operator":"odrl:eq","rightOperand":"2025-10-05T22:33:48.907185+00:00"}},"permission":{"action":"odrl:use","constraint":{"leftOperand":"odrl:purpose","operator":"odrl:eq","rightOperand":"shopping_assist"},"target":"https://example.org/insight/delfour"},"profile":"Delfour-Insight-Policy","prohibition":{"action":"odrl:distribute","constraint":{"leftOperand":"odrl:purpose","operator":"odrl:eq","rightOperand":"marketing"},"target":"https://example.org/insight/delfour"},"type":"odrl:Policy"}}`

// Dataset is the Go equivalent of the JavaScript dataset() object.
//
// Names such as "odrl:use" are kept as strings because this is a concrete
// translation of one example, not a namespace-aware RDF implementation.
type Dataset struct {
	Case       CaseFacts
	Products   []Product
	Household  HouseholdFacts
	Scan       ScanFacts
	Insight    Insight
	Policy     Policy
	Signature  Signature
	ReasonText string
}

// CaseFacts contains the runtime request and audit facts for the fixture.
type CaseFacts struct {
	CaseName       string
	RequestPurpose string
	RequestAction  string
	ScannerAuthAt  string
	ScannerDutyAt  string
	FilesWritten   int
	AuditEntries   int
}

// Product is a small product catalog entry. SugarTenths stores grams-per-serving
// as tenths of a gram so integer comparisons can mirror the N3 rules exactly.
type Product struct {
	ID              string
	Name            string
	SugarTenths     int
	SugarPerServing float64
}

// HouseholdFacts contains the sensitive local fact that should not be sent to
// the scanner directly.
type HouseholdFacts struct {
	Condition string
}

// ScanFacts identifies the product currently scanned by the shopper.
type ScanFacts struct {
	ScannedProductID string
}

// Insight is the minimized, scoped envelope derived from the sensitive local
// fact. DerivedFromNeed is filled by Infer when the desensitization rule fires.
type Insight struct {
	ID                  string
	Metric              string
	ThresholdTenths     int
	ThresholdDisplay    string
	ThresholdG          float64
	SuggestionPolicy    string
	ScopeDevice         string
	ScopeEvent          string
	Retailer            string
	CreatedAt           string
	ExpiresAt           string
	SerializedLowercase string
	DerivedFromNeed     string
}

// Constraint captures the simple ODRL-style constraints used by this fixture.
type Constraint struct {
	LeftOperand  string
	Operator     string
	RightOperand string
}

// PolicyRule is used for permission and prohibition blocks.
type PolicyRule struct {
	Action     string
	Target     string
	Constraint Constraint
}

// DutyRule is separated from PolicyRule because the fixture's duty has no
// target field; it only requires a delete action by a given time.
type DutyRule struct {
	Action     string
	Constraint Constraint
}

// Policy is the minimal policy envelope that constrains how the scanner may use
// the insight.
type Policy struct {
	Profile     string
	Permission  PolicyRule
	Prohibition PolicyRule
	Duty        DutyRule
}

// Signature stores fixture metadata. The HMAC is trusted precomputed input in
// this example; only the payload SHA-256 is recomputed locally.
type Signature struct {
	Alg                  string
	KeyID                string
	Created              string
	PayloadHashSHA256    string
	SignatureHMAC        string
	HMACVerificationMode string
}

// Decision records the authorization outcome that the renderer prints.
type Decision struct {
	At      string
	Outcome string
	Target  string
}

// Banner is the shopper-facing message inferred from a valid insight and a
// high-sugar scanned product.
type Banner struct {
	Headline             string
	Note                 string
	SuggestedAlternative string
}

// Checks mirrors the Check section of the N3 output. Each field is a named proof
// condition rather than an anonymous boolean.
type Checks struct {
	PayloadHashMatches               bool
	SignatureVerifies                bool
	MinimizationStripsSensitiveTerms bool
	ScopeComplete                    bool
	AuthorizationAllowed             bool
	BannerFlagsHighSugar             bool
	AlternativeIsLowerSugar          bool
	DutyTimingConsistent             bool
	MarketingProhibited              bool
	FilesWrittenExpected             bool
}

// InferenceResult is the combined result of running the translated rules.
type InferenceResult struct {
	NeedsLowSugar        bool
	Decision             Decision
	Scanned              Product
	SuggestedAlternative Product
	Banner               Banner
	Checks               Checks
}

// dataset returns the hard-coded facts translated from delfour.js / delfour.n3.
func dataset() Dataset {
	return Dataset{
		Case: CaseFacts{
			CaseName:       "delfour",
			RequestPurpose: "shopping_assist",
			RequestAction:  "odrl:use",
			ScannerAuthAt:  "2025-10-05T20:35:48.907163+00:00",
			ScannerDutyAt:  "2025-10-05T20:37:48.907163+00:00",
			FilesWritten:   6,
			AuditEntries:   1,
		},

		Products: []Product{
			{
				ID:              "prod:BIS_001",
				Name:            "Classic Tea Biscuits",
				SugarTenths:     120,
				SugarPerServing: 12.0,
			},
			{
				ID:              "prod:BIS_101",
				Name:            "Low-Sugar Tea Biscuits",
				SugarTenths:     30,
				SugarPerServing: 3.0,
			},
			{
				ID:              "prod:CHOC_050",
				Name:            "Milk Chocolate Bar",
				SugarTenths:     150,
				SugarPerServing: 15.0,
			},
			{
				ID:              "prod:CHOC_150",
				Name:            "85% Dark Chocolate",
				SugarTenths:     60,
				SugarPerServing: 6.0,
			},
		},

		Household: HouseholdFacts{
			Condition: "Diabetes",
		},

		Scan: ScanFacts{
			ScannedProductID: "prod:BIS_001",
		},

		Insight: Insight{
			ID:                  insightID,
			Metric:              "sugar_g_per_serving",
			ThresholdTenths:     100,
			ThresholdDisplay:    "10.0",
			ThresholdG:          10.0,
			SuggestionPolicy:    "lower_metric_first_higher_price_ok",
			ScopeDevice:         "self-scanner",
			ScopeEvent:          "pick_up_scanner",
			Retailer:            "Delfour",
			CreatedAt:           "2025-10-05T20:33:48.907163+00:00",
			ExpiresAt:           "2025-10-05T22:33:48.907185+00:00",
			SerializedLowercase: `{"createdat":"2025-10-05t20:33:48.907163+00:00","expiresat":"2025-10-05t22:33:48.907185+00:00","id":"https://example.org/insight/delfour","metric":"sugar_g_per_serving","retailer":"delfour","scopedevice":"self-scanner","scopeevent":"pick_up_scanner","suggestionpolicy":"lower_metric_first_higher_price_ok","threshold":10.0,"type":"ins:insight"}`,
		},

		Policy: Policy{
			Profile: "Delfour-Insight-Policy",
			Permission: PolicyRule{
				Action: "odrl:use",
				Target: insightID,
				Constraint: Constraint{
					LeftOperand:  "odrl:purpose",
					Operator:     "odrl:eq",
					RightOperand: "shopping_assist",
				},
			},
			Prohibition: PolicyRule{
				Action: "odrl:distribute",
				Target: insightID,
				Constraint: Constraint{
					LeftOperand:  "odrl:purpose",
					Operator:     "odrl:eq",
					RightOperand: "marketing",
				},
			},
			Duty: DutyRule{
				Action: "odrl:delete",
				Constraint: Constraint{
					LeftOperand:  "odrl:dateTime",
					Operator:     "odrl:eq",
					RightOperand: "2025-10-05T22:33:48.907185+00:00",
				},
			},
		},

		Signature: Signature{
			Alg:                  "HMAC-SHA256",
			KeyID:                "demo-shared-secret",
			Created:              "2025-10-05T20:33:48.907163+00:00",
			PayloadHashSHA256:    expectedPayloadSHA256,
			SignatureHMAC:        "b21d0072d90112a9f820aced0286889f4b6ef92b145e6fdef1011f3bfa4608c2",
			HMACVerificationMode: trustedPrecomputedHMAC,
		},

		ReasonText: "Household requires low-sugar guidance (diabetes in POD). A neutral Insight is scoped to device 'self-scanner', event 'pick_up_scanner', retailer 'Delfour', and expires soon; the policy confines use to shopping assistance.\n",
	}
}

// infer runs the translated N3 rules in a deterministic order and returns the
// facts needed by the output renderer.
func infer(data Dataset) (InferenceResult, error) {
	scanned, ok := productByID(data, data.Scan.ScannedProductID)
	if !ok {
		return InferenceResult{}, fmt.Errorf("product not found: %s", data.Scan.ScannedProductID)
	}

	// Rule: desensitize(&["Diabetes"]) -> true.
	needsLowSugar := data.Household.Condition == "Diabetes"

	// Rule: derive_insight(...). Copy the fixture insight so inference can add
	// derived fields without mutating the original dataset.
	insight := data.Insight
	if needsLowSugar {
		insight.DerivedFromNeed = "low_sugar"
	}

	// Rule: choose the lower-sugar alternative with the smallest sugar score.
	suggestedAlternative, ok := chooseLowestSugarAlternative(data.Products, scanned)
	if !ok {
		return InferenceResult{}, errors.New("no lower-sugar alternative found")
	}

	// The N3 file has hard checks that fail the proof immediately. Keep them as
	// early Go errors so a broken fixture cannot silently render.
	if err := hardChecks(data, insight, scanned, suggestedAlternative); err != nil {
		return InferenceResult{}, err
	}

	checks := Checks{
		PayloadHashMatches:               sha256Hex(canonicalJSONEscaped()) == data.Signature.PayloadHashSHA256,
		SignatureVerifies:                data.Signature.HMACVerificationMode == trustedPrecomputedHMAC,
		MinimizationStripsSensitiveTerms: minimizationStripsSensitiveTerms(insight),
		ScopeComplete:                    scopeComplete(insight),
		AuthorizationAllowed:             authorizationAllowed(data, insight),
		BannerFlagsHighSugar:             authorizationAllowed(data, insight) && scanned.SugarPerServing >= insight.ThresholdG,
		AlternativeIsLowerSugar:          suggestedAlternative.SugarTenths < scanned.SugarTenths,
		DutyTimingConsistent:             notGreaterThanISOUTC(data.Case.ScannerDutyAt, insight.ExpiresAt),
		MarketingProhibited:              marketingProhibited(data.Policy),
		FilesWrittenExpected:             data.Case.FilesWritten == expectedFilesWritten,
	}

	if !checks.AllPass() {
		return InferenceResult{}, fmt.Errorf("not all checks passed: %+v", checks)
	}

	outcome := "Denied"
	if checks.AuthorizationAllowed {
		outcome = "Allowed"
	}

	return InferenceResult{
		NeedsLowSugar: needsLowSugar,
		Decision: Decision{
			At:      data.Case.ScannerAuthAt,
			Outcome: outcome,
			Target:  insightID,
		},
		Scanned:              scanned,
		SuggestedAlternative: suggestedAlternative,
		Banner: Banner{
			Headline:             "Track sugar per serving while you scan",
			Note:                 "High sugar",
			SuggestedAlternative: suggestedAlternative.Name,
		},
		Checks: checks,
	}, nil
}

// productByID finds the product referenced by the scanner facts.
func productByID(data Dataset, id string) (Product, bool) {
	for _, candidate := range data.Products {
		if candidate.ID == id {
			return candidate, true
		}
	}
	return Product{}, false
}

// chooseLowestSugarAlternative finds any product with less sugar than the
// scanned item and keeps the lowest-sugar candidate among them.
func chooseLowestSugarAlternative(products []Product, scanned Product) (Product, bool) {
	var best Product
	found := false

	for _, candidate := range products {
		if candidate.SugarTenths >= scanned.SugarTenths {
			continue
		}
		if !found || candidate.SugarTenths < best.SugarTenths {
			best = candidate
			found = true
		}
	}

	return best, found
}

// hardChecks translates N3 checks that are intended to stop the proof if core
// fixture invariants no longer hold.
func hardChecks(data Dataset, insight Insight, scanned Product, suggestedAlternative Product) error {
	if data.Case.FilesWritten != expectedFilesWritten {
		return fmt.Errorf("filesWritten was %d, expected %d", data.Case.FilesWritten, expectedFilesWritten)
	}

	if greaterThanISOUTC(data.Case.ScannerAuthAt, insight.ExpiresAt) {
		return errors.New("scanner authorization happened after insight expiry")
	}

	if suggestedAlternative.SugarTenths >= scanned.SugarTenths {
		return errors.New("suggested alternative is not lower sugar")
	}

	if data.Policy.Prohibition.Action != "odrl:distribute" {
		return errors.New("policy prohibition action must be odrl:distribute")
	}

	actualHash := sha256Hex(canonicalJSONEscaped())
	if actualHash != data.Signature.PayloadHashSHA256 {
		return fmt.Errorf("payload hash mismatch: actual %s, expected %s", actualHash, data.Signature.PayloadHashSHA256)
	}

	return nil
}

// minimizationStripsSensitiveTerms ensures the exported insight text does not
// expose medical terms such as the original household condition.
func minimizationStripsSensitiveTerms(insight Insight) bool {
	return !strings.Contains(insight.SerializedLowercase, "diabetes") &&
		!strings.Contains(insight.SerializedLowercase, "medical")
}

// scopeComplete checks that the insight is tied to a device, an event, and an
// expiry time before the scanner can use it.
func scopeComplete(insight Insight) bool {
	return insight.ScopeDevice != "" && insight.ScopeEvent != "" && insight.ExpiresAt != ""
}

// authorizationAllowed applies the ODRL-style permission rule to the current
// request and also requires that authorization happens before expiry.
func authorizationAllowed(data Dataset, insight Insight) bool {
	permission := data.Policy.Permission

	return permission.Action == data.Case.RequestAction &&
		permission.Target == insight.ID &&
		permission.Constraint.LeftOperand == "odrl:purpose" &&
		permission.Constraint.Operator == "odrl:eq" &&
		permission.Constraint.RightOperand == data.Case.RequestPurpose &&
		notGreaterThanISOUTC(data.Case.ScannerAuthAt, insight.ExpiresAt)
}

// marketingProhibited verifies that the policy explicitly blocks distribution
// for marketing purposes.
func marketingProhibited(policy Policy) bool {
	return policy.Prohibition.Action == "odrl:distribute" &&
		policy.Prohibition.Target == insightID &&
		policy.Prohibition.Constraint.RightOperand == "marketing"
}

// AllPass returns true only when every named check succeeded.
func (checks Checks) AllPass() bool {
	return checks.PayloadHashMatches &&
		checks.SignatureVerifies &&
		checks.MinimizationStripsSensitiveTerms &&
		checks.ScopeComplete &&
		checks.AuthorizationAllowed &&
		checks.BannerFlagsHighSugar &&
		checks.AlternativeIsLowerSugar &&
		checks.DutyTimingConsistent &&
		checks.MarketingProhibited &&
		checks.FilesWrittenExpected
}

// CountPassed returns how many named checks are true and the total number of
// checks represented by the struct.
func (checks Checks) CountPassed() (int, int) {
	values := []bool{
		checks.PayloadHashMatches,
		checks.SignatureVerifies,
		checks.MinimizationStripsSensitiveTerms,
		checks.ScopeComplete,
		checks.AuthorizationAllowed,
		checks.BannerFlagsHighSugar,
		checks.AlternativeIsLowerSugar,
		checks.DutyTimingConsistent,
		checks.MarketingProhibited,
		checks.FilesWrittenExpected,
	}

	passed := 0
	for _, value := range values {
		if value {
			passed++
		}
	}

	return passed, len(values)
}

// lowerSugarCandidates returns every catalog product that could satisfy the
// lower-sugar alternative rule before the lowest-sugar tie-break is applied.
func lowerSugarCandidates(products []Product, scanned Product) []Product {
	var candidates []Product
	for _, product := range products {
		if product.SugarTenths < scanned.SugarTenths {
			candidates = append(candidates, product)
		}
	}
	return candidates
}

// derivedNeedLabel keeps the audit output explicit about whether the
// desensitization rule produced the neutral low-sugar need.
func derivedNeedLabel(needsLowSugar bool) string {
	if needsLowSugar {
		return "low_sugar"
	}
	return "none"
}

// The fixture timestamps are all RFC3339 strings with the same UTC offset and
// precision. For that constrained data shape, lexicographic ordering is the same
// as chronological ordering, so no date/time package is required.
func notGreaterThanISOUTC(left, right string) bool {
	return left <= right
}

func greaterThanISOUTC(left, right string) bool {
	return left > right
}

// canonicalJSONEscaped returns the exact lexical string whose SHA-256 is stored
// by the source fixture.
func canonicalJSONEscaped() string {
	return escapeQuotesForN3LexicalForm(canonicalJSON)
}

// escapeQuotesForN3LexicalForm mirrors the JavaScript helper used before
// hashing the JSON text.
func escapeQuotesForN3LexicalForm(text string) string {
	return strings.ReplaceAll(text, `"`, `\"`)
}

// sha256Hex computes a lower-case hexadecimal SHA-256 digest.
func sha256Hex(text string) string {
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

// renderArcOutput renders the same answer/check style as the N3 log:outputString
// section, with a final Go-specific audit detail block.
func renderArcOutput(data Dataset, result InferenceResult) {
	fmt.Println("=== Answer ===")
	fmt.Printf("The scanner is allowed to use a neutral shopping insight and recommends %s instead of %s.\n", result.SuggestedAlternative.Name, result.Scanned.Name)
	fmt.Printf("case : %s\n", data.Case.CaseName)
	fmt.Printf("decision : %s\n", result.Decision.Outcome)
	fmt.Printf("scanned product : %s\n", result.Scanned.Name)
	fmt.Printf("suggested alternative: %s\n", result.SuggestedAlternative.Name)

	fmt.Println()
	fmt.Println("=== Reason Why ===")
	fmt.Println("The phone desensitizes a diabetes-related household condition into a scoped low-sugar need, wraps it in an expiring Insight + Policy envelope, signs it, and the scanner consumes that envelope for shopping assistance.")
	fmt.Printf("metric : %s\n", data.Insight.Metric)
	fmt.Printf("threshold : %s\n", data.Insight.ThresholdDisplay)
	fmt.Printf("scope : %s @ %s\n", data.Insight.ScopeDevice, data.Insight.ScopeEvent)
	fmt.Printf("retailer : %s\n", data.Insight.Retailer)
	fmt.Printf("signature alg : %s\n", data.Signature.Alg)
	fmt.Printf("banner headline : %s\n", result.Banner.Headline)
	fmt.Printf("expires at : %s\n", data.Insight.ExpiresAt)
	fmt.Printf("reason.txt : %s", data.ReasonText)
	fmt.Printf("audit entries : %d\n", data.Case.AuditEntries)
	fmt.Printf("bus files written : %d\n", data.Case.FilesWritten)

	fmt.Println()
	fmt.Println("=== Check ===")
	fmt.Printf("signature verifies : %s\n", yesNo(result.Checks.SignatureVerifies))
	fmt.Printf("payload hash matches : %s\n", yesNo(result.Checks.PayloadHashMatches))
	fmt.Printf("minimization strips sensitive terms: %s\n", yesNo(result.Checks.MinimizationStripsSensitiveTerms))
	fmt.Printf("scope complete : %s\n", yesNo(result.Checks.ScopeComplete))
	fmt.Printf("authorization allowed : %s\n", yesNo(result.Checks.AuthorizationAllowed))
	fmt.Printf("high-sugar banner : %s\n", yesNo(result.Checks.BannerFlagsHighSugar))
	fmt.Printf("alternative lowers sugar : %s\n", yesNo(result.Checks.AlternativeIsLowerSugar))
	fmt.Printf("duty timing consistent : %s\n", yesNo(result.Checks.DutyTimingConsistent))
	fmt.Printf("marketing prohibited : %s\n", yesNo(result.Checks.MarketingProhibited))

	renderGoAuditDetails(data, result)
}

// renderGoAuditDetails prints Go-side diagnostics that are not part of the N3
// answer, but help verify the translated facts, rule decisions, and checks.
func renderGoAuditDetails(data Dataset, result InferenceResult) {
	candidates := lowerSugarCandidates(data.Products, result.Scanned)
	canonicalEscaped := canonicalJSONEscaped()
	passedChecks, totalChecks := result.Checks.CountPassed()
	sugarDropTenths := result.Scanned.SugarTenths - result.SuggestedAlternative.SugarTenths

	fmt.Println()
	fmt.Println("=== Go audit details ===")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case facts : case=%s requestAction=%s requestPurpose=%s\n", data.Case.CaseName, data.Case.RequestAction, data.Case.RequestPurpose)
	fmt.Printf("catalog products : %d\n", len(data.Products))
	fmt.Printf("scanned product id : %s\n", data.Scan.ScannedProductID)
	fmt.Printf("scanned sugar : %.1fg per serving (%d tenths)\n", result.Scanned.SugarPerServing, result.Scanned.SugarTenths)
	fmt.Printf("threshold sugar : %.1fg per serving (%d tenths)\n", data.Insight.ThresholdG, data.Insight.ThresholdTenths)
	fmt.Printf("lower-sugar candidates : %d\n", len(candidates))
	for _, candidate := range candidates {
		fmt.Printf(" - %s sugar=%.1fg (%d tenths)\n", candidate.Name, candidate.SugarPerServing, candidate.SugarTenths)
	}
	fmt.Printf("selected alternative : %s\n", result.SuggestedAlternative.Name)
	fmt.Printf("sugar reduction : %.1fg per serving (%d tenths)\n", float64(sugarDropTenths)/10.0, sugarDropTenths)
	fmt.Printf("needs low sugar : %s\n", yesNo(result.NeedsLowSugar))
	fmt.Printf("derived need : %s\n", derivedNeedLabel(result.NeedsLowSugar))
	fmt.Printf("authorization window : %s <= %s -> %s\n", result.Decision.At, data.Insight.ExpiresAt, yesNo(notGreaterThanISOUTC(result.Decision.At, data.Insight.ExpiresAt)))
	fmt.Printf("duty deadline check : %s <= %s -> %s\n", data.Case.ScannerDutyAt, data.Insight.ExpiresAt, yesNo(result.Checks.DutyTimingConsistent))
	fmt.Printf("decision target : %s\n", result.Decision.Target)
	fmt.Printf("policy profile : %s\n", data.Policy.Profile)
	fmt.Printf("permission rule : action=%s target=%s purpose=%s\n", data.Policy.Permission.Action, data.Policy.Permission.Target, data.Policy.Permission.Constraint.RightOperand)
	fmt.Printf("prohibition rule : action=%s target=%s purpose=%s\n", data.Policy.Prohibition.Action, data.Policy.Prohibition.Target, data.Policy.Prohibition.Constraint.RightOperand)
	fmt.Printf("delete duty : action=%s due=%s\n", data.Policy.Duty.Action, data.Policy.Duty.Constraint.RightOperand)
	fmt.Printf("signature mode : %s\n", data.Signature.HMACVerificationMode)
	fmt.Printf("signature key id : %s\n", data.Signature.KeyID)
	fmt.Printf("signature hmac : %s\n", data.Signature.SignatureHMAC)
	fmt.Printf("canonical json bytes : %d\n", len(canonicalJSON))
	fmt.Printf("escaped payload bytes : %d\n", len(canonicalEscaped))
	fmt.Printf("payload sha256 : %s\n", sha256Hex(canonicalEscaped))
	fmt.Printf("expected sha256 : %s\n", data.Signature.PayloadHashSHA256)
	fmt.Printf("checks passed : %d/%d\n", passedChecks, totalChecks)
	fmt.Printf("all checks pass : %s\n", yesNo(result.Checks.AllPass()))
	fmt.Printf("audit entries : %d\n", data.Case.AuditEntries)
	fmt.Printf("files written : %d/%d\n", data.Case.FilesWritten, expectedFilesWritten)
}

// yesNo prints booleans in the N3 example's human-readable style.
func yesNo(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}

func main() {
	data := exampleinput.Load(eyelingoExampleName, dataset())
	result, err := infer(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delfour inference failed: %s\n", err)
		os.Exit(1)
	}

	renderArcOutput(data, result)
}
