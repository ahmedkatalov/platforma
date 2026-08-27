package middleware

import (
	"context"
	"net/http"
	"strings"

	"platforma/backend/internal/auth"
	"platforma/backend/internal/domain"
)

type ctxKey string

const (
	CtxUserID ctxKey = "user_id"
	CtxRole   ctxKey = "role"
	CtxEmail  ctxKey = "email"
)

// Auth проверяет access-токен в заголовке Authorization: Bearer <token>.
func Auth(tm *auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(strings.ToLower(header), "bearer ") {
				unauthorized(w)
				return
			}

			claims, err := tm.Parse(strings.TrimSpace(header[7:]))
			if err != nil {
				unauthorized(w)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin пропускает только администраторов.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Role(r.Context()) != domain.RoleAdmin {
			http.Error(w, `{"message":"Доступ только для администратора"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserID(ctx context.Context) string {
	v, _ := ctx.Value(CtxUserID).(string)
	return v
}

func Role(ctx context.Context) string {
	v, _ := ctx.Value(CtxRole).(string)
	return v
}

func Email(ctx context.Context) string {
	v, _ := ctx.Value(CtxEmail).(string)
	return v
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"message":"Требуется авторизация"}`))
}
