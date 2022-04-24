package rand

import (
	"math/rand"
)

// Float64Range returns a random float64 value in the specified range.
func Float64Range(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
