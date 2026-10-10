// Package climbingstairs — решение leetcode 70. Climbing Stairs.
package climbingstairs

func climbStairs(n int) int {
	prev, cur := 1, 1

	for range n - 1 {
		prev, cur = cur, prev+cur
	}

	return cur
}
