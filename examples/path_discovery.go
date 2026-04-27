// path_discovery.go
//
// A self-contained Go translation of path-discovery.n3 from the Eyeling
// examples.
//
// The original N3 file contains a large Neptune air-routes graph and a bounded
// recursive rule:
//
//	(?source ?destination () 0 2) :route ?airports
//
// This Go version translates the query, the recursive no-revisit route rule,
// and the exact bounded search slice needed for that query. The full N3 source
// has 7,698 airport labels and 37,505 outbound-route facts; only the 338
// outbound facts reachable from Ostend-Bruges within the two-stopover bound are
// needed to prove the same answer.
//
// This is intentionally not a generic RDF/N3 reasoner. The concrete facts and
// rules are represented as Go structs and ordinary functions so the route
// derivation remains visible and directly runnable.
//
// Run:
//
//	go run path_discovery.go
//
// The program has no third-party dependencies.
package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
)

const (
	sourceGraphAirportLabels = 7698
	sourceGraphOutboundFacts = 37505
	maxStopovers             = 2
	maxHops                  = maxStopovers + 1

	sourceAirport      = "res:AIRPORT_310"
	destinationAirport = "res:AIRPORT_1587"
	mandatoryFirstHop  = "res:AIRPORT_309"
)

type Dataset struct {
	Question      string
	SourceID      string
	DestinationID string
	Labels        map[string]string
	Edges         []Edge
}

type Edge struct {
	From string
	To   string
}

type Route struct {
	Airports []string
}

type SearchStats struct {
	RecursiveCalls   int
	EdgeTests        int
	EdgesExtended    int
	RevisitPrunes    int
	DepthLimitLeaves int
	DeadEnds         int
	RoutesEmitted    int
	MaxDepth         int
}

type Checks struct {
	SourceAndDestinationKnown bool
	FirstHopMatchesN3Facts    bool
	RouteSetMatchesN3Query    bool
	NoShorterRouteExists      bool
	RoutesWithinStopoverLimit bool
	EveryHopHasFact           bool
	NoAirportRevisited        bool
	CompleteBoundedSliceUsed  bool
	RoutesSortedDeterministic bool
}

type InferenceResult struct {
	Routes           []Route
	Stats            SearchStats
	Checks           Checks
	RelevantAirports int
	ExpandedAirports []string
	SourceOut        []string
	FirstHopOut      []string
	DirectRoutes     int
	OneStopRoutes    int
	TwoStopRoutes    int
}

func fixture() Dataset {
	return Dataset{
		Question:      "Find routes from Ostend-Bruges International Airport to Václav Havel Airport Prague with at most 2 stopovers.",
		SourceID:      sourceAirport,
		DestinationID: destinationAirport,
		Labels: map[string]string{
			"res:AIRPORT_1064": "Al Massira Airport",
			"res:AIRPORT_1321": "Bastia-Poretta Airport",
			"res:AIRPORT_1324": "Ajaccio-Napoléon Bonaparte Airport",
			"res:AIRPORT_1399": "Lille-Lesquin Airport",
			"res:AIRPORT_1452": "Heraklion International Nikos Kazantzakis Airport",
			"res:AIRPORT_1472": "Diagoras Airport",
			"res:AIRPORT_1587": "Václav Havel Airport Prague",
			"res:AIRPORT_309":  "Liège Airport",
			"res:AIRPORT_310":  "Ostend-Bruges International Airport",
			"res:AIRPORT_3998": "Palma De Mallorca Airport",
		},
		Edges: []Edge{
			{From: "res:AIRPORT_310", To: "res:AIRPORT_309"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_1324"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_1064"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_1321"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_1472"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_1452"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_1399"},
			{From: "res:AIRPORT_309", To: "res:AIRPORT_3998"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1361"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1264"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1403"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1278"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1412"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1285"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1382"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1423"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1665"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1399"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_309"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_502"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_629"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1335"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1353"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1359"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1418"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1354"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1520"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_644"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1386"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1268"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_737"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1435"},
			{From: "res:AIRPORT_1324", To: "res:AIRPORT_1273"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_580"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_337"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_302"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_304"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1382"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_5670"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_345"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_340"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1054"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_5672"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_348"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1399"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_502"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_629"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1335"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_478"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1353"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1074"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_346"},
			{From: "res:AIRPORT_1064", To: "res:AIRPORT_1386"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_351"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1264"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1403"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_302"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1382"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_344"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_345"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_340"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1665"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_342"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1399"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_309"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_502"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_629"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1335"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1353"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1418"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1354"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1386"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1350"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1435"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_350"},
			{From: "res:AIRPORT_1321", To: "res:AIRPORT_1273"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_580"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_337"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_351"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_2939"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_494"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_302"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_304"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1382"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1450"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_344"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_4029"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_338"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_345"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_523"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_585"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_3941"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1423"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_340"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_534"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_342"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_352"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1452"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1525"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1459"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_4196"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_3956"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1458"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_517"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_348"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_491"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_309"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_502"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_548"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_629"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_478"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_657"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_346"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_521"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_347"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_2948"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1476"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_2985"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_699"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_350"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1486"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1613"},
			{From: "res:AIRPORT_1472", To: "res:AIRPORT_1587"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_580"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1590"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_337"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_351"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_469"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1264"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_2939"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_490"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_302"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1382"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_344"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1472"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_4029"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_338"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_345"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_523"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_535"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_586"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_3941"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1423"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_340"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1665"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_534"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1609"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_342"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_352"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_4191"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1458"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1197"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_517"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_348"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_415"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1399"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1611"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_502"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_492"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_629"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1335"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1524"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_478"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1353"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_346"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1469"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1418"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_521"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_347"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1386"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_2948"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_3953"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_591"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1612"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_2985"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_350"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1486"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1613"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_3959"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_2988"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1587"},
			{From: "res:AIRPORT_1452", To: "res:AIRPORT_1678"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1324"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1064"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1070"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1218"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1321"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1280"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1264"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1472"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_293"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1200"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_231"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1626"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1323"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1636"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1051"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1665"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1452"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_210"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1460"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_309"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1335"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1353"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1075"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1359"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1418"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1354"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1206"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1435"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1273"},
			{From: "res:AIRPORT_1399", To: "res:AIRPORT_1551"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_628"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_607"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1229"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1212"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1213"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_580"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1214"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1218"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_465"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_636"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_337"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_351"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1676"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1216"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_608"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_469"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_514"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1538"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1264"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_494"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_353"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_490"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_302"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_304"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_488"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1553"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_344"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_609"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_596"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_4029"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_373"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_338"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_599"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_345"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_523"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_535"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_586"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_585"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_339"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1423"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_552"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1626"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1223"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1636"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_340"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_355"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_382"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1665"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_467"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1222"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_534"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_537"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_691"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_687"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1609"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_342"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_352"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_421"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_210"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1225"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1525"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1610"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1226"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_4166"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_400"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_3956"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_517"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_348"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1611"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_491"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_309"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_7459"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_503"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_502"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_507"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_492"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_548"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_629"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1335"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_364"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1745"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_582"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1524"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_478"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1353"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_3986"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1231"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_8414"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_657"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_346"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1230"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_341"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_5673"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1418"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_521"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_347"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_644"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_310"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_371"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1367"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1386"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1236"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_5562"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_772"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_591"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_349"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1612"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_664"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1251"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1243"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1253"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_603"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1194"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_495"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_508"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_699"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_737"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_350"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1273"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1246"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1613"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1587"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_4198"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1252"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_393"},
			{From: "res:AIRPORT_3998", To: "res:AIRPORT_1678"},
		},
	}
}

func infer(data Dataset) InferenceResult {
	adj := buildAdjacency(data)
	edgeSet := buildEdgeSet(data.Edges)
	stats := SearchStats{}
	routes := []Route{}
	dfs(data.SourceID, data.DestinationID, adj, []string{data.SourceID}, &routes, &stats)
	sortRoutes(data, routes)

	direct, oneStop, twoStop := routeDistribution(routes)
	expanded := expandedAirports(data.SourceID, adj, 2)
	firstHopOut := append([]string{}, adj[mandatoryFirstHop]...)
	sortTermsByLabel(data, firstHopOut)

	checks := Checks{
		SourceAndDestinationKnown: data.label(data.SourceID) != data.SourceID && data.label(data.DestinationID) != data.DestinationID,
		FirstHopMatchesN3Facts:    len(adj[data.SourceID]) == 1 && adj[data.SourceID][0] == mandatoryFirstHop,
		RouteSetMatchesN3Query:    routeSetMatches(data, routes, expectedRoutes()),
		NoShorterRouteExists:      direct == 0 && oneStop == 0,
		RoutesWithinStopoverLimit: routesWithinStopovers(routes, maxStopovers),
		EveryHopHasFact:           everyHopHasFact(routes, edgeSet),
		NoAirportRevisited:        noAirportRevisited(routes),
		CompleteBoundedSliceUsed:  stats.EdgeTests == len(data.Edges),
		RoutesSortedDeterministic: routesSorted(data, routes),
	}

	return InferenceResult{
		Routes:           routes,
		Stats:            stats,
		Checks:           checks,
		RelevantAirports: countAirportsInEdges(data.Edges),
		ExpandedAirports: expanded,
		SourceOut:        append([]string{}, adj[data.SourceID]...),
		FirstHopOut:      firstHopOut,
		DirectRoutes:     direct,
		OneStopRoutes:    oneStop,
		TwoStopRoutes:    twoStop,
	}
}

func buildAdjacency(data Dataset) map[string][]string {
	adj := map[string][]string{}
	for _, edge := range data.Edges {
		adj[edge.From] = append(adj[edge.From], edge.To)
	}
	for from := range adj {
		sortTermsByLabel(data, adj[from])
	}
	return adj
}

func buildEdgeSet(edges []Edge) map[string]bool {
	set := map[string]bool{}
	for _, edge := range edges {
		set[edge.From+" -> "+edge.To] = true
	}
	return set
}

func dfs(current, destination string, adj map[string][]string, path []string, routes *[]Route, stats *SearchStats) {
	depth := len(path) - 1
	stats.RecursiveCalls++
	if depth > stats.MaxDepth {
		stats.MaxDepth = depth
	}

	if len(path) > 1 && current == destination {
		stats.RoutesEmitted++
		copyPath := append([]string{}, path...)
		*routes = append(*routes, Route{Airports: copyPath})
		return
	}

	if depth >= maxHops {
		stats.DepthLimitLeaves++
		return
	}

	outbound := adj[current]
	stats.EdgeTests += len(outbound)
	if len(outbound) == 0 {
		stats.DeadEnds++
	}

	for _, next := range outbound {
		if contains(path, next) {
			stats.RevisitPrunes++
			continue
		}
		stats.EdgesExtended++
		nextPath := append(append([]string{}, path...), next)
		dfs(next, destination, adj, nextPath, routes, stats)
	}
}

func expectedRoutes() [][]string {
	return [][]string{
		{sourceAirport, mandatoryFirstHop, "res:AIRPORT_1472", destinationAirport},
		{sourceAirport, mandatoryFirstHop, "res:AIRPORT_1452", destinationAirport},
		{sourceAirport, mandatoryFirstHop, "res:AIRPORT_3998", destinationAirport},
	}
}

func routeSetMatches(data Dataset, actual []Route, expected [][]string) bool {
	actualKeys := make([]string, 0, len(actual))
	for _, route := range actual {
		actualKeys = append(actualKeys, strings.Join(route.Airports, "|"))
	}
	expectedKeys := make([]string, 0, len(expected))
	for _, route := range expected {
		expectedKeys = append(expectedKeys, strings.Join(route, "|"))
	}
	sort.Strings(actualKeys)
	sort.Strings(expectedKeys)
	if len(actualKeys) != len(expectedKeys) {
		return false
	}
	for i := range actualKeys {
		if actualKeys[i] != expectedKeys[i] {
			return false
		}
	}
	return routesSorted(data, actual)
}

func routesWithinStopovers(routes []Route, limit int) bool {
	for _, route := range routes {
		if route.Stopovers() > limit {
			return false
		}
	}
	return true
}

func everyHopHasFact(routes []Route, edgeSet map[string]bool) bool {
	for _, route := range routes {
		for i := 0; i < len(route.Airports)-1; i++ {
			key := route.Airports[i] + " -> " + route.Airports[i+1]
			if !edgeSet[key] {
				return false
			}
		}
	}
	return true
}

func noAirportRevisited(routes []Route) bool {
	for _, route := range routes {
		seen := map[string]bool{}
		for _, airport := range route.Airports {
			if seen[airport] {
				return false
			}
			seen[airport] = true
		}
	}
	return true
}

func routesSorted(data Dataset, routes []Route) bool {
	for i := 1; i < len(routes); i++ {
		if routeLabel(data, routes[i-1]) > routeLabel(data, routes[i]) {
			return false
		}
	}
	return true
}

func sortRoutes(data Dataset, routes []Route) {
	sort.Slice(routes, func(i, j int) bool {
		return routeLabel(data, routes[i]) < routeLabel(data, routes[j])
	})
}

func sortTermsByLabel(data Dataset, terms []string) {
	sort.Slice(terms, func(i, j int) bool {
		left := data.label(terms[i])
		right := data.label(terms[j])
		if left == right {
			return terms[i] < terms[j]
		}
		return left < right
	})
}

func routeDistribution(routes []Route) (direct, oneStop, twoStop int) {
	for _, route := range routes {
		switch route.Stopovers() {
		case 0:
			direct++
		case 1:
			oneStop++
		case 2:
			twoStop++
		}
	}
	return direct, oneStop, twoStop
}

func expandedAirports(source string, adj map[string][]string, maxDepth int) []string {
	seen := map[string]bool{}
	ordered := []string{}
	var walk func(string, []string)
	walk = func(current string, path []string) {
		depth := len(path) - 1
		if depth > maxDepth {
			return
		}
		if !seen[current] {
			seen[current] = true
			ordered = append(ordered, current)
		}
		if depth == maxDepth {
			return
		}
		for _, next := range adj[current] {
			if contains(path, next) {
				continue
			}
			walk(next, append(append([]string{}, path...), next))
		}
	}
	walk(source, []string{source})
	return ordered
}

func countAirportsInEdges(edges []Edge) int {
	seen := map[string]bool{}
	for _, edge := range edges {
		seen[edge.From] = true
		seen[edge.To] = true
	}
	return len(seen)
}

func countChecks(checks Checks) (passed, total int) {
	values := []bool{
		checks.SourceAndDestinationKnown,
		checks.FirstHopMatchesN3Facts,
		checks.RouteSetMatchesN3Query,
		checks.NoShorterRouteExists,
		checks.RoutesWithinStopoverLimit,
		checks.EveryHopHasFact,
		checks.NoAirportRevisited,
		checks.CompleteBoundedSliceUsed,
		checks.RoutesSortedDeterministic,
	}
	for _, ok := range values {
		if ok {
			passed++
		}
	}
	return passed, len(values)
}

func contains(path []string, term string) bool {
	for _, item := range path {
		if item == term {
			return true
		}
	}
	return false
}

func (data Dataset) label(term string) string {
	if label, ok := data.Labels[term]; ok {
		return label
	}
	return term
}

func (route Route) Stopovers() int {
	if len(route.Airports) < 2 {
		return 0
	}
	return len(route.Airports) - 2
}

func (route Route) Hops() int {
	if len(route.Airports) == 0 {
		return 0
	}
	return len(route.Airports) - 1
}

func routeLabel(data Dataset, route Route) string {
	labels := make([]string, 0, len(route.Airports))
	for _, airport := range route.Airports {
		labels = append(labels, data.label(airport))
	}
	return strings.Join(labels, " -> ")
}

func routeTerms(route Route) string {
	return strings.Join(route.Airports, " -> ")
}

func status(ok bool) string {
	if ok {
		return "OK"
	}
	return "FAIL"
}

func main() {
	data := fixture()
	result := infer(data)
	checksPassed, checksTotal := countChecks(result.Checks)

	fmt.Println("=== Answer ===")
	fmt.Printf("The path discovery query finds %d air routes with at most %d stopovers.\n", len(result.Routes), maxStopovers)
	fmt.Printf("from : %s\n", data.label(data.SourceID))
	fmt.Printf("to : %s\n", data.label(data.DestinationID))
	fmt.Printf("max stopovers : %d\n", maxStopovers)
	fmt.Println()
	fmt.Println("Discovered routes:")
	for i, route := range result.Routes {
		fmt.Printf(" - route %d (%d stopovers): %s\n", i+1, route.Stopovers(), routeLabel(data, route))
	}

	fmt.Println()
	fmt.Println("=== Reason Why ===")
	fmt.Println("The N3 source defines a recursive :route relation over nepo:hasOutboundRouteTo facts. A route can use a direct edge when the current length is within the maximum, or extend through a non-visited intermediate airport and recurse with length+1. The final log:collectAllIn query collects the labels of each airport in every route from the source to the destination.")
	fmt.Printf("source N3 airport labels : %d\n", sourceGraphAirportLabels)
	fmt.Printf("source N3 outbound-route facts : %d\n", sourceGraphOutboundFacts)
	fmt.Printf("translated bounded-search airports : %d\n", result.RelevantAirports)
	fmt.Printf("translated bounded-search outbound facts : %d\n", len(data.Edges))
	fmt.Printf("frontier airports expanded : %d\n", len(result.ExpandedAirports))
	fmt.Printf("source outbound candidates : %d\n", len(result.SourceOut))
	fmt.Printf("Liège outbound candidates : %d\n", len(result.FirstHopOut))
	fmt.Printf("direct routes : %d\n", result.DirectRoutes)
	fmt.Printf("one-stop routes : %d\n", result.OneStopRoutes)
	fmt.Printf("two-stopover routes : %d\n", result.TwoStopRoutes)
	fmt.Printf("search recursive calls : %d\n", result.Stats.RecursiveCalls)
	fmt.Printf("search edge tests : %d\n", result.Stats.EdgeTests)
	fmt.Printf("search depth-limit leaves : %d\n", result.Stats.DepthLimitLeaves)
	fmt.Println("Second-hop candidates from Liège:")
	for _, airport := range result.FirstHopOut {
		fmt.Printf(" - %s (%s)\n", data.label(airport), airport)
	}

	fmt.Println()
	fmt.Println("=== Check ===")
	fmt.Printf("C1 %s - source and destination airport labels are known.\n", status(result.Checks.SourceAndDestinationKnown))
	fmt.Printf("C2 %s - Ostend-Bruges has one outbound route in the bounded proof slice, to Liège Airport.\n", status(result.Checks.FirstHopMatchesN3Facts))
	fmt.Printf("C3 %s - the discovered route set matches the N3 query answer.\n", status(result.Checks.RouteSetMatchesN3Query))
	fmt.Printf("C4 %s - no direct or one-stop route exists under the same bound.\n", status(result.Checks.NoShorterRouteExists))
	fmt.Printf("C5 %s - every discovered route has at most two stopovers.\n", status(result.Checks.RoutesWithinStopoverLimit))
	fmt.Printf("C6 %s - every hop is backed by a nepo:hasOutboundRouteTo fact.\n", status(result.Checks.EveryHopHasFact))
	fmt.Printf("C7 %s - no route revisits an airport.\n", status(result.Checks.NoAirportRevisited))
	fmt.Printf("C8 %s - the search scanned the complete bounded slice of 338 outbound facts.\n", status(result.Checks.CompleteBoundedSliceUsed))
	fmt.Printf("C9 %s - route output is sorted deterministically by airport labels.\n", status(result.Checks.RoutesSortedDeterministic))

	fmt.Println()
	fmt.Println("=== Go audit details ===")
	fmt.Printf("go runtime : %s\n", runtime.Version())
	fmt.Printf("go os/arch : %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("question : %s\n", data.Question)
	fmt.Printf("source airport : %s (%s)\n", data.label(data.SourceID), data.SourceID)
	fmt.Printf("destination airport : %s (%s)\n", data.label(data.DestinationID), data.DestinationID)
	fmt.Printf("source graph airport labels : %d\n", sourceGraphAirportLabels)
	fmt.Printf("source graph outbound facts : %d\n", sourceGraphOutboundFacts)
	fmt.Printf("translated bounded slice airports : %d\n", result.RelevantAirports)
	fmt.Printf("translated bounded slice outbound facts : %d\n", len(data.Edges))
	fmt.Printf("max stopovers : %d\n", maxStopovers)
	fmt.Printf("max hops : %d\n", maxHops)
	fmt.Printf("routes discovered : %d\n", len(result.Routes))
	fmt.Printf("mandatory first hop : %s (%s)\n", data.label(mandatoryFirstHop), mandatoryFirstHop)
	fmt.Println("expanded airports:")
	for _, airport := range result.ExpandedAirports {
		fmt.Printf(" - %s (%s)\n", data.label(airport), airport)
	}
	for i, route := range result.Routes {
		fmt.Printf("route %d terms : %s\n", i+1, routeTerms(route))
		fmt.Printf("route %d labels : %s\n", i+1, routeLabel(data, route))
		fmt.Printf("route %d hops : %d\n", i+1, route.Hops())
		fmt.Printf("route %d stopovers : %d\n", i+1, route.Stopovers())
	}
	fmt.Printf("search recursive calls : %d\n", result.Stats.RecursiveCalls)
	fmt.Printf("search edge tests : %d\n", result.Stats.EdgeTests)
	fmt.Printf("search edges extended : %d\n", result.Stats.EdgesExtended)
	fmt.Printf("search revisit prunes : %d\n", result.Stats.RevisitPrunes)
	fmt.Printf("search depth-limit leaves : %d\n", result.Stats.DepthLimitLeaves)
	fmt.Printf("search dead ends : %d\n", result.Stats.DeadEnds)
	fmt.Printf("search routes emitted : %d\n", result.Stats.RoutesEmitted)
	fmt.Printf("search max depth : %d\n", result.Stats.MaxDepth)
	fmt.Printf("checks passed : %d/%d\n", checksPassed, checksTotal)
	fmt.Printf("all checks pass : %s\n", map[bool]string{true: "yes", false: "no"}[checksPassed == checksTotal])

	if checksPassed != checksTotal {
		os.Exit(1)
	}
}
