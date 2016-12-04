package phi

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

func Dump(m Monitor) fmt.Stringer {
	switch x := m.(type) {
	case *monitor:
		return dumpLinear{x.linear}
	}
	return nil
}

type (
	dumpLinear struct{ *linear }
)

func (m dumpLinear) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "=== DUMP %s\n",
		fmt.Sprintf("phi.Monitor(%f %f)", m.lf, m.tf),
	)
	if m.trends == nil || len(m.trends) <= 0 {
		return w.String()
	}

	func(max int) {
		t := tabwriter.NewWriter(w, 0, 8, 0, '\t', 0)
		defer t.Flush()

		fmt.Fprintf(t, "\tSignal(i)\tLevel(i)\tTrend(i)\tForecast(i)\n")
		for i := 0; i < max; i++ {
			level := m.level(i)
			trend := m.trend(i)
			fmt.Fprintf(t, "% 5d\t%+f\t%+f\t%+f\t% f\n",
				i, m.signals[i], level, trend, level+trend,
			)
		}
	}(len(m.signals))
	return w.String()
}
