package main

import (
	"fmt"
)

func checkDeepTaxonomy(ctx *Context) []Check {
	d := ctx.M()
	depth := integer(d["TaxonomyDepth"])
	totalRules := depth + 2 + 7
	typeFacts := (depth + 1) + (2 * depth) + 1
	return []Check{
		{"the requested taxonomy depth is 100000", depth == 100000},
		{"the terminal class N100000 is reported as reached", contains(ctx.Answer, fmt.Sprintf(":ind a :N%d", depth))},
		{"the terminal A2 class and success flag are reported", contains(ctx.Answer, ":ind a :A2") && contains(ctx.Answer, ":test :is true")},
		{"taxonomy-step rule count matches the JSON depth", grabInt(ctx.Reason, `source N3 taxonomy-step rules : (\d+)`) == depth},
		{"total source-rule count recomputes from start, step, terminal, success, and report rules", grabInt(ctx.Reason, `source N3 total rules counted : (\d+)`) == totalRules},
		{"agenda pops match one pop per N-class fact in the main chain", grabInt(ctx.Reason, `agenda pops : (\d+)`) == depth+1},
		{"taxonomy-step applications match the depth", grabInt(ctx.Reason, `taxonomy-step rule applications : (\d+)`) == depth},
		{"terminal and success rules fire exactly once", grabInt(ctx.Reason, `terminal rule applications : (\d+)`) == 1 && grabInt(ctx.Reason, `success rule applications : (\d+)`) == 1},
		{"classification fact total accounts for N, I/J side labels, and A2", contains(ctx.Reason, fmt.Sprint(typeFacts)) && contains(ctx.Reason, "300002 type facts")},
		{"midpoint and endpoint checkpoints are present", contains(ctx.Answer, ":N50000 plus :I50000/:J50000 present : yes") && contains(ctx.Answer, ":N99999 and :N100000 present : yes")},
	}
}
