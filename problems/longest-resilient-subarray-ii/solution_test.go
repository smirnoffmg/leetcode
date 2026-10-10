package longestresilientsubarrayii

import (
	"math/rand/v2"
	"testing"
)

func TestResilientSubarray(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{name: "example 1", nums: []int{2, 4, 6, 3}, k: 2, want: 3},
		{name: "example 2", nums: []int{1, 4, 7, 1}, k: 3, want: 4},
		{name: "example 3", nums: []int{5, 5, 4, 8}, k: 4, want: 2},

		{name: "single", nums: []int{7}, k: 5, want: 1},
		{name: "k is one", nums: []int{3, 1, 4, 1, 5}, k: 1, want: 5},
		// остаток 2 при k=4: шаг 2, из серии длины 4 годится только длина 3
		{name: "step by gcd", nums: []int{2, 6, 10, 14}, k: 4, want: 3},
		{name: "no pair fits", nums: []int{1, 2, 3}, k: 5, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resilientSubarray(tt.nums, tt.k); got != tt.want {
				t.Errorf("resilientSubarray(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

// Эталон — прямо по определению: каждый подмассив, каждое удаление.
func bruteResilient(nums []int, k int) int {
	best := 0
	for i := range nums {
		for j := i; j < len(nums); j++ {
			sum := 0
			for p := i; p <= j; p++ {
				sum += nums[p]
			}

			ok := true
			for p := i; p <= j && ok; p++ {
				ok = (sum-nums[p])%k == 0
			}
			if ok {
				best = max(best, j-i+1)
			}
		}
	}

	return best
}

func TestResilientSubarrayRandom(t *testing.T) {
	rng := rand.New(rand.NewPCG(4, 4))

	for range 20000 {
		n := 1 + rng.IntN(9)
		k := 1 + rng.IntN(6)
		// два остатка на вход: серии с одним остатком длинные, а остаток
		// произвольный, в том числе с gcd(r, k) > 1
		base := [2]int{1 + rng.IntN(2*k), 1 + rng.IntN(2*k)}
		nums := make([]int, n)
		for i := range nums {
			nums[i] = base[rng.IntN(2)] + rng.IntN(3)*k
		}

		if got, want := resilientSubarray(nums, k), bruteResilient(nums, k); got != want {
			t.Fatalf("nums=%v k=%d: got %d, want %d", nums, k, got, want)
		}
	}
}

func BenchmarkResilientSubarray(b *testing.B) {
	const n = 100_000

	rng := rand.New(rand.NewPCG(4, 1))
	nums := make([]int, n)
	for i := range nums {
		nums[i] = 1 + rng.IntN(100_000)
	}

	b.ResetTimer()
	for b.Loop() {
		resilientSubarray(nums, 12)
	}
}
