package utils

import (
	"math"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// RoundDown rounds down a float value to the specified decimal places
func RoundDown(value float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(value*shift) / shift
}

// ValidateNAB calculates Net Asset Value (NAB)
func ValidateNAB(totalBalance, totalUnits float64) float64 {
	if totalUnits <= 0 {
		return 1.0
	}

	return RoundDown(totalBalance/totalUnits, 4)
}
