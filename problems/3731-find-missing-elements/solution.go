// Package findmissingelements — решение leetcode 3731. Find Missing Elements.
package findmissingelements

func findMissingElements(nums []int) []int {
	var seen [101]bool
	res := []int{}
	minNum := 101
	maxNum := 0

	for _, n := range nums {
		if n > maxNum {
			maxNum = n
		}

		if n < minNum {
			minNum = n
		}

		seen[n] = true
	}

	for i := minNum; i < maxNum+1; i++ {
		if !seen[i] {
			res = append(res, i)
		}

	}

	return res

}
