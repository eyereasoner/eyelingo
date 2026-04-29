// digital_product_passport.go
//
// A self-contained Go translation inspired by Eyeling's
// `examples/digital-product-passport.n3`.
//
// The scenario models a Digital Product Passport for a smartphone. It derives
// compact public indicators from component, material, document, lifecycle, and
// footprint facts: total mass, recycled content, critical raw material exposure,
// lifecycle footprint, and a repair-friendly circularity hint.
//
// Run:
//
//	go run examples/digital_product_passport.go
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

const eyelingoExampleName = "digital_product_passport"

// Dataset contains the concrete DPP fixture: passport metadata, product facts,
// components, materials, documents, lifecycle events, and expected indicators.
type Dataset struct {
	CaseName     string
	Question     string
	Organization Organization
	Site         Site
	Repairer     Organization
	Passport     Passport
	AccessPolicy AccessPolicy
	Product      Product
	Components   []Component
	Materials    []Material
	Documents    []Document
	Lifecycle    []LifecycleEvent
	Footprint    Footprint
	Expected     ExpectedIndicators
}

type Organization struct {
	ID        string
	LegalName string
	VATID     string
	Country   string
}

type Site struct {
	ID      string
	Name    string
	Country string
}

type Passport struct {
	ID              string
	Version         string
	IssuedAt        string
	PublicEndpoint  string
	RestrictedScope string
}

type AccessPolicy struct {
	PublicSection       string
	PublicAudience      string
	RestrictedSection   string
	RestrictedAudience  string
	RestrictedDocTypes  []string
	PublicDocTypes      []string
	RequiredPublicClaim string
}

type Product struct {
	ID              string
	Model           string
	SerialNumber    string
	BatchID         string
	MadeBy          string
	MadeAtSite      string
	ManufactureDate string
	Category        string
	DigitalLink     string
}

type Component struct {
	ID               string
	Type             string
	MassG            int
	RecycledMassG    int
	Replaceable      bool
	ContainsMaterial []string
}

type Material struct {
	ID                  string
	CriticalRawMaterial bool
}

type Document struct {
	ID       string
	DocType  string
	URL      string
	Section  string
	Declares []string
}

type LifecycleEvent struct {
	ID                string
	Type              string
	ForProduct        string
	PerformedBy       string
	AtSite            string
	OnDate            string
	ReplacedComponent string
	Section           string
}

type Footprint struct {
	ManufacturingGCO2e int
	TransportGCO2e     int
	UsePhaseGCO2e      int
}

type ExpectedIndicators struct {
	TotalMassG          int
	RecycledMassG       int
	RecycledContentPct  int
	LifecycleGCO2e      int
	CircularityHint     string
	CriticalRawMaterial bool
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Decision          string
	TotalMassG        int
	RecycledMassG     int
	RecycledPct       int
	LifecycleGCO2e    int
	CriticalMaterials []string
	RepairFriendly    bool
	PublicDocs        []Document
	RestrictedDocs    []Document
	LatestEvent       LifecycleEvent
	Checks            []Check
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
	materials := map[string]Material{}
	for _, material := range ds.Materials {
		materials[material.ID] = material
	}

	totalMass := 0
	recycledMass := 0
	criticalSet := map[string]bool{}
	batteryReplaceable := false
	for _, component := range ds.Components {
		totalMass += component.MassG
		recycledMass += component.RecycledMassG
		if strings.EqualFold(component.Type, "Battery") && component.Replaceable {
			batteryReplaceable = true
		}
		for _, materialID := range component.ContainsMaterial {
			if material, ok := materials[materialID]; ok && material.CriticalRawMaterial {
				criticalSet[materialID] = true
			}
		}
	}

	criticalMaterials := make([]string, 0, len(criticalSet))
	for materialID := range criticalSet {
		criticalMaterials = append(criticalMaterials, materialID)
	}
	sort.Strings(criticalMaterials)

	recycledPct := 0
	if totalMass > 0 {
		recycledPct = recycledMass * 100 / totalMass
	}
	lifecycle := ds.Footprint.ManufacturingGCO2e + ds.Footprint.TransportGCO2e + ds.Footprint.UsePhaseGCO2e

	publicDocs, restrictedDocs := splitDocuments(ds.Documents, ds.AccessPolicy.PublicSection, ds.AccessPolicy.RestrictedSection)
	repairFriendly := batteryReplaceable && hasDocument(publicDocs, "RepairGuide") && hasDocument(publicDocs, "SparePartsCatalog") && declares(publicDocs, ds.AccessPolicy.RequiredPublicClaim)
	latest := latestLifecycleEvent(ds.Lifecycle)

	checks := []Check{
		{ID: "C1", OK: totalMass == ds.Expected.TotalMassG, Text: fmt.Sprintf("component masses fold to %d g", totalMass)},
		{ID: "C2", OK: recycledMass == ds.Expected.RecycledMassG, Text: fmt.Sprintf("recycled component masses fold to %d g", recycledMass)},
		{ID: "C3", OK: recycledPct == ds.Expected.RecycledContentPct, Text: fmt.Sprintf("integer recycled-content percentage is %d%%", recycledPct)},
		{ID: "C4", OK: lifecycle == ds.Expected.LifecycleGCO2e, Text: fmt.Sprintf("lifecycle footprint totals %d gCO2e", lifecycle)},
		{ID: "C5", OK: len(criticalMaterials) > 0 && ds.Expected.CriticalRawMaterial, Text: fmt.Sprintf("critical raw material exposure is %s", strings.Join(criticalMaterials, ", "))},
		{ID: "C6", OK: repairFriendly && ds.Expected.CircularityHint == "repairFriendly", Text: "replaceable battery plus repair and spare-parts documents derive repairFriendly"},
		{ID: "C7", OK: allDocTypesPresent(publicDocs, ds.AccessPolicy.PublicDocTypes), Text: "public section contains user manual, repair guide, and spare-parts catalog"},
		{ID: "C8", OK: allDocTypesPresent(restrictedDocs, ds.AccessPolicy.RestrictedDocTypes), Text: "restricted declarations stay in the restricted section"},
		{ID: "C9", OK: lifecycleOrderOK(ds.Lifecycle), Text: "manufacturing, sale, and repair events are chronologically consistent"},
		{ID: "C10", OK: ds.Passport.PublicEndpoint != "" && ds.Product.DigitalLink == ds.Passport.PublicEndpoint, Text: "passport endpoint matches the product digital link"},
	}

	decision := "REVIEW"
	if allChecksOK(checks) {
		decision = "PASS"
	}

	return Analysis{
		Decision:          decision,
		TotalMassG:        totalMass,
		RecycledMassG:     recycledMass,
		RecycledPct:       recycledPct,
		LifecycleGCO2e:    lifecycle,
		CriticalMaterials: criticalMaterials,
		RepairFriendly:    repairFriendly,
		PublicDocs:        publicDocs,
		RestrictedDocs:    restrictedDocs,
		LatestEvent:       latest,
		Checks:            checks,
	}
}

func splitDocuments(documents []Document, publicSection, restrictedSection string) ([]Document, []Document) {
	var publicDocs []Document
	var restrictedDocs []Document
	for _, document := range documents {
		switch document.Section {
		case publicSection:
			publicDocs = append(publicDocs, document)
		case restrictedSection:
			restrictedDocs = append(restrictedDocs, document)
		}
	}
	return publicDocs, restrictedDocs
}

func hasDocument(documents []Document, docType string) bool {
	for _, document := range documents {
		if document.DocType == docType {
			return true
		}
	}
	return false
}

func declares(documents []Document, claim string) bool {
	for _, document := range documents {
		for _, declared := range document.Declares {
			if declared == claim {
				return true
			}
		}
	}
	return false
}

func allDocTypesPresent(documents []Document, docTypes []string) bool {
	available := map[string]bool{}
	for _, document := range documents {
		available[document.DocType] = true
	}
	for _, docType := range docTypes {
		if !available[docType] {
			return false
		}
	}
	return true
}

func lifecycleOrderOK(events []LifecycleEvent) bool {
	order := map[string]int{"ManufacturingEvent": 0, "SaleEvent": 1, "RepairEvent": 2}
	lastDate := time.Time{}
	lastOrder := -1
	for _, event := range events {
		date := mustParseDate(event.ID, event.OnDate)
		if !lastDate.IsZero() && date.Before(lastDate) {
			return false
		}
		if eventOrder, ok := order[event.Type]; ok {
			if eventOrder < lastOrder {
				return false
			}
			lastOrder = eventOrder
		}
		lastDate = date
	}
	return true
}

func latestLifecycleEvent(events []LifecycleEvent) LifecycleEvent {
	if len(events) == 0 {
		return LifecycleEvent{}
	}
	latest := events[0]
	latestDate := mustParseDate(latest.ID, latest.OnDate)
	for _, event := range events[1:] {
		date := mustParseDate(event.ID, event.OnDate)
		if date.After(latestDate) {
			latest = event
			latestDate = date
		}
	}
	return latest
}

func mustParseDate(label, value string) time.Time {
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse %s date %q: %v\n", label, value, err)
		os.Exit(1)
	}
	return date
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
	fmt.Println("=== Digital Product Passport ===")
	fmt.Println()
	fmt.Println("=== Answer ===")
	fmt.Printf("Passport decision : %s for %s %s.\n", analysis.Decision, ds.Product.Model, ds.Product.SerialNumber)
	fmt.Printf("recycled content : %d%%\n", analysis.RecycledPct)
	fmt.Printf("lifecycle footprint : %d gCO2e\n", analysis.LifecycleGCO2e)
	fmt.Printf("total component mass : %d g\n", analysis.TotalMassG)
	fmt.Printf("critical raw materials : %s\n", strings.Join(analysis.CriticalMaterials, ", "))
	fmt.Println("circularity hint : repairFriendly")
	fmt.Printf("public endpoint : %s\n", ds.Passport.PublicEndpoint)
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("=== Reason Why ===")
	fmt.Println("The passport folds the explicit component list to derive total mass and recycled mass, then computes an integer recycled-content percentage.")
	fmt.Println("Lifecycle footprint is derived by summing manufacturing, transport, and use-phase emissions.")
	fmt.Println("The product is repair-friendly because the battery is replaceable and the public passport section exposes both repair and spare-parts documentation.")
	fmt.Printf("Critical raw-material exposure is derived from component-material links: %s.\n", strings.Join(analysis.CriticalMaterials, ", "))
	fmt.Println()
	fmt.Println("Component roll-up:")
	for _, component := range ds.Components {
		fmt.Printf(" - %s %s mass=%dg recycled=%dg materials=%s replaceable=%v\n",
			component.ID,
			component.Type,
			component.MassG,
			component.RecycledMassG,
			strings.Join(component.ContainsMaterial, ", "),
			component.Replaceable,
		)
	}
	fmt.Println("Public documents:")
	for _, document := range analysis.PublicDocs {
		fmt.Printf(" - %s %s %s\n", document.ID, document.DocType, document.URL)
	}
	fmt.Println()
}

func printChecks(analysis Analysis) {
	fmt.Println("=== Check ===")
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
	fmt.Println("=== Go audit details ===")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("passport : %s version=%s issued=%s endpoint=%s\n", ds.Passport.ID, ds.Passport.Version, ds.Passport.IssuedAt, ds.Passport.PublicEndpoint)
	fmt.Printf("product : %s model=%s serial=%s batch=%s category=%s\n", ds.Product.ID, ds.Product.Model, ds.Product.SerialNumber, ds.Product.BatchID, ds.Product.Category)
	fmt.Printf("manufacturer : %s (%s, %s) site=%s\n", ds.Organization.ID, ds.Organization.LegalName, ds.Organization.Country, ds.Site.ID)
	fmt.Printf("components : %d totalMassG=%d recycledMassG=%d recycledPct=%d%%\n", len(ds.Components), analysis.TotalMassG, analysis.RecycledMassG, analysis.RecycledPct)
	fmt.Printf("critical raw materials : %s\n", strings.Join(analysis.CriticalMaterials, ", "))
	fmt.Printf("documents : public=%d restricted=%d\n", len(analysis.PublicDocs), len(analysis.RestrictedDocs))
	fmt.Printf("lifecycle events : %d latest=%s %s %s\n", len(ds.Lifecycle), analysis.LatestEvent.ID, analysis.LatestEvent.Type, analysis.LatestEvent.OnDate)
	fmt.Printf("footprint : manufacturing=%d transport=%d use=%d total=%d\n", ds.Footprint.ManufacturingGCO2e, ds.Footprint.TransportGCO2e, ds.Footprint.UsePhaseGCO2e, analysis.LifecycleGCO2e)
	fmt.Printf("checks passed : %d/%d\n", countChecksOK(analysis.Checks), len(analysis.Checks))
	fmt.Printf("decision : %s\n", analysis.Decision)
}
