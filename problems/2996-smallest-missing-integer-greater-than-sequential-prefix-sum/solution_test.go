package smallestmissingintegergreaterthansequentialprefixsum

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestMissingInteger(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{name: "example 1", nums: []int{1, 2, 3, 2, 5}, want: 6},
		{name: "example 2", nums: []int{3, 4, 5, 1, 12, 14, 13}, want: 15},

		// массив короче 50 элементов
		{name: "single one", nums: []int{1}, want: 2},
		{name: "single max value", nums: []int{50}, want: 51},
		{name: "two elements, sum occupied", nums: []int{2, 5}, want: 3},
		{name: "two elements, sum occupied by first", nums: []int{4, 9}, want: 5},

		// префикс — только nums[0], равные соседи последовательность не продолжают
		{name: "equal neighbours", nums: []int{3, 3}, want: 4},
		{name: "all equal", nums: []int{2, 2, 2}, want: 3},
		{name: "decreasing", nums: []int{5, 4, 3, 2, 1}, want: 6},
		{name: "step of two", nums: []int{1, 3, 5}, want: 2},

		// последовательный участок не в начале в счёт не идёт
		{name: "sequential run after first element", nums: []int{1, 3, 4, 5}, want: 2},
		{name: "longer run later", nums: []int{7, 1, 2, 3, 4, 5}, want: 8},

		// сумма префикса больше любого элемента: ответ равен сумме
		{name: "prefix sum exceeds values", nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, want: 55},
		{name: "whole array sequential", nums: sequential(1, 50), want: 1275},
		{name: "sequential from 30", nums: sequential(30, 21), want: 840},

		// ответ выходит за 50, хотя сумма префикса маленькая
		{name: "all values 1..50 present", nums: append([]int{50}, sequential(1, 49)...), want: 51},
		{name: "answer exactly fifty", nums: []int{49, 48}, want: 50},

		// цепочка занятых значений начиная с суммы префикса
		{name: "long occupied run above sum", nums: []int{3, 4, 5, 12, 13, 14, 15, 16, 1}, want: 17},
		{name: "sum itself missing", nums: []int{3, 4, 5, 13, 14, 15}, want: 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := missingInteger(slices.Clone(tt.nums))
			if got != tt.want {
				t.Errorf("missingInteger(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestMissingIntegerMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2996, 2024))

	for range 5000 {
		nums := randomNums(rng)

		want := bruteForceMissing(nums)
		if got := missingInteger(slices.Clone(nums)); got != want {
			t.Fatalf("missingInteger(%v) = %v, want %v", nums, got, want)
		}
	}
}

// Крайние размеры по условию: длина 50, значения 1..50.
func TestMissingIntegerMatchesBruteForceMaxSize(t *testing.T) {
	rng := rand.New(rand.NewPCG(2996, 50))

	for range 2000 {
		nums := make([]int, 50)
		for i := range nums {
			nums[i] = 1 + rng.IntN(50)
		}

		want := bruteForceMissing(nums)
		if got := missingInteger(slices.Clone(nums)); got != want {
			t.Fatalf("missingInteger(%v) = %v, want %v", nums, got, want)
		}
	}
}

func TestMissingIntegerDoesNotMutateInput(t *testing.T) {
	nums := []int{3, 4, 5, 1, 12, 14, 13}
	before := slices.Clone(nums)

	missingInteger(nums)

	if !slices.Equal(nums, before) {
		t.Fatalf("вход изменён: got %v, want %v", nums, before)
	}
}

func BenchmarkMissingIntegerSequential(b *testing.B) {
	nums := sequential(1, 50)

	for b.Loop() {
		missingInteger(nums)
	}
}

func BenchmarkMissingIntegerDenseValues(b *testing.B) {
	nums := append([]int{50}, sequential(1, 49)...)

	for b.Loop() {
		missingInteger(nums)
	}
}

func BenchmarkMissingIntegerRandom(b *testing.B) {
	rng := rand.New(rand.NewPCG(2996, 7))

	nums := make([]int, 50)
	for i := range nums {
		nums[i] = 1 + rng.IntN(50)
	}

	for b.Loop() {
		missingInteger(nums)
	}
}

func BenchmarkMissingIntegerSingle(b *testing.B) {
	nums := []int{50}

	for b.Loop() {
		missingInteger(nums)
	}
}

// sequential возвращает [start, start+1, ..., start+n-1].
func sequential(start, n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = start + i
	}

	return nums
}

// randomNums генерирует короткие входы из узкого алфавита: так последовательные
// префиксы и занятые значения около суммы встречаются достаточно часто.
func randomNums(rng *rand.Rand) []int {
	nums := make([]int, 1+rng.IntN(6))
	for i := range nums {
		nums[i] = 1 + rng.IntN(6)
	}

	return nums
}

// bruteForceMissing — эталон по определению задачи: длина префикса ищется прямой
// проверкой условия nums[j] == nums[j-1]+1, ответ — линейным перебором с поиском
// вхождения по всему массиву.
func bruteForceMissing(nums []int) int {
	sum := nums[0]

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1]+1 {
			break
		}

		sum += nums[i]
	}

	for x := sum; ; x++ {
		if !slices.Contains(nums, x) {
			return x
		}
	}
}
