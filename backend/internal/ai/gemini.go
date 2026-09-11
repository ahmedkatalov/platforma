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
	"strings"
	"time"
)

const defaultModel = "gemini-2.5-flash"

// Message — одна реплика диалога. Role: "user" или "model".
type Message struct {
	Role string
	Text string
}

// Client вызывает generateContent у Gemini.
type Client struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewClient(apiKey, model string) *Client {
	if strings.TrimSpace(model) == "" {
		model = defaultModel
	}
	return &Client{
		apiKey: strings.TrimSpace(apiKey),
		model:  model,
		// Таймаут меньше WriteTimeout сервера (30с), чтобы успеть ответить.
		http: &http.Client{Timeout: 25 * time.Second},
	}
}

func (c *Client) Model() string { return c.model }

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
	Temperature     float64 `json:"temperature"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
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
func (c *Client) Ask(ctx context.Context, system string, history []Message) (string, error) {
	contents := make([]genContent, 0, len(history))
	for _, m := range history {
		role := m.Role
		if role != "model" {
			role = "user"
		}
		contents = append(contents, genContent{Role: role, Parts: []genPart{{Text: m.Text}}})
	}

	reqBody := genRequest{
		Contents: contents,
		GenerationConfig: genConfig{
			Temperature:     0.4,
			MaxOutputTokens: 1024,
		},
	}
	if strings.TrimSpace(system) != "" {
		reqBody.SystemInstruction = &genContent{Parts: []genPart{{Text: system}}}
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent", c.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	// Ключ — в заголовке, а не в URL: так он не попадёт в логи прокси/доступа.
	req.Header.Set("x-goog-api-key", c.apiKey)

	resp, err := c.http.Do(req)
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
