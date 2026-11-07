// Package countvalidsubstringssubsequence — решение контестной задачи LeetCode о подстроках без подпоследовательности "abc".
package countvalidsubstringssubsequence

// starts[k] — самый правый старт префикса "abc" длины k; пустой префикс начинается
// в самом right, поэтому переход един для всех букв и обходится без ветвлений.
func countValidSubstrings(s string) int64 {
	var total int64

	starts := [4]int{-1, -1, -1, -1}

	for right := range s {
		starts[0] = right
		letter := s[right] - 'a'
		starts[letter+1] = starts[letter]
		total += int64(right - starts[3])
	}

	return total
}
