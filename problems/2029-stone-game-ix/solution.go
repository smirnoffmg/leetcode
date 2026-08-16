// Package stonegameix — решение leetcode 2029. Stone Game IX.
package stonegameix

func stoneGameIX(stones []int) bool {
	var cnt [3]int
	for _, v := range stones {
		cnt[v%3]++
	}

	if cnt[0]%2 == 0 {
		return cnt[1] > 0 && cnt[2] > 0
	}

	return cnt[1] > cnt[2]+2 || cnt[2] > cnt[1]+2
}
