// eco_route_insight.go
//
// A compact insight-economy example: Go computes a local fuel-saving route
// suggestion and emits a small signed envelope, while Python independently
// checks the derivation during testing.
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"sort"
	"time"
)

const eyelingoExampleName = "eco_route_insight"

type Dataset struct {
	CaseName          string           `json:"caseName"`
	Scenario          Scenario         `json:"scenario"`
	Shipment          Shipment         `json:"shipment"`
	CurrentRoute      Route            `json:"currentRoute"`
	AlternativeRoutes []Route          `json:"alternativeRoutes"`
	Policy            Policy           `json:"policy"`
	DataMinimization  DataMinimization `json:"dataMinimization"`
	Signing           Signing          `json:"signing"`
}

type Scenario struct {
	Driver     string `json:"driver"`
	Depot      string `json:"depot"`
	AllowedUse string `json:"allowedUse"`
	IssuedAt   string `json:"issuedAt"`
	TTLHours   int    `json:"ttlHours"`
}

type Shipment struct {
	ID        string  `json:"id"`
	PayloadKg float64 `json:"payloadKg"`
}

type Route struct {
	ID             string  `json:"id"`
	Label          string  `json:"label"`
	DistanceKm     float64 `json:"distanceKm"`
	GradientFactor float64 `json:"gradientFactor"`
	ETAMinutes     int     `json:"etaMinutes"`
}

type Policy struct {
	FuelIndexThreshold               float64 `json:"fuelIndexThreshold"`
	MaxETADelayMinutes               int     `json:"maxEtaDelayMinutes"`
	RequireAlternativeBelowThreshold bool    `json:"requireAlternativeBelowThreshold"`
	AllowedUse                       string  `json:"allowedUse"`
	SignatureAlgorithm               string  `json:"signatureAlgorithm"`
}

type DataMinimization struct {
	ForbiddenEnvelopeTerms []string `json:"forbiddenEnvelopeTerms"`
}

type Signing struct {
	KeyID  string `json:"keyId"`
	Secret string `json:"secret"`
}

type RouteScore struct {
	Route     Route
	FuelIndex float64
	Saving    float64
	Delay     int
	Eligible  bool
}

type Envelope struct {
	Audience   string             `json:"audience"`
	AllowedUse string             `json:"allowedUse"`
	IssuedAt   string             `json:"issuedAt"`
	Expiry     string             `json:"expiry"`
	KeyID      string             `json:"keyId"`
	Assertions EnvelopeAssertions `json:"assertions"`
}

type EnvelopeAssertions struct {
	ShowEcoBanner      bool    `json:"showEcoBanner"`
	SuggestedRoute     string  `json:"suggestedRoute"`
	CurrentFuelIndex   float64 `json:"currentFuelIndex"`
	SuggestedFuelIndex float64 `json:"suggestedFuelIndex"`
	EstimatedSaving    float64 `json:"estimatedSaving"`
	RawDataExported    bool    `json:"rawDataExported"`
}

type Analysis struct {
	PayloadTonnes float64
	CurrentFuel   float64
	Alternatives  []RouteScore
	Selected      RouteScore
	IssueInsight  bool
	Envelope      Envelope
	Canonical     string
	PayloadDigest string
	Signature     string
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	analysis := derive(ds)
	printReport(ds, analysis)
}

func derive(ds Dataset) Analysis {
	payloadTonnes := ds.Shipment.PayloadKg / 1000.0
	currentFuel := fuelIndex(ds.CurrentRoute, payloadTonnes)
	alternatives := make([]RouteScore, 0, len(ds.AlternativeRoutes))
	for _, route := range ds.AlternativeRoutes {
		fi := fuelIndex(route, payloadTonnes)
		saving := currentFuel - fi
		delay := route.ETAMinutes - ds.CurrentRoute.ETAMinutes
		eligible := saving > 0 && delay <= ds.Policy.MaxETADelayMinutes
		if ds.Policy.RequireAlternativeBelowThreshold {
			eligible = eligible && fi <= ds.Policy.FuelIndexThreshold
		}
		alternatives = append(alternatives, RouteScore{Route: route, FuelIndex: fi, Saving: saving, Delay: delay, Eligible: eligible})
	}
	sort.Slice(alternatives, func(i, j int) bool {
		if alternatives[i].Eligible != alternatives[j].Eligible {
			return alternatives[i].Eligible
		}
		if !almostEqual(alternatives[i].Saving, alternatives[j].Saving) {
			return alternatives[i].Saving > alternatives[j].Saving
		}
		return alternatives[i].Route.ID < alternatives[j].Route.ID
	})

	selected := alternatives[0]
	issue := currentFuel > ds.Policy.FuelIndexThreshold && selected.Eligible
	expiry := expiryTime(ds.Scenario.IssuedAt, ds.Scenario.TTLHours)
	envelope := Envelope{
		Audience:   ds.Scenario.Depot,
		AllowedUse: ds.Policy.AllowedUse,
		IssuedAt:   ds.Scenario.IssuedAt,
		Expiry:     expiry,
		KeyID:      ds.Signing.KeyID,
		Assertions: EnvelopeAssertions{
			ShowEcoBanner:      issue,
			SuggestedRoute:     selected.Route.ID,
			CurrentFuelIndex:   round2(currentFuel),
			SuggestedFuelIndex: round2(selected.FuelIndex),
			EstimatedSaving:    round2(selected.Saving),
			RawDataExported:    false,
		},
	}
	canonical := stableJSON(envelope)
	digestBytes := sha256.Sum256([]byte(canonical))
	signature := sign(ds.Signing.Secret, canonical)

	return Analysis{
		PayloadTonnes: payloadTonnes,
		CurrentFuel:   currentFuel,
		Alternatives:  alternatives,
		Selected:      selected,
		IssueInsight:  issue,
		Envelope:      envelope,
		Canonical:     canonical,
		PayloadDigest: hex.EncodeToString(digestBytes[:]),
		Signature:     signature,
	}
}

func fuelIndex(route Route, payloadTonnes float64) float64 {
	return route.DistanceKm * payloadTonnes * route.GradientFactor
}

func expiryTime(issuedAt string, ttlHours int) string {
	parsed, err := time.Parse(time.RFC3339, issuedAt)
	if err != nil {
		panic(err)
	}
	return parsed.Add(time.Duration(ttlHours) * time.Hour).Format(time.RFC3339)
}

func stableJSON(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(encoded)
}

func sign(secret string, canonical string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonical))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func round2(x float64) float64 {
	return math.Round(x*100) / 100
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.0000001
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}

func printReport(ds Dataset, a Analysis) {
	fmt.Println("# Eco Route Insight")
	fmt.Println()
	fmt.Println("## Answer")
	if a.IssueInsight {
		fmt.Println("insight status : issue")
	} else {
		fmt.Println("insight status : hold")
	}
	fmt.Printf("show eco banner : %s\n", yesNo(a.Envelope.Assertions.ShowEcoBanner))
	fmt.Printf("audience : %s\n", a.Envelope.Audience)
	fmt.Printf("allowed use : %s\n", a.Envelope.AllowedUse)
	fmt.Printf("suggested route : %s\n", a.Envelope.Assertions.SuggestedRoute)
	fmt.Printf("current fuel index : %.2f\n", a.Envelope.Assertions.CurrentFuelIndex)
	fmt.Printf("suggested fuel index : %.2f\n", a.Envelope.Assertions.SuggestedFuelIndex)
	fmt.Printf("estimated saving : %.2f\n", a.Envelope.Assertions.EstimatedSaving)
	fmt.Printf("expires at : %s\n", a.Envelope.Expiry)
	fmt.Printf("raw data exported : %s\n", yesNo(a.Envelope.Assertions.RawDataExported))
	fmt.Printf("signature algorithm : %s\n", ds.Policy.SignatureAlgorithm)
	fmt.Printf("payload digest : %s\n", a.PayloadDigest)
	fmt.Printf("signature key : %s\n", a.Envelope.KeyID)
	fmt.Printf("signature : %s\n", a.Signature)
	fmt.Println()
	fmt.Println("## Reason")
	fmt.Println("The current route uses fuel index = distanceKm × (payloadKg / 1000) × gradientFactor.")
	fmt.Printf("For %s, %s gives %.2f × %.2f × %.2f = %.2f.\n", ds.Shipment.ID, ds.CurrentRoute.Label, ds.CurrentRoute.DistanceKm, a.PayloadTonnes, ds.CurrentRoute.GradientFactor, a.CurrentFuel)
	fmt.Printf("The policy threshold is %.2f, so a local eco banner is justified.\n", ds.Policy.FuelIndexThreshold)
	fmt.Printf("The selected alternative %s gives %.2f × %.2f × %.2f = %.2f, saving %.2f while staying within the ETA delay limit.\n", a.Selected.Route.ID, a.Selected.Route.DistanceKm, a.PayloadTonnes, a.Selected.Route.GradientFactor, a.Selected.FuelIndex, a.Selected.Saving)
	fmt.Println("The signed envelope exposes audience, use, expiry, route suggestion, and compact fuel indices, but not raw payload, GPS trace, driver behavior, or raw telemetry.")
	fmt.Println("This follows the insight pattern: ship the decision, not the data.")
	fmt.Println()
}
