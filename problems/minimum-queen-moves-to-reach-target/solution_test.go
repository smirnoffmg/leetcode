package minimumqueenmovestoreachtarget

import "testing"

func TestMinQueenMoves(t *testing.T) {
	tests := []struct {
		name           string
		source, target []int
		want           int
	}{
		{name: "example 1", source: []int{8, 1}, target: []int{1, 8}, want: 1},
		{name: "example 2", source: []int{4, 2}, target: []int{1, 3}, want: 2},

		{name: "same square", source: []int{5, 5}, target: []int{5, 5}, want: 0},
		{name: "same row", source: []int{3, 1}, target: []int{3, 8}, want: 1},
		{name: "same column", source: []int{1, 6}, target: []int{7, 6}, want: 1},
		{name: "main diagonal", source: []int{1, 1}, target: []int{8, 8}, want: 1},
		{name: "anti-diagonal", source: []int{2, 7}, target: []int{7, 2}, want: 1},
		// соседние по «ходу коня» клетки не лежат ни на одной линии
		{name: "knight jump", source: []int{1, 1}, target: []int{2, 3}, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minQueenMoves(tt.source, tt.target); got != tt.want {
				t.Errorf("minQueenMoves(%v, %v) = %d, want %d", tt.source, tt.target, got, tt.want)
			}
		})
	}
}

// Доска маленькая, поэтому сверка идёт со всеми 64×64 парами, а не со случайной выборкой.
func TestMinQueenMovesMatchesBFS(t *testing.T) {
	for sr := 1; sr <= 8; sr++ {
		for sc := 1; sc <= 8; sc++ {
			dist := queenBFS(sr, sc)
			for tr := 1; tr <= 8; tr++ {
				for tc := 1; tc <= 8; tc++ {
					got := minQueenMoves([]int{sr, sc}, []int{tr, tc})
					if want := dist[tr][tc]; got != want {
						t.Fatalf("minQueenMoves([%d,%d], [%d,%d]) = %d, want %d", sr, sc, tr, tc, got, want)
					}
				}
			}
		}
	}
}

// queenBFS — эталон от определения хода: расстояния от (sr, sc) до всех клеток доски.
func queenBFS(sr, sc int) [9][9]int {
	var dist [9][9]int
	for i := range dist {
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	dist[sr][sc] = 0
	queue := [][2]int{{sr, sc}}

	for i := 0; i < len(queue); i++ {
		p := queue[i]
		for _, d := range dirs {
			for r, c := p[0]+d[0], p[1]+d[1]; r >= 1 && r <= 8 && c >= 1 && c <= 8; r, c = r+d[0], c+d[1] {
				if dist[r][c] < 0 {
					dist[r][c] = dist[p[0]][p[1]] + 1
					queue = append(queue, [2]int{r, c})
				}
			}
		}
	}

	return dist
}

func BenchmarkMinQueenMoves(b *testing.B) {
	var pairs [][2][]int
	for sr := 1; sr <= 8; sr++ {
		for sc := 1; sc <= 8; sc++ {
			for tr := 1; tr <= 8; tr++ {
				for tc := 1; tc <= 8; tc++ {
					pairs = append(pairs, [2][]int{{sr, sc}, {tr, tc}})
				}
			}
		}
	}

	b.ResetTimer()
	for b.Loop() {
		for _, p := range pairs {
			minQueenMoves(p[0], p[1])
		}
	}
}
