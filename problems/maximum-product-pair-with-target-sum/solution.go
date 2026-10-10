// Package maximumproductpairwithtargetsum — решение leetcode Maximum Product Pair With Target Sum.
package maximumproductpairwithtargetsum

const (
	minVal = -100
	maxVal = 100
)

func maxProductPair(nums []int, target int) []int {
	var pos [maxVal - minVal + 1]int
	for i := range pos {
		pos[i] = -1
	}
	for i, v := range nums {
		pos[v-minVal] = i
	}

	bi, bj, bestProduct := -1, -1, 0
	for i, a := range nums {
		b := target - a
		if b >= a || b < minVal || pos[b-minVal] < 0 {
			continue
		}
		if p := a * b; bi < 0 || p > bestProduct {
			bi, bj, bestProduct = i, pos[b-minVal], p
		}
	}

	return []int{bi, bj}
}
