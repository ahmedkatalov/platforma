package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Asset struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Original  string    `json:"original"`
	Mime      string    `json:"mime"`
	SizeBytes int64     `json:"sizeBytes"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}

type AssetRepo struct{ db *pgxpool.Pool }

func NewAssetRepo(db *pgxpool.Pool) *AssetRepo { return &AssetRepo{db: db} }

type AssetInput struct {
	Filename   string
	Original   string
	Mime       string
	SizeBytes  int64
	UploadedBy string
}

func (r *AssetRepo) Create(ctx context.Context, in AssetInput) (*Asset, error) {
	var by *string
	if in.UploadedBy != "" {
		by = &in.UploadedBy
	}

	var a Asset
	err := r.db.QueryRow(ctx, `
		INSERT INTO assets (filename, original, mime, size_bytes, uploaded_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, filename, original, mime, size_bytes, created_at`,
		in.Filename, in.Original, in.Mime, in.SizeBytes, by).
		Scan(&a.ID, &a.Filename, &a.Original, &a.Mime, &a.SizeBytes, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	a.URL = "/uploads/" + a.Filename
	return &a, nil
}

func (r *AssetRepo) List(ctx context.Context, limit int) ([]Asset, error) {
	if limit <= 0 || limit > 200 {
		limit = 60
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, filename, original, mime, size_bytes, created_at
		  FROM assets ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Asset, 0, limit)
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.Filename, &a.Original, &a.Mime, &a.SizeBytes, &a.CreatedAt); err != nil {
			return nil, err
		}
		a.URL = "/uploads/" + a.Filename
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AssetRepo) Delete(ctx context.Context, id string) (*Asset, error) {
	var a Asset
	err := r.db.QueryRow(ctx, `
		DELETE FROM assets WHERE id = $1
		RETURNING id, filename, original, mime, size_bytes, created_at`, id).
		Scan(&a.ID, &a.Filename, &a.Original, &a.Mime, &a.SizeBytes, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	a.URL = "/uploads/" + a.Filename
	return &a, nil
}
