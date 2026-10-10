package longestsubarraydivisiblebykwithatmostonenegationi

import (
	"math/rand/v2"
	"testing"
)

func TestLongestSubarray(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{name: "example 1", nums: []int{4, 1, 2}, k: 3, want: 3},
		{name: "example 2", nums: []int{5, 3, 4}, k: 7, want: 2},
		{name: "example 3", nums: []int{2, 2, 5}, k: 6, want: 2},

		{name: "k is 1", nums: []int{5, -3, 7}, k: 1, want: 3},
		{name: "single valid", nums: []int{4}, k: 4, want: 1},
		// -x ≡ 0 тогда и только тогда, когда x ≡ 0: инверсия одиночке не помогает
		{name: "single invalid", nums: []int{3}, k: 4, want: 0},
		{name: "negative values", nums: []int{-4, -1, -2}, k: 3, want: 3},
		// инвертировать надо последний элемент отрезка — ловит проверку до учёта элемента
		{name: "negate last", nums: []int{1, 1, 2}, k: 5, want: 3},
		// лучший отрезок начинается с нулевой позиции: ловит проверку first > 0 вместо >= 0
		{name: "from the start", nums: []int{2, 2, 5}, k: 6, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestSubarray(tt.nums, tt.k); got != tt.want {
				t.Errorf("longestSubarray(%v, %d) = %d, want %d", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

func TestLongestSubarrayMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 1))

	for range 200000 {
		n := 1 + rng.IntN(9)
		k := 1 + rng.IntN(8)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(21) - 10
		}

		if got, want := longestSubarray(nums, k), bruteForce(nums, k); got != want {
			t.Fatalf("longestSubarray(%v, %d) = %d, want %d", nums, k, got, want)
		}
	}
}

// bruteForce — эталон от определения: каждый отрезок, без инверсии и с инверсией каждого элемента.
func bruteForce(nums []int, k int) int {
	best := 0
	for l := range nums {
		for r := l; r < len(nums); r++ {
			sum := 0
			for i := l; i <= r; i++ {
				sum += nums[i]
			}

			ok := sum%k == 0
			for i := l; i <= r && !ok; i++ {
				ok = (sum-2*nums[i])%k == 0
			}

			if ok {
				best = max(best, r-l+1)
			}
		}
	}

	return best
}

func BenchmarkLongestSubarray(b *testing.B) {
	const n = 1000

	rng := rand.New(rand.NewPCG(3, 2))
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.IntN(200_001) - 100_000
	}

	b.ResetTimer()
	for b.Loop() {
		longestSubarray(nums, 99_991)
	}
}
