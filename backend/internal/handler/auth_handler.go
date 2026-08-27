package handler

import (
	"errors"
	"net"
	"net/http"
	"strings"

	"platforma/backend/internal/domain"
	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"
	"platforma/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler { return &AuthHandler{auth: a} }

// Routes — публичные маршруты /api/auth.
func (h *AuthHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/send-code", h.sendCode)
	r.Post("/register", h.register)
	r.Post("/login", h.login)
	r.Post("/refresh", h.refresh)
	r.Post("/logout", h.logout)
	r.Post("/reset-password", h.resetPassword)
	return r
}

func (h *AuthHandler) sendCode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email   string `json:"email"`
		Purpose string `json:"purpose"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if body.Purpose == "" {
		body.Purpose = domain.PurposeRegistration
	}

	if err := h.auth.SendCode(r.Context(), body.Email, body.Purpose); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, service.ErrTooManyRequests) {
			status = http.StatusTooManyRequests
		}
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Код отправлен на указанную почту",
	})
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var body service.RegisterInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.auth.Register(r.Context(), body, r.UserAgent(), clientIP(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, session)
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.auth.Login(r.Context(), body.Email, body.Password, r.UserAgent(), clientIP(r))
	if err != nil {
		status := http.StatusUnauthorized
		if errors.Is(err, service.ErrUserBlocked) {
			status = http.StatusForbidden
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	session, err := h.auth.Refresh(r.Context(), body.RefreshToken, r.UserAgent(), clientIP(r))
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (h *AuthHandler) logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.auth.Logout(r.Context(), body.RefreshToken); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось выйти")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Сессия завершена"})
}

func (h *AuthHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var body service.ResetPasswordInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.auth.ResetPassword(r.Context(), body); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusBadRequest, "Аккаунт не найден")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Пароль обновлён"})
}

// ChangePassword — защищённый маршрут, монтируется в /api/me.
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.UserID(r.Context())
	if err := h.auth.ChangePassword(r.Context(), userID, body.CurrentPassword, body.NewPassword); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Пароль изменён"})
}

func clientIP(r *http.Request) string {
	if fwd := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
