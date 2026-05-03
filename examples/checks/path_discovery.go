package main

import (
	"sort"
	"strings"
)

func checkPathDiscovery(ctx *Context) []Check {
	d := ctx.M()
	routes, graph, _, _ := pathDiscover(d)
	labels := asMap(d["Labels"])
	labelRoutes := [][]string{}
	for _, r := range routes {
		lr := []string{}
		for _, id := range r {
			lr = append(lr, str(labels[id]))
		}
		labelRoutes = append(labelRoutes, lr)
	}
	reported := [][]string{}
	for _, m := range reAll(ctx.Answer, `route (\d+) \(2 stopovers\): ([^\n]+)`) {
		parts := []string{}
		for _, p := range strings.Split(m[2], "->") {
			parts = append(parts, strings.TrimSpace(p))
		}
		reported = append(reported, parts)
	}
	second := []string{}
	firstOut := graph[str(d["SourceID"])][0]
	for _, x := range graph[firstOut] {
		second = append(second, str(labels[x]))
	}
	sort.Strings(second)
	noRevisit := true
	backed := true
	for _, r := range routes {
		seen := map[string]bool{}
		for i, id := range r {
			if seen[id] {
				noRevisit = false
			}
			seen[id] = true
			if i+1 < len(r) {
				found := false
				for _, to := range graph[id] {
					if to == r[i+1] {
						found = true
					}
				}
				backed = backed && found
			}
		}
	}
	sortedOK := true
	for i := 1; i < len(labelRoutes); i++ {
		if strings.Join(labelRoutes[i-1], "|") > strings.Join(labelRoutes[i], "|") {
			sortedOK = false
		}
	}
	expectedSecond := []string{"Ajaccio-Napoléon Bonaparte Airport", "Al Massira Airport", "Bastia-Poretta Airport", "Diagoras Airport", "Heraklion International Nikos Kazantzakis Airport", "Lille-Lesquin Airport", "Palma De Mallorca Airport"}
	return []Check{{"source and destination airport labels are known", str(labels[str(d["SourceID"])]) == "Ostend-Bruges International Airport" && str(labels[str(d["DestinationID"])]) == "Václav Havel Airport Prague"}, {"Ostend-Bruges has one outbound route in the full graph, to Liège Airport", len(graph[str(d["SourceID"])]) == 1 && str(labels[graph[str(d["SourceID"])][0]]) == "Liège Airport"}, {"bounded DFS independently finds exactly three two-stopover routes", len(routes) == 3 && allPathLen(routes, 4)}, {"reported route labels match the independently discovered route set", routesEq(reported, labelRoutes)}, {"no direct or one-stop route exists under the same bound", noShort(routes)}, {"every discovered hop is backed by an outbound-route fact", backed}, {"no discovered route revisits an airport", noRevisit}, {"the translated graph size matches the full source counts", len(labels) == 7698 && countGraphEdges(graph) == 37505}, {"the second-hop candidates from Liège are independently recovered", sliceEq(second, expectedSecond)}, {"route output is sorted deterministically by airport labels", sortedOK}}
}
