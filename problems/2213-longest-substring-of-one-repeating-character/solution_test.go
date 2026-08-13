package longestsubstringofonerepeatingcharacter

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestLongestRepeating(t *testing.T) {
	tests := []struct {
		name            string
		s               string
		queryCharacters string
		queryIndices    []int
		want            []int
	}{
		{name: "example 1", s: "babacc", queryCharacters: "bcb", queryIndices: []int{1, 3, 3}, want: []int{3, 3, 4}},
		{name: "example 2", s: "abyzz", queryCharacters: "aa", queryIndices: []int{2, 1}, want: []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestRepeating(tt.s, tt.queryCharacters, tt.queryIndices)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("longestRepeating(%v, %v, %v) = %v, want %v", tt.s, tt.queryCharacters, tt.queryIndices, got, tt.want)
			}
		})
	}
}

func bruteLongestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	chars := []byte(s)
	lengths := make([]int, len(queryIndices))
	for i, pos := range queryIndices {
		chars[pos] = queryCharacters[i]

		best, run := 1, 1
		for j := 1; j < len(chars); j++ {
			if chars[j] == chars[j-1] {
				run++
			} else {
				run = 1
			}
			best = max(best, run)
		}
		lengths[i] = best
	}

	return lengths
}

func TestLongestRepeatingAgainstBrute(t *testing.T) {
	rnd := rand.New(rand.NewSource(2213))
	alphabet := "abc"

	for range 3000 {
		n := 1 + rnd.Intn(9)
		k := 1 + rnd.Intn(9)

		s := make([]byte, n)
		for i := range s {
			s[i] = alphabet[rnd.Intn(len(alphabet))]
		}

		qc := make([]byte, k)
		qi := make([]int, k)
		for i := range k {
			qc[i] = alphabet[rnd.Intn(len(alphabet))]
			qi[i] = rnd.Intn(n)
		}

		got := longestRepeating(string(s), string(qc), qi)
		want := bruteLongestRepeating(string(s), string(qc), qi)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("longestRepeating(%q, %q, %v) = %v, want %v", s, qc, qi, got, want)
		}
	}
}

func BenchmarkLongestRepeating(b *testing.B) {
	const n, k = 100_000, 100_000

	rnd := rand.New(rand.NewSource(2213))
	s := make([]byte, n)
	for i := range s {
		s[i] = 'a' + byte(rnd.Intn(26))
	}

	qc := make([]byte, k)
	qi := make([]int, k)
	for i := range k {
		qc[i] = 'a' + byte(rnd.Intn(26))
		qi[i] = rnd.Intn(n)
	}

	b.ResetTimer()

	for b.Loop() {
		longestRepeating(string(s), string(qc), qi)
	}
}
