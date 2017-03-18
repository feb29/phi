package phi

type signals struct {
	sigs []float64
	// level factor, trend factor
	lf     float64
	levels []float64
	tf     float64
	trends []float64
}

func initsig(size int, lf, tf float64) *signals {
	return &signals{
		sigs: make([]float64, 0, size),
		lf:   bounded(lf),
		tf:   bounded(tf),
	}
}

func bounded(x float64) float64 {
	if x <= 0 {
		x = 0.100
	}
	if x >= 1 {
		x = 0.999
	}
	return x
}

// i: index, x: signals[i]
func (xs *signals) level(i int, x float64) float64 {
	if i == 0 {
		return x
	}
	return xs.lf*x + (1-xs.lf)*(xs.levels[i-1]+xs.trends[i-1])
}

// i: index, x: levels[i]
func (xs *signals) trend(i int, x float64) float64 {
	if i == 0 {
		return 0
	}
	return xs.tf*(x-xs.levels[i-1]) + (1-xs.tf)*xs.trends[i-1]
}

func (xs *signals) put(x float64) float64 {
	xs.sigs = append(xs.sigs, x)
	i := len(xs.sigs) - 1
	return xs.putn(i, x)
}

func (xs *signals) putn(i int, x float64) (d float64) {
	xs.levels = append(xs.levels, xs.level(i, x))
	d = xs.levels[i]
	xs.trends = append(xs.trends, xs.trend(i, d))
	return
}

func (xs *signals) reset() {
	xs.levels = xs.levels[:]
	xs.trends = xs.trends[:]
}
