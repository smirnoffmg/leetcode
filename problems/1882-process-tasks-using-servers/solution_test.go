package processtasksusingservers

import (
	"container/heap"
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"
)

const maxConstraint = 200_000

func TestAssignTasks(t *testing.T) {
	tests := []struct {
		name    string
		servers []int
		tasks   []int
		want    []int
	}{
		{name: "example 1", servers: []int{3, 3, 2}, tasks: []int{1, 2, 3, 2, 1, 2}, want: []int{2, 2, 0, 2, 1, 2}},
		{name: "example 2", servers: []int{5, 1, 4, 3, 2}, tasks: []int{2, 1, 2, 4, 5, 2, 1}, want: []int{1, 4, 1, 4, 1, 3, 2}},

		{name: "single server queues everything", servers: []int{7}, tasks: []int{5, 1, 1}, want: []int{0, 0, 0}},
		{name: "more servers than tasks", servers: []int{9, 8, 7}, tasks: []int{1}, want: []int{2}},
		{name: "equal weights pick smallest index", servers: []int{2, 2}, tasks: []int{1, 1}, want: []int{0, 0}},
		{name: "busy lighter server is skipped", servers: []int{1, 2}, tasks: []int{2, 1}, want: []int{0, 1}},
		// задача 2 ждёт сервер 0 до секунды 3, задача 3 — до секунды 4, когда свободны оба
		{name: "queue waits for each release", servers: []int{1, 2}, tasks: []int{3, 3, 1, 1}, want: []int{0, 1, 0, 0}},
		// оба освобождаются в секунду 5: задачи из очереди идут по весу, а не по индексу сервера
		{name: "simultaneous release serves queue by weight", servers: []int{2, 1}, tasks: []int{5, 4, 1, 1}, want: []int{1, 0, 1, 0}},
		// задача 2 ждёт до секунды 10, задача 3 — до 11, когда свободны оба
		{name: "time jumps while queue waits", servers: []int{2, 1}, tasks: []int{10, 10, 1, 1}, want: []int{1, 0, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assignTasks(tt.servers, tt.tasks)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("assignTasks(%v, %v) = %v, want %v", tt.servers, tt.tasks, got, tt.want)
			}
		})
	}
}

func TestAssignTasksMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(1882, 2021))

	for range 20_000 {
		servers := randomInts(rng, 1+rng.IntN(5), 3)
		tasks := randomInts(rng, 1+rng.IntN(9), 4)

		want := bruteForceAssign(servers, tasks)
		if got := assignTasks(servers, tasks); !reflect.DeepEqual(got, want) {
			t.Fatalf("assignTasks(%v, %v) = %v, want %v", servers, tasks, got, want)
		}
	}
}

// Сверка с вариантом на структурах, где ключи не упакованы: на предельных значениях
// переполнение упакованного ключа дало бы расхождение.
func TestAssignTasksAgreesWithUnpackedOnLimits(t *testing.T) {
	rng := rand.New(rand.NewPCG(1882, 7))

	inputs := []struct {
		name           string
		servers, tasks []int
	}{
		{name: "one server, longest total time", servers: []int{maxConstraint}, tasks: constantInts(maxConstraint, maxConstraint)},
		{name: "max equal weights", servers: constantInts(maxConstraint, maxConstraint), tasks: constantInts(maxConstraint, maxConstraint)},
		{name: "two servers, longest tasks", servers: []int{maxConstraint, maxConstraint - 1}, tasks: constantInts(maxConstraint, maxConstraint)},
		{name: "random", servers: randomInts(rng, maxConstraint, maxConstraint), tasks: randomInts(rng, maxConstraint, maxConstraint)},
	}

	for _, in := range inputs {
		t.Run(in.name, func(t *testing.T) {
			want := assignTasksContainerStruct(in.servers, in.tasks)
			if got := assignTasks(in.servers, in.tasks); !slices.Equal(got, want) {
				t.Fatalf("assignTasks расходится с вариантом на структурах")
			}

			if got := assignTasksContainerPacked(in.servers, in.tasks); !slices.Equal(got, want) {
				t.Fatalf("assignTasksContainerPacked расходится с вариантом на структурах")
			}
		})
	}
}

func TestAssignTasksDoesNotMutateInput(t *testing.T) {
	servers, tasks := []int{3, 1, 2}, []int{2, 1, 3, 1}
	serversCopy, tasksCopy := slices.Clone(servers), slices.Clone(tasks)

	assignTasks(servers, tasks)

	if !slices.Equal(servers, serversCopy) || !slices.Equal(tasks, tasksCopy) {
		t.Fatalf("вход изменён: servers = %v, tasks = %v", servers, tasks)
	}
}

func BenchmarkAssignTasks(b *testing.B) {
	rng := rand.New(rand.NewPCG(1882, 1))

	shapes := []struct {
		name           string
		servers, tasks []int
	}{
		// длительности со средним около n/2 держат обе кучи глубокими одновременно
		{name: "shape=deep-heaps", servers: randomInts(rng, maxConstraint, maxConstraint), tasks: randomInts(rng, maxConstraint, maxConstraint)},
		{name: "shape=short-tasks", servers: randomInts(rng, maxConstraint, maxConstraint), tasks: randomInts(rng, maxConstraint, 4)},
		{name: "shape=few-servers", servers: randomInts(rng, 16, maxConstraint), tasks: randomInts(rng, maxConstraint, maxConstraint)},
	}

	impls := []struct {
		name string
		fn   func(servers, tasks []int) []int
	}{
		{name: "impl=manual-packed", fn: assignTasks},
		{name: "impl=container-packed", fn: assignTasksContainerPacked},
		{name: "impl=container-struct", fn: assignTasksContainerStruct},
	}

	for _, shape := range shapes {
		for _, impl := range impls {
			b.Run(shape.name+"/"+impl.name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					impl.fn(shape.servers, shape.tasks)
				}
			})
		}
	}
}

// bruteForceAssign проигрывает условие посекундно, перебирая все серверы на каждое назначение.
func bruteForceAssign(servers, tasks []int) []int {
	freeAt := make([]int, len(servers))
	ans := make([]int, 0, len(tasks))

	for second := 0; len(ans) < len(tasks); second++ {
		inserted := min(second+1, len(tasks))

		for len(ans) < inserted {
			best := -1

			for i, weight := range servers {
				if freeAt[i] <= second && (best == -1 || weight < servers[best]) {
					best = i
				}
			}

			if best == -1 {
				break
			}

			freeAt[best] = second + tasks[len(ans)]
			ans = append(ans, best)
		}
	}

	return ans
}

func randomInts(rng *rand.Rand, n, maxValue int) []int {
	values := make([]int, n)
	for i := range values {
		values[i] = 1 + rng.IntN(maxValue)
	}

	return values
}

func constantInts(n, value int) []int {
	values := make([]int, n)
	for i := range values {
		values[i] = value
	}

	return values
}

type packedHeap []int

func (h packedHeap) Len() int           { return len(h) }
func (h packedHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h packedHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *packedHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *packedHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]

	return last
}

func assignTasksContainerPacked(servers, tasks []int) []int {
	free := make(packedHeap, len(servers))
	for i, weight := range servers {
		free[i] = weight<<indexBits | i
	}

	heap.Init(&free)

	busy := make(packedHeap, 0, len(servers))
	ans := make([]int, len(tasks))
	now := 0

	for j, duration := range tasks {
		now = max(now, j)
		if len(free) == 0 {
			now = busy[0] >> indexBits
		}

		for len(busy) > 0 && busy[0]>>indexBits <= now {
			server := heap.Pop(&busy).(int) & indexMask
			heap.Push(&free, servers[server]<<indexBits|server)
		}

		server := heap.Pop(&free).(int) & indexMask
		heap.Push(&busy, (now+duration)<<indexBits|server)
		ans[j] = server
	}

	return ans
}

type serverEntry struct{ key, index int }

type entryHeap []serverEntry

func (h entryHeap) Len() int { return len(h) }
func (h entryHeap) Less(i, j int) bool {
	if h[i].key != h[j].key {
		return h[i].key < h[j].key
	}

	return h[i].index < h[j].index
}
func (h entryHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *entryHeap) Push(x any)   { *h = append(*h, x.(serverEntry)) }
func (h *entryHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]

	return last
}

func assignTasksContainerStruct(servers, tasks []int) []int {
	free := make(entryHeap, len(servers))
	for i, weight := range servers {
		free[i] = serverEntry{key: weight, index: i}
	}

	heap.Init(&free)

	busy := make(entryHeap, 0, len(servers))
	ans := make([]int, len(tasks))
	now := 0

	for j, duration := range tasks {
		now = max(now, j)
		if len(free) == 0 {
			now = busy[0].key
		}

		for len(busy) > 0 && busy[0].key <= now {
			server := heap.Pop(&busy).(serverEntry).index
			heap.Push(&free, serverEntry{key: servers[server], index: server})
		}

		server := heap.Pop(&free).(serverEntry).index
		heap.Push(&busy, serverEntry{key: now + duration, index: server})
		ans[j] = server
	}

	return ans
}
