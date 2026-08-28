package repository

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Certificate struct {
	ID               string     `json:"id"`
	Serial           string     `json:"serial"`
	UserID           string     `json:"userId"`
	CourseID         string     `json:"courseId"`
	HolderName       string     `json:"holderName"`
	CourseTitle      string     `json:"courseTitle"`
	Score            float64    `json:"score"`
	LessonsTotal     int        `json:"lessonsTotal"`
	LessonsCompleted int        `json:"lessonsCompleted"`
	RevokedAt        *time.Time `json:"revokedAt"`
	IssuedAt         time.Time  `json:"issuedAt"`

	// Заполняется только при выдаче — нужен для письма, в базе не хранится.
	HolderEmail string `json:"-"`
}

type CertificateRepo struct{ db *pgxpool.Pool }

func NewCertificateRepo(db *pgxpool.Pool) *CertificateRepo { return &CertificateRepo{db: db} }

const certColumns = `id, serial, user_id, course_id, holder_name, course_title, score,
	lessons_total, lessons_completed, revoked_at, issued_at`

func scanCertificate(row pgx.Row) (*Certificate, error) {
	var c Certificate
	err := row.Scan(&c.ID, &c.Serial, &c.UserID, &c.CourseID, &c.HolderName, &c.CourseTitle,
		&c.Score, &c.LessonsTotal, &c.LessonsCompleted, &c.RevokedAt, &c.IssuedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CourseCompletion — сводка по прохождению курса: нужна для решения о выдаче.
type CourseCompletion struct {
	CourseTitle string
	HolderName  string
	HolderEmail string
	Total       int
	Completed   int
	AvgScore    float64
}

// Completion считает, сколько уроков курса студент прошёл и с каким средним баллом.
func (r *CertificateRepo) Completion(ctx context.Context, userID, courseID string) (*CourseCompletion, error) {
	var out CourseCompletion

	err := r.db.QueryRow(ctx, `
		SELECT c.title,
		       COALESCE(NULLIF(u.full_name, ''), u.email),
		       u.email,
		       (SELECT count(*) FROM lessons l
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id),
		       (SELECT count(*) FROM lesson_progress p
		          JOIN lessons l ON l.id = p.lesson_id
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id AND p.user_id = u.id AND p.status = 'completed'),
		       COALESCE((SELECT avg(COALESCE(p.best_score, 100)) FROM lesson_progress p
		          JOIN lessons l ON l.id = p.lesson_id
		          JOIN modules m ON m.id = l.module_id
		         WHERE m.course_id = c.id AND p.user_id = u.id AND p.status = 'completed'), 0)
		  FROM courses c, users u
		 WHERE c.id = $1 AND u.id = $2`, courseID, userID).
		Scan(&out.CourseTitle, &out.HolderName, &out.HolderEmail, &out.Total, &out.Completed, &out.AvgScore)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// Issue выдаёт сертификат. Повторная выдача возвращает уже существующий.
func (r *CertificateRepo) Issue(ctx context.Context, userID, courseID string, c CourseCompletion) (*Certificate, error) {
	serial, err := newSerial()
	if err != nil {
		return nil, err
	}

	cert, err := scanCertificate(r.db.QueryRow(ctx, `
		INSERT INTO certificates (serial, user_id, course_id, holder_name, course_title,
		                          score, lessons_total, lessons_completed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id, course_id) DO UPDATE
		   SET score = EXCLUDED.score,
		       lessons_completed = EXCLUDED.lessons_completed,
		       lessons_total = EXCLUDED.lessons_total,
		       holder_name = EXCLUDED.holder_name
		RETURNING `+certColumns,
		serial, userID, courseID, c.HolderName, c.CourseTitle,
		c.AvgScore, c.Total, c.Completed))
	if cert != nil {
		cert.HolderEmail = c.HolderEmail
	}
	return cert, err
}

func (r *CertificateRepo) GetBySerial(ctx context.Context, serial string) (*Certificate, error) {
	return scanCertificate(r.db.QueryRow(ctx,
		`SELECT `+certColumns+` FROM certificates WHERE upper(serial) = upper($1)`, serial))
}

func (r *CertificateRepo) ForUser(ctx context.Context, userID string) ([]Certificate, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+certColumns+` FROM certificates WHERE user_id = $1 ORDER BY issued_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Certificate, 0, 4)
	for rows.Next() {
		cert, err := scanCertificate(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *cert)
	}
	return out, rows.Err()
}

func (r *CertificateRepo) List(ctx context.Context, limit int) ([]Certificate, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := r.db.Query(ctx,
		`SELECT `+certColumns+` FROM certificates ORDER BY issued_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Certificate, 0, limit)
	for rows.Next() {
		cert, err := scanCertificate(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *cert)
	}
	return out, rows.Err()
}

func (r *CertificateRepo) Revoke(ctx context.Context, id string, revoked bool) error {
	var tag string
	if revoked {
		tag = `UPDATE certificates SET revoked_at = now() WHERE id = $1`
	} else {
		tag = `UPDATE certificates SET revoked_at = NULL WHERE id = $1`
	}
	result, err := r.db.Exec(ctx, tag, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

const serialAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// newSerial формирует читаемый номер вида DP-2026-K3M9QF.
func newSerial() (string, error) {
	buf := make([]byte, 6)
	max := big.NewInt(int64(len(serialAlphabet)))
	for i := range buf {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		buf[i] = serialAlphabet[n.Int64()]
	}
	return fmt.Sprintf("DP-%d-%s", time.Now().Year(), string(buf)), nil
}
