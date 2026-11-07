package countvalidsubstringssubsequence

import (
	"math/rand/v2"
	"strings"
	"testing"
)

const maxLength = 100_000

func TestCountValidSubstrings(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int64
	}{
		{name: "example 1", s: "abac", want: 9},
		{name: "example 2", s: "abc", want: 5},
		{name: "example 3", s: "ac", want: 3},

		{name: "single char", s: "c", want: 1},
		{name: "reversed letters", s: "cba", want: 6},
		{name: "subsequence with gaps", s: "aabbcc", want: 17},
		// a после b не помогает: для старта важна a левее b
		{name: "a after b", s: "bac", want: 6},
		{name: "repeated abc", s: "abcabc", want: 14},
		{name: "max length without c", s: strings.Repeat("ab", maxLength/2), want: 5_000_050_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countValidSubstrings(tt.s); got != tt.want {
				t.Errorf("countValidSubstrings(%.32q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestCountValidSubstringsMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2026, 1109))

	for range 20_000 {
		s := randomString(rng, 1+rng.IntN(12))
		want := bruteForceCount(s)

		for _, impl := range implementations() {
			if got := impl.fn(s); got != want {
				t.Fatalf("%s(%q) = %v, want %v", impl.name, s, got, want)
			}
		}
	}
}

func BenchmarkCountValidSubstrings(b *testing.B) {
	rng := rand.New(rand.NewPCG(2026, 3))

	inputs := []struct{ name, s string }{
		{name: "shape=random", s: randomString(rng, maxLength)},
		{name: "shape=all-abc", s: strings.Repeat("abc", maxLength/3)},
	}

	for _, input := range inputs {
		for _, impl := range implementations() {
			b.Run(input.name+"/"+impl.name, func(b *testing.B) {
				for b.Loop() {
					impl.fn(input.s)
				}
			})
		}
	}
}

func bruteForceCount(s string) int64 {
	var count int64

	for left := range s {
		for right := left + 1; right <= len(s); right++ {
			if !containsABCSubsequence(s[left:right]) {
				count++
			}
		}
	}

	return count
}

func containsABCSubsequence(s string) bool {
	matched := 0

	for i := range s {
		if matched < 3 && s[i] == "abc"[matched] {
			matched++
		}
	}

	return matched == 3
}

type implementation struct {
	name string
	fn   func(string) int64
}

func implementations() []implementation {
	return []implementation{
		{name: "impl=switch", fn: countSwitch},
		{name: "impl=table", fn: countValidSubstrings},
	}
}

func randomString(rng *rand.Rand, n int) string {
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = 'a' + byte(rng.IntN(3))
	}

	return string(letters)
}

func countSwitch(s string) int64 {
	var total int64

	lastA, lastAB, lastABC := -1, -1, -1

	for right := range s {
		switch s[right] {
		case 'a':
			lastA = right
		case 'b':
			lastAB = lastA
		case 'c':
			lastABC = lastAB
		}

		total += int64(right - lastABC)
	}

	return total
}
