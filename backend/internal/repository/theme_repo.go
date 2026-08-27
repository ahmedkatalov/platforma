package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThemeRepo struct{ db *pgxpool.Pool }

func NewThemeRepo(db *pgxpool.Pool) *ThemeRepo { return &ThemeRepo{db: db} }

// GetPlatform — общее оформление, заданное администратором (nil, если не задано).
func (r *ThemeRepo) GetPlatform(ctx context.Context) (json.RawMessage, error) {
	var raw json.RawMessage
	err := r.db.QueryRow(ctx, `SELECT settings FROM platform_theme WHERE id = TRUE`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (r *ThemeRepo) SetPlatform(ctx context.Context, settings json.RawMessage, updatedBy string) error {
	var by *string
	if updatedBy != "" {
		by = &updatedBy
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO platform_theme (id, settings, updated_by, updated_at)
		VALUES (TRUE, $1, $2, now())
		ON CONFLICT (id) DO UPDATE
		   SET settings = EXCLUDED.settings,
		       updated_by = EXCLUDED.updated_by,
		       updated_at = now()`, settings, by)
	return err
}

// GetUser — персональные настройки интерфейса (nil, если пользователь ничего не менял).
func (r *ThemeRepo) GetUser(ctx context.Context, userID string) (json.RawMessage, error) {
	var raw json.RawMessage
	err := r.db.QueryRow(ctx,
		`SELECT theme FROM user_preferences WHERE user_id = $1`, userID).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (r *ThemeRepo) SetUser(ctx context.Context, userID string, settings json.RawMessage) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO user_preferences (user_id, theme, updated_at)
		VALUES ($1, $2, now())
		ON CONFLICT (user_id) DO UPDATE
		   SET theme = EXCLUDED.theme, updated_at = now()`, userID, settings)
	return err
}

func (r *ThemeRepo) ResetUser(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM user_preferences WHERE user_id = $1`, userID)
	return err
}
