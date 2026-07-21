// Package main — lcgen: генератор каталога задачи по её слагу на leetcode.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var (
		daily      = flag.Bool("daily", false, "взять задачу дня")
		slug       = flag.String("slug", "", "titleSlug задачи")
		num        = flag.Int("num", 0, "номер задачи (фолбэк, если leetcode недоступен)")
		title      = flag.String("title", "", "название задачи (фолбэк)")
		difficulty = flag.String("difficulty", "", "сложность (фолбэк)")
		force      = flag.Bool("force", false, "перезаписать существующий каталог")
		root       = flag.String("root", "", "корень репозитория (по умолчанию — рядом с бинарём)")
	)
	flag.Parse()

	if err := run(*daily, *slug, *num, *title, *difficulty, *force, *root); err != nil {
		fmt.Fprintln(os.Stderr, "lcgen:", err)
		os.Exit(1)
	}
}

func run(daily bool, slug string, num int, title, difficulty string, force bool, root string) error {
	if !daily && slug == "" {
		return fmt.Errorf("нужен -daily или -slug")
	}

	date := ""
	if daily {
		var err error
		if date, slug, err = fetchDaily(); err != nil {
			return fmt.Errorf("задача дня недоступна: %w", err)
		}
	}

	q, err := fetchQuestion(slug)
	if err != nil {
		if daily {
			return err
		}
		fmt.Fprintf(os.Stderr, "lcgen: leetcode недоступен (%v), делаю заглушку\n", err)
		if q, err = offlineQuestion(slug, num, title, difficulty); err != nil {
			return err
		}
	}

	if root == "" {
		if root, err = repoRoot(); err != nil {
			return err
		}
	}

	dir, created, err := scaffold(q, root, force)
	if err != nil {
		return err
	}

	if date != "" {
		fmt.Printf("date:       %s\n", date)
	}
	fmt.Printf("problem:    %s. %s (%s)\n", q.FrontendID, q.Title, strings.ToLower(q.Difficulty))
	fmt.Printf("link:       https://leetcode.com/problems/%s/\n", q.TitleSlug)
	if len(q.TopicTags) > 0 {
		var tags []string
		for _, t := range q.TopicTags {
			tags = append(tags, t.Name)
		}
		fmt.Printf("topics:     %s\n", strings.Join(tags, ", "))
	}
	fmt.Printf("dir:        %s (%s)\n", dir, created)
	return nil
}

func scaffold(q Question, root string, force bool) (dir, status string, err error) {
	id, err := strconv.Atoi(q.FrontendID)
	if err != nil {
		return "", "", fmt.Errorf("странный номер задачи %q: %w", q.FrontendID, err)
	}

	pkg := packageName(q.TitleSlug)
	dir = filepath.Join(root, "problems", fmt.Sprintf("%04d-%s", id, q.TitleSlug))

	if _, err := os.Stat(dir); err == nil && !force {
		return dir, "already exists, left untouched", nil
	}

	statement := htmlToMarkdown(q.Content)
	files := map[string]string{
		"README.md":        renderReadme(q, statement),
		"solution.go":      renderSolution(q, pkg),
		"solution_test.go": renderTest(pkg, q.meta(), q.ExampleTestcases, expectedOutputs(statement)),
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", err
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			return "", "", err
		}
	}
	return dir, "created", nil
}

func offlineQuestion(slug string, num int, title, difficulty string) (Question, error) {
	if num == 0 || title == "" {
		return Question{}, fmt.Errorf("нет сети и не переданы -num/-title для заглушки")
	}
	if difficulty == "" {
		difficulty = "unknown"
	}
	return Question{
		FrontendID: strconv.Itoa(num),
		Title:      title,
		TitleSlug:  slug,
		Difficulty: difficulty,
	}, nil
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
