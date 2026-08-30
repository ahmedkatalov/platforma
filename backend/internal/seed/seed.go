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
		Subtitle: "С нуля до практики: Linux, Git, сети, Docker, CI/CD, IaC, Kubernetes, облако, мониторинг, безопасность и надёжность",
		Description: "Курс для начинающих: предварительных знаний не требуется. Начинаем с первых команд " +
			"в терминале и доходим до Kubernetes, мониторинга, безопасности, облака и надёжности. Каждая тема объясняется " +
			"простыми словами, с примерами и разбором частых ошибок новичка. " +
			"В каждом модуле: теория со ссылками на первоисточники, тренажёр терминала прямо в браузере, " +
			"практика с настоящими конфигурациями и квиз с проверкой на сервере. " +
			"В конце — аттестация и сертификат с публичной страницей проверки.",
		Level: "beginner",
		Tags: []string{
			"devops", "linux", "git", "docker", "ci-cd",
			"kubernetes", "terraform", "monitoring", "security",
		},
		Modules: loadModules(),
	}
}
