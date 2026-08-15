// Package longestsubsequencewithnonzerobitwisexor — решение leetcode 3702. Longest Subsequence With Non-Zero Bitwise XOR.
package longestsubsequencewithnonzerobitwisexor

func longestSubsequence(nums []int) int {
	acc := 0
	hasNonZero := false

	for _, v := range nums {
		acc ^= v
		hasNonZero = hasNonZero || v != 0
	}

	if acc != 0 {
		return len(nums)
	}

	if hasNonZero {
		return len(nums) - 1
	}

	return 0
}
