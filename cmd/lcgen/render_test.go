package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGoType(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "scalar", in: "integer", want: "int"},
		{name: "slice", in: "integer[]", want: "[]int"},
		{name: "matrix", in: "character[][]", want: "[][]byte"},
		{name: "list", in: "list<string>", want: "[]string"},
		{name: "node", in: "TreeNode", want: "*TreeNode"},
		{name: "unknown", in: "whatever", want: "any"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := goType(tt.in); got != tt.want {
				t.Errorf("goType(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestGoLiteral(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		gotype  string
		want    string
		wantErr bool
	}{
		{name: "int", raw: "9", gotype: "int", want: "9"},
		{name: "ints", raw: "[2,7,11,15]", gotype: "[]int", want: "[]int{2, 7, 11, 15}"},
		{name: "nested elides inner type", raw: "[[1,2],[3]]", gotype: "[][]int", want: "[][]int{{1, 2}, {3}}"},
		{name: "string", raw: `"ab\"c"`, gotype: "string", want: `"ab\"c"`},
		{name: "chars", raw: `[["a","b"]]`, gotype: "[][]byte", want: "[][]byte{{'a', 'b'}}"},
		{name: "bool", raw: "true", gotype: "bool", want: "true"},
		{name: "null slice", raw: "null", gotype: "[]int", want: "nil"},
		{name: "tree not expressible", raw: "[1,null,2]", gotype: "*TreeNode", wantErr: true},
		{name: "type mismatch", raw: `"x"`, gotype: "int", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := goLiteral(tt.raw, tt.gotype)
			if (err != nil) != tt.wantErr {
				t.Fatalf("goLiteral(%q, %q) err = %v, wantErr %v", tt.raw, tt.gotype, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("goLiteral(%q, %q) = %q, want %q", tt.raw, tt.gotype, got, tt.want)
			}
		})
	}
}

func TestPackageName(t *testing.T) {
	tests := []struct {
		name string
		slug string
		want string
	}{
		{name: "words", slug: "two-sum", want: "twosum"},
		{name: "leading digit", slug: "3sum-closest", want: "p3sumclosest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := packageName(tt.slug); got != tt.want {
				t.Errorf("packageName(%q) = %q, want %q", tt.slug, got, tt.want)
			}
		})
	}
}

func TestRenderTestUsesExamplesAndOutputs(t *testing.T) {
	meta := unmarshalMeta(t, `{"name":"twoSum",
		"params":[{"name":"nums","type":"integer[]"},{"name":"target","type":"integer"}],
		"return":{"type":"integer[]"}}`)

	got := renderTest("twosum", meta, "[2,7,11,15]\n9\n[3,2,4]\n6", []string{"[0,1]", "[1,2]"})

	for _, want := range []string{
		"package twosum",
		`"reflect"`,
		"func TestTwoSum(t *testing.T) {",
		"nums   []int",
		"target int",
		"want   []int",
		`{name: "example 1", nums: []int{2, 7, 11, 15}, target: 9, want: []int{0, 1}},`,
		`{name: "example 2", nums: []int{3, 2, 4}, target: 6, want: []int{1, 2}},`,
		"got := twoSum(tt.nums, tt.target)",
		"if !reflect.DeepEqual(got, tt.want) {",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("сгенерированный тест не содержит %q:\n%s", want, got)
		}
	}
}

func TestRenderTestScalarReturnSkipsReflect(t *testing.T) {
	meta := unmarshalMeta(t, `{"name":"isPalindrome",
		"params":[{"name":"x","type":"integer"}],"return":{"type":"boolean"}}`)

	got := renderTest("palindromenumber", meta, "121\n-121", []string{"true", "false"})

	if strings.Contains(got, "reflect") {
		t.Errorf("для скалярного результата reflect не нужен:\n%s", got)
	}
	if !strings.Contains(got, "if got != tt.want {") {
		t.Errorf("ожидалось сравнение через !=:\n%s", got)
	}
	if !strings.Contains(got, `{name: "example 2", x: -121, want: false},`) {
		t.Errorf("ожидался второй пример:\n%s", got)
	}
}

func TestRenderTestMarksUnexpressibleValues(t *testing.T) {
	meta := unmarshalMeta(t, `{"name":"maxDepth",
		"params":[{"name":"root","type":"TreeNode"}],"return":{"type":"integer"}}`)

	got := renderTest("maximumdepth", meta, "[3,9,20,null,null,15,7]", []string{"3"})

	if !strings.Contains(got, "// TODO: root = [3,9,20,null,null,15,7]") {
		t.Errorf("дерево должно остаться TODO с исходной строкой:\n%s", got)
	}
	if !strings.Contains(got, "root: nil") {
		t.Errorf("ожидалось нулевое значение для root:\n%s", got)
	}
}

func TestRenderTestFallsBackToStub(t *testing.T) {
	got := renderTest("designproblem", Meta{}, "[\"MyStack\",\"push\"]\n[[],[1]]", nil)

	if !strings.Contains(got, `t.Skip("TODO: write the failing test first")`) {
		t.Errorf("ожидалась заглушка:\n%s", got)
	}
	if !strings.Contains(got, `// ["MyStack","push"]`) {
		t.Errorf("примеры должны остаться в комментарии:\n%s", got)
	}
}

func TestRenderSolutionUnfoldsNodeDefinition(t *testing.T) {
	q := questionWithSnippet("104", "Maximum Depth of Binary Tree", `/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func maxDepth(root *TreeNode) int {

}`)

	got := renderSolution(q, "maximumdepth")

	for _, want := range []string{
		"// Package maximumdepth — решение leetcode 104. Maximum Depth of Binary Tree.",
		"package maximumdepth",
		"type TreeNode struct {",
		"Left  *TreeNode",
		"func maxDepth(root *TreeNode) int {",
		`panic("not implemented")`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("solution.go не содержит %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Definition for a binary tree node") {
		t.Errorf("комментарий-определение должен был превратиться в код:\n%s", got)
	}
}

func TestRenderSolutionDiscardsStubParams(t *testing.T) {
	tests := []struct {
		name    string
		snippet string
		want    string
	}{
		{
			name:    "named params",
			snippet: "func twoSum(nums []int, target int) []int {\n\n}",
			want:    "\t_, _ = nums, target\n\tpanic(\"not implemented\")",
		},
		{
			name:    "grouped params share a type",
			snippet: "func gcd(a, b int) int {\n\n}",
			want:    "\t_, _ = a, b\n\tpanic(\"not implemented\")",
		},
		{
			name:    "method on a design-problem receiver",
			snippet: "func (this *MyStack) Push(x int) {\n\n}",
			want:    "\t_ = x\n\tpanic(\"not implemented\")",
		},
		{
			name:    "no params",
			snippet: "func Constructor() MyStack {\n\n}",
			want:    "MyStack {\n\tpanic(\"not implemented\")",
		},
		{
			name:    "blank param needs no discard",
			snippet: "func solve(_ int) int {\n\n}",
			want:    "int {\n\tpanic(\"not implemented\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderSolution(questionWithSnippet("1", "Some Problem", tt.snippet), "pkg")
			if !strings.Contains(got, tt.want) {
				t.Errorf("solution.go не содержит %q:\n%s", tt.want, got)
			}
		})
	}
}

func TestRenderReadmeCollectsMetadata(t *testing.T) {
	q := Question{
		FrontendID:       "1",
		Title:            "Two Sum",
		TitleSlug:        "two-sum",
		Difficulty:       "Easy",
		Hints:            []string{"Use a <code>map</code>."},
		Stats:            `{"acRate":"55.7%"}`,
		SimilarQuestions: `[{"title":"3Sum","titleSlug":"3sum","difficulty":"Medium"}]`,
		TopicTags: []struct {
			Name string `json:"name"`
		}{{Name: "Array"}, {Name: "Hash Table"}},
	}

	got := renderReadme(q, "Given an array...")

	for _, want := range []string{
		"# 1. Two Sum",
		"- Difficulty: easy",
		"- Link: https://leetcode.com/problems/two-sum/",
		"- Topics: Array, Hash Table",
		"- Acceptance: 55.7%",
		"## Statement\n\nGiven an array...",
		"<details><summary>Hint 1</summary>",
		"Use a `map`.",
		"- [3Sum](https://leetcode.com/problems/3sum/) (medium)",
		"- Time: O(?)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("README не содержит %q:\n%s", want, got)
		}
	}
}

func questionWithSnippet(id, title, snippet string) Question {
	q := Question{FrontendID: id, Title: title}
	q.CodeSnippets = append(q.CodeSnippets, struct {
		LangSlug string `json:"langSlug"`
		Code     string `json:"code"`
	}{LangSlug: "golang", Code: snippet})
	return q
}

func unmarshalMeta(t *testing.T, raw string) Meta {
	t.Helper()

	var m Meta
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		t.Fatalf("не разобрал metaData: %v", err)
	}
	return m
}
