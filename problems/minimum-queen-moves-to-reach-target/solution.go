// Package minimumqueenmovestoreachtarget — решение leetcode Minimum Queen Moves to Reach Target.
package minimumqueenmovestoreachtarget

func minQueenMoves(source []int, target []int) int {
	sr, sc, tr, tc := source[0], source[1], target[0], target[1]

	switch {
	case sr == tr && sc == tc:
		return 0
	case sr == tr, sc == tc, sr-sc == tr-tc, sr+sc == tr+tc:
		return 1
	default:
		return 2
	}
}
