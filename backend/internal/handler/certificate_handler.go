package handler

import (
	"net/http"
	"strings"

	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

// CertificateHandler — выдача и проверка сертификатов.
// Проверка по номеру доступна без авторизации: ссылку можно показать работодателю.
type CertificateHandler struct {
	certs *repository.CertificateRepo
}

func NewCertificateHandler(certs *repository.CertificateRepo) *CertificateHandler {
	return &CertificateHandler{certs: certs}
}

// PublicRoutes — /api/certificates.
func (h *CertificateHandler) PublicRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/{serial}", h.verify)
	return r
}

// AdminRoutes — /api/admin/certificates.
func (h *CertificateHandler) AdminRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/{id}/revoke", h.revoke)
	r.Post("/{id}/restore", h.restore)
	return r
}

// GET /api/certificates/{serial} — проверка подлинности.
func (h *CertificateHandler) verify(w http.ResponseWriter, r *http.Request) {
	serial := strings.TrimSpace(chi.URLParam(r, "serial"))

	cert, err := h.certs.GetBySerial(r.Context(), serial)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]any{
			"valid":   false,
			"message": "Сертификат с таким номером не найден",
		})
		return
	}

	if cert.RevokedAt != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"valid":       false,
			"message":     "Сертификат отозван",
			"certificate": cert,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"valid":       true,
		"certificate": cert,
	})
}

func (h *CertificateHandler) list(w http.ResponseWriter, r *http.Request) {
	certs, err := h.certs.List(r.Context(), queryInt(r, "limit", 100, 1, 500))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить сертификаты")
		return
	}
	writeJSON(w, http.StatusOK, certs)
}

func (h *CertificateHandler) revoke(w http.ResponseWriter, r *http.Request) {
	if err := h.certs.Revoke(r.Context(), chi.URLParam(r, "id"), true); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось отозвать сертификат")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Сертификат отозван"})
}

func (h *CertificateHandler) restore(w http.ResponseWriter, r *http.Request) {
	if err := h.certs.Revoke(r.Context(), chi.URLParam(r, "id"), false); err != nil {
		writeError(w, statusForRepoError(err), "Не удалось восстановить сертификат")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Сертификат восстановлен"})
}
