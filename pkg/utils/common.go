package utils

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
