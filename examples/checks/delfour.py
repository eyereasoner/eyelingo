# Independent Python checks for the delfour example.
from __future__ import annotations

import hashlib
import json
import re

from .common import run_checks


def canonical_payload(data: dict) -> str:
    insight = data["Insight"]
    policy = data["Policy"]
    payload = {
        "insight": {
            "createdAt": insight["CreatedAt"],
            "expiresAt": insight["ExpiresAt"],
            "id": insight["ID"],
            "metric": insight["Metric"],
            "retailer": insight["Retailer"],
            "scopeDevice": insight["ScopeDevice"],
            "scopeEvent": insight["ScopeEvent"],
            "suggestionPolicy": insight["SuggestionPolicy"],
            "threshold": float(insight["ThresholdG"]),
            "type": "ins:Insight",
        },
        "policy": {
            "duty": {
                "action": policy["Duty"]["Action"],
                "constraint": {
                    "leftOperand": policy["Duty"]["Constraint"]["LeftOperand"],
                    "operator": policy["Duty"]["Constraint"]["Operator"],
                    "rightOperand": policy["Duty"]["Constraint"]["RightOperand"],
                },
            },
            "permission": {
                "action": policy["Permission"]["Action"],
                "constraint": {
                    "leftOperand": policy["Permission"]["Constraint"]["LeftOperand"],
                    "operator": policy["Permission"]["Constraint"]["Operator"],
                    "rightOperand": policy["Permission"]["Constraint"]["RightOperand"],
                },
                "target": policy["Permission"]["Target"],
            },
            "profile": policy["Profile"],
            "prohibition": {
                "action": policy["Prohibition"]["Action"],
                "constraint": {
                    "leftOperand": policy["Prohibition"]["Constraint"]["LeftOperand"],
                    "operator": policy["Prohibition"]["Constraint"]["Operator"],
                    "rightOperand": policy["Prohibition"]["Constraint"]["RightOperand"],
                },
                "target": policy["Prohibition"]["Target"],
            },
            "type": "odrl:Policy",
        },
    }
    # Match the fixture's lexical JSON: compact separators with the same key order.
    return json.dumps(payload, separators=(",", ":"), ensure_ascii=False)


def answer_field(answer: str, label: str) -> str | None:
    match = re.search(rf"{re.escape(label)}\s*:\s*(.+)", answer)
    return match.group(1).strip() if match else None


def run(ctx):
    data = ctx.load_input()
    case = data["Case"]
    insight = data["Insight"]
    policy = data["Policy"]
    signature = data["Signature"]
    products = {product["ID"]: product for product in data["Products"]}
    scanned = products[data["Scan"]["ScannedProductID"]]
    lower = [product for product in products.values() if product["SugarTenths"] < scanned["SugarTenths"]]
    suggested = min(lower, key=lambda product: (product["SugarTenths"], product["Name"])) if lower else None

    payload = canonical_payload(data)
    escaped_payload = payload.replace('"', r'\"')
    payload_hash = hashlib.sha256(escaped_payload.encode("utf-8")).hexdigest()
    serialized = insight["SerializedLowercase"].lower()

    permission = policy["Permission"]
    prohibition = policy["Prohibition"]
    duty = policy["Duty"]
    allowed = (
        permission["Action"] == case["RequestAction"]
        and permission["Target"] == insight["ID"]
        and permission["Constraint"]["LeftOperand"] == "odrl:purpose"
        and permission["Constraint"]["Operator"] == "odrl:eq"
        and permission["Constraint"]["RightOperand"] == case["RequestPurpose"]
        and case["ScannerAuthAt"] <= insight["ExpiresAt"]
    )
    duty_consistent = (
        duty["Action"] == "odrl:delete"
        and duty["Constraint"]["LeftOperand"] == "odrl:dateTime"
        and duty["Constraint"]["Operator"] == "odrl:eq"
        and duty["Constraint"]["RightOperand"] == insight["ExpiresAt"]
        and case["ScannerDutyAt"] <= insight["ExpiresAt"]
    )
    marketing_blocked = (
        prohibition["Action"] == "odrl:distribute"
        and prohibition["Target"] == insight["ID"]
        and prohibition["Constraint"]["RightOperand"] == "marketing"
    )

    checks = [
        ("the scanner request satisfies the ODRL permission independently from Go", allowed and "decision : Allowed" in ctx.answer),
        ("the policy prohibition independently blocks marketing distribution", marketing_blocked),
        ("the delete duty is tied to the insight expiry and is still before expiry", duty_consistent),
        ("the minimized serialized insight omits the sensitive medical condition", "diabetes" not in serialized and "medical" not in serialized),
        ("the scanned product exceeds the sugar threshold, so a banner is warranted", scanned["SugarTenths"] > insight["ThresholdTenths"] and "banner headline" in ctx.reason),
        ("the selected alternative is the lowest-sugar lower-sugar catalog item", suggested is not None and suggested["Name"] == answer_field(ctx.answer, "suggested alternative") and suggested["SugarTenths"] == 30),
        ("the suggested alternative reduces sugar grams per serving by 9 g", suggested is not None and scanned["SugarPerServing"] - suggested["SugarPerServing"] == 9),
        ("the payload SHA-256 is recomputed from the canonical escaped JSON", payload_hash == signature["PayloadHashSHA256"]),
        ("the signature metadata is structurally valid for the trusted precomputed HMAC mode", signature["Alg"] == "HMAC-SHA256" and signature["HMACVerificationMode"] == "trusted-precomputed-input" and len(signature["SignatureHMAC"]) == 64),
        ("the reported bus and audit counts match the independent input fixture", f"bus files written : {case['FilesWritten']}" in ctx.reason and f"audit entries : {case['AuditEntries']}" in ctx.reason),
    ]
    return run_checks(checks)
