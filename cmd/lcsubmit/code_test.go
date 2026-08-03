package main

import (
	"strings"
	"testing"
)

func TestPrepareCode(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		want     []string
		wantGone []string
		wantErr  bool
	}{
		{
			name: "срезает doc-комментарий и package",
			src: "// Package twosum — решение leetcode 1. Two Sum.\npackage twosum\n\n" +
				"func twoSum(nums []int) []int {\n\treturn nums\n}\n",
			want:     []string{"func twoSum(nums []int) []int {"},
			wantGone: []string{"package", "Package twosum"},
		},
		{
			name: "оставляет импорты и комментарии внутри кода",
			src: "package p\n\nimport \"sort\"\n\n" +
				"// подсказка про инвариант\nfunc f(a []int) { sort.Ints(a) }\n",
			want: []string{"import \"sort\"", "// подсказка про инвариант", "sort.Ints(a)"},
		},
		{
			name: "срезает объявления типов leetcode",
			src: "package p\n\n// ListNode из шаблона.\ntype ListNode struct {\n\tVal  int\n\tNext *ListNode\n}\n\n" +
				"type TreeNode struct{ Val int }\n\nfunc f(h *ListNode) int { return h.Val }\n",
			want:     []string{"func f(h *ListNode) int"},
			wantGone: []string{"type ListNode", "type TreeNode", "ListNode из шаблона"},
		},
		{
			name:     "оставляет собственные типы",
			src:      "package p\n\ntype pair struct{ a, b int }\n\nfunc f() pair { return pair{} }\n",
			want:     []string{"type pair struct{ a, b int }", "func f() pair"},
			wantGone: nil,
		},
		{
			name:    "не отправляет незаполненную заготовку",
			src:     "package p\n\nfunc f() int {\n\tpanic(\"not implemented\")\n}\n",
			wantErr: true,
		},
		{
			name:    "не отправляет то, что не парсится",
			src:     "package p\n\nfunc f( {\n",
			wantErr: true,
		},
		{
			name:    "не отправляет пустой файл",
			src:     "package p\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prepareCode([]byte(tt.src))
			if (err != nil) != tt.wantErr {
				t.Fatalf("prepareCode() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("в коде нет %q:\n%s", want, got)
				}
			}
			for _, gone := range tt.wantGone {
				if strings.Contains(got, gone) {
					t.Errorf("в коде осталось %q:\n%s", gone, got)
				}
			}
		})
	}
}
