package phi

import (
	"math/rand"
	"time"
)

func RandomFloat64(stddev, mean float64) float64 {
	return rand.NormFloat64()*stddev + mean
}

func RandomDuration(stddev, mean float64) time.Duration {
	return time.Duration(RandomFloat64(stddev, mean))
}

var Truncate = truncate
