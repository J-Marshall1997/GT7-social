package utils

import "math"

// Why would math.Min only work for floats...
func MinInt(x ...int) (int) {
	min := math.MaxInt
	for _, num := range(x){
		if num < min {
			min = num
		}
	}
	return min
}

func MaxInt(x ...int) (int) {
	max := math.MinInt
	for _, num := range(x){
		if num > max {
			max = num
		}
	}
	return max
}