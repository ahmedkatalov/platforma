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
	access   *repository.AccessRepo
	contacts *repository.ContactsRepo
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
	access *repository.AccessRepo,
	contacts *repository.ContactsRepo,
	userSvc *service.UserService,
) *AdminHandler {
	return &AdminHandler{users: users, courses: courses, stats: stats,
		activity: activity, audit: audit, theme: theme, progress: progress,
		access: access, contacts: contacts, userSvc: userSvc}
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
		r.Get("/{id}/module-access", h.getModuleAccess)
		r.Post("/{id}/module-access", h.setModuleAccess)
	})

	r.Get("/students-progress", h.studentsProgress)
	r.Get("/requests-count", h.requestsCount)

	r.Route("/access-requests", func(r chi.Router) {
		r.Get("/", h.listAccessRequests)
		r.Post("/{id}/approve", h.approveAccessRequest)
		r.Post("/{id}/reject", h.rejectAccessRequest)
	})

	r.Route("/course-requests", func(r chi.Router) {
		r.Get("/", h.listCourseRequests)
		r.Post("/{id}/approve", h.approveCourseRequest)
		r.Post("/{id}/reject", h.rejectCourseRequest)
	})

	r.Route("/theme", func(r chi.Router) {
		r.Get("/", h.getTheme)
		r.Put("/", h.putTheme)
	})

	r.Route("/contacts", func(r chi.Router) {
		r.Get("/", h.getContacts)
		r.Put("/", h.putContacts)
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
	// Сразу открываем первую главу, чтобы студенту было с чего начать.
	if err := h.courses.GrantFirstModuleAccess(r.Context(), userID, body.CourseID, middleware.UserID(r.Context())); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось открыть первую главу")
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

// --- Пошаговый доступ к главам ---

func (h *AdminHandler) getModuleAccess(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	courseID := r.URL.Query().Get("courseId")
	if courseID == "" {
		writeError(w, http.StatusBadRequest, "Не указан курс")
		return
	}
	m, err := h.courses.ModuleAccessMap(r.Context(), userID, courseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить доступы")
		return
	}
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	writeJSON(w, http.StatusOK, map[string]any{"granted": ids})
}

func (h *AdminHandler) setModuleAccess(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var body struct {
		ModuleID string `json:"moduleId"`
		Granted  bool   `json:"granted"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if body.ModuleID == "" {
		writeError(w, http.StatusBadRequest, "Не указана глава")
		return
	}

	actor := middleware.UserID(r.Context())
	action := "module.revoke"
	if body.Granted {
		action = "module.grant"
		if err := h.courses.GrantModuleAccess(r.Context(), userID, body.ModuleID, actor); err != nil {
			writeError(w, http.StatusInternalServerError, "Не удалось открыть главу")
			return
		}
	} else {
		if err := h.courses.RevokeModuleAccess(r.Context(), userID, body.ModuleID); err != nil {
			writeError(w, http.StatusInternalServerError, "Не удалось закрыть главу")
			return
		}
	}
	h.audit.Log(r.Context(), actor, action, "user", userID, map[string]any{"moduleId": body.ModuleID})
	writeJSON(w, http.StatusOK, map[string]any{"granted": body.Granted})
}

// --- Заявки на доступ ---

func (h *AdminHandler) listAccessRequests(w http.ResponseWriter, r *http.Request) {
	items, err := h.access.List(r.Context(), r.URL.Query().Get("status"), queryInt(r, "limit", 200, 1, 500))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить заявки")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *AdminHandler) approveAccessRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, moduleID, status, err := h.access.Get(r.Context(), id)
	if err != nil {
		writeError(w, statusForRepoError(err), "Заявка не найдена")
		return
	}
	if status != "pending" {
		writeError(w, http.StatusConflict, "Заявка уже обработана")
		return
	}

	actor := middleware.UserID(r.Context())
	if err := h.courses.GrantModuleAccess(r.Context(), userID, moduleID, actor); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось открыть главу")
		return
	}
	if err := h.access.Decide(r.Context(), id, "approved", actor); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось обновить заявку")
		return
	}
	h.audit.Log(r.Context(), actor, "access.approve", "user", userID, map[string]any{"moduleId": moduleID})
	writeJSON(w, http.StatusOK, map[string]string{"message": "Доступ открыт"})
}

func (h *AdminHandler) rejectAccessRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actor := middleware.UserID(r.Context())
	if err := h.access.Decide(r.Context(), id, "rejected", actor); err != nil {
		writeError(w, statusForRepoError(err), "Заявка уже обработана или не найдена")
		return
	}
	h.audit.Log(r.Context(), actor, "access.reject", "request", id, nil)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Заявка отклонена"})
}

// --- Заявки на курсы (запись на курс) ---

// requestsCount — сколько заявок ждут решения (для бейджа-уведомления).
func (h *AdminHandler) requestsCount(w http.ResponseWriter, r *http.Request) {
	chapters, courses, err := h.access.PendingCounts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось посчитать заявки")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{
		"chapters": chapters,
		"courses":  courses,
		"total":    chapters + courses,
	})
}

func (h *AdminHandler) listCourseRequests(w http.ResponseWriter, r *http.Request) {
	items, err := h.access.ListCourseRequests(r.Context(), r.URL.Query().Get("status"), queryInt(r, "limit", 200, 1, 500))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить заявки")
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *AdminHandler) approveCourseRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	userID, courseID, status, err := h.access.GetCourseRequest(r.Context(), id)
	if err != nil {
		writeError(w, statusForRepoError(err), "Заявка не найдена")
		return
	}
	if status != "pending" {
		writeError(w, http.StatusConflict, "Заявка уже обработана")
		return
	}

	actor := middleware.UserID(r.Context())
	// Записываем студента на курс и открываем первую главу.
	if err := h.courses.Enroll(r.Context(), userID, courseID, actor, nil); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось записать на курс")
		return
	}
	if err := h.courses.GrantFirstModuleAccess(r.Context(), userID, courseID, actor); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось открыть первую главу")
		return
	}
	if err := h.access.DecideCourseRequest(r.Context(), id, "approved", actor); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось обновить заявку")
		return
	}
	h.audit.Log(r.Context(), actor, "course.request.approve", "user", userID, map[string]any{"courseId": courseID})
	writeJSON(w, http.StatusOK, map[string]string{"message": "Студент записан на курс"})
}

func (h *AdminHandler) rejectCourseRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	actor := middleware.UserID(r.Context())
	if err := h.access.DecideCourseRequest(r.Context(), id, "rejected", actor); err != nil {
		writeError(w, statusForRepoError(err), "Заявка уже обработана или не найдена")
		return
	}
	h.audit.Log(r.Context(), actor, "course.request.reject", "request", id, nil)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Заявка отклонена"})
}

// --- Контакты для связи (Telegram/WhatsApp) ---

func (h *AdminHandler) getContacts(w http.ResponseWriter, r *http.Request) {
	settings, err := h.contacts.Get(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить контакты")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"settings": rawOrNil(settings)})
}

func (h *AdminHandler) putContacts(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Settings json.RawMessage `json:"settings"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(body.Settings) == 0 || !json.Valid(body.Settings) {
		writeError(w, http.StatusBadRequest, "Некорректные контакты")
		return
	}
	actor := middleware.UserID(r.Context())
	if err := h.contacts.Set(r.Context(), body.Settings, actor); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить контакты")
		return
	}
	h.audit.Log(r.Context(), actor, "contacts.update", "platform", "", nil)
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
