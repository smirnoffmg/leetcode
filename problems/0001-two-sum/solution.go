// Package twosum — решение leetcode 1. Two Sum.
package twosum

func twoSum(nums []int, target int) []int {
	hm := make(map[int]int, len(nums))

	for i := range nums {
		if _, ok := hm[nums[i]]; !ok {
			hm[target-nums[i]] = i
		} else {
			return []int{hm[nums[i]], i}
		}
	}

	return []int{}
}
