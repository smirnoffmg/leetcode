// Package stonegameiv — решение leetcode 1510. Stone Game IV.
package stonegameiv

const maxN = 100_000

// Позиция «осталось i камней» не зависит от исходного n, поэтому таблица
// переиспользуется между вызовами и достраивается только вверх.
var (
	win      [maxN + 1]bool
	computed int
)

func winnerSquareGame(n int) bool {
	for i := computed + 1; i <= n; i++ {
		for k := 1; k*k <= i; k++ {
			if !win[i-k*k] {
				win[i] = true
				break
			}
		}
	}

	if n > computed {
		computed = n
	}

	return win[n]
}
