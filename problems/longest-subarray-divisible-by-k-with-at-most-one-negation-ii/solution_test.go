package longestsubarraydivisiblebykwithatmostonenegationii

import (
	"math/rand/v2"
	"strconv"
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

// На больших входах эталоном служит квадратичное решение версии I — оно само сверено с перебором.
func TestLongestSubarrayMatchesQuadratic(t *testing.T) {
	rng := rand.New(rand.NewPCG(4, 1))

	for range 300 {
		n := 1 + rng.IntN(1000)
		k := 1 + rng.IntN(3000)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(200_001) - 100_000
		}

		if got, want := longestSubarray(nums, k), quadratic(nums, k); got != want {
			t.Fatalf("n=%d k=%d: longestSubarray = %d, want %d", n, k, got, want)
		}
	}
}

func quadratic(nums []int, k int) int {
	mod := func(x int) int { return (x%k + k) % k }
	stamp := make([]int, k)
	best := 0

	for l := range nums {
		sum := 0
		for r := l; r < len(nums); r++ {
			sum += nums[r]
			stamp[mod(2*nums[r])] = l + 1
			if s := mod(sum); s == 0 || stamp[s] == l+1 {
				best = max(best, r-l+1)
			}
		}
	}

	return best
}

func BenchmarkLongestSubarray(b *testing.B) {
	const n = 100_000

	rng := rand.New(rand.NewPCG(4, 2))
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.IntN(200_001) - 100_000
	}

	// нечётное k: 2x пробегает все остатки, и пар (u, d) набирается под k²
	for _, k := range []int{2999, 3000} {
		b.Run("k="+strconv.Itoa(k), func(b *testing.B) {
			for b.Loop() {
				longestSubarray(nums, k)
			}
		})
	}
}
