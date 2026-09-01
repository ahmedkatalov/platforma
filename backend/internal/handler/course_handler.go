package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"platforma/backend/internal/domain"
	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

var slugRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// CourseHandler отдаёт курсы студентам (только опубликованные и доступные)
// и обслуживает редактор курсов у администратора.
type CourseHandler struct {
	courses  *repository.CourseRepo
	audit    *repository.AuditRepo
	progress *repository.ProgressRepo
}

func NewCourseHandler(
	courses *repository.CourseRepo,
	audit *repository.AuditRepo,
	progress *repository.ProgressRepo,
) *CourseHandler {
	return &CourseHandler{courses: courses, audit: audit, progress: progress}
}

// StudentRoutes — /api/courses для любого авторизованного пользователя.
func (h *CourseHandler) StudentRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.listForStudent)
	r.Get("/{slug}", h.getForStudent)
	return r
}

// AdminRoutes — /api/admin/courses, полный доступ к структуре курса.
func (h *CourseHandler) AdminRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.listAll)
	r.Post("/", h.createCourse)
	r.Post("/import", h.importCourse)
	r.Get("/{id}/export", h.exportCourse)
	r.Get("/{id}", h.getFull)
	r.Put("/{id}", h.updateCourse)
	r.Delete("/{id}", h.deleteCourse)

	r.Post("/{id}/modules", h.createModule)
	r.Put("/modules/{moduleId}", h.updateModule)
	r.Delete("/modules/{moduleId}", h.deleteModule)

	r.Post("/modules/{moduleId}/lessons", h.createLesson)
	r.Get("/lessons/{lessonId}", h.getLesson)
	r.Put("/lessons/{lessonId}", h.updateLesson)
	r.Delete("/lessons/{lessonId}", h.deleteLesson)
	return r
}

// --- Студент ---

func (h *CourseHandler) listForStudent(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courses.List(r.Context(), "published")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курсы")
		return
	}

	enrollments, err := h.courses.EnrollmentsForUser(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить записи на курсы")
		return
	}
	enrolled := make(map[string]bool, len(enrollments))
	for _, e := range enrollments {
		enrolled[e.CourseID] = true
	}

	// Сколько уроков студент уже прошёл в каждом курсе — для прогресс-баров.
	completed, err := h.progress.CompletedCountByCourse(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить прогресс")
		return
	}

	items := make([]map[string]any, 0, len(courses))
	for _, c := range courses {
		items = append(items, map[string]any{
			"course":           c,
			"enrolled":         enrolled[c.ID],
			"completedLessons": completed[c.ID],
		})
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *CourseHandler) getForStudent(w http.ResponseWriter, r *http.Request) {
	course, err := h.courses.GetBySlug(r.Context(), chi.URLParam(r, "slug"))
	if err != nil {
		writeError(w, http.StatusNotFound, "Курс не найден")
		return
	}

	userID := middleware.UserID(r.Context())
	isAdmin := middleware.Role(r.Context()) == domain.RoleAdmin

	if course.Status != "published" && !isAdmin {
		writeError(w, http.StatusNotFound, "Курс не найден")
		return
	}

	enrolled, err := h.courses.IsEnrolled(r.Context(), userID, course.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось проверить доступ")
		return
	}
	if !enrolled && !isAdmin {
		// Структуру показываем всем, чтобы студент видел программу курса.
		if err := h.courses.WithContent(r.Context(), course, false); err != nil {
			writeError(w, http.StatusInternalServerError, "Не удалось загрузить курс")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"course":   course,
			"enrolled": false,
			"progress": []any{},
		})
		return
	}

	// Тела уроков сюда не отдаём — они приходят по одному через /api/lessons/{id}.
	if err := h.courses.WithContent(r.Context(), course, false); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курс")
		return
	}

	progress, err := h.progress.ForCourse(r.Context(), userID, course.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить прогресс")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"course":   course,
		"enrolled": true,
		"progress": progress,
	})
}

// --- Администратор ---

func (h *CourseHandler) listAll(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courses.List(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курсы")
		return
	}
	writeJSON(w, http.StatusOK, courses)
}

func (h *CourseHandler) getFull(w http.ResponseWriter, r *http.Request) {
	course, err := h.courses.GetByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "Курс не найден")
		return
	}
	if err := h.courses.WithContent(r.Context(), course, true); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить структуру курса")
		return
	}
	writeJSON(w, http.StatusOK, course)
}

func (h *CourseHandler) createCourse(w http.ResponseWriter, r *http.Request) {
	var body repository.CourseInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateCourse(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actor := middleware.UserID(r.Context())
	course, err := h.courses.Create(r.Context(), body, actor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	h.audit.Log(r.Context(), actor, "course.create", "course", course.ID,
		map[string]any{"slug": course.Slug})
	writeJSON(w, http.StatusCreated, course)
}

func (h *CourseHandler) updateCourse(w http.ResponseWriter, r *http.Request) {
	var body repository.CourseInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateCourse(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id := chi.URLParam(r, "id")
	course, err := h.courses.Update(r.Context(), id, body)
	if err != nil {
		writeError(w, statusForRepoError(err), err.Error())
		return
	}
	h.audit.Log(r.Context(), middleware.UserID(r.Context()), "course.update", "course", id, nil)
	writeJSON(w, http.StatusOK, course)
}

func (h *CourseHandler) deleteCourse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.courses.Delete(r.Context(), id); err != nil {
		writeError(w, statusForRepoError(err), err.Error())
		return
	}
	h.audit.Log(r.Context(), middleware.UserID(r.Context()), "course.delete", "course", id, nil)
	writeJSON(w, http.StatusOK, map[string]string{"message": "Курс удалён"})
}

func (h *CourseHandler) createModule(w http.ResponseWriter, r *http.Request) {
	var body repository.ModuleInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if strings.TrimSpace(body.Title) == "" {
		writeError(w, http.StatusBadRequest, "Укажите название модуля")
		return
	}

	module, err := h.courses.CreateModule(r.Context(), chi.URLParam(r, "id"), body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось создать модуль")
		return
	}
	writeJSON(w, http.StatusCreated, module)
}

func (h *CourseHandler) updateModule(w http.ResponseWriter, r *http.Request) {
	var body repository.ModuleInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	module, err := h.courses.UpdateModule(r.Context(), chi.URLParam(r, "moduleId"), body)
	if err != nil {
		writeError(w, statusForRepoError(err), "Не удалось обновить модуль")
		return
	}
	writeJSON(w, http.StatusOK, module)
}

func (h *CourseHandler) deleteModule(w http.ResponseWriter, r *http.Request) {
	if err := h.courses.DeleteModule(r.Context(), chi.URLParam(r, "moduleId")); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось удалить модуль")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Модуль удалён"})
}

func (h *CourseHandler) createLesson(w http.ResponseWriter, r *http.Request) {
	var body repository.LessonInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateLesson(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lesson, err := h.courses.CreateLesson(r.Context(), chi.URLParam(r, "moduleId"), body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось создать урок")
		return
	}
	writeJSON(w, http.StatusCreated, lesson)
}

func (h *CourseHandler) getLesson(w http.ResponseWriter, r *http.Request) {
	lesson, err := h.courses.GetLesson(r.Context(), chi.URLParam(r, "lessonId"))
	if err != nil {
		writeError(w, http.StatusNotFound, "Урок не найден")
		return
	}
	writeJSON(w, http.StatusOK, lesson)
}

func (h *CourseHandler) updateLesson(w http.ResponseWriter, r *http.Request) {
	var body repository.LessonInput
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateLesson(&body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lesson, err := h.courses.UpdateLesson(r.Context(), chi.URLParam(r, "lessonId"), body)
	if err != nil {
		writeError(w, statusForRepoError(err), "Не удалось обновить урок")
		return
	}
	writeJSON(w, http.StatusOK, lesson)
}

func (h *CourseHandler) deleteLesson(w http.ResponseWriter, r *http.Request) {
	if err := h.courses.DeleteLesson(r.Context(), chi.URLParam(r, "lessonId")); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось удалить урок")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Урок удалён"})
}

func validateCourse(in *repository.CourseInput) error {
	in.Slug = strings.ToLower(strings.TrimSpace(in.Slug))
	in.Title = strings.TrimSpace(in.Title)

	if in.Title == "" {
		return errors.New("укажите название курса")
	}
	if !slugRe.MatchString(in.Slug) {
		return errors.New("slug должен состоять из латиницы, цифр и дефисов: например devops-basics")
	}
	switch in.Level {
	case "", "beginner":
		in.Level = "beginner"
	case "intermediate", "advanced":
	default:
		return errors.New("недопустимый уровень курса")
	}
	switch in.Status {
	case "", "draft":
		in.Status = "draft"
	case "published", "archived":
	default:
		return errors.New("недопустимый статус курса")
	}
	return nil
}

func validateLesson(in *repository.LessonInput) error {
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		return errors.New("укажите название урока")
	}
	switch in.Kind {
	case "", domain.LessonText:
		in.Kind = domain.LessonText
	case domain.LessonQuiz, domain.LessonTerminal, domain.LessonCode:
	default:
		return errors.New("недопустимый тип урока")
	}
	if in.DurationMin <= 0 {
		in.DurationMin = 10
	}
	return nil
}
