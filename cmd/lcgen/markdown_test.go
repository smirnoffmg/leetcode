package main

import (
	"strings"
	"testing"
)

func TestHTMLToMarkdown(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "paragraph with inline markup",
			in:   `<p>Return <code>nums[i]</code> if <strong>valid</strong>.</p>`,
			want: "Return `nums[i]` if **valid**.",
		},
		{
			name: "example block stays verbatim",
			in:   "<pre><strong>Input:</strong> nums = [2,7]\n<strong>Output:</strong> [0,1]\n</pre>",
			want: "```\nInput: nums = [2,7]\nOutput: [0,1]\n```",
		},
		{
			name: "list becomes bullets",
			in:   `<ul><li>1 &lt;= n &lt;= 10<sup>4</sup></li><li>unique</li></ul>`,
			want: "- 1 <= n <= 10^4\n- unique",
		},
		{
			name: "entities and nbsp",
			in:   `<p>a&nbsp;&amp;&nbsp;b</p>`,
			want: "a & b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := htmlToMarkdown(tt.in); got != tt.want {
				t.Errorf("htmlToMarkdown(%q) =\n%q\nwant\n%q", tt.in, got, tt.want)
			}
		})
	}
}

func TestHTMLToMarkdownKeepsExamplesReadable(t *testing.T) {
	in := `<p>Example 1:</p><pre><strong>Input:</strong> s = "abc"
<strong>Output:</strong> 3
<strong>Explanation:</strong> whole string.
</pre><p><strong>Constraints:</strong></p><ul><li>0 &lt;= s.length</li></ul>`

	got := htmlToMarkdown(in)

	if strings.Contains(got, "<") && !strings.Contains(got, "0 <= s.length") {
		t.Errorf("теги должны быть вырезаны:\n%s", got)
	}
	if !strings.Contains(got, "**Constraints:**") {
		t.Errorf("заголовок ограничений потерян:\n%s", got)
	}
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("лишние пустые строки:\n%q", got)
	}
}

func TestExpectedOutputs(t *testing.T) {
	statement := "```\nInput: nums = [2,7]\nOutput: [0,1]\n```\n```\nInput: nums = [3,3]\nOutput: [0,1]\nExplanation: x\n```"

	got := expectedOutputs(statement)
	want := []string{"[0,1]", "[0,1]"}

	assertOutputs(t, got, want)
}

// Новые условия отдают примеры не в <pre>, а абзацами: "**Output:** 1".
func TestExpectedOutputsBoldFormat(t *testing.T) {
	statement := "**Input:** s = \"01\"\n\n**Output:** 1\n\n**Input:** s = \"0100\"\n\n**Output:** 4\n"

	assertOutputs(t, expectedOutputs(statement), []string{"1", "4"})
}

func assertOutputs(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expectedOutputs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("expectedOutputs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
