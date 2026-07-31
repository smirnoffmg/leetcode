// Package minimumnumberofpushestotypewordii — решение leetcode 3016. Minimum Number of Pushes to Type Word II.
package minimumnumberofpushestotypewordii

import (
	"slices"
)

func minimumPushes(word string) int {
	var freq [26]int

	for i := 0; i < len(word); i++ {
		freq[word[i]-'a']++
	}

	slices.SortFunc(freq[:], func(a, b int) int { return b - a })

	var res int

	// Буква ранга r стоит r/8+1 нажатий, поэтому её счётчик суммируется
	// по одному разу за каждый пройденный ярус из 8 клавиш.
	for start := 0; start < len(freq); start += 8 {
		for _, cnt := range freq[start:] {
			res += cnt
		}
	}

	return res
}
