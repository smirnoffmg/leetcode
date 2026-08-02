// Package stonegame — решение leetcode 877. Stone Game.
package stonegame

func stoneGame(piles []int) bool {
	// dp[j] — лучшая разница очков (ходящий минус соперник) на отрезке piles[i..j].
	// Строка i зависит только от строки i+1, поэтому хватает одного массива.
	dp := make([]int, len(piles))
	copy(dp, piles)

	for i := len(piles) - 2; i >= 0; i-- {
		for j := i + 1; j < len(piles); j++ {
			// dp[j] ещё от прошлой итерации — это отрезок [i+1..j];
			// dp[j-1] уже обновлён на текущей — это отрезок [i..j-1].
			dp[j] = max(piles[i]-dp[j], piles[j]-dp[j-1])
		}
	}

	return dp[len(piles)-1] > 0
}
