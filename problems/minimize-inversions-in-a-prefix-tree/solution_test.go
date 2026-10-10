package minimizeinversionsinaprefixtree

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestMinInversions(t *testing.T) {
	tests := []struct {
		name   string
		words  []string
		target []int
		want   int
	}{
		{name: "example 1", words: []string{"bad", "bag", "fig", "fed"}, target: []int{3, 1, 0, 2}, want: 2},
		{name: "example 2", words: []string{"ab", "ac", "ba", "bc"}, target: []int{0, 3, 1, 2}, want: 1},

		{name: "single", words: []string{"j"}, target: []int{0}, want: 0},
		// метка на пути: "a" записывается раньше "ab", а в target стоит позже
		{name: "prefix word", words: []string{"a", "ab"}, target: []int{1, 0}, want: 1},
		{name: "chain sorted", words: []string{"a", "ab", "abc"}, target: []int{0, 1, 2}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minInversions(tt.words, tt.target); got != tt.want {
				t.Errorf("minInversions(%v, %v) = %d, want %d", tt.words, tt.target, got, tt.want)
			}
		})
	}
}

// Эталон — по определению: все обходы при всех порядках детей, минимум инверсий.
type bruteNode struct {
	word int
	kids map[byte]*bruteNode
}

func traversals(v *bruteNode) [][]int {
	keys := make([]byte, 0, len(v.kids))
	for c := range v.kids {
		keys = append(keys, c)
	}

	var res [][]int
	permute(keys, 0, func(order []byte) {
		combos := [][]int{{}}
		for _, c := range order {
			var next [][]int
			for _, prefix := range combos {
				for _, tail := range traversals(v.kids[c]) {
					next = append(next, append(slices.Clone(prefix), tail...))
				}
			}
			combos = next
		}
		for _, cmb := range combos {
			if v.word >= 0 {
				cmb = append([]int{v.word}, cmb...)
			}
			res = append(res, cmb)
		}
	})
	return res
}

func permute(a []byte, k int, visit func([]byte)) {
	if k == len(a) {
		visit(a)
		return
	}
	for i := k; i < len(a); i++ {
		a[k], a[i] = a[i], a[k]
		permute(a, k+1, visit)
		a[k], a[i] = a[i], a[k]
	}
}

func bruteMin(words []string, target []int) int {
	root := &bruteNode{word: -1, kids: map[byte]*bruteNode{}}
	for i, w := range words {
		v := root
		for j := range len(w) {
			if v.kids[w[j]] == nil {
				v.kids[w[j]] = &bruteNode{word: -1, kids: map[byte]*bruteNode{}}
			}
			v = v.kids[w[j]]
		}
		v.word = i
	}

	pos := make([]int, len(target))
	for p, w := range target {
		pos[w] = p
	}

	best := -1
	for _, tr := range traversals(root) {
		inv := 0
		for a := range tr {
			for b := a + 1; b < len(tr); b++ {
				if pos[tr[a]] > pos[tr[b]] {
					inv++
				}
			}
		}
		if best < 0 || inv < best {
			best = inv
		}
	}
	return best
}

func randomInput(rng *rand.Rand, n, letters, maxLen int) ([]string, []int) {
	seen := map[string]bool{}
	var words []string
	for len(words) < n {
		b := make([]byte, 1+rng.IntN(maxLen))
		for i := range b {
			b[i] = byte('a' + rng.IntN(letters))
		}
		if w := string(b); !seen[w] {
			seen[w] = true
			words = append(words, w)
		}
	}
	return words, rng.Perm(n)
}

func TestMinInversionsRandom(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 5))

	for range 3000 {
		words, target := randomInput(rng, 1+rng.IntN(7), 2+rng.IntN(3), 3)
		if got, want := minInversions(words, target), bruteMin(words, target); got != want {
			t.Fatalf("words=%q target=%v: got %d, want %d", words, target, got, want)
		}
	}
}

// Много детей у корня: перебор порядков ещё посилен, а DP по маскам работает в полную ширину.
func TestMinInversionsWideRoot(t *testing.T) {
	rng := rand.New(rand.NewPCG(5, 6))

	for range 30 {
		words, target := randomInput(rng, 7, 7, 1)
		if got, want := minInversions(words, target), bruteMin(words, target); got != want {
			t.Fatalf("words=%q target=%v: got %d, want %d", words, target, got, want)
		}
	}
}

func BenchmarkMinInversions(b *testing.B) {
	const n = 10_000

	rng := rand.New(rand.NewPCG(5, 1))
	words, target := randomInput(rng, n, alphabet, 20)

	b.ResetTimer()
	for b.Loop() {
		minInversions(words, target)
	}
}

// Худший случай для DP по маскам: полное 10-арное дерево, у каждого внутреннего узла все 10 детей.
func BenchmarkMinInversionsWide(b *testing.B) {
	const n = 10_000

	words := make([]string, 0, n)
	for i := range n {
		words = append(words, fmt.Sprintf("%04d", i))
	}
	for i := range words {
		w := []byte(words[i])
		for j := range w {
			w[j] += 'a' - '0'
		}
		words[i] = string(w)
	}
	target := rand.New(rand.NewPCG(5, 2)).Perm(n)

	b.ResetTimer()
	for b.Loop() {
		minInversions(words, target)
	}
}
