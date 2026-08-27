package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepo struct{ db *pgxpool.Pool }

func NewStatsRepo(db *pgxpool.Pool) *StatsRepo { return &StatsRepo{db: db} }

// AdminOverview — сводка для главной страницы администратора.
type AdminOverview struct {
	Students        int `json:"students"`
	ActiveStudents  int `json:"activeStudents"`
	BlockedStudents int `json:"blockedStudents"`
	InvitedStudents int `json:"invitedStudents"`
	Admins          int `json:"admins"`
	Courses         int `json:"courses"`
	PublishedCourse int `json:"publishedCourses"`
	Lessons         int `json:"lessons"`
	Enrollments     int `json:"enrollments"`
	ActiveToday     int `json:"activeToday"`
	ActiveWeek      int `json:"activeWeek"`
}

func (r *StatsRepo) Overview(ctx context.Context) (*AdminOverview, error) {
	var o AdminOverview
	err := r.db.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM users WHERE role = 'student'),
			(SELECT count(*) FROM users WHERE role = 'student' AND status = 'active'),
			(SELECT count(*) FROM users WHERE role = 'student' AND status = 'blocked'),
			(SELECT count(*) FROM users WHERE role = 'student' AND status = 'invited'),
			(SELECT count(*) FROM users WHERE role = 'admin'),
			(SELECT count(*) FROM courses),
			(SELECT count(*) FROM courses WHERE status = 'published'),
			(SELECT count(*) FROM lessons),
			(SELECT count(*) FROM enrollments),
			(SELECT count(*) FROM activity_days WHERE day = CURRENT_DATE),
			(SELECT count(DISTINCT user_id) FROM activity_days WHERE day > CURRENT_DATE - 7)
	`).Scan(&o.Students, &o.ActiveStudents, &o.BlockedStudents, &o.InvitedStudents, &o.Admins,
		&o.Courses, &o.PublishedCourse, &o.Lessons, &o.Enrollments, &o.ActiveToday, &o.ActiveWeek)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// StudentSummary — краткая статистика студента для списков и его дашборда.
type StudentSummary struct {
	UserID           string     `json:"userId"`
	Email            string     `json:"email"`
	FullName         string     `json:"fullName"`
	Status           string     `json:"status"`
	LastLoginAt      *time.Time `json:"lastLoginAt"`
	Courses          int        `json:"courses"`
	LessonsTotal     int        `json:"lessonsTotal"`
	LessonsCompleted int        `json:"lessonsCompleted"`
	DaysVisited      int        `json:"daysVisited"`
	MinutesSpent     int        `json:"minutesSpent"`
	Progress         float64    `json:"progress"` // 0..100
}

func (r *StatsRepo) StudentSummary(ctx context.Context, userID string) (*StudentSummary, error) {
	var s StudentSummary
	err := r.db.QueryRow(ctx, `
		SELECT u.id, u.email, u.full_name, u.status, u.last_login_at,
			(SELECT count(*) FROM enrollments e WHERE e.user_id = u.id),
			(SELECT count(*) FROM lessons l
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id
			  WHERE e.user_id = u.id),
			(SELECT count(*) FROM lesson_progress p
			  WHERE p.user_id = u.id AND p.status = 'completed'),
			(SELECT count(*) FROM activity_days a WHERE a.user_id = u.id),
			(SELECT COALESCE(sum(a.seconds_spent), 0) / 60 FROM activity_days a WHERE a.user_id = u.id)
		  FROM users u
		 WHERE u.id = $1`, userID).
		Scan(&s.UserID, &s.Email, &s.FullName, &s.Status, &s.LastLoginAt,
			&s.Courses, &s.LessonsTotal, &s.LessonsCompleted, &s.DaysVisited, &s.MinutesSpent)
	if err != nil {
		return nil, err
	}
	if s.LessonsTotal > 0 {
		s.Progress = float64(s.LessonsCompleted) / float64(s.LessonsTotal) * 100
	}
	return &s, nil
}

// StudentsSummary — та же статистика по всем студентам (для таблицы успеваемости).
func (r *StatsRepo) StudentsSummary(ctx context.Context, limit int) ([]StudentSummary, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.email, u.full_name, u.status, u.last_login_at,
			(SELECT count(*) FROM enrollments e WHERE e.user_id = u.id),
			(SELECT count(*) FROM lessons l
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id
			  WHERE e.user_id = u.id),
			(SELECT count(*) FROM lesson_progress p
			  WHERE p.user_id = u.id AND p.status = 'completed'),
			(SELECT count(*) FROM activity_days a WHERE a.user_id = u.id),
			(SELECT COALESCE(sum(a.seconds_spent), 0) / 60 FROM activity_days a WHERE a.user_id = u.id)
		  FROM users u
		 WHERE u.role = 'student'
		 ORDER BY u.created_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]StudentSummary, 0, limit)
	for rows.Next() {
		var s StudentSummary
		if err := rows.Scan(&s.UserID, &s.Email, &s.FullName, &s.Status, &s.LastLoginAt,
			&s.Courses, &s.LessonsTotal, &s.LessonsCompleted, &s.DaysVisited, &s.MinutesSpent); err != nil {
			return nil, err
		}
		if s.LessonsTotal > 0 {
			s.Progress = float64(s.LessonsCompleted) / float64(s.LessonsTotal) * 100
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
