package main

import (
	"math"
)

func checkBayesTherapy(ctx *Context) []Check {
	d := ctx.M()
	diseases, raw, post, total := recomputeBayesPosteriors(d)
	therapy := recomputeTherapies(d, diseases, post)
	reportedPost := parseBayesDistribution(ctx.Answer)
	reportedTher := parseTherapies(ctx.Answer)
	best, bestUtil := parseRecommendation(ctx.Answer)
	reportedTotal := math.NaN()
	if m := reFind(ctx.Reason, `Evidence total \(normalizing constant\) = ([0-9]+(?:\.[0-9]+)?)`); m != nil {
		reportedTotal = parseFloat(m[0])
	}
	expectedBest := ""
	expUtil := -1e99
	for name, t := range therapy {
		if t.utility > expUtil {
			expUtil = t.utility
			expectedBest = name
		}
	}
	priorSum := 0.0
	priorOK := true
	for _, di := range maps(d["Diseases"]) {
		p := num(di["Prior"])
		priorOK = priorOK && p >= 0 && p <= 1
		priorSum += p
	}
	condOK := true
	for _, vv := range asMap(d["ProbGiven"]) {
		for _, p := range asMap(vv) {
			condOK = condOK && num(p) >= 0 && num(p) <= 1
		}
	}
	postRawOK, postOK := true, true
	rsum := 0.0
	for dise, r := range reportedPost {
		postRawOK = postRawOK && math.Abs(r.raw-raw[dise]) < 5e-8
		postOK = postOK && math.Abs(r.post-post[dise]) < 5e-6
		rsum += r.post
	}
	therOK := len(reportedTher) == len(therapy)
	successOK, utilOK := true, true
	for name, r := range reportedTher {
		t := therapy[name]
		successOK = successOK && math.Abs(r.success-t.success) < 5e-6 && math.Abs(r.adverse-t.adverse) < 5e-8
		utilOK = utilOK && math.Abs(r.utility-t.utility) < 5e-6
	}
	return []Check{{"all disease priors are probabilities and the model is deliberately incomplete", priorOK && priorSum < 1}, {"every conditional probability is in [0, 1]", condOK}, {"the evidence normalizing constant is recomputed independently", !math.IsNaN(reportedTotal) && math.Abs(reportedTotal-total) < 5e-9}, {"reported posteriors and unnormalized likelihoods match recomputation", len(reportedPost) == len(post) && postRawOK && postOK}, {"reported posteriors sum to one after rounding", math.Abs(rsum-1.0) < 2e-6}, {"all therapy rows are reported", therOK}, {"expected success and adverse-risk inputs match independent utility computation", successOK}, {"each reported utility is benefit-weighted success minus harm-weighted adverse risk", utilOK}, {"the recommended therapy is independently selected as maximum utility", best == expectedBest && math.Abs(bestUtil-expUtil) < 5e-6 && contains(ctx.Answer, "recommended therapy is "+expectedBest)}}
}
