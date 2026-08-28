package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MaxQuoteLen ограничивает длину сохранённой цитаты.
const MaxQuoteLen = 2000

// Note — заметка вместе с контекстом: из какого курса, модуля и урока она сделана.
type Note struct {
	ID          string    `json:"id"`
	LessonID    string    `json:"lessonId"`
	LessonTitle string    `json:"lessonTitle"`
	LessonKind  string    `json:"lessonKind"`
	ModuleTitle string    `json:"moduleTitle"`
	CourseSlug  string    `json:"courseSlug"`
	CourseTitle string    `json:"courseTitle"`
	Quote       string    `json:"quote"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type NoteRepo struct{ db *pgxpool.Pool }

func NewNoteRepo(db *pgxpool.Pool) *NoteRepo { return &NoteRepo{db: db} }

const noteSelect = `
	SELECT n.id, n.lesson_id, l.title, l.kind, m.title, c.slug, c.title,
	       n.quote, n.body, n.created_at, n.updated_at
	  FROM notes n
	  JOIN lessons l ON l.id = n.lesson_id
	  JOIN modules m ON m.id = l.module_id
	  JOIN courses c ON c.id = m.course_id`

func scanNote(row pgx.Row) (*Note, error) {
	var n Note
	err := row.Scan(&n.ID, &n.LessonID, &n.LessonTitle, &n.LessonKind, &n.ModuleTitle,
		&n.CourseSlug, &n.CourseTitle, &n.Quote, &n.Body, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// Create сохраняет заметку и возвращает её с контекстом урока.
func (r *NoteRepo) Create(ctx context.Context, userID, lessonID, quote, body string) (*Note, error) {
	quote = strings.TrimSpace(quote)
	if quote == "" {
		return nil, errors.New("пустая цитата")
	}
	if len(quote) > MaxQuoteLen {
		quote = quote[:MaxQuoteLen]
	}

	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO notes (user_id, lesson_id, quote, body)
		VALUES ($1, $2, $3, $4)
		RETURNING id`, userID, lessonID, quote, strings.TrimSpace(body)).Scan(&id)
	if err != nil {
		return nil, err
	}
	return scanNote(r.db.QueryRow(ctx, noteSelect+` WHERE n.id = $1`, id))
}

// ForUser — все заметки студента, свежие сверху.
func (r *NoteRepo) ForUser(ctx context.Context, userID string) ([]Note, error) {
	rows, err := r.db.Query(ctx, noteSelect+`
		 WHERE n.user_id = $1
		 ORDER BY n.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Note, 0, 16)
	for rows.Next() {
		note, err := scanNote(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *note)
	}
	return out, rows.Err()
}

// UpdateBody меняет комментарий к заметке. Чужую заметку тронуть нельзя.
func (r *NoteRepo) UpdateBody(ctx context.Context, userID, noteID, body string) (*Note, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE notes SET body = $3, updated_at = now()
		 WHERE id = $1 AND user_id = $2`, noteID, userID, strings.TrimSpace(body))
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}
	return scanNote(r.db.QueryRow(ctx, noteSelect+` WHERE n.id = $1`, noteID))
}

func (r *NoteRepo) Delete(ctx context.Context, userID, noteID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM notes WHERE id = $1 AND user_id = $2`, noteID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
