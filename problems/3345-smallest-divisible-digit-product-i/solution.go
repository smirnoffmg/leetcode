// Package smallestdivisibledigitproducti — решение leetcode 3345. Smallest Divisible Digit Product I.
package smallestdivisibledigitproducti

func digProduct(n int) (res int) {
	res = 1

	for n > 0 {
		res *= (n % 10)
		n /= 10
	}

	return
}

func smallestNumber(n int, t int) int {
	for {
		if digProduct(n)%t == 0 {
			return n
		}

		n++
	}
}
