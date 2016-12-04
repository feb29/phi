package phi

// Linear implements `Holt's Linear Method` (DoubleExponentialSmoothing).
// `Holt's Linear Method` is good for non-seasonal data with a trend.
//
//    Level:    L[i]   = a*X[i] + (1−a)*(L[i-1] + T[i-1])
//    Trend:    T[i]   = b*(L[i] − L[i−1]) + (1−b)*T[i−1]
//    Forecast: F[i+1] = L[i] + T[i]
type linear struct {
	signals
	// level factor, trend factor
	lf, tf float64
	levels,
	trends []float64
}

func newLinear(size int, leveling, trending float64) *linear {
	return &linear{
		signals: make(signals, 0, size),
		lf:      bounded(leveling),
		tf:      bounded(trending),
	}
}

func (m *linear) signal(i int) float64 { return m.signals[i] }
func (m *linear) level(i int) float64  { return m.levels[i] }
func (m *linear) trend(i int) float64  { return m.trends[i] }

func (m *linear) put(x float64) float64 {
	m.signals = append(m.signals, x)
	i := len(m.signals) - 1
	return m.putn(i, x)
}

func (m *linear) putn(i int, x float64) (diff float64) {
	switch {
	case i == 0:
		if m.levels == nil || len(m.levels) == 0 {
			m.levels = make([]float64, 1, len(m.signals))
		}
		diff, m.levels[0] = x, x

		if m.trends == nil || len(m.trends) == 0 {
			m.trends = make([]float64, 1, len(m.signals))
		}
		m.trends[0] = 0

	case i >= 1:
		diff = m.lf * m.signal(i)
		diff += (1 - m.lf) * (m.level(i-1) + m.trend(i-1))
		m.levels = append(m.levels, diff)

		trend := m.tf * (m.level(i) - m.level(i-1))
		trend += (1 - m.tf) * m.trend(i-1)
		m.trends = append(m.trends, trend)
	}
	return
}

func (m *linear) shrink(to int) {
	n := shrink(&m.signals, to)
	if n == 0 {
		return
	}
	if to == 0 {
		m.levels = nil
		m.trends = nil
		return
	}
}
