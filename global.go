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

// Register initialize Monitor with a specified name.
// size, lf, tf are the same arguments with NewMonitor.
func Register(name string, size int, lf, tf float64) {
	global.Lock()
	named[name] = NewMonitor(size, lf, tf)
	global.Unlock()
}

// FailureOf returns Failure of monitor with given name.
func FailureOf(name string) float64 {
	global.RLock()
	defer global.RUnlock()
	if m, ok := named[name]; ok {
		return m.Failure()
	}
	return 0
}

// Observed collect d, for later Failure estimation.
func Observed(name string, d time.Duration) {
	global.Lock()
	defer global.Unlock()
	if m, ok := named[name]; ok {
		m.Observed(d)
	}
}

// Dump dump a detailed information of monitor.
func Dump(name string, w io.Writer) {
	global.RLock()
	defer global.RUnlock()
	if m, ok := named[name]; ok {
		fmt.Fprintf(w, "%v", dumpmonitor(*m.signals))
	}
}
