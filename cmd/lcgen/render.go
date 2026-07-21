package main

import (
	"fmt"
	"go/format"
	"strings"
	"unicode"
)

func packageName(slug string) string {
	pkg := strings.ReplaceAll(slug, "-", "")
	if pkg == "" {
		return "solution"
	}
	// Пакет не может начинаться с цифры, а слаги вроде "3sum" встречаются.
	if unicode.IsDigit(rune(pkg[0])) {
		return "p" + pkg
	}
	return pkg
}

func renderReadme(q Question, statement string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s. %s\n\n", q.FrontendID, q.Title)
	fmt.Fprintf(&b, "- Difficulty: %s\n", strings.ToLower(q.Difficulty))
	fmt.Fprintf(&b, "- Link: https://leetcode.com/problems/%s/\n", q.TitleSlug)

	var tags []string
	for _, t := range q.TopicTags {
		tags = append(tags, t.Name)
	}
	if len(tags) > 0 {
		fmt.Fprintf(&b, "- Topics: %s\n", strings.Join(tags, ", "))
	}
	if rate := q.acceptance(); rate != "" {
		fmt.Fprintf(&b, "- Acceptance: %s\n", rate)
	}

	if statement != "" {
		fmt.Fprintf(&b, "\n## Statement\n\n%s\n", statement)
	}

	if len(q.Hints) > 0 {
		b.WriteString("\n## Hints\n")
		for i, hint := range q.Hints {
			fmt.Fprintf(&b, "\n<details><summary>Hint %d</summary>\n\n%s\n\n</details>\n",
				i+1, htmlToMarkdown(hint))
		}
	}

	if similar := q.similar(); len(similar) > 0 {
		b.WriteString("\n## Similar\n\n")
		for _, s := range similar {
			fmt.Fprintf(&b, "- [%s](https://leetcode.com/problems/%s/) (%s)\n",
				s.Title, s.TitleSlug, strings.ToLower(s.Difficulty))
		}
	}

	b.WriteString("\n## Idea\n\nTODO\n\n- Time: O(?)\n- Space: O(?)\n")
	return b.String()
}

func renderSolution(q Question, pkg string) string {
	defs, code := splitSnippet(q.goSnippet())

	var b strings.Builder
	fmt.Fprintf(&b, "// Package %s — решение leetcode %s. %s.\n", pkg, q.FrontendID, q.Title)
	fmt.Fprintf(&b, "package %s\n\n", pkg)
	for _, def := range defs {
		b.WriteString(def + "\n\n")
	}
	b.WriteString(fillBodies(code) + "\n")
	return gofmt(b.String())
}

// splitSnippet вытаскивает определения ListNode/TreeNode, которые leetcode прячет
// в комментарии: по CLAUDE.md они объявляются локально в пакете задачи.
func splitSnippet(snippet string) (defs []string, code string) {
	var kept []string
	lines := strings.Split(snippet, "\n")

	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(strings.TrimSpace(lines[i]), "/*") {
			kept = append(kept, lines[i])
			continue
		}

		var comment []string
		for ; i < len(lines); i++ {
			trimmed := strings.TrimSpace(lines[i])
			trimmed = strings.TrimPrefix(trimmed, "/**")
			trimmed = strings.TrimPrefix(trimmed, "/*")
			trimmed = strings.TrimPrefix(trimmed, "*/")
			trimmed = strings.TrimPrefix(trimmed, "*")
			comment = append(comment, strings.TrimPrefix(trimmed, " "))
			if strings.HasSuffix(strings.TrimSpace(lines[i]), "*/") {
				break
			}
		}
		if def := typeDefinition(comment); def != "" {
			defs = append(defs, def)
		}
	}

	return defs, strings.TrimSpace(strings.Join(kept, "\n"))
}

func typeDefinition(comment []string) string {
	start := -1
	for i, line := range comment {
		if strings.HasPrefix(line, "type ") {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}

	for i := start; i < len(comment); i++ {
		if strings.TrimSpace(comment[i]) == "}" {
			return strings.Join(comment[start:i+1], "\n")
		}
	}
	return ""
}

// fillBodies подставляет panic в пустые тела функций из сниппета — иначе пакет не собирается.
func fillBodies(code string) string {
	lines := strings.Split(code, "\n")
	var out []string

	for i := 0; i < len(lines); i++ {
		out = append(out, lines[i])
		if !strings.HasPrefix(lines[i], "func ") || !strings.HasSuffix(strings.TrimSpace(lines[i]), "{") {
			continue
		}

		j := i + 1
		for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
			j++
		}
		if j < len(lines) && strings.TrimSpace(lines[j]) == "}" {
			if names := paramNames(lines[i]); len(names) > 0 {
				blanks := strings.TrimSuffix(strings.Repeat("_, ", len(names)), ", ")
				out = append(out, fmt.Sprintf("\t%s = %s", blanks, strings.Join(names, ", ")))
			}
			out = append(out, "\tpanic(\"not implemented\")", "}")
			i = j
		}
	}

	return strings.Join(out, "\n")
}

// paramNames достаёт имена параметров из строки сигнатуры, чтобы заглушка могла
// их «использовать»: пока задача не решена, они иначе висят неиспользованными.
func paramNames(signature string) []string {
	params := paramList(signature)
	if params == "" {
		return nil
	}

	groups := splitTopLevel(params)
	named := false
	for _, g := range groups {
		if len(strings.Fields(g)) > 1 {
			named = true
			break
		}
	}
	// Смешивать именованные и анонимные параметры Go не даёт, так что либо
	// имена есть у всех, либо сигнатура из одних типов — и брать нечего.
	if !named {
		return nil
	}

	var names []string
	for _, g := range groups {
		if name := strings.Fields(g)[0]; name != "_" {
			names = append(names, name)
		}
	}
	return names
}

func paramList(signature string) string {
	rest := strings.TrimPrefix(signature, "func ")
	if strings.HasPrefix(rest, "(") {
		rest = rest[matchingParen(rest)+1:] // пропускаем получателя метода
	}

	open := strings.Index(rest, "(")
	if open < 0 {
		return ""
	}
	rest = rest[open:]
	return strings.TrimSpace(rest[1:matchingParen(rest)])
}

func matchingParen(s string) int {
	depth := 0
	for i, r := range s {
		switch r {
		case '(':
			depth++
		case ')':
			if depth--; depth == 0 {
				return i
			}
		}
	}
	return len(s) - 1
}

func splitTopLevel(params string) []string {
	var groups []string
	depth, start := 0, 0

	for i, r := range params {
		switch r {
		case '(', '[', '{':
			depth++
		case ')', ']', '}':
			depth--
		case ',':
			if depth == 0 {
				groups = append(groups, strings.TrimSpace(params[start:i]))
				start = i + 1
			}
		}
	}

	return append(groups, strings.TrimSpace(params[start:]))
}

func renderTest(pkg string, meta Meta, examples string, outputs []string) string {
	if meta.Name == "" || len(meta.Params) == 0 {
		return gofmt(renderTestStub(pkg, examples))
	}

	caseField := "name"
	for _, p := range meta.Params {
		if p.Name == "name" {
			caseField = "caseName"
		}
	}

	returnType := ""
	if meta.Return != nil {
		returnType = goType(meta.Return.Type)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\nimport (\n", pkg)
	if needsDeepEqual(returnType) {
		b.WriteString("\t\"reflect\"\n")
	}
	b.WriteString("\t\"testing\"\n)\n\n")

	fmt.Fprintf(&b, "func Test%s(t *testing.T) {\n\ttests := []struct {\n", exported(meta.Name))
	fmt.Fprintf(&b, "\t\t%s string\n", caseField)
	for _, p := range meta.Params {
		fmt.Fprintf(&b, "\t\t%s %s\n", p.Name, goType(p.Type))
	}
	if returnType != "" {
		fmt.Fprintf(&b, "\t\twant %s\n", returnType)
	}
	b.WriteString("\t}{\n")

	for _, c := range testCases(meta, examples, outputs, caseField, returnType) {
		b.WriteString("\t\t" + c + "\n")
	}
	b.WriteString("\t}\n\n")

	args := make([]string, 0, len(meta.Params))
	for _, p := range meta.Params {
		args = append(args, "tt."+p.Name)
	}
	argList := strings.Join(args, ", ")
	format := strings.TrimSuffix(strings.Repeat("%v, ", len(args)), ", ")

	fmt.Fprintf(&b, "\tfor _, tt := range tests {\n\t\tt.Run(tt.%s, func(t *testing.T) {\n", caseField)
	switch {
	case returnType == "":
		fmt.Fprintf(&b, "\t\t\t%s(%s)\n", meta.Name, argList)
	case needsDeepEqual(returnType):
		fmt.Fprintf(&b, "\t\t\tgot := %s(%s)\n", meta.Name, argList)
		fmt.Fprintf(&b, "\t\t\tif !reflect.DeepEqual(got, tt.want) {\n")
		fmt.Fprintf(&b, "\t\t\t\tt.Errorf(\"%s(%s) = %%v, want %%v\", %s, got, tt.want)\n", meta.Name, format, argList)
		b.WriteString("\t\t\t}\n")
	default:
		fmt.Fprintf(&b, "\t\t\tgot := %s(%s)\n", meta.Name, argList)
		b.WriteString("\t\t\tif got != tt.want {\n")
		fmt.Fprintf(&b, "\t\t\t\tt.Errorf(\"%s(%s) = %%v, want %%v\", %s, got, tt.want)\n", meta.Name, format, argList)
		b.WriteString("\t\t\t}\n")
	}
	b.WriteString("\t\t})\n\t}\n}\n")

	return gofmt(b.String())
}

func testCases(meta Meta, examples string, outputs []string, caseField, returnType string) []string {
	groups := groupExamples(examples, len(meta.Params))
	if len(groups) == 0 {
		return []string{fmt.Sprintf("{%s: \"example 1\"}, // TODO: заполнить примерами из условия", caseField)}
	}

	cases := make([]string, 0, len(groups))
	for i, group := range groups {
		fields := []string{fmt.Sprintf("%s: \"example %d\"", caseField, i+1)}
		var todo []string

		for j, p := range meta.Params {
			lit, err := goLiteral(group[j], goType(p.Type))
			if err != nil {
				lit = zeroValue(goType(p.Type))
				todo = append(todo, fmt.Sprintf("%s = %s", p.Name, group[j]))
			}
			fields = append(fields, fmt.Sprintf("%s: %s", p.Name, lit))
		}

		if returnType != "" {
			lit := zeroValue(returnType)
			switch {
			case i >= len(outputs):
				todo = append(todo, "want")
			default:
				if parsed, err := goLiteral(outputs[i], returnType); err == nil {
					lit = parsed
				} else {
					todo = append(todo, fmt.Sprintf("want = %s", outputs[i]))
				}
			}
			fields = append(fields, "want: "+lit)
		}

		line := "{" + strings.Join(fields, ", ") + "},"
		if len(todo) > 0 {
			line += " // TODO: " + strings.Join(todo, "; ")
		}
		cases = append(cases, line)
	}

	return cases
}

func groupExamples(examples string, size int) [][]string {
	if size <= 0 || strings.TrimSpace(examples) == "" {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(examples), "\n")
	var groups [][]string
	for i := 0; i+size <= len(lines); i += size {
		groups = append(groups, lines[i:i+size])
	}
	return groups
}

func renderTestStub(pkg, examples string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\nimport \"testing\"\n\n", pkg)

	if strings.TrimSpace(examples) != "" {
		b.WriteString("// Примеры из leetcode:\n")
		for _, line := range strings.Split(strings.TrimSpace(examples), "\n") {
			fmt.Fprintf(&b, "// %s\n", line)
		}
		b.WriteString("\n")
	}

	b.WriteString(`func TestSolve(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "example 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("TODO: write the failing test first")
		})
	}
}
`)
	return b.String()
}

func needsDeepEqual(gotype string) bool {
	return strings.HasPrefix(gotype, "[]") || strings.HasPrefix(gotype, "*") || gotype == "any"
}

func exported(name string) string {
	if name == "" {
		return "Solve"
	}
	return strings.ToUpper(name[:1]) + name[1:]
}

// gofmt форматирует сгенерированный код; если он не парсится (сниппет с сюрпризом),
// отдаём как есть — пусть человек увидит и поправит.
func gofmt(src string) string {
	out, err := format.Source([]byte(src))
	if err != nil {
		return src
	}
	return string(out)
}
