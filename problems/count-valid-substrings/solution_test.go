package countvalidsubstrings

import (
	"math/rand/v2"
	"strings"
	"testing"
)

const maxLength = 100_000

func TestCountValidSubstrings(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int64
	}{
		{name: "example 1", s: "abac", want: 10},
		{name: "example 2", s: "abc", want: 5},
		{name: "example 3", s: "ac", want: 3},

		{name: "single char", s: "a", want: 1},
		{name: "two chars ab", s: "ab", want: 3},
		{name: "reversed letters are valid", s: "cba", want: 6},
		{name: "abc in the middle", s: "aabcc", want: 11},
		{name: "adjacent occurrences", s: "abcabc", want: 14},
		// вхождения разделены одним символом: левая граница сдвигается скачком
		{name: "separated occurrences", s: "abcbabc", want: 19},
		{name: "abc at the very end", s: "ccabc", want: 12},
		{name: "abc at the very start", s: "abcaa", want: 12},
		// 5 000 050 000 не влезает в int32
		{name: "max length without abc", s: strings.Repeat("a", maxLength), want: 5_000_050_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countValidSubstrings(tt.s); got != tt.want {
				t.Errorf("countValidSubstrings(%q) = %v, want %v", truncated(tt.s), got, tt.want)
			}
		})
	}
}

func TestCountValidSubstringsMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2026, 1107))

	for range 20_000 {
		s := randomString(rng, 1+rng.IntN(12))
		want := bruteForceCount(s)

		for _, impl := range implementations() {
			if got := impl.fn(s); got != want {
				t.Fatalf("%s(%q) = %v, want %v", impl.name, s, got, want)
			}
		}
	}
}

func TestImplementationsAgreeOnMaxLength(t *testing.T) {
	for _, input := range benchInputs() {
		t.Run(input.name, func(t *testing.T) {
			want := countValidSubstrings(input.s)

			for _, impl := range implementations() {
				if got := impl.fn(input.s); got != want {
					t.Fatalf("%s = %v, want %v", impl.name, got, want)
				}
			}
		})
	}
}

func BenchmarkCountValidSubstrings(b *testing.B) {
	for _, input := range benchInputs() {
		for _, impl := range implementations() {
			b.Run(input.name+"/"+impl.name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					impl.fn(input.s)
				}
			})
		}
	}
}

// bruteForceCount проверяет определение напрямую: перебирает все подстроки.
func bruteForceCount(s string) int64 {
	var count int64

	for left := range s {
		for right := left + 1; right <= len(s); right++ {
			if !strings.Contains(s[left:right], "abc") {
				count++
			}
		}
	}

	return count
}

type implementation struct {
	name string
	fn   func(string) int64
}

func implementations() []implementation {
	return []implementation{
		{name: "impl=scan-slice", fn: countValidSubstrings},
		{name: "impl=scan-bytes", fn: countScanBytes},
		{name: "impl=strings-index", fn: countStringsIndex},
	}
}

type benchInput struct {
	name string
	s    string
}

func benchInputs() []benchInput {
	rng := rand.New(rand.NewPCG(2026, 1))

	return []benchInput{
		{name: "shape=random", s: randomString(rng, maxLength)},
		{name: "shape=all-abc", s: strings.Repeat("abc", maxLength/3)},
		{name: "shape=no-abc-ab", s: strings.Repeat("ab", maxLength/2)},
		{name: "shape=all-a", s: strings.Repeat("a", maxLength)},
	}
}

func randomString(rng *rand.Rand, n int) string {
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = 'a' + byte(rng.IntN(3))
	}

	return string(letters)
}

func truncated(s string) string {
	if len(s) <= 32 {
		return s
	}

	return s[:32] + "..."
}

func countScanBytes(s string) int64 {
	var total int64

	last := -1

	for right := range s {
		if right >= 2 && s[right] == 'c' && s[right-1] == 'b' && s[right-2] == 'a' {
			last = right - 2
		}

		total += int64(right - last)
	}

	return total
}

// countStringsIndex вычитает из всех подстрок содержащие "abc": пока самое правое
// завершившееся вхождение начинается в start, у каждого правого конца ровно start+1
// плохих левых концов. Вхождения "abc" не перекрываются, поэтому поиск идёт с start+3.
func countStringsIndex(s string) int64 {
	n := int64(len(s))
	total := n * (n + 1) / 2

	start := strings.Index(s, "abc")
	for start >= 0 {
		searchFrom := start + 3

		next := strings.Index(s[searchFrom:], "abc")
		if next < 0 {
			total -= int64(start+1) * (n - int64(start) - 2)
			break
		}

		next += searchFrom
		total -= int64(start+1) * int64(next-start)
		start = next
	}

	return total
}
