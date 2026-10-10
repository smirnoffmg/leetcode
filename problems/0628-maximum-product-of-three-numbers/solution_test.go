package maximumproductofthreenumbers

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func TestMaximumProduct(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{name: "example 1", nums: []int{1, 2, 3}, want: 6},
		{name: "example 2", nums: []int{1, 2, 3, 4}, want: 24},
		{name: "example 3", nums: []int{-1, -2, -3}, want: -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maximumProduct(tt.nums)
			if got != tt.want {
				t.Errorf("maximumProduct(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestMaximumProductMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(628, 2026))

	for range 2000 {
		nums := make([]int, 3+rng.IntN(7))
		for i := range nums {
			nums[i] = rng.IntN(11) - 5
		}

		want := bruteForceProduct(nums)
		if got := maximumProduct(slices.Clone(nums)); got != want {
			t.Fatalf("maximumProduct(%v) = %v, want %v", nums, got, want)
		}
	}
}

func TestMaximumProductDoesNotMutateInput(t *testing.T) {
	nums := []int{3, -4, 1, -2, 5}
	before := slices.Clone(nums)

	maximumProduct(nums)

	if !slices.Equal(nums, before) {
		t.Fatalf("input mutated: got %v, want %v", nums, before)
	}
}

func BenchmarkMaximumProduct(b *testing.B) {
	rng := rand.New(rand.NewPCG(628, 1))

	nums := make([]int, 10_000)
	for i := range nums {
		nums[i] = rng.IntN(2001) - 1000
	}

	// копия входа оставлена ради сравнимости с прежней версией, которая сортировала nums на месте
	work := make([]int, len(nums))

	b.ResetTimer()

	for b.Loop() {
		copy(work, nums)
		maximumProduct(work)
	}
}

// bruteForceProduct — эталон от определения: перебор всех троек индексов.
func bruteForceProduct(nums []int) int {
	best := math.MinInt

	for i := range nums {
		for j := i + 1; j < len(nums); j++ {
			for k := j + 1; k < len(nums); k++ {
				best = max(best, nums[i]*nums[j]*nums[k])
			}
		}
	}

	return best
}
