package maximumproductoftwodigits

import (
	"math/rand/v2"
	"testing"
)

func TestMaxProduct(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "example 1", n: 31, want: 3},
		{name: "example 2", n: 22, want: 4},
		{name: "example 3", n: 124, want: 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxProduct(tt.n)
			if got != tt.want {
				t.Errorf("maxProduct(%v) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestMaxProductMatchesBruteForce(t *testing.T) {
	for n := 10; n <= 100_000; n++ {
		if got, want := maxProduct(n), bruteForceProduct(n); got != want {
			t.Fatalf("maxProduct(%d) = %d, want %d", n, got, want)
		}
	}

	rng := rand.New(rand.NewPCG(3536, 2026))
	for range 10_000 {
		n := 10 + rng.IntN(1_000_000_000-9)
		if got, want := maxProduct(n), bruteForceProduct(n); got != want {
			t.Fatalf("maxProduct(%d) = %d, want %d", n, got, want)
		}
	}
}

func BenchmarkMaxProduct(b *testing.B) {
	for b.Loop() {
		maxProduct(1_000_000_000)
	}
}

// bruteForceProduct — эталон от определения: перебор всех пар позиций цифр.
func bruteForceProduct(n int) int {
	var digits []int
	for ; n > 0; n /= 10 {
		digits = append(digits, n%10)
	}

	best := 0
	for i := range digits {
		for j := i + 1; j < len(digits); j++ {
			best = max(best, digits[i]*digits[j])
		}
	}

	return best
}
