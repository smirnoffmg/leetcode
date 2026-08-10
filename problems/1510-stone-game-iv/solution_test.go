package stonegameiv

import (
	"testing"
)

func TestWinnerSquareGame(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want bool
	}{
		{name: "example 1", n: 1, want: true},
		{name: "example 2", n: 2, want: false},
		{name: "example 3", n: 4, want: true},
		{name: "greedy trap", n: 8, want: true},
		{name: "losing 7", n: 7, want: false},
		{name: "losing 17", n: 17, want: false},
		{name: "gap before 34", n: 22, want: false},
		{name: "between 22 and 34", n: 30, want: true},
		{name: "losing 34", n: 34, want: false},
		{name: "max n", n: 100000, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := winnerSquareGame(tt.n)
			if got != tt.want {
				t.Errorf("winnerSquareGame(%v) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

func TestWinnerSquareGameLosingPositions(t *testing.T) {
	// A030193: проигрышные позиции игры «вычти квадрат».
	losing := map[int]bool{
		0: true, 2: true, 5: true, 7: true, 10: true, 12: true, 15: true,
		17: true, 20: true, 22: true, 34: true, 39: true, 44: true, 52: true,
		57: true, 62: true, 65: true, 67: true, 72: true, 85: true, 95: true,
	}

	for n := 1; n <= 100; n++ {
		if got := winnerSquareGame(n); got == losing[n] {
			t.Fatalf("winnerSquareGame(%d) = %v, want %v", n, got, !losing[n])
		}
	}
}

// bruteForce — определение задачи без мемоизации: ходящий выигрывает,
// если существует ход в позицию, проигрышную для соперника.
func bruteForce(n int) bool {
	for k := 1; k*k <= n; k++ {
		if !bruteForce(n - k*k) {
			return true
		}
	}

	return false
}

func TestWinnerSquareGameAgainstBruteForce(t *testing.T) {
	for n := 0; n <= 28; n++ {
		want := bruteForce(n)
		if got := winnerSquareGame(n); got != want {
			t.Fatalf("winnerSquareGame(%d) = %v, brute force = %v", n, got, want)
		}
	}
}

// Таблица живёт между вызовами: порядок аргументов не должен на неё влиять.
func TestWinnerSquareGameTableReuse(t *testing.T) {
	first := make([]bool, 201)
	for n := 200; n >= 0; n-- {
		first[n] = winnerSquareGame(n)
	}

	for n := 0; n <= 200; n++ {
		if got := winnerSquareGame(n); got != first[n] {
			t.Fatalf("winnerSquareGame(%d) = %v on second call, %v on first", n, got, first[n])
		}
	}
}

func BenchmarkWinnerSquareGame(b *testing.B) {
	for b.Loop() {
		computed = 0
		win = [maxN + 1]bool{}
		winnerSquareGame(maxN)
	}
}
