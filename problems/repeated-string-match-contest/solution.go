// Package repeatedstringmatchcontest — решение контестной задачи LeetCode о вхождении b*y в a*x.
package repeatedstringmatchcontest

import "strings"

func repeatedStringMatch(a string, b string, x int, y int) bool {
	total, pattern := len(a)*x, len(b)*y
	if pattern > total {
		return false
	}

	// Файн — Вильф: строки с периодами |a| и |b|, совпавшие на |a|+|b|-gcd символах,
	// совпадают и дальше, поэтому сравнивать все 10^11 символов не нужно.
	compared := min(pattern, len(a)+len(b)-gcd(len(a), len(b)))

	// Сдвиг на |a| даёт ту же строку, поэтому хватает сдвигов меньше |a|. Текст обрезан
	// так, чтобы найденный сдвиг не превысил total-pattern: образец усечён, и без этого
	// вхождение могло бы оказаться там, где полному b*y уже не хватает места.
	textLen := min(len(a)-1, total-pattern) + compared
	text := strings.Repeat(a, textLen/len(a)+1)[:textLen]
	needle := strings.Repeat(b, compared/len(b)+1)[:compared]

	return strings.Contains(text, needle)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}
