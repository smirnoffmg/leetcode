// Package lengthoflongestsubarraywithatmostkfrequency — решение leetcode 2958. Length of Longest Subarray With at Most K Frequency.
package lengthoflongestsubarraywithatmostkfrequency

func maxSubarrayLength(nums []int, k int) int {
	freq := make(map[int]int)
	best, left := 0, 0

	for right, x := range nums {
		freq[x]++
		// нарушить лимит могла только частота x, поэтому достаточно
		// сжимать окно, пока именно она не вернётся к k
		for freq[x] > k {
			freq[nums[left]]--
			left++
		}
		best = max(best, right-left+1)
	}

	return best
}
