package utils

import (
	"math"
	"math/rand"
	"time"
)

// Seed random generator.
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// Generate random number in the full open range (min, max)
func RandInRange(min, max float64) float64 {
again:
	val := min + r.Float64()*(max-min)
	if val == min || val == max {
		goto again
	}
	return val
}

// RandGaussianInRange generates a random number with Gaussian distribution within a given range (min, max)
func RandGaussianInRange(min, max, stddev float64) float64 {
	mean := (min + max) / 2
	// Generate a random number based on a Gaussian distribution (mean, stddev)
	val := mean + r.NormFloat64()*stddev
	// Ensure the value is within the specified range
	val = math.Max(val, min)
	val = math.Min(val, max)
	return val
}
