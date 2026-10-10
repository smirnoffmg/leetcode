// Package maximumproductofthreenumbers — решение leetcode 628. Maximum Product of Three Numbers.
package maximumproductofthreenumbers

import "math"

// Кандидатов в ответ ровно два, и оба собираются из пяти крайних значений:
// три наибольших или два наименьших (отрицательных) с наибольшим.
func maximumProduct(nums []int) int {
	hi1, hi2, hi3 := math.MinInt, math.MinInt, math.MinInt
	lo1, lo2 := math.MaxInt, math.MaxInt

	for _, v := range nums {
		switch {
		case v >= hi1:
			hi1, hi2, hi3 = v, hi1, hi2
		case v >= hi2:
			hi2, hi3 = v, hi2
		case v > hi3:
			hi3 = v
		}

		switch {
		case v <= lo1:
			lo1, lo2 = v, lo1
		case v < lo2:
			lo2 = v
		}
	}

	return max(lo1*lo2*hi1, hi1*hi2*hi3)
}
