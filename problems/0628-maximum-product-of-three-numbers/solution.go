// Package maximumproductofthreenumbers — решение leetcode 628. Maximum Product of Three Numbers.
package maximumproductofthreenumbers

import "slices"

func maximumProduct(nums []int) int {
	n := len(nums)

	slices.Sort(nums)

	return max(nums[0]*nums[1]*nums[n-1], nums[n-3]*nums[n-2]*nums[n-1])
}
