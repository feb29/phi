package phi

import "time"

func Estimate(m Monitor, x time.Duration) Value {
	switch t := m.(type) {
	case *monitor:
		return t.estimate(milliSeconds(x))
	}
	return 0
}

func Duration(pm Monitor, threshold Value) time.Duration {
	switch t := pm.(type) {
	case *monitor:
		d := 10 * time.Millisecond
		for t.estimate(milliSeconds(d)) < threshold {
			d += 10 * time.Millisecond
		}
		return d
	}
	return 0
}
