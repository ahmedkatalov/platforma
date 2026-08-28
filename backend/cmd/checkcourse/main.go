// Команда checkcourse проверяет, что курс в базе пригоден для прохождения:
// у каждого вопроса есть верный ответ и пояснение, у каждого задания терминала —
// эталонная команда, а эталонное решение практики проходит все проверки.
//
//	go run ./cmd/checkcourse
//	go run ./cmd/checkcourse -slug devops-engineer
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"platforma/backend/internal/config"
	"platforma/backend/internal/db"
	"platforma/backend/internal/domain"
	"platforma/backend/internal/dotenv"
	"platforma/backend/internal/repository"
)

const minBodyLength = 400

func main() {
	slug := flag.String("slug", "devops-engineer", "адрес курса (slug)")
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

	courses := repository.NewCourseRepo(pool)

	course, err := courses.GetBySlug(ctx, *slug)
	if err != nil {
		log.Fatalf("курс %q не найден: %v", *slug, err)
	}
	if err := courses.WithContent(ctx, course, true); err != nil {
		log.Fatalf("загрузка структуры курса: %v", err)
	}

	var problems []string
	counts := map[string]int{}
	minutes := 0

	for _, module := range course.Modules {
		if len(module.Lessons) == 0 {
			problems = append(problems, fmt.Sprintf("модуль «%s»: нет уроков", module.Title))
		}

		for _, lesson := range module.Lessons {
			counts[lesson.Kind]++
			minutes += lesson.DurationMin
			problems = append(problems, checkLesson(lesson)...)
		}
	}

	fmt.Printf("Курс «%s»: %d модулей, %d уроков, ~%d ч %d мин\n",
		course.Title, len(course.Modules), course.LessonsCount, minutes/60, minutes%60)
	fmt.Printf("  теория: %d · квизы: %d · терминал: %d · практика: %d\n\n",
		counts[domain.LessonText], counts[domain.LessonQuiz],
		counts[domain.LessonTerminal], counts[domain.LessonCode])

	if len(problems) == 0 {
		fmt.Println("Замечаний нет: курс готов к прохождению.")
		return
	}

	fmt.Printf("Найдено замечаний: %d\n", len(problems))
	for _, problem := range problems {
		fmt.Println(" -", problem)
	}
	os.Exit(1)
}

func checkLesson(lesson domain.Lesson) []string {
	var problems []string
	label := fmt.Sprintf("«%s»", lesson.Title)

	switch lesson.Kind {
	case domain.LessonText:
		var content struct {
			Body string `json:"body"`
		}
		if err := json.Unmarshal(lesson.Content, &content); err != nil {
			return append(problems, label+": содержимое не разбирается")
		}
		if len([]rune(content.Body)) < minBodyLength {
			problems = append(problems, fmt.Sprintf("%s: слишком короткая теория (%d символов)",
				label, len([]rune(content.Body))))
		}

	case domain.LessonQuiz:
		quiz, err := domain.ParseQuiz(lesson.Content)
		if err != nil {
			return append(problems, label+": квиз не разбирается")
		}
		if len(quiz.Questions) < 3 {
			problems = append(problems, fmt.Sprintf("%s: меньше трёх вопросов", label))
		}
		for _, question := range quiz.Questions {
			correct := 0
			for _, option := range question.Options {
				if option.Correct {
					correct++
				}
			}
			switch {
			case correct == 0:
				problems = append(problems, fmt.Sprintf("%s / %s: нет правильного ответа", label, question.ID))
			case correct > 1 && !question.Multiple:
				problems = append(problems, fmt.Sprintf("%s / %s: несколько верных вариантов без multiple", label, question.ID))
			}
			if len(question.Options) < 2 {
				problems = append(problems, fmt.Sprintf("%s / %s: меньше двух вариантов", label, question.ID))
			}
			if question.Explanation == "" {
				problems = append(problems, fmt.Sprintf("%s / %s: нет пояснения", label, question.ID))
			}
		}

	case domain.LessonTerminal:
		terminal, err := domain.ParseTerminal(lesson.Content)
		if err != nil {
			return append(problems, label+": задания не разбираются")
		}
		if len(terminal.Tasks) == 0 {
			problems = append(problems, label+": нет заданий")
		}
		for _, task := range terminal.Tasks {
			if len(task.Expected) == 0 && task.Pattern == "" {
				problems = append(problems, fmt.Sprintf("%s / %s: не задана эталонная команда", label, task.ID))
			}
			if task.Prompt == "" {
				problems = append(problems, fmt.Sprintf("%s / %s: пустая формулировка", label, task.ID))
			}
			if task.Hint == "" {
				problems = append(problems, fmt.Sprintf("%s / %s: нет подсказки", label, task.ID))
			}
		}

	case domain.LessonCode:
		code, err := domain.ParseCode(lesson.Content)
		if err != nil {
			return append(problems, label+": практика не разбирается")
		}
		if code.Task == "" {
			problems = append(problems, label+": нет текста задания")
		}
		if len(code.Checks) == 0 {
			problems = append(problems, label+": нет проверок решения")
		}
		if code.Solution == "" {
			problems = append(problems, label+": нет эталонного решения")
			break
		}
		// Главная проверка: эталон обязан проходить все условия задания.
		for _, check := range code.Checks {
			if !domain.RunCodeCheck(check, code.Solution) {
				problems = append(problems, fmt.Sprintf("%s: эталон не проходит проверку «%s»",
					label, checkLabel(check)))
			}
		}
	}

	return problems
}

func checkLabel(check domain.CodeCheck) string {
	if check.Message != "" {
		return check.Message
	}
	return check.Type + " " + check.Value
}
