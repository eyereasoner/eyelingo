package main

import (
	"math"
)

func checkBayesDiagnosis(ctx *Context) []Check {
	d := ctx.M()
	priors := map[string]float64{}
	for _, di := range maps(d["Diseases"]) {
		priors[str(di["Name"])] = num(di["Prior"])
	}
	cond := asMap(d["ProbGiven"])
	evidence := maps(d["Evidence"])
	raw := map[string]float64{}
	for disease, prior := range priors {
		lik := prior
		probs := asMap(cond[disease])
		for _, item := range evidence {
			p := num(probs[str(item["Symptom"])])
			if boolean(item["Present"]) {
				lik *= p
			} else {
				lik *= 1 - p
			}
		}
		raw[disease] = lik
	}
	total := 0.0
	for _, v := range raw {
		total += v
	}
	post := map[string]float64{}
	for d, v := range raw {
		post[d] = v / total
	}
	reported := parseBayesDistribution(ctx.Answer)
	best := ""
	bestv := -1.0
	for k, v := range post {
		if v > bestv {
			bestv = v
			best = k
		}
	}
	reportedTotal := math.NaN()
	if m := reFind(ctx.Reason, `Evidence total \(normalizing constant\) = ([0-9]+(?:\.[0-9]+)?)`); m != nil {
		reportedTotal = parseFloat(m[0])
	}
	priorsOK := true
	psum := 0.0
	for _, p := range priors {
		priorsOK = priorsOK && p >= 0 && p <= 1
		psum += p
	}
	condOK := true
	for _, vv := range cond {
		for _, p := range asMap(vv) {
			condOK = condOK && num(p) >= 0 && num(p) <= 1
		}
	}
	evidenceOK := true
	for disease := range priors {
		probs := asMap(cond[disease])
		for _, item := range evidence {
			_, ok := probs[str(item["Symptom"])]
			evidenceOK = evidenceOK && ok
		}
	}
	reportOK := len(reported) == len(priors)
	rawOK, postOK := true, true
	rsum := 0.0
	for disease, r := range reported {
		rawOK = rawOK && math.Abs(r.raw-raw[disease]) < 5e-8
		postOK = postOK && math.Abs(r.post-post[disease]) < 5e-6
		rsum += r.post
	}
	covidRaw := 0.05 * 0.7 * 0.65 * 0.4 * (1 - 0.15) * 0.2
	return []Check{{"all priors are probabilities and the prior mass is less than one", priorsOK && psum < 1}, {"every conditional probability is in [0, 1]", condOK}, {"all evidence symptoms are available for every disease", evidenceOK}, {"the absent Sneezing evidence uses the complement likelihood", close(raw["COVID19"], covidRaw, 1e-15)}, {"the Bayesian normalizing constant is recomputed independently", !math.IsNaN(reportedTotal) && math.Abs(reportedTotal-total) < 5e-9}, {"the reported distribution contains one posterior for each disease", reportOK}, {"each reported unnormalized likelihood matches the Go recomputation", rawOK}, {"each reported posterior matches likelihood divided by evidence total", postOK}, {"the reported posteriors sum to one after rounding", math.Abs(rsum-1.0) < 2e-6}, {"COVID19 is independently selected as the maximum-posterior disease", best == "COVID19" && contains(ctx.Answer, "most likely disease is COVID19")}}
}
