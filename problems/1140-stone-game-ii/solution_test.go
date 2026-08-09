package stonegameii

import (
	"math/rand/v2"
	"testing"
)

func TestStoneGameII(t *testing.T) {
	tests := []struct {
		name  string
		piles []int
		want  int
	}{
		{name: "example 1", piles: []int{2, 7, 9, 4, 4}, want: 10},
		{name: "example 2", piles: []int{1, 2, 3, 4, 5, 100}, want: 104},
		{name: "single pile", piles: []int{5}, want: 5},
		// при M = 1 доступно X = 2, то есть весь ряд
		{name: "two piles taken at once", piles: []int{1, 2}, want: 3},
		{name: "two piles reversed", piles: []int{100, 1}, want: 101},
		// взять две кучи (2) лучше, чем одну: иначе Боб забирает оставшиеся две
		{name: "three ones", piles: []int{1, 1, 1}, want: 2},
		{name: "four ones", piles: []int{1, 1, 1, 1}, want: 2},
		// взять меньше выгоднее: после X = 1 у Боба обе альтернативы отдают 100 Алисе
		{name: "small move keeps the big pile", piles: []int{1, 1, 1, 100}, want: 101},
		// на куче короче Боб дотягивается до 100 при любом первом ходе Алисы
		{name: "big pile out of reach", piles: []int{1, 1, 100}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stoneGameII(tt.piles)
			if got != tt.want {
				t.Errorf("stoneGameII(%v) = %v, want %v", tt.piles, got, tt.want)
			}
		})
	}
}

func TestStoneGameIIMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(1140, 2026))

	for range 300 {
		piles := make([]int, 1+rng.IntN(9))
		for j := range piles {
			piles[j] = 1 + rng.IntN(9)
		}

		want := bruteForceAlice(piles)
		if got := stoneGameII(piles); got != want {
			t.Fatalf("stoneGameII(%v) = %v, want %v", piles, got, want)
		}
	}
}

// таблица переживает вызовы, поэтому записи прошлого расчёта не должны читаться
func TestStoneGameIIReusesTableAcrossCalls(t *testing.T) {
	big := make([]int, 100)
	for i := range big {
		big[i] = 10_000
	}

	first := stoneGameII(big)

	if got := stoneGameII([]int{2, 7, 9, 4, 4}); got != 10 {
		t.Errorf("после большого входа stoneGameII([2 7 9 4 4]) = %v, want 10", got)
	}

	if got := stoneGameII(big); got != first {
		t.Errorf("повторный вызов на том же входе = %v, want %v", got, first)
	}
}

func TestStoneGameIIDoesNotMutateInput(t *testing.T) {
	piles := []int{2, 7, 9, 4, 4}
	before := make([]int, len(piles))
	copy(before, piles)

	stoneGameII(piles)

	for i := range piles {
		if piles[i] != before[i] {
			t.Fatalf("input mutated: got %v, want %v", piles, before)
		}
	}
}

func TestStoneGameIIMaxLength(t *testing.T) {
	piles := make([]int, 100)
	for i := range piles {
		piles[i] = 10_000
	}

	// ряд одинаковых куч: ответ — доля Алисы от 10^6, она заведомо не пуста
	// и не равна всему ряду
	got := stoneGameII(piles)
	if got <= 0 || got >= 1_000_000 {
		t.Errorf("stoneGameII(100 piles of 10000) = %v, want strictly between 0 and 1000000", got)
	}
}

func BenchmarkStoneGameII(b *testing.B) {
	piles := make([]int, 100)
	for i := range piles {
		piles[i] = (i*37)%10_000 + 1
	}

	b.ResetTimer()

	for b.Loop() {
		stoneGameII(piles)
	}
}

// bruteForceAlice разыгрывает все партии без мемоизации — эталон для сверки.
func bruteForceAlice(piles []int) int {
	alice, _ := play(piles, 0, 1)
	return alice
}

// play возвращает очки ходящего и его соперника на суффиксе piles[i:] при данном m.
// Каждый игрок максимизирует собственную добычу — так задача и сформулирована.
func play(piles []int, i, m int) (mover, opponent int) {
	taken, best := 0, -1

	for x := 1; x <= 2*m && i+x <= len(piles); x++ {
		taken += piles[i+x-1]

		// после хода очередь переходит к сопернику: он становится ходящим на остатке
		nextMover, nextOpponent := play(piles, i+x, max(m, x))

		if score := taken + nextOpponent; score > best {
			best = score
			mover, opponent = score, nextMover
		}
	}

	return mover, opponent
}
