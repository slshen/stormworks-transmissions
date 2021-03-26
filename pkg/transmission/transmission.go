package transmission

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"sort"
	"strings"

	"github.com/montanaflynn/stats"
)

func NewRatio(a, b int64) *big.Rat {
	return big.NewRat(a, b)
}

var (
	ONE_ONE    = NewRatio(1, 1)
	SIX_FIVE   = NewRatio(6, 5)
	THREE_TWO  = NewRatio(3, 2)
	NINE_FIVE  = NewRatio(9, 5)
	TWO_ONE    = NewRatio(2, 1)
	FIVE_TWO   = NewRatio(5, 2)
	THREE_ONE  = NewRatio(3, 1)
	gearRatios = []*big.Rat{
		ONE_ONE,
		SIX_FIVE,
		THREE_TWO,
		NINE_FIVE,
		TWO_ONE,
		FIVE_TWO,
		THREE_ONE,
	}
)

func ToRatioString(r *big.Rat) string {
	return fmt.Sprintf("%d:%d", r.Num(), r.Denom())
}

type Transmission struct {
	Gearboxes     []*Gearbox
	FinalRatios   []*FinalRatio
	Steps         []float64
	Overall       float64
	MaxStep       float64
	StepDeviation float64
}

type FinalRatio struct {
	Bits  uint
	Ratio *big.Rat
}

func NewTransmission(gearboxes []*Gearbox) *Transmission {
	t := &Transmission{
		Gearboxes: gearboxes,
	}
	t.calculateFinalRatios()
	return t
}

func (t *Transmission) Key() string {
	var s strings.Builder
	for i, r := range t.FinalRatios {
		if i > 0 {
			s.WriteRune(' ')
		}
		s.WriteString(ToRatioString(r.Ratio))
	}
	return s.String()
}

func (t *Transmission) getFinalRatio(bits uint) *big.Rat {
	ratio := big.NewRat(1, 1)
	for i, gearbox := range t.Gearboxes {
		var gear *big.Rat
		set := (bits>>i)&1 == 1
		if set {
			gear = gearbox.On
		} else {
			gear = gearbox.Off
		}
		ratio.Mul(ratio, gear)
	}
	return ratio
}

func (t *Transmission) calculateFinalRatios() {
	m := uint(1) << uint(len(t.Gearboxes))
bitloop:
	for bits := uint(0); bits < m; bits++ {
		ratio := t.getFinalRatio(bits)
		for _, fr := range t.FinalRatios {
			if fr.Ratio.Cmp(ratio) == 0 {
				continue bitloop
			}
		}
		t.FinalRatios = append(t.FinalRatios, &FinalRatio{
			Bits:  bits,
			Ratio: ratio,
		})
	}
	sort.Slice(t.FinalRatios, func(i, j int) bool {
		return t.FinalRatios[i].Ratio.Cmp(t.FinalRatios[j].Ratio) < 0
	})
	high, _ := t.FinalRatios[len(t.FinalRatios)-1].Ratio.Float64()
	low, _ := t.FinalRatios[0].Ratio.Float64()
	t.Overall = high / low
	t.Steps = make([]float64, len(t.FinalRatios)-1)
	p := 0.
	for i, c := range t.FinalRatios {
		f, _ := c.Ratio.Float64()
		if i > 0 {
			step := f/p - 1.
			t.Steps[i-1] = step
		}
		p = f
	}
	t.MaxStep, _ = stats.Max(t.Steps)
	nSteps := float64(len(t.Steps))
	idealStep := math.Exp(math.Log(t.Overall)/nSteps) - 1
	for i := range t.Steps {
		t.StepDeviation += math.Abs(t.Steps[i] - idealStep)
	}
	t.StepDeviation /= nSteps
}

// return true if the maximum step between any gearCount consecutive gears
// is less than or equal to maxStep
func (t *Transmission) HasStepRange(gearCount int, maxStep float64) bool {
	if gearCount > len(t.FinalRatios) {
		return false
	}
startGears:
	for startGear := 0; startGear <= len(t.FinalRatios)-gearCount; startGear++ {
		p, _ := t.FinalRatios[startGear].Ratio.Float64()
		for gear := startGear + 1; gear < len(t.FinalRatios); gear++ {
			f, _ := t.FinalRatios[gear].Ratio.Float64()
			step := f/p - 1.
			if step > maxStep {
				continue startGears
			}
			p = f
		}
		return true
	}
	return false
}

func (t *Transmission) Format(f fmt.State, c rune) {
	for i, gearbox := range t.Gearboxes {
		if i > 0 {
			fmt.Fprint(f, " -> ")
		}
		fmt.Fprint(f, gearbox.String())
	}
	fmt.Fprintf(f, " OVERALL %.f%% MAX %3.0f DEV %3.1f G %d", 100*t.Overall, 100*t.MaxStep, 100*t.StepDeviation,
		len(t.FinalRatios))
	for i, r := range t.FinalRatios {
		fmt.Fprintln(f)
		rf, _ := r.Ratio.Float64()
		bitsFormat := fmt.Sprintf("%%0%db", len(t.Gearboxes))
		bitSettings := fmt.Sprintf(bitsFormat, bits.Reverse(r.Bits)>>(64-len(t.Gearboxes)))
		fmt.Fprintf(f, "%2d %s ", i, bitSettings)
		for j := range t.Gearboxes {
			if j > 0 {
				fmt.Fprint(f, " x")
			}
			var gear *big.Rat
			if r.Bits&(1<<j) == 0 {
				gear = t.Gearboxes[j].Off
			} else {
				gear = t.Gearboxes[j].On
			}
			fmt.Fprintf(f, "%5s", ToRatioString(gear))
		}
		fmt.Fprintf(f, " = %5s ~ %4.2f", ToRatioString(r.Ratio), rf)
		// %5s = %4.2f ", i, bitSettings, ToRatioString(r.ratio), rf)
		if i > 0 {
			fmt.Fprintf(f, " +%3.0f%%", 100*t.Steps[i-1])
		}
		fmt.Fprint(f)
	}
}
