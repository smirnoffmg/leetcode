package maximumlengthsubstringwithtwooccurrences

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func TestMaximumLengthSubstring(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int
	}{
		{name: "example 1", s: "bcbbbcba", want: 4},
		{name: "example 2", s: "aaaa", want: 2},

		// минимальная длина по ограничениям
		{name: "two same", s: "aa", want: 2},
		{name: "two distinct", s: "ab", want: 2},

		// вся строка подходит целиком
		{name: "all distinct", s: "abcdef", want: 6},
		{name: "each exactly twice", s: "aabbcc", want: 6},
		{name: "whole string with tail", s: "aabbc", want: 5},

		// одна буква
		{name: "triple letter", s: "aaa", want: 2},
		{name: "long single letter", s: strings.Repeat("a", 100), want: 2},

		// третье вхождение обрезает окно
		{name: "two blocks of three", s: "aaabbb", want: 4},
		{name: "third occurrence in middle", s: "aabaab", want: 4},
		{name: "period three", s: "abcabcabc", want: 6},

		// ответ упирается в край строки
		{name: "answer at start", s: "bbaaa", want: 4},
		{name: "answer at end", s: "aaabb", want: 4},
		{name: "answer includes last char", s: "aaab", want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maximumLengthSubstring(tt.s)
			if got != tt.want {
				t.Errorf("maximumLengthSubstring(%v) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestMaximumLengthSubstringMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(3090, 2026))

	for range 2000 {
		n := 2 + rng.IntN(8)
		alphabet := 2 + rng.IntN(2)

		letters := make([]byte, n)
		for i := range letters {
			letters[i] = byte('a' + rng.IntN(alphabet))
		}

		s := string(letters)
		want := bruteForceMaxLen(s)

		if got := maximumLengthSubstring(s); got != want {
			t.Fatalf("maximumLengthSubstring(%q) = %v, want %v", s, got, want)
		}
	}
}

func BenchmarkMaximumLengthSubstring(b *testing.B) {
	rng := rand.New(rand.NewPCG(3090, 100))

	letters := make([]byte, 100)
	for i := range letters {
		letters[i] = byte('a' + rng.IntN(26))
	}

	s := string(letters)

	b.ResetTimer()

	for b.Loop() {
		maximumLengthSubstring(s)
	}
}

// bruteForceMaxLen — эталон от определения: перебор всех подстрок с проверкой,
// что каждая буква входит не больше двух раз.
func bruteForceMaxLen(s string) int {
	best := 0

	for i := range s {
		for j := i + 1; j <= len(s); j++ {
			if atMostTwoOfEach(s[i:j]) && j-i > best {
				best = j - i
			}
		}
	}

	return best
}

func atMostTwoOfEach(sub string) bool {
	var count [26]int

	for i := range len(sub) {
		count[sub[i]-'a']++
		if count[sub[i]-'a'] > 2 {
			return false
		}
	}

	return true
}
