package transmission

import (
	"encoding/gob"
	"fmt"
	"math/big"
	"os"
	"time"
)

func GenerateCache(gearboxCount int, clear bool) (result []*Transmission) {
	path := fmt.Sprintf("transmissions-%d.gob", gearboxCount)
	if !clear {
		f, err := os.Open(path)
		if err == nil {
			err = func() error {
				defer f.Close()
				dec := gob.NewDecoder(f)
				return dec.Decode(&result)
			}()
			if err == nil {
				return
			}
			if err != nil {
				fmt.Printf("Cache file %s is invalid: %s\n", path, err)
			}
		}
	}
	result = Generate(gearboxCount)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("Could not cache results into %s: %s\n", path, err)
		return
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(result)
	if err != nil {
		fmt.Printf("Could not write results into %s: %s\n", path, err)
	}
	return
}

func Generate(gearboxCount int) []*Transmission {
	t0 := time.Now()
	var counter = &Counter{}
	var n int
	transmissionsMap := map[string]*Transmission{}
	transmissions := []*Transmission{}

	for i := 0; i < gearboxCount; i++ {
		counter.AddDigit(len(gearRatios))
		counter.AddDigit(len(gearRatios))
		counter.AddDigit(2)
	}
	for {
		gearboxes := make([]*Gearbox, gearboxCount)
		for i := 0; i < gearboxCount; i++ {
			var (
				mult bool
				on   *big.Rat
				off  *big.Rat
			)
			on = gearRatios[counter.digits[3*i+0]]
			off = gearRatios[counter.digits[3*i+1]]
			if counter.digits[3*i+2] == 1 {
				mult = true
			}
			gearboxes[i] = NewGearbox(mult, on, off)
		}
		t := NewTransmission(gearboxes)
		key := t.Key()
		if transmissionsMap[key] == nil {
			transmissionsMap[key] = t
			transmissions = append(transmissions, t)
			n++
		}
		if !counter.Increment() {
			break
		}
	}
	fmt.Printf("Generated %d unique transmission with %d gearboxes in %s\n", n, gearboxCount, time.Since(t0).Round(time.Millisecond))
	return transmissions
}
