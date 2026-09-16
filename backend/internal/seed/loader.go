package seed

import (
	"embed"
	"encoding/json"
	"io/fs"
	"sort"
)

// Контент курса хранится как данные в content/*.json и встраивается в бинарник.
// Числовой префикс имени файла задаёт порядок глав.
//
//go:embed content/*.json
var contentFS embed.FS

// Второй курс (Go Backend Engineer) — в content-go/*.json.
//
//go:embed content-go/*.json
var goContentFS embed.FS

type lessonFile struct {
	Title       string         `json:"title"`
	Kind        string         `json:"kind"`
	Summary     string         `json:"summary"`
	DurationMin int            `json:"durationMin"`
	Content     map[string]any `json:"content"`
}

type moduleFile struct {
	Title   string       `json:"title"`
	Summary string       `json:"summary"`
	Lessons []lessonFile `json:"lessons"`
}

// loadModules читает главы курса DevOps из content/*.json.
func loadModules() []ModuleSeed { return readModules(contentFS, "content") }

// loadGoModules читает главы курса Go из content-go/*.json.
func loadGoModules() []ModuleSeed { return readModules(goContentFS, "content-go") }

// readModules читает главы курса из встроенных JSON-файлов каталога dir,
// отсортированных по имени.
func readModules(contentFS embed.FS, dir string) []ModuleSeed {
	entries, err := fs.ReadDir(contentFS, dir)
	if err != nil {
		panic("seed: не удалось прочитать " + dir + "/: " + err.Error())
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	modules := make([]ModuleSeed, 0, len(names))
	for _, name := range names {
		data, err := contentFS.ReadFile(dir + "/" + name)
		if err != nil {
			panic("seed: не удалось прочитать " + dir + "/" + name + ": " + err.Error())
		}
		var mf moduleFile
		if err := json.Unmarshal(data, &mf); err != nil {
			panic("seed: битый JSON " + dir + "/" + name + ": " + err.Error())
		}
		m := ModuleSeed{Title: mf.Title, Summary: mf.Summary}
		for _, lf := range mf.Lessons {
			m.Lessons = append(m.Lessons, LessonSeed{
				Title:       lf.Title,
				Kind:        lf.Kind,
				Summary:     lf.Summary,
				DurationMin: lf.DurationMin,
				Content:     lf.Content,
			})
		}
		modules = append(modules, m)
	}
	return modules
}
