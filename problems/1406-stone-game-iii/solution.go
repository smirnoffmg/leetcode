// Package stonegameiii — решение leetcode 1406. Stone Game III.
package stonegameiii

func stoneGameIII(stoneValue []int) string {
	n := len(stoneValue)
	// dp[i] — лучшая разница очков (ходящий минус соперник) на суффиксе stoneValue[i:];
	// dp[n] = 0 — камней не осталось.
	dp := make([]int, n+1)

	for i := n - 1; i >= 0; i-- {
		taken := stoneValue[i]
		best := taken - dp[i+1]

		for k := 1; k < 3 && i+k < n; k++ {
			taken += stoneValue[i+k]
			best = max(best, taken-dp[i+k+1])
		}

		dp[i] = best
	}

	switch {
	case dp[0] > 0:
		return "Alice"
	case dp[0] < 0:
		return "Bob"
	default:
		return "Tie"
	}
}
