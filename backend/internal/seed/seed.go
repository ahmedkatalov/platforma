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
		Title:    "DevOps-инженер: с нуля до уверенной практики",
		Subtitle: "Linux, сети, Git, Bash, Python, Go, Docker, CI/CD, облако, Terraform, Ansible, Kubernetes, наблюдаемость, безопасность и SRE",
		Description: "Полный практический курс для начинающих без предварительных знаний. Материал идёт по зависимости: " +
			"сначала рабочая среда, терминал, Linux и сети; затем Git, автоматизация на Bash, Python и практический минимум Go; " +
			"после — приложения, базы данных, Docker, CI/CD, облако, IaC, Kubernetes, мониторинг, безопасность и SRE. " +
			"Новая команда или конструкция сначала объясняется на простом примере, затем отрабатывается в терминале или редакторе. " +
			"В каждом большом блоке есть проверки, диагностика типовых ошибок и практический результат для портфолио. " +
			"Финал — сквозной капстоун от репозитория до наблюдаемого и восстанавливаемого окружения.",
		Level: "beginner",
		Tags: []string{
			"devops", "linux", "git", "docker", "ci-cd",
			"kubernetes", "terraform", "monitoring", "security",
			"python", "go", "bash", "ansible", "sre", "observability",
		},
		Modules: loadModules(),
	}
}
