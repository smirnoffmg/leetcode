package smallestdivisibledigitproductii

import (
	"strconv"
	"strings"
	"testing"
)

func TestSmallestNumber(t *testing.T) {
	tests := []struct {
		name string
		num  string
		t    int64
		want string
	}{
		{name: "example 1", num: "1234", t: 256, want: "1488"},
		{name: "example 2", num: "12355", t: 50, want: "12355"},
		{name: "example 3", num: "11111", t: 26, want: "-1"},
		{name: "простой делитель 11 — нет ответа", num: "99", t: 11, want: "-1"},
		{name: "t=1, num уже zero-free", num: "42", t: 1, want: "42"},
		{name: "t=1, ноль в num", num: "20", t: 1, want: "21"},
		{name: "ноль в num, нужен суффикс", num: "105", t: 3, want: "113"},
		{name: "ни одна цифра не покрывает 10", num: "10", t: 10, want: "25"},
		{name: "все девятки — ответ длиннее", num: "999", t: 8, want: "1118"},
		{name: "t требует больше цифр, чем в num", num: "11", t: 244140625, want: "555555555555"},
		{name: "жадный хвост: 26 меньше 34", num: "11", t: 12, want: "26"},
		{name: "максимальный t", num: "11", t: 100_000_000_000_000, want: "4555555555555558888"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smallestNumber(tt.num, tt.t)
			if got != tt.want {
				t.Errorf("smallestNumber(%v, %v) = %v, want %v", tt.num, tt.t, got, tt.want)
			}
		})
	}
}

func TestSmallestNumberLongInput(t *testing.T) {
	num := strings.Repeat("9", 200_000)
	want := strings.Repeat("1", 200_001-16) + "2" + strings.Repeat("8", 15)

	got := smallestNumber(num, 1<<46)
	if got != want {
		t.Errorf("smallestNumber(9*200000, 2^46): len = %d, want len %d", len(got), len(want))
	}
}

// bruteforce — перебор вверх от start; для гладких t ответ близко, лимита хватает с запасом.
func bruteforce(start int, t int64) string {
	for x := start; x <= start+100_000; x++ {
		p, zero := int64(1), false
		for y := x; y > 0; y /= 10 {
			d := int64(y % 10)
			if d == 0 {
				zero = true

				break
			}
			p *= d
		}

		if !zero && p%t == 0 {
			return strconv.Itoa(x)
		}
	}

	return ""
}

func TestSmallestNumberBruteforce(t *testing.T) {
	divisors := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 16, 24, 35, 49, 50, 63, 81, 98, 100, 343}

	for num := 10; num <= 999; num++ {
		for _, tv := range divisors {
			want := bruteforce(num, tv)
			if got := smallestNumber(strconv.Itoa(num), tv); got != want {
				t.Fatalf("smallestNumber(%d, %d) = %v, want %v", num, tv, got, want)
			}
		}
	}
}

func BenchmarkSmallestNumber(b *testing.B) {
	// худший случай: все девятки — цикл проходит всю строку, ответ длиннее num
	num := strings.Repeat("9", 200_000)
	for b.Loop() {
		smallestNumber(num, 1<<46)
	}
}

func BenchmarkSmallestNumberNoChange(b *testing.B) {
	// num уже подходит: измеряет только префиксные суммы степеней
	num := strings.Repeat("8", 200_000)
	for b.Loop() {
		smallestNumber(num, 1<<46)
	}
}
