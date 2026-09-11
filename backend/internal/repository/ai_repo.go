package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AISettingsRepo хранит настройки ИИ-помощника: вкл/выкл, ключ Gemini и модель.
// Ключ живёт только на сервере и наружу не отдаётся.
type AISettingsRepo struct{ db *pgxpool.Pool }

func NewAISettingsRepo(db *pgxpool.Pool) *AISettingsRepo { return &AISettingsRepo{db: db} }

// AISettings — сохранённые администратором настройки. APIKey/Model пустые,
// если админ их не задавал (тогда берётся значение из окружения).
type AISettings struct {
	Enabled bool
	APIKey  string
	Model   string
}

// Get возвращает настройки. Если строки ещё нет — фича включена «из коробки»,
// ключ/модель пустые (сработает fallback на переменные окружения).
func (r *AISettingsRepo) Get(ctx context.Context) (AISettings, error) {
	var s AISettings
	var key, model *string
	err := r.db.QueryRow(ctx,
		`SELECT enabled, api_key, model FROM platform_ai WHERE id = TRUE`).
		Scan(&s.Enabled, &key, &model)
	if errors.Is(err, pgx.ErrNoRows) {
		return AISettings{Enabled: true}, nil
	}
	if err != nil {
		return AISettings{}, err
	}
	if key != nil {
		s.APIKey = strings.TrimSpace(*key)
	}
	if model != nil {
		s.Model = strings.TrimSpace(*model)
	}
	return s, nil
}

// Update сохраняет настройки. apiKey == nil — ключ не трогаем (админ не
// перевводил его); apiKey != nil — ставим новое значение (пустая строка стирает).
func (r *AISettingsRepo) Update(ctx context.Context, enabled bool, model string, apiKey *string, updatedBy string) error {
	cur, err := r.Get(ctx)
	if err != nil {
		return err
	}
	key := cur.APIKey
	if apiKey != nil {
		key = strings.TrimSpace(*apiKey)
	}

	var by *string
	if updatedBy != "" {
		by = &updatedBy
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO platform_ai (id, enabled, api_key, model, updated_by, updated_at)
		VALUES (TRUE, $1, NULLIF($2, ''), NULLIF($3, ''), $4, now())
		ON CONFLICT (id) DO UPDATE
		   SET enabled    = EXCLUDED.enabled,
		       api_key    = EXCLUDED.api_key,
		       model      = EXCLUDED.model,
		       updated_by = EXCLUDED.updated_by,
		       updated_at = now()`,
		enabled, key, strings.TrimSpace(model), by)
	return err
}
