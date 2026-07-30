// Package minimumnumberofpushestotypewordi — решение leetcode 3014. Minimum Number of Pushes to Type Word I.
package minimumnumberofpushestotypewordi

func minimumPushes(word string) int {
	var res int

	for l := len(word); l > 0; l -= 8 {
		res += l
	}

	return res
}
