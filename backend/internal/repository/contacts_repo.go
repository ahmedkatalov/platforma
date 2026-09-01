package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ContactsRepo хранит контакты для связи (Telegram/WhatsApp), заданные админом.
type ContactsRepo struct{ db *pgxpool.Pool }

func NewContactsRepo(db *pgxpool.Pool) *ContactsRepo { return &ContactsRepo{db: db} }

// Get — настройки контактов (nil, если не заданы).
func (r *ContactsRepo) Get(ctx context.Context) (json.RawMessage, error) {
	var raw json.RawMessage
	err := r.db.QueryRow(ctx, `SELECT settings FROM platform_contacts WHERE id = TRUE`).Scan(&raw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (r *ContactsRepo) Set(ctx context.Context, settings json.RawMessage, updatedBy string) error {
	var by *string
	if updatedBy != "" {
		by = &updatedBy
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO platform_contacts (id, settings, updated_by, updated_at)
		VALUES (TRUE, $1, $2, now())
		ON CONFLICT (id) DO UPDATE
		   SET settings = EXCLUDED.settings,
		       updated_by = EXCLUDED.updated_by,
		       updated_at = now()`, settings, by)
	return err
}
