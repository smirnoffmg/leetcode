package twosum

import (
	"reflect"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name   string
		nums   []int
		target int
		want   []int
	}{
		{name: "example 1", nums: []int{2, 7, 11, 15}, target: 9, want: []int{0, 1}},
		{name: "example 2", nums: []int{3, 2, 4}, target: 6, want: []int{1, 2}},
		{name: "example 3", nums: []int{3, 3}, target: 6, want: []int{0, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := twoSum(tt.nums, tt.target)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("twoSum(%v, %v) = %v, want %v", tt.nums, tt.target, got, tt.want)
			}
		})
	}
}

// Худший случай: пара — два последних элемента, поэтому в таблицу успевает
// попасть весь массив. Значения подряд дают единственную пару с суммой 2n−3.
func BenchmarkTwoSum(b *testing.B) {
	const n = 10_000

	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}

	b.ResetTimer()

	for b.Loop() {
		twoSum(nums, 2*n-3)
	}
}
