// Package jobs содержит фоновые задачи платформы.
package jobs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"platforma/backend/internal/config"
	"platforma/backend/internal/mailer"
	"platforma/backend/internal/repository"
)

// Reminders раз в несколько часов проверяет сроки прохождения курсов
// и напоминает студентам о дедлайнах и о том, что они давно не заходили.
type Reminders struct {
	repo *repository.ReminderRepo
	mail *mailer.Mailer
	cfg  *config.Config
}

func NewReminders(repo *repository.ReminderRepo, mail *mailer.Mailer, cfg *config.Config) *Reminders {
	return &Reminders{repo: repo, mail: mail, cfg: cfg}
}

// Run запускает цикл проверки. Первый прогон — через минуту после старта.
func (j *Reminders) Run(ctx context.Context, every time.Duration) {
	if every <= 0 {
		every = 6 * time.Hour
	}

	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			j.RunOnce(ctx)
			timer.Reset(every)
		}
	}
}

// RunOnce делает один проход — удобно вызывать и вручную.
func (j *Reminders) RunOnce(ctx context.Context) {
	sent := 0

	deadlines, err := j.repo.DueSoon(ctx)
	if err != nil {
		log.Printf("напоминания: не удалось выбрать дедлайны: %v", err)
	}
	for _, item := range deadlines {
		if j.notify(ctx, item) {
			sent++
		}
	}

	idle, err := j.repo.Idle(ctx)
	if err != nil {
		log.Printf("напоминания: не удалось выбрать неактивных: %v", err)
	}
	for _, item := range idle {
		if j.notify(ctx, item) {
			sent++
		}
	}

	if sent > 0 {
		log.Printf("напоминания: отправлено писем: %d", sent)
	}
}

func (j *Reminders) notify(ctx context.Context, item repository.Reminder) bool {
	subject, message := composeReminder(item, j.cfg.PublicBaseURL)
	link := courseLink(j.cfg.PublicBaseURL, item.CourseSlug)

	if !j.mail.Enabled() {
		log.Printf("[DEV] напоминание для %s: %s", item.Email, subject)
	} else if err := j.mail.SendNotice(ctx, item.Email, subject, message, link); err != nil {
		log.Printf("напоминания: письмо для %s не отправлено: %v", item.Email, err)
		return false
	}

	if err := j.repo.MarkNotified(ctx, item.EnrollmentID, item.Kind); err != nil {
		log.Printf("напоминания: не удалось отметить отправку: %v", err)
	}
	return true
}

// composeReminder формирует текст письма по типу напоминания.
func composeReminder(item repository.Reminder, baseURL string) (string, string) {
	left := item.Total - item.Completed
	link := courseLink(baseURL, item.CourseSlug)

	if item.Kind == "idle" {
		return fmt.Sprintf("Вы давно не заглядывали: %s", item.CourseTitle),
			fmt.Sprintf(
				"%s, вы не заходили на платформу %d дней. По курсу «%s» осталось пройти %d из %d уроков. Продолжить: %s",
				item.Name, item.IdleDays, item.CourseTitle, left, item.Total, link)
	}

	switch {
	case item.DaysLeft > 1:
		return fmt.Sprintf("До конца курса %d дня: %s", item.DaysLeft, item.CourseTitle),
			fmt.Sprintf(
				"%s, до срока прохождения курса «%s» осталось %d дня. Не пройдено уроков: %d из %d. Продолжить: %s",
				item.Name, item.CourseTitle, item.DaysLeft, left, item.Total, link)

	case item.DaysLeft == 1:
		return fmt.Sprintf("Завтра дедлайн: %s", item.CourseTitle),
			fmt.Sprintf(
				"%s, завтра истекает срок прохождения курса «%s». Осталось уроков: %d из %d. Продолжить: %s",
				item.Name, item.CourseTitle, left, item.Total, link)

	case item.DaysLeft == 0:
		return fmt.Sprintf("Сегодня последний день: %s", item.CourseTitle),
			fmt.Sprintf(
				"%s, сегодня последний день прохождения курса «%s». Осталось уроков: %d из %d. Продолжить: %s",
				item.Name, item.CourseTitle, left, item.Total, link)

	default:
		return fmt.Sprintf("Срок прохождения истёк: %s", item.CourseTitle),
			fmt.Sprintf(
				"%s, срок прохождения курса «%s» истёк %d дней назад. Осталось уроков: %d из %d. Курс по-прежнему доступен: %s",
				item.Name, item.CourseTitle, -item.DaysLeft, left, item.Total, link)
	}
}

func courseLink(baseURL, slug string) string {
	return strings.TrimRight(baseURL, "/") + "/learn/courses/" + slug
}
