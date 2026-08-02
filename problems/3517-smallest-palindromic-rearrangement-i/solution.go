// Package smallestpalindromicrearrangementi — решение leetcode 3517. Smallest Palindromic Rearrangement I.
package smallestpalindromicrearrangementi

import "slices"

func smallestPalindrome(s string) string {
	partition := len(s) / 2

	res := []byte(s)
	slices.Sort(res[:partition])

	for i := range partition {
		res[len(s)-1-i] = res[i]
	}

	return string(res)
}
