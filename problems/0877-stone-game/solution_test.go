package stonegame

import (
	"math/rand/v2"
	"testing"
)

func TestStoneGame(t *testing.T) {
	tests := []struct {
		name  string
		piles []int
		want  bool
	}{
		{name: "example 1", piles: []int{5, 3, 4, 5}, want: true},
		{name: "example 2", piles: []int{3, 7, 2, 3}, want: true},
		{name: "minimum length", piles: []int{1, 2}, want: true},
		{name: "minimum length reversed", piles: []int{2, 1}, want: true},
		{name: "big pile at the edge", piles: []int{500, 1, 1, 1}, want: true},
		{name: "big pile inside", piles: []int{1, 500, 1, 3}, want: true},
		// жадный первый ход (взять 2) проигрывает 3:6; оптимальный — взять 1
		// и забрать все чётные индексы, 6:3
		{name: "greedy trap", piles: []int{1, 1, 5, 2}, want: true},
		{name: "descending", piles: []int{9, 7, 5, 3, 1, 2}, want: true},
		{name: "ascending", piles: []int{1, 2, 3, 4, 5, 6}, want: true},
		{name: "alternating", piles: []int{1, 100, 1, 100, 1, 100}, want: true},
		{name: "max values", piles: []int{500, 500, 500, 499}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stoneGame(tt.piles)
			if got != tt.want {
				t.Errorf("stoneGame(%v) = %v, want %v", tt.piles, got, tt.want)
			}
		})
	}
}

// Ограничения задачи (чётная длина, нечётная сумма) гарантируют победу Алисы,
// поэтому на допустимых входах ответ всегда true и ничего не различает.
// Эти входы выходят за ограничения и проверяют сам расчёт функции ценности.
func TestStoneGameValueFunction(t *testing.T) {
	tests := []struct {
		name  string
		piles []int
		want  bool
	}{
		{name: "single pile", piles: []int{7}, want: true},
		{name: "tie on two equal piles", piles: []int{1, 1}, want: false},
		{name: "tie on four equal piles", piles: []int{3, 3, 3, 3}, want: false},
		{name: "tie by parity groups", piles: []int{2, 4, 4, 2}, want: false},
		// нечётная длина снимает чётностный аргумент: ходящий первым проигрывает
		{name: "odd length loss", piles: []int{2, 4, 55, 6, 8}, want: false},
		{name: "odd length three piles", piles: []int{1, 5, 2}, want: false},
		{name: "odd length big middle loss", piles: []int{1, 5, 233, 7, 1}, want: false},
		{name: "odd length win", piles: []int{9, 7, 5, 3, 1}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stoneGame(tt.piles)
			if got != tt.want {
				t.Errorf("stoneGame(%v) = %v, want %v", tt.piles, got, tt.want)
			}
		})
	}
}

func TestStoneGameMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(877, 2024))

	for range 300 {
		piles := make([]int, 1+rng.IntN(9))
		for j := range piles {
			piles[j] = 1 + rng.IntN(50)
		}

		want := bruteForceWins(piles)
		if got := stoneGame(piles); got != want {
			t.Fatalf("stoneGame(%v) = %v, want %v", piles, got, want)
		}
	}
}

// Чётностный аргумент: при чётной длине Алиса забирает целиком либо все чётные
// индексы, либо все нечётные, а нечётная сумма исключает равенство групп.
func TestStoneGameAliceAlwaysWinsUnderConstraints(t *testing.T) {
	rng := rand.New(rand.NewPCG(877, 42))

	for range 500 {
		piles := randomValidPiles(rng, 250)

		var even, odd int

		for i, p := range piles {
			if i%2 == 0 {
				even += p
			} else {
				odd += p
			}
		}

		if even == odd {
			t.Fatalf("группы по чётности совпали при нечётной сумме: %v", piles)
		}

		if !stoneGame(piles) {
			t.Fatalf("stoneGame(%v) = false, want true", piles)
		}
	}
}

func TestStoneGameMaxLength(t *testing.T) {
	rng := rand.New(rand.NewPCG(877, 500))

	piles := randomValidPiles(rng, 250)
	if !stoneGame(piles) {
		t.Fatalf("stoneGame на входе длины %d = false, want true", len(piles))
	}
}

func TestStoneGameDoesNotMutateInput(t *testing.T) {
	piles := []int{5, 3, 4, 5}
	before := make([]int, len(piles))
	copy(before, piles)

	stoneGame(piles)

	for i := range piles {
		if piles[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", piles, before)
		}
	}
}

func BenchmarkStoneGame(b *testing.B) {
	rng := rand.New(rand.NewPCG(877, 1))
	piles := randomValidPiles(rng, 250)

	for b.Loop() {
		stoneGame(piles)
	}
}

// randomValidPiles генерирует вход по ограничениям задачи: чётная длина
// (до 2*maxPairs куч), значения 1..500 и нечётная сумма.
func randomValidPiles(rng *rand.Rand, maxPairs int) []int {
	piles := make([]int, 2*(1+rng.IntN(maxPairs)))

	sum := 0

	for i := range piles {
		piles[i] = 1 + rng.IntN(500)
		sum += piles[i]
	}

	if sum%2 == 0 {
		if piles[0] == 500 {
			piles[0]--
		} else {
			piles[0]++
		}
	}

	return piles
}

// bruteForceWins разыгрывает все варианты без мемоизации — эталон для сверки.
func bruteForceWins(piles []int) bool {
	mover, opponent := play(piles, 0, len(piles)-1)

	return mover > opponent
}

// play возвращает очки ходящего и его соперника на отрезке piles[i..j].
func play(piles []int, i, j int) (mover, opponent int) {
	if i == j {
		return piles[i], 0
	}

	leftMover, leftOpponent := play(piles, i+1, j)
	rightMover, rightOpponent := play(piles, i, j-1)

	takeLeft := piles[i] + leftOpponent
	takeRight := piles[j] + rightOpponent

	if takeLeft-leftMover >= takeRight-rightMover {
		return takeLeft, leftMover
	}

	return takeRight, rightMover
}
