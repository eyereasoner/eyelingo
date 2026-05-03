package main

func checkParcel(ctx *Context) []Check {
	d := ctx.M()
	computed := map[string]string{}
	details := map[string]map[string]bool{}
	for _, r := range maps(d["Requests"]) {
		dec, ch := parcelEval(d, r)
		computed[str(r["Key"])] = dec
		details[str(r["Key"])] = ch
	}
	reported := map[string]string{}
	for _, m := range reAll(ctx.Answer, `(?m)^\s+([a-z-]+)\s+: (PERMIT|DENY)`) {
		reported[m[1]] = m[2]
	}
	pickupTrue := 0
	for _, v := range details["pickup"] {
		if v {
			pickupTrue++
		}
	}
	guard := 0
	for k, dec := range computed {
		if k != "pickup" && dec == "DENY" {
			guard++
		}
	}
	return []Check{{"the source pickup request satisfies all ten authorization conditions", computed["pickup"] == "PERMIT" && pickupTrue == 10}, {"all reported request decisions match independent policy evaluation", mapStrEq(reported, computed)}, {"billing access is denied by the privacy guardrail", computed["billing"] == "DENY" && !details["billing"]["C9"]}, {"redirect is denied by the parcel-redirection guardrail", computed["redirect"] == "DENY" && !details["redirect"]["C10"]}, {"wrong-person use is denied because requester must match the delegate", computed["wrong-person"] == "DENY" && !details["wrong-person"]["C1"]}, {"wrong-locker use is denied because locker must match the token", computed["wrong-locker"] == "DENY" && !details["wrong-locker"]["C3"]}, {"single-use reuse is denied after the token is already consumed", computed["reuse"] == "DENY" && !details["reuse"]["C6"]}, {"guardrail denials recompute to five out of five", guard == 5 && contains(ctx.Answer, "guardrail denials : 5/5")}, {"the release text matches parcel owner, delegate, parcel, locker, and site", contains(ctx.Answer, "Noah may collect parcel123 for Maya from locker B17 at Station West")}}
}
