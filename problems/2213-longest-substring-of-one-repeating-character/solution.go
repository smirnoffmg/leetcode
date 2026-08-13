// Package longestsubstringofonerepeatingcharacter — решение leetcode 2213. Longest Substring of One Repeating Character.
package longestsubstringofonerepeatingcharacter

type segNode struct {
	size, pref, suf, best int
	left, right           byte
}

func mergeNodes(a, b segNode) segNode {
	c := segNode{
		size: a.size + b.size,
		pref: a.pref,
		suf:  b.suf,
		best: max(a.best, b.best),
		left: a.left, right: b.right,
	}
	if a.right == b.left {
		c.best = max(c.best, a.suf+b.pref)
		if a.pref == a.size {
			c.pref = a.size + b.pref
		}
		if b.suf == b.size {
			c.suf = b.size + a.suf
		}
	}

	return c
}

type segTree struct {
	nodes []segNode
	chars []byte
}

func (t *segTree) build(v, lo, hi int) {
	if lo == hi {
		ch := t.chars[lo]
		t.nodes[v] = segNode{size: 1, pref: 1, suf: 1, best: 1, left: ch, right: ch}

		return
	}

	mid := (lo + hi) / 2
	t.build(2*v, lo, mid)
	t.build(2*v+1, mid+1, hi)
	t.nodes[v] = mergeNodes(t.nodes[2*v], t.nodes[2*v+1])
}

func (t *segTree) update(v, lo, hi, pos int) {
	if lo == hi {
		ch := t.chars[pos]
		t.nodes[v].left, t.nodes[v].right = ch, ch

		return
	}

	mid := (lo + hi) / 2
	if pos <= mid {
		t.update(2*v, lo, mid, pos)
	} else {
		t.update(2*v+1, mid+1, hi, pos)
	}
	t.nodes[v] = mergeNodes(t.nodes[2*v], t.nodes[2*v+1])
}

func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	n := len(s)
	t := &segTree{nodes: make([]segNode, 4*n), chars: []byte(s)}
	t.build(1, 0, n-1)

	lengths := make([]int, len(queryIndices))
	for i, pos := range queryIndices {
		t.chars[pos] = queryCharacters[i]
		t.update(1, 0, n-1, pos)
		lengths[i] = t.nodes[1].best
	}

	return lengths
}
