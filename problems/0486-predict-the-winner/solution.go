// Package predictthewinner — решение leetcode 486. Predict the Winner.
package predictthewinner

func predictTheWinner(nums []int) bool {
	// dp[j] — лучшая разница очков (ходящий минус соперник) на отрезке nums[i..j].
	// Строка i зависит только от строки i+1, поэтому хватает одного массива.
	dp := make([]int, len(nums))
	copy(dp, nums)

	for i := len(nums) - 2; i >= 0; i-- {
		for j := i + 1; j < len(nums); j++ {
			// dp[j] ещё от прошлой итерации — это отрезок [i+1..j];
			// dp[j-1] уже обновлён на текущей — это отрезок [i..j-1].
			dp[j] = max(nums[i]-dp[j], nums[j]-dp[j-1])
		}
	}

	return dp[len(nums)-1] >= 0
}
