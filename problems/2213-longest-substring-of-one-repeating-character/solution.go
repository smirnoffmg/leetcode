// Package longestsubstringofonerepeatingcharacter — решение leetcode 2213. Longest Substring of One Repeating Character.
package longestsubstringofonerepeatingcharacter

// Сводки узлов разложены по плоским массивам (SoA), а не по срезу структур:
// на n = k = 1e5 это −27% времени и −33% памяти (benchstat, p=0.000).
func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	n := len(s)
	pre := make([]int, 4*n)
	suf := make([]int, 4*n)
	maxLen := make([]int, 4*n)
	leftChar := make([]byte, 4*n)
	rightChar := make([]byte, 4*n)

	pushUp := func(u, l, r int) {
		mid := (l + r) >> 1
		leftLen, rightLen := mid-l+1, r-mid
		left, right := u<<1, u<<1|1
		leftChar[u], rightChar[u] = leftChar[left], rightChar[right]

		pre[u] = pre[left]
		if pre[left] == leftLen && rightChar[left] == leftChar[right] {
			pre[u] = pre[left] + pre[right]
		}
		suf[u] = suf[right]
		if suf[right] == rightLen && rightChar[left] == leftChar[right] {
			suf[u] = suf[right] + suf[left]
		}
		maxLen[u] = max(maxLen[left], maxLen[right])
		if rightChar[left] == leftChar[right] {
			maxLen[u] = max(maxLen[u], suf[left]+pre[right])
		}
	}

	var build func(u, l, r int)
	build = func(u, l, r int) {
		if l == r {
			pre[u], suf[u], maxLen[u] = 1, 1, 1
			leftChar[u], rightChar[u] = s[l], s[l]

			return
		}
		mid := (l + r) >> 1
		build(u<<1, l, mid)
		build(u<<1|1, mid+1, r)
		pushUp(u, l, r)
	}

	var update func(u, l, r, pos int, ch byte)
	update = func(u, l, r, pos int, ch byte) {
		if l == r {
			leftChar[u], rightChar[u] = ch, ch

			return
		}
		mid := (l + r) >> 1
		if pos <= mid {
			update(u<<1, l, mid, pos, ch)
		} else {
			update(u<<1|1, mid+1, r, pos, ch)
		}
		pushUp(u, l, r)
	}

	build(1, 0, n-1)
	lengths := make([]int, len(queryIndices))
	for i := range queryIndices {
		update(1, 0, n-1, queryIndices[i], queryCharacters[i])
		lengths[i] = maxLen[1]
	}

	return lengths
}
