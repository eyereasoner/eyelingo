package main

func checkAckermann(ctx *Context) []Check {
	data := maps(ctx.Load())
	values := map[string]*bigInt{}
	small := map[string]*bigInt{}
	// Represent huge values by decimal summaries when necessary.
	type fp struct {
		digits                 int
		leading, trailing, sha string
	}
	summarize := func(v *bigInt) fp {
		s := v.String()
		l := len(s)
		lead := s
		if len(lead) > 24 {
			lead = lead[:24]
		}
		tail := s
		if len(tail) > 24 {
			tail = tail[len(tail)-24:]
		}
		return fp{l, lead, tail, sha256Hex(s)}
	}
	for _, item := range data {
		values[str(item["ID"])] = ackermannBinary(integer(item["X"]), integer(item["Y"]))
	}
	for _, m := range reAll(ctx.Answer, `(?m)^(A\d+) ackermann\(\d+,\d+\) = (\d+)$`) {
		small[m[1]] = newBigIntFromString(m[2])
	}
	fps := map[string]fp{}
	for _, m := range reAll(ctx.Answer, `(?m)^(A\d+) digits=(\d+) leading=([0-9]+) trailing=([0-9]+) sha256=([0-9a-f]+)$`) {
		fps[m[1]] = fp{parseInt(m[2]), m[3], m[4], m[5]}
	}
	smallOK := true
	for k, v := range small {
		if values[k] == nil || values[k].Cmp(v) != 0 {
			smallOK = false
		}
	}
	return []Check{
		{"all twelve JSON Ackermann queries are recomputed", len(values) == 12},
		{"x=0 queries reduce to successor after the +3 binary offset", values["A0"].EqInt(1) && values["A1"].EqInt(7)},
		{"x=1 queries reduce to addition after the +3 binary offset", values["A2"].EqInt(4) && values["A3"].EqInt(9)},
		{"x=2 queries reduce to multiplication after the +3 binary offset", values["A4"].EqInt(7) && values["A5"].EqInt(21)},
		{"x=3 queries reduce to exact base-2 exponentiation", values["A6"].EqInt(125) && summarize(values["A7"]) == fps["A7"]},
		{"A(4,0) and A(4,1) match the first tetration cases", values["A8"].EqInt(13) && values["A9"].EqInt(65533)},
		{"A(4,2) is held exactly as 2^65536-3 with the reported fingerprint", summarize(values["A10"]) == fps["A10"]},
		{"A(5,0) lands on the same value as A(4,1)", values["A11"].Cmp(values["A9"]) == 0 && values["A9"].EqInt(65533)},
		{"all non-huge reported exact values match recomputation", smallOK},
		{"the reported proof statistics match the query and memo structure", contains(ctx.Reason, "primitive test queries : 12") && contains(ctx.Reason, "binary reductions : 12") && contains(ctx.Reason, "distinct ternary facts : 23")},
	}
}
