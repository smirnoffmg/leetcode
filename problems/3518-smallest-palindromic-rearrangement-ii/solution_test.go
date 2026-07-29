package smallestpalindromicrearrangementii

import (
	"testing"
)

func TestSmallestPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		k    int
		want string
	}{
		{name: "example 1", s: "abba", k: 2, want: "baab"},
		{name: "example 2", s: "aa", k: 2, want: ""},
		{name: "example 3", s: "bacab", k: 1, want: "abcba"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smallestPalindrome(tt.s, tt.k)
			if got != tt.want {
				t.Errorf("smallestPalindrome(%v, %v) = %v, want %v", tt.s, tt.k, got, tt.want)
			}
		})
	}
}
