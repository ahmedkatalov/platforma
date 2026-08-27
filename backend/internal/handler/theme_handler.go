package handler

import (
	"net/http"

	"platforma/backend/internal/repository"
)

// ThemeHandler отдаёт оформление платформы без авторизации — чтобы экран входа
// сразу открывался в фирменных цветах, заданных администратором.
type ThemeHandler struct {
	theme *repository.ThemeRepo
}

func NewThemeHandler(theme *repository.ThemeRepo) *ThemeHandler {
	return &ThemeHandler{theme: theme}
}

func (h *ThemeHandler) Public(w http.ResponseWriter, r *http.Request) {
	settings, err := h.theme.GetPlatform(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить оформление")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"settings": rawOrNil(settings)})
}
