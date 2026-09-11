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

// AISettings — сохранённые администратором настройки. APIKey/Model/Proxy пустые,
// если админ их не задавал (тогда берётся значение из окружения).
type AISettings struct {
	Enabled bool
	APIKey  string
	Model   string
	Proxy   string
}

// Get возвращает настройки. Если строки ещё нет — фича включена «из коробки»,
// ключ/модель/прокси пустые (сработает fallback на переменные окружения).
func (r *AISettingsRepo) Get(ctx context.Context) (AISettings, error) {
	var s AISettings
	var key, model, proxy *string
	err := r.db.QueryRow(ctx,
		`SELECT enabled, api_key, model, proxy FROM platform_ai WHERE id = TRUE`).
		Scan(&s.Enabled, &key, &model, &proxy)
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
	if proxy != nil {
		s.Proxy = strings.TrimSpace(*proxy)
	}
	return s, nil
}

// AIUpdate — что изменить. Для APIKey/Proxy: nil — не трогать (админ не
// перевводил), не-nil — задать новое значение (пустая строка стирает).
type AIUpdate struct {
	Enabled bool
	Model   string
	APIKey  *string
	Proxy   *string
}

func (r *AISettingsRepo) Update(ctx context.Context, u AIUpdate, updatedBy string) error {
	cur, err := r.Get(ctx)
	if err != nil {
		return err
	}
	key := cur.APIKey
	if u.APIKey != nil {
		key = strings.TrimSpace(*u.APIKey)
	}
	proxy := cur.Proxy
	if u.Proxy != nil {
		proxy = strings.TrimSpace(*u.Proxy)
	}

	var by *string
	if updatedBy != "" {
		by = &updatedBy
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO platform_ai (id, enabled, api_key, model, proxy, updated_by, updated_at)
		VALUES (TRUE, $1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''), $5, now())
		ON CONFLICT (id) DO UPDATE
		   SET enabled    = EXCLUDED.enabled,
		       api_key    = EXCLUDED.api_key,
		       model      = EXCLUDED.model,
		       proxy      = EXCLUDED.proxy,
		       updated_by = EXCLUDED.updated_by,
		       updated_at = now()`,
		u.Enabled, key, strings.TrimSpace(u.Model), proxy, by)
	return err
}
