package lengthoflongestsubarraywithatmostkfrequency

import (
	"math/rand/v2"
	"testing"
)

func TestMaxSubarrayLength(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want int
	}{
		{name: "example 1", nums: []int{1, 2, 3, 1, 2, 3, 1, 2}, k: 2, want: 6},
		{name: "example 2", nums: []int{1, 2, 1, 2, 1, 2, 1, 2}, k: 1, want: 2},
		{name: "example 3", nums: []int{5, 5, 5, 5, 5, 5, 5}, k: 4, want: 4},

		{name: "single element", nums: []int{7}, k: 1, want: 1},
		{name: "k covers whole array", nums: []int{1, 1, 2, 2}, k: 2, want: 4},
		{name: "all distinct", nums: []int{1, 2, 3, 4}, k: 1, want: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxSubarrayLength(tt.nums, tt.k)
			if got != tt.want {
				t.Errorf("maxSubarrayLength(%v, %v) = %v, want %v", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}

func TestMaxSubarrayLengthMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2958, 2026))

	for range 2000 {
		n := 1 + rng.IntN(9)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = 1 + rng.IntN(3)
		}
		k := 1 + rng.IntN(n)

		got := maxSubarrayLength(nums, k)
		want := bruteForce(nums, k)
		if got != want {
			t.Fatalf("maxSubarrayLength(%v, %v) = %v, want %v", nums, k, got, want)
		}
	}
}

const benchN = 100_000

func BenchmarkMaxSubarrayLengthAllSame(b *testing.B) {
	// один элемент, k=1: окно сжимается на каждом шаге
	nums := make([]int, benchN)
	for i := range nums {
		nums[i] = 42
	}

	for b.Loop() {
		maxSubarrayLength(nums, 1)
	}
}

func BenchmarkMaxSubarrayLengthAllDistinct(b *testing.B) {
	// все элементы различны: окно не сжимается, счётчик растёт до n ключей
	nums := make([]int, benchN)
	for i := range nums {
		nums[i] = i
	}

	for b.Loop() {
		maxSubarrayLength(nums, 1)
	}
}

func BenchmarkMaxSubarrayLengthFewValues(b *testing.B) {
	// три значения по кругу: типичный режим — окно и сжимается, и растёт
	nums := make([]int, benchN)
	for i := range nums {
		nums[i] = i % 3
	}

	for b.Loop() {
		maxSubarrayLength(nums, 100)
	}
}

// bruteForce перебирает все подотрезки и проверяет каждый по определению
// «хорошего» массива из условия.
func bruteForce(nums []int, k int) int {
	best := 0
	for i := range nums {
		for j := i; j < len(nums); j++ {
			freq := make(map[int]int)
			good := true
			for _, x := range nums[i : j+1] {
				freq[x]++
				if freq[x] > k {
					good = false
					break
				}
			}
			if good && j-i+1 > best {
				best = j - i + 1
			}
		}
	}
	return best
}
