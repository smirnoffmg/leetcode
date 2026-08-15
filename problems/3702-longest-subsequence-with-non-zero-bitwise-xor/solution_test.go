package longestsubsequencewithnonzerobitwisexor

import (
	"math/bits"
	"math/rand/v2"
	"testing"
)

func TestLongestSubsequence(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{name: "example 1", nums: []int{1, 2, 3}, want: 2},
		{name: "example 2", nums: []int{2, 3, 4}, want: 3},

		// один элемент
		{name: "single zero", nums: []int{0}, want: 0},
		{name: "single non-zero", nums: []int{1}, want: 1},
		{name: "single max value", nums: []int{1_000_000_000}, want: 1},

		// все нули — ненулевой XOR не собрать ни из чего
		{name: "two zeros", nums: []int{0, 0}, want: 0},
		{name: "many zeros", nums: []int{0, 0, 0, 0, 0}, want: 0},

		// XOR всего массива нулевой, но есть ненулевые элементы
		{name: "pair cancels", nums: []int{1, 1}, want: 1},
		{name: "zeros around pair", nums: []int{0, 1, 1, 0}, want: 3},
		{name: "triple cancels", nums: []int{1, 2, 3, 0}, want: 3},
		{name: "four equal", nums: []int{1, 1, 1, 1}, want: 3},
		{name: "pair with zero tail", nums: []int{2, 2, 0}, want: 2},

		// XOR всего массива ненулевой — берём всё
		{name: "all equal odd count", nums: []int{7, 7, 7}, want: 3},
		{name: "zeros do not spoil total", nums: []int{0, 0, 5}, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestSubsequence(tt.nums)
			if got != tt.want {
				t.Errorf("longestSubsequence(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestLongestSubsequenceMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(3702, 2026))

	for range 5000 {
		n := 1 + rng.IntN(9)

		// узкий алфавит {0, 1, 2}: часто даёт нулевые XOR и массивы из одних нулей
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(3)
		}

		want := bruteForceLongest(nums)

		if got := longestSubsequence(nums); got != want {
			t.Fatalf("longestSubsequence(%v) = %v, want %v", nums, got, want)
		}
	}
}

func BenchmarkLongestSubsequence(b *testing.B) {
	const n = 100_000

	rng := rand.New(rand.NewPCG(3702, 1))

	random := make([]int, n)
	for i := range random {
		random[i] = rng.IntN(1_000_000_001)
	}

	// XOR всего массива нулевой: пары равных элементов вперемешку с нулями
	xorZero := make([]int, n)
	for i := 0; i < n; i += 2 {
		v := rng.IntN(1_000_000_001)
		xorZero[i], xorZero[i+1] = v, v
	}

	allZeros := make([]int, n)

	for _, bm := range []struct {
		name string
		nums []int
	}{
		{name: "random", nums: random},
		{name: "xor zero", nums: xorZero},
		{name: "all zeros", nums: allZeros},
	} {
		b.Run(bm.name, func(b *testing.B) {
			for b.Loop() {
				longestSubsequence(bm.nums)
			}
		})
	}
}

// bruteForceLongest — эталон от определения задачи: перебор всех непустых
// подмножеств, максимум длины среди тех, чей XOR ненулевой.
func bruteForceLongest(nums []int) int {
	best := 0

	for mask := 1; mask < 1<<len(nums); mask++ {
		xor := 0
		for i, v := range nums {
			if mask&(1<<i) != 0 {
				xor ^= v
			}
		}

		if xor != 0 {
			best = max(best, bits.OnesCount(uint(mask)))
		}
	}

	return best
}
