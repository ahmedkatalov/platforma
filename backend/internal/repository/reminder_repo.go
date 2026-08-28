package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Reminder — кандидат на письмо-напоминание.
type Reminder struct {
	EnrollmentID string
	UserID       string
	CourseID     string
	Email        string
	Name         string
	CourseTitle  string
	CourseSlug   string
	DueDate      *time.Time
	DaysLeft     int
	IdleDays     int
	Total        int
	Completed    int
	Kind         string // deadline | idle
}

type ReminderRepo struct{ db *pgxpool.Pool }

func NewReminderRepo(db *pgxpool.Pool) *ReminderRepo { return &ReminderRepo{db: db} }

// DueSoon возвращает записи, по которым пора напомнить о сроке:
// за 3 дня, за день, в день дедлайна и после просрочки. Не чаще раза в сутки.
func (r *ReminderRepo) DueSoon(ctx context.Context) ([]Reminder, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.user_id, e.course_id, u.email,
		       COALESCE(NULLIF(u.full_name, ''), u.email), c.title, c.slug, e.due_date,
		       (e.due_date - CURRENT_DATE) AS days_left,
		       (SELECT count(*) FROM lessons l
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id),
		       (SELECT count(*) FROM lesson_progress p
		          JOIN lessons l ON l.id = p.lesson_id
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id AND p.user_id = u.id AND p.status = 'completed')
		  FROM enrollments e
		  JOIN users u ON u.id = e.user_id
		  JOIN courses c ON c.id = e.course_id
		 WHERE e.due_date IS NOT NULL
		   AND e.status <> 'completed'
		   AND u.status = 'active'
		   AND (e.due_date - CURRENT_DATE) IN (3, 1, 0, -1, -7)
		   AND (e.deadline_notified_at IS NULL OR e.deadline_notified_at < now() - interval '20 hours')
		   AND NOT EXISTS (
		       SELECT 1 FROM certificates cert
		        WHERE cert.user_id = e.user_id AND cert.course_id = e.course_id
		   )`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReminders(rows, "deadline")
}

// Idle возвращает студентов, которые начали курс, но не заходили неделю.
func (r *ReminderRepo) Idle(ctx context.Context) ([]Reminder, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.user_id, e.course_id, u.email,
		       COALESCE(NULLIF(u.full_name, ''), u.email), c.title, c.slug, e.due_date,
		       COALESCE(EXTRACT(DAY FROM now() - u.last_login_at)::int, 0),
		       (SELECT count(*) FROM lessons l
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id),
		       (SELECT count(*) FROM lesson_progress p
		          JOIN lessons l ON l.id = p.lesson_id
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id AND p.user_id = u.id AND p.status = 'completed')
		  FROM enrollments e
		  JOIN users u ON u.id = e.user_id
		  JOIN courses c ON c.id = e.course_id
		 WHERE e.status <> 'completed'
		   AND u.status = 'active'
		   AND u.last_login_at IS NOT NULL
		   AND u.last_login_at < now() - interval '7 days'
		   AND (e.idle_notified_at IS NULL OR e.idle_notified_at < now() - interval '7 days')
		   AND NOT EXISTS (
		       SELECT 1 FROM certificates cert
		        WHERE cert.user_id = e.user_id AND cert.course_id = e.course_id
		   )`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReminders(rows, "idle")
}

func scanReminders(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}, kind string) ([]Reminder, error) {
	out := make([]Reminder, 0, 8)

	for rows.Next() {
		var item Reminder
		var metric int

		if err := rows.Scan(&item.EnrollmentID, &item.UserID, &item.CourseID, &item.Email,
			&item.Name, &item.CourseTitle, &item.CourseSlug, &item.DueDate, &metric,
			&item.Total, &item.Completed); err != nil {
			return nil, err
		}

		item.Kind = kind
		if kind == "deadline" {
			item.DaysLeft = metric
		} else {
			item.IdleDays = metric
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

// MarkNotified отмечает, что письмо отправлено, чтобы не слать повторно.
func (r *ReminderRepo) MarkNotified(ctx context.Context, enrollmentID, kind string) error {
	column := "deadline_notified_at"
	if kind == "idle" {
		column = "idle_notified_at"
	}
	_, err := r.db.Exec(ctx,
		`UPDATE enrollments SET `+column+` = now() WHERE id = $1`, enrollmentID)
	return err
}
