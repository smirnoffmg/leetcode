package findthelexicographicallysmallestvalidsequence

import (
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func TestValidSequence(t *testing.T) {
	tests := []struct {
		name  string
		word1 string
		word2 string
		want  []int
	}{
		{name: "example 1", word1: "vbcca", word2: "abc", want: []int{0, 1, 2}},
		{name: "example 2", word1: "bacdc", word2: "abc", want: []int{1, 2, 4}},
		{name: "example 3", word1: "aaaaaa", word2: "aaabc", want: []int{}},
		{name: "example 4", word1: "abc", word2: "ab", want: []int{0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validSequence(tt.word1, tt.word2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validSequence(%v, %v) = %v, want %v", tt.word1, tt.word2, got, tt.want)
			}
		})
	}
}

func bruteValidSequence(word1, word2 string) []int {
	n, m := len(word1), len(word2)
	best := []int(nil)
	seq := make([]int, m)

	var rec func(pos, start, diff int)
	rec = func(pos, start, diff int) {
		if pos == m {
			if best == nil || slices.Compare(seq, best) < 0 {
				best = slices.Clone(seq)
			}

			return
		}

		for i := start; i < n; i++ {
			d := diff
			if word1[i] != word2[pos] {
				d++
			}

			if d <= 1 {
				seq[pos] = i
				rec(pos+1, i+1, d)
			}
		}
	}
	rec(0, 0, 0)

	if best == nil {
		return []int{}
	}

	return best
}

func TestValidSequenceAgainstBrute(t *testing.T) {
	rnd := rand.New(rand.NewSource(3302))
	alphabet := "abc"

	for range 3000 {
		n := 1 + rnd.Intn(9)
		m := 1 + rnd.Intn(n)

		w1 := make([]byte, n)
		for i := range w1 {
			w1[i] = alphabet[rnd.Intn(len(alphabet))]
		}

		w2 := make([]byte, m)
		for i := range w2 {
			w2[i] = alphabet[rnd.Intn(len(alphabet))]
		}

		word1, word2 := string(w1), string(w2)

		got, want := validSequence(word1, word2), bruteValidSequence(word1, word2)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("validSequence(%q, %q) = %v, want %v", word1, word2, got, want)
		}
	}
}

// Максимум по ограничениям: len(word1) = 3*10^5, len(word2) = len(word1) - 1.
const benchN = 300_000

func BenchmarkValidSequenceExactMatch(b *testing.B) {
	// word2 — подпоследовательность word1, замена не нужна: только suf и жадный проход
	word1 := strings.Repeat("a", benchN)
	word2 := strings.Repeat("a", benchN-1)

	for b.Loop() {
		validSequence(word1, word2)
	}
}

func BenchmarkValidSequenceChangeAtEnd(b *testing.B) {
	// единственное несовпадение — на последнем символе word2, замена тратится в самом конце
	word1 := strings.Repeat("a", benchN)
	word2 := strings.Repeat("a", benchN-2) + "b"

	for b.Loop() {
		validSequence(word1, word2)
	}
}

func BenchmarkValidSequenceNoAnswer(b *testing.B) {
	// два несовпадения на одну замену: жадность досматривает word1 до конца и возвращает []
	word1 := strings.Repeat("a", benchN)
	word2 := strings.Repeat("a", benchN-3) + "bc"

	for b.Loop() {
		validSequence(word1, word2)
	}
}
