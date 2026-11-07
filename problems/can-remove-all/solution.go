// Package canremoveall — решение контестной задачи LeetCode об удалении групп элементов.
package canremoveall

// Тип < 10^5 и цвет <= 10^5 меньше 2^17, поэтому пара (тип, цвет) упаковывается без коллизий.
const colorBits = 17

type run struct{ key, count int }

func canRemoveAll(types []int, colors []int, groupSize []int) bool {
	stack := make([]run, 0, len(types))

	for i, kind := range types {
		key := kind<<colorBits | colors[i]

		if top := len(stack) - 1; top >= 0 && stack[top].key == key {
			stack[top].count++
		} else {
			stack = append(stack, run{key: key, count: 1})
		}

		if stack[len(stack)-1].count == groupSize[kind] {
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
