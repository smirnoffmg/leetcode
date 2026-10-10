package smallestpalindromicrearrangementii

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestSmallestPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		k    int
		want string
	}{
		{name: "example 1", s: "abba", k: 2, want: "baab"},
		{name: "example 2", s: "aa", k: 2, want: ""},
		{name: "example 3", s: "bacab", k: 1, want: "abcba"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := smallestPalindrome(tt.s, tt.k)
			if got != tt.want {
				t.Errorf("smallestPalindrome(%v, %v) = %v, want %v", tt.s, tt.k, got, tt.want)
			}
		})
	}
}

// Эталон строится от определения: все различные перестановки s, из них палиндромы,
// по возрастанию. Узкий алфавит даёт повторы букв, на которых и ошибаются в подсчёте.
func TestSmallestPalindromeMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(3518, 2026))

	for range 1500 {
		s := randomPalindrome(rng, rng.IntN(5), 2+rng.IntN(2))
		all := palindromicPermutations(s)
		k := 1 + rng.IntN(len(all)+2)

		want := ""
		if k <= len(all) {
			want = all[k-1]
		}

		if got := smallestPalindrome(s, k); got != want {
			t.Fatalf("smallestPalindrome(%q, %d) = %q, want %q", s, k, got, want)
		}
	}
}

// Число перестановок половины здесь больше 10⁶, поэтому countArrangements упирается
// в потолок permLimit — эта ветка на коротких входах не исполняется.
func TestSmallestPalindromeSaturatedCount(t *testing.T) {
	rng := rand.New(rand.NewPCG(3518, 7))

	for range 20 {
		s := randomPalindrome(rng, 10+rng.IntN(3), 6+rng.IntN(4))
		k := 1 + rng.IntN(1_000_000)
		want := kthHalfByEnumeration(s, k)

		if got := smallestPalindrome(s, k); got != want {
			t.Fatalf("smallestPalindrome(%q, %d) = %q, want %q", s, k, got, want)
		}
	}
}

func BenchmarkSmallestPalindrome(b *testing.B) {
	rng := rand.New(rand.NewPCG(3518, 1))
	s := randomPalindrome(rng, 5_000, 26)

	b.ResetTimer()

	for b.Loop() {
		smallestPalindrome(s, 1_000_000)
	}
}

// randomPalindrome — палиндром с половиной длины half над первыми alphabet буквами
// и случайной серединой либо без неё.
func randomPalindrome(rng *rand.Rand, half, alphabet int) string {
	left := make([]byte, half)
	for i := range left {
		left[i] = byte('a' + rng.IntN(alphabet))
	}

	right := slices.Clone(left)
	slices.Reverse(right)

	mid := ""
	if rng.IntN(2) == 0 {
		mid = string(byte('a' + rng.IntN(alphabet)))
	}

	return string(left) + mid + string(right)
}

func palindromicPermutations(s string) []string {
	seen := map[string]bool{}
	letters := []byte(s)

	var permute func(i int)
	permute = func(i int) {
		if i == len(letters) {
			p := string(letters)
			if isPalindrome(p) {
				seen[p] = true
			}

			return
		}

		for j := i; j < len(letters); j++ {
			letters[i], letters[j] = letters[j], letters[i]
			permute(i + 1)
			letters[i], letters[j] = letters[j], letters[i]
		}
	}
	permute(0)

	all := make([]string, 0, len(seen))
	for p := range seen {
		all = append(all, p)
	}
	slices.Sort(all)

	return all
}

func isPalindrome(s string) bool {
	for i := range len(s) / 2 {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}

	return true
}

// kthHalfByEnumeration перечисляет перестановки левой половины в лексикографическом
// порядке и останавливается на k-й — без подсчёта числа перестановок вовсе.
func kthHalfByEnumeration(s string, k int) string {
	var count [26]int
	for i := range len(s) / 2 {
		count[s[i]-'a']++
	}

	half := make([]byte, 0, len(s)/2)
	found := ""

	var walk func() bool
	walk = func() bool {
		if len(half) == cap(half) {
			k--
			if k == 0 {
				right := slices.Clone(half)
				slices.Reverse(right)
				found = string(half) + s[len(s)/2:len(s)-len(s)/2] + string(right)

				return true
			}

			return false
		}

		for c := range count {
			if count[c] == 0 {
				continue
			}

			count[c]--
			half = append(half, byte('a'+c))
			done := walk()
			half = half[:len(half)-1]
			count[c]++

			if done {
				return true
			}
		}

		return false
	}
	walk()

	return found
}
