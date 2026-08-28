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
		Subtitle: "Linux, Git, сети, Docker, CI/CD, Kubernetes, IaC и мониторинг — с практикой в тренажёре",
		Description: "Полный практический курс для тех, кто хочет собирать, доставлять и эксплуатировать " +
			"приложения как инженер DevOps. Десять модулей от командной строки до Kubernetes и мониторинга: " +
			"в каждом — теория, тренажёр терминала, практика с конфигурациями и квиз с проверкой на сервере. " +
			"Курс завершается аттестацией: экзаменом, диагностикой инцидента и итоговой практикой, " +
			"после которой платформа выдаёт сертификат.",
		Level: "beginner",
		Tags: []string{
			"devops", "linux", "git", "docker", "ci-cd",
			"kubernetes", "terraform", "monitoring", "security",
		},
		Modules: []ModuleSeed{
			moduleIntro(),
			moduleLinux(),
			moduleGit(),
			moduleNetwork(),
			moduleDocker(),
			moduleCICD(),
			moduleIaC(),
			moduleKubernetes(),
			moduleObservability(),
			moduleSecurity(),
			moduleExam(),
		},
	}
}
