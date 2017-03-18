// Package phi implements `The Phi Accrual Failure Detector` defined in:
//    http://ddg.jaist.ac.jp/pub/HDY+04.pdf
package phi

import (
	"errors"
	"math"
	"sort"
	"time"
)

// FailureDetector can be used to detect failure of remote server.
type FailureDetector interface {
	Failure() float64
	Observed(time.Duration)
}

var (
	// DefaultThreshold is a default threshold to determine healthy or unhealthy.
	DefaultThreshold = 16.0

	_ FailureDetector = (*Monitor)(nil)
)

// TODO: Mean Squared Error

// Monitor implements `Holt's Linear Method` (DoubleExponentialSmoothing).
// `Holt's Linear Method` is good for non-seasonal data with a trend.
//
//    Level:    L[i]   = a*X[i] + (1−a)*(L[i-1] + T[i-1])
//    Trend:    T[i]   = b*(L[i] − L[i−1]) + (1−b)*T[i−1]
//    Forecast: F[i+1] = L[i] + T[i]
//
// a and b are factor for smoothing observed signals.
type Monitor struct {
	*signals
	sum, sos float64 // sum of square
}

// NewMonitor returns a new Monitor with given factors.
func NewMonitor(size int, level, trend float64) *Monitor {
	return &Monitor{signals: initsig(size, level, trend)}
}

// Failure return the value calculated as:
//
//    -math.Log10(1 - CDF(timeSinceLastHeartbeat)
//
// where CDF is the cumulative distribution function of a normal distribution
// with mean and stddev estimated from historical heartbeat inter-arrival durations
func (m *Monitor) Failure() float64 {
	i := len(m.sigs) - 1
	x := m.levels[i] + m.trends[i]
	e := m.estimator()
	// millisec to time.Duration(nanosec, int64)
	return e.estimate(frommillis(x))
}

func (m *Monitor) estimator() estimator {
	var (
		length   = float64(len(m.levels))
		mean     = m.sum / length
		variance = m.sos/length - mean*mean
		stddev   = math.Sqrt(variance)
	)
	return estimator{mean, stddev}
}

type estimator struct {
	mean, stddev float64
}

func (e estimator) estimate(d time.Duration) float64 {
	return e.calculate(tomillis(d))
}

func (e estimator) calculate(x float64) (val float64) {
	y := (x - e.mean) / e.stddev
	p := math.Exp(-y * (1.5976 + 0.070566*y*y))
	if x > e.mean {
		val = -math.Log10(p / (1.0 + p))
	} else {
		val = -math.Log10(1.0 - 1.0/(1.0+p))
	}
	val = math.Max(0, val)
	return
}

func (e estimator) duration(sec int, threshold float64) time.Duration {
	if threshold <= 0 {
		threshold = DefaultThreshold
	}

	// assume that estimated values are monotonic
	// if a given durations are monotonic.
	j := sort.Search(sec*1000, func(i int) bool {
		return threshold <= e.estimate(fromint(i))
	})
	return fromint(j - 1)
}

// Observed store duration and update internal state.
func (m *Monitor) Observed(d time.Duration) {
	x := m.put(tomillis(d))
	m.sum += x
	m.sos += x * x
}

var errOutOfRange = errors.New("phi: truncation out of range")

// Truncate discards all but the last n signals. t is a truncated size.
func (m *Monitor) Truncate(n int) (t int) {
	if n < 0 || len(m.sigs) <= n {
		panic(errOutOfRange)
	}
	t = truncate(&m.sigs, n)
	m.reset()
	for i := 0; i < len(m.sigs); i++ {
		d := m.putn(i, m.sigs[i])
		m.sum += d
		m.sos += d * d
	}
	return
}

func (m *Monitor) reset() {
	// reset slices and accumlators
	m.signals.reset()
	m.sum, m.sos = 0, 0
}

func truncate(xs *[]float64, size int) (n int) {
	n = len(*xs) - size
	copy((*xs)[:size], (*xs)[n:])
	*xs = (*xs)[:size]
	return
}
