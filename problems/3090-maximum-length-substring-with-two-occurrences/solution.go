// Package maximumlengthsubstringwithtwooccurrences — решение leetcode 3090. Maximum Length Substring With Two Occurrences.
package maximumlengthsubstringwithtwooccurrences

func allLessThan(arr []int16, target int16) bool {
	for _, val := range arr {
		if val >= target {
			return false
		}
	}

	return true
}

func maximumLengthSubstring(s string) int {
	buf := make([]int16, 26)
	l, p, maxLen, currentLen := 0, 0, 1, 1

	buf[s[0]-'a']++

	for p < len(s)-1 {
		p++
		buf[s[p]-'a']++
		currentLen++

		for !allLessThan(buf, 3) {
			buf[s[l]-'a']--
			l++
			currentLen--
		}

		maxLen = max(maxLen, currentLen)

	}

	return maxLen
}
