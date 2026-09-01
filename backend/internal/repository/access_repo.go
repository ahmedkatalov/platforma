package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AccessRepo — заявки студентов на открытие следующей главы.
type AccessRepo struct{ db *pgxpool.Pool }

func NewAccessRepo(db *pgxpool.Pool) *AccessRepo { return &AccessRepo{db: db} }

// AccessRequest — заявка в списке администратора (с данными студента и главы).
type AccessRequest struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	UserName    string    `json:"userName"`
	UserEmail   string    `json:"userEmail"`
	ModuleID    string    `json:"moduleId"`
	ModuleTitle string    `json:"moduleTitle"`
	ChapterNo   int       `json:"chapterNo"` // порядковый номер главы (1-based), как видит студент
	CourseID    string    `json:"courseId"`
	CourseTitle string    `json:"courseTitle"`
	CourseSlug  string    `json:"courseSlug"`
	Status      string    `json:"status"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Create создаёт заявку. Повторная висящая заявка на ту же главу игнорируется.
func (r *AccessRepo) Create(ctx context.Context, userID, moduleID, note string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO access_requests (user_id, module_id, note)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, module_id) WHERE status = 'pending' DO NOTHING`,
		userID, moduleID, note)
	return err
}

// List возвращает заявки, при status != "" — только с этим статусом.
func (r *AccessRepo) List(ctx context.Context, status string, limit int) ([]AccessRequest, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := r.db.Query(ctx, `
		SELECT ar.id, ar.user_id, COALESCE(u.full_name, ''), u.email,
		       ar.module_id, m.title,
		       (SELECT count(*) FROM modules m2
		         WHERE m2.course_id = m.course_id
		           AND (m2.position, m2.created_at) <= (m.position, m.created_at)) AS module_rank,
		       c.id, c.title, c.slug,
		       ar.status, ar.note, ar.created_at
		  FROM access_requests ar
		  JOIN users u ON u.id = ar.user_id
		  JOIN modules m ON m.id = ar.module_id
		  JOIN courses c ON c.id = m.course_id
		 WHERE ($1 = '' OR ar.status = $1)
		 ORDER BY ar.created_at DESC
		 LIMIT $2`, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]AccessRequest, 0, limit)
	for rows.Next() {
		var a AccessRequest
		if err := rows.Scan(&a.ID, &a.UserID, &a.UserName, &a.UserEmail,
			&a.ModuleID, &a.ModuleTitle, &a.ChapterNo,
			&a.CourseID, &a.CourseTitle, &a.CourseSlug,
			&a.Status, &a.Note, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Get возвращает студента и главу заявки (для одобрения/отклонения).
func (r *AccessRepo) Get(ctx context.Context, id string) (userID, moduleID, status string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT user_id, module_id, status FROM access_requests WHERE id = $1`, id).
		Scan(&userID, &moduleID, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", ErrNotFound
	}
	return userID, moduleID, status, err
}

// Decide переводит висящую заявку в approved/rejected.
func (r *AccessRepo) Decide(ctx context.Context, id, status, decidedBy string) error {
	var by *string
	if decidedBy != "" {
		by = &decidedBy
	}
	tag, err := r.db.Exec(ctx, `
		UPDATE access_requests
		   SET status = $2, decided_at = now(), decided_by = $3
		 WHERE id = $1 AND status = 'pending'`, id, status, by)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// StatusMap — последний статус заявки по каждой главе курса (id главы -> статус).
// Нужен студенту, чтобы показать «заявка на рассмотрении» на закрытой главе.
func (r *AccessRepo) StatusMap(ctx context.Context, userID, courseID string) (map[string]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT ON (ar.module_id) ar.module_id, ar.status
		  FROM access_requests ar
		  JOIN modules m ON m.id = ar.module_id
		 WHERE ar.user_id = $1 AND m.course_id = $2
		 ORDER BY ar.module_id, ar.created_at DESC`, userID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string, 8)
	for rows.Next() {
		var id, status string
		if err := rows.Scan(&id, &status); err != nil {
			return nil, err
		}
		out[id] = status
	}
	return out, rows.Err()
}

// PendingCount — сколько заявок ждут решения (для бейджа в админке).
func (r *AccessRepo) PendingCount(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRow(ctx,
		`SELECT count(*) FROM access_requests WHERE status = 'pending'`).Scan(&n)
	return n, err
}

// PendingCounts — ожидающие заявки по главам и по курсам разом (для бейджа).
func (r *AccessRepo) PendingCounts(ctx context.Context) (chapters, courses int, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT
			(SELECT count(*) FROM access_requests WHERE status = 'pending'),
			(SELECT count(*) FROM course_requests WHERE status = 'pending')`).
		Scan(&chapters, &courses)
	return chapters, courses, err
}

// --- Заявки на доступ к КУРСУ (запись на курс) ---

// CourseRequest — заявка на курс в списке администратора.
type CourseRequest struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	UserName    string    `json:"userName"`
	UserEmail   string    `json:"userEmail"`
	CourseID    string    `json:"courseId"`
	CourseTitle string    `json:"courseTitle"`
	CourseSlug  string    `json:"courseSlug"`
	Status      string    `json:"status"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (r *AccessRepo) CreateCourseRequest(ctx context.Context, userID, courseID, note string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO course_requests (user_id, course_id, note)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, course_id) WHERE status = 'pending' DO NOTHING`,
		userID, courseID, note)
	return err
}

func (r *AccessRepo) ListCourseRequests(ctx context.Context, status string, limit int) ([]CourseRequest, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := r.db.Query(ctx, `
		SELECT cr.id, cr.user_id, COALESCE(u.full_name, ''), u.email,
		       cr.course_id, c.title, c.slug,
		       cr.status, cr.note, cr.created_at
		  FROM course_requests cr
		  JOIN users u ON u.id = cr.user_id
		  JOIN courses c ON c.id = cr.course_id
		 WHERE ($1 = '' OR cr.status = $1)
		 ORDER BY cr.created_at DESC
		 LIMIT $2`, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]CourseRequest, 0, limit)
	for rows.Next() {
		var a CourseRequest
		if err := rows.Scan(&a.ID, &a.UserID, &a.UserName, &a.UserEmail,
			&a.CourseID, &a.CourseTitle, &a.CourseSlug,
			&a.Status, &a.Note, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AccessRepo) GetCourseRequest(ctx context.Context, id string) (userID, courseID, status string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT user_id, course_id, status FROM course_requests WHERE id = $1`, id).
		Scan(&userID, &courseID, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", "", ErrNotFound
	}
	return userID, courseID, status, err
}

func (r *AccessRepo) DecideCourseRequest(ctx context.Context, id, status, decidedBy string) error {
	var by *string
	if decidedBy != "" {
		by = &decidedBy
	}
	tag, err := r.db.Exec(ctx, `
		UPDATE course_requests
		   SET status = $2, decided_at = now(), decided_by = $3
		 WHERE id = $1 AND status = 'pending'`, id, status, by)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// CourseRequestStatusForUser — последний статус заявки по каждому курсу (id курса -> статус).
// Нужен витрине, чтобы показать «заявка на рассмотрении» на закрытом курсе.
func (r *AccessRepo) CourseRequestStatusForUser(ctx context.Context, userID string) (map[string]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT DISTINCT ON (course_id) course_id, status
		  FROM course_requests
		 WHERE user_id = $1
		 ORDER BY course_id, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string, 8)
	for rows.Next() {
		var id, status string
		if err := rows.Scan(&id, &status); err != nil {
			return nil, err
		}
		out[id] = status
	}
	return out, rows.Err()
}
