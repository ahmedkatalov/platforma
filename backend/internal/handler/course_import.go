package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

// Формат «курс как файл»: один JSON-пакет со всей структурой курса,
// который админ загружает в интерфейсе и получает готовый курс.

const maxPackageSize = 64 << 20 // 64 МБ — курс целиком с квизами и практиками

type pkgLesson struct {
	Title       string          `json:"title"`
	Kind        string          `json:"kind"`
	Summary     string          `json:"summary"`
	Content     json.RawMessage `json:"content"`
	DurationMin int             `json:"durationMin"`
}

type pkgModule struct {
	Title   string      `json:"title"`
	Summary string      `json:"summary"`
	Lessons []pkgLesson `json:"lessons"`
}

type pkgCourseMeta struct {
	Slug        string   `json:"slug"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	CoverURL    string   `json:"coverUrl"`
	Level       string   `json:"level"`
	Tags        []string `json:"tags"`
}

// coursePackage — формат обмена курсом. format/version для будущей совместимости.
type coursePackage struct {
	Format  string        `json:"format"`
	Version int           `json:"version"`
	Course  pkgCourseMeta `json:"course"`
	Modules []pkgModule   `json:"modules"`
}

const packageFormat = "platforma-course"
const packageVersion = 1

// importCourse принимает JSON-пакет курса и создаёт курс целиком.
// POST /api/admin/courses/import?replace=true — если курс с таким slug есть, пересоздать.
func (h *CourseHandler) importCourse(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxPackageSize)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Файл слишком большой или не удалось его прочитать")
		return
	}
	// Убираем возможный BOM и пробелы в начале/конце.
	body = bytes.TrimSpace(bytes.TrimPrefix(body, []byte("\xEF\xBB\xBF")))
	if len(body) == 0 {
		writeError(w, http.StatusBadRequest, "Файл пустой")
		return
	}
	if body[0] != '{' {
		writeError(w, http.StatusBadRequest,
			"Это не JSON-пакет курса: файл начинается не с «{». Загрузите файл вида «...course.json»")
		return
	}
	var pkg coursePackage
	if err := json.Unmarshal(body, &pkg); err != nil {
		writeError(w, http.StatusBadRequest, "Не удалось разобрать JSON курса: "+err.Error())
		return
	}
	if pkg.Format != "" && pkg.Format != packageFormat {
		writeError(w, http.StatusBadRequest, "Неизвестный формат файла курса")
		return
	}
	if len(pkg.Modules) == 0 {
		writeError(w, http.StatusBadRequest, "В файле нет модулей курса")
		return
	}

	in := repository.CourseInput{
		Slug:        pkg.Course.Slug,
		Title:       pkg.Course.Title,
		Subtitle:    pkg.Course.Subtitle,
		Description: pkg.Course.Description,
		CoverURL:    pkg.Course.CoverURL,
		Level:       pkg.Course.Level,
		Tags:        pkg.Course.Tags,
		Status:      "draft", // грузим как черновик — публикует админ вручную
	}
	if err := validateCourse(&in); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	actor := middleware.UserID(r.Context())
	replace := r.URL.Query().Get("replace") == "true"

	if existing, err := h.courses.GetBySlug(r.Context(), in.Slug); err == nil {
		if !replace {
			writeError(w, http.StatusConflict,
				"Курс с адресом «"+in.Slug+"» уже существует. Отметьте «заменить», чтобы пересоздать его из файла.")
			return
		}
		if err := h.courses.Delete(r.Context(), existing.ID); err != nil {
			writeError(w, http.StatusInternalServerError, "Не удалось удалить старую версию курса")
			return
		}
	}

	// Новый курс добавляем в конец списка.
	if all, err := h.courses.List(r.Context(), ""); err == nil {
		in.Position = len(all) + 1
	}

	course, err := h.courses.Create(r.Context(), in, actor)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// При любой ошибке ниже откатываем частично созданный курс, чтобы не оставить обрубок.
	fail := func(status int, message string) {
		_ = h.courses.Delete(r.Context(), course.ID)
		writeError(w, status, message)
	}

	lessonCount := 0
	for mi, m := range pkg.Modules {
		if strings.TrimSpace(m.Title) == "" {
			fail(http.StatusBadRequest, "У модуля нет названия")
			return
		}
		module, err := h.courses.CreateModule(r.Context(), course.ID, repository.ModuleInput{
			Title: m.Title, Summary: m.Summary, Position: mi + 1,
		})
		if err != nil {
			fail(http.StatusBadRequest, "Не удалось создать модуль «"+m.Title+"»")
			return
		}
		for li, l := range m.Lessons {
			lin := repository.LessonInput{
				Title: l.Title, Kind: l.Kind, Summary: l.Summary,
				Content: l.Content, DurationMin: l.DurationMin, Position: li + 1,
			}
			if err := validateLesson(&lin); err != nil {
				fail(http.StatusBadRequest, "Урок «"+l.Title+"»: "+err.Error())
				return
			}
			if len(strings.TrimSpace(string(lin.Content))) == 0 {
				lin.Content = json.RawMessage("{}")
			}
			if _, err := h.courses.CreateLesson(r.Context(), module.ID, lin); err != nil {
				fail(http.StatusBadRequest, "Не удалось создать урок «"+l.Title+"»")
				return
			}
			lessonCount++
		}
	}

	h.audit.Log(r.Context(), actor, "course.import", "course", course.ID,
		map[string]any{"slug": course.Slug, "modules": len(pkg.Modules), "lessons": lessonCount})
	writeJSON(w, http.StatusCreated, map[string]any{
		"course":  course,
		"modules": len(pkg.Modules),
		"lessons": lessonCount,
		"message": "Курс загружен как черновик. Опубликуйте его, когда будете готовы.",
	})
}

// exportCourse отдаёт курс целиком одним JSON-пакетом (для скачивания файлом).
// GET /api/admin/courses/{id}/export
func (h *CourseHandler) exportCourse(w http.ResponseWriter, r *http.Request) {
	course, err := h.courses.GetByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "Курс не найден")
		return
	}
	if err := h.courses.WithContent(r.Context(), course, true); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить курс")
		return
	}

	pkg := coursePackage{
		Format:  packageFormat,
		Version: packageVersion,
		Course: pkgCourseMeta{
			Slug: course.Slug, Title: course.Title, Subtitle: course.Subtitle,
			Description: course.Description, CoverURL: course.CoverURL,
			Level: course.Level, Tags: course.Tags,
		},
	}
	for _, m := range course.Modules {
		mod := pkgModule{Title: m.Title, Summary: m.Summary}
		for _, l := range m.Lessons {
			mod.Lessons = append(mod.Lessons, pkgLesson{
				Title: l.Title, Kind: l.Kind, Summary: l.Summary,
				Content: l.Content, DurationMin: l.DurationMin,
			})
		}
		pkg.Modules = append(pkg.Modules, mod)
	}

	filename := course.Slug
	if filename == "" {
		filename = "course"
	}
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`.course.json"`)
	writeJSON(w, http.StatusOK, pkg)
}
