package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

// ReportHandler — выгрузка отчётов по успеваемости в CSV для Excel и Google Sheets.
type ReportHandler struct {
	stats *repository.StatsRepo
	certs *repository.CertificateRepo
}

func NewReportHandler(stats *repository.StatsRepo, certs *repository.CertificateRepo) *ReportHandler {
	return &ReportHandler{stats: stats, certs: certs}
}

func (h *ReportHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/students.csv", h.students)
	r.Get("/certificates.csv", h.certificates)
	return r
}

func (h *ReportHandler) students(w http.ResponseWriter, r *http.Request) {
	items, err := h.stats.StudentsSummary(r.Context(), 500)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать отчёт")
		return
	}

	writer := startCSV(w, "students")
	defer writer.Flush()

	_ = writer.Write([]string{
		"Имя", "Почта", "Статус", "Курсов", "Уроков пройдено", "Уроков всего",
		"Прогресс, %", "Дней посещения", "Минут обучения", "Последний вход",
	})

	for _, item := range items {
		_ = writer.Write([]string{
			item.FullName,
			item.Email,
			statusLabel(item.Status),
			fmt.Sprint(item.Courses),
			fmt.Sprint(item.LessonsCompleted),
			fmt.Sprint(item.LessonsTotal),
			fmt.Sprintf("%.0f", item.Progress),
			fmt.Sprint(item.DaysVisited),
			fmt.Sprint(item.MinutesSpent),
			formatTime(item.LastLoginAt),
		})
	}
}

func (h *ReportHandler) certificates(w http.ResponseWriter, r *http.Request) {
	items, err := h.certs.List(r.Context(), 500)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось собрать отчёт")
		return
	}

	writer := startCSV(w, "certificates")
	defer writer.Flush()

	_ = writer.Write([]string{"Номер", "Студент", "Курс", "Балл", "Уроков", "Выдан", "Отозван"})

	for _, cert := range items {
		revoked := ""
		if cert.RevokedAt != nil {
			revoked = cert.RevokedAt.Format("02.01.2006")
		}
		_ = writer.Write([]string{
			cert.Serial,
			cert.HolderName,
			cert.CourseTitle,
			fmt.Sprintf("%.0f", cert.Score),
			fmt.Sprintf("%d/%d", cert.LessonsCompleted, cert.LessonsTotal),
			cert.IssuedAt.Format("02.01.2006"),
			revoked,
		})
	}
}

// startCSV готовит ответ: заголовки скачивания, BOM и точка с запятой —
// чтобы Excel корректно открыл кириллицу и разбил на колонки.
func startCSV(w http.ResponseWriter, name string) *csv.Writer {
	filename := fmt.Sprintf("%s-%s.csv", name, time.Now().Format("2006-01-02"))

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(w)
	writer.Comma = ';'
	return writer
}

func statusLabel(status string) string {
	switch status {
	case "active":
		return "активен"
	case "invited":
		return "приглашён"
	case "blocked":
		return "заблокирован"
	default:
		return status
	}
}

func formatTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("02.01.2006 15:04")
}
