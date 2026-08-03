package stonegameiii

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestStoneGameIII(t *testing.T) {
	tests := []struct {
		name       string
		stoneValue []int
		want       string
	}{
		{name: "example 1", stoneValue: []int{1, 2, 3, 7}, want: "Bob"},
		{name: "example 2", stoneValue: []int{1, 2, 3, -9}, want: "Alice"},
		{name: "example 3", stoneValue: []int{1, 2, 3, 6}, want: "Tie"},
		{name: "single positive", stoneValue: []int{5}, want: "Alice"},
		{name: "single zero", stoneValue: []int{0}, want: "Tie"},
		{name: "single negative", stoneValue: []int{-5}, want: "Bob"},
		{name: "two stones", stoneValue: []int{1, 1}, want: "Alice"},
		{name: "three stones taken at once", stoneValue: []int{1, 1, 1}, want: "Alice"},
		// Алиса забирает три камня, Бобу остаётся один
		{name: "four ones", stoneValue: []int{1, 1, 1, 1}, want: "Alice"},
		// Алиса берёт −1 и −2, Боб вынужден взять −3: −3 против −3
		{name: "all negative", stoneValue: []int{-1, -2, -3}, want: "Tie"},
		{name: "all zeros", stoneValue: []int{0, 0, 0, 0, 0}, want: "Tie"},
		// у ряда единиц значение позиции периодично с периодом 6 и равно нулю
		// только при длине, кратной 6
		{name: "six ones", stoneValue: []int{1, 1, 1, 1, 1, 1}, want: "Tie"},
		{name: "seven ones", stoneValue: []int{1, 1, 1, 1, 1, 1, 1}, want: "Alice"},
		// какой бы префикс Алиса ни взяла, доступ к 7 достаётся Бобу
		{name: "mixed signs", stoneValue: []int{-1, -2, -3, 7}, want: "Bob"},
		// взять больше не значит лучше: 7+1+1 оставляет Бобу 9 и −2 по разнице,
		// лучший ход — взять 7+1 (−1), но и он проигрывает
		{name: "greedy trap", stoneValue: []int{7, 1, 1, 9, 1, 1}, want: "Bob"},
		{name: "boundary values", stoneValue: []int{1000, -1000, 1000, -1000}, want: "Alice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stoneGameIII(tt.stoneValue)
			if got != tt.want {
				t.Errorf("stoneGameIII(%v) = %v, want %v", tt.stoneValue, got, tt.want)
			}
		})
	}
}

func TestStoneGameIIIMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(1406, 2024))

	for range 300 {
		stoneValue := make([]int, 1+rng.IntN(12))
		for j := range stoneValue {
			stoneValue[j] = rng.IntN(21) - 10
		}

		want := bruteForceWinner(stoneValue)
		if got := stoneGameIII(stoneValue); got != want {
			t.Fatalf("stoneGameIII(%v) = %v, want %v", stoneValue, got, want)
		}
	}
}

func TestStoneGameIIIDoesNotMutateInput(t *testing.T) {
	stoneValue := []int{1, 2, 3, 7}
	before := make([]int, len(stoneValue))
	copy(before, stoneValue)

	stoneGameIII(stoneValue)

	for i := range stoneValue {
		if stoneValue[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", stoneValue, before)
		}
	}
}

func TestStoneGameIIIMaxLength(t *testing.T) {
	stoneValue := make([]int, 50_000)
	for i := range stoneValue {
		stoneValue[i] = 1
	}

	// 50000 не кратно 6, поэтому ничьей у ряда единиц не будет
	if got := stoneGameIII(stoneValue); got != "Alice" {
		t.Errorf("stoneGameIII(50000 ones) = %v, want Alice", got)
	}
}

// bruteForceWinner разыгрывает все варианты без мемоизации — эталон для сверки.
func bruteForceWinner(stoneValue []int) string {
	alice, bob := play(stoneValue, 0)

	switch {
	case alice > bob:
		return "Alice"
	case alice < bob:
		return "Bob"
	default:
		return "Tie"
	}
}

// play возвращает очки ходящего и его соперника на суффиксе stoneValue[i:].
func play(stoneValue []int, i int) (mover, opponent int) {
	if i >= len(stoneValue) {
		return 0, 0
	}

	taken, bestDiff := 0, math.MinInt

	for k := 0; k < 3 && i+k < len(stoneValue); k++ {
		taken += stoneValue[i+k]

		// после хода очередь переходит к сопернику: он становится ходящим на остатке
		nextMover, nextOpponent := play(stoneValue, i+k+1)
		score := taken + nextOpponent

		if diff := score - nextMover; diff > bestDiff {
			bestDiff = diff
			mover, opponent = score, nextMover
		}
	}

	return mover, opponent
}
