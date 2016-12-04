package phi

import (
	"math"
	"time"
)

// TODO: Mean Squared Error
type monitor struct {
	*linear
	sum, sos float64 // sum of square
}

func newMonitor(size int, lf, tf float64) *monitor {
	return &monitor{linear: newLinear(size, lf, tf)}
}

func (m *monitor) Put(d time.Duration) {
	x := m.put(milliSeconds(d))
	m.sum += x
	m.sos += x * x
}

func (m *monitor) Phi() Value {
	i := len(m.signals) - 1
	x := m.level(i) + m.trend(i)
	return m.estimate(x)
}

func (m *monitor) estimate(x float64) Value {
	var (
		length   = float64(len(m.levels))
		mean     = m.sum / length
		variance = m.sos/length - mean*mean
		stddev   = math.Sqrt(variance)
	)
	return Value(calculate(x, mean, stddev))
}

// The value of phi is calculated as:
//
//    -math.Log10(1 - CDF(timeSinceLastHeartbeat)
//
// where CDF is the cumulative distribution function of a normal distribution
// with mean and stddev estimated from historical heartbeat inter-arrival durations
func calculate(x, mean, stddev float64) (val float64) {
	y := (x - mean) / stddev
	e := math.Exp(-y * (1.5976 + 0.070566*y*y))
	if x > mean {
		val = math.Max(0, -math.Log10(e/(1.0+e)))
	} else {
		val = math.Max(0, -math.Log10(1.0-1.0/(1.0+e)))
	}
	return
}

func (x *monitor) shrink(to int) {
	n := shrink(&x.signals, to)
	if n == 0 {
		return
	}
	x.sum, x.sos = 0, 0
	if to == 0 {
		x.levels, x.trends = nil, nil
		return
	}
	for i := 0; i < len(x.signals); i++ {
		x.putn(i, x.signals[i])
	}
}
