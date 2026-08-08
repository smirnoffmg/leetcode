// Package findthelexicographicallysmallestvalidsequence — решение leetcode 3302. Find the Lexicographically Smallest Valid Sequence.
package findthelexicographicallysmallestvalidsequence

func validSequence(word1 string, word2 string) []int {
	n, m := int32(len(word1)), int32(len(word2))

	// last[j] — самый правый индекс word1, на который можно поставить word2[j]
	// при жадном матчинге справа; -1, если word2[j:] не укладывается в word1 вовсе.
	// Отсюда: word2[j:] — подпоследовательность word1[i:] ⟺ i <= last[j].
	last := make([]int32, m+1)
	for j := range m {
		last[j] = -1
	}

	last[m] = n

	for i, j := n-1, m-1; i >= 0 && j >= 0; i-- {
		if word1[i] == word2[j] {
			last[j] = i
			j--
		}
	}

	res := make([]int, 0, m)
	changed := false

	for i, j := int32(0), int32(0); i < n && j < m; i++ {
		switch {
		case word1[i] == word2[j]:
			res = append(res, int(i))
			j++
		// Замену берём как можно левее: индекс i меньше любого следующего, а last
		// гарантирует, что остаток word2 добирается точными совпадениями.
		case !changed && i < last[j+1]:
			res = append(res, int(i))
			changed = true
			j++
		}
	}

	if int32(len(res)) < m {
		return []int{}
	}

	return res
}
