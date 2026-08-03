package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolveProblem(t *testing.T) {
	problems := t.TempDir()
	touch := func(name string, age time.Duration) {
		dir := filepath.Join(problems, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(dir, "solution.go")
		if err := os.WriteFile(file, []byte("package p\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		when := time.Now().Add(-age)
		if err := os.Chtimes(file, when, when); err != nil {
			t.Fatal(err)
		}
	}
	touch("0001-two-sum", time.Hour)
	touch("0877-stone-game", time.Minute)

	t.Run("по слагу", func(t *testing.T) {
		dir, slug, err := resolveProblem(problems, "two-sum")
		if err != nil {
			t.Fatal(err)
		}
		if want := filepath.Join(problems, "0001-two-sum"); dir != want {
			t.Errorf("dir = %q, want %q", dir, want)
		}
		if slug != "two-sum" {
			t.Errorf("slug = %q, want two-sum", slug)
		}
	})

	t.Run("без слага — самая свежая правка", func(t *testing.T) {
		_, slug, err := resolveProblem(problems, "")
		if err != nil {
			t.Fatal(err)
		}
		if slug != "stone-game" {
			t.Errorf("slug = %q, want stone-game", slug)
		}
	})

	t.Run("неизвестный слаг", func(t *testing.T) {
		if _, _, err := resolveProblem(problems, "nope"); err == nil {
			t.Error("ожидалась ошибка")
		}
	})
}
