package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AISettingsRepo хранит флаг доступности ИИ-помощника (вкл/выкл админом).
type AISettingsRepo struct{ db *pgxpool.Pool }

func NewAISettingsRepo(db *pgxpool.Pool) *AISettingsRepo { return &AISettingsRepo{db: db} }

// Enabled — включён ли ИИ-помощник. Если строки нет — считаем включённым
// (фича «из коробки» работает, как только задан ключ); админ может выключить.
func (r *AISettingsRepo) Enabled(ctx context.Context) (bool, error) {
	var enabled bool
	err := r.db.QueryRow(ctx, `SELECT enabled FROM platform_ai WHERE id = TRUE`).Scan(&enabled)
	if errors.Is(err, pgx.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return enabled, nil
}

// SetEnabled — сохранить флаг (кто изменил — для аудита).
func (r *AISettingsRepo) SetEnabled(ctx context.Context, enabled bool, updatedBy string) error {
	var by *string
	if updatedBy != "" {
		by = &updatedBy
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO platform_ai (id, enabled, updated_by, updated_at)
		VALUES (TRUE, $1, $2, now())
		ON CONFLICT (id) DO UPDATE
		   SET enabled = EXCLUDED.enabled,
		       updated_by = EXCLUDED.updated_by,
		       updated_at = now()`, enabled, by)
	return err
}
