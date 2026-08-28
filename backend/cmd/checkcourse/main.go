// Команда checkcourse проверяет качество курса перед публикацией:
// решаемость заданий, корректность вопросов, наличие материалов и читаемость текста.
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
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"platforma/backend/internal/config"
	"platforma/backend/internal/db"
	"platforma/backend/internal/domain"
	"platforma/backend/internal/dotenv"
	"platforma/backend/internal/repository"
)

// maxSentence — предел длины предложения в символах.
// Всё, что длиннее, новичку читать тяжело.
const maxSentence = 160

func main() {
	slug := flag.String("slug", "devops-engineer", "адрес курса для проверки")
	verbose := flag.Bool("v", false, "показать разбор по каждому уроку")
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
		log.Fatalf("загрузка содержимого: %v", err)
	}

	var problems []string
	kinds := map[string]int{}
	var minutes, links, sentences, sentenceLen int

	for _, module := range course.Modules {
		if len(module.Lessons) == 0 {
			problems = append(problems, fmt.Sprintf("модуль «%s»: нет уроков", module.Title))
		}

		for _, lesson := range module.Lessons {
			kinds[lesson.Kind]++
			minutes += lesson.DurationMin

			report := checkLesson(lesson)
			problems = append(problems, report.problems...)
			links += report.links
			sentences += report.sentences
			sentenceLen += report.sentenceLen

			if *verbose {
				status := "ok"
				if len(report.problems) > 0 {
					status = fmt.Sprintf("%d замечаний", len(report.problems))
				}
				fmt.Printf("  %-46s %-9s %s\n", trim(lesson.Title, 45), lesson.Kind, status)
			}
		}
	}

	avg := 0
	if sentences > 0 {
		avg = sentenceLen / sentences
	}

	fmt.Printf("\nКурс: %s\n", course.Title)
	fmt.Printf("Модулей: %d · уроков: %d (теория %d, квизы %d, тренажёры %d, практики %d)\n",
		len(course.Modules), kinds["text"]+kinds["quiz"]+kinds["terminal"]+kinds["code"],
		kinds["text"], kinds["quiz"], kinds["terminal"], kinds["code"])
	fmt.Printf("Объём: ~%d ч %d мин · ссылок на материалы: %d · средняя длина предложения: %d знаков\n",
		minutes/60, minutes%60, links, avg)

	if len(problems) == 0 {
		fmt.Println("\n✓ Замечаний нет: задания решаемы, вопросы корректны, материалы на месте.")
		return
	}

	fmt.Printf("\n✗ Замечаний: %d\n", len(problems))
	for _, p := range problems {
		fmt.Println("  -", p)
	}
	os.Exit(1)
}

type lessonReport struct {
	problems    []string
	links       int
	sentences   int
	sentenceLen int
}

func checkLesson(lesson domain.Lesson) lessonReport {
	var out lessonReport
	add := func(format string, args ...any) {
		out.problems = append(out.problems, fmt.Sprintf("%s: %s", lesson.Title, fmt.Sprintf(format, args...)))
	}

	// Материалы по теме нужны в каждом уроке.
	var meta struct {
		Resources []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"resources"`
	}
	_ = json.Unmarshal(lesson.Content, &meta)
	out.links = len(meta.Resources)
	if out.links == 0 {
		add("нет блока «Материалы по теме»")
	}
	for _, item := range meta.Resources {
		if !strings.HasPrefix(item.URL, "https://") {
			add("ссылка не по https — %s", item.URL)
		}
		if strings.TrimSpace(item.Title) == "" {
			add("у ссылки нет названия")
		}
	}

	switch lesson.Kind {
	case domain.LessonText:
		var content struct {
			Body string `json:"body"`
		}
		if err := json.Unmarshal(lesson.Content, &content); err != nil {
			add("не удалось разобрать содержимое")
			return out
		}
		if utf8.RuneCountInString(content.Body) < 900 {
			add("слишком короткая теория (%d символов)", utf8.RuneCountInString(content.Body))
		}
		if !strings.Contains(content.Body, "## Зачем это нужно") {
			add("нет блока «Зачем это нужно»")
		}
		if !strings.Contains(content.Body, "## Запомнить") {
			add("нет блока «Запомнить»")
		}

		count, total, long := analyzeText(content.Body)
		out.sentences, out.sentenceLen = count, total
		for _, sentence := range long {
			add("длинное предложение (%d знаков): %s…", utf8.RuneCountInString(sentence), trim(sentence, 60))
		}

	case domain.LessonQuiz:
		quiz, err := domain.ParseQuiz(lesson.Content)
		if err != nil {
			add("не удалось разобрать квиз")
			return out
		}
		if len(quiz.Questions) < 4 {
			add("мало вопросов (%d)", len(quiz.Questions))
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
				add("вопрос %s без правильного ответа", question.ID)
			case correct > 1 && !question.Multiple:
				add("вопрос %s: несколько верных вариантов без multiple", question.ID)
			}
			if len(question.Options) < 3 {
				add("вопрос %s: меньше трёх вариантов", question.ID)
			}
			if strings.TrimSpace(question.Explanation) == "" {
				add("вопрос %s без пояснения", question.ID)
			}
		}

		// Эталонные ответы обязаны проходить порог.
		answers := make([]domain.QuizAnswer, 0, len(quiz.Questions))
		for _, question := range quiz.Questions {
			ids := make([]string, 0, 2)
			for _, option := range question.Options {
				if option.Correct {
					ids = append(ids, option.ID)
				}
			}
			answers = append(answers, domain.QuizAnswer{QuestionID: question.ID, OptionIDs: ids})
		}
		if result := domain.GradeQuiz(quiz, answers); !result.Passed {
			add("эталонные ответы не проходят порог (%.0f%%)", result.Score)
		}

	case domain.LessonTerminal:
		terminal, err := domain.ParseTerminal(lesson.Content)
		if err != nil {
			add("не удалось разобрать задания")
			return out
		}
		if len(terminal.Tasks) == 0 {
			add("нет заданий")
		}
		for _, task := range terminal.Tasks {
			if len(task.Expected) == 0 && task.Pattern == "" {
				add("задание %s без эталонной команды", task.ID)
				continue
			}
			if strings.TrimSpace(task.Hint) == "" {
				add("задание %s без подсказки", task.ID)
			}
			for _, expected := range task.Expected {
				if !domain.MatchTerminalCommand(&task, expected) {
					add("задание %s: команда «%s» не проходит собственную проверку", task.ID, expected)
				}
			}
		}

	case domain.LessonCode:
		code, err := domain.ParseCode(lesson.Content)
		if err != nil {
			add("не удалось разобрать практику")
			return out
		}
		if strings.TrimSpace(code.Task) == "" {
			add("нет текста задания")
		}
		if strings.TrimSpace(code.Solution) == "" {
			add("нет эталонного решения")
			return out
		}
		for _, check := range code.Checks {
			if !domain.RunCodeCheck(check, code.Solution) {
				message := check.Message
				if message == "" {
					message = domain.DescribeCodeCheck(check)
				}
				add("эталонное решение не проходит проверку «%s»", message)
			}
		}
	}

	return out
}

var (
	codeBlock = regexp.MustCompile("(?s)```.*?```")
	tableRow  = regexp.MustCompile(`(?m)^\|.*\|$`)
	listItem  = regexp.MustCompile(`(?m)^\s*([-*]|\d+\.)\s`)
	heading   = regexp.MustCompile(`(?m)^#{1,6}\s`)
	markup    = regexp.MustCompile("[#>*`]")
)

// analyzeText считает предложения обычного текста: блоки кода, таблицы,
// списки и заголовки не участвуют — их длина ничего не говорит о читаемости.
func analyzeText(body string) (count, total int, long []string) {
	text := codeBlock.ReplaceAllString(body, "")
	text = tableRow.ReplaceAllString(text, "")

	for _, paragraph := range strings.Split(text, "\n") {
		line := strings.TrimSpace(paragraph)
		if line == "" || heading.MatchString(line) || listItem.MatchString(line) {
			continue
		}
		line = markup.ReplaceAllString(line, "")

		for _, sentence := range strings.FieldsFunc(line, func(r rune) bool {
			return r == '.' || r == '!' || r == '?'
		}) {
			sentence = strings.TrimSpace(sentence)
			if utf8.RuneCountInString(sentence) < 15 {
				continue
			}
			length := utf8.RuneCountInString(sentence)
			count++
			total += length
			if length > maxSentence {
				long = append(long, sentence)
			}
		}
	}
	return count, total, long
}

func trim(s string, limit int) string {
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit])
}
