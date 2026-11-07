// Package processtasksusingservers — решение leetcode 1882. Process Tasks Using Servers.
package processtasksusingservers

// Ключ кучи — (приоритет << indexBits) | индекс сервера: сравнение одного int сразу даёт
// порядок «приоритет, затем индекс». Индекс < 2^18, вес <= 2*10^5 < 2^18, а время
// освобождения <= 2*10^5 + 2*10^5 * 2*10^5 < 2^36, так что ключ укладывается в 54 бита.
const (
	indexBits = 18
	indexMask = 1<<indexBits - 1
)

func assignTasks(servers []int, tasks []int) []int {
	free := make([]int, len(servers))
	for i, weight := range servers {
		free[i] = weight<<indexBits | i
	}

	heapify(free)

	busy := make([]int, 0, len(servers))
	ans := make([]int, len(tasks))
	now := 0

	for j, duration := range tasks {
		now = max(now, j)
		// Без max: все в busy освобождаются позже прошлого now, а j <= прошлый now + 1,
		// поэтому время здесь назад не идёт.
		if len(free) == 0 {
			now = busy[0] >> indexBits
		}

		for len(busy) > 0 && busy[0]>>indexBits <= now {
			server := busy[0] & indexMask
			busy = popMin(busy)
			free = pushMin(free, servers[server]<<indexBits|server)
		}

		server := free[0] & indexMask
		free = popMin(free)
		busy = pushMin(busy, (now+duration)<<indexBits|server)
		ans[j] = server
	}

	return ans
}

func heapify(h []int) {
	for i := len(h)/2 - 1; i >= 0; i-- {
		siftDown(h, i)
	}
}

func pushMin(h []int, key int) []int {
	h = append(h, key)
	i := len(h) - 1

	for i > 0 {
		parent := (i - 1) / 2
		if h[parent] <= key {
			break
		}

		h[i] = h[parent]
		i = parent
	}

	h[i] = key

	return h
}

func popMin(h []int) []int {
	last := len(h) - 1
	h[0] = h[last]
	h = h[:last]
	siftDown(h, 0)

	return h
}

func siftDown(h []int, i int) {
	if i >= len(h) {
		return
	}

	key := h[i]

	for {
		child := 2*i + 1
		if child >= len(h) {
			break
		}

		if right := child + 1; right < len(h) && h[right] < h[child] {
			child = right
		}

		if h[child] >= key {
			break
		}

		h[i] = h[child]
		i = child
	}

	h[i] = key
}
