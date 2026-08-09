// Package climbingstairs — решение leetcode 70. Climbing Stairs.
package climbingstairs

func climbStairs(n int) int {
	if n < 4 {
		return n
	}
	dp := make([]int, n)
	dp[0] = 1
	dp[1] = 2

	for i := 2; i < n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}

	return dp[n-1]

}
