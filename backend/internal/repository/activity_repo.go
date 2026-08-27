package repository

import (
	"context"
	"encoding/json"
	"time"

	"platforma/backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepo struct{ db *pgxpool.Pool }

func NewActivityRepo(db *pgxpool.Pool) *ActivityRepo { return &ActivityRepo{db: db} }

// Touch отмечает, что пользователь заходил сегодня, и добавляет проведённое время.
func (r *ActivityRepo) Touch(ctx context.Context, userID string, seconds int) error {
	if seconds < 0 {
		seconds = 0
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO activity_days (user_id, day, visits, seconds_spent)
		VALUES ($1, CURRENT_DATE, 1, $2)
		ON CONFLICT (user_id, day) DO UPDATE
		   SET visits = activity_days.visits + 1,
		       seconds_spent = activity_days.seconds_spent + EXCLUDED.seconds_spent,
		       last_seen_at = now()`, userID, seconds)
	return err
}

// Days возвращает дни посещения за последние N дней.
func (r *ActivityRepo) Days(ctx context.Context, userID string, days int) ([]domain.ActivityDay, error) {
	if days <= 0 {
		days = 30
	}
	rows, err := r.db.Query(ctx, `
		SELECT to_char(day, 'YYYY-MM-DD'), visits, seconds_spent
		  FROM activity_days
		 WHERE user_id = $1 AND day > CURRENT_DATE - $2::int
		 ORDER BY day`, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.ActivityDay, 0, days)
	for rows.Next() {
		var d domain.ActivityDay
		if err := rows.Scan(&d.Day, &d.Visits, &d.SecondsSpent); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// Streak — сколько дней подряд (включая сегодня или вчера) студент заходил.
func (r *ActivityRepo) Streak(ctx context.Context, userID string) (int, error) {
	rows, err := r.db.Query(ctx, `
		SELECT day FROM activity_days WHERE user_id = $1 ORDER BY day DESC LIMIT 400`, userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	streak := 0
	expected := time.Now().UTC().Truncate(24 * time.Hour)
	first := true

	for rows.Next() {
		var day time.Time
		if err := rows.Scan(&day); err != nil {
			return 0, err
		}
		day = day.UTC().Truncate(24 * time.Hour)

		if first {
			// Серия не рвётся, если сегодня студент ещё не заходил.
			if day.Equal(expected) || day.Equal(expected.AddDate(0, 0, -1)) {
				expected = day
			} else {
				return 0, nil
			}
			first = false
		}
		if !day.Equal(expected) {
			break
		}
		streak++
		expected = expected.AddDate(0, 0, -1)
	}
	return streak, rows.Err()
}

type AuditRepo struct{ db *pgxpool.Pool }

func NewAuditRepo(db *pgxpool.Pool) *AuditRepo { return &AuditRepo{db: db} }

// Log пишет действие администратора. Ошибки не критичны для запроса.
func (r *AuditRepo) Log(ctx context.Context, actorID, action, entity, entityID string, payload any) {
	raw := []byte("{}")
	if payload != nil {
		if b, err := json.Marshal(payload); err == nil {
			raw = b
		}
	}
	var actor *string
	if actorID != "" {
		actor = &actorID
	}
	_, _ = r.db.Exec(ctx, `
		INSERT INTO audit_log (actor_id, action, entity, entity_id, payload)
		VALUES ($1, $2, $3, $4, $5)`, actor, action, entity, entityID, raw)
}

type AuditEntry struct {
	ID        int64           `json:"id"`
	ActorID   *string         `json:"actorId"`
	ActorName string          `json:"actorName"`
	Action    string          `json:"action"`
	Entity    string          `json:"entity"`
	EntityID  string          `json:"entityId"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"createdAt"`
}

func (r *AuditRepo) List(ctx context.Context, limit int) ([]AuditEntry, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.actor_id, COALESCE(u.full_name, u.email, ''), a.action, a.entity,
		       a.entity_id, a.payload, a.created_at
		  FROM audit_log a
		  LEFT JOIN users u ON u.id = a.actor_id
		 ORDER BY a.created_at DESC
		 LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]AuditEntry, 0, limit)
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.ActorID, &e.ActorName, &e.Action, &e.Entity,
			&e.EntityID, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
