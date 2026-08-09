// Package stonegameii — решение leetcode 1140. Stone Game II.
package stonegameii

const maxPiles = 100

var (
	suffix [maxPiles + 1]int
	dp     [maxPiles + 1][maxPiles + 1]int
	// поколение записи в dp: номер вызова вместо обнуления таблицы между вызовами
	seen       [maxPiles + 1][maxPiles + 1]int
	generation int
)

func stoneGameII(piles []int) int {
	n := len(piles)

	generation++

	suffix[n] = 0
	for i := n - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] + piles[i]
	}

	return solve(0, 1, n)
}

// solve — максимум камней ходящего на суффиксе piles[i:] при текущем M = m.
func solve(i, m, n int) int {
	// одного хода хватает на весь остаток; при i == n отсюда же приходит ноль
	if 2*m >= n-i {
		return suffix[i]
	}

	if seen[i][m] == generation {
		return dp[i][m]
	}

	best := 0

	for x := 1; x <= 2*m; x++ {
		if got := suffix[i] - solve(i+x, max(m, x), n); got > best {
			best = got
		}
	}

	dp[i][m], seen[i][m] = best, generation

	return best
}
