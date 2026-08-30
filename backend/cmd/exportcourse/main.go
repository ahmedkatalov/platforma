// Команда exportcourse выгружает текущий курс в один JSON-пакет «курс как файл»,
// который потом загружают в админке кнопкой «Загрузить курс из файла».
//
//	go run ./cmd/exportcourse                 # -> devops-engineer.course.json
//	go run ./cmd/exportcourse -o /tmp/x.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"platforma/backend/internal/seed"
)

type pkgLesson struct {
	Title       string         `json:"title"`
	Kind        string         `json:"kind"`
	Summary     string         `json:"summary"`
	Content     map[string]any `json:"content"`
	DurationMin int            `json:"durationMin"`
}

type pkgModule struct {
	Title   string      `json:"title"`
	Summary string      `json:"summary"`
	Lessons []pkgLesson `json:"lessons"`
}

type pkgCourseMeta struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	CoverURL    string   `json:"coverUrl"`
	Level       string   `json:"level"`
	Tags        []string `json:"tags"`
}

type coursePackage struct {
	Format  string        `json:"format"`
	Version int           `json:"version"`
	Course  pkgCourseMeta `json:"course"`
	Modules []pkgModule   `json:"modules"`
}

func main() {
	out := flag.String("o", "", "путь выходного файла (по умолчанию <slug>.course.json)")
	flag.Parse()

	c := seed.DevOpsCourse()
	pkg := coursePackage{
		Format:  "platforma-course",
		Version: 1,
		Course: pkgCourseMeta{
			Slug: c.Slug, Title: c.Title, Subtitle: c.Subtitle,
			Description: c.Description, Level: c.Level, Tags: c.Tags,
		},
	}
	lessons := 0
	for _, m := range c.Modules {
		mod := pkgModule{Title: m.Title, Summary: m.Summary}
		for _, l := range m.Lessons {
			mod.Lessons = append(mod.Lessons, pkgLesson{
				Title: l.Title, Kind: l.Kind, Summary: l.Summary,
				Content: l.Content, DurationMin: l.DurationMin,
			})
			lessons++
		}
		pkg.Modules = append(pkg.Modules, mod)
	}

	data, err := json.MarshalIndent(pkg, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	path := *out
	if path == "" {
		path = c.Slug + ".course.json"
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Пакет курса: %s\n", path)
	fmt.Printf("  %d глав · %d уроков · %.1f МБ\n", len(pkg.Modules), lessons, float64(len(data))/1024/1024)
}
