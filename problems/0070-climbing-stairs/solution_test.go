package climbingstairs

import (
	"testing"
)

func TestClimbStairs(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{name: "example 1", n: 2, want: 2},
		{name: "example 2", n: 3, want: 3},
		// единственный способ — один шаг на одну ступеньку
		{name: "single step", n: 1, want: 1},
		{name: "four steps", n: 4, want: 5},
		{name: "five steps", n: 5, want: 8},
		{name: "ten steps", n: 10, want: 89},
		// n = 45  — верхняя граница условия, последнее значение, влезающее в int32
		{name: "upper bound", n: 45, want: 1_836_311_903},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := climbStairs(tt.n)
			if got != tt.want {
				t.Errorf("climbStairs(%v) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

// Весь диапазон условия сверяется с комбинаторным эталоном — он выведен
// из определения задачи, а не из рекурренты, поэтому общей ошибки у них быть не может.
func TestClimbStairsMatchesCombinatorics(t *testing.T) {
	for n := 1; n <= 45; n++ {
		if got, want := climbStairs(n), countByCombinatorics(n); got != want {
			t.Errorf("climbStairs(%d) = %d, want %d", n, got, want)
		}
	}
}

// На малых n тот же ответ подтверждается прямым перечислением последовательностей шагов.
func TestClimbStairsMatchesEnumeration(t *testing.T) {
	for n := 1; n <= 20; n++ {
		if got, want := climbStairs(n), enumerate(n); got != want {
			t.Errorf("climbStairs(%d) = %d, want %d", n, got, want)
		}
	}
}

func BenchmarkClimbStairs(b *testing.B) {
	for b.Loop() {
		climbStairs(45)
	}
}

// countByCombinatorics считает способы по числу двоек: путь с k двойками — это
// перестановка k двоек и n−2k единиц, то есть выбор позиций двоек среди n−k шагов.
func countByCombinatorics(n int) int {
	ways := 0
	for k := 0; 2*k <= n; k++ {
		ways += binomial(n-k, k)
	}

	return ways
}

func binomial(n, k int) int {
	if k < 0 || k > n {
		return 0
	}

	result := 1
	// произведение i+1 подряд идущих чисел делится на i+1, поэтому деление всегда точное
	for i := range k {
		result = result * (n - i) / (i + 1)
	}

	return result
}

// enumerate перебирает все последовательности шагов по 1 и 2 в сумме дающие n.
func enumerate(remaining int) int {
	if remaining == 0 {
		return 1
	}

	ways := 0

	for _, step := range []int{1, 2} {
		if step <= remaining {
			ways += enumerate(remaining - step)
		}
	}

	return ways
}
