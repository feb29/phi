package phi

import (
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() { rand.Seed(time.Now().UnixNano()) }

const (
	LF = 0.625 // leveling factor
	TF = 0.625 // trending factor
)

func initMonitor(max int, pm FailureDetector) {
	const stddev, mean = 66.295, 230.855
	for j := 0; j < max; j++ {
		diff := math.Max(100, RandomFloat64(stddev, mean))
		pm.Observed(time.Duration(diff * 1e6))
	}
}

func TestDumpMonitor(t *testing.T) {
	testDumpMonitor(t, 20)
}

func testDumpMonitor(t *testing.T, n int) {
	pm := NewMonitor(n, LF, TF)
	initMonitor(n, pm)

	var (
		est = pm.estimator()
		dur = est.duration(10, 8)
	)

	t.Logf("Duration:%v Estimate:%f", dur, est.estimate(dur))

	pm.Truncate(n / 2)

	est = pm.estimator()
	dur = est.duration(10, 8)
	t.Logf("Duration:%v Estimate:%f", dur, est.estimate(dur))
}

func TestTruncate(t *testing.T) {
	var (
		n  = 1000
		e  = n
		xs = make([]float64, n)
	)
	for n > 1 {
		e /= 2
		n /= 2
		Truncate(&xs, n)
		if len(xs) != e {
			t.Errorf("expected: %d, got: %d", n, len(xs))
		}
	}
}

func TestFailureOf(t *testing.T) {
	const (
		size = 50
		name = "test"
	)

	Register(name, size, 0.76, 0.75)
	for i := 0; i < size; i++ {
		Observed(name, RandomDuration(120*e6, 240*e6))
	}
	Dump(name, os.Stdout)
	t.Logf("%f", FailureOf(name))
}
