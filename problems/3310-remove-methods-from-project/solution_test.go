package removemethodsfromproject

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestRemainingMethods(t *testing.T) {
	tests := []struct {
		name        string
		n           int
		k           int
		invocations [][]int
		want        []int
	}{
		{name: "example 1", n: 4, k: 1, invocations: [][]int{{1, 2}, {0, 1}, {3, 2}}, want: []int{0, 1, 2, 3}},
		{name: "example 2", n: 5, k: 0, invocations: [][]int{{1, 2}, {0, 2}, {0, 1}, {3, 4}}, want: []int{3, 4}},
		{name: "example 3", n: 3, k: 2, invocations: [][]int{{1, 2}, {0, 1}, {2, 0}}, want: []int{}},

		// вырожденные размеры
		{name: "single method removed", n: 1, k: 0, invocations: [][]int{}, want: []int{}},
		{name: "no invocations at all", n: 4, k: 1, invocations: [][]int{}, want: []int{0, 2, 3}},
		{name: "two methods, bug in caller", n: 2, k: 0, invocations: [][]int{{0, 1}}, want: []int{}},
		{name: "two methods, bug in callee", n: 2, k: 1, invocations: [][]int{{0, 1}}, want: []int{0, 1}},

		// подозрительная группа замкнута — удаляется целиком
		{name: "chain from k", n: 5, k: 0, invocations: [][]int{{0, 1}, {1, 2}, {2, 3}}, want: []int{4}},
		{name: "star from k", n: 4, k: 0, invocations: [][]int{{0, 1}, {0, 2}, {0, 3}}, want: []int{}},
		{name: "diamond from k", n: 5, k: 0, invocations: [][]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}}, want: []int{4}},
		{name: "cycle inside suspicious group", n: 4, k: 0, invocations: [][]int{{0, 1}, {1, 2}, {2, 0}}, want: []int{3}},
		// обратное ребро 1→0 ведёт в саму группу, а не извне
		{name: "back edge into k", n: 3, k: 0, invocations: [][]int{{0, 1}, {1, 0}}, want: []int{2}},

		// вызов извне — не удаляем ничего
		{name: "outsider calls k", n: 3, k: 1, invocations: [][]int{{0, 1}}, want: []int{0, 1, 2}},
		{name: "outsider calls deep callee", n: 4, k: 0, invocations: [][]int{{0, 1}, {1, 2}, {3, 2}}, want: []int{0, 1, 2, 3}},
		// вызов в последнюю вершину длинной цепочки: проверять надо всю группу, а не только k
		{name: "outsider calls tail of chain", n: 6, k: 0, invocations: [][]int{{0, 1}, {1, 2}, {2, 3}, {5, 3}}, want: []int{0, 1, 2, 3, 4, 5}},
		{name: "outsider calls into cycle", n: 4, k: 0, invocations: [][]int{{0, 1}, {1, 0}, {3, 1}}, want: []int{0, 1, 2, 3}},
		{name: "several outside callers", n: 5, k: 2, invocations: [][]int{{0, 2}, {1, 2}, {3, 2}}, want: []int{0, 1, 2, 3, 4}},

		// вызовы вне группы на ответ не влияют
		{name: "outside edges ignored", n: 6, k: 0, invocations: [][]int{{0, 1}, {2, 3}, {3, 4}, {4, 2}}, want: []int{2, 3, 4, 5}},
		{name: "isolated k with busy rest", n: 5, k: 3, invocations: [][]int{{0, 1}, {1, 2}, {4, 0}}, want: []int{0, 1, 2, 4}},
		// подозрительная 1 вызывает 2, но саму 1 никто извне не зовёт
		{name: "suspicious calls shared method", n: 4, k: 1, invocations: [][]int{{1, 2}, {3, 0}}, want: []int{0, 3}},

		// вся программа подозрительна
		{name: "everything suspicious via cycle", n: 3, k: 0, invocations: [][]int{{0, 1}, {1, 2}, {2, 0}}, want: []int{}},
		{name: "everything suspicious via tree", n: 7, k: 0, invocations: [][]int{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 5}, {2, 6}}, want: []int{}},

		// k — последний номер, крайние индексы
		{name: "k is last method", n: 4, k: 3, invocations: [][]int{{3, 0}, {0, 1}}, want: []int{2}},
		{name: "k is last and called from first", n: 4, k: 3, invocations: [][]int{{0, 3}}, want: []int{0, 1, 2, 3}},

		// две компоненты: баг только в одной
		{name: "two components", n: 6, k: 0, invocations: [][]int{{0, 1}, {2, 3}, {3, 4}}, want: []int{2, 3, 4, 5}},
		{name: "two components, outsider in other one", n: 6, k: 0, invocations: [][]int{{0, 1}, {2, 3}, {4, 5}}, want: []int{2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := remainingMethods(tt.n, tt.k, tt.invocations)

			if !equalAsSets(got, tt.want) {
				t.Errorf("remainingMethods(%v, %v, %v) = %v, want %v (в любом порядке)",
					tt.n, tt.k, tt.invocations, got, tt.want)
			}
		})
	}
}

func TestRemainingMethodsMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(3310, 2024))

	for range 2000 {
		n := 1 + rng.IntN(7)
		k := rng.IntN(n)
		invocations := randomInvocations(rng, n)

		want := bruteForceRemaining(n, k, invocations)

		if got := remainingMethods(n, k, invocations); !equalAsSets(got, want) {
			t.Fatalf("remainingMethods(%v, %v, %v) = %v, want %v", n, k, invocations, got, want)
		}
	}
}

func TestRemainingMethodsReturnsValidSet(t *testing.T) {
	rng := rand.New(rand.NewPCG(3310, 42))

	for range 500 {
		n := 1 + rng.IntN(10)
		k := rng.IntN(n)
		invocations := randomInvocations(rng, n)

		got := remainingMethods(n, k, invocations)
		seen := make(map[int]bool, len(got))

		for _, method := range got {
			if method < 0 || method >= n {
				t.Fatalf("remainingMethods(%v, %v, %v) = %v: номер %v вне диапазона", n, k, invocations, got, method)
			}

			if seen[method] {
				t.Fatalf("remainingMethods(%v, %v, %v) = %v: дубликат %v", n, k, invocations, got, method)
			}

			seen[method] = true
		}

		// k подозрителен всегда, значит остаться он может только вместе со всеми
		if seen[k] && len(got) != n {
			t.Fatalf("remainingMethods(%v, %v, %v) = %v: k остался, но удалено не пусто", n, k, invocations, got)
		}
	}
}

func TestRemainingMethodsDoesNotMutateInput(t *testing.T) {
	invocations := [][]int{{1, 2}, {0, 2}, {0, 1}, {3, 4}}
	before := make([][]int, len(invocations))

	for i, inv := range invocations {
		before[i] = slices.Clone(inv)
	}

	remainingMethods(5, 0, invocations)

	if !slices.EqualFunc(invocations, before, slices.Equal) {
		t.Fatalf("вход изменён: got %v, want %v", invocations, before)
	}
}

// Длинная цепочка вызовов: рекурсивный обход на такой глубине рискует переполнить стек.
func TestRemainingMethodsLongChain(t *testing.T) {
	const n = 100_000

	invocations := make([][]int, 0, n-1)
	for i := range n - 1 {
		invocations = append(invocations, []int{i, i + 1})
	}

	if got := remainingMethods(n, 0, invocations); len(got) != 0 {
		t.Fatalf("цепочка из %d методов: осталось %d методов, want 0", n, len(got))
	}

	// один вызов извне в хвост цепочки сохраняет весь проект
	withOutsider := append(slices.Clone(invocations), []int{n - 1, 1})
	if got := remainingMethods(n, 1, withOutsider); len(got) != n {
		t.Fatalf("цепочка с вызовом извне: осталось %d методов, want %d", len(got), n)
	}
}

// Максимальные по условию размеры: n = 10^5, invocations = 2·10^5.
func TestRemainingMethodsMaxSize(t *testing.T) {
	const n = 100_000

	rng := rand.New(rand.NewPCG(3310, 100_000))
	invocations := make([][]int, 0, 2*n)
	seen := make(map[[2]int]bool, 2*n)

	for len(invocations) < 2*n {
		caller, callee := rng.IntN(n), rng.IntN(n)
		if caller == callee || seen[[2]int{caller, callee}] {
			continue
		}

		seen[[2]int{caller, callee}] = true
		invocations = append(invocations, []int{caller, callee})
	}

	// в случайном плотном графе почти наверняка найдётся вызов извне
	if got := remainingMethods(n, 0, invocations); len(got) != n {
		t.Fatalf("плотный граф: осталось %d методов, want %d", len(got), n)
	}
}

func BenchmarkRemainingMethods(b *testing.B) {
	const n = 100_000

	invocations := make([][]int, 0, 2*n)
	for i := range n - 1 {
		invocations = append(invocations, []int{i, i + 1})
	}

	b.ResetTimer()

	for b.Loop() {
		remainingMethods(n, 0, invocations)
	}
}

// randomInvocations выбирает подмножество различных упорядоченных пар (a, b), a != b.
func randomInvocations(rng *rand.Rand, n int) [][]int {
	invocations := make([][]int, 0, n*(n-1)/3)

	for caller := range n {
		for callee := range n {
			if caller != callee && rng.IntN(3) == 0 {
				invocations = append(invocations, []int{caller, callee})
			}
		}
	}

	return invocations
}

// bruteForceRemaining — эталон: достижимость из k ищется наивным пересчётом до
// неподвижной точки, независимо от обхода в решении.
func bruteForceRemaining(n int, k int, invocations [][]int) []int {
	suspicious := make([]bool, n)
	suspicious[k] = true

	for changed := true; changed; {
		changed = false

		for _, inv := range invocations {
			if suspicious[inv[0]] && !suspicious[inv[1]] {
				suspicious[inv[1]] = true
				changed = true
			}
		}
	}

	for _, inv := range invocations {
		if !suspicious[inv[0]] && suspicious[inv[1]] {
			suspicious = make([]bool, n)
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

// equalAsSets сравнивает ответы без учёта порядка: условие разрешает любой.
func equalAsSets(got, want []int) bool {
	gotSorted, wantSorted := slices.Clone(got), slices.Clone(want)
	slices.Sort(gotSorted)
	slices.Sort(wantSorted)

	return slices.Equal(gotSorted, wantSorted)
}
