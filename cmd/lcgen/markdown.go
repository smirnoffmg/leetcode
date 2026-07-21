package main

import (
	"html"
	"regexp"
	"strings"
)

var (
	inlineTags    = map[string]string{"strong": "**", "b": "**", "em": "_", "i": "_", "code": "`"}
	blockTags     = map[string]bool{"p": true, "ul": true, "ol": true, "div": true, "blockquote": true}
	trailingSpace = regexp.MustCompile(`[ \t]+\n`)
	blankLines    = regexp.MustCompile(`\n{3,}`)
	emptyMarks    = regexp.MustCompile(`\*\*[ \t]*\*\*|_[ \t]*_`)
	outputLine    = regexp.MustCompile(`(?m)^\*{0,2}Output:\*{0,2}\s*(.*?)\s*$`)
)

// htmlToMarkdown переводит условие задачи в markdown. Разметка leetcode — узкое
// подмножество HTML, поэтому обходимся сканером тегов вместо полноценного парсера.
func htmlToMarkdown(src string) string {
	var buf []byte
	write := func(s string) { buf = append(buf, s...) }
	pre, trimLead := 0, false

	// Текстовый узел: внутри <pre> сохраняется как есть, снаружи переносы строк из
	// вёрстки значения не имеют — но склеивать соседние теги без пробела нельзя.
	writeText := func(s string) {
		if pre > 0 {
			if trimLead {
				s, trimLead = strings.TrimLeft(s, " \t\n"), false
			}
			write(text(s))
			return
		}
		if s != "" && strings.TrimSpace(s) == "" {
			if len(buf) > 0 && buf[len(buf)-1] != '\n' {
				write(" ")
			}
			return
		}
		write(text(s))
	}

	for len(src) > 0 {
		open := strings.IndexByte(src, '<')
		if open < 0 {
			writeText(src)
			break
		}
		writeText(src[:open])

		end := strings.IndexByte(src[open:], '>')
		if end < 0 {
			writeText(src[open:])
			break
		}

		name, closing := parseTag(src[open+1 : open+end])
		src = src[open+end+1:]

		switch {
		case name == "pre":
			if closing {
				pre--
				buf = trimRightNewlines(buf)
			} else {
				pre++
				trimLead = true
			}
			write("\n```\n")
		case pre > 0:
			if name == "br" {
				write("\n")
			}
		case inlineTags[name] != "":
			if closing {
				// "**слово **" markdown уже не считает выделением, поэтому пробел
				// переезжает за закрывающий маркер.
				trimmed := strings.TrimRight(string(buf), " ")
				spaces := len(buf) - len(trimmed)
				buf = []byte(trimmed)
				write(inlineTags[name])
				write(strings.Repeat(" ", spaces))
			} else {
				write(inlineTags[name])
			}
		case name == "li" && !closing:
			write("\n- ")
		case name == "li":
			// Перевод строки поставит следующий <li> или закрытие списка.
		case blockTags[name]:
			write("\n\n")
		case name == "br":
			write("\n")
		case name == "sup" && !closing:
			write("^")
		case name == "sub" && !closing:
			write("_")
		}
	}

	out := trailingSpace.ReplaceAllString(string(buf), "\n")
	out = blankLines.ReplaceAllString(out, "\n\n")
	out = emptyMarks.ReplaceAllString(out, "")
	return strings.TrimSpace(out)
}

func trimRightNewlines(buf []byte) []byte {
	return []byte(strings.TrimRight(string(buf), " \t\n"))
}

func parseTag(raw string) (name string, closing bool) {
	raw = strings.TrimSpace(raw)
	closing = strings.HasPrefix(raw, "/")
	raw = strings.TrimPrefix(raw, "/")
	if i := strings.IndexAny(raw, " \t\n/"); i >= 0 {
		raw = raw[:i]
	}
	return strings.ToLower(raw), closing
}

func text(s string) string {
	return strings.ReplaceAll(html.UnescapeString(s), "\u00a0", " ")
}

// expectedOutputs собирает строки "Output: ..." из примеров — API отдаёт только входы.
func expectedOutputs(statement string) []string {
	var out []string
	for _, m := range outputLine.FindAllStringSubmatch(statement, -1) {
		out = append(out, strings.TrimSpace(m[1]))
	}
	return out
}
