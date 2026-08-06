package smallestdivisibledigitproducti

import (
	"testing"
)

func TestSmallestNumber(t *testing.T) {
	tests := []struct {
		name string
		n    int
		t    int
		want int
	}{
		{name: "example 1", n: 10, t: 2, want: 10},
		{name: "example 2", n: 15, t: 3, want: 16},
		{name: "минимальные значения", n: 1, t: 1, want: 1},
		{name: "n уже подходит", n: 5, t: 5, want: 5},
		{name: "худший случай: 10 итераций до нуля в конце", n: 11, t: 10, want: 20},
		{name: "переход через сотню", n: 99, t: 7, want: 100},
		{name: "верхняя граница", n: 100, t: 10, want: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smallestNumber(tt.n, tt.t)
			if got != tt.want {
				t.Errorf("smallestNumber(%v, %v) = %v, want %v", tt.n, tt.t, got, tt.want)
			}
		})
	}
}

func BenchmarkSmallestNumber(b *testing.B) {
	// худший случай: продукт цифр 11..19 не делится на 10, ответ 20
	for b.Loop() {
		smallestNumber(11, 10)
	}
}
