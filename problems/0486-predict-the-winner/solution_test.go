package predictthewinner

import (
	"math/rand/v2"
	"testing"
)

func TestPredictTheWinner(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{name: "example 1", nums: []int{1, 5, 2}, want: false},
		{name: "example 2", nums: []int{1, 5, 233, 7}, want: true},
		{name: "single element", nums: []int{5}, want: true},
		{name: "single zero", nums: []int{0}, want: true},
		{name: "two equal", nums: []int{2, 2}, want: true},
		{name: "two different", nums: []int{1, 2}, want: true},
		{name: "all equal even count", nums: []int{3, 3, 3, 3}, want: true},
		{name: "all equal odd count", nums: []int{1, 1, 1}, want: true},
		// жадный выбор большего конца даёт здесь 65:10 в пользу первого,
		// но при оптимальной игре обоих он проигрывает
		{name: "greedy trap", nums: []int{2, 4, 55, 6, 8}, want: false},
		{name: "descending", nums: []int{9, 7, 5, 3, 1}, want: true},
		{name: "zeros around one value", nums: []int{0, 0, 7, 0}, want: true},
		{name: "max values", nums: []int{10000000, 0, 10000000}, want: true},
		{name: "alternating", nums: []int{1, 100, 1, 100, 1, 100}, want: true},
		{
			name: "max length",
			nums: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := predictTheWinner(tt.nums)
			if got != tt.want {
				t.Errorf("predictTheWinner(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

func TestPredictTheWinnerMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(486, 2024))

	for range 300 {
		nums := make([]int, 1+rng.IntN(9))
		for j := range nums {
			nums[j] = rng.IntN(50)
		}

		want := bruteForceWins(nums)
		if got := predictTheWinner(nums); got != want {
			t.Fatalf("predictTheWinner(%v) = %v, want %v", nums, got, want)
		}
	}
}

func TestPredictTheWinnerDoesNotMutateInput(t *testing.T) {
	nums := []int{1, 5, 233, 7}
	before := make([]int, len(nums))
	copy(before, nums)

	predictTheWinner(nums)

	for i := range nums {
		if nums[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", nums, before)
		}
	}
}

// bruteForceWins разыгрывает все варианты без мемоизации — эталон для сверки.
func bruteForceWins(nums []int) bool {
	first, second := play(nums, 0, len(nums)-1)

	return first >= second
}

// play возвращает очки ходящего и его соперника на отрезке nums[i..j].
func play(nums []int, i, j int) (mover, opponent int) {
	if i == j {
		return nums[i], 0
	}

	leftMover, leftOpponent := play(nums, i+1, j)
	rightMover, rightOpponent := play(nums, i, j-1)

	takeLeft := nums[i] + leftOpponent
	takeRight := nums[j] + rightOpponent

	if takeLeft-leftMover >= takeRight-rightMover {
		return takeLeft, leftMover
	}

	return takeRight, rightMover
}
