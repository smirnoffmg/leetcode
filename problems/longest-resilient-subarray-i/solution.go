// Package longestresilientsubarrayi — решение leetcode Longest Resilient Subarray I.
package longestresilientsubarrayi

// Подмассив устойчив, когда все его элементы дают один остаток r по модулю k
// и (L−1)·r ≡ 0 (mod k), то есть L−1 кратно k/gcd(r, k). Поэтому в каждой
// максимальной серии с одним остатком берётся самая длинная допустимая длина.
func resilientSubarray(nums []int, k int) int {
	best := 0
	for start := 0; start < len(nums); {
		r := nums[start] % k
		end := start + 1
		for end < len(nums) && nums[end]%k == r {
			end++
		}

		step := k / gcd(r, k)
		best = max(best, 1+(end-start-1)/step*step)
		start = end
	}

	return best
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
