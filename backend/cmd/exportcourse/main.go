// Команда exportcourse выгружает текущий курс в один JSON-пакет «курс как файл»,
// который потом загружают в админке кнопкой «Загрузить курс из файла».
//
//	go run ./cmd/exportcourse                 # -> devops-engineer.course.json
//	go run ./cmd/exportcourse -o /tmp/x.json
//	go run ./cmd/exportcourse -o x.json -ids prev.course.json   # перенести id из прошлого экспорта
package main

import (
	"crypto/rand"
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

// graftIDs переносит id курса, глав и уроков из прошлого экспорта, сопоставляя главы
// и уроки ПО НАЗВАНИЮ (состав уроков между этапами меняется). Не найденным в эталоне
// главам/урокам выдаётся новый UUID v4 — так импорт в админке обновляет старое на месте
// (SyncContent ищет по id), а новое создаёт, не теряя прогресс студентов.
func graftIDs(pkg *coursePackage, refPath string) error {
	raw, err := os.ReadFile(refPath)
	if err != nil {
		return fmt.Errorf("ids: %w", err)
	}
	var ref coursePackage
	if err := json.Unmarshal(raw, &ref); err != nil {
		return fmt.Errorf("ids: разбор %s: %w", refPath, err)
	}
	pkg.Course.ID = ref.Course.ID
	if pkg.Course.ID == "" {
		pkg.Course.ID = newUUID()
	}
	refMods := map[string]pkgModule{}
	for _, m := range ref.Modules {
		refMods[m.Title] = m
	}
	kept, created, orphan := 0, 0, 0
	for i := range pkg.Modules {
		m := &pkg.Modules[i]
		rm, ok := refMods[m.Title]
		if ok {
			m.ID = rm.ID
		} else {
			m.ID = newUUID()
			fmt.Printf("  новая глава без эталона: «%s» → новый id\n", m.Title)
		}
		refLes := map[string]string{}
		for _, l := range rm.Lessons {
			refLes[l.Title] = l.ID
		}
		seen := map[string]bool{}
		for j := range m.Lessons {
			l := &m.Lessons[j]
			if id, ok := refLes[l.Title]; ok && !seen[l.Title] {
				l.ID, seen[l.Title] = id, true
				kept++
			} else {
				l.ID = newUUID()
				created++
			}
		}
		for t := range refLes {
			if !seen[t] {
				orphan++
				fmt.Printf("  ВНИМАНИЕ: урок из эталона не найден (переименован/удалён?): глава «%s» → «%s»\n", m.Title, t)
			}
		}
	}
	fmt.Printf("  id: сохранено %d, новых UUID %d, потеряно из эталона %d\n", kept, created, orphan)
	return nil
}

// newUUID — UUID v4 из crypto/rand (без внешних зависимостей).
func newUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
