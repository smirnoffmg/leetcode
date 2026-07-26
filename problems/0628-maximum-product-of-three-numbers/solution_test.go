package maximumproductofthreenumbers

import (
	"testing"
)

func TestMaximumProduct(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{name: "example 1", nums: []int{1, 2, 3}, want: 6},
		{name: "example 2", nums: []int{1, 2, 3, 4}, want: 24},
		{name: "example 3", nums: []int{-1, -2, -3}, want: -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maximumProduct(tt.nums)
			if got != tt.want {
				t.Errorf("maximumProduct(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
