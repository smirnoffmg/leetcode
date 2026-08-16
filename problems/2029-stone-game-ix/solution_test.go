package stonegameix

import (
	"math/rand/v2"
	"testing"
)

func TestStoneGameIX(t *testing.T) {
	tests := []struct {
		name   string
		stones []int
		want   bool
	}{
		{name: "example 1", stones: []int{2, 1}, want: true},
		{name: "example 2", stones: []int{2}, want: false},
		{name: "example 3", stones: []int{5, 1, 2, 4, 3}, want: false},

		// вырожденные размеры
		{name: "single zero stone", stones: []int{3}, want: false},
		{name: "single one", stones: []int{1}, want: false},
		{name: "pair of ones", stones: []int{1, 1}, want: false},

		// без нулей: оба остатка есть — Алиса начинает с меньшинства
		{name: "both residues, ones majority", stones: []int{1, 2, 1}, want: true},
		{name: "one residue missing", stones: []int{1, 1, 1, 1}, want: false},

		// чётные нули взаимно уничтожаются
		{name: "even zeros, both residues", stones: []int{3, 3, 1, 2}, want: true},
		{name: "even zeros, residue missing", stones: []int{3, 3, 2, 2}, want: false},

		// нечётный нуль: нужен перекос больше двух
		{name: "odd zero, balanced", stones: []int{3, 1, 2}, want: false},
		{name: "odd zero, skew of two", stones: []int{3, 1, 1}, want: false},
		{name: "odd zero, skew of three", stones: []int{3, 1, 1, 1}, want: true},
		{name: "odd zero, skew of four", stones: []int{3, 1, 1, 1, 1}, want: true},
		{name: "all zeros", stones: []int{3, 6, 9}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stoneGameIX(tt.stones)
			if got != tt.want {
				t.Errorf("stoneGameIX(%v) = %v, want %v", tt.stones, got, tt.want)
			}
		})
	}
}

func TestStoneGameIXMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2029, 2026))

	for range 2000 {
		n := 1 + rng.IntN(9)
		stones := make([]int, n)
		for i := range stones {
			stones[i] = 1 + rng.IntN(9)
		}

		want := bruteForceAliceWins(stones)

		if got := stoneGameIX(stones); got != want {
			t.Fatalf("stoneGameIX(%v) = %v, want %v", stones, got, want)
		}
	}
}

func BenchmarkStoneGameIX(b *testing.B) {
	const n = 100_000

	rng := rand.New(rand.NewPCG(2029, 1))
	stones := make([]int, n)
	for i := range stones {
		stones[i] = 1 + rng.IntN(10_000)
	}

	b.ResetTimer()

	for b.Loop() {
		stoneGameIX(stones)
	}
}

// bruteForceAliceWins — минимакс от определения игры, независимый от идеи решения:
// состояние — счётчики остатков, сумма по модулю 3 и чей ход.
func bruteForceAliceWins(stones []int) bool {
	var cnt [3]int
	for _, v := range stones {
		cnt[v%3]++
	}

	type state struct {
		c0, c1, c2, sum, turn int
	}
	memo := map[state]bool{}

	var aliceWins func(c [3]int, sum, turn int) bool
	aliceWins = func(c [3]int, sum, turn int) bool {
		if c[0]+c[1]+c[2] == 0 {
			return false
		}

		key := state{c[0], c[1], c[2], sum, turn}
		if v, ok := memo[key]; ok {
			return v
		}

		aliceMoves := turn == 0
		result := !aliceMoves

		for r := range 3 {
			if c[r] == 0 {
				continue
			}

			var outcome bool
			if (sum+r)%3 == 0 {
				outcome = !aliceMoves
			} else {
				c[r]--
				outcome = aliceWins(c, (sum+r)%3, 1-turn)
				c[r]++
			}

			if outcome == aliceMoves {
				result = outcome
				break
			}
		}

		memo[key] = result

		return result
	}

	return aliceWins(cnt, 0, 0)
}
