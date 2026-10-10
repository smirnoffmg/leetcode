// Package maximumlengthsubstringwithtwooccurrences — решение leetcode 3090. Maximum Length Substring With Two Occurrences.
package maximumlengthsubstringwithtwooccurrences

func maximumLengthSubstring(s string) int {
	var count [26]int
	best, left := 0, 0

	for right := range len(s) {
		c := s[right] - 'a'
		count[c]++
		// нарушить лимит могла только буква c, поэтому достаточно
		// сжимать окно, пока именно её счётчик не вернётся к двум
		for count[c] > 2 {
			count[s[left]-'a']--
			left++
		}
		best = max(best, right-left+1)
	}

	return best
}
