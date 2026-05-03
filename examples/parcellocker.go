// parcellocker.go
//
// A self-contained Go translation of the Eyeling ParcelLocker delegation
// example.
//
// The source N3 file models a narrow, one-time permission: Maya lets Noah pick
// up one specific parcel from one specific locker for the pickup purpose only.
// The same token must not reveal billing details, redirect the parcel, work for
// another person, work at another locker, or be used after it has been consumed.
//
// This program is intentionally not a generic RDF or N3 reasoner. It translates
// the concrete facts and rules for this example into ordinary Go structs and
//
// Run:
//
//	go run parcellocker.go
//
// The program has no third-party dependencies.
package main

import (
	"see/internal/exampleinput"
	"fmt"
	"sort"
	"strings"
)

const seeExampleName = "parcellocker"

const sourceFile = "parcellocker.n3"

const (
	ActionCollectParcel = "CollectParcel"
	ActionViewBilling   = "ViewBillingDetails"
	ActionRedirect      = "RedirectParcel"

	PurposePickupOnly = "PickupOnly"
	PurposeBilling    = "BillingAccess"
	PurposeRedirect   = "RedirectDelivery"

	StateActive   = "Active"
	StateConsumed = "Consumed"

	ReuseSingleUse = "SingleUse"

	BillingNone = "None"
	RedirectNo  = "No"

	StatusReady = "ReadyForPickup"
)

type Person struct {
	ID   string
	Name string
}

type Parcel struct {
	ID     string
	Owner  string
	Status string
}

type Locker struct {
	ID   string
	Site string
}

type Authorization struct {
	ID              string
	IssuedBy        string
	Delegate        string
	Parcel          string
	Locker          string
	Action          string
	Purpose         string
	State           string
	Reuse           string
	BillingAccess   string
	RedirectAllowed string
}

type Request struct {
	Key       string
	Label     string
	Requester string
	Parcel    string
	Locker    string
	Action    string
	Purpose   string
	UsedOnce  bool
	Expected  string
}

type Dataset struct {
	CaseName      string
	Question      string
	People        map[string]Person
	Parcel        Parcel
	Locker        Locker
	Authorization Authorization
	Requests      []Request
}

type Condition struct {
	Code string
	OK   bool
	Text string
}

type RequestResult struct {
	Request    Request
	Decision   string
	Reason     string
	Conditions []Condition
}

type Check struct {
	Label string
	OK    bool
	Text  string
}

type Analysis struct {
	Results        []RequestResult
	Primary        RequestResult
	GuardrailDeny  int
	ExpectedPass   int
	ConditionsPass int
	Checks         []Check
}

func main() {
	ds := exampleinput.Load(seeExampleName, fixture())
	analysis := derive(ds)
	printAnswer(ds, analysis)
	printReason(ds, analysis)
}

func fixture() Dataset {
	return Dataset{
		CaseName: "parcellocker",
		Question: "May Noah use Maya's one-time pickup token to collect parcel123 from locker B17?",
		People: map[string]Person{
			"maya": {ID: "maya", Name: "Maya"},
			"noah": {ID: "noah", Name: "Noah"},
		},
		Parcel: Parcel{ID: "parcel123", Owner: "maya", Status: StatusReady},
		Locker: Locker{ID: "lockerB17", Site: "Station West"},
		Authorization: Authorization{
			ID:              "pickupToken",
			IssuedBy:        "maya",
			Delegate:        "noah",
			Parcel:          "parcel123",
			Locker:          "lockerB17",
			Action:          ActionCollectParcel,
			Purpose:         PurposePickupOnly,
			State:           StateActive,
			Reuse:           ReuseSingleUse,
			BillingAccess:   BillingNone,
			RedirectAllowed: RedirectNo,
		},
		Requests: []Request{
			{
				Key:       "pickup",
				Label:     "Noah collects the parcel",
				Requester: "noah",
				Parcel:    "parcel123",
				Locker:    "lockerB17",
				Action:    ActionCollectParcel,
				Purpose:   PurposePickupOnly,
				Expected:  "PERMIT",
			},
			{
				Key:       "billing",
				Label:     "Noah views billing details",
				Requester: "noah",
				Parcel:    "parcel123",
				Locker:    "lockerB17",
				Action:    ActionViewBilling,
				Purpose:   PurposeBilling,
				Expected:  "DENY",
			},
			{
				Key:       "redirect",
				Label:     "Noah redirects the parcel",
				Requester: "noah",
				Parcel:    "parcel123",
				Locker:    "lockerB17",
				Action:    ActionRedirect,
				Purpose:   PurposeRedirect,
				Expected:  "DENY",
			},
			{
				Key:       "wrong-person",
				Label:     "Maya uses Noah's pickup token",
				Requester: "maya",
				Parcel:    "parcel123",
				Locker:    "lockerB17",
				Action:    ActionCollectParcel,
				Purpose:   PurposePickupOnly,
				Expected:  "DENY",
			},
			{
				Key:       "wrong-locker",
				Label:     "Noah tries another locker",
				Requester: "noah",
				Parcel:    "parcel123",
				Locker:    "lockerA04",
				Action:    ActionCollectParcel,
				Purpose:   PurposePickupOnly,
				Expected:  "DENY",
			},
			{
				Key:       "reuse",
				Label:     "Noah reuses the token",
				Requester: "noah",
				Parcel:    "parcel123",
				Locker:    "lockerB17",
				Action:    ActionCollectParcel,
				Purpose:   PurposePickupOnly,
				UsedOnce:  true,
				Expected:  "DENY",
			},
		},
	}
}

func derive(ds Dataset) Analysis {
	results := make([]RequestResult, 0, len(ds.Requests))
	primary := RequestResult{}
	guardrailDeny := 0
	expectedPass := 0
	conditionsPass := 0

	for _, request := range ds.Requests {
		result := evaluate(ds, request)
		results = append(results, result)
		if request.Key == "pickup" {
			primary = result
		}
		if request.Key != "pickup" && result.Decision == "DENY" {
			guardrailDeny++
		}
		if result.Decision == request.Expected {
			expectedPass++
		}
		for _, condition := range result.Conditions {
			if condition.OK {
				conditionsPass++
			}
		}
	}

	analysis := Analysis{
		Results:        results,
		Primary:        primary,
		GuardrailDeny:  guardrailDeny,
		ExpectedPass:   expectedPass,
		ConditionsPass: conditionsPass,
	}
	analysis.Checks = buildChecks(ds, analysis)
	return analysis
}

func evaluate(ds Dataset, request Request) RequestResult {
	token := ds.Authorization
	parcel := ds.Parcel

	conditions := []Condition{
		{Code: "C1", OK: request.Requester == token.Delegate, Text: "requester must match the named delegate"},
		{Code: "C2", OK: request.Parcel == token.Parcel, Text: "requested parcel must match the authorized parcel"},
		{Code: "C3", OK: request.Locker == token.Locker, Text: "requested locker must match the authorized locker"},
		{Code: "C4", OK: request.Action == token.Action, Text: "requested action must be parcel collection"},
		{Code: "C5", OK: request.Purpose == token.Purpose, Text: "requested purpose must be pickup only"},
		{Code: "C6", OK: token.State == StateActive && !request.UsedOnce, Text: "authorization must be active and not already consumed"},
		{Code: "C7", OK: token.Reuse == ReuseSingleUse, Text: "authorization must be single-use"},
		{Code: "C8", OK: parcel.Status == StatusReady, Text: "parcel must be ready for pickup"},
		{Code: "C9", OK: token.BillingAccess == BillingNone && request.Action != ActionViewBilling, Text: "billing details must stay hidden"},
		{Code: "C10", OK: token.RedirectAllowed == RedirectNo && request.Action != ActionRedirect, Text: "parcel redirection must stay blocked"},
	}

	result := RequestResult{Request: request, Conditions: conditions}
	failing := failedConditions(conditions)
	if len(failing) == 0 {
		result.Decision = "PERMIT"
		result.Reason = "Permit: the requester, parcel, locker, action, purpose, active state, single-use limit, parcel readiness, and privacy restrictions all match."
		return result
	}

	result.Decision = "DENY"
	result.Reason = "Deny: " + strings.Join(failing, "; ") + "."
	return result
}

func failedConditions(conditions []Condition) []string {
	var out []string
	for _, condition := range conditions {
		if !condition.OK {
			out = append(out, condition.Code+" "+condition.Text)
		}
	}
	return out
}

func buildChecks(ds Dataset, analysis Analysis) []Check {
	byKey := resultsByKey(analysis.Results)
	sourceConditions := byKey["pickup"].Decision == "PERMIT" && allConditionsOK(byKey["pickup"].Conditions)
	guardrails := byKey["billing"].Decision == "DENY" && byKey["redirect"].Decision == "DENY" && byKey["wrong-person"].Decision == "DENY" && byKey["wrong-locker"].Decision == "DENY" && byKey["reuse"].Decision == "DENY"
	expected := analysis.ExpectedPass == len(ds.Requests)
	singleUse := conditionOK(byKey["pickup"], "C7") && !conditionOK(byKey["reuse"], "C6")
	privacy := !conditionOK(byKey["billing"], "C9") && !conditionOK(byKey["redirect"], "C10")

	return []Check{
		{Label: "C1", OK: sourceConditions, Text: "the source pickup request satisfies all ten authorization conditions"},
		{Label: "C2", OK: guardrails, Text: "the same token is denied for billing, redirect, wrong person, wrong locker, and reuse"},
		{Label: "C3", OK: expected, Text: "every request matches its expected PERMIT or DENY result"},
		{Label: "C4", OK: singleUse, Text: "single-use state permits the first pickup but rejects reuse"},
		{Label: "C5", OK: privacy, Text: "billing details stay hidden and parcel redirection remains blocked"},
	}
}

func resultsByKey(results []RequestResult) map[string]RequestResult {
	out := map[string]RequestResult{}
	for _, result := range results {
		out[result.Request.Key] = result
	}
	return out
}

func allConditionsOK(conditions []Condition) bool {
	for _, condition := range conditions {
		if !condition.OK {
			return false
		}
	}
	return true
}

func conditionOK(result RequestResult, code string) bool {
	for _, condition := range result.Conditions {
		if condition.Code == code {
			return condition.OK
		}
	}
	return false
}

func printAnswer(ds Dataset, analysis Analysis) {
	primary := analysis.Primary
	parcel := ds.Parcel
	locker := ds.Locker
	owner := ds.People[parcel.Owner].Name
	delegate := ds.People[ds.Authorization.Delegate].Name

	fmt.Println("# Parcel Locker")

	fmt.Println()

	fmt.Println("## Answer")
	fmt.Println(ds.Question)
	fmt.Printf("decision : %s\n", primary.Decision)
	fmt.Printf("release : %s may collect %s for %s from locker %s at %s\n", delegate, parcel.ID, owner, lockerLabel(locker.ID), locker.Site)
	fmt.Printf("guardrail denials : %d/%d\n", analysis.GuardrailDeny, len(ds.Requests)-1)
	fmt.Println()
	fmt.Println("Request decisions:")
	for _, result := range analysis.Results {
		marker := ""
		if result.Request.Key == "pickup" {
			marker = " (source request)"
		}
		fmt.Printf("  %-13s : %s%s\n", result.Request.Key, result.Decision, marker)
	}
	fmt.Println()
}

func printReason(ds Dataset, analysis Analysis) {
	fmt.Println("## Reason")
	token := ds.Authorization
	fmt.Printf("token : delegate=%s parcel=%s locker=%s action=%s purpose=%s state=%s reuse=%s\n", token.Delegate, token.Parcel, token.Locker, token.Action, token.Purpose, token.State, token.Reuse)
	fmt.Printf("privacy : billingAccess=%s redirectAllowed=%s\n", token.BillingAccess, token.RedirectAllowed)
	fmt.Printf("parcel : %s status=%s\n", ds.Parcel.ID, ds.Parcel.Status)
	fmt.Println()

	for _, result := range analysis.Results {
		request := result.Request
		fmt.Printf("%s\n", request.Label)
		fmt.Printf("  request : requester=%s parcel=%s locker=%s action=%s purpose=%s\n", request.Requester, request.Parcel, request.Locker, request.Action, request.Purpose)
		if request.UsedOnce {
			fmt.Println("  state override : token has already been used once")
		}
		fmt.Printf("  decision : %s\n", result.Decision)
		fmt.Printf("  reason : %s\n", result.Reason)
		fmt.Printf("  passed conditions : %d/%d\n", countPassed(result.Conditions), len(result.Conditions))
		fmt.Println()
	}
}

func lockerLabel(id string) string {
	return strings.TrimPrefix(id, "locker")
}

func personNames(people map[string]Person) string {
	names := make([]string, 0, len(people))
	for _, person := range people {
		names = append(names, person.Name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

func countPassed(conditions []Condition) int {
	n := 0
	for _, condition := range conditions {
		if condition.OK {
			n++
		}
	}
	return n
}

func passedChecks(checks []Check) int {
	n := 0
	for _, check := range checks {
		if check.OK {
			n++
		}
	}
	return n
}

func allChecksOK(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func yesNo(ok bool) string {
	if ok {
		return "yes"
	}
	return "no"
}
