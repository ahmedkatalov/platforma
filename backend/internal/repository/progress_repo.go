package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"platforma/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepo struct{ db *pgxpool.Pool }

func NewProgressRepo(db *pgxpool.Pool) *ProgressRepo { return &ProgressRepo{db: db} }

// LessonContext — урок вместе с курсом, которому он принадлежит.
type LessonContext struct {
	Lesson      domain.Lesson `json:"lesson"`
	CourseID    string        `json:"courseId"`
	CourseSlug  string        `json:"courseSlug"`
	CourseTitle string        `json:"courseTitle"`
	ModuleTitle string        `json:"moduleTitle"`
	Published   bool          `json:"published"`
	PrevID      *string       `json:"prevLessonId"`
	NextID      *string       `json:"nextLessonId"`
}

// LessonWithCourse возвращает урок, его курс и соседние уроки для навигации.
func (r *ProgressRepo) LessonWithCourse(ctx context.Context, lessonID string) (*LessonContext, error) {
	var out LessonContext
	var status string

	err := r.db.QueryRow(ctx, `
		SELECT l.id, l.module_id, l.title, l.kind, l.summary, l.content, l.duration_min, l.position,
		       l.created_at, l.updated_at, c.id, c.slug, c.title, m.title, c.status
		  FROM lessons l
		  JOIN modules m ON m.id = l.module_id
		  JOIN courses c ON c.id = m.course_id
		 WHERE l.id = $1`, lessonID).
		Scan(&out.Lesson.ID, &out.Lesson.ModuleID, &out.Lesson.Title, &out.Lesson.Kind,
			&out.Lesson.Summary, &out.Lesson.Content, &out.Lesson.DurationMin, &out.Lesson.Position,
			&out.Lesson.CreatedAt, &out.Lesson.UpdatedAt,
			&out.CourseID, &out.CourseSlug, &out.CourseTitle, &out.ModuleTitle, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	out.Published = status == "published"

	// Соседние уроки в порядке прохождения курса.
	rows, err := r.db.Query(ctx, `
		SELECT l.id
		  FROM lessons l
		  JOIN modules m ON m.id = l.module_id
		 WHERE m.course_id = $1
		 ORDER BY m.position, m.created_at, l.position, l.created_at`, out.CourseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordered := make([]string, 0, 16)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ordered = append(ordered, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, id := range ordered {
		if id != lessonID {
			continue
		}
		if i > 0 {
			out.PrevID = &ordered[i-1]
		}
		if i < len(ordered)-1 {
			out.NextID = &ordered[i+1]
		}
		break
	}

	return &out, nil
}

// LessonProgress — прогресс студента по одному уроку.
type LessonProgress struct {
	LessonID     string     `json:"lessonId"`
	Status       string     `json:"status"`
	Score        *float64   `json:"score"`
	BestScore    *float64   `json:"bestScore"`
	Attempts     int        `json:"attempts"`
	SecondsSpent int        `json:"secondsSpent"`
	CompletedAt  *time.Time `json:"completedAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

// Start отмечает, что студент открыл урок.
func (r *ProgressRepo) Start(ctx context.Context, userID, lessonID string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO lesson_progress (user_id, lesson_id, status)
		VALUES ($1, $2, 'in_progress')
		ON CONFLICT (user_id, lesson_id) DO NOTHING`, userID, lessonID)
	return err
}

// Complete отмечает урок пройденным и обновляет лучший результат.
func (r *ProgressRepo) Complete(ctx context.Context, userID, lessonID string, score *float64, seconds int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO lesson_progress (user_id, lesson_id, status, score, best_score, seconds_spent,
		                             attempts, completed_at, updated_at)
		VALUES ($1, $2, 'completed', $3, $3, $4, 1, now(), now())
		ON CONFLICT (user_id, lesson_id) DO UPDATE
		   SET status = 'completed',
		       score = EXCLUDED.score,
		       best_score = GREATEST(COALESCE(lesson_progress.best_score, 0), COALESCE(EXCLUDED.score, 0)),
		       seconds_spent = lesson_progress.seconds_spent + EXCLUDED.seconds_spent,
		       attempts = lesson_progress.attempts + 1,
		       completed_at = COALESCE(lesson_progress.completed_at, now()),
		       updated_at = now()`, userID, lessonID, score, seconds)
	return err
}

// TouchAttempt засчитывает попытку, не отмечая урок пройденным.
func (r *ProgressRepo) TouchAttempt(ctx context.Context, userID, lessonID string, score *float64, seconds int) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO lesson_progress (user_id, lesson_id, status, score, best_score, seconds_spent,
		                             attempts, updated_at)
		VALUES ($1, $2, 'in_progress', $3, $3, $4, 1, now())
		ON CONFLICT (user_id, lesson_id) DO UPDATE
		   SET score = EXCLUDED.score,
		       best_score = GREATEST(COALESCE(lesson_progress.best_score, 0), COALESCE(EXCLUDED.score, 0)),
		       seconds_spent = lesson_progress.seconds_spent + EXCLUDED.seconds_spent,
		       attempts = lesson_progress.attempts + 1,
		       updated_at = now()`, userID, lessonID, score, seconds)
	return err
}

// ForCourse — прогресс студента по всем урокам курса.
func (r *ProgressRepo) ForCourse(ctx context.Context, userID, courseID string) ([]LessonProgress, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.lesson_id, p.status, p.score, p.best_score, p.attempts, p.seconds_spent,
		       p.completed_at, p.updated_at
		  FROM lesson_progress p
		  JOIN lessons l ON l.id = p.lesson_id
		  JOIN modules m ON m.id = l.module_id
		 WHERE p.user_id = $1 AND m.course_id = $2`, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProgress(rows)
}

// ForUser — весь прогресс студента (для кабинета).
func (r *ProgressRepo) ForUser(ctx context.Context, userID string) ([]LessonProgress, error) {
	rows, err := r.db.Query(ctx, `
		SELECT lesson_id, status, score, best_score, attempts, seconds_spent, completed_at, updated_at
		  FROM lesson_progress WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProgress(rows)
}

func scanProgress(rows pgx.Rows) ([]LessonProgress, error) {
	out := make([]LessonProgress, 0, 16)
	for rows.Next() {
		var p LessonProgress
		if err := rows.Scan(&p.LessonID, &p.Status, &p.Score, &p.BestScore, &p.Attempts,
			&p.SecondsSpent, &p.CompletedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// SaveAttempt записывает попытку прохождения и возвращает её id.
func (r *ProgressRepo) SaveAttempt(
	ctx context.Context,
	userID, lessonID, kind string,
	score float64,
	correct, total int,
	passed bool,
	seconds int,
	details any,
) (string, error) {
	raw := []byte("{}")
	if details != nil {
		if encoded, err := json.Marshal(details); err == nil {
			raw = encoded
		}
	}

	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO lesson_attempts (user_id, lesson_id, kind, score, correct_count, total_count,
		                             passed, duration_seconds, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		userID, lessonID, kind, score, correct, total, passed, seconds, raw).Scan(&id)
	return id, err
}

// SaveQuizAnswers сохраняет ответы по вопросам — из них считается статистика.
func (r *ProgressRepo) SaveQuizAnswers(
	ctx context.Context,
	attemptID, userID, lessonID string,
	results []domain.QuestionResult,
	timings map[string]int,
) error {
	if len(results) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, result := range results {
		batch.Queue(`
			INSERT INTO quiz_answers (attempt_id, user_id, lesson_id, question_id, correct, seconds_spent)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			attemptID, userID, lessonID, result.QuestionID, result.Correct, timings[result.QuestionID])
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for range results {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

// Attempt — попытка в списке истории.
type Attempt struct {
	ID           string          `json:"id"`
	LessonID     string          `json:"lessonId"`
	LessonTitle  string          `json:"lessonTitle"`
	CourseTitle  string          `json:"courseTitle"`
	Kind         string          `json:"kind"`
	Score        float64         `json:"score"`
	CorrectCount int             `json:"correctCount"`
	TotalCount   int             `json:"totalCount"`
	Passed       bool            `json:"passed"`
	Seconds      int             `json:"durationSeconds"`
	Details      json.RawMessage `json:"details"`
	CreatedAt    time.Time       `json:"createdAt"`
}

func (r *ProgressRepo) Attempts(ctx context.Context, userID string, limit int) ([]Attempt, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.lesson_id, l.title, c.title, a.kind, a.score, a.correct_count, a.total_count,
		       a.passed, a.duration_seconds, a.details, a.created_at
		  FROM lesson_attempts a
		  JOIN lessons l ON l.id = a.lesson_id
		  JOIN modules m ON m.id = l.module_id
		  JOIN courses c ON c.id = m.course_id
		 WHERE a.user_id = $1
		 ORDER BY a.created_at DESC
		 LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Attempt, 0, limit)
	for rows.Next() {
		var a Attempt
		if err := rows.Scan(&a.ID, &a.LessonID, &a.LessonTitle, &a.CourseTitle, &a.Kind, &a.Score,
			&a.CorrectCount, &a.TotalCount, &a.Passed, &a.Seconds, &a.Details, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// QuizStats — сводка по квизам студента.
type QuizStats struct {
	Attempts        int     `json:"attempts"`
	Passed          int     `json:"passed"`
	AverageScore    float64 `json:"averageScore"`
	BestScore       float64 `json:"bestScore"`
	Accuracy        float64 `json:"accuracy"` // доля верных ответов, %
	AvgSecondsPerQ  float64 `json:"avgSecondsPerQuestion"`
	FastestSeconds  int     `json:"fastestSeconds"`
	AnsweredTotal   int     `json:"answeredTotal"`
	AnsweredCorrect int     `json:"answeredCorrect"`
}

func (r *ProgressRepo) QuizStats(ctx context.Context, userID string) (*QuizStats, error) {
	var s QuizStats

	err := r.db.QueryRow(ctx, `
		SELECT count(*),
		       count(*) FILTER (WHERE passed),
		       COALESCE(avg(score), 0),
		       COALESCE(max(score), 0)
		  FROM lesson_attempts
		 WHERE user_id = $1 AND kind = 'quiz'`, userID).
		Scan(&s.Attempts, &s.Passed, &s.AverageScore, &s.BestScore)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `
		SELECT count(*),
		       count(*) FILTER (WHERE correct),
		       COALESCE(avg(seconds_spent), 0),
		       COALESCE(min(NULLIF(seconds_spent, 0)), 0)
		  FROM quiz_answers
		 WHERE user_id = $1`, userID).
		Scan(&s.AnsweredTotal, &s.AnsweredCorrect, &s.AvgSecondsPerQ, &s.FastestSeconds)
	if err != nil {
		return nil, err
	}

	if s.AnsweredTotal > 0 {
		s.Accuracy = float64(s.AnsweredCorrect) / float64(s.AnsweredTotal) * 100
	}
	return &s, nil
}

// TaskState — состояние одного задания внутри урока.
type TaskState struct {
	TaskID      string     `json:"taskId"`
	Attempts    int        `json:"attempts"`
	HintsUsed   int        `json:"hintsUsed"`
	CompletedAt *time.Time `json:"completedAt"`
}

func (r *ProgressRepo) Tasks(ctx context.Context, userID, lessonID string) ([]TaskState, error) {
	rows, err := r.db.Query(ctx, `
		SELECT task_id, attempts, hints_used, completed_at
		  FROM task_progress
		 WHERE user_id = $1 AND lesson_id = $2`, userID, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]TaskState, 0, 8)
	for rows.Next() {
		var t TaskState
		if err := rows.Scan(&t.TaskID, &t.Attempts, &t.HintsUsed, &t.CompletedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// MarkTask отмечает попытку по заданию и, если решено, фиксирует выполнение.
func (r *ProgressRepo) MarkTask(ctx context.Context, userID, lessonID, taskID string, solved, usedHint bool) error {
	hint := 0
	if usedHint {
		hint = 1
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO task_progress (user_id, lesson_id, task_id, attempts, hints_used, completed_at, updated_at)
		VALUES ($1, $2, $3, 1, $5, CASE WHEN $4 THEN now() END, now())
		ON CONFLICT (user_id, lesson_id, task_id) DO UPDATE
		   SET attempts = task_progress.attempts + 1,
		       hints_used = task_progress.hints_used + $5,
		       completed_at = CASE
		           WHEN $4 THEN COALESCE(task_progress.completed_at, now())
		           ELSE task_progress.completed_at
		       END,
		       updated_at = now()`, userID, lessonID, taskID, solved, hint)
	return err
}

// CompletedTaskIDs — какие задания урока уже решены.
func (r *ProgressRepo) CompletedTaskIDs(ctx context.Context, userID, lessonID string) (map[string]bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT task_id FROM task_progress
		 WHERE user_id = $1 AND lesson_id = $2 AND completed_at IS NOT NULL`, userID, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]bool, 8)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out[id] = true
	}
	return out, rows.Err()
}

// QuizCard — квиз в списке «Квизы» у студента.
type QuizCard struct {
	LessonID    string     `json:"lessonId"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	CourseID    string     `json:"courseId"`
	CourseSlug  string     `json:"courseSlug"`
	CourseTitle string     `json:"courseTitle"`
	ModuleTitle string     `json:"moduleTitle"`
	Position    int        `json:"position"`
	Questions   int        `json:"questions"`
	PassScore   float64    `json:"passScore"`
	DurationMin int        `json:"durationMin"`
	Status      string     `json:"status"`
	BestScore   *float64   `json:"bestScore"`
	Attempts    int        `json:"attempts"`
	LastTriedAt *time.Time `json:"lastTriedAt"`
}

// Quizzes собирает все квизы курсов, на которые записан студент,
// вместе с его результатами. Нужен для отдельной страницы «Квизы».
func (r *ProgressRepo) Quizzes(ctx context.Context, userID string) ([]QuizCard, error) {
	rows, err := r.db.Query(ctx, `
		SELECT l.id, l.title, l.summary, l.content, l.duration_min,
		       c.id, c.slug, c.title, m.title,
		       row_number() OVER (ORDER BY m.position, m.created_at, l.position, l.created_at),
		       COALESCE(p.status, 'not_started'), p.best_score,
		       COALESCE(p.attempts, 0),
		       (SELECT max(a.created_at) FROM lesson_attempts a
		         WHERE a.user_id = $1 AND a.lesson_id = l.id)
		  FROM lessons l
		  JOIN modules m ON m.id = l.module_id
		  JOIN courses c ON c.id = m.course_id
		  JOIN enrollments e ON e.course_id = c.id AND e.user_id = $1
		  LEFT JOIN lesson_progress p ON p.lesson_id = l.id AND p.user_id = $1
		 WHERE l.kind = 'quiz' AND c.status = 'published'
		 ORDER BY c.position, m.position, m.created_at, l.position, l.created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]QuizCard, 0, 16)
	for rows.Next() {
		var card QuizCard
		var content json.RawMessage

		if err := rows.Scan(&card.LessonID, &card.Title, &card.Summary, &content, &card.DurationMin,
			&card.CourseID, &card.CourseSlug, &card.CourseTitle, &card.ModuleTitle, &card.Position,
			&card.Status, &card.BestScore, &card.Attempts, &card.LastTriedAt); err != nil {
			return nil, err
		}

		// Количество вопросов и порог читаем прямо из содержимого урока.
		var quiz struct {
			PassScore float64           `json:"passScore"`
			Questions []json.RawMessage `json:"questions"`
		}
		if err := json.Unmarshal(content, &quiz); err == nil {
			card.Questions = len(quiz.Questions)
			card.PassScore = quiz.PassScore
		}
		if card.PassScore <= 0 {
			card.PassScore = 70
		}

		out = append(out, card)
	}
	return out, rows.Err()
}
