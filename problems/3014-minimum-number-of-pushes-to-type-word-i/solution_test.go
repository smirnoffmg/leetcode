package minimumnumberofpushestotypewordi

import (
	"testing"
)

func TestMinimumPushes(t *testing.T) {
	tests := []struct {
		name string
		word string
		want int
	}{
		{name: "example 1", word: "abcde", want: 5},
		{name: "example 2", word: "xycdefghij", want: 12},
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
