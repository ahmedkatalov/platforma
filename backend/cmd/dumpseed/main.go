// Команда dumpseed выгружает текущий курс из Go-модулей в JSON-файлы content/,
// чтобы дальше наполнять курс контентом как данными, а не кодом.
//
//	go run ./cmd/dumpseed
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"platforma/backend/internal/seed"
)

type lessonOut struct {
	Title       string         `json:"title"`
	Kind        string         `json:"kind"`
	Summary     string         `json:"summary"`
	DurationMin int            `json:"durationMin"`
	Content     map[string]any `json:"content"`
}

type moduleOut struct {
	Title   string      `json:"title"`
	Summary string      `json:"summary"`
	Lessons []lessonOut `json:"lessons"`
}

func main() {
	course := seed.DevOpsCourse()
	dir := filepath.Join("internal", "seed", "content")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatal(err)
	}

	for i, m := range course.Modules {
		out := moduleOut{Title: m.Title, Summary: m.Summary}
		for _, l := range m.Lessons {
			out.Lessons = append(out.Lessons, lessonOut{
				Title: l.Title, Kind: l.Kind, Summary: l.Summary,
				DurationMin: l.DurationMin, Content: l.Content,
			})
		}
		data, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		// Шаг нумерации 10 — чтобы позже вставлять главы между существующими.
		name := fmt.Sprintf("%03d.json", (i+1)*10)
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s  ←  %s (%d уроков)\n", name, m.Title, len(m.Lessons))
	}
	fmt.Printf("\nВыгружено модулей: %d\n", len(course.Modules))
}
