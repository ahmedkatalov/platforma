package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"platforma/backend/internal/config"
	"platforma/backend/internal/domain"
	"platforma/backend/internal/mailer"
	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

// LessonHandler — прохождение уроков студентом: теория, квизы, терминал и код.
type LessonHandler struct {
	progress *repository.ProgressRepo
	courses  *repository.CourseRepo
	activity *repository.ActivityRepo
	certs    *repository.CertificateRepo
	mail     *mailer.Mailer
	cfg      *config.Config
}

func NewLessonHandler(
	progress *repository.ProgressRepo,
	courses *repository.CourseRepo,
	activity *repository.ActivityRepo,
	certs *repository.CertificateRepo,
	mail *mailer.Mailer,
	cfg *config.Config,
) *LessonHandler {
	return &LessonHandler{progress: progress, courses: courses, activity: activity,
		certs: certs, mail: mail, cfg: cfg}
}

// finishLesson отмечает урок пройденным и, если курс закрыт целиком, выдаёт сертификат.
// Возвращает сертификат, если он выдан именно сейчас.
func (h *LessonHandler) finishLesson(
	r *http.Request,
	userID, lessonID, courseID string,
	score *float64,
	seconds int,
) (*repository.Certificate, error) {
	if err := h.progress.Complete(r.Context(), userID, lessonID, score, seconds); err != nil {
		return nil, err
	}
	_ = h.activity.Touch(r.Context(), userID, seconds)

	completion, err := h.certs.Completion(r.Context(), userID, courseID)
	if err != nil || completion.Total == 0 || completion.Completed < completion.Total {
		return nil, nil
	}

	// Курс пройден целиком — выдаём сертификат (повторная выдача просто обновит данные).
	existing, err := h.certs.ForUser(r.Context(), userID)
	if err == nil {
		for _, cert := range existing {
			if cert.CourseID == courseID {
				return nil, nil
			}
		}
	}

	cert, err := h.certs.Issue(r.Context(), userID, courseID, *completion)
	if err != nil {
		log.Printf("сертификат: не удалось выдать: %v", err)
		return nil, nil
	}

	go h.notifyCertificate(cert)
	return cert, nil
}

// notifyCertificate отправляет письмо о выдаче сертификата, не задерживая ответ.
func (h *LessonHandler) notifyCertificate(cert *repository.Certificate) {
	if !h.mail.Enabled() {
		log.Printf("сертификат %s выдан для %s (письмо не отправлено — EmailJS не настроен)",
			cert.Serial, cert.HolderName)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	link := strings.TrimRight(h.cfg.PublicBaseURL, "/") + "/certificates/" + cert.Serial
	subject := "Курс пройден: " + cert.CourseTitle
	message := fmt.Sprintf(
		"Поздравляем! Курс «%s» пройден полностью. Ваш сертификат № %s доступен по ссылке: %s",
		cert.CourseTitle, cert.Serial, link)

	if err := h.mail.SendNotice(ctx, cert.HolderEmail, subject, message, link); err != nil {
		log.Printf("сертификат: письмо не отправлено: %v", err)
	}
}

func (h *LessonHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", h.get)
	r.Post("/{id}/start", h.start)
	r.Post("/{id}/complete", h.complete)
	r.Post("/{id}/quiz", h.submitQuiz)
	r.Post("/{id}/terminal", h.checkTerminal)
	r.Post("/{id}/code", h.checkCode)
	return r
}

// access проверяет, что студент записан на курс урока (админу доступно всё).
func (h *LessonHandler) access(r *http.Request, lessonID string) (*repository.LessonContext, error) {
	lesson, err := h.progress.LessonWithCourse(r.Context(), lessonID)
	if err != nil {
		return nil, err
	}

	if middleware.Role(r.Context()) == domain.RoleAdmin {
		return lesson, nil
	}
	if !lesson.Published {
		return nil, repository.ErrNotFound
	}

	enrolled, err := h.courses.IsEnrolled(r.Context(), middleware.UserID(r.Context()), lesson.CourseID)
	if err != nil {
		return nil, err
	}
	if !enrolled {
		return nil, errNoAccess
	}
	return lesson, nil
}

var errNoAccess = errors.New("курс вам ещё не открыт — обратитесь к администратору")

func (h *LessonHandler) writeAccessError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "Урок не найден")
	case errors.Is(err, errNoAccess):
		writeError(w, http.StatusForbidden, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "Не удалось открыть урок")
	}
}

// GET /api/lessons/{id} — урок без правильных ответов, плюс прогресс студента.
func (h *LessonHandler) get(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "id")

	lesson, err := h.access(r, lessonID)
	if err != nil {
		h.writeAccessError(w, err)
		return
	}

	userID := middleware.UserID(r.Context())
	lesson.Lesson.Content = domain.SanitizeContent(lesson.Lesson.Kind, lesson.Lesson.Content)

	progress, err := h.progress.ForCourse(r.Context(), userID, lesson.CourseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить прогресс")
		return
	}

	tasks, err := h.progress.Tasks(r.Context(), userID, lessonID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить задания")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"lesson":       lesson.Lesson,
		"courseId":     lesson.CourseID,
		"courseSlug":   lesson.CourseSlug,
		"courseTitle":  lesson.CourseTitle,
		"moduleTitle":  lesson.ModuleTitle,
		"prevLessonId": lesson.PrevID,
		"nextLessonId": lesson.NextID,
		"progress":     progress,
		"tasks":        tasks,
	})
}

func (h *LessonHandler) start(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "id")
	if _, err := h.access(r, lessonID); err != nil {
		h.writeAccessError(w, err)
		return
	}

	if err := h.progress.Start(r.Context(), middleware.UserID(r.Context()), lessonID); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось начать урок")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

// POST /api/lessons/{id}/complete — отметить теорию пройденной.
func (h *LessonHandler) complete(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Seconds int `json:"seconds"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lessonID := chi.URLParam(r, "id")
	lesson, err := h.access(r, lessonID)
	if err != nil {
		h.writeAccessError(w, err)
		return
	}

	userID := middleware.UserID(r.Context())
	seconds := clampSeconds(body.Seconds)

	cert, err := h.finishLesson(r, userID, lessonID, lesson.CourseID, nil, seconds)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить прогресс")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message":     "Урок отмечен пройденным",
		"certificate": cert,
	})
}

// POST /api/lessons/{id}/quiz — проверка ответов на сервере.
func (h *LessonHandler) submitQuiz(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Answers []domain.QuizAnswer `json:"answers"`
		Seconds int                 `json:"seconds"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lessonID := chi.URLParam(r, "id")
	lesson, err := h.access(r, lessonID)
	if err != nil {
		h.writeAccessError(w, err)
		return
	}
	if lesson.Lesson.Kind != domain.LessonQuiz {
		writeError(w, http.StatusBadRequest, "Этот урок не является квизом")
		return
	}

	quiz, err := domain.ParseQuiz(lesson.Lesson.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Некорректное содержимое квиза")
		return
	}

	result := domain.GradeQuiz(quiz, body.Answers)
	userID := middleware.UserID(r.Context())
	seconds := clampSeconds(body.Seconds)

	attemptID, err := h.progress.SaveAttempt(r.Context(), userID, lessonID, domain.LessonQuiz,
		result.Score, result.CorrectCount, result.TotalCount, result.Passed, seconds,
		map[string]any{"questions": result.Questions})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить попытку")
		return
	}

	timings := make(map[string]int, len(body.Answers))
	for _, answer := range body.Answers {
		timings[answer.QuestionID] = clampSeconds(answer.SecondsSpent)
	}
	if err := h.progress.SaveQuizAnswers(r.Context(), attemptID, userID, lessonID,
		result.Questions, timings); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить ответы")
		return
	}

	score := result.Score
	var cert *repository.Certificate

	if result.Passed {
		cert, err = h.finishLesson(r, userID, lessonID, lesson.CourseID, &score, seconds)
	} else {
		err = h.progress.TouchAttempt(r.Context(), userID, lessonID, &score, seconds)
		_ = h.activity.Touch(r.Context(), userID, seconds)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить прогресс")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"score":        result.Score,
		"passed":       result.Passed,
		"correctCount": result.CorrectCount,
		"totalCount":   result.TotalCount,
		"passScore":    result.PassScore,
		"questions":    result.Questions,
		"certificate":  cert,
	})
}

// POST /api/lessons/{id}/terminal — проверка одной команды из задания.
func (h *LessonHandler) checkTerminal(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TaskID   string `json:"taskId"`
		Command  string `json:"command"`
		UsedHint bool   `json:"usedHint"`
		Seconds  int    `json:"seconds"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lessonID := chi.URLParam(r, "id")
	lesson, err := h.access(r, lessonID)
	if err != nil {
		h.writeAccessError(w, err)
		return
	}
	if lesson.Lesson.Kind != domain.LessonTerminal {
		writeError(w, http.StatusBadRequest, "Этот урок не является тренажёром терминала")
		return
	}

	terminal, err := domain.ParseTerminal(lesson.Lesson.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Некорректное содержимое урока")
		return
	}

	var task *domain.TerminalTask
	for i := range terminal.Tasks {
		if terminal.Tasks[i].ID == body.TaskID {
			task = &terminal.Tasks[i]
			break
		}
	}
	if task == nil {
		writeError(w, http.StatusNotFound, "Задание не найдено")
		return
	}

	solved := domain.MatchTerminalCommand(task, body.Command)
	userID := middleware.UserID(r.Context())

	if err := h.progress.MarkTask(r.Context(), userID, lessonID, task.ID, solved, body.UsedHint); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить попытку")
		return
	}

	completed, err := h.progress.CompletedTaskIDs(r.Context(), userID, lessonID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось проверить прогресс")
		return
	}

	allDone := len(terminal.Tasks) > 0
	for _, t := range terminal.Tasks {
		if !completed[t.ID] {
			allDone = false
			break
		}
	}

	seconds := clampSeconds(body.Seconds)
	var cert *repository.Certificate

	if allDone {
		score := 100.0
		cert, err = h.finishLesson(r, userID, lessonID, lesson.CourseID, &score, seconds)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Не удалось сохранить прогресс")
			return
		}
	}

	message := task.Success
	if solved && message == "" {
		message = "Верно! Задание выполнено."
	}
	if !solved {
		message = "Пока не то. Проверьте команду и попробуйте ещё раз."
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"solved":         solved,
		"message":        message,
		"hint":           task.Hint,
		"completedTasks": keys(completed),
		"lessonComplete": allDone,
		"certificate":    cert,
	})
}

// POST /api/lessons/{id}/code — проверка решения в редакторе кода.
func (h *LessonHandler) checkCode(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code    string `json:"code"`
		Seconds int    `json:"seconds"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	lessonID := chi.URLParam(r, "id")
	lesson, err := h.access(r, lessonID)
	if err != nil {
		h.writeAccessError(w, err)
		return
	}
	if lesson.Lesson.Kind != domain.LessonCode {
		writeError(w, http.StatusBadRequest, "Этот урок не является практикой с кодом")
		return
	}

	code, err := domain.ParseCode(lesson.Lesson.Content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Некорректное содержимое урока")
		return
	}

	type checkResult struct {
		OK      bool   `json:"ok"`
		Message string `json:"message"`
	}

	results := make([]checkResult, 0, len(code.Checks))
	passedAll := true

	for _, check := range code.Checks {
		ok := domain.RunCodeCheck(check, body.Code)
		if !ok {
			passedAll = false
		}
		message := check.Message
		if message == "" {
			message = domain.DescribeCodeCheck(check)
		}
		results = append(results, checkResult{OK: ok, Message: message})
	}

	// Урок без проверок засчитывается по факту отправки решения.
	if len(code.Checks) == 0 {
		passedAll = strings.TrimSpace(body.Code) != ""
	}

	okCount := 0
	for _, result := range results {
		if result.OK {
			okCount++
		}
	}

	score := 0.0
	if len(results) > 0 {
		score = float64(okCount) / float64(len(results)) * 100
	} else if passedAll {
		score = 100
	}

	userID := middleware.UserID(r.Context())
	seconds := clampSeconds(body.Seconds)

	if _, err := h.progress.SaveAttempt(r.Context(), userID, lessonID, domain.LessonCode,
		score, okCount, len(results), passedAll, seconds,
		map[string]any{"checks": results}); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить попытку")
		return
	}

	var cert *repository.Certificate
	if passedAll {
		cert, err = h.finishLesson(r, userID, lessonID, lesson.CourseID, &score, seconds)
	} else {
		err = h.progress.TouchAttempt(r.Context(), userID, lessonID, &score, seconds)
		_ = h.activity.Touch(r.Context(), userID, seconds)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить прогресс")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"passed":      passedAll,
		"score":       score,
		"checks":      results,
		"hint":        code.Hint,
		"certificate": cert,
	})
}

func clampSeconds(seconds int) int {
	if seconds < 0 {
		return 0
	}
	if seconds > 3600 {
		return 3600
	}
	return seconds
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
