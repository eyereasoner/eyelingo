package main

import (
	"fmt"
	"math"
	"math/big"
	"math/cmplx"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type bigInt struct{ v *big.Int }

func newBigInt(x int64) *bigInt { return &bigInt{big.NewInt(x)} }

func newBigIntFromString(s string) *bigInt { z := new(big.Int); z.SetString(s, 10); return &bigInt{z} }

func (b *bigInt) String() string { return b.v.String() }

func (b *bigInt) Cmp(c *bigInt) int { return b.v.Cmp(c.v) }

func (b *bigInt) EqInt(n int64) bool { return b.v.Cmp(big.NewInt(n)) == 0 }

func ackermannBinary(x, y int) *bigInt {
	ty := y + 3
	var z *big.Int
	switch {
	case x == 0:
		z = big.NewInt(int64(ty + 1))
	case x == 1:
		z = big.NewInt(int64(ty + 2))
	case x == 2:
		z = big.NewInt(int64(2 * ty))
	case x == 3:
		z = new(big.Int).Lsh(big.NewInt(1), uint(ty))
	case x == 4 && ty == 3:
		z = big.NewInt(16)
	case x == 4 && ty == 4:
		z = big.NewInt(65536)
	case x == 4 && ty == 5:
		z = new(big.Int).Lsh(big.NewInt(1), 65536)
	case x == 5 && ty == 3:
		z = big.NewInt(65536)
	default:
		z = big.NewInt(0)
	}
	z.Sub(z, big.NewInt(3))
	return &bigInt{z}
}

type auroraDecision struct{ decision, uid, why string }

func auroraEvaluate(d anymap) map[string]auroraDecision {
	reqs := map[string]anymap{}
	for id, r := range asMap(d["Requesters"]) {
		reqs[id] = asMap(r)
	}
	subs := map[string]anymap{}
	for id, s := range asMap(d["Subjects"]) {
		subs[id] = asMap(s)
	}
	policies := maps(d["Policies"])
	out := map[string]auroraDecision{}
	for _, sc := range maps(d["Scenarios"]) {
		id := str(sc["Key"])
		requester := reqs[str(sc["RequesterID"])]
		subject := subs[str(sc["SubjectID"])]
		if setOf(sarr(subject["ConsentDeny"]))[str(sc["Purpose"])] {
			out[id] = auroraDecision{"DENY", "", "subject_opted_out"}
			continue
		}
		matchedPerm := ""
		prohibited := ""
		for _, p := range policies {
			if auroraPolicyMatches(p, sc, requester, subject) {
				if str(p["Kind"]) == "prohibition" {
					prohibited = str(p["UID"])
				} else if str(p["Kind"]) == "permission" && matchedPerm == "" {
					matchedPerm = str(p["UID"])
				}
			}
		}
		if prohibited != "" {
			out[id] = auroraDecision{"DENY", prohibited, "prohibition"}
		} else if matchedPerm != "" {
			out[id] = auroraDecision{"PERMIT", matchedPerm, "permission"}
		} else {
			out[id] = auroraDecision{"DENY", "", "no_permission"}
		}
	}
	return out
}

func auroraPolicyMatches(p, sc, req, sub anymap) bool {
	purpose := str(sc["Purpose"])
	cats := setOf(sarr(sc["Categories"]))
	if str(p["Kind"]) == "permission" {
		if len(sliceAny(p["AllowedPurposes"])) > 0 && !setOf(sarr(p["AllowedPurposes"]))[purpose] {
			return false
		}
		if str(p["RequiredRole"]) != "" && str(sc["Role"]) != str(p["RequiredRole"]) {
			return false
		}
		if str(p["RequiredEnvironment"]) != "" && str(sc["Environment"]) != str(p["RequiredEnvironment"]) {
			return false
		}
		if str(p["RequiredTOM"]) != "" && str(sc["TOM"]) != str(p["RequiredTOM"]) {
			return false
		}
		if len(sliceAny(p["AllowAnyCategories"])) > 0 {
			ok := false
			for _, c := range sarr(p["AllowAnyCategories"]) {
				if cats[c] {
					ok = true
				}
			}
			if !ok {
				return false
			}
		}
		if len(sliceAny(p["RequireAllCategories"])) > 0 {
			for _, c := range sarr(p["RequireAllCategories"]) {
				if !cats[c] {
					return false
				}
			}
		}
		if str(p["UID"]) == "urn:policy:primary-care-001" && str(req["LinkedTo"]) != str(sc["SubjectID"]) {
			return false
		}
		if str(p["UID"]) == "urn:policy:research-aurora-diabetes" && !setOf(sarr(sub["ConsentAllow"]))[purpose] {
			return false
		}
		return true
	}
	if str(p["Kind"]) == "prohibition" {
		return setOf(sarr(p["ProhibitedPurposes"]))[purpose]
	}
	return false
}

func auroraParseRows(answer string) map[string]auroraDecision {
	out := map[string]auroraDecision{}
	for _, m := range reAll(answer, `(?m)^\s+([A-Z]) – .*? : (PERMIT|DENY) \((.*?)\)$`) {
		uid := m[3]
		if uid == "no policy matched" {
			uid = ""
		}
		out[m[1]] = auroraDecision{m[2], uid, ""}
	}
	return out
}

func barleyBlockers(world, line anymap) []string {
	b := []string{}
	if !boolean(line["digitalHeredity"]) {
		b = append(b, "digital-heredity")
	}
	if !boolean(line["repair"]) {
		b = append(b, "repair")
	}
	if !boolean(line["dormancyProtection"]) {
		b = append(b, "protected-dormancy")
	}
	if !boolean(line["heritableVariation"]) {
		b = append(b, "heritable-variation")
	}
	return b
}

func sliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type bayesRow struct{ post, raw float64 }

func parseBayesDistribution(ans string) map[string]bayesRow {
	out := map[string]bayesRow{}
	for _, m := range reAll(ans, `(?m)^\s*([A-Za-z0-9_]+)\s+posterior = ([0-9.]+)\s+\(unnormalized = ([0-9.]+)\)`) {
		out[m[1]] = bayesRow{parseFloat(m[2]), parseFloat(m[3])}
	}
	return out
}

func recomputeBayesPosteriors(d anymap) ([]string, map[string]float64, map[string]float64, float64) {
	diseases := []string{}
	priors := map[string]float64{}
	for _, di := range maps(d["Diseases"]) {
		name := str(di["Name"])
		diseases = append(diseases, name)
		priors[name] = num(di["Prior"])
	}
	cond := asMap(d["ProbGiven"])
	raw := map[string]float64{}
	for _, dise := range diseases {
		lik := priors[dise]
		probs := asMap(cond[dise])
		for _, item := range maps(d["Evidence"]) {
			p := num(probs[str(item["Symptom"])])
			if boolean(item["Present"]) {
				lik *= p
			} else {
				lik *= 1 - p
			}
		}
		raw[dise] = lik
	}
	total := 0.0
	for _, v := range raw {
		total += v
	}
	post := map[string]float64{}
	for k, v := range raw {
		post[k] = v / total
	}
	return diseases, raw, post, total
}

type therapyRow struct{ success, adverse, utility float64 }

func recomputeTherapies(d anymap, diseases []string, post map[string]float64) map[string]therapyRow {
	benefit := num(d["BenefitWeight"])
	harm := num(d["HarmWeight"])
	out := map[string]therapyRow{}
	for _, th := range maps(d["Therapies"]) {
		vals := farr(th["SuccessByDisease"])
		exp := 0.0
		for i, dise := range diseases {
			exp += post[dise] * vals[i]
		}
		adv := num(th["Adverse"])
		out[str(th["Name"])] = therapyRow{exp, adv, benefit*exp - harm*adv}
	}
	return out
}

func parseTherapies(ans string) map[string]therapyRow {
	out := map[string]therapyRow{}
	for _, m := range reAll(ans, `(?m)^\s*([A-Za-z0-9_]+)\s+expectedSuccess = ([0-9.]+)\s+adverse = ([0-9.]+)\s+utility = ([0-9.\-]+)`) {
		out[m[1]] = therapyRow{parseFloat(m[2]), parseFloat(m[3]), parseFloat(m[4])}
	}
	return out
}

func parseRecommendation(ans string) (string, float64) {
	if m := reFind(ans, `recommended therapy is ([^(]+) \(utility = ([0-9.\-]+)\)`); m != nil {
		return strings.TrimSpace(m[0]), parseFloat(m[1])
	}
	return "", math.NaN()
}

func roundHalfUp(v float64, places int) float64 {
	scale := math.Pow10(places)
	return math.Floor(v*scale+0.5) / scale
}

func bmiClass(v float64) string {
	if v < 18.5 {
		return "Underweight"
	}
	if v < 25 {
		return "Normal"
	}
	if v < 30 {
		return "Overweight"
	}
	if v < 35 {
		return "Obesity I"
	}
	if v < 40 {
		return "Obesity II"
	}
	return "Obesity III"
}

func parseBMIAnswer(ans string) (float64, string, float64, float64) {
	bmi := math.NaN()
	cat := ""
	lo, hi := math.NaN(), math.NaN()
	if m := reFind(ans, `BMI = ([0-9]+(?:\.[0-9]+)?)`); m != nil {
		bmi = parseFloat(m[0])
	}
	if m := reFind(ans, `Category = ([^\n]+)`); m != nil {
		cat = strings.TrimSpace(m[0])
	}
	if m := reFind(ans, `range is about ([0-9]+(?:\.[0-9]+)?)–([0-9]+(?:\.[0-9]+)?) kg`); m != nil {
		lo = parseFloat(m[0])
		hi = parseFloat(m[1])
	}
	return bmi, cat, lo, hi
}

func calidorActiveNeeds(d anymap) map[string]bool {
	return map[string]bool{"heat_alert": integer(d["CurrentAlertLevel"]) >= integer(d["AlertLevelAtLeast"]), "unsafe_indoor_heat": num(d["CurrentIndoorTempC"]) >= num(d["IndoorTempCAtLeast"]) && integer(d["HoursAtOrAboveThreshold"]) >= integer(d["HoursAtOrAboveThresholdAtLeast"]), "vulnerability_present": len(asSlice(d["VulnerabilityFlags"])) > 0, "energy_constraint": num(d["RemainingPrepaidCreditEur"]) <= num(d["EnergyCreditEurAtMost"])}
}

func chooseCapabilityPackage(pkgs []anymap, maxCost float64, required []string) anymap {
	req := setOf(required)
	candidates := []anymap{}
	for _, p := range pkgs {
		caps := setOf(sarr(p["Capabilities"]))
		ok := true
		for c := range req {
			ok = ok && caps[c]
		}
		if ok && num(p["CostEur"]) <= maxCost {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) == 0 {
		return nil
	}
	sort.Slice(candidates, func(i, j int) bool {
		if num(candidates[i]["CostEur"]) == num(candidates[j]["CostEur"]) {
			return str(candidates[i]["PackageID"]) < str(candidates[j]["PackageID"])
		}
		return num(candidates[i]["CostEur"]) < num(candidates[j]["CostEur"])
	})
	return candidates[0]
}

func parseCalidorAnswer(ans string) map[string]string {
	vals := map[string]string{}
	patterns := map[string]string{"active": `Active need count: (\d+)/`, "pkg": `Recommended package: ([^\n]+)`, "cost": `Package cost: €(\d+)`, "hash": `Payload SHA-256: ([0-9a-f]+)`, "hmac": `Envelope HMAC-SHA-256: ([0-9a-f]+)`}
	for k, p := range patterns {
		if m := reFind(ans, p); m != nil {
			vals[k] = strings.TrimSpace(m[0])
		}
	}
	return vals
}

func toComplex(m anymap) complex128 { return complex(num(m["Re"]), num(m["Im"])) }

func complexClose(a, b complex128, tol float64) bool {
	return math.Abs(real(a)-real(b)) <= tol && math.Abs(imag(a)-imag(b)) <= tol
}

func parseComplexAnswer(ans string) map[string]complex128 {
	out := map[string]complex128{}
	for _, m := range reAll(ans, `(?m)^(C[0-9]+).*?=\s+([-0-9.]+)\s+([+-])\s+([0-9.]+)i\s*$`) {
		im := parseFloat(m[4])
		if m[3] == "-" {
			im = -im
		}
		out[m[1]] = complex(parseFloat(m[2]), im)
	}
	return out
}

func measurement10(pair []float64) (float64, string) {
	a, b := pair[0], pair[1]
	if a < b {
		return math.Sqrt(b - a), "lessThan"
	}
	return a, "notLessThan"
}

func parseControlAnswer(ans string) map[string]float64 {
	patterns := map[string]string{"actuator1": `actuator1 control1 = (-?[0-9.]+)`, "actuator2": `actuator2 control1 = (-?[0-9.]+)`, "input1_m10": `input1 measurement10 = (-?[0-9.]+)`, "disturbance2_m10": `disturbance2 measurement10 = (-?[0-9.]+)`}
	out := map[string]float64{}
	for k, p := range patterns {
		if m := reFind(ans, p); m != nil {
			out[k] = parseFloat(m[0])
		}
	}
	return out
}

func grabInt(text, pat string) int {
	if m := reFind(text, pat); m != nil {
		return parseInt(m[0])
	}
	return -1
}

func fieldInt(text, label string) int {
	if m := reFind(text, regexp.QuoteMeta(label)+`\s*:\s*([0-9]+)`); m != nil {
		return parseInt(m[0])
	}
	return -1
}

func fieldFloat(text, label string) float64 {
	if m := reFind(text, regexp.QuoteMeta(label)+`\s*:\s*([-+0-9.eE]+)`); m != nil {
		return parseFloat(m[0])
	}
	return math.NaN()
}

func doctorClosure(job string, sub anymap) map[string]bool {
	out := map[string]bool{job: true}
	for {
		v, ok := sub[job]
		if !ok {
			break
		}
		job = str(v)
		out[job] = true
	}
	return out
}

func doctorEval(d anymap) map[string]map[string]string {
	sick := str(asMap(d["Person"])["Condition"]) == "Flu"
	doctorJobs := setOf(sarr(d["DoctorCanDoJobs"]))
	policies := asMap(d["ConflictPolicies"])
	sub := asMap(d["SubclassOf"])
	res := map[string]map[string]string{}
	for _, req := range maps(d["Requests"]) {
		raw := []string{}
		fam := doctorClosure(str(req["Job"]), sub)
		if doctorJobs[str(req["Job"])] {
			raw = append(raw, "Permit")
		}
		if sick && fam["Work"] && str(policies["SickWorkDefault"]) == "Deny" {
			raw = append(raw, "Deny")
		}
		sort.Strings(raw)
		status := "None"
		if len(raw) == 1 {
			status = raw[0]
		} else if len(raw) == 2 {
			status = "BothPermitDeny"
		}
		eff := status
		if status == "BothPermitDeny" && str(req["Location"]) == "Home" && str(policies["HomeProgrammingWork"]) == "Permit" {
			eff = "Permit"
		} else if status == "BothPermitDeny" && str(req["Location"]) == "Office" && str(policies["SickOfficeWorkDefault"]) == "Deny" {
			eff = "Deny"
		}
		res[str(req["ID"])] = map[string]string{"raw": strings.Join(raw, "+"), "status": status, "effective": eff}
	}
	return res
}

func parseDoctorRows(ans string) map[string]map[string]string {
	out := map[string]map[string]string{}
	for _, m := range reAll(ans, `(?m)^(Request_[A-Za-z0-9_]+) : raw=([A-Za-z+]+) status=([A-Za-z]+) effective=([A-Za-z]+)`) {
		out[m[1]] = map[string]string{"raw": m[2], "status": m[3], "effective": m[4]}
	}
	return out
}

var droneFields = []string{"Location", "Battery", "Permit"}

func droneStateKey(st anymap) string {
	return str(st["Location"]) + "|" + str(st["Battery"]) + "|" + str(st["Permit"])
}

func droneMatches(p, st anymap) bool {
	for _, f := range droneFields {
		if str(p[f]) != "*" && str(p[f]) != str(st[f]) {
			return false
		}
	}
	return true
}

func droneApply(p, cur anymap) anymap {
	out := anymap{}
	for _, f := range droneFields {
		if str(p[f]) == "*" {
			out[f] = cur[f]
		} else {
			out[f] = p[f]
		}
	}
	return out
}

type dronePlan struct {
	actions               []string
	end                   anymap
	duration              int
	cost, belief, comfort float64
	fuel                  int
}

func droneSearch(d anymap) []dronePlan {
	plans := []dronePlan{}
	start := asMap(d["Start"])
	actions := maps(d["Actions"])
	goal := str(d["GoalLocation"])
	th := asMap(d["Thresholds"])
	var walk func(anymap, int, []string, int, float64, float64, float64, map[string]bool)
	walk = func(state anymap, fuel int, path []string, dur int, cost, belief, comfort float64, seen map[string]bool) {
		if str(state["Location"]) == goal {
			if belief > num(th["MinBelief"]) && cost < num(th["MaxCost"]) {
				cp := append([]string(nil), path...)
				plans = append(plans, dronePlan{cp, state, dur, cost, belief, comfort, fuel})
			}
			return
		}
		if fuel == 0 {
			return
		}
		for _, a := range actions {
			if !droneMatches(asMap(a["From"]), state) {
				continue
			}
			nxt := droneApply(asMap(a["To"]), state)
			key := droneStateKey(nxt)
			if seen[key] {
				continue
			}
			ns := map[string]bool{}
			for k, v := range seen {
				ns[k] = v
			}
			ns[key] = true
			walk(nxt, fuel-1, append(append([]string(nil), path...), str(a["Name"])), dur+integer(a["DurationSec"]), cost+num(a["Cost"]), belief*num(a["Belief"]), comfort*num(a["Comfort"]), ns)
		}
	}
	walk(start, integer(d["Fuel"]), nil, 0, 0, 1, 1, map[string]bool{droneStateKey(start): true})
	sort.Slice(plans, func(i, j int) bool {
		if math.Abs(plans[i].cost-plans[j].cost) > 1e-12 {
			return plans[i].cost < plans[j].cost
		}
		if plans[i].duration != plans[j].duration {
			return plans[i].duration < plans[j].duration
		}
		if math.Abs(plans[i].belief-plans[j].belief) > 1e-12 {
			return plans[i].belief > plans[j].belief
		}
		return strings.Join(plans[i].actions, ",") < strings.Join(plans[j].actions, ",")
	})
	return plans
}

func parseDroneAnswer(ans string) dronePlan {
	plan := dronePlan{}
	if m := reFind(ans, `selected plan : ([^\n]+)`); m != nil {
		for _, p := range strings.Split(m[0], "->") {
			plan.actions = append(plan.actions, strings.TrimSpace(p))
		}
	}
	plan.duration = grabInt(ans, `duration : (\d+) s`)
	plan.cost = fieldFloat(ans, "cost")
	plan.belief = fieldFloat(ans, "belief")
	plan.comfort = fieldFloat(ans, "comfort")
	plan.fuel = grabInt(ans, `surviving plans : (\d+)`)
	return plan
}

func grayN(n int) int { return n ^ (n >> 1) }

func bitCount(n int) int {
	c := 0
	for n > 0 {
		c += n & 1
		n >>= 1
	}
	return c
}

func maxInt(xs []int) int {
	m := xs[0]
	for _, x := range xs {
		if x > m {
			m = x
		}
	}
	return m
}

func kapStep(n int) int {
	ds := []byte(fmt.Sprintf("%04d", n))
	sort.Slice(ds, func(i, j int) bool { return ds[i] < ds[j] })
	asc, _ := strconv.Atoi(string(ds))
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })
	desc, _ := strconv.Atoi(string(ds))
	return desc - asc
}

func kapChain(n, max, target, zero int) []int {
	if n == target || n == zero {
		return []int{n}
	}
	out := []int{}
	cur := n
	for i := 0; i < max; i++ {
		cur = kapStep(cur)
		out = append(out, cur)
		if cur == target || cur == zero {
			return out
		}
	}
	return out
}

func intsEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func allKapChains(m map[int][]int, target, max int) bool {
	for _, ch := range m {
		if ch[len(ch)-1] != target || len(ch) > max+1 {
			return false
		}
	}
	return true
}

func fibPair(n int) (*big.Int, *big.Int) {
	if n == 0 {
		return big.NewInt(0), big.NewInt(1)
	}
	a, b := fibPair(n / 2)
	twoB := new(big.Int).Mul(b, big.NewInt(2))
	twoB.Sub(twoB, a)
	c := new(big.Int).Mul(a, twoB)
	d := new(big.Int).Add(new(big.Int).Mul(a, a), new(big.Int).Mul(b, b))
	if n%2 == 0 {
		return c, d
	}
	return d, new(big.Int).Add(c, d)
}

func fibBig(n int) *big.Int { a, _ := fibPair(n); return a }

func dftReal(samples []float64) []complex128 {
	n := len(samples)
	out := make([]complex128, n)
	for k := 0; k < n; k++ {
		sum := complex(0, 0)
		for t, x := range samples {
			angle := -2 * math.Pi * float64(k*t) / float64(n)
			sum += complex(x, 0) * cmplx.Exp(complex(0, angle))
		}
		out[k] = sum
	}
	return out
}

func allDominantMag(mags []float64, dom []int, target, tol float64) bool {
	for _, k := range dom {
		if !close(mags[k], target, tol) {
			return false
		}
	}
	return true
}

func fftWave(spec anymap, n int) []float64 {
	kind := str(spec["kind"])
	out := make([]float64, n)
	switch kind {
	case "alternating":
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				out[i] = 1
			} else {
				out[i] = -1
			}
		}
	case "constant":
		for i := range out {
			out[i] = 1
		}
	case "cosine":
		b := integer(spec["bin"])
		for i := 0; i < n; i++ {
			out[i] = math.Cos(2 * math.Pi * float64(b*i) / float64(n))
		}
	case "sine":
		b := integer(spec["bin"])
		for i := 0; i < n; i++ {
			out[i] = math.Sin(2 * math.Pi * float64(b*i) / float64(n))
		}
	case "impulse":
		out[0] = 1
	case "ramp":
		for i := 0; i < n; i++ {
			out[i] = float64(i)
		}
	}
	return out
}

func dominantBins(spec []complex128, tol float64) []int {
	peak := 0.0
	mags := make([]float64, len(spec))
	for i, z := range spec {
		mags[i] = cmplx.Abs(z)
		if mags[i] > peak {
			peak = mags[i]
		}
	}
	out := []int{}
	for i, m := range mags {
		if math.Abs(m-peak) <= tol {
			out = append(out, i)
		}
	}
	return out
}

func parseFFT32Rows(ans string) map[string]string {
	rows := map[string]string{}
	for _, line := range strings.Split(ans, "\n") {
		if strings.Contains(line, " : ") && (strings.Contains(line, "k=") || strings.Contains(line, "all 32 bins")) {
			parts := strings.SplitN(line, " : ", 2)
			rows[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return rows
}

func primeInt(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for d := 3; d*d <= n; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func countQueens(n int) int {
	total := 0
	var search func(int, int, int, int)
	search = func(row, cols, d1, d2 int) {
		if row == n {
			total++
			return
		}
		available := ((1 << n) - 1) & ^(cols | d1 | d2)
		for available != 0 {
			bit := available & -available
			available -= bit
			search(row+1, cols|bit, (d1|bit)<<1, (d2|bit)>>1)
		}
	}
	search(0, 0, 0, 0)
	return total
}

func parseSudokuGrid(ans string) [][]int {
	idx := strings.Index(ans, "Completed grid")
	if idx < 0 {
		return nil
	}
	out := [][]int{}
	for _, line := range strings.Split(ans[idx+len("Completed grid"):], "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ok := true
		row := []int{}
		for _, ch := range line {
			if ch >= '1' && ch <= '9' {
				row = append(row, int(ch-'0'))
			} else if ch != ' ' && ch != '.' && ch != '|' {
				ok = false
			}
		}
		if !ok {
			break
		}
		if len(row) == 9 {
			out = append(out, row)
		}
		if len(out) == 9 {
			break
		}
	}
	return out
}

func digitSet(xs []int) bool {
	if len(xs) != 9 {
		return false
	}
	seen := map[int]bool{}
	for _, x := range xs {
		seen[x] = true
	}
	for i := 1; i <= 9; i++ {
		if !seen[i] {
			return false
		}
	}
	return true
}

func sudokuLegal(g [][]int, r, c, v int) bool {
	for cc := 0; cc < 9; cc++ {
		if cc != c && g[r][cc] == v {
			return false
		}
	}
	for rr := 0; rr < 9; rr++ {
		if rr != r && g[rr][c] == v {
			return false
		}
	}
	br := r / 3 * 3
	bc := c / 3 * 3
	for rr := br; rr < br+3; rr++ {
		for cc := bc; cc < bc+3; cc++ {
			if (rr != r || cc != c) && g[rr][cc] == v {
				return false
			}
		}
	}
	return true
}

func gridEq(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !intsEq(a[i], b[i]) {
			return false
		}
	}
	return true
}

func allRowsLen(g [][]int, n int) bool {
	for _, r := range g {
		if len(r) != n {
			return false
		}
	}
	return true
}

func windPower(v, cut, rated, cutout, rp float64) (string, float64) {
	if v < cut || v >= cutout {
		return "stopped", 0
	}
	if v >= rated {
		return "rated", rp
	}
	return "partial", rp * (math.Pow(v, 3) - math.Pow(cut, 3)) / (math.Pow(rated, 3) - math.Pow(cut, 3))
}

func parseRFC3339ish(s string) int64 { return parseTime(strings.ReplaceAll(s, "Z", "+00:00")).Unix() }

func photoEff(c anymap, rc string) bool {
	return str(c["excitonCoupling"]) == "Strong" && str(c["stateDelocalization"]) == "Present" && str(c["vibronicBridge"]) == "Tuned" && str(c["energyLandscape"]) == "DownhillToReactionCenter" && (str(c["dephasing"]) == "Moderate" || str(c["dephasing"]) == "Low") && (str(c["electronicCoherenceLifetime"]) == "Short" || str(c["electronicCoherenceLifetime"]) == "Long") && str(c["connectedTo"]) == rc
}

func parcelEval(d, req anymap) (string, map[string]bool) {
	auth := asMap(d["Authorization"])
	parcel := asMap(d["Parcel"])
	checks := map[string]bool{"C1": str(req["Requester"]) == str(auth["Delegate"]), "C2": str(req["Parcel"]) == str(auth["Parcel"]), "C3": str(req["Locker"]) == str(auth["Locker"]), "C4": str(req["Action"]) == str(auth["Action"]), "C5": str(req["Purpose"]) == str(auth["Purpose"]), "C6": str(auth["State"]) == "Active" && !boolean(req["UsedOnce"]), "C7": str(auth["Reuse"]) == "SingleUse", "C8": str(parcel["Status"]) == "ReadyForPickup", "C9": !(str(req["Action"]) == "ViewBillingDetails" || str(req["Purpose"]) == "BillingAccess") || str(auth["BillingAccess"]) != "None", "C10": !(str(req["Action"]) == "RedirectParcel" || str(req["Purpose"]) == "RedirectDelivery") || str(auth["RedirectAllowed"]) != "No"}
	ok := true
	for _, v := range checks {
		ok = ok && v
	}
	if ok {
		return "PERMIT", checks
	}
	return "DENY", checks
}

func mapStrEq(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func pathDiscover(d anymap) ([][]string, map[string][]string, int, int) {
	graph := map[string][]string{}
	labels := asMap(d["Labels"])
	for _, e := range maps(d["Edges"]) {
		graph[str(e["From"])] = append(graph[str(e["From"])], str(e["To"]))
	}
	for node := range graph {
		sort.Slice(graph[node], func(i, j int) bool { return str(labels[graph[node][i]]) < str(labels[graph[node][j]]) })
	}
	routes := [][]string{}
	calls, edgeTests := 0, 0
	maxEdges := 3
	var dfs func(string, []string)
	dfs = func(node string, path []string) {
		calls++
		if node == str(d["DestinationID"]) {
			routes = append(routes, append([]string(nil), path...))
			return
		}
		if len(path)-1 == maxEdges {
			return
		}
		for _, nxt := range graph[node] {
			edgeTests++
			seen := false
			for _, p := range path {
				if p == nxt {
					seen = true
				}
			}
			if seen {
				continue
			}
			dfs(nxt, append(append([]string(nil), path...), nxt))
		}
	}
	dfs(str(d["SourceID"]), []string{str(d["SourceID"])})
	sort.Slice(routes, func(i, j int) bool {
		for k := 0; k < len(routes[i]) && k < len(routes[j]); k++ {
			ai, aj := str(labels[routes[i][k]]), str(labels[routes[j][k]])
			if ai != aj {
				return ai < aj
			}
		}
		return len(routes[i]) < len(routes[j])
	})
	return routes, graph, calls, edgeTests
}

func allPathLen(routes [][]string, n int) bool {
	for _, r := range routes {
		if len(r) != n {
			return false
		}
	}
	return true
}

func routesEq(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !sliceEq(a[i], b[i]) {
			return false
		}
	}
	return true
}

func noShort(routes [][]string) bool {
	for _, r := range routes {
		if len(r) == 2 || len(r) == 3 {
			return false
		}
	}
	return true
}

func countGraphEdges(g map[string][]string) int {
	n := 0
	for _, v := range g {
		n += len(v)
	}
	return n
}

func parseRiskPath(ans string) ([]string, float64, float64, float64, int) {
	path := []string{}
	if m := reFind(ans, `selected path\s*:\s*(.+)`); m != nil {
		for _, p := range strings.Split(m[0], "->") {
			path = append(path, strings.TrimSpace(p))
		}
	}
	return path, fieldFloat(ans, "raw cost"), fieldFloat(ans, "risk sum"), fieldFloat(ans, "risk-adjusted score"), grabInt(ans, `edges in selected path\s*:\s*(\d+)`)
}

func edgeRiskScore(e anymap, w float64) float64 { return num(e["cost"]) + w*num(e["risk"]) }

func riskPathTotals(edges []anymap, path []string, w float64) (float64, float64, float64, bool) {
	raw, risk, score := 0.0, 0.0, 0.0
	by := map[string]anymap{}
	for _, e := range edges {
		by[str(e["from"])+"|"+str(e["to"])] = e
	}
	for i := 0; i+1 < len(path); i++ {
		e, ok := by[path[i]+"|"+path[i+1]]
		if !ok {
			return 0, 0, 0, false
		}
		raw += num(e["cost"])
		risk += num(e["risk"])
		score += edgeRiskScore(e, w)
	}
	return raw, risk, score, true
}

func riskAllPaths(d anymap) [][]string {
	graph := map[string][]string{}
	for _, e := range maps(d["edges"]) {
		graph[str(e["from"])] = append(graph[str(e["from"])], str(e["to"]))
	}
	out := [][]string{}
	var dfs func(string, []string)
	dfs = func(node string, path []string) {
		if node == str(d["goal"]) {
			out = append(out, append([]string(nil), path...))
			return
		}
		for _, n := range graph[node] {
			seen := false
			for _, p := range path {
				if p == n {
					seen = true
				}
			}
			if !seen {
				dfs(n, append(append([]string(nil), path...), n))
			}
		}
	}
	dfs(str(d["start"]), []string{str(d["start"])})
	return out
}

var philos = []string{"P1", "P2", "P3", "P4", "P5"}

var forks = []string{"F12", "F23", "F34", "F45", "F51"}

var leftFork = map[string]string{"P1": "F51", "P2": "F12", "P3": "F23", "P4": "F34", "P5": "F45"}

var rightFork = map[string]string{"P1": "F12", "P2": "F23", "P3": "F34", "P4": "F45", "P5": "F51"}

type forkState struct{ holder, clean string }

type mealTrace struct {
	meals     []string
	forksUsed []string
	transfers [][3]string
	requests  int
}

func diningDerive(schedule []any) ([]mealTrace, map[string]int, map[string]forkState, int, int) {
	state := map[string]forkState{"F12": {"P1", "Dirty"}, "F23": {"P2", "Dirty"}, "F34": {"P3", "Dirty"}, "F45": {"P4", "Dirty"}, "F51": {"P1", "Dirty"}}
	counts := map[string]int{}
	traces := []mealTrace{}
	reqTotal, transTotal := 0, 0
	order := map[string]int{}
	for i, p := range philos {
		order[p] = i
	}
	for _, rr := range schedule {
		r := asMap(rr)
		hungry := sarr(r["Hungry"])
		sort.Slice(hungry, func(i, j int) bool { return order[hungry[i]] < order[hungry[j]] })
		type req struct {
			p, holder, fork, side string
			dirty                 bool
		}
		reqs := []req{}
		for _, p := range hungry {
			for _, fork := range []string{leftFork[p], rightFork[p]} {
				st := state[fork]
				if st.holder != p {
					reqs = append(reqs, req{p, st.holder, fork, "", st.clean == "Dirty"})
				}
			}
		}
		transfers := [][3]string{}
		for _, rq := range reqs {
			if rq.dirty {
				transfers = append(transfers, [3]string{rq.holder, rq.p, rq.fork})
				state[rq.fork] = forkState{rq.p, "Clean"}
			}
		}
		meals := []string{}
		used := []string{}
		for _, p := range hungry {
			lf, rf := leftFork[p], rightFork[p]
			if state[lf].holder == p && state[rf].holder == p {
				counts[p]++
				meals = append(meals, p)
				used = append(used, lf, rf)
			}
		}
		for _, f := range forks {
			st := state[f]
			st.clean = "Dirty"
			state[f] = st
		}
		traces = append(traces, mealTrace{meals, used, transfers, len(reqs)})
		reqTotal += len(reqs)
		transTotal += len(transfers)
	}
	return traces, counts, state, reqTotal, transTotal
}

func parseDiningRows(ans string) [][]string {
	rows := [][]string{}
	for _, m := range reAll(ans, `round (\d+) cycle \d+ : ([^\n]+)`) {
		ps := []string{}
		for _, x := range reAll(m[2], `(P\d)#\d+ uses`) {
			ps = append(ps, x[1])
		}
		rows = append(rows, ps)
	}
	return rows
}

func ftaPrime(n int) bool { return primeInt(n) }

func ftaFactor(n int) []int {
	out := []int{}
	for n%2 == 0 {
		out = append(out, 2)
		n /= 2
	}
	for d := 3; d*d <= n; d += 2 {
		for n%d == 0 {
			out = append(out, d)
			n /= d
		}
	}
	if n > 1 {
		out = append(out, n)
	}
	return out
}

func ftaFormat(fs []int) string {
	counts := map[int]int{}
	for _, p := range fs {
		counts[p]++
	}
	keys := []int{}
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	parts := []string{}
	for _, p := range keys {
		if counts[p] == 1 {
			parts = append(parts, fmt.Sprint(p))
		} else {
			parts = append(parts, fmt.Sprintf("%d^%d", p, counts[p]))
		}
	}
	return strings.Join(parts, " * ")
}

func knapEval(genome string, items []anymap, cap int) map[string]any {
	w, v := 0, 0
	selected := []string{}
	for i, ch := range genome {
		if ch == '1' {
			w += integer(items[i]["Weight"])
			v += integer(items[i]["Value"])
			selected = append(selected, str(items[i]["Name"]))
		}
	}
	fit := 1000000 - v
	if w > cap {
		fit = 2000000 + (w - cap)
	}
	return map[string]any{"genome": genome, "weight": w, "value": v, "fitness": fit, "items": selected}
}

func mutants(genome string) []string {
	out := []string{}
	for i, ch := range genome {
		b := '1'
		if ch == '1' {
			b = '0'
		}
		out = append(out, genome[:i]+string(b)+genome[i+1:])
	}
	return out
}

func knapBest(genomes []string, items []anymap, cap int) map[string]any {
	var best map[string]any
	for _, g := range genomes {
		e := knapEval(g, items, cap)
		if best == nil || integer(e["fitness"]) < integer(best["fitness"]) || (integer(e["fitness"]) == integer(best["fitness"]) && str(e["genome"]) < str(best["genome"])) {
			best = e
		}
	}
	return best
}

func gpsPaths(edges []anymap, start, goal string) [][]anymap {
	by := map[string][]anymap{}
	for _, e := range edges {
		by[str(e["From"])] = append(by[str(e["From"])], e)
	}
	out := [][]anymap{}
	var dfs func(string, []anymap, map[string]bool)
	dfs = func(node string, path []anymap, seen map[string]bool) {
		if node == goal {
			out = append(out, append([]anymap(nil), path...))
			return
		}
		for _, e := range by[node] {
			to := str(e["To"])
			if seen[to] {
				continue
			}
			ns := map[string]bool{}
			for k, v := range seen {
				ns[k] = v
			}
			ns[to] = true
			dfs(to, append(append([]anymap(nil), path...), e), ns)
		}
	}
	dfs(start, nil, map[string]bool{start: true})
	return out
}

func gpsMetrics(path []anymap) (float64, float64, float64, float64) {
	dur, cost, belief, comfort := 0.0, 0.0, 1.0, 1.0
	for _, e := range path {
		dur += num(e["Duration"])
		cost += num(e["Cost"])
		belief *= num(e["Belief"])
		comfort *= num(e["Comfort"])
	}
	return dur, cost, belief, comfort
}

func gpsLabel(path []anymap) string {
	if len(path) == 0 {
		return ""
	}
	nodes := []string{str(path[0]["From"])}
	for _, e := range path {
		nodes = append(nodes, str(e["To"]))
	}
	return strings.Join(nodes, " → ")
}

var evFields = []string{"At", "Battery", "Pass"}

func evKey(st anymap) string { return str(st["At"]) + "|" + str(st["Battery"]) + "|" + str(st["Pass"]) }

func evMatches(pattern, state anymap) bool {
	for _, f := range evFields {
		if str(pattern[f]) != "*" && str(pattern[f]) != str(state[f]) {
			return false
		}
	}
	return true
}

func evApply(pattern, state anymap) anymap {
	out := anymap{}
	for _, f := range evFields {
		if str(pattern[f]) == "*" {
			out[f] = state[f]
		} else {
			out[f] = pattern[f]
		}
	}
	return out
}

type evPlan struct {
	actions                         []string
	state                           anymap
	duration, cost, belief, comfort float64
	fuel                            int
}

func evSearch(d anymap) ([]evPlan, int) {
	plans := []evPlan{}
	maxDepth := 0
	start := anymap{}
	veh := asMap(d["Vehicle"])
	for _, f := range evFields {
		start[f] = veh[f]
	}
	goal := asMap(d["Goal"])
	th := asMap(d["Thresholds"])
	actions := maps(d["Actions"])
	var walk func(anymap, []string, float64, float64, float64, float64, int, map[string]bool)
	walk = func(state anymap, path []string, dur, cost, belief, comfort float64, fuel int, seen map[string]bool) {
		if len(path) > maxDepth {
			maxDepth = len(path)
		}
		if evMatches(goal, state) {
			if belief > num(th["MinBelief"]) && cost < num(th["MaxCost"]) && dur < num(th["MaxDuration"]) {
				plans = append(plans, evPlan{append([]string(nil), path...), state, dur, cost, belief, comfort, fuel})
			}
			return
		}
		if fuel == 0 {
			return
		}
		for _, a := range actions {
			if !evMatches(asMap(a["From"]), state) {
				continue
			}
			nxt := evApply(asMap(a["To"]), state)
			k := evKey(nxt)
			if seen[k] && k != evKey(state) {
				continue
			}
			ns := map[string]bool{}
			for kk, v := range seen {
				ns[kk] = v
			}
			ns[k] = true
			walk(nxt, append(append([]string(nil), path...), str(a["Name"])), dur+num(a["Duration"]), cost+num(a["Cost"]), belief*num(a["Belief"]), comfort*num(a["Comfort"]), fuel-1, ns)
		}
	}
	walk(start, nil, 0, 0, 1, 1, integer(d["FuelSteps"]), map[string]bool{evKey(start): true})
	sort.Slice(plans, func(i, j int) bool {
		if plans[i].duration != plans[j].duration {
			return plans[i].duration < plans[j].duration
		}
		if plans[i].cost != plans[j].cost {
			return plans[i].cost < plans[j].cost
		}
		return strings.Join(plans[i].actions, "/") < strings.Join(plans[j].actions, "/")
	})
	return plans, maxDepth
}

func odrlNoticeDays(p anymap) (int, bool) {
	for _, d := range maps(p["Duties"]) {
		if str(d["Action"]) == "odrl:inform" {
			for _, c := range maps(d["Constraints"]) {
				if str(c["LeftOperand"]) == "tosl:noticeDays" {
					return integer(asMap(c["RightOperand"])["Int"]), true
				}
			}
		}
	}
	return 0, false
}

func odrlConsent(p anymap) bool {
	for _, c := range maps(p["Constraints"]) {
		lo := str(c["LeftOperand"])
		if lo == "dpv:Consent" || lo == "tosl:explicitConsent" {
			return true
		}
	}
	return false
}

type riskRow struct {
	clause      string
	score       int
	level       string
	mitigations int
}

func odrlRiskRows(d anymap) []riskRow {
	ag := asMap(d["Agreement"])
	pol := asMap(ag["Policy"])
	perms := asMap(pol["Permissions"])
	pros := asMap(pol["Prohibitions"])
	clauses := asMap(ag["Clauses"])
	needs := asMap(asMap(d["Consumer"])["Needs"])
	rows := []riskRow{}
	del := asMap(perms[":PermDeleteAccount"])
	if str(del["Action"]) == "tosl:removeAccount" {
		if _, ok := odrlNoticeDays(del); !ok {
			rows = append(rows, riskRow{str(asMap(clauses[str(del["ClauseID"])])["ID"]), 100, "HighRisk", 2})
		}
	}
	share := asMap(perms[":PermShareData"])
	if str(share["Action"]) == "tosl:shareData" && !odrlConsent(share) {
		rows = append(rows, riskRow{str(asMap(clauses[str(share["ClauseID"])])["ID"]), 97, "HighRisk", 1})
	}
	change := asMap(perms[":PermChangeTerms"])
	got, ok := odrlNoticeDays(change)
	req := integer(asMap(needs[":Need_ChangeOnlyWithPriorNotice"])["MinNoticeDays"])
	if ok && got < req {
		rows = append(rows, riskRow{str(asMap(clauses[str(change["ClauseID"])])["ID"]), 85, "HighRisk", 1})
	}
	exp := asMap(pros[":ProhibitExportData"])
	if str(exp["Action"]) == "tosl:exportData" {
		rows = append(rows, riskRow{str(asMap(clauses[str(exp["ClauseID"])])["ID"]), 70, "ModerateRisk", 1})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].score != rows[j].score {
			return rows[i].score > rows[j].score
		}
		return rows[i].clause < rows[j].clause
	})
	return rows
}

func anyRisk(rows []riskRow, clause string, score int, level string) bool {
	for _, r := range rows {
		if r.clause == clause && r.score == score && r.level == level {
			return true
		}
	}
	return false
}

type mobPair struct{ a, b bool }

func rel(rows []any) map[mobPair]bool {
	out := map[mobPair]bool{}
	for _, rr := range rows {
		arr := asSlice(rr)
		out[mobPair{boolean(arr[0]), boolean(arr[1])}] = true
		out[mobPair{boolean(arr[2]), boolean(arr[3])}] = true
	}
	return out
}

func composeRel(first, second map[mobPair]bool) map[mobPair]bool {
	out := map[mobPair]bool{}
	for p := range first {
		for q := range second {
			if p.b == q.a {
				out[mobPair{p.a, q.b}] = true
			}
		}
	}
	return out
}

func relEq(a, b map[mobPair]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func mapIntEq(a, b map[int]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func delfourPayload(d anymap) string {
	insight := asMap(d["Insight"])
	policy := asMap(d["Policy"])
	perm := asMap(policy["Permission"])
	pc := asMap(perm["Constraint"])
	pro := asMap(policy["Prohibition"])
	prc := asMap(pro["Constraint"])
	duty := asMap(policy["Duty"])
	dc := asMap(duty["Constraint"])
	return fmt.Sprintf(`{"insight":{"createdAt":%q,"expiresAt":%q,"id":%q,"metric":%q,"retailer":%q,"scopeDevice":%q,"scopeEvent":%q,"suggestionPolicy":%q,"threshold":%.1f,"type":"ins:Insight"},"policy":{"duty":{"action":%q,"constraint":{"leftOperand":%q,"operator":%q,"rightOperand":%q}},"permission":{"action":%q,"constraint":{"leftOperand":%q,"operator":%q,"rightOperand":%q},"target":%q},"profile":%q,"prohibition":{"action":%q,"constraint":{"leftOperand":%q,"operator":%q,"rightOperand":%q},"target":%q},"type":"odrl:Policy"}}`, str(insight["CreatedAt"]), str(insight["ExpiresAt"]), str(insight["ID"]), str(insight["Metric"]), str(insight["Retailer"]), str(insight["ScopeDevice"]), str(insight["ScopeEvent"]), str(insight["SuggestionPolicy"]), num(insight["ThresholdG"]), str(duty["Action"]), str(dc["LeftOperand"]), str(dc["Operator"]), str(dc["RightOperand"]), str(perm["Action"]), str(pc["LeftOperand"]), str(pc["Operator"]), str(pc["RightOperand"]), str(perm["Target"]), str(policy["Profile"]), str(pro["Action"]), str(prc["LeftOperand"]), str(prc["Operator"]), str(prc["RightOperand"]), str(pro["Target"]))
}

func ecoAnswer(ans string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(ans, "\n") {
		if strings.Contains(line, " : ") {
			p := strings.SplitN(line, " : ", 2)
			out[strings.TrimSpace(p[0])] = strings.TrimSpace(p[1])
		}
	}
	return out
}

func ecoFuel(route anymap, payload float64) float64 {
	return num(route["distanceKm"]) * payload * num(route["gradientFactor"])
}

func ecoExpiry(issued string, hours int) string {
	return parseTime(strings.ReplaceAll(issued, "Z", "+00:00")).Add(time.Duration(hours) * time.Hour).UTC().Format("2006-01-02T15:04:05Z")
}

func ecoGoNum(v float64) any {
	r := math.Round(v*100) / 100
	if math.Abs(r-math.Round(r)) < 1e-12 {
		return int(math.Round(r))
	}
	return r
}

type ecoAssertions struct {
	ShowEcoBanner      any `json:"showEcoBanner"`
	SuggestedRoute     any `json:"suggestedRoute"`
	CurrentFuelIndex   any `json:"currentFuelIndex"`
	SuggestedFuelIndex any `json:"suggestedFuelIndex"`
	EstimatedSaving    any `json:"estimatedSaving"`
	RawDataExported    any `json:"rawDataExported"`
}

type ecoEnv struct {
	Audience   string        `json:"audience"`
	AllowedUse string        `json:"allowedUse"`
	IssuedAt   string        `json:"issuedAt"`
	Expiry     string        `json:"expiry"`
	KeyID      string        `json:"keyId"`
	Assertions ecoAssertions `json:"assertions"`
}

func schoolRank(student anymap, school string) int {
	prefs := sarr(student["preferences"])
	for i, p := range prefs {
		if p == school {
			return i
		}
	}
	return 999
}

type schoolAssign struct {
	student     anymap
	school      string
	distance    anymap
	score, rank int
}

func chooseSchool(d anymap, st anymap, audited bool) schoolAssign {
	var best *schoolAssign
	penalty := integer(asMap(d["policy"])["preferencePenaltyMeters"])
	for _, row := range maps(d["distances"]) {
		if str(row["student"]) != str(st["id"]) {
			continue
		}
		rank := schoolRank(st, str(row["school"]))
		score := integer(row["straightMeters"])
		if audited {
			score = integer(row["walkingMeters"]) + rank*penalty
		}
		cand := schoolAssign{st, str(row["school"]), row, score, rank}
		if best == nil || cand.score < best.score || (cand.score == best.score && cand.rank < best.rank) || (cand.score == best.score && cand.rank == best.rank && cand.school < best.school) {
			tmp := cand
			best = &tmp
		}
	}
	return *best
}

func hiddenDetour(a schoolAssign) int {
	return integer(a.distance["walkingMeters"]) - integer(a.distance["straightMeters"])
}

func keysBool(m map[string]bool) []string {
	out := []string{}
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
