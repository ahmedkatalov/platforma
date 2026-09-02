// Команда exportcourse выгружает текущий курс в один JSON-пакет «курс как файл»,
// который потом загружают в админке кнопкой «Загрузить курс из файла».
//
//	go run ./cmd/exportcourse                 # -> devops-engineer.course.json
//	go run ./cmd/exportcourse -o /tmp/x.json
//	go run ./cmd/exportcourse -o x.json -ids prev.course.json   # перенести id из прошлого экспорта
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
	ID          string         `json:"id,omitempty"` // для обновления на месте (переносится из -ids)
	Title       string         `json:"title"`
	Kind        string         `json:"kind"`
	Summary     string         `json:"summary"`
	Content     map[string]any `json:"content"`
	DurationMin int            `json:"durationMin"`
}

type pkgModule struct {
	ID      string      `json:"id,omitempty"`
	Title   string      `json:"title"`
	Summary string      `json:"summary"`
	Lessons []pkgLesson `json:"lessons"`
}

type pkgCourseMeta struct {
	ID          string   `json:"id,omitempty"`
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
	ids := flag.String("ids", "", "прошлый экспорт (platforma-course): перенести из него id курса/глав/уроков по позиции")
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

	if *ids != "" {
		if err := graftIDs(&pkg, *ids); err != nil {
			log.Fatal(err)
		}
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

// graftIDs переносит id курса, глав и уроков из прошлого экспорта по позиции.
// Идентификаторы нужны, чтобы импорт в админке обновлял курс на месте (SyncContent
// сначала ищет по id) и не терял прогресс студентов. Порядок и названия должны совпадать —
// иначе отказываемся, чтобы не приклеить чужой id.
func graftIDs(pkg *coursePackage, refPath string) error {
	raw, err := os.ReadFile(refPath)
	if err != nil {
		return fmt.Errorf("ids: %w", err)
	}
	var ref coursePackage
	if err := json.Unmarshal(raw, &ref); err != nil {
		return fmt.Errorf("ids: разбор %s: %w", refPath, err)
	}
	if len(ref.Modules) != len(pkg.Modules) {
		return fmt.Errorf("ids: глав в эталоне %d, в выгрузке %d", len(ref.Modules), len(pkg.Modules))
	}
	pkg.Course.ID = ref.Course.ID
	moved := 0
	for i := range pkg.Modules {
		rm, m := ref.Modules[i], &pkg.Modules[i]
		if rm.Title != m.Title {
			return fmt.Errorf("ids: глава %d: «%s» ≠ «%s»", i+1, rm.Title, m.Title)
		}
		if len(rm.Lessons) != len(m.Lessons) {
			return fmt.Errorf("ids: глава %d: уроков в эталоне %d, в выгрузке %d", i+1, len(rm.Lessons), len(m.Lessons))
		}
		m.ID = rm.ID
		for j := range m.Lessons {
			if rm.Lessons[j].Title != m.Lessons[j].Title {
				return fmt.Errorf("ids: глава %d урок %d: «%s» ≠ «%s»", i+1, j+1, rm.Lessons[j].Title, m.Lessons[j].Title)
			}
			m.Lessons[j].ID = rm.Lessons[j].ID
			moved++
		}
	}
	fmt.Printf("  id перенесены из %s: курс, %d глав, %d уроков\n", refPath, len(pkg.Modules), moved)
	return nil
}
