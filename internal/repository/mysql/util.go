package mysql

import "math"

func roundDown(value float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(value*shift) / shift
}

func validateNAB(totalBalance, totalUnits float64) float64 {
	if totalUnits <= 0 {
		return 1.0
	}

	return roundDown(totalBalance/totalUnits, 4)
}
