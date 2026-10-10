package maximumproductpairwithtargetsum

import (
	"math/rand/v2"
	"testing"
)

// Ответ неоднозначен («any of them»), поэтому тест сверяет не индексы, а их свойства:
// пара допустима и её произведение равно эталонному максимуму.
func check(t *testing.T, nums []int, target int, got []int) {
	t.Helper()

	want, ok := bruteMax(nums, target)
	if len(got) != 2 {
		t.Fatalf("nums=%v target=%d: got %v, want two indices", nums, target, got)
	}
	if !ok {
		if got[0] != -1 || got[1] != -1 {
			t.Fatalf("nums=%v target=%d: got %v, want [-1 -1]", nums, target, got)
		}
		return
	}

	i, j := got[0], got[1]
	if i < 0 || j < 0 || i >= len(nums) || j >= len(nums) || i == j ||
		nums[i]+nums[j] != target || nums[i] <= nums[j] {
		t.Fatalf("nums=%v target=%d: got %v, pair is not valid", nums, target, got)
	}
	if p := nums[i] * nums[j]; p != want {
		t.Fatalf("nums=%v target=%d: got %v with product %d, want %d", nums, target, got, p, want)
	}
}

func bruteMax(nums []int, target int) (int, bool) {
	best, ok := 0, false
	for i := range nums {
		for j := range nums {
			if i == j || nums[i]+nums[j] != target || nums[i] <= nums[j] {
				continue
			}
			if p := nums[i] * nums[j]; !ok || p > best {
				best, ok = p, true
			}
		}
	}
	return best, ok
}

func TestMaxProductPair(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
	}{
		{name: "example 1", nums: []int{1, 2, 3, 4}, target: 5},
		{name: "example 2", nums: []int{-3, -1, 4, 2}, target: 1},
		{name: "example 3", nums: []int{3, 3, 5}, target: 6},

		{name: "both negative", nums: []int{-100, -1, -99, -2}, target: -101},
		{name: "zero product", nums: []int{0, 7}, target: 7},
		{name: "all equal", nums: []int{4, 4, 4}, target: 8},
		{name: "extremes", nums: []int{100, -100}, target: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check(t, tt.nums, tt.target, maxProductPair(tt.nums, tt.target))
		})
	}

	// Ответы примеров из условия: там пара с максимумом единственна.
	for _, ex := range []struct {
		nums   []int
		target int
		want   [2]int
	}{
		{[]int{1, 2, 3, 4}, 5, [2]int{2, 1}},
		{[]int{-3, -1, 4, 2}, 1, [2]int{3, 1}},
		{[]int{3, 3, 5}, 6, [2]int{-1, -1}},
	} {
		if got := maxProductPair(ex.nums, ex.target); got[0] != ex.want[0] || got[1] != ex.want[1] {
			t.Errorf("maxProductPair(%v, %d) = %v, want %v", ex.nums, ex.target, got, ex.want)
		}
	}
}

func TestMaxProductPairRandom(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 3))

	for range 20000 {
		n := 2 + rng.IntN(8)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(9) - 4
		}
		target := rng.IntN(17) - 8

		check(t, nums, target, maxProductPair(nums, target))
	}
}

func BenchmarkMaxProductPair(b *testing.B) {
	const n = 100

	rng := rand.New(rand.NewPCG(3, 1))
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.IntN(201) - 100
	}

	b.ResetTimer()
	for b.Loop() {
		maxProductPair(nums, 0)
	}
}
