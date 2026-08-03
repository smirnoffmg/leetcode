package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Эти типы leetcode объявляет сам в шаблоне задачи; lcgen достаёт их из комментария,
// чтобы пакет компилировался локально, но при отправке повторное объявление ломает сборку.
var templateTypes = map[string]bool{"ListNode": true, "TreeNode": true, "Node": true}

const unimplemented = `panic("not implemented")`

// prepareCode превращает solution.go в сниппет для leetcode: без package clause
// и без объявлений типов из шаблона, но с импортами и комментариями по месту.
func prepareCode(src []byte) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "solution.go", src, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("solution.go не парсится: %w", err)
	}

	type span struct{ from, to int }
	cuts := []span{{0, fset.Position(file.Name.End()).Offset}}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || !isTemplateType(gen) {
			continue
		}
		from := gen.Pos()
		if gen.Doc != nil {
			from = gen.Doc.Pos()
		}
		cuts = append(cuts, span{fset.Position(from).Offset, fset.Position(gen.End()).Offset})
	}

	var kept strings.Builder
	prev := 0
	for _, cut := range cuts {
		kept.Write(src[prev:cut.from])
		prev = cut.to
	}
	kept.Write(src[prev:])

	code := strings.TrimSpace(kept.String())
	if code == "" {
		return "", fmt.Errorf("в solution.go нечего отправлять")
	}
	if strings.Contains(code, unimplemented) {
		return "", fmt.Errorf("решение ещё не написано: в коде остался %s", unimplemented)
	}
	return code + "\n", nil
}

func isTemplateType(gen *ast.GenDecl) bool {
	if gen.Tok != token.TYPE || len(gen.Specs) == 0 {
		return false
	}
	for _, spec := range gen.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok || !templateTypes[ts.Name.Name] {
			return false
		}
	}
	return true
}
