package handler

import (
	"encoding/json"
	"net/http"

	"platforma/backend/internal/repository"
)

// ContactsHandler отдаёт контакты для связи без авторизации — чтобы кнопки
// «Написать в Telegram/WhatsApp» работали и на экране входа.
type ContactsHandler struct {
	contacts *repository.ContactsRepo
}

func NewContactsHandler(contacts *repository.ContactsRepo) *ContactsHandler {
	return &ContactsHandler{contacts: contacts}
}

// contactsEnabled — включена ли связь администратором.
func contactsEnabled(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var v struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.Unmarshal(raw, &v); err != nil {
		return false
	}
	return v.Enabled
}

func (h *ContactsHandler) Public(w http.ResponseWriter, r *http.Request) {
	raw, err := h.contacts.Get(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить контакты")
		return
	}
	// Выключено администратором — контакты не раскрываем.
	if !contactsEnabled(raw) {
		writeJSON(w, http.StatusOK, map[string]any{"contacts": map[string]any{"enabled": false}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"contacts": rawOrNil(raw)})
}
