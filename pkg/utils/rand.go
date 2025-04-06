package utils

import (
	"math"
	"math/rand"
	"time"
)

// Seed random generator.
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandInRange(min, max float64) float64 {
	// Generate random price in the limited range
	val := min + r.Float64()*(max-min)
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
