package findmissingelements

import (
	"math/rand/v2"
	"reflect"
	"testing"
)

func TestFindMissingElements(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want []int
	}{
		{name: "example 1", nums: []int{1, 4, 2, 5}, want: []int{3}},
		{name: "example 2", nums: []int{7, 8, 6, 9}, want: []int{}},
		{name: "example 3", nums: []int{5, 1}, want: []int{2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMissingElements(tt.nums)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findMissingElements(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

// Крайние значения 1 и 100 присутствуют всегда, поэтому просматривается весь диапазон;
// остальные — различные значения из 2..99, как требует условие.
func BenchmarkFindMissingElements(b *testing.B) {
	rng := rand.New(rand.NewPCG(3731, 1))

	nums := []int{1, 100}
	for _, v := range rng.Perm(98)[:48] {
		nums = append(nums, v+2)
	}

	b.ResetTimer()

	for b.Loop() {
		findMissingElements(nums)
	}
}
