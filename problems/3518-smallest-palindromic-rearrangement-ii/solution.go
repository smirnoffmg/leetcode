// Package smallestpalindromicrearrangementii — решение leetcode 3518. Smallest Palindromic Rearrangement II.
package smallestpalindromicrearrangementii

// permLimit — потолок для подсчёта перестановок: k <= 1e6, большие значения различать не нужно.
const permLimit = 1_000_001

func smallestPalindrome(s string, k int) string {
	var half [26]int
	for i := 0; i < len(s); i++ {
		half[s[i]-'a']++
	}

	var mid byte
	hasMid := len(s)%2 == 1
	for c := range half {
		if half[c]%2 == 1 {
			mid = byte('a' + c)
		}
		half[c] /= 2
	}

	if countArrangements(half) < k {
		return ""
	}

	m := len(s) / 2
	left := make([]byte, 0, m)

	for len(left) < m && k > 1 {
		for c := range half {
			if half[c] == 0 {
				continue
			}
			half[c]--
			cnt := countArrangements(half)
			if cnt >= k {
				left = append(left, byte('a'+c))
				break
			}
			k -= cnt
			half[c]++
		}
	}
	for c := range half {
		for ; half[c] > 0; half[c]-- {
			left = append(left, byte('a'+c))
		}
	}

	res := make([]byte, len(s))
	copy(res, left)
	if hasMid {
		res[m] = mid
	}
	for i, ch := range left {
		res[len(s)-1-i] = ch
	}

	return string(res)
}

// countArrangements — число различных перестановок мультимножества, обрезанное до permLimit.
func countArrangements(counts [26]int) int {
	total := 0
	for _, c := range counts {
		total += c
	}

	res := 1
	for _, c := range counts {
		if c == 0 {
			continue
		}
		res *= binom(total, c)
		if res >= permLimit {
			return permLimit
		}
		total -= c
	}

	return res
}

// binom — C(n, r), обрезанное до permLimit; промежуточные значения растут монотонно,
// поэтому выход по потолку безопасен.
func binom(n, r int) int {
	if r > n-r {
		r = n - r
	}

	res := 1
	for i := 1; i <= r; i++ {
		res = res * (n - r + i) / i
		if res >= permLimit {
			return permLimit
		}
	}

	return res
}
