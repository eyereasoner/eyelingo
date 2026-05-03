package main

import (
	"fmt"
	"strings"
)

func checkDelfour(ctx *Context) []Check {
	d := ctx.M()
	caseM := asMap(d["Case"])
	insight := asMap(d["Insight"])
	policy := asMap(d["Policy"])
	sig := asMap(d["Signature"])
	products := map[string]anymap{}
	for _, p := range maps(d["Products"]) {
		products[str(p["ID"])] = p
	}
	scanned := products[str(asMap(d["Scan"])["ScannedProductID"])]
	var suggested anymap
	for _, p := range products {
		if num(p["SugarTenths"]) < num(scanned["SugarTenths"]) {
			if suggested == nil || num(p["SugarTenths"]) < num(suggested["SugarTenths"]) || (num(p["SugarTenths"]) == num(suggested["SugarTenths"]) && str(p["Name"]) < str(suggested["Name"])) {
				suggested = p
			}
		}
	}
	escaped := strings.ReplaceAll(delfourPayload(d), `"`, `\"`)
	hash := sha256Hex(escaped)
	perm := asMap(policy["Permission"])
	pc := asMap(perm["Constraint"])
	pro := asMap(policy["Prohibition"])
	prc := asMap(pro["Constraint"])
	duty := asMap(policy["Duty"])
	dc := asMap(duty["Constraint"])
	allowed := str(perm["Action"]) == str(caseM["RequestAction"]) && str(perm["Target"]) == str(insight["ID"]) && str(pc["LeftOperand"]) == "odrl:purpose" && str(pc["Operator"]) == "odrl:eq" && str(pc["RightOperand"]) == str(caseM["RequestPurpose"]) && str(caseM["ScannerAuthAt"]) <= str(insight["ExpiresAt"])
	dutyOK := str(duty["Action"]) == "odrl:delete" && str(dc["LeftOperand"]) == "odrl:dateTime" && str(dc["Operator"]) == "odrl:eq" && str(dc["RightOperand"]) == str(insight["ExpiresAt"]) && str(caseM["ScannerDutyAt"]) <= str(insight["ExpiresAt"])
	blocked := str(pro["Action"]) == "odrl:distribute" && str(pro["Target"]) == str(insight["ID"]) && str(prc["RightOperand"]) == "marketing"
	serialized := strings.ToLower(str(insight["SerializedLowercase"]))
	return []Check{{"the scanner request satisfies the ODRL permission independently", allowed && contains(ctx.Answer, "decision : Allowed")}, {"the policy prohibition independently blocks marketing distribution", blocked}, {"the delete duty is tied to the insight expiry and is still before expiry", dutyOK}, {"the minimized serialized insight omits the sensitive medical condition", !contains(serialized, "diabetes") && !contains(serialized, "medical")}, {"the scanned product exceeds the sugar threshold, so a banner is warranted", num(scanned["SugarTenths"]) > num(insight["ThresholdTenths"]) && contains(ctx.Reason, "banner headline")}, {"the selected alternative is the lowest-sugar lower-sugar catalog item", suggested != nil && str(suggested["Name"]) == answerField(ctx.Answer, "suggested alternative") && integer(suggested["SugarTenths"]) == 30}, {"the suggested alternative reduces sugar grams per serving by 9 g", suggested != nil && integer(scanned["SugarPerServing"])-integer(suggested["SugarPerServing"]) == 9}, {"the payload SHA-256 is recomputed from the canonical escaped JSON", hash == str(sig["PayloadHashSHA256"])}, {"the signature metadata is structurally valid for the trusted precomputed HMAC mode", str(sig["Alg"]) == "HMAC-SHA256" && str(sig["HMACVerificationMode"]) == "trusted-precomputed-input" && len(str(sig["SignatureHMAC"])) == 64}, {"the reported bus and audit counts match the independent input fixture", contains(ctx.Reason, fmt.Sprintf("bus files written : %d", integer(caseM["FilesWritten"]))) && contains(ctx.Reason, fmt.Sprintf("audit entries : %d", integer(caseM["AuditEntries"])))}}
}
