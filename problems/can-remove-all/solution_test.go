package canremoveall

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

const maxConstraint = 100_000

func TestCanRemoveAll(t *testing.T) {
	tests := []struct {
		name      string
		types     []int
		colors    []int
		groupSize []int
		want      bool
	}{
		{name: "example 1", types: []int{0, 1, 1, 0}, colors: []int{5, 3, 3, 5}, groupSize: []int{2, 2}, want: true},
		{name: "example 2", types: []int{0, 1, 0}, colors: []int{9, 9, 9}, groupSize: []int{3, 2}, want: false},
		{name: "example 3", types: []int{2, 0, 1, 2}, colors: []int{7, 5, 4, 7}, groupSize: []int{3, 1, 2}, want: false},

		{name: "single element, group of one", types: []int{0}, colors: []int{1}, groupSize: []int{1}, want: true},
		{name: "single element, group of two", types: []int{0}, colors: []int{1}, groupSize: []int{2}, want: false},
		{name: "same type, different colors", types: []int{0, 0}, colors: []int{1, 2}, groupSize: []int{2}, want: false},
		{name: "different types, same color", types: []int{0, 1}, colors: []int{5, 5}, groupSize: []int{2, 2}, want: false},
		{name: "run length is not a multiple", types: []int{0, 0, 0}, colors: []int{1, 1, 1}, groupSize: []int{2}, want: false},
		{name: "run length is a multiple", types: []int{0, 0, 0, 0}, colors: []int{1, 1, 1, 1}, groupSize: []int{2}, want: true},
		{name: "group of one removes everything", types: []int{0, 0, 0}, colors: []int{1, 2, 3}, groupSize: []int{1}, want: true},
		// x (y y) x x: части группы x смыкаются только после удаления y
		{name: "group closes after inner removal", types: []int{0, 1, 1, 0, 0}, colors: []int{1, 1, 1, 1, 1}, groupSize: []int{3, 2}, want: true},
		{name: "two parts before inner group", types: []int{0, 0, 1, 1, 0}, colors: []int{1, 1, 1, 1, 1}, groupSize: []int{3, 2}, want: true},
		{name: "group of one between parts", types: []int{0, 1, 0}, colors: []int{1, 1, 1}, groupSize: []int{2, 1}, want: true},
		// x y (x x) y → x y y → x: одна x остаётся без пары
		{name: "leftover after cascade", types: []int{0, 1, 0, 0, 1}, colors: []int{1, 1, 1, 1, 1}, groupSize: []int{2, 2}, want: false},
		{name: "balanced counts in wrong order", types: []int{0, 1, 0, 1}, colors: []int{1, 1, 1, 1}, groupSize: []int{2, 2}, want: false},
		{name: "limits of type and color", types: []int{maxConstraint - 1, maxConstraint - 1}, colors: []int{maxConstraint, maxConstraint}, groupSize: limitGroupSizes(), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canRemoveAll(tt.types, tt.colors, tt.groupSize); got != tt.want {
				t.Errorf("canRemoveAll(%v, %v, %d групп) = %v, want %v", tt.types, tt.colors, len(tt.groupSize), got, tt.want)
			}
		})
	}
}

func TestCanRemoveAllMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2026, 1108))

	trueCount, falseCount := 0, 0

	for range 30_000 {
		types, colors, groupSize := randomTinyInput(rng)
		want := bruteForceCanRemove(types, colors, groupSize)

		if want {
			trueCount++
		} else {
			falseCount++
		}

		for _, impl := range implementations() {
			if got := impl.fn(types, colors, groupSize); got != want {
				t.Fatalf("%s(%v, %v, %v) = %v, want %v", impl.name, types, colors, groupSize, got, want)
			}
		}
	}

	// без этого перекос генератора в один ответ сделал бы сверку бессмысленной
	if trueCount < 5000 || falseCount < 5000 {
		t.Fatalf("генератор перекошен: true = %d, false = %d", trueCount, falseCount)
	}
}

func TestImplementationsAgreeOnMaxInputs(t *testing.T) {
	for _, input := range benchInputs() {
		t.Run(input.name, func(t *testing.T) {
			if got := canRemoveAll(input.types, input.colors, input.groupSize); got != input.want {
				t.Fatalf("canRemoveAll = %v, want %v", got, input.want)
			}

			for _, impl := range implementations() {
				if got := impl.fn(input.types, input.colors, input.groupSize); got != input.want {
					t.Fatalf("%s = %v, want %v", impl.name, got, input.want)
				}
			}
		})
	}
}

func TestCanRemoveAllDoesNotMutateInput(t *testing.T) {
	types, colors, groupSize := []int{0, 1, 1, 0}, []int{5, 3, 3, 5}, []int{2, 2}
	typesCopy, colorsCopy, groupSizeCopy := slices.Clone(types), slices.Clone(colors), slices.Clone(groupSize)

	canRemoveAll(types, colors, groupSize)

	if !slices.Equal(types, typesCopy) || !slices.Equal(colors, colorsCopy) || !slices.Equal(groupSize, groupSizeCopy) {
		t.Fatalf("вход изменён: %v, %v, %v", types, colors, groupSize)
	}
}

func BenchmarkCanRemoveAll(b *testing.B) {
	for _, input := range benchInputs() {
		for _, impl := range implementations() {
			b.Run(input.name+"/"+impl.name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					impl.fn(input.types, input.colors, input.groupSize)
				}
			})
		}
	}
}

// bruteForceCanRemove пробует все порядки удалений из условия, запоминая уже виденные
// последовательности. Порядок здесь не зафиксирован, поэтому эталон проверяет и то,
// что жадное удаление в стеке ничего не теряет.
func bruteForceCanRemove(types, colors, groupSize []int) bool {
	seen := map[string]bool{}

	var search func(types, colors []int) bool

	search = func(types, colors []int) bool {
		if len(types) == 0 {
			return true
		}

		state := fmt.Sprint(types, colors)
		if seen[state] {
			return false
		}

		seen[state] = true

		for start := range types {
			end := start + groupSize[types[start]]
			if end > len(types) || !sameGroup(types[start:end], colors[start:end]) {
				continue
			}

			if search(slices.Delete(slices.Clone(types), start, end), slices.Delete(slices.Clone(colors), start, end)) {
				return true
			}
		}

		return false
	}

	return search(types, colors)
}

func sameGroup(types, colors []int) bool {
	for i := range types {
		if types[i] != types[0] || colors[i] != colors[0] {
			return false
		}
	}

	return true
}

// randomTinyInput вставляет целые группы в случайные места: такие входы заведомо
// удаляются полностью. Случайная порча половины входов даёт и отрицательные примеры,
// а равномерно случайные последовательности почти всегда были бы false.
func randomTinyInput(rng *rand.Rand) (types, colors, groupSize []int) {
	groupSize = make([]int, 1+rng.IntN(3))
	for i := range groupSize {
		groupSize[i] = 1 + rng.IntN(3)
	}

	for range 1 + rng.IntN(5) {
		kind, color := rng.IntN(len(groupSize)), 1+rng.IntN(2)
		at := rng.IntN(len(types) + 1)

		for range groupSize[kind] {
			types = slices.Insert(types, at, kind)
			colors = slices.Insert(colors, at, color)
		}
	}

	if rng.IntN(2) == 0 {
		at := rng.IntN(len(types))
		types[at] = rng.IntN(len(groupSize))
		colors[at] = 1 + rng.IntN(2)
	}

	return types, colors, groupSize
}

func limitGroupSizes() []int {
	sizes := make([]int, maxConstraint)
	for i := range sizes {
		sizes[i] = 2
	}

	return sizes
}

type implementation struct {
	name string
	fn   func(types, colors, groupSize []int) bool
}

func implementations() []implementation {
	return []implementation{
		{name: "impl=runs-packed", fn: canRemoveAll},
		{name: "impl=runs-packed-grow", fn: canRemoveAllGrow},
		{name: "impl=runs-struct", fn: canRemoveAllStruct},
		{name: "impl=elements", fn: canRemoveAllElements},
	}
}

type benchInput struct {
	name                     string
	types, colors, groupSize []int
	want                     bool
}

func benchInputs() []benchInput {
	return []benchInput{
		nestedPairsInput(),
		allSameInput(),
		tripletsInput(),
		randomStuckInput(),
	}
}

// nestedPairsInput: x0 x1 ... x1 x0 — стек растёт до n/2 и схлопывается целиком.
func nestedPairsInput() benchInput {
	input := newBenchInput("shape=nested-pairs", true)

	half := maxConstraint / 2
	for i := range half {
		input.types[i], input.types[maxConstraint-1-i] = i, i
		input.colors[i], input.colors[maxConstraint-1-i] = 1+i%maxConstraint, 1+i%maxConstraint
		input.groupSize[i] = 2
	}

	return input
}

func allSameInput() benchInput {
	input := newBenchInput("shape=all-same", true)
	input.groupSize = []int{maxConstraint}

	for i := range input.types {
		input.colors[i] = 1
	}

	return input
}

func tripletsInput() benchInput {
	input := newBenchInput("shape=triplets", true)

	for i := range input.types {
		group := min(i/3, maxConstraint/3-1)
		input.types[i] = group
		input.colors[i] = 1 + group
	}

	for i := range input.groupSize {
		input.groupSize[i] = 3
	}

	// 10^5 не делится на 3: последний элемент получает свой тип с группой из одного
	last := maxConstraint - 1
	input.types[last] = last
	input.colors[last] = 1
	input.groupSize[last] = 1

	return input
}

// randomStuckInput: соседние элементы различны, группы по 2 — удалить нечего, стек полон.
func randomStuckInput() benchInput {
	input := newBenchInput("shape=random-stuck", false)
	rng := rand.New(rand.NewPCG(2026, 2))

	for i := range input.types {
		input.types[i] = rng.IntN(maxConstraint)
		input.colors[i] = 1 + rng.IntN(maxConstraint)

		if i > 0 && input.types[i] == input.types[i-1] && input.colors[i] == input.colors[i-1] {
			input.colors[i] = 1 + input.colors[i]%maxConstraint
		}

		input.groupSize[i] = 2
	}

	return input
}

func newBenchInput(name string, want bool) benchInput {
	return benchInput{
		name:      name,
		types:     make([]int, maxConstraint),
		colors:    make([]int, maxConstraint),
		groupSize: make([]int, maxConstraint),
		want:      want,
	}
}

func canRemoveAllGrow(types, colors, groupSize []int) bool {
	var stack []run

	for i, kind := range types {
		key := kind<<colorBits | colors[i]

		if len(stack) > 0 && stack[len(stack)-1].key == key {
			stack[len(stack)-1].count++
		} else {
			stack = append(stack, run{key: key, count: 1})
		}

		if stack[len(stack)-1].count == groupSize[kind] {
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

type structRun struct{ kind, color, count int }

func canRemoveAllStruct(types, colors, groupSize []int) bool {
	stack := make([]structRun, 0, len(types))

	for i, kind := range types {
		color := colors[i]

		if top := len(stack) - 1; top >= 0 && stack[top].kind == kind && stack[top].color == color {
			stack[top].count++
		} else {
			stack = append(stack, structRun{kind: kind, color: color, count: 1})
		}

		if stack[len(stack)-1].count == groupSize[kind] {
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

// canRemoveAllElements хранит каждый элемент и длину блока до него включительно:
// группа снимается срезом сразу на groupSize элементов.
func canRemoveAllElements(types, colors, groupSize []int) bool {
	keys := make([]int, 0, len(types))
	counts := make([]int, 0, len(types))

	for i, kind := range types {
		key := kind<<colorBits | colors[i]
		count := 1

		if top := len(keys) - 1; top >= 0 && keys[top] == key {
			count = counts[top] + 1
		}

		if count == groupSize[kind] {
			keys = keys[:len(keys)-(count-1)]
			counts = counts[:len(counts)-(count-1)]

			continue
		}

		keys = append(keys, key)
		counts = append(counts, count)
	}

	return len(keys) == 0
}
