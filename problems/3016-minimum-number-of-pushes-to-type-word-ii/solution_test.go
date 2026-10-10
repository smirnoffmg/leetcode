package minimumnumberofpushestotypewordii

import (
	"math/rand/v2"
	"testing"
)

func TestMinimumPushes(t *testing.T) {
	tests := []struct {
		name string
		word string
		want int
	}{
		{name: "example 1", word: "abcde", want: 5},
		{name: "example 2", word: "xyzxyzxyzxyz", want: 12},
		{name: "example 3", word: "aabbccddeeffgghhiiiiii", want: 24},
		{name: "full alphabet", word: "abcdefghijklmnopqrstuvwxyz", want: 56},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := minimumPushes(tt.word)
			if got != tt.want {
				t.Errorf("minimumPushes(%v) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

func BenchmarkMinimumPushes(b *testing.B) {
	rng := rand.New(rand.NewPCG(3016, 1))

	letters := make([]byte, 100_000)
	for i := range letters {
		letters[i] = byte('a' + rng.IntN(26))
	}

	word := string(letters)

	b.ResetTimer()

	for b.Loop() {
		minimumPushes(word)
	}
}
