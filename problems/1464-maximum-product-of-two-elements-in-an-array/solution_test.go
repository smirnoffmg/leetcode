package maximumproductoftwoelementsinanarray

import (
	"math/rand/v2"
	"testing"
)

func TestMaxProduct(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{name: "example 1", nums: []int{3, 4, 5, 2}, want: 12},
		{name: "example 2", nums: []int{1, 5, 4, 5}, want: 16},
		{name: "example 3", nums: []int{3, 7}, want: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxProduct(tt.nums)
			if got != tt.want {
				t.Errorf("maxProduct(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func BenchmarkMaxProduct(b *testing.B) {
	rng := rand.New(rand.NewPCG(1464, 1))

	nums := make([]int, 500)
	for i := range nums {
		nums[i] = 1 + rng.IntN(1000)
	}

	b.ResetTimer()

	for b.Loop() {
		maxProduct(nums)
	}
}
