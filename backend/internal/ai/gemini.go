// Package ai — тонкий клиент к Google Gemini (AI Studio) для ИИ-помощника.
// Ключ живёт только здесь, на сервере; наружу отдаём лишь текст ответа.
package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const defaultModel = "gemini-2.5-flash"

// Message — одна реплика диалога. Role: "user" или "model".
type Message struct {
	Role string
	Text string
}

// Client вызывает generateContent у Gemini. Ключ, модель и прокси передаются в
// каждый вызов Ask (их источник — БД с fallback на окружение), а не хранятся
// в клиенте. http-клиенты кэшируются по адресу прокси.
type Client struct {
	mu      sync.Mutex
	clients map[string]*http.Client // ключ — адрес прокси ("" = напрямую)
}

func NewClient() *Client {
	return &Client{clients: make(map[string]*http.Client)}
}

// httpClient — http.Client для заданного прокси (кэшируется). Пустой proxy —
// прямое соединение. Поддерживает http/https/socks5.
func (c *Client) httpClient(proxy string) (*http.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cl, ok := c.clients[proxy]; ok {
		return cl, nil
	}
	tr := &http.Transport{}
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil || u.Host == "" {
			return nil, errors.New("ai: неверный адрес прокси")
		}
		tr.Proxy = http.ProxyURL(u)
	}
	// Таймаут меньше WriteTimeout сервера (30с), чтобы успеть ответить.
	cl := &http.Client{Timeout: 25 * time.Second, Transport: tr}
	c.clients[proxy] = cl
	return cl, nil
}

// --- форма запроса/ответа Gemini ---

type genPart struct {
	Text string `json:"text"`
}
type genContent struct {
	Role  string    `json:"role,omitempty"`
	Parts []genPart `json:"parts"`
}
type genRequest struct {
	SystemInstruction *genContent  `json:"systemInstruction,omitempty"`
	Contents          []genContent `json:"contents"`
	GenerationConfig  genConfig    `json:"generationConfig"`
}
type genConfig struct {
	Temperature     float64         `json:"temperature"`
	MaxOutputTokens int             `json:"maxOutputTokens"`
	ThinkingConfig  *thinkingConfig `json:"thinkingConfig,omitempty"`
}

// thinkingConfig управляет «размышлениями» моделей Gemini 2.5. Они тратят токены
// из общего лимита ответа, и при небольшом лимите видимый ответ обрывается
// (например, на пустом блоке кода). Для 2.5-flash отключаем их (budget 0).
type thinkingConfig struct {
	ThinkingBudget int `json:"thinkingBudget"`
}
type genResponse struct {
	Candidates []struct {
		Content      genContent `json:"content"`
		FinishReason string     `json:"finishReason"`
	} `json:"candidates"`
	PromptFeedback struct {
		BlockReason string `json:"blockReason"`
	} `json:"promptFeedback"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// ErrBlocked — ответ отфильтрован политиками безопасности Gemini.
var ErrBlocked = errors.New("ai: ответ заблокирован фильтрами безопасности")

// Ask отправляет системную инструкцию и историю диалога, возвращает текст ответа.
// apiKey, model и proxy берутся из настроек (БД или окружения) на каждый запрос.
func (c *Client) Ask(ctx context.Context, apiKey, model, proxy, system string, history []Message) (string, error) {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", errors.New("ai: не задан ключ Gemini")
	}
	if strings.TrimSpace(model) == "" {
		model = defaultModel
	}

	httpClient, err := c.httpClient(strings.TrimSpace(proxy))
	if err != nil {
		return "", err
	}

	contents := make([]genContent, 0, len(history))
	for _, m := range history {
		role := m.Role
		if role != "model" {
			role = "user"
		}
		contents = append(contents, genContent{Role: role, Parts: []genPart{{Text: m.Text}}})
	}

	cfg := genConfig{
		Temperature:     0.4,
		MaxOutputTokens: 2048,
	}
	// 2.5-flash по умолчанию «думает», и размышления съедают лимит ответа —
	// из-за этого ответ обрывался. Отключаем размышления для flash-моделей 2.5.
	if strings.Contains(model, "2.5-flash") {
		cfg.ThinkingConfig = &thinkingConfig{ThinkingBudget: 0}
	}
	reqBody := genRequest{
		Contents:         contents,
		GenerationConfig: cfg,
	}
	if strings.TrimSpace(system) != "" {
		reqBody.SystemInstruction = &genContent{Parts: []genPart{{Text: system}}}
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// Ключ — в заголовке, а не в URL: так он не попадёт в логи прокси/доступа.
	req.Header.Set("x-goog-api-key", apiKey)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}

	var parsed genResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("ai: не разобрать ответ (%d)", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		// Сообщение Gemini полезно для логов, но ключа в нём нет.
		msg := parsed.Error.Message
		if msg == "" {
			msg = strings.TrimSpace(string(raw))
		}
		return "", fmt.Errorf("ai: gemini вернул %d: %s", resp.StatusCode, msg)
	}

	if parsed.PromptFeedback.BlockReason != "" {
		return "", ErrBlocked
	}

	for _, cand := range parsed.Candidates {
		var sb strings.Builder
		for _, p := range cand.Content.Parts {
			sb.WriteString(p.Text)
		}
		if text := strings.TrimSpace(sb.String()); text != "" {
			return text, nil
		}
		if cand.FinishReason == "SAFETY" {
			return "", ErrBlocked
		}
	}
	return "", errors.New("ai: пустой ответ модели")
}
