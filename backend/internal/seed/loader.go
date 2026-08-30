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

// loadModules читает главы курса из встроенных JSON-файлов, отсортированных по имени.
func loadModules() []ModuleSeed {
	entries, err := fs.ReadDir(contentFS, "content")
	if err != nil {
		panic("seed: не удалось прочитать content/: " + err.Error())
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
		data, err := contentFS.ReadFile("content/" + name)
		if err != nil {
			panic("seed: не удалось прочитать content/" + name + ": " + err.Error())
		}
		var mf moduleFile
		if err := json.Unmarshal(data, &mf); err != nil {
			panic("seed: битый JSON content/" + name + ": " + err.Error())
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
