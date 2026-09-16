package seed

// GoBackendCourse — второй курс платформы: Go Backend Engineer.
// Модули берутся из content-go/*.json, метаданные — из блока course файла-источника.
func GoBackendCourse() CourseSeed {
	return CourseSeed{
		Slug:     "go-backend-engineer-2026",
		Title:    "Go Backend Engineer 2026: с нуля до production",
		Subtitle: "Go 1.27, REST/gRPC, PostgreSQL 18, Redis, Kafka/NATS, безопасность, observability, Docker, Kubernetes и системный дизайн",
		Description: "Практический курс backend-разработки на Go для начинающих. От основ языка и стандартной " +
			"библиотеки — к боевому production: REST и gRPC, работа с PostgreSQL (pgx, sqlc) и Redis, очереди " +
			"Kafka/NATS, конкурентность, тестирование, безопасность, наблюдаемость (OpenTelemetry), Docker, " +
			"Kubernetes и системный дизайн. Каждая тема — с примерами, кодом в редакторе, тренажёром и квизом " +
			"с проверкой на сервере. В конце — сквозной проект и аттестация.",
		Level: "beginner",
		Tags: []string{
			"golang", "go-1.27", "backend", "rest-api", "grpc",
			"postgresql-18", "pgx", "sqlc", "redis", "kafka", "nats",
			"opentelemetry", "docker", "kubernetes", "security", "system-design",
		},
		Modules: loadGoModules(),
	}
}

// AllCourses — все курсы, которыми наполняется платформа командой seedcourse.
func AllCourses() []CourseSeed {
	return []CourseSeed{
		DevOpsCourse(),
		GoBackendCourse(),
	}
}
