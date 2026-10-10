// Package longestsubarraydivisiblebykwithatmostonenegationi — решение leetcode Longest Subarray Divisible by K with At Most One Negation I.
package longestsubarraydivisiblebykwithatmostonenegationi

func longestSubarray(nums []int, k int) int {
	mod := func(x int) int { return (x%k + k) % k }

	// stamp[c] == l+1 — остаток c встретился среди 2x при текущем l; так массив не надо чистить на каждом l
	stamp := make([]int, k)
	best := 0

	for l := range nums {
		sum := 0
		for r := l; r < len(nums); r++ {
			sum += nums[r]
			stamp[mod(2*nums[r])] = l + 1

			s := mod(sum)
			if (s == 0 || stamp[s] == l+1) && r-l+1 > best {
				best = r - l + 1
			}
		}
	}

	return best
}
