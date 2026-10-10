// Package minimizeinversionsinaprefixtree — решение leetcode Minimize Inversions in a Prefix Tree.
package minimizeinversionsinaprefixtree

import (
	"math"
	"math/bits"
)

const alphabet = 10

type trie struct {
	next [][alphabet]int32
	mark []int32 // позиция слова в target или -1
}

func (t *trie) add() int32 {
	t.next = append(t.next, [alphabet]int32{})
	t.mark = append(t.mark, -1)
	return int32(len(t.mark) - 1)
}

func minInversions(words []string, target []int) int {
	pos := make([]int32, len(target))
	for p, w := range target {
		pos[w] = int32(p)
	}

	t := &trie{}
	t.add()
	for i, w := range words {
		v := int32(0)
		for j := range len(w) {
			c := w[j] - 'a'
			if t.next[v][c] == 0 {
				u := t.add()
				t.next[v][c] = u
			}
			v = t.next[v][c]
		}
		t.mark[v] = pos[i]
	}

	total := 0
	t.solve(0, &total)
	return total
}

// solve возвращает отсортированные позиции поддерева v и копит в total
// инверсии, которые решаются внутри v.
func (t *trie) solve(v int32, total *int) []int32 {
	var kids [][]int32
	for _, u := range t.next[v] {
		if u != 0 {
			kids = append(kids, t.solve(u, total))
		}
	}

	// Порядок детей выбирается перебором подмножеств, а не сортировкой:
	// попарное сравнение cross[a][b] против cross[b][a] бывает нетранзитивным.
	*total += bestOrder(kids)

	merged := mergeAll(kids)
	if p := t.mark[v]; p >= 0 {
		// Метка записывается раньше всех потомков: инверсию дают потомки левее в target.
		*total += lowerCount(merged, p)
		merged = insertSorted(merged, p)
	}
	return merged
}

func bestOrder(kids [][]int32) int {
	c := len(kids)
	if c < 2 {
		return 0
	}

	// cross[a][b] — инверсии между a и b, если a обходится раньше b.
	var cross [alphabet][alphabet]int
	for a := range c {
		for b := range c {
			if a != b {
				cross[a][b] = pairsGreater(kids[a], kids[b])
			}
		}
	}

	// after[b][mask] — цена поставить b после всех детей из mask.
	full := 1 << c
	after := make([][]int, c)
	for b := range c {
		after[b] = make([]int, full)
		for m := 1; m < full; m++ {
			low := m & -m
			after[b][m] = after[b][m^low] + cross[bits.TrailingZeros(uint(low))][b]
		}
	}

	dp := make([]int, full)
	for m := 1; m < full; m++ {
		dp[m] = math.MaxInt
	}
	for m := range full - 1 {
		for b := range c {
			if m&(1<<b) == 0 {
				dp[m|1<<b] = min(dp[m|1<<b], dp[m]+after[b][m])
			}
		}
	}

	return dp[full-1]
}

// pairsGreater считает пары x из a, y из b с x > y; оба списка отсортированы.
func pairsGreater(a, b []int32) int {
	res, j := 0, 0
	for _, x := range a {
		for j < len(b) && b[j] < x {
			j++
		}
		res += j
	}
	return res
}

func lowerCount(s []int32, p int32) int {
	n := 0
	for _, x := range s {
		if x >= p {
			break
		}
		n++
	}
	return n
}

func insertSorted(s []int32, p int32) []int32 {
	i := lowerCount(s, p)
	s = append(s, 0)
	copy(s[i+1:], s[i:])
	s[i] = p
	return s
}

func mergeAll(lists [][]int32) []int32 {
	for len(lists) > 1 {
		var next [][]int32
		for i := 0; i+1 < len(lists); i += 2 {
			next = append(next, merge(lists[i], lists[i+1]))
		}
		if len(lists)%2 == 1 {
			next = append(next, lists[len(lists)-1])
		}
		lists = next
	}
	if len(lists) == 0 {
		return nil
	}
	return lists[0]
}

func merge(a, b []int32) []int32 {
	out := make([]int32, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			out = append(out, a[i])
			i++
		} else {
			out = append(out, b[j])
			j++
		}
	}
	out = append(out, a[i:]...)
	return append(out, b[j:]...)
}
