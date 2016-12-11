package phi

import (
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	global sync.RWMutex
	named  = make(map[string]*Monitor)
)

func Register(name string, size int, lf, tf float64) {
	global.Lock()
	named[name] = NewMonitor(size, lf, tf)
	global.Unlock()
}

func FailureOf(name string) float64 {
	global.RLock()
	defer global.RUnlock()
	if m, ok := named[name]; ok {
		return m.Failure()
	}
	return 0
}

func Observed(name string, d time.Duration) {
	global.Lock()
	defer global.Unlock()
	if m, ok := named[name]; ok {
		m.Observed(d)
	}
}

func Dump(name string, w io.Writer) {
	global.RLock()
	defer global.RUnlock()
	if m, ok := named[name]; ok {
		fmt.Fprintf(w, "%v", dumpmonitor(*m.signals))
	}
}
