package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"platforma/backend/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound      = errors.New("не найдено")
	ErrEmailTaken    = errors.New("пользователь с такой почтой уже существует")
	ErrLastAdminGone = errors.New("нельзя удалить последнего администратора")
)

const userColumns = `id, email, full_name, password_hash, role, status,
	email_verified, avatar_url, last_login_at, created_at, updated_at`

type UserRepo struct{ db *pgxpool.Pool }

func NewUserRepo(db *pgxpool.Pool) *UserRepo { return &UserRepo{db: db} }

func scanUser(row pgx.Row) (*domain.User, error) {
	var u domain.User
	err := row.Scan(&u.ID, &u.Email, &u.FullName, &u.PasswordHash, &u.Role, &u.Status,
		&u.EmailVerified, &u.AvatarURL, &u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

type CreateUserInput struct {
	Email         string
	FullName      string
	PasswordHash  string
	Role          string
	Status        string
	EmailVerified bool
	CreatedBy     *string
}

func (r *UserRepo) Create(ctx context.Context, in CreateUserInput) (*domain.User, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO users (email, full_name, password_hash, role, status, email_verified, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING `+userColumns,
		strings.TrimSpace(in.Email), strings.TrimSpace(in.FullName), in.PasswordHash,
		in.Role, in.Status, in.EmailVerified, in.CreatedBy)

	u, err := scanUser(row)
	if err != nil && isUniqueViolation(err) {
		return nil, ErrEmailTaken
	}
	return u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id))
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return scanUser(r.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE email_lower = lower($1)`, strings.TrimSpace(email)))
}

func (r *UserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM users WHERE email_lower = lower($1))`,
		strings.TrimSpace(email)).Scan(&exists)
	return exists, err
}

type ListUsersFilter struct {
	Search string
	Role   string
	Status string
	Limit  int
	Offset int
}

// ListUsers возвращает страницу пользователей и общее количество по фильтру.
func (r *UserRepo) ListUsers(ctx context.Context, f ListUsersFilter) ([]domain.User, int, error) {
	where := []string{"1 = 1"}
	args := []any{}

	if s := strings.TrimSpace(f.Search); s != "" {
		args = append(args, "%"+strings.ToLower(s)+"%")
		where = append(where, fmt.Sprintf("(email_lower LIKE $%d OR lower(full_name) LIKE $%d)", len(args), len(args)))
	}
	if s := strings.TrimSpace(f.Role); s != "" {
		args = append(args, s)
		where = append(where, fmt.Sprintf("role = $%d", len(args)))
	}
	if s := strings.TrimSpace(f.Status); s != "" {
		args = append(args, s)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}

	clause := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, `SELECT count(*) FROM users WHERE `+clause, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if f.Limit <= 0 {
		f.Limit = 50
	}
	args = append(args, f.Limit, f.Offset)
	query := fmt.Sprintf(`SELECT %s FROM users WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		userColumns, clause, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]domain.User, 0, f.Limit)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, *u)
	}
	return users, total, rows.Err()
}

type UpdateUserInput struct {
	FullName *string
	Email    *string
	Role     *string
	Status   *string
}

func (r *UserRepo) Update(ctx context.Context, id string, in UpdateUserInput) (*domain.User, error) {
	sets := []string{"updated_at = now()"}
	args := []any{}

	if in.FullName != nil {
		args = append(args, strings.TrimSpace(*in.FullName))
		sets = append(sets, fmt.Sprintf("full_name = $%d", len(args)))
	}
	if in.Email != nil {
		args = append(args, strings.TrimSpace(*in.Email))
		sets = append(sets, fmt.Sprintf("email = $%d", len(args)))
	}
	if in.Role != nil {
		args = append(args, *in.Role)
		sets = append(sets, fmt.Sprintf("role = $%d", len(args)))
	}
	if in.Status != nil {
		args = append(args, *in.Status)
		sets = append(sets, fmt.Sprintf("status = $%d", len(args)))
	}

	args = append(args, id)
	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d RETURNING %s`,
		strings.Join(sets, ", "), len(args), userColumns)

	u, err := scanUser(r.db.QueryRow(ctx, query, args...))
	if err != nil && isUniqueViolation(err) {
		return nil, ErrEmailTaken
	}
	return u, err
}

func (r *UserRepo) SetPassword(ctx context.Context, id, passwordHash string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1`, id, passwordHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ActivateByEmail помечает почту подтверждённой и включает аккаунт.
func (r *UserRepo) ActivateByEmail(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	return scanUser(r.db.QueryRow(ctx, `
		UPDATE users
		   SET password_hash = $2,
		       email_verified = TRUE,
		       status = 'active',
		       updated_at = now()
		 WHERE email_lower = lower($1)
		RETURNING `+userColumns, strings.TrimSpace(email), passwordHash))
}

func (r *UserRepo) TouchLogin(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET last_login_at = now() WHERE id = $1`, id)
	return err
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserRepo) CountByRole(ctx context.Context, role string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT count(*) FROM users WHERE role = $1`, role).Scan(&n)
	return n, err
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "SQLSTATE 23505")
}
