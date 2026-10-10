package smallestpalindromicrearrangementi

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestSmallestPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{name: "example 1", s: "z", want: "z"},
		{name: "example 2", s: "babab", want: "abbba"},
		{name: "example 3", s: "daccad", want: "acddca"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smallestPalindrome(tt.s)
			if got != tt.want {
				t.Errorf("smallestPalindrome(%v) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func BenchmarkSmallestPalindrome(b *testing.B) {
	rng := rand.New(rand.NewPCG(3517, 1))

	left := make([]byte, 50_000)
	for i := range left {
		left[i] = byte('a' + rng.IntN(26))
	}

	right := slices.Clone(left)
	slices.Reverse(right)
	s := string(left) + string(right)

	b.ResetTimer()

	for b.Loop() {
		smallestPalindrome(s)
	}
}
