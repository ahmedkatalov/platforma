// Команда seedcourse наполняет платформу демонстрационным курсом по DevOps:
// модули, теория, квизы, задания для тренажёра терминала и практика с кодом.
//
//	go run ./cmd/seedcourse
//	go run ./cmd/seedcourse -force     # пересоздать курс, если он уже есть
//	go run ./cmd/seedcourse -publish=false
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"platforma/backend/internal/config"
	"platforma/backend/internal/db"
	"platforma/backend/internal/dotenv"
	"platforma/backend/internal/repository"
	"platforma/backend/internal/seed"
)

func main() {
	force := flag.Bool("force", false, "пересоздать курс, если он уже существует")
	publish := flag.Bool("publish", true, "сразу опубликовать курс")
	flag.Parse()

	dotenv.Load(".env")
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("подключение к базе: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("миграции: %v", err)
	}

	courses := repository.NewCourseRepo(pool)
	data := seed.DevOpsCourse()

	existing, err := courses.GetBySlug(ctx, data.Slug)
	switch {
	case err == nil:
		if !*force {
			log.Fatalf("курс %q уже существует — запустите с флагом -force, чтобы пересоздать", data.Slug)
		}
		if err := courses.Delete(ctx, existing.ID); err != nil {
			log.Fatalf("удаление старого курса: %v", err)
		}
		fmt.Println("Старая версия курса удалена")
	case errors.Is(err, repository.ErrNotFound):
		// первая загрузка — всё в порядке
	default:
		log.Fatalf("поиск курса: %v", err)
	}

	status := "published"
	if !*publish {
		status = "draft"
	}

	course, err := courses.Create(ctx, repository.CourseInput{
		Slug:        data.Slug,
		Title:       data.Title,
		Subtitle:    data.Subtitle,
		Description: data.Description,
		Level:       data.Level,
		Tags:        data.Tags,
		Status:      status,
		Position:    1,
	}, "")
	if err != nil {
		log.Fatalf("создание курса: %v", err)
	}

	lessonCount := 0
	for moduleIndex, moduleSeed := range data.Modules {
		module, err := courses.CreateModule(ctx, course.ID, repository.ModuleInput{
			Title:    moduleSeed.Title,
			Summary:  moduleSeed.Summary,
			Position: moduleIndex + 1,
		})
		if err != nil {
			log.Fatalf("создание модуля %q: %v", moduleSeed.Title, err)
		}

		for lessonIndex, lessonSeed := range moduleSeed.Lessons {
			content, err := json.Marshal(lessonSeed.Content)
			if err != nil {
				log.Fatalf("сериализация урока %q: %v", lessonSeed.Title, err)
			}

			if _, err := courses.CreateLesson(ctx, module.ID, repository.LessonInput{
				Title:       lessonSeed.Title,
				Kind:        lessonSeed.Kind,
				Summary:     lessonSeed.Summary,
				Content:     content,
				DurationMin: lessonSeed.DurationMin,
				Position:    lessonIndex + 1,
			}); err != nil {
				log.Fatalf("создание урока %q: %v", lessonSeed.Title, err)
			}
			lessonCount++
		}
	}

	fmt.Printf("\n✓ Курс «%s» загружен: %d модулей, %d уроков (статус: %s)\n",
		course.Title, len(data.Modules), lessonCount, status)
	fmt.Printf("  Адрес курса: /learn/courses/%s\n", course.Slug)
	fmt.Println("  Назначьте курс студентам в разделе «Студенты» админки.")
}
