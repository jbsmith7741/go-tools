package gtools

import (
	"math"
)

// PreciseFloat truncates a float to the specified number of digits
func PreciseFloat(f float64, digits int) float64 {
	pow := math.Pow(10.0, float64(digits))
	return math.Trunc(pow*f) / pow
}
