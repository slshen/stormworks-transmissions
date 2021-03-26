package main

import (
	"flag"
	"fmt"

	"github.com/slshen/stormworks-transmission/pkg/transmission"
)

func main() {
	var (
		gearboxCount int
		search       transmission.SearchParams
		clearCache   bool
	)
	flag.IntVar(&gearboxCount, "gearboxcount", 2, "The number of gearboxes")
	flag.IntVar(&search.MaxResults, "max-results", 10, "Maximum number of results")
	flag.IntVar(&search.UniqueGearMin, "min-gears", 2, "The minimum number or unique gears")
	flag.Float64Var(&search.LowGearMin, "low-gear-min", 0, "Minimum low gear ratio")
	flag.Float64Var(&search.LowGearMax, "low-gear-max", 0, "Maximum low gear ratio")
	flag.IntVar(&search.MaxStep, "max-step", 0, "Maximum step in pct between gears")
	flag.BoolVar(&clearCache, "clear-cache", false, "Clear the transmission cache")
	flag.Parse()
	ts := transmission.GenerateCache(gearboxCount, clearCache)
	ts = transmission.Search(ts, &search)
	for _, t := range ts {
		fmt.Printf("%s\n", t)
	}
}
