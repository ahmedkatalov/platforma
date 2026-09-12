package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Bookmark — сохранённый урок (тема) или модуль (глава) с контекстом для показа.
type Bookmark struct {
	Kind        string    `json:"kind"` // lesson | module
	RefID       string    `json:"refId"`
	Title       string    `json:"title"`
	LessonKind  string    `json:"lessonKind,omitempty"`  // только для урока
	ModuleTitle string    `json:"moduleTitle,omitempty"` // только для урока
	CourseSlug  string    `json:"courseSlug"`
	CourseTitle string    `json:"courseTitle"`
	CreatedAt   time.Time `json:"createdAt"`
}

type BookmarkRepo struct{ db *pgxpool.Pool }

func NewBookmarkRepo(db *pgxpool.Pool) *BookmarkRepo { return &BookmarkRepo{db: db} }

var errBadBookmarkKind = errors.New("некорректный тип закладки")

func validKind(kind string) bool { return kind == "lesson" || kind == "module" }

// exists проверяет, что цель закладки реально существует.
func (r *BookmarkRepo) exists(ctx context.Context, kind, refID string) (bool, error) {
	table := "lessons"
	if kind == "module" {
		table = "modules"
	}
	var ok bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM `+table+` WHERE id = $1)`, refID).Scan(&ok)
	return ok, err
}

// Add сохраняет закладку (повторное сохранение не ошибка).
func (r *BookmarkRepo) Add(ctx context.Context, userID, kind, refID string) error {
	if !validKind(kind) {
		return errBadBookmarkKind
	}
	ok, err := r.exists(ctx, kind, refID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO bookmarks (user_id, kind, ref_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, kind, ref_id) DO NOTHING`, userID, kind, refID)
	return err
}

// Remove убирает закладку.
func (r *BookmarkRepo) Remove(ctx context.Context, userID, kind, refID string) error {
	if !validKind(kind) {
		return errBadBookmarkKind
	}
	_, err := r.db.Exec(ctx,
		`DELETE FROM bookmarks WHERE user_id = $1 AND kind = $2 AND ref_id = $3`,
		userID, kind, refID)
	return err
}

// List — все закладки студента (уроки и модули), свежие сверху. Осиротевшие
// (цель удалили) join отсекает автоматически.
func (r *BookmarkRepo) List(ctx context.Context, userID string) ([]Bookmark, error) {
	rows, err := r.db.Query(ctx, `
		SELECT b.kind, b.ref_id, l.title, l.kind, m.title, c.slug, c.title, b.created_at
		  FROM bookmarks b
		  JOIN lessons l ON b.kind = 'lesson' AND l.id = b.ref_id
		  JOIN modules m ON m.id = l.module_id
		  JOIN courses c ON c.id = m.course_id
		 WHERE b.user_id = $1
		UNION ALL
		SELECT b.kind, b.ref_id, m.title, '', '', c.slug, c.title, b.created_at
		  FROM bookmarks b
		  JOIN modules m ON b.kind = 'module' AND m.id = b.ref_id
		  JOIN courses c ON c.id = m.course_id
		 WHERE b.user_id = $1
		 ORDER BY 8 DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Bookmark, 0, 16)
	for rows.Next() {
		var b Bookmark
		if err := rows.Scan(&b.Kind, &b.RefID, &b.Title, &b.LessonKind, &b.ModuleTitle,
			&b.CourseSlug, &b.CourseTitle, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}
