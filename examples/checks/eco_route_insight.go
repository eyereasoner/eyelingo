package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math"
	"strings"
)

func checkEcoRoute(ctx *Context) []Check {
	d := ctx.M()
	ans := ecoAnswer(ctx.Answer)
	payload := num(asMap(d["shipment"])["payloadKg"]) / 1000
	current := ecoFuel(asMap(d["currentRoute"]), payload)
	policy := asMap(d["policy"])
	curETA := integer(asMap(d["currentRoute"])["etaMinutes"])
	var selected anymap
	selectedFI, selectedSaving := 0.0, 0.0
	selectedDelay := 0
	selectedElig := false
	for _, r := range maps(d["alternativeRoutes"]) {
		fi := ecoFuel(r, payload)
		saving := current - fi
		delay := integer(r["etaMinutes"]) - curETA
		elig := saving > 0 && delay <= integer(policy["maxEtaDelayMinutes"])
		if boolean(policy["requireAlternativeBelowThreshold"]) {
			elig = elig && fi <= num(policy["fuelIndexThreshold"])
		}
		if selected == nil || (!selectedElig && elig) || (elig == selectedElig && saving > selectedSaving) || (elig == selectedElig && saving == selectedSaving && str(r["id"]) < str(selected["id"])) {
			selected = r
			selectedFI = fi
			selectedSaving = saving
			selectedDelay = delay
			selectedElig = elig
		}
	}
	scenario := asMap(d["scenario"])
	signing := asMap(d["signing"])
	exp := ecoExpiry(str(scenario["issuedAt"]), integer(scenario["ttlHours"]))
	issue := current > num(policy["fuelIndexThreshold"]) && selectedElig
	env := ecoEnv{str(scenario["depot"]), str(policy["allowedUse"]), str(scenario["issuedAt"]), exp, str(signing["keyId"]), ecoAssertions{issue, str(selected["id"]), ecoGoNum(current), ecoGoNum(selectedFI), ecoGoNum(selectedSaving), false}}
	b, _ := json.Marshal(env)
	canonical := string(b)
	digest := sha256Hex(canonical)
	mac := hmac.New(sha256.New, []byte(str(signing["secret"])))
	mac.Write([]byte(canonical))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	forbiddenOK := true
	for _, term := range sarr(asMap(d["dataMinimization"])["forbiddenEnvelopeTerms"]) {
		forbiddenOK = forbiddenOK && !contains(canonical, term)
	}
	reasonOK := contains(ctx.Reason, "42") && contains(ctx.Reason, "120.75") && contains(ctx.Reason, "alt-low-fuel") && contains(ctx.Reason, "99.75") && contains(ctx.Reason, "ship the decision, not the data")
	return []Check{{"the current fuel index is recomputed from distance, payload tonnes, and gradient", math.Round(current*100)/100 == 120.75 && ans["current fuel index"] == "120.75"}, {"the policy threshold triggers a local eco banner", current > num(policy["fuelIndexThreshold"]) && strings.ToLower(ans["show eco banner"]) == "yes"}, {"the selected alternative is the best eligible lower-fuel route", str(selected["id"]) == "alt-low-fuel" && selectedElig && ans["suggested route"] == "alt-low-fuel"}, {"the alternative fuel index and saving are independently recomputed", math.Round(selectedFI*100)/100 == 99.75 && math.Round(selectedSaving*100)/100 == 21.00 && ans["suggested fuel index"] == "99.75" && ans["estimated saving"] == "21.00"}, {"the selected route stays within the allowed ETA delay", selectedDelay == 2 && selectedDelay <= integer(policy["maxEtaDelayMinutes"])}, {"the envelope audience, allowed use, expiry, and raw-data flag match the policy", ans["audience"] == env.Audience && ans["allowed use"] == env.AllowedUse && ans["expires at"] == env.Expiry && strings.ToLower(ans["raw data exported"]) == "no"}, {"the canonical envelope omits forbidden raw logistics terms", forbiddenOK}, {"the payload digest is SHA-256 over the independently rebuilt envelope", ans["payload digest"] == digest}, {"the signature is the expected base64url HMAC-SHA256 value", ans["signature algorithm"] == str(policy["signatureAlgorithm"]) && ans["signature key"] == str(signing["keyId"]) && ans["signature"] == sig}, {"the Reason text reports the same arithmetic and the insight pattern", reasonOK}}
}
