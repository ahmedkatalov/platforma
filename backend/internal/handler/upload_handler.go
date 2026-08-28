package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

const maxUploadSize = 5 << 20 // 5 MB

// Разрешаем только картинки — файлы отдаются статикой, исполняемому контенту здесь не место.
var allowedMime = map[string]string{
	"image/png":     ".png",
	"image/jpeg":    ".jpg",
	"image/gif":     ".gif",
	"image/webp":    ".webp",
	"image/svg+xml": ".svg",
}

// UploadHandler — картинки и схемы для уроков.
type UploadHandler struct {
	assets *repository.AssetRepo
	dir    string
}

func NewUploadHandler(assets *repository.AssetRepo, dir string) *UploadHandler {
	return &UploadHandler{assets: assets, dir: dir}
}

func (h *UploadHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.upload)
	r.Get("/", h.list)
	r.Delete("/{id}", h.remove)
	return r
}

func (h *UploadHandler) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+1024)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "Файл слишком большой (максимум 5 МБ)")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Файл не передан")
		return
	}
	defer file.Close()

	// Тип определяем по содержимому, а не по расширению из запроса.
	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	mime := http.DetectContentType(head[:n])
	if strings.HasPrefix(mime, "text/xml") || strings.HasPrefix(mime, "text/plain") {
		// SVG определяется как текст — проверяем по содержимому.
		if strings.Contains(strings.ToLower(string(head[:n])), "<svg") {
			mime = "image/svg+xml"
		}
	}

	ext, ok := allowedMime[strings.Split(mime, ";")[0]]
	if !ok {
		writeError(w, http.StatusBadRequest, "Можно загружать только изображения (png, jpg, gif, webp, svg)")
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось прочитать файл")
		return
	}

	name, err := randomName(ext)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить файл")
		return
	}

	if err := os.MkdirAll(h.dir, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось создать каталог загрузок")
		return
	}

	target, err := os.Create(filepath.Join(h.dir, name))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить файл")
		return
	}
	defer target.Close()

	written, err := io.Copy(target, io.LimitReader(file, maxUploadSize))
	if err != nil {
		_ = os.Remove(filepath.Join(h.dir, name))
		writeError(w, http.StatusInternalServerError, "Не удалось записать файл")
		return
	}

	asset, err := h.assets.Create(r.Context(), repository.AssetInput{
		Filename:   name,
		Original:   filepath.Base(header.Filename),
		Mime:       mime,
		SizeBytes:  written,
		UploadedBy: middleware.UserID(r.Context()),
	})
	if err != nil {
		_ = os.Remove(filepath.Join(h.dir, name))
		writeError(w, http.StatusInternalServerError, "Не удалось сохранить сведения о файле")
		return
	}

	writeJSON(w, http.StatusCreated, asset)
}

func (h *UploadHandler) list(w http.ResponseWriter, r *http.Request) {
	assets, err := h.assets.List(r.Context(), queryInt(r, "limit", 60, 1, 200))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Не удалось загрузить список файлов")
		return
	}
	writeJSON(w, http.StatusOK, assets)
}

func (h *UploadHandler) remove(w http.ResponseWriter, r *http.Request) {
	asset, err := h.assets.Delete(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, statusForRepoError(err), "Не удалось удалить файл")
		return
	}

	// Имя файла сгенерировано нами и не содержит разделителей пути.
	_ = os.Remove(filepath.Join(h.dir, asset.Filename))
	writeJSON(w, http.StatusOK, map[string]string{"message": "Файл удалён"})
}

func randomName(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", hex.EncodeToString(buf), ext), nil
}
