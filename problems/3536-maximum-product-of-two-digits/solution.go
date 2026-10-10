// Package maximumproductoftwodigits — решение leetcode 3536. Maximum Product of Two Digits.
package maximumproductoftwodigits

func maxProduct(n int) int {
	var a, b int

	for ; n > 0; n /= 10 {
		d := n % 10

		if d > a {
			a, b = d, a
		} else if d > b {
			b = d
		}
	}

	return a * b
}
