package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"platforma/backend/internal/domain"
	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"
	"platforma/backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// AdminHandler — всё, что доступно только администратору:
// аккаунты студентов, сводная статистика, журнал действий и оформление платформы.
type AdminHandler struct {
	users    *repository.UserRepo
	courses  *repository.CourseRepo
	stats    *repository.StatsRepo
	activity *repository.ActivityRepo
	audit    *repository.AuditRepo
	theme    *repository.ThemeRepo
	progress *repository.ProgressRepo
	userSvc  *service.UserService
}

func NewAdminHandler(
	users *repository.UserRepo,
	courses *repository.CourseRepo,
	stats *repository.StatsRepo,
	activity *repository.ActivityRepo,
	audit *repository.AuditRepo,
	theme *repository.ThemeRepo,
	progress *repository.ProgressRepo,
	userSvc *service.UserService,
) *AdminHandler {
	return &AdminHandler{users: users, courses: courses, stats: stats,
		activity: activity, audit: audit, theme: theme, progress: progress, userSvc: userSvc}
}

// Routes собирает /api/admin. Редактор курсов монтируется сюда же, чтобы не
// плодить пересекающиеся Mount-пути в главном роутере.
func (h *AdminHandler) Routes(courses, certificates, reports, uploads http.Handler) http.Handler {
	r := chi.NewRouter()

	r.Mount("/courses", courses)
	r.Mount("/certificates", certificates)
	r.Mount("/reports", reports)
	r.Mount("/uploads", uploads)

	r.Get("/overview", h.overview)
	r.Get("/audit", h.auditLog)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", h.listUsers)
		r.Post("/", h.createUser)
		r.Get("/{id}", h.getUser)
		r.Patch("/{id}", h.updateUser)
		r.Delete("/{id}", h.deleteUser)
		r.Post("/{id}/reset-password", h.resetUserPassword)
		r.Post("/{id}/enrollments", h.enroll)
		r.Patch("/{id}/enrollments/{courseId}", h.setDueDate)
		r.Delete("/{id}/enrollments/{courseId}", h.unenroll)
	})

	r.Get("/students-progress", h.studentsProgress)

	r.Route("/theme", func(r chi.Router) {
		r.Get("/", h.getTheme)
		r.Put("/", h.putTheme)
	})

	return r
}

func (h *AdminHandler) overview(w http.ResponseWriter, r *http.Request) {
	overview, err := h.stats.Overview(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать сводку")
		return
	}
	writeJSON(w, http.StatusOK, overview)
}

func (h *AdminHandler) auditLog(w http.ResponseWriter, r *http.Request) {
	entries, err := h.audit.List(r.Context(), queryInt(r, "limit", 50, 1, 200))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить журнал")
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (h *AdminHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 50, 1, 200)
	page := queryInt(r, "page", 1, 1, 10000)

	users, total, err := h.users.ListUsers(r.Context(), repository.ListUsersFilter{
		Search: r.URL.Query().Get("search"),
		Role:   r.URL.Query().Get("role"),
		Status: r.URL.Query().Get("status"),
		Limit:  limit,
		Offset: (page - 1) * limit,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить пользователей")
		return
	}

	// Дополняем каждого пользователя признаком «онлайн» и временем последней активности.
	ids := make([]string, len(users))
	for i := range users {
		ids[i] = users[i].ID
	}
	lastSeen, err := h.activity.LastSeenByUsers(r.Context(), ids)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить активность")
		return
	}

	now := time.Now()
	items := make([]userListItem, len(users))
	for i, u := range users {
		var ls *time.Time
		if t, ok := lastSeen[u.ID]; ok {
			ls = &t
		}
		items[i] = userListItem{User: u, LastSeenAt: ls, Online: repository.IsOnline(ls, now)}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// userListItem — пользователь в списке админки с признаком онлайна.
// Анонимное встраивание domain.User поднимает его поля в JSON на верхний уровень.
type userListItem struct {
	domain.User
	LastSeenAt *time.Time `json:"lastSeenAt"`
	Online     bool       `json:"online"`
}

func (h *AdminHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var body service.CreateStudentInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	created, err := h.userSvc.CreateStudent(r.Context(), middleware.UserID(r.Context()), body)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, repository.ErrEmailTaken) {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *AdminHandler) getUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	enrollments, err := h.courses.EnrollmentsForUser(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курсы")
		return
	}
	summary, err := h.stats.StudentSummary(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать статистику")
		return
	}
	activity, err := h.activity.Days(r.Context(), id, queryInt(r, "days", 30, 7, 365))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить активность")
		return
	}

	attempts, err := h.progress.Attempts(r.Context(), id, 30)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить попытки")
		return
	}
	quiz, err := h.progress.QuizStats(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать статистику квизов")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user":        user,
		"enrollments": enrollments,
		"summary":     summary,
		"activity":    activity,
		"attempts":    attempts,
		"quiz":        quiz,
	})
}

func (h *AdminHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FullName *string `json:"fullName"`
		Email    *string `json:"email"`
		Role     *string `json:"role"`
		Status   *string `json:"status"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if body.Role != nil && *body.Role != domain.RoleAdmin && *body.Role != domain.RoleStudent {
		writeError(w, http.StatusBadRequest, "Недопустимая роль")
		return
	}
	if body.Status != nil {
		switch *body.Status {
		case domain.StatusActive, domain.StatusBlocked, domain.StatusInvited:
		default:
			writeError(w, http.StatusBadRequest, "Недопустимый статус")
			return
		}
	}

	user, err := h.userSvc.Update(r.Context(), middleware.UserID(r.Context()), chi.URLParam(r, "id"),
		repository.UpdateUserInput{
			FullName: body.FullName,
			Email:    body.Email,
			Role:     body.Role,
			Status:   body.Status,
		})
	if err != nil {
		writeError(w, statusForRepoError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func (h *AdminHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	err := h.userSvc.Delete(r.Context(), middleware.UserID(r.Context()), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, statusForRepoError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Аккаунт удалён"})
}

func (h *AdminHandler) resetUserPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SendMail bool `json:"sendMail"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.userSvc.ResetPasswordByAdmin(r.Context(), middleware.UserID(r.Context()),
		chi.URLParam(r, "id"), body.SendMail)
	if err != nil {
		writeError(w, statusForRepoError(err), err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *AdminHandler) enroll(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CourseID string `json:"courseId"`
		DueDate  string `json:"dueDate"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	dueDate, err := parseDate(body.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректная дата дедлайна")
		return
	}

	userID := chi.URLParam(r, "id")
	if err := h.courses.Enroll(r.Context(), userID, body.CourseID, middleware.UserID(r.Context()), dueDate); err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось записать студента на курс")
		return
	}
	h.audit.Log(r.Context(), middleware.UserID(r.Context()), "enrollment.create", "user", userID,
		map[string]any{"courseId": body.CourseID})

	enrollments, err := h.courses.EnrollmentsForUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курсы")
		return
	}
	writeJSON(w, http.StatusOK, enrollments)
}

// PATCH /admin/users/{id}/enrollments/{courseId} — изменить срок прохождения.
func (h *AdminHandler) setDueDate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DueDate string `json:"dueDate"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	dueDate, err := parseDate(body.DueDate)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректная дата дедлайна")
		return
	}

	userID := chi.URLParam(r, "id")
	courseID := chi.URLParam(r, "courseId")

	if err := h.courses.SetDueDate(r.Context(), userID, courseID, dueDate); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось изменить срок")
		return
	}

	enrollments, err := h.courses.EnrollmentsForUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курсы")
		return
	}
	writeJSON(w, http.StatusOK, enrollments)
}

func (h *AdminHandler) unenroll(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	courseID := chi.URLParam(r, "courseId")

	if err := h.courses.Unenroll(r.Context(), userID, courseID); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось отписать студента")
		return
	}
	h.audit.Log(r.Context(), middleware.UserID(r.Context()), "enrollment.delete", "user", userID,
		map[string]any{"courseId": courseID})
	writeJSON(w, http.StatusOK, map[string]string{"message": "Студент отписан от курса"})
}

func (h *AdminHandler) studentsProgress(w http.ResponseWriter, r *http.Request) {
	items, err := h.stats.StudentsSummary(r.Context(), queryInt(r, "limit", 100, 1, 500))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать успеваемость")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *AdminHandler) getTheme(w http.ResponseWriter, r *http.Request) {
	settings, err := h.theme.GetPlatform(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить оформление")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"settings": rawOrNil(settings)})
}

func (h *AdminHandler) putTheme(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Settings json.RawMessage `json:"settings"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(body.Settings) == 0 || !json.Valid(body.Settings) {
		writeError(w, http.StatusBadRequest, "Некорректные настройки оформления")
		return
	}

	actor := middleware.UserID(r.Context())
	if err := h.theme.SetPlatform(r.Context(), body.Settings, actor); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить оформление")
		return
	}
	h.audit.Log(r.Context(), actor, "theme.update", "platform", "", nil)
	writeJSON(w, http.StatusOK, map[string]any{"settings": body.Settings})
}

func statusForRepoError(err error) int {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, repository.ErrEmailTaken):
		return http.StatusConflict
	case errors.Is(err, repository.ErrLastAdminGone):
		return http.StatusConflict
	default:
		return http.StatusBadRequest
	}
}
