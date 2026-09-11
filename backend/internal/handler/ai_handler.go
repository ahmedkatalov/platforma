package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"platforma/backend/internal/ai"
	"platforma/backend/internal/middleware"
	"platforma/backend/internal/repository"

	"github.com/go-chi/chi/v5"
)

const (
	maxContextChars = 4000
	maxMessageChars = 2000
	maxMessages     = 20
)

// AIHandler — ИИ-помощник по урокам (Google Gemini). Ключ живёт на сервере,
// студенту отдаём только текст ответа. Есть лимиты на частоту запросов.
type AIHandler struct {
	client    *ai.Client
	settings  *repository.AISettingsRepo
	envKey    string // fallback-ключ из окружения (GEMINI_API_KEY)
	envModel  string // fallback-модель из окружения (GEMINI_MODEL)
	perMinute *ai.RateLimiter
	perHour   *ai.RateLimiter
}

func NewAIHandler(client *ai.Client, settings *repository.AISettingsRepo, envKey, envModel string) *AIHandler {
	return &AIHandler{
		client:    client,
		settings:  settings,
		envKey:    strings.TrimSpace(envKey),
		envModel:  strings.TrimSpace(envModel),
		perMinute: ai.NewRateLimiter(6, time.Minute),
		perHour:   ai.NewRateLimiter(60, time.Hour),
	}
}

func (h *AIHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/status", h.status)
	r.Post("/ask", h.ask)
	return r
}

// resolve читает настройки один раз и возвращает эффективные ключ/модель
// (приоритет у заданных в админке, иначе — из окружения) и признак доступности:
// клиент есть, фича включена и ключ задан. При ошибке чтения — fail-closed,
// чтобы status и ask не могли разойтись из-за двух отдельных запросов к БД.
func (h *AIHandler) resolve(ctx context.Context) (key, model string, ok bool) {
	key, model = h.envKey, h.envModel
	if h.client == nil {
		return key, model, false
	}
	s, err := h.settings.Get(ctx)
	if err != nil {
		return key, model, false
	}
	if s.APIKey != "" {
		key = s.APIKey
	}
	if s.Model != "" {
		model = s.Model
	}
	return key, model, s.Enabled && key != ""
}

// TestConnection делает пробный запрос к Gemini текущими ключом/моделью и
// возвращает реальную ошибку (для диагностики в админке). Ключ игнорирует
// флаг «включено»: проверять связь можно и при выключенном помощнике.
func (h *AIHandler) TestConnection(ctx context.Context) (model string, err error) {
	key, model := h.envKey, h.envModel
	if h.client == nil {
		return model, errors.New("клиент ИИ не инициализирован")
	}
	if s, e := h.settings.Get(ctx); e == nil {
		if s.APIKey != "" {
			key = s.APIKey
		}
		if s.Model != "" {
			model = s.Model
		}
	}
	if key == "" {
		return model, errors.New("ключ Gemini не задан")
	}
	_, err = h.client.Ask(ctx, key, model,
		"Ты — проверка связи. Ответь ровно одним словом.",
		[]ai.Message{{Role: "user", Text: "Ответь одним словом: OK"}})
	return model, err
}

func (h *AIHandler) status(w http.ResponseWriter, r *http.Request) {
	_, _, ok := h.resolve(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{"enabled": ok})
}

type askMessage struct {
	Role string `json:"role"` // user | assistant
	Text string `json:"text"`
}

func (h *AIHandler) ask(w http.ResponseWriter, r *http.Request) {
	key, model, ok := h.resolve(r.Context())
	if !ok {
		writeError(w, http.StatusServiceUnavailable, "ИИ-помощник сейчас недоступен")
		return
	}

	userID := middleware.UserID(r.Context())
	if !h.perHour.Allow(userID) || !h.perMinute.Allow(userID) {
		writeError(w, http.StatusTooManyRequests, "Слишком много вопросов к ИИ подряд. Подождите немного.")
		return
	}

	var body struct {
		LessonID string       `json:"lessonId"`
		Context  string       `json:"context"`
		Messages []askMessage `json:"messages"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Собираем историю: режем по количеству и длине, роли приводим к модели.
	msgs := body.Messages
	if len(msgs) > maxMessages {
		msgs = msgs[len(msgs)-maxMessages:]
	}
	history := make([]ai.Message, 0, len(msgs))
	for _, m := range msgs {
		text := strings.TrimSpace(m.Text)
		if text == "" {
			continue
		}
		history = append(history, ai.Message{Role: mapRole(m.Role), Text: clip(text, maxMessageChars)})
	}
	if len(history) == 0 || history[len(history)-1].Role != "user" {
		writeError(w, http.StatusBadRequest, "Нет вопроса для ИИ")
		return
	}

	answer, err := h.client.Ask(r.Context(), key, model, systemPrompt(body.Context), history)
	if err != nil {
		if errors.Is(err, ai.ErrBlocked) {
			writeError(w, http.StatusUnprocessableEntity,
				"ИИ не смог ответить на этот запрос. Переформулируйте вопрос.")
			return
		}
		// Детали (в т.ч. ответ Gemini) — только в лог сервера, не студенту.
		log.Printf("ai ask: %v", err)
		writeError(w, http.StatusBadGateway, "ИИ временно недоступен, попробуйте позже")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"answer": answer})
}

func mapRole(role string) string {
	if role == "assistant" || role == "model" {
		return "model"
	}
	return "user"
}

func clip(s string, n int) string {
	if len([]rune(s)) <= n {
		return s
	}
	return string([]rune(s)[:n])
}

// systemPrompt — задаёт роль наставника и вставляет выделенный фрагмент урока.
func systemPrompt(context string) string {
	var b strings.Builder
	b.WriteString(`Ты — дружелюбный наставник по DevOps на учебной платформе Okvion Learning.
Помогаешь студенту понять материал урока. Отвечай на русском языке, кратко и по делу,
с примерами команд и кода, где это уместно (в блоках ` + "```" + `).
Если вопрос выходит за рамки темы урока — мягко верни студента к теме.
Не выполняй за студента проверяемые задания целиком: подсказывай направление и объясняй,
но не давай готовый ответ для зачёта. Не выдумывай факты: если чего-то не знаешь, так и скажи.`)

	if ctx := clip(strings.TrimSpace(context), maxContextChars); ctx != "" {
		b.WriteString("\n\nСтудент читает урок и выделил этот фрагмент — опирайся на него:\n\"\"\"\n")
		b.WriteString(ctx)
		b.WriteString("\n\"\"\"")
	}
	return b.String()
}
