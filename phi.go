// Package phi implements `The Phi Accrual Failure Detector` defined in:
//    http://ddg.jaist.ac.jp/pub/HDY+04.pdf
package phi

import "time"

type (
	Value float64

	Monitor interface {
		Phi() (phi Value)
		Put(time.Duration)
	}
)

var DefaultThreshold = Value(16.0)

var _ Monitor = (*monitor)(nil)

func NewMonitor(size int, lf, tf float64) Monitor {
	return newMonitor(size, lf, tf)
}

type signals []float64

func shrink(xs *signals, to int) int {
	return cutoff(xs, len(*xs)-to)
}

func cutoff(xs *signals, n int) int {
	if n < 0 {
		return 0
	}
	if len(*xs) <= n {
		n = len(*xs)
		*xs = signals(nil)
		return n
	}
	*xs = append(signals(nil), (*xs)[n:]...)
	return n
}
