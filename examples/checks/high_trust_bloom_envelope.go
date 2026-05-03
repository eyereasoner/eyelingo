package main

import (
	"math"
)

func checkBloom(ctx *Context) []Check {
	d := ctx.M()
	a := asMap(d["Artifact"])
	lam := num(a["HashFunctions"]) * num(a["CanonicalTripleCount"]) / num(a["BloomBits"])
	fpLow := math.Pow(1-num(a["ExpMinusLambdaUpper"]), num(a["HashFunctions"]))
	fpHigh := math.Pow(1-num(a["ExpMinusLambdaLower"]), num(a["HashFunctions"]))
	extra := fpHigh * num(a["NegativeLookupsPerBatch"])
	exact := math.Exp(-lam)
	fp1, fp2 := math.NaN(), math.NaN()
	if m := reFind(ctx.Answer, `false-positive envelope : ([0-9.]+) \.\. ([0-9.]+)`); m != nil {
		fp1 = parseFloat(m[0])
		fp2 = parseFloat(m[1])
	}
	pos := true
	for _, k := range []string{"CanonicalTripleCount", "BloomBits", "HashFunctions", "NegativeLookupsPerBatch"} {
		pos = pos && num(a[k]) > 0
	}
	return []Check{{"numeric Bloom parameters are positive", pos}, {"canonical graph and SPO index triple counts agree", num(a["CanonicalTripleCount"]) == num(a["SPOIndexTripleCount"])}, {"lambda recomputes as k*n/m", close(fieldFloat(ctx.Answer, "lambda"), lam, 1e-12) && close(lam, num(a["CertifiedLambda"]), 1e-12)}, {"the exp(-lambda) decimal certificate brackets the exact value", num(a["ExpMinusLambdaLower"]) <= exact && exact <= num(a["ExpMinusLambdaUpper"])}, {"false-positive envelope is recomputed from the certificate", close(fp1, fpLow, 5e-10) && close(fp2, fpHigh, 5e-10)}, {"false-positive upper bound stays below the policy budget", fpHigh < num(a["FPRateBudget"])}, {"expected extra exact lookups stay below budget", close(fieldFloat(ctx.Answer, "expected extra exact lookups upper"), extra, 5e-3) && extra < num(a["ExtraExactLookupsBudget"])}, {"maybe-positive answers must be confirmed against the canonical graph", str(asMap(d["Policies"])["MaybePositivePolicy"]) == "ConfirmAgainstCanonicalGraph" && contains(ctx.Answer, "maybe-positive policy : ConfirmAgainstCanonicalGraph")}, {"definite negatives may return absent without exact lookup", str(asMap(d["Policies"])["DefiniteNegativePolicy"]) == "ReturnAbsent" && contains(ctx.Answer, "definite-negative policy : ReturnAbsent")}, {"the deployment decision matches the recomputed envelope", contains(ctx.Answer, str(asMap(d["Expected"])["Decision"])) && fpHigh < num(a["FPRateBudget"]) && extra < num(a["ExtraExactLookupsBudget"])}}
}
