// Package main — lcsubmit: отправка solution.go на leetcode и ожидание вердикта.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var (
		slug    = flag.String("slug", "", "titleSlug задачи (по умолчанию — свежая правка в problems/)")
		lang    = flag.String("lang", "golang", "язык решения на leetcode")
		timeout = flag.Duration("timeout", time.Minute, "сколько ждать вердикт")
		root    = flag.String("root", "", "корень репозитория (по умолчанию — вверх по дереву до go.mod)")
		dry     = flag.Bool("dry-run", false, "показать код для отправки и выйти")
	)
	flag.Parse()

	if err := run(*slug, *lang, *root, *timeout, *dry); err != nil {
		fmt.Fprintln(os.Stderr, "lcsubmit:", err)
		os.Exit(1)
	}
}

var errRejected = errors.New("решение не принято")

func run(slug, lang, root string, timeout time.Duration, dry bool) error {
	if root == "" {
		var err error
		if root, err = repoRoot(); err != nil {
			return err
		}
	}

	dir, slug, err := resolveProblem(filepath.Join(root, "problems"), slug)
	if err != nil {
		return err
	}

	src, err := os.ReadFile(filepath.Join(dir, "solution.go"))
	if err != nil {
		return err
	}
	code, err := prepareCode(src)
	if err != nil {
		return err
	}
	if dry {
		fmt.Print(code)
		return nil
	}

	session, csrf := os.Getenv("LEETCODE_SESSION"), os.Getenv("LEETCODE_CSRF")
	if session == "" || csrf == "" {
		return errors.New("нужны переменные LEETCODE_SESSION и LEETCODE_CSRF — это cookies с leetcode.com")
	}

	c := newClient(session, csrf)
	q, err := c.question(slug)
	if err != nil {
		return err
	}

	id, err := c.submit(slug, q.ID, lang, code)
	if err != nil {
		return err
	}

	fmt.Printf("problem:    %s. %s\n", q.FrontendID, q.Title)
	fmt.Printf("submission: %s/submissions/detail/%d/\n", origin, id)

	result, err := c.waitResult(slug, id, timeout)
	if err != nil {
		return err
	}

	printResult(result)
	if !result.accepted() {
		return errRejected
	}
	return nil
}

func printResult(r checkResult) {
	verdict := r.StatusMsg
	if r.TotalTestcases > 0 {
		verdict = fmt.Sprintf("%s (%d/%d тестов)", verdict, r.TotalCorrect, r.TotalTestcases)
	}
	fmt.Printf("verdict:    %s\n", verdict)

	if r.accepted() {
		fmt.Printf("runtime:    %s (быстрее %.2f%% решений)\n", r.StatusRuntime, r.RuntimePercentile)
		fmt.Printf("memory:     %s (экономнее %.2f%% решений)\n", r.StatusMemory, r.MemoryPercentile)
		return
	}

	if r.CompileError != "" {
		fmt.Printf("compile:    %s\n", indent(r.CompileError))
	}
	if r.RuntimeError != "" {
		fmt.Printf("panic:      %s\n", indent(r.RuntimeError))
	}
	if r.LastTestcase != "" {
		fmt.Printf("input:      %s\n", indent(r.LastTestcase))
		fmt.Printf("expected:   %s\n", r.ExpectedOutput)
		fmt.Printf("got:        %s\n", r.CodeOutput)
	}
}

func indent(text string) string {
	return strings.ReplaceAll(strings.TrimSpace(text), "\n", "\n            ")
}

// resolveProblem находит каталог задачи. Без слага берётся тот, где solution.go правился последним.
func resolveProblem(problems, slug string) (dir, resolved string, err error) {
	entries, err := os.ReadDir(problems)
	if err != nil {
		return "", "", err
	}

	var newest time.Time
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		parts := strings.SplitN(name, "-", 2)
		if len(parts) != 2 {
			continue
		}

		if slug != "" {
			if parts[1] == slug {
				return filepath.Join(problems, name), slug, nil
			}
			continue
		}

		info, err := os.Stat(filepath.Join(problems, name, "solution.go"))
		if err != nil {
			continue
		}
		if info.ModTime().After(newest) {
			newest, dir, resolved = info.ModTime(), filepath.Join(problems, name), parts[1]
		}
	}

	if dir == "" {
		if slug != "" {
			return "", "", fmt.Errorf("нет каталога задачи %q — сначала `make new S=%s`", slug, slug)
		}
		return "", "", fmt.Errorf("в %s нет ни одной задачи", problems)
	}
	return dir, resolved, nil
}

func repoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("не нашёл go.mod вверх по дереву от %s", dir)
		}
		dir = parent
	}
}
