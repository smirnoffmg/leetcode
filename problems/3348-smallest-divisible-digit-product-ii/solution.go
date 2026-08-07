// Package smallestdivisibledigitproductii — решение leetcode 3348. Smallest Divisible Digit Product II.
package smallestdivisibledigitproductii

import "slices"

// exps[d] — степени простых 2, 3, 5, 7 в разложении цифры d; нулевая строка для 0 не используется.
var exps = [10][4]int{
	1: {0, 0, 0, 0},
	2: {1, 0, 0, 0},
	3: {0, 1, 0, 0},
	4: {2, 0, 0, 0},
	5: {0, 0, 1, 0},
	6: {1, 1, 0, 0},
	7: {0, 0, 0, 1},
	8: {3, 0, 0, 0},
	9: {0, 2, 0, 0},
}

// minLen23 — минимальное число цифр из {2,3,4,6,8,9}, произведение которых делится на 2^a * 3^b.
// База — восьмёрки и девятки, остатки a%3 и b%2 закрываются одной цифрой (2, 4, 3 или 6),
// кроме пары остатков 2^2 * 3: её одной цифрой не покрыть.
func minLen23(a, b int) int {
	n := a/3 + b/2

	switch {
	case a%3 == 0 && b%2 == 0:
	case a%3 == 2 && b%2 == 1:
		n += 2
	default:
		n++
	}

	return n
}

func minLen(need [4]int) int {
	return minLen23(need[0], need[1]) + need[2] + need[3]
}

func sub(need, have [4]int) [4]int {
	for j := range need {
		need[j] = max(0, need[j]-have[j])
	}

	return need
}

// fill — наименьшая строка ровно из slots ненулевых цифр, произведение которых делится
// на 2^need[0] * 3^need[1] * 5^need[2] * 7^need[3]. Требует minLen(need) <= slots.
func fill(need [4]int, slots int) []byte {
	res := make([]byte, 0, slots)
	for range slots - minLen(need) {
		res = append(res, '1')
	}

	// Жадно по возрастанию: цифра годится, если после неё остаток всё ещё
	// укладывается в оставшиеся позиции. Каждая взятая цифра уменьшает minLen23
	// ровно на 1, поэтому цикл сходится за minLen23 шагов.
	a, b := need[0], need[1]
	for rem := minLen23(a, b); rem > 0; rem-- {
		for _, d := range [...]byte{2, 3, 4, 6, 8, 9} {
			na := max(0, a-exps[d][0])
			nb := max(0, b-exps[d][1])

			if minLen23(na, nb) < rem {
				a, b = na, nb
				res = append(res, '0'+d)

				break
			}
		}
	}

	for range need[2] {
		res = append(res, '5')
	}
	for range need[3] {
		res = append(res, '7')
	}

	slices.Sort(res)

	return res
}

func smallestNumber(num string, t int64) string {
	var need [4]int
	for j, p := range [...]int64{2, 3, 5, 7} {
		for t%p == 0 {
			t /= p
			need[j]++
		}
	}

	if t > 1 {
		return "-1"
	}

	n := len(num)
	pre := make([][4]int, n+1)
	zeroIdx := n

	for i := range n {
		d := num[i] - '0'
		if d == 0 && zeroIdx == n {
			zeroIdx = i
		}

		pre[i+1] = pre[i]
		for j := range 4 {
			pre[i+1][j] += exps[d][j]
		}
	}

	if zeroIdx == n && minLen(sub(need, pre[n])) == 0 {
		return num
	}

	// Ответ той же длины: сохраняем максимально длинный префикс num (без нулей),
	// в позиции i ставим цифру больше num[i], хвост добиваем минимальной строкой.
	// Чем позже расходимся с num, тем меньше результат.
	for i := min(n-1, zeroIdx); i >= 0; i-- {
		slots := n - 1 - i
		for d := num[i] - '0' + 1; d <= 9; d++ {
			rem := sub(sub(need, pre[i]), exps[d])
			if minLen(rem) <= slots {
				res := make([]byte, 0, n)
				res = append(res, num[:i]...)
				res = append(res, '0'+d)
				res = append(res, fill(rem, slots)...)

				return string(res)
			}
		}
	}

	// Длиннее num: любое число из n+1 и более ненулевых цифр уже больше num,
	// поэтому берём кратчайшую длину и минимальное заполнение.
	return string(fill(need, max(n+1, minLen(need))))
}
