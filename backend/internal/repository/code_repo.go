package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const MaxCodeAttempts = 5

type VerificationCode struct {
	ID        string
	Email     string
	Purpose   string
	CodeHash  string
	Attempts  int
	ExpiresAt time.Time
}

type CodeRepo struct{ db *pgxpool.Pool }

func NewCodeRepo(db *pgxpool.Pool) *CodeRepo { return &CodeRepo{db: db} }

func (r *CodeRepo) Create(ctx context.Context, email, purpose, codeHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO verification_codes (email, purpose, code_hash, expires_at)
		VALUES ($1, $2, $3, $4)`, email, purpose, codeHash, expiresAt)
	return err
}

// GetActive возвращает последний неиспользованный код для почты и цели.
func (r *CodeRepo) GetActive(ctx context.Context, email, purpose string) (*VerificationCode, error) {
	var c VerificationCode
	err := r.db.QueryRow(ctx, `
		SELECT id, email, purpose, code_hash, attempts, expires_at
		  FROM verification_codes
		 WHERE lower(email) = lower($1) AND purpose = $2
		       AND consumed_at IS NULL AND expires_at > now()
		 ORDER BY created_at DESC
		 LIMIT 1`, email, purpose).Scan(&c.ID, &c.Email, &c.Purpose, &c.CodeHash, &c.Attempts, &c.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CodeRepo) IncAttempts(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE verification_codes SET attempts = attempts + 1 WHERE id = $1`, id)
	return err
}

func (r *CodeRepo) Consume(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE verification_codes SET consumed_at = now() WHERE id = $1`, id)
	return err
}

// SentRecently защищает от спама: сколько кодов отправлено за интервал.
func (r *CodeRepo) SentRecently(ctx context.Context, email, purpose string, within time.Duration) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `
		SELECT count(*)
		  FROM verification_codes
		 WHERE lower(email) = lower($1) AND purpose = $2
		       AND created_at > now() - $3::interval`,
		email, purpose, within.String()).Scan(&n)
	return n, err
}

func (r *CodeRepo) DeleteExpired(ctx context.Context) (int64, error) {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM verification_codes WHERE expires_at < now() - interval '1 day'`)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}
