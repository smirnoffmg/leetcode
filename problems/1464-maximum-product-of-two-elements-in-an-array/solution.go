// Package maximumproductoftwoelementsinanarray — решение leetcode 1464. Maximum Product of Two Elements in an Array.
package maximumproductoftwoelementsinanarray

func maxProduct(nums []int) int {
	a, b := 1, 1

	for _, n := range nums {
		if n >= a {
			a, b = n, a
		} else if n > b {
			b = n
		}
	}

	return (a - 1) * (b - 1)

}
