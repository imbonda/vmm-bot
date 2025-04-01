package utils

import (
	"math/rand"
	"time"
)

func RandInRange(min, max float64) float64 {
	// Seed random generator.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generate random price in the limited range
	val := min + r.Float64()*(max-min)
	return val
}
