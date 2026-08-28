package handler

import (
	"encoding/json"
	"net/http"

	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

// MeHandler — данные текущего пользователя: профиль, курсы, активность, тема.
type MeHandler struct {
	users    *repository.UserRepo
	courses  *repository.CourseRepo
	activity *repository.ActivityRepo
	stats    *repository.StatsRepo
	theme    *repository.ThemeRepo
	progress *repository.ProgressRepo
	certs    *repository.CertificateRepo
	auth     *AuthHandler
}

func NewMeHandler(
	users *repository.UserRepo,
	courses *repository.CourseRepo,
	activity *repository.ActivityRepo,
	stats *repository.StatsRepo,
	theme *repository.ThemeRepo,
	progress *repository.ProgressRepo,
	certs *repository.CertificateRepo,
	auth *AuthHandler,
) *MeHandler {
	return &MeHandler{users: users, courses: courses, activity: activity,
		stats: stats, theme: theme, progress: progress, certs: certs, auth: auth}
}

func (h *MeHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.profile)
	r.Patch("/", h.updateProfile)
	r.Post("/change-password", h.auth.ChangePassword)
	r.Get("/stats", h.myStats)
	r.Get("/attempts", h.myAttempts)
	r.Get("/certificates", h.myCertificates)
	r.Post("/activity", h.trackActivity)
	r.Get("/preferences", h.getPreferences)
	r.Put("/preferences", h.putPreferences)
	r.Delete("/preferences", h.resetPreferences)
	return r
}

func (h *MeHandler) profile(w http.ResponseWriter, r *http.Request) {
	user, err := h.users.GetByID(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	enrollments, err := h.courses.EnrollmentsForUser(r.Context(), user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курсы")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user":        user,
		"enrollments": enrollments,
	})
}

func (h *MeHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FullName *string `json:"fullName"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.users.Update(r.Context(), middleware.UserID(r.Context()),
		repository.UpdateUserInput{FullName: body.FullName})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *MeHandler) myStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserID(r.Context())
	days := queryInt(r, "days", 30, 7, 365)

	summary, err := h.stats.StudentSummary(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать статистику")
		return
	}
	activity, err := h.activity.Days(r.Context(), userID, days)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить активность")
		return
	}
	streak, err := h.activity.Streak(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось посчитать серию дней")
		return
	}
	quiz, err := h.progress.QuizStats(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать статистику квизов")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"summary":  summary,
		"activity": activity,
		"streak":   streak,
		"quiz":     quiz,
	})
}

// myAttempts — история попыток по квизам, терминалу и коду.
func (h *MeHandler) myAttempts(w http.ResponseWriter, r *http.Request) {
	attempts, err := h.progress.Attempts(r.Context(), middleware.UserID(r.Context()),
		queryInt(r, "limit", 50, 1, 200))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить историю")
		return
	}
	writeJSON(w, http.StatusOK, attempts)
}

// myCertificates — сертификаты, полученные студентом.
func (h *MeHandler) myCertificates(w http.ResponseWriter, r *http.Request) {
	certs, err := h.certs.ForUser(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить сертификаты")
		return
	}
	writeJSON(w, http.StatusOK, certs)
}

// trackActivity — фронтенд периодически шлёт проведённое время.
func (h *MeHandler) trackActivity(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Seconds int `json:"seconds"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if body.Seconds > 3600 {
		body.Seconds = 3600 // защита от накрутки одним запросом
	}
	if err := h.activity.Touch(r.Context(), middleware.UserID(r.Context()), body.Seconds); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось записать активность")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

func (h *MeHandler) getPreferences(w http.ResponseWriter, r *http.Request) {
	personal, err := h.theme.GetUser(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить настройки")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"theme": rawOrNil(personal)})
}

func (h *MeHandler) putPreferences(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Theme json.RawMessage `json:"theme"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(body.Theme) == 0 || !json.Valid(body.Theme) {
		writeError(w, http.StatusBadRequest, "Некорректные настройки оформления")
		return
	}
	if err := h.theme.SetUser(r.Context(), middleware.UserID(r.Context()), body.Theme); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить настройки")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"theme": body.Theme})
}

func (h *MeHandler) resetPreferences(w http.ResponseWriter, r *http.Request) {
	if err := h.theme.ResetUser(r.Context(), middleware.UserID(r.Context())); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сбросить настройки")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"theme": nil})
}

func rawOrNil(raw json.RawMessage) any {
	if len(raw) == 0 {
		return nil
	}
	return raw
}
