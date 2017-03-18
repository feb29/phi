package phi

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

const (
	e3 = 1e3
	e6 = 1e6
)

func tomillis(d time.Duration) float64   { return d.Seconds() * e3 }
func frommillis(d float64) time.Duration { return time.Duration(d * e6) }
func fromint(n int) time.Duration        { return time.Duration(n) * time.Millisecond }

type dumpmonitor signals

func (m dumpmonitor) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "=== DUMP %s\n", fmt.Sprintf("LEVEL=%f TREND=%f", m.lf, m.tf))
	if m.trends == nil || len(m.trends) <= 0 {
		return w.String()
	}

	func(max int) {
		t := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
		defer t.Flush()

		fmt.Fprintf(t, "\tSignal(i)\tLevel(i)\tTrend(i)\tForecast(i)\t\n")
		for i := 0; i < max; i++ {
			level := m.levels[i]
			trend := m.trends[i]
			fmt.Fprintf(t, "% 5d\t% f\t% f\t%+f\t% f\t\n",
				i, m.sigs[i], level, trend, level+trend,
			)
		}
	}(len(m.sigs))
	return w.String()
}
