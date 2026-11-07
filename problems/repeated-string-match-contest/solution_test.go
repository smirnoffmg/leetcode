package repeatedstringmatchcontest

import (
	"math/rand/v2"
	"strings"
	"testing"
)

const maxRepeats = 1_000_000_000

func TestRepeatedStringMatch(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		x, y int
		want bool
	}{
		{name: "example 1", a: "ab", b: "abab", x: 2, y: 1, want: true},
		{name: "example 2", a: "abc", b: "bc", x: 2, y: 3, want: false},
		{name: "example 3", a: "abcd", b: "cdab", x: 2, y: 1, want: true},

		{name: "identical single chars", a: "a", b: "a", x: 1, y: 1, want: true},
		{name: "pattern longer by one repeat", a: "a", b: "a", x: 1, y: 2, want: false},
		{name: "different letters", a: "a", b: "b", x: 5, y: 1, want: false},
		// сдвиг 1 есть в бесконечном повторе, но в "ab" места для него уже нет
		{name: "rotation fits only with room", a: "ab", b: "ba", x: 1, y: 1, want: false},
		{name: "rotation with room", a: "ab", b: "ba", x: 2, y: 1, want: true},
		{name: "offset limited by total length", a: "ab", b: "ba", x: maxRepeats, y: maxRepeats, want: false},
		{name: "one repeat of room at limits", a: "ab", b: "ba", x: maxRepeats, y: maxRepeats - 1, want: true},
		{name: "max equal lengths", a: "a", b: "aa", x: maxRepeats, y: maxRepeats / 2, want: true},
		{name: "one char too long at limits", a: "a", b: "aa", x: maxRepeats, y: maxRepeats/2 + 1, want: false},
		// различие видно только на втором повторе b
		{name: "mismatch after first repeat of b", a: "aab", b: "aa", x: 3, y: 2, want: false},
		{name: "b is longer than a", a: "ab", b: "bababa", x: 10, y: 2, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := repeatedStringMatch(tt.a, tt.b, tt.x, tt.y); got != tt.want {
				t.Errorf("repeatedStringMatch(%q, %q, %d, %d) = %v, want %v", tt.a, tt.b, tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestRepeatedStringMatchMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(2026, 1110))

	trueCount, falseCount := 0, 0

	for range 50_000 {
		a := randomString(rng, 1+rng.IntN(6))
		b := randomPattern(rng, a)
		x, y := 1+rng.IntN(6), 1+rng.IntN(6)
		want := strings.Contains(strings.Repeat(a, x), strings.Repeat(b, y))

		if want {
			trueCount++
		} else {
			falseCount++
		}

		for _, impl := range implementations() {
			if got := impl.fn(a, b, x, y); got != want {
				t.Fatalf("%s(%q, %q, %d, %d) = %v, want %v", impl.name, a, b, x, y, got, want)
			}
		}
	}

	if trueCount < 5000 || falseCount < 5000 {
		t.Fatalf("генератор перекошен: true = %d, false = %d", trueCount, falseCount)
	}
}

// Граница Файна — Вильфа сверяется с lcm: там периоды заведомо проходят полный цикл.
func TestFineWilfBoundAgreesWithLCMOnLongStrings(t *testing.T) {
	rng := rand.New(rand.NewPCG(2026, 1111))

	for range 5_000 {
		a := randomString(rng, 1+rng.IntN(100))
		b := randomPattern(rng, a)
		x, y := 1+rng.IntN(maxRepeats), 1+rng.IntN(maxRepeats)

		for _, impl := range implementations() {
			want := repeatedStringMatchLCM(a, b, x, y)
			if got := impl.fn(a, b, x, y); got != want {
				t.Fatalf("%s(%q, %q, %d, %d) = %v, want %v", impl.name, a, b, x, y, got, want)
			}
		}
	}
}

func BenchmarkRepeatedStringMatch(b *testing.B) {
	inputs := []struct {
		name         string
		text, needle string
		x, y         int
	}{
		{name: "shape=all-same-true", text: strings.Repeat("a", 100), needle: strings.Repeat("a", 99), x: maxRepeats, y: maxRepeats},
		{name: "shape=late-mismatch", text: strings.Repeat("a", 100), needle: strings.Repeat("a", 99) + "b", x: maxRepeats, y: maxRepeats / 2},
		{name: "shape=coprime-false", text: strings.Repeat("a", 99) + "b", needle: strings.Repeat("a", 99), x: maxRepeats, y: maxRepeats / 2},
	}

	for _, input := range inputs {
		for _, impl := range implementations() {
			b.Run(input.name+"/"+impl.name, func(b *testing.B) {
				b.ReportAllocs()

				for b.Loop() {
					impl.fn(input.text, input.needle, input.x, input.y)
				}
			})
		}
	}
}

type implementation struct {
	name string
	fn   func(a, b string, x, y int) bool
}

func implementations() []implementation {
	return []implementation{
		{name: "impl=contains", fn: repeatedStringMatch},
		{name: "impl=fine-wilf-modulo", fn: repeatedStringMatchModulo},
		{name: "impl=lcm-modulo", fn: repeatedStringMatchLCM},
	}
}

func randomString(rng *rand.Rand, n int) string {
	letters := make([]byte, n)
	for i := range letters {
		letters[i] = 'a' + byte(rng.IntN(2))
	}

	return string(letters)
}

// randomPattern в половине случаев вырезает кусок из повтора a, иначе случайные строки
// почти всегда давали бы false.
func randomPattern(rng *rand.Rand, a string) string {
	n := 1 + rng.IntN(max(len(a), 4))
	if rng.IntN(2) == 0 {
		return randomString(rng, n)
	}

	offset := rng.IntN(len(a))

	return strings.Repeat(a, n/len(a)+2)[offset : offset+n]
}

func repeatedStringMatchLCM(a, b string, x, y int) bool {
	return matchesByModulo(a, b, x, y, len(a)/gcd(len(a), len(b))*len(b))
}

func repeatedStringMatchModulo(a, b string, x, y int) bool {
	return matchesByModulo(a, b, x, y, len(a)+len(b)-gcd(len(a), len(b)))
}

func matchesByModulo(a, b string, x, y, bound int) bool {
	total, pattern := len(a)*x, len(b)*y
	if pattern > total {
		return false
	}

	compared := min(pattern, bound)

	for offset := range min(len(a)-1, total-pattern) + 1 {
		matched := true

		for i := range compared {
			if a[(offset+i)%len(a)] != b[i%len(b)] {
				matched = false

				break
			}
		}

		if matched {
			return true
		}
	}

	return false
}
