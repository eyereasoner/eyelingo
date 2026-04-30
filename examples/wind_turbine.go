// wind_turbine.go
//
// Inspired by Eyeling's `examples/wind-turbine.n3`.
//
// The example classifies wind-speed samples for a turbine power curve and
// computes the certified energy contribution of the usable intervals.
package main

import (
	"eyelingo/internal/exampleinput"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
)

const eyelingoExampleName = "wind_turbine"

type Dataset struct {
	CaseName        string    `json:"caseName"`
	Question        string    `json:"question"`
	CutInMS         float64   `json:"cutInMS"`
	RatedMS         float64   `json:"ratedMS"`
	CutOutMS        float64   `json:"cutOutMS"`
	RatedPowerMW    float64   `json:"ratedPowerMW"`
	IntervalMinutes float64   `json:"intervalMinutes"`
	WindSpeedsMS    []float64 `json:"windSpeedsMS"`
	Expected        struct {
		UsableIntervals  int `json:"usableIntervals"`
		RatedIntervals   int `json:"ratedIntervals"`
		StoppedIntervals int `json:"stoppedIntervals"`
	} `json:"expected"`
}

type Interval struct {
	Index   int
	Speed   float64
	State   string
	PowerMW float64
	Energy  float64
}

type Check struct {
	ID   string
	OK   bool
	Text string
}

type Analysis struct {
	Intervals []Interval
	TotalMWh  float64
	Checks    []Check
}

func main() {
	ds := exampleinput.Load(eyelingoExampleName, Dataset{})
	a := derive(ds)
	printReport(ds, a)
	if !allOK(a.Checks) {
		os.Exit(1)
	}
}

func derive(ds Dataset) Analysis {
	intervals := []Interval{}
	total := 0.0
	for i, v := range ds.WindSpeedsMS {
		state, power := classify(ds, v)
		energy := power * ds.IntervalMinutes / 60.0
		intervals = append(intervals, Interval{Index: i + 1, Speed: v, State: state, PowerMW: power, Energy: energy})
		total += energy
	}
	checks := []Check{
		{"C1", ds.CutInMS < ds.RatedMS && ds.RatedMS < ds.CutOutMS, "cut-in, rated, and cut-out thresholds are ordered"},
		{"C2", countState(intervals, "partial")+countState(intervals, "rated") == ds.Expected.UsableIntervals, "usable intervals are exactly the samples inside the operating envelope"},
		{"C3", countState(intervals, "rated") == ds.Expected.RatedIntervals, "two intervals reach rated power"},
		{"C4", countState(intervals, "stopped") == ds.Expected.StoppedIntervals, "two intervals are stopped by low or cut-out wind"},
		{"C5", intervals[1].PowerMW > 0 && intervals[1].PowerMW < ds.RatedPowerMW, "a below-rated usable wind speed follows the cubic power curve"},
		{"C6", total > 1.5 && total < 1.6, "the six ten-minute samples yield about 1.57 MWh"},
	}
	return Analysis{Intervals: intervals, TotalMWh: total, Checks: checks}
}

func classify(ds Dataset, speed float64) (string, float64) {
	if speed < ds.CutInMS || speed >= ds.CutOutMS {
		return "stopped", 0
	}
	if speed >= ds.RatedMS {
		return "rated", ds.RatedPowerMW
	}
	numer := math.Pow(speed, 3) - math.Pow(ds.CutInMS, 3)
	denom := math.Pow(ds.RatedMS, 3) - math.Pow(ds.CutInMS, 3)
	return "partial", ds.RatedPowerMW * numer / denom
}

func countState(xs []Interval, state string) int {
	n := 0
	for _, x := range xs {
		if x.State == state {
			n++
		}
	}
	return n
}

func intervalLine(xs []Interval) string {
	parts := []string{}
	for _, x := range xs {
		parts = append(parts, fmt.Sprintf("t%d %.1f m/s %s %.3f MW", x.Index, x.Speed, x.State, x.PowerMW))
	}
	return strings.Join(parts, "; ")
}

func allOK(checks []Check) bool {
	for _, c := range checks {
		if !c.OK {
			return false
		}
	}
	return true
}

func countOK(checks []Check) int {
	n := 0
	for _, c := range checks {
		if c.OK {
			n++
		}
	}
	return n
}

func printReport(ds Dataset, a Analysis) {
	fmt.Println("# Wind Turbine Envelope")
	fmt.Println()
	fmt.Println("## Answer")
	fmt.Printf("operating thresholds : cut-in %.1f m/s, rated %.1f m/s, cut-out %.1f m/s\n", ds.CutInMS, ds.RatedMS, ds.CutOutMS)
	fmt.Printf("rated power : %.1f MW\n", ds.RatedPowerMW)
	fmt.Printf("interval classifications : %s\n", intervalLine(a.Intervals))
	fmt.Printf("usable intervals : %d\n", countState(a.Intervals, "partial")+countState(a.Intervals, "rated"))
	fmt.Printf("total energy : %.3f MWh\n", a.TotalMWh)
	fmt.Println()
	fmt.Println("## Reason why")
	fmt.Println("Wind below cut-in and at or above cut-out is stopped for production and safety.")
	fmt.Println("Wind between cut-in and rated speed follows a cubic power curve normalized to the rated point.")
	fmt.Println("Wind between rated speed and cut-out is capped at rated power.")
	fmt.Println("Energy is accumulated by multiplying each interval power by the ten-minute interval duration.")
	fmt.Println()
	fmt.Println("## Check")
	for _, c := range a.Checks {
		status := "FAIL"
		if c.OK {
			status = "OK"
		}
		fmt.Printf("%s %s - %s\n", c.ID, status, c.Text)
	}
	fmt.Println()
	fmt.Println("## Go audit details")
	fmt.Printf("platform : %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Printf("case : %s\n", ds.CaseName)
	fmt.Printf("question : %s\n", ds.Question)
	fmt.Printf("samples evaluated : %d\n", len(ds.WindSpeedsMS))
	fmt.Printf("interval minutes : %.1f\n", ds.IntervalMinutes)
	fmt.Printf("partial intervals : %d\n", countState(a.Intervals, "partial"))
	fmt.Printf("rated intervals : %d\n", countState(a.Intervals, "rated"))
	fmt.Printf("stopped intervals : %d\n", countState(a.Intervals, "stopped"))
	fmt.Printf("checks passed : %d/%d\n", countOK(a.Checks), len(a.Checks))
}
