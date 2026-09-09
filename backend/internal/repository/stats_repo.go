package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatsRepo struct{ db *pgxpool.Pool }

func NewStatsRepo(db *pgxpool.Pool) *StatsRepo { return &StatsRepo{db: db} }

// OnlineWindow — насколько недавно должна быть активность, чтобы считать студента «онлайн».
const OnlineWindow = 5 * time.Minute

// IsOnline проставляет признак «онлайн» по времени последней активности.
func IsOnline(lastSeen *time.Time, now time.Time) bool {
	return lastSeen != nil && now.Sub(*lastSeen) < OnlineWindow
}

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
	OnlineNow       int `json:"onlineNow"`
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
			(SELECT count(DISTINCT user_id) FROM activity_days WHERE day > CURRENT_DATE - 7),
			(SELECT count(*) FROM (
				SELECT user_id, max(last_seen_at) AS ls FROM activity_days GROUP BY user_id
			) t WHERE t.ls > now() - make_interval(mins => $1))
	`, int(OnlineWindow.Minutes())).Scan(&o.Students, &o.ActiveStudents, &o.BlockedStudents, &o.InvitedStudents, &o.Admins,
		&o.Courses, &o.PublishedCourse, &o.Lessons, &o.Enrollments, &o.ActiveToday, &o.ActiveWeek, &o.OnlineNow)
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
	LastSeenAt       *time.Time `json:"lastSeenAt"`
	Online           bool       `json:"online"`
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
			(SELECT max(a.last_seen_at) FROM activity_days a WHERE a.user_id = u.id),
			(SELECT count(*) FROM enrollments e WHERE e.user_id = u.id),
			(SELECT count(*) FROM lessons l
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id
			  WHERE e.user_id = u.id),
			(SELECT count(*) FROM lesson_progress p
			   JOIN lessons l ON l.id = p.lesson_id
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id AND e.user_id = u.id
			  WHERE p.user_id = u.id AND p.status = 'completed'),
			(SELECT count(*) FROM activity_days a WHERE a.user_id = u.id),
			(SELECT COALESCE(sum(a.seconds_spent), 0) / 60 FROM activity_days a WHERE a.user_id = u.id)
		  FROM users u
		 WHERE u.id = $1`, userID).
		Scan(&s.UserID, &s.Email, &s.FullName, &s.Status, &s.LastLoginAt, &s.LastSeenAt,
			&s.Courses, &s.LessonsTotal, &s.LessonsCompleted, &s.DaysVisited, &s.MinutesSpent)
	if err != nil {
		return nil, err
	}
	s.Online = IsOnline(s.LastSeenAt, time.Now())
	if s.LessonsTotal > 0 {
		s.Progress = float64(s.LessonsCompleted) / float64(s.LessonsTotal) * 100
	}
	return &s, nil
}

// CurrentActivity — чем студент занят: последний (или текущий) урок и когда.
type CurrentActivity struct {
	Course string     `json:"course"`
	Lesson string     `json:"lesson"`
	Status string     `json:"status"` // in_progress | completed
	At     *time.Time `json:"at"`
}

// LeaderboardEntry — безопасный (без e-mail) срез статистики студента
// для общей доски достижений, которую видят сами студенты.
type LeaderboardEntry struct {
	UserID           string           `json:"userId"`
	FullName         string           `json:"fullName"`
	LessonsTotal     int              `json:"lessonsTotal"`
	LessonsCompleted int              `json:"lessonsCompleted"`
	Progress         float64          `json:"progress"` // 0..100
	DaysVisited      int              `json:"daysVisited"`
	MinutesSpent     int              `json:"minutesSpent"`
	Certificates     int              `json:"certificates"`
	QuizzesPassed    int              `json:"quizzesPassed"`
	AvgQuizScore     float64          `json:"avgQuizScore"` // 0..100
	Online           bool             `json:"online"`
	LastSeenAt       *time.Time       `json:"lastSeenAt"`
	Current          *CurrentActivity `json:"current"`
}

// CommunityStats — сводные показатели по всем студентам платформы.
type CommunityStats struct {
	Students         int `json:"students"`
	OnlineNow        int `json:"onlineNow"`
	ActiveWeek       int `json:"activeWeek"`
	LessonsCompleted int `json:"lessonsCompleted"`
	Certificates     int `json:"certificates"`
}

// Community — общие показатели платформы, безопасные для показа студентам.
func (r *StatsRepo) Community(ctx context.Context) (*CommunityStats, error) {
	var c CommunityStats
	// «Онлайн» и «за неделю» считаем по той же выборке, что и рейтинг
	// (активные студенты), иначе карточка «Сейчас онлайн» учитывала бы админов
	// и расходилась бы с блоком «Сейчас на платформе».
	err := r.db.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM users WHERE role = 'student'),
			(SELECT count(*) FROM (
				SELECT a.user_id, max(a.last_seen_at) AS ls
				  FROM activity_days a
				  JOIN users u ON u.id = a.user_id
				 WHERE u.role = 'student' AND u.status = 'active'
				 GROUP BY a.user_id
			) t WHERE t.ls > now() - make_interval(mins => $1)),
			(SELECT count(DISTINCT a.user_id) FROM activity_days a
			   JOIN users u ON u.id = a.user_id
			  WHERE u.role = 'student' AND u.status = 'active' AND a.day > CURRENT_DATE - 7),
			(SELECT count(*) FROM lesson_progress WHERE status = 'completed'),
			(SELECT count(*) FROM certificates WHERE revoked_at IS NULL)
	`, int(OnlineWindow.Minutes())).Scan(&c.Students, &c.OnlineNow, &c.ActiveWeek, &c.LessonsCompleted, &c.Certificates)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Leaderboard — рейтинг студентов по числу пройденных уроков.
// Возвращает только безопасные поля (без e-mail и статуса аккаунта):
// имя, прогресс, достижения, признак онлайн, последний визит и чем занят.
func (r *StatsRepo) Leaderboard(ctx context.Context, limit int) ([]LeaderboardEntry, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.full_name,
			ls.last_seen,
			(SELECT count(*) FROM lessons l
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id
			  WHERE e.user_id = u.id) AS lessons_total,
			(SELECT count(*) FROM lesson_progress p
			   JOIN lessons l ON l.id = p.lesson_id
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id AND e.user_id = u.id
			  WHERE p.user_id = u.id AND p.status = 'completed') AS lessons_done,
			(SELECT count(*) FROM activity_days a WHERE a.user_id = u.id),
			(SELECT COALESCE(sum(a.seconds_spent), 0) / 60 FROM activity_days a WHERE a.user_id = u.id),
			(SELECT count(*) FROM certificates c WHERE c.user_id = u.id AND c.revoked_at IS NULL),
			(SELECT count(*) FROM lesson_attempts la WHERE la.user_id = u.id AND la.kind = 'quiz' AND la.passed),
			(SELECT COALESCE(avg(la.score), 0) FROM lesson_attempts la WHERE la.user_id = u.id AND la.kind = 'quiz'),
			cur.course_title, cur.lesson_title, cur.status, cur.updated_at
		  FROM users u
		  LEFT JOIN LATERAL (
		      SELECT max(a.last_seen_at) AS last_seen FROM activity_days a WHERE a.user_id = u.id
		  ) ls ON true
		  LEFT JOIN LATERAL (
		      SELECT c.title AS course_title, l.title AS lesson_title, p.status, p.updated_at
		        FROM lesson_progress p
		        JOIN lessons l ON l.id = p.lesson_id
		        JOIN modules m ON m.id = l.module_id
		        JOIN courses c ON c.id = m.course_id
		       WHERE p.user_id = u.id
		       ORDER BY p.updated_at DESC
		       LIMIT 1
		  ) cur ON true
		 WHERE u.role = 'student' AND u.status = 'active'
		 ORDER BY lessons_done DESC, u.created_at ASC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	now := time.Now()
	out := make([]LeaderboardEntry, 0, limit)
	for rows.Next() {
		var e LeaderboardEntry
		var curCourse, curLesson, curStatus *string
		var curAt *time.Time
		if err := rows.Scan(&e.UserID, &e.FullName, &e.LastSeenAt,
			&e.LessonsTotal, &e.LessonsCompleted, &e.DaysVisited, &e.MinutesSpent,
			&e.Certificates, &e.QuizzesPassed, &e.AvgQuizScore,
			&curCourse, &curLesson, &curStatus, &curAt); err != nil {
			return nil, err
		}
		e.Online = IsOnline(e.LastSeenAt, now)
		if e.LessonsTotal > 0 {
			e.Progress = float64(e.LessonsCompleted) / float64(e.LessonsTotal) * 100
		}
		if curLesson != nil {
			e.Current = &CurrentActivity{
				Course: strOrEmpty(curCourse),
				Lesson: *curLesson,
				Status: strOrEmpty(curStatus),
				At:     curAt,
			}
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// AdminActivityRow — подробная строка живой активности для администратора.
// В отличие от доски достижений содержит e-mail и точное время.
type AdminActivityRow struct {
	UserID       string           `json:"userId"`
	Email        string           `json:"email"`
	FullName     string           `json:"fullName"`
	Status       string           `json:"status"` // статус аккаунта
	Online       bool             `json:"online"`
	LastSeenAt   *time.Time       `json:"lastSeenAt"`
	MinutesToday int              `json:"minutesToday"`
	DaysVisited  int              `json:"daysVisited"`
	LessonsDone  int              `json:"lessonsCompleted"`
	Current      *CurrentActivity `json:"current"`
}

// LiveActivity — все студенты по свежести активности: кто онлайн и чем занят.
func (r *StatsRepo) LiveActivity(ctx context.Context, limit int) ([]AdminActivityRow, error) {
	if limit <= 0 || limit > 1000 {
		limit = 300
	}
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.email, u.full_name, u.status,
			ls.last_seen,
			(SELECT COALESCE(sum(a.seconds_spent), 0) / 60 FROM activity_days a
			  WHERE a.user_id = u.id AND a.day = CURRENT_DATE) AS minutes_today,
			(SELECT count(*) FROM activity_days a WHERE a.user_id = u.id) AS days_visited,
			(SELECT count(*) FROM lesson_progress p WHERE p.user_id = u.id AND p.status = 'completed') AS lessons_done,
			cur.course_title, cur.lesson_title, cur.status, cur.updated_at
		  FROM users u
		  LEFT JOIN LATERAL (
		      SELECT max(a.last_seen_at) AS last_seen FROM activity_days a WHERE a.user_id = u.id
		  ) ls ON true
		  LEFT JOIN LATERAL (
		      SELECT c.title AS course_title, l.title AS lesson_title, p.status, p.updated_at
		        FROM lesson_progress p
		        JOIN lessons l ON l.id = p.lesson_id
		        JOIN modules m ON m.id = l.module_id
		        JOIN courses c ON c.id = m.course_id
		       WHERE p.user_id = u.id
		       ORDER BY p.updated_at DESC
		       LIMIT 1
		  ) cur ON true
		 WHERE u.role = 'student'
		 ORDER BY ls.last_seen DESC NULLS LAST, u.created_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	now := time.Now()
	out := make([]AdminActivityRow, 0, limit)
	for rows.Next() {
		var a AdminActivityRow
		var curCourse, curLesson, curStatus *string
		var curAt *time.Time
		if err := rows.Scan(&a.UserID, &a.Email, &a.FullName, &a.Status,
			&a.LastSeenAt, &a.MinutesToday, &a.DaysVisited, &a.LessonsDone,
			&curCourse, &curLesson, &curStatus, &curAt); err != nil {
			return nil, err
		}
		a.Online = IsOnline(a.LastSeenAt, now)
		if curLesson != nil {
			a.Current = &CurrentActivity{
				Course: strOrEmpty(curCourse),
				Lesson: *curLesson,
				Status: strOrEmpty(curStatus),
				At:     curAt,
			}
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// StudentsSummary — та же статистика по всем студентам (для таблицы успеваемости).
func (r *StatsRepo) StudentsSummary(ctx context.Context, limit int) ([]StudentSummary, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.email, u.full_name, u.status, u.last_login_at,
			(SELECT max(a.last_seen_at) FROM activity_days a WHERE a.user_id = u.id),
			(SELECT count(*) FROM enrollments e WHERE e.user_id = u.id),
			(SELECT count(*) FROM lessons l
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id
			  WHERE e.user_id = u.id),
			(SELECT count(*) FROM lesson_progress p
			   JOIN lessons l ON l.id = p.lesson_id
			   JOIN modules m ON m.id = l.module_id
			   JOIN enrollments e ON e.course_id = m.course_id AND e.user_id = u.id
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

	now := time.Now()
	out := make([]StudentSummary, 0, limit)
	for rows.Next() {
		var s StudentSummary
		if err := rows.Scan(&s.UserID, &s.Email, &s.FullName, &s.Status, &s.LastLoginAt, &s.LastSeenAt,
			&s.Courses, &s.LessonsTotal, &s.LessonsCompleted, &s.DaysVisited, &s.MinutesSpent); err != nil {
			return nil, err
		}
		s.Online = IsOnline(s.LastSeenAt, now)
		if s.LessonsTotal > 0 {
			s.Progress = float64(s.LessonsCompleted) / float64(s.LessonsTotal) * 100
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
