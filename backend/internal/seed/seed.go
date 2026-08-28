// Package seed содержит демонстрационный курс по DevOps, которым наполняется
// платформа командой `go run ./cmd/seedcourse`.
package seed

// LessonSeed — урок в заготовке курса.
type LessonSeed struct {
	Title       string
	Kind        string
	Summary     string
	DurationMin int
	Content     map[string]any
}

// ModuleSeed — модуль с уроками.
type ModuleSeed struct {
	Title   string
	Summary string
	Lessons []LessonSeed
}

// CourseSeed — курс целиком.
type CourseSeed struct {
	Slug        string
	Title       string
	Subtitle    string
	Description string
	Level       string
	Tags        []string
	Modules     []ModuleSeed
}

// DevOpsCourse собирает курс из модулей, описанных в соседних файлах.
func DevOpsCourse() CourseSeed {
	return CourseSeed{
		Slug:     "devops-engineer",
		Title:    "DevOps с нуля до практики",
		Subtitle: "Linux, Docker, CI/CD и Kubernetes — с квизами и тренажёром терминала",
		Description: "Практический курс для тех, кто хочет собирать, доставлять и эксплуатировать " +
			"приложения как инженер DevOps. Каждый модуль заканчивается квизом, а ключевые темы " +
			"закрепляются в тренажёре терминала и редакторе конфигураций.",
		Level: "beginner",
		Tags:  []string{"devops", "linux", "docker", "ci-cd", "kubernetes"},
		Modules: []ModuleSeed{
			moduleIntro(),
			moduleLinux(),
			moduleDocker(),
			moduleCICD(),
			moduleKubernetes(),
		},
	}
}
