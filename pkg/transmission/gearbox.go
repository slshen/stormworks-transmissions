package transmission

import (
	"fmt"
	"math/big"
)

type Gearbox struct {
	Off *big.Rat
	On  *big.Rat
}

func NewGearbox(mult bool, on, off *big.Rat) *Gearbox {
	g := &Gearbox{
		Off: &big.Rat{},
		On:  &big.Rat{},
	}
	if mult {
		g.Off.Set(off)
		g.On.Set(on)
	} else {
		g.Off.Inv(off)
		g.On.Inv(on)
	}
	return g
}

func (g *Gearbox) String() string {
	return fmt.Sprintf("[%s,%s]", ToRatioString(g.Off), ToRatioString(g.On))
}
