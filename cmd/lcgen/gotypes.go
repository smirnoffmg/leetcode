package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

var baseTypes = map[string]string{
	"integer":       "int",
	"int":           "int",
	"long":          "int64",
	"double":        "float64",
	"float":         "float64",
	"string":        "string",
	"boolean":       "bool",
	"character":     "byte",
	"void":          "",
	"ListNode":      "*ListNode",
	"TreeNode":      "*TreeNode",
	"Node":          "*Node",
	"NestedInteger": "*NestedInteger",
}

func goType(t string) string {
	t = strings.TrimSpace(t)
	if inner, ok := strings.CutPrefix(t, "list<"); ok {
		if inner, ok := strings.CutSuffix(inner, ">"); ok {
			return "[]" + goType(inner)
		}
	}
	if inner, ok := strings.CutSuffix(t, "[]"); ok {
		return "[]" + goType(inner)
	}
	if got, ok := baseTypes[t]; ok {
		return got
	}
	return "any"
}

func zeroValue(gotype string) string {
	switch {
	case gotype == "string":
		return `""`
	case gotype == "bool":
		return "false"
	case strings.HasPrefix(gotype, "[]"), strings.HasPrefix(gotype, "*"), gotype == "any":
		return "nil"
	default:
		return "0"
	}
}

// goLiteral переводит значение примера (JSON, как его отдаёт leetcode) в литерал Go
// нужного типа. Указательные типы (деревья, списки) литералом не выразить — вызывающий
// получит ошибку и оставит TODO с исходной строкой.
func goLiteral(raw, gotype string) (string, error) {
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.UseNumber()

	var value any
	if err := dec.Decode(&value); err != nil {
		return "", err
	}
	return literal(value, gotype, false)
}

func literal(value any, gotype string, elideType bool) (string, error) {
	if strings.HasPrefix(gotype, "*") || gotype == "any" {
		return "", fmt.Errorf("тип %s литералом не выражается", gotype)
	}

	if inner, ok := strings.CutPrefix(gotype, "[]"); ok {
		if value == nil {
			return "nil", nil
		}
		items, ok := value.([]any)
		if !ok {
			return "", fmt.Errorf("ожидался список для %s", gotype)
		}

		parts := make([]string, 0, len(items))
		for _, item := range items {
			part, err := literal(item, inner, strings.HasPrefix(inner, "[]"))
			if err != nil {
				return "", err
			}
			parts = append(parts, part)
		}

		body := "{" + strings.Join(parts, ", ") + "}"
		if elideType {
			return body, nil
		}
		return gotype + body, nil
	}

	switch gotype {
	case "string":
		s, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("ожидалась строка, получено %T", value)
		}
		return strconv.Quote(s), nil
	case "byte":
		s, ok := value.(string)
		if !ok || len(s) != 1 {
			return "", fmt.Errorf("ожидался символ, получено %v", value)
		}
		return strconv.QuoteRune(rune(s[0])), nil
	case "bool":
		b, ok := value.(bool)
		if !ok {
			return "", fmt.Errorf("ожидался bool, получено %T", value)
		}
		return fmt.Sprintf("%t", b), nil
	default:
		n, ok := value.(json.Number)
		if !ok {
			return "", fmt.Errorf("ожидалось число, получено %T", value)
		}
		return n.String(), nil
	}
}
