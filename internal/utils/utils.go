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