package transmission

import (
	"sort"
)

type SearchParams struct {
	LowGearMin    float64
	LowGearMax    float64
	HighGearMin   float64
	HighGearMax   float64
	UniqueGearMin int
	UniqueGearMax int
	MaxResults    int
	MaxStep       int
}

func Search(transmissions []*Transmission, p *SearchParams) (result []*Transmission) {
	for _, t := range transmissions {
		low, _ := t.FinalRatios[0].Ratio.Float64()
		if p.LowGearMin != 0 && low < p.LowGearMin {
			//fmt.Printf("p.LowGearMin=%f low=%f\n", p.LowGearMin, low)
			continue
		}
		if p.LowGearMax != 0 && low > p.LowGearMax {
			continue
		}
		high, _ := t.FinalRatios[len(t.FinalRatios)-1].Ratio.Float64()
		if p.HighGearMin != 0 && high < p.HighGearMin {
			continue
		}
		if p.HighGearMax != 0 && high > p.HighGearMax {
			continue
		}
		if len(t.FinalRatios) < p.UniqueGearMin {
			continue
		}
		if p.UniqueGearMax > 0 && len(t.FinalRatios) > p.UniqueGearMax {
			continue
		}
		if p.MaxStep > 0 {
			gearCount := len(t.FinalRatios)
			if p.UniqueGearMin > 0 {
				gearCount = p.UniqueGearMin
			}
			if !t.HasStepRange(gearCount, float64(p.MaxStep)/100.0) {
				continue
			}
		}
		result = append(result, t)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].StepDeviation < result[j].StepDeviation
	})
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].MaxStep < result[j].MaxStep
	})
	if p.MaxResults != 0 && len(result) > p.MaxResults {
		result = result[0:p.MaxResults]
	}
	return
}
