package utils

import (
	"strconv"
	"strings"
)

func MaxFloat64(a []float64) float64 {
	max := a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func MinFloat64(a []float64) float64 {
	min := a[0]
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

func StringToBool(str string) bool {
	if strings.ToLower(str) == "true" {
		return true
	} else {
		return false
	}
}

func StringToInt(str string) int {
	if str == "" {
		return 0
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}
