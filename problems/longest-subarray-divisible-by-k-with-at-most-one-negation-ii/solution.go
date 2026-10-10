// Package longestsubarraydivisiblebykwithatmostonenegationii — решение leetcode Longest Subarray Divisible by K with At Most One Negation II.
package longestsubarraydivisiblebykwithatmostonenegationii

func longestSubarray(nums []int, k int) int {
	mod := func(x int) int { return (x%k + k) % k }

	// first[u] — самый ранний префикс с остатком u; best[c] — самый ранний l, для которого
	// уже есть i >= l с P[l] + 2·nums[i] ≡ c
	first := make([]int, k)
	best := make([]int, k)
	for c := range k {
		first[c] = -1
		best[c] = -1
	}
	first[0] = 0

	// order — остатки префиксов в порядке первого появления; applied[d] — сколько из них
	// уже учтено для сдвига d. Так каждая пара (u, d) обрабатывается ровно один раз.
	order := []int{0}
	applied := make([]int, k)

	ans, p := 0, 0
	for i, x := range nums {
		d := mod(2 * x)
		for _, u := range order[applied[d]:] {
			c := (u + d) % k
			if best[c] < 0 || first[u] < best[c] {
				best[c] = first[u]
			}
		}
		applied[d] = len(order)

		p = mod(p + x)
		r := i + 1

		if first[p] >= 0 {
			ans = max(ans, r-first[p])
		}
		if best[p] >= 0 {
			ans = max(ans, r-best[p])
		}
		if first[p] < 0 {
			first[p] = r
			order = append(order, p)
		}
	}

	return ans
}
