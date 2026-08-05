// Package removemethodsfromproject — решение leetcode 3310. Remove Methods From Project.
package removemethodsfromproject

func remainingMethods(n int, k int, invocations [][]int) []int {
	starts, targets := buildCSR(n, invocations)

	// очередь только растёт, обход — курсор по ней: снимать вершины не нужно,
	// нас интересует само множество достижимых, а не порядок их посещения
	suspicious := make([]bool, n)
	suspicious[k] = true
	queue := make([]int, 1, n)
	queue[0] = k

	for i := 0; i < len(queue); i++ {
		caller := queue[i]

		for _, callee := range targets[starts[caller]:starts[caller+1]] {
			if !suspicious[callee] {
				suspicious[callee] = true
				queue = append(queue, callee)
			}
		}
	}

	// вызов из-вне группы запрещает удаление целиком, а «удалять нечего»
	// и «подозрительных нет» дают один и тот же ответ — все методы
	for _, inv := range invocations {
		caller, callee := inv[0], inv[1]
		if !suspicious[caller] && suspicious[callee] {
			clear(suspicious)

			break
		}
	}

	remaining := make([]int, 0, n)

	for method := range n {
		if !suspicious[method] {
			remaining = append(remaining, method)
		}
	}

	return remaining
}

// buildCSR укладывает граф вызовов в две плоские таблицы вместо среза срезов:
// вызываемые методом v лежат в targets[starts[v]:starts[v+1]].
func buildCSR(n int, invocations [][]int) (starts, targets []int) {
	starts = make([]int, n+1)
	for _, inv := range invocations {
		starts[inv[0]+1]++
	}

	for caller := range n {
		starts[caller+1] += starts[caller]
	}

	targets = make([]int, len(invocations))
	for _, inv := range invocations {
		caller, callee := inv[0], inv[1]
		targets[starts[caller]] = callee
		starts[caller]++
	}

	// раскладка использовала starts как курсоры и сдвинула их на длину блока:
	// теперь starts[v] — конец блока v, то есть прежнее starts[v+1]
	for caller := n; caller > 0; caller-- {
		starts[caller] = starts[caller-1]
	}

	starts[0] = 0

	return starts, targets
}
