package phi

import (
	"fmt"
	"time"
)

func bounded(x float64) float64 {
	if x <= 0 {
		x = 0.100
	}
	if x >= 1 {
		x = 0.999
	}
	return x
}

const (
	e3 = 1e3
	e6 = 1e6
)

func milliSeconds(d time.Duration) float64 { return d.Seconds() * e3 }

// panic if pred return false
func bugon(text string, pred bool) {
	if pred {
		return
	}
	panic(fmt.Sprintf("BUGON: %s", text))
}
