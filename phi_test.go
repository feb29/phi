package phi_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/section09/phi"
)

func init() { rand.Seed(time.Now().UnixNano()) }

const (
	LF = 0.625 // leveling factor
	TF = 0.625 // trending factor
)

func norm(stddev, mean float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

func initMonitor(max int, pm phi.Monitor) {
	const stddev, mean = 66.295, 230.855
	for j := 0; j < max; j++ {
		diff := math.Max(100, norm(stddev, mean))
		pm.Put(time.Duration(diff * 1e6))
	}
}

func TestMonitor(t *testing.T) {
	const n = 20
	pm := phi.NewMonitor(n, LF, TF)
	initMonitor(n, pm)

	threshold := phi.DefaultThreshold
	d := phi.Duration(pm, threshold)
	e := phi.Estimate(pm, d)
	t.Logf("%v %f", d, e)
	t.Logf("%v\n", phi.Dump(pm))
}
