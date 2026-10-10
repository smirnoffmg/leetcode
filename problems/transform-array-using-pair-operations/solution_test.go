package transformarrayusingpairoperations

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestCanTransform(t *testing.T) {
	tests := []struct {
		name           string
		source, target []int
		want           bool
	}{
		{name: "example 1", source: []int{1, 2, 3}, target: []int{0, 2, 4}, want: true},
		{name: "example 2", source: []int{-5, -5}, target: []int{-15, 5}, want: true},
		{name: "example 3", source: []int{1, 2, 1}, target: []int{0, 2, 5}, want: false},

		{name: "already equal", source: []int{7, -3}, target: []int{7, -3}, want: true},
		{name: "swap", source: []int{1, 2}, target: []int{2, 1}, want: true},
		{name: "off by one", source: []int{0, 0}, target: []int{0, 1}, want: false},
		// крайние значения по модулю
		{name: "extremes", source: []int{1e9, 1e9, -1e9}, target: []int{-1e9, -1e9, -1e9}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canTransform(tt.source, tt.target); got != tt.want {
				t.Errorf("canTransform(%v, %v) = %v, want %v", tt.source, tt.target, got, tt.want)
			}
		})
	}
}

// Достаточность проверяется честно: к source применяются сами операции из разбора,
// и результат обязан совпасть с target.
func TestCanTransformConstruction(t *testing.T) {
	rng := rand.New(rand.NewPCG(2, 2))

	for range 10000 {
		n := 2 + rng.IntN(7)
		source, target := make([]int, n), make([]int, n)
		for i := range n {
			source[i] = rng.IntN(7) - 3
			target[i] = rng.IntN(7) - 3
		}

		if !canTransform(source, target) {
			continue
		}

		got := slices.Clone(source)
		last := n - 1
		for j := range last {
			apply(got, last, j, target[j])
		}

		if !slices.Equal(got, target) {
			t.Fatalf("source=%v target=%v: операции дали %v", source, target, got)
		}
	}
}

func apply(a []int, i, j, delta int) {
	a[i], a[j] = a[i]+a[j]-delta, delta
}

func BenchmarkCanTransform(b *testing.B) {
	const n = 100_000

	rng := rand.New(rand.NewPCG(2, 1))
	source, target := make([]int, n), make([]int, n)
	for i := range n {
		source[i] = rng.IntN(2_000_000_001) - 1_000_000_000
		target[i] = rng.IntN(2_000_000_001) - 1_000_000_000
	}

	b.ResetTimer()
	for b.Loop() {
		canTransform(source, target)
	}
}
