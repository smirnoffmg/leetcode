// Package countvalidsubstrings — решение контестной задачи LeetCode о подстроках без "abc".
package countvalidsubstrings

func countValidSubstrings(s string) int64 {
	var total int64

	last := -1

	for right := range s {
		if right >= 2 && s[right-2:right+1] == "abc" {
			last = right - 2
		}

		total += int64(right - last)
	}

	return total
}
