// Package maximumproductoftwodigits — решение leetcode 3536. Maximum Product of Two Digits.
package maximumproductoftwodigits

func maxProduct(n int) int {
	var a, b, next int

	for {
		next = n % 10

		if next > a {
			a, b = next, a
		} else if next > b {
			b = next
		}

		n -= next
		n /= 10

		if n == 0 {
			break
		}
	}
	return a * b
}
