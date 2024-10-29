package util

import (
	"math"
	"math/rand"
)

func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// Returns a random real in [0,1)
func RandomFloat() float64 {
	return rand.Float64()
}

// Returns a random float in [min, max)
func RandomFloatFromRange(min, max float64) float64 {
	return min + (max-min)*RandomFloat()
}
